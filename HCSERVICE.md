# Healthcheck Service

## What

The healthcheck service (`hc-service`) is a standalone Go binary that performs active health checks against katran backends. It runs as a sidecar alongside the katran server and exposes a REST API (documented in [HC_API.md](HC_API.md)) that the katran server uses to register VIPs, manage reals, and poll health state.

The service supports four healthcheck types:

| Type | Behavior |
|------|----------|
| `tcp` | TCP connect to VIP:port. Connection success = healthy. |
| `http` | HTTP GET to `http://VIP:port/path`. Match expected status code. |
| `https` | HTTPS GET with configurable TLS settings. Same as HTTP otherwise. |
| `dummy` | Always healthy. No checks are performed. |

## Why

Katran is an XDP/BPF-based L4 load balancer. It forwards packets at kernel level and has no application-layer visibility into backend health. The healthcheck service fills this gap by:

1. **Probing backends through katran itself** -- health checks dial the VIP address with `SO_MARK` set per-real. Katran's BPF program recognizes the mark and routes the check packet to the correct real via tunnel encapsulation. This validates the full forwarding path, not just direct reachability.

2. **Decoupling health checking from the control plane** -- the katran server manages BPF maps and configuration. Health checking is I/O-heavy (many concurrent TCP/HTTP connections with timeouts). Running it in a separate process keeps the control plane responsive and allows independent scaling, restart, and resource limits.

3. **Providing a clean separation of concerns** -- the katran server registers targets, polls health state, and applies transitions (adding/removing reals from BPF maps). The HC service only checks and reports. Neither needs to know the other's internals.

## How

### Architecture

```
                         +-----------------+
                         |  katran server  |
                         |  (port 8080)    |
                         +----+-------+----+
                              |       ^
          register targets,   |       |  poll health state
          manage reals        |       |  GET /api/v1/health
                              v       |
                         +-----------------+
                         |   hc-service    |
                         |  (port 9000)    |
                         +----+-------+----+
                              |       ^
          dial VIP:port       |       |  connection result
          with SO_MARK        |       |
                              v       |
                         +-----------------+
                         |   katran BPF    |
                         |  (XDP/TC)       |
                         +--------+--------+
                                  |
                                  v
                            real backends
```

### SO_MARK Routing

Each unique real server gets a somark value allocated from a configurable range (`base_somark` to `base_somark + max_reals - 1`). The service:

1. Allocates a somark when a real is first seen across any VIP.
2. Registers the somark-to-destination mapping with katran via `POST /api/v1/healthcheck/dsts`.
3. Uses `SO_MARK` on health check sockets so katran's BPF program routes the check to the correct real.
4. Reference-counts somarks -- the same real across multiple VIPs shares one somark. The mapping is deregistered only when the last VIP referencing that real is removed.

### Scheduler

The service uses a tick-and-sweep scheduler with a bounded worker pool:

- A tick loop sweeps all check targets every `tick_interval_ms` (default 100ms).
- When a target's next check time has arrived, it is dispatched to a worker via a buffered channel.
- Workers run the appropriate checker (TCP/HTTP/HTTPS) and update the health state.
- New targets are staggered over `spread_interval_ms` to avoid burst load.
- If all workers are busy, overdue checks are skipped and retried on the next interval.

### Interaction Flow

1. **Startup**: The katran server calls `POST /api/v1/targets` for each VIP that has a healthcheck configuration.
2. **Polling**: The katran server periodically calls `GET /api/v1/health` and applies health transitions to its BPF maps.
3. **Real changes**: `POST /api/v1/targets/reals` and `DELETE /api/v1/targets/reals` keep the HC service in sync when reals are added or removed.
4. **VIP removal**: `DELETE /api/v1/targets` stops all checks and cleans up somark allocations.
5. **HC service down**: Reals stay in their last known state. The katran server retries registrations on reconnect.

### Building

```bash
go build ./go/cmd/hc-service/
```

### Running

```bash
# With defaults (port 9000, katran at localhost:8080)
./hc-service

# With config file
./hc-service -config /etc/hc-service/config.yaml
```

## Configuration

The service is configured via a YAML file. All fields are optional -- sensible defaults are applied.

### Full Example

```yaml
server:
  host: ""              # Bind address ("" = all interfaces)
  port: 9000            # Listen port
  read_timeout: 30      # HTTP read timeout in seconds
  write_timeout: 30     # HTTP write timeout in seconds

katran:
  server_url: "http://localhost:8080"   # Katran server base URL
  timeout: 10                           # HTTP client timeout in seconds

somark:
  base_somark: 10000    # Starting somark value
  max_reals: 4096       # Max unique reals tracked simultaneously

scheduler:
  spread_interval_ms: 3000   # Stagger new checks over this window
  worker_count: 64           # Max concurrent health check goroutines
  tick_interval_ms: 100      # Scheduler sweep granularity
```

### Configuration Reference

#### `server`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `host` | string | `""` | Address to bind to. Empty binds to all interfaces. |
| `port` | int | `9000` | Port to listen on. |
| `read_timeout` | int | `30` | HTTP read timeout in seconds. |
| `write_timeout` | int | `30` | HTTP write timeout in seconds. |

#### `katran`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `server_url` | string | `"http://localhost:8080"` | Base URL of the katran server. **Required** (must be non-empty). |
| `timeout` | int | `10` | HTTP client timeout in seconds for calls to katran. |

#### `somark`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `base_somark` | uint32 | `10000` | Starting somark value. All allocated somarks are in `[base_somark, base_somark + max_reals)`. Must be > 0. |
| `max_reals` | uint32 | `4096` | Maximum number of unique real addresses that can be tracked. Must be > 0. |

#### `scheduler`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spread_interval_ms` | int | `3000` | Window in ms over which newly registered checks are staggered. Must be > 0. |
| `worker_count` | int | `64` | Maximum concurrent health check goroutines. Must be > 0. |
| `tick_interval_ms` | int | `100` | Scheduler sweep interval in ms. Controls how closely checks fire to their schedule. Must be > 0. |

### Minimal Config

Most deployments only need to set the katran server URL:

```yaml
katran:
  server_url: "http://katran:8080"
```

### Production Config

A production deployment with higher concurrency and a larger real pool:

```yaml
server:
  port: 9000

katran:
  server_url: "http://katran:8080"
  timeout: 5

somark:
  base_somark: 20000
  max_reals: 16384

scheduler:
  spread_interval_ms: 5000
  worker_count: 256
  tick_interval_ms: 50
```

## Package Layout

```
go/hcservice/
    config.go          - Config types, YAML parsing, defaults, validation
    somark.go          - Somark allocator with reference counting
    katran_client.go   - HTTP client for katran POST/DELETE /api/v1/healthcheck/dsts
    dialer.go          - Custom dial with SO_MARK via syscall.SetsockoptInt
    checker.go         - Checker interface + HTTP/HTTPS/TCP/Dummy implementations
    state.go           - Thread-safe VIP/real state + health tracking
    scheduler.go       - Tick-and-sweep scheduler with worker pool
    handlers.go        - HTTP handlers implementing HC_API.md endpoints
    routes.go          - Route registration on http.ServeMux
    server.go          - Server lifecycle with graceful shutdown
go/cmd/hc-service/
    main.go            - Entry point, flag parsing, config loading
```

## API Reference

See [HC_API.md](go/HC_API.md) for the full REST API specification.
