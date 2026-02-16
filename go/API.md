# Katran Load Balancer REST API

This document describes the REST API for the Katran Load Balancer HTTP server.

## Base URL

All API endpoints are prefixed with `/api/v1`.

## Response Format

All responses follow a standard format:

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | The request was malformed or invalid |
| `NOT_FOUND` | 404 | The requested resource was not found |
| `ALREADY_EXISTS` | 409 | The resource already exists |
| `SPACE_EXHAUSTED` | 507 | Maximum capacity was reached |
| `BPF_FAILED` | 500 | A BPF operation failed |
| `FEATURE_DISABLED` | 501 | The requested feature is not enabled |
| `INTERNAL_ERROR` | 500 | An internal server error occurred |
| `LB_NOT_INITIALIZED` | 503 | The load balancer is not initialized |
| `UNAUTHORIZED` | 401 | The request lacks valid authentication |
| `HC_SERVICE_UNAVAILABLE` | 502 | The healthcheck service is unreachable |

---

## Health Check

### GET /health

Check if the server is running.

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

---

## Lifecycle Management

### POST /api/v1/lb/create

Create a new load balancer instance.

**Request Body:**
```json
{
  "main_interface": "eth0",
  "balancer_prog_path": "/path/to/balancer.o",
  "healthchecking_prog_path": "/path/to/healthchecker.o",
  "v4_tun_interface": "ipip0",
  "v6_tun_interface": "ipip6tnl0",
  "hc_interface": "eth0",
  "default_mac": "aa:bb:cc:dd:ee:ff",
  "local_mac": "11:22:33:44:55:66",
  "root_map_path": "/sys/fs/bpf/root_map",
  "root_map_pos": 2,
  "use_root_map": true,
  "max_vips": 512,
  "max_reals": 4096,
  "ch_ring_size": 65537,
  "lru_size": 8000000,
  "max_lpm_src_size": 3000000,
  "max_decap_dst": 6,
  "global_lru_size": 100000,
  "enable_hc": true,
  "tunnel_based_hc_encap": true,
  "testing": false,
  "memlock_unlimited": true,
  "flow_debug": false,
  "enable_cid_v3": false,
  "cleanup_on_shutdown": true,
  "forwarding_cores": [0, 1, 2, 3],
  "numa_nodes": [0, 0, 1, 1],
  "xdp_attach_flags": 0,
  "priority": 2307,
  "main_interface_index": 0,
  "hc_interface_index": 0,
  "katran_src_v4": "10.0.0.1",
  "katran_src_v6": "fc00::1",
  "hash_function": "maglev_v2"
}
```
`hash_function` defaults to `"maglev_v2"` when omitted.

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/lb/close

Close the load balancer instance.

**Response:**
```json
{
  "success": true
}
```

### GET /api/v1/lb/status

Get the load balancer status.

**Response:**
```json
{
  "success": true,
  "data": {
    "initialized": true,
    "ready": true
  }
}
```

### POST /api/v1/lb/load-bpf-progs

Load BPF programs into the kernel.

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/lb/attach-bpf-progs

Attach loaded BPF programs to network interfaces.

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/lb/reload

Reload the balancer BPF program at runtime.

**Request Body:**
```json
{
  "path": "/path/to/new_balancer.o",
  "config": null
}
```

**Response:**
```json
{
  "success": true
}
```

---

## VIP Management

### GET /api/v1/vips

List all configured VIPs.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "address": "10.0.0.1",
      "port": 80,
      "proto": 6
    }
  ]
}
```

### POST /api/v1/vips

Add a new VIP.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6,
  "flags": 0
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/vips

Delete a VIP.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true
}
```

### GET /api/v1/vips/flags

Get VIP flags.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "flags": 0
  }
}
```

### PUT /api/v1/vips/flags

Modify VIP flags.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6,
  "flag": 1,
  "set": true
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT /api/v1/vips/hash-function

Change the hash function for a VIP.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6,
  "hash_function": 0
}
```

**Response:**
```json
{
  "success": true
}
```

---

## Real Server Management

### GET /api/v1/vips/reals

Get all real servers for a VIP. Returns both healthy and unhealthy reals with their health state.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "address": "192.168.1.1",
      "weight": 100,
      "flags": 0,
      "healthy": true
    }
  ]
}
```

The `healthy` field indicates whether the real server is currently receiving traffic. When `healthy` is `false`, the real is tracked by the server but is not programmed into katran's forwarding plane.

### POST /api/v1/vips/reals

Add a real server to a VIP.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "real": {
    "address": "192.168.1.1",
    "weight": 100,
    "flags": 0
  }
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/vips/reals

Remove a real server from a VIP.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "real": {
    "address": "192.168.1.1",
    "weight": 0,
    "flags": 0
  }
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT /api/v1/vips/reals/batch

Batch modify real servers for a VIP.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "action": 0,
  "reals": [
    {
      "address": "192.168.1.1",
      "weight": 100,
      "flags": 0
    },
    {
      "address": "192.168.1.2",
      "weight": 100,
      "flags": 0
    }
  ]
}
```

| Action | Description |
|--------|-------------|
| 0 | Add |
| 1 | Delete |

**Response:**
```json
{
  "success": true
}
```

### GET /api/v1/reals/index

Get the internal index for a real server.

**Query Parameters:**
- address (string)

**Response:**
```json
{
  "success": true,
  "data": {
    "index": 1
  }
}
```

### PUT /api/v1/reals/flags

Modify real server flags.

**Request Body:**
```json
{
  "address": "192.168.1.1",
  "flags": 1,
  "set": true
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT /api/v1/vips/reals/health

Update the health state of a single real server for a VIP. When a real transitions from unhealthy to healthy, it is added to katran's forwarding plane. When it transitions from healthy to unhealthy, it is removed.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "address": "192.168.1.1",
  "healthy": true
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT /api/v1/vips/reals/health/batch

Batch update health states for multiple real servers of a VIP. Transitions are batched for efficiency.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "reals": [
    {
      "address": "192.168.1.1",
      "healthy": true
    },
    {
      "address": "192.168.1.2",
      "healthy": false
    }
  ]
}
```

**Response:**
```json
{
  "success": true
}
```

### Health Default Behavior

The default health state of newly added reals depends on the `healthchecker_endpoint` configuration:

- **No `healthchecker_endpoint` configured**: Reals default to **healthy** and immediately receive traffic when added.
- **`healthchecker_endpoint` configured**: Reals default to **unhealthy** and do not receive traffic until explicitly marked healthy via the health update endpoints.

The `healthy` field can be set explicitly per backend in the YAML configuration to override the default behavior.

---

## Statistics

### GET /api/v1/stats/vip

Get VIP statistics.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Counters:**
- v1: packets sent to the VIP
- v2: bytes sent to the VIP

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000000,
    "v2": 500000000
  }
}
```

### GET /api/v1/stats/vip/decap

Get VIP decapsulation statistics.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Counters:**
- v1: packets decapsulated to the VIP
- v2: unused

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 100,
    "v2": 0
  }
}
```

### GET /api/v1/stats/real

Get real server statistics.

**Query Parameters:**
- index (integer)

**Counters:**
- v1: packets sent to the real
- v2: bytes sent to the real

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 500000,
    "v2": 250000000
  }
}
```

### GET /api/v1/stats/lru

Get LRU cache statistics.

**Counters:**
- v1: total packets
- v2: LRU hits

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000000,
    "v2": 900000
  }
}
```

### GET /api/v1/stats/lru/miss

Get LRU miss statistics.

**Counters:**
- v1: TCP SYN misses
- v2: non-SYN misses

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000,
    "v2": 500
  }
}
```

### GET /api/v1/stats/lru/fallback

Get LRU fallback statistics.

**Counters:**
- v1: fallback LRU hits
- v2: unused

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 100,
    "v2": 0
  }
}
```

### GET /api/v1/stats/lru/global

Get global LRU statistics.

**Counters:**
- v1: map lookup failures
- v2: global LRU routed

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 10,
    "v2": 99990
  }
}
```

### GET /api/v1/stats/icmp-too-big

Get ICMP too big statistics.

**Counters:**
- v1: ICMPv4 count
- v2: ICMPv6 count

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 10,
    "v2": 5
  }
}
```

### GET /api/v1/stats/ch-drop

Get consistent hash drop statistics.

**Counters:**
- v1: real ID out of bounds
- v2: real #0 (unmapped)

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 0,
    "v2": 0
  }
}
```

### GET /api/v1/stats/src-routing

Get source routing statistics.

**Counters:**
- v1: packets sent to local backends
- v2: packets sent to remote destinations (LPM matched)

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 100,
    "v2": 50
  }
}
```

### GET /api/v1/stats/inline-decap

Get inline decapsulation statistics.

**Counters:**
- v1: packets decapsulated inline
- v2: unused

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000,
    "v2": 0
  }
}
```

### GET /api/v1/stats/decap

Get decapsulation statistics.

**Counters:**
- v1: IPv4 packets decapsulated
- v2: IPv6 packets decapsulated

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000,
    "v2": 500
  }
}
```

### GET /api/v1/stats/quic-icmp

Get QUIC ICMP statistics.

**Counters:**
- v1: QUIC ICMP messages
- v2: QUIC ICMP messages dropped by Shiv

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 10,
    "v2": 5
  }
}
```

### GET /api/v1/stats/quic-packets

Get QUIC packet routing statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "ch_routed": 1000,
    "cid_initial": 100,
    "cid_invalid_server_id": 5,
    "cid_invalid_server_id_sample": 1,
    "cid_routed": 900,
    "cid_unknown_real_dropped": 0,
    "cid_v0": 100,
    "cid_v1": 200,
    "cid_v2": 300,
    "cid_v3": 300,
    "dst_match_in_lru": 800,
    "dst_mismatch_in_lru": 10,
    "dst_not_found_in_lru": 90
  }
}
```

### GET /api/v1/stats/tcp-server-id-routing

Get TCP server ID routing (TPR) statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "ch_routed": 1000,
    "dst_mismatch_in_lru": 10,
    "sid_routed": 9000,
    "tcp_syn": 100
  }
}
```

### GET /api/v1/stats/xdp/total

Get XDP total statistics.

**Counters:**
- v1: packets
- v2: bytes

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 10000000,
    "v2": 5000000000
  }
}
```

### GET /api/v1/stats/xdp/tx

Get XDP TX statistics.

**Counters:**
- v1: packets (XDP_TX)
- v2: bytes (XDP_TX)

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 9000000,
    "v2": 4500000000
  }
}
```

### GET /api/v1/stats/xdp/drop

Get XDP drop statistics.

**Counters:**
- v1: packets dropped (XDP_DROP)
- v2: bytes dropped (XDP_DROP)

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 1000,
    "v2": 500000
  }
}
```

### GET /api/v1/stats/xdp/pass

Get XDP pass statistics.

**Counters:**
- v1: packets passed to kernel (XDP_PASS)
- v2: bytes passed to kernel (XDP_PASS)

**Response:**
```json
{
  "success": true,
  "data": {
    "v1": 999000,
    "v2": 499500000
  }
}
```

### GET /api/v1/stats/hc-prog

Get healthcheck program statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "packets_processed": 10000,
    "packets_dropped": 10,
    "packets_skipped": 100,
    "packets_too_big": 5
  }
}
```

### GET /api/v1/stats/bpf-map

Get BPF map statistics.

**Query Parameters:**
- map_name (string)

**Response:**
```json
{
  "success": true,
  "data": {
    "max_entries": 512,
    "current_entries": 10
  }
}
```

### GET /api/v1/stats/userspace

Get userspace library statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "bpf_failed_calls": 0,
    "addr_validation_failed": 2
  }
}
```

### GET /api/v1/stats/per-core-packets

Get per-core packet statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "counts": [1000000, 1000000, 1000000, 1000000]
  }
}
```

### GET /api/v1/stats/flood-status

Get flood status.

**Response:**
```json
{
  "success": true,
  "data": {
    "under_flood": false
  }
}
```

### GET /api/v1/stats/monitor

Get monitor statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "limit": 1000,
    "amount": 500,
    "buffer_full": 0
  }
}
```

---

## QUIC Management

### GET /api/v1/quic/reals

Get all QUIC real server mappings.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "address": "192.168.1.1",
      "id": 1
    }
  ]
}
```

### PUT /api/v1/quic/reals

Modify QUIC real server mappings.

**Request Body:**
```json
{
  "action": 0,
  "reals": [
    {
      "address": "192.168.1.1",
      "id": 1
    }
  ]
}
```

**Response:**
```json
{
  "success": true
}
```

---

## Routing

### GET /api/v1/routing/src-rules

Get all source routing rules.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "src": "10.0.0.0/8",
      "dst": "192.168.1.1"
    }
  ]
}
```

### POST /api/v1/routing/src-rules

Add source routing rules.

**Request Body:**
```json
{
  "src_prefixes": ["10.0.0.0/8", "172.16.0.0/12"],
  "dst": "192.168.1.1"
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/routing/src-rules

Delete source routing rules.

**Request Body:**
```json
{
  "src_prefixes": ["10.0.0.0/8"]
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/routing/src-rules/all

Clear all source routing rules.

**Response:**
```json
{
  "success": true
}
```

### GET /api/v1/routing/src-rules/size

Get the number of source routing rules.

**Response:**
```json
{
  "success": true,
  "data": {
    "size": 10
  }
}
```

### GET /api/v1/routing/decap/inline

Get all inline decapsulation destinations.

**Response:**
```json
{
  "success": true,
  "data": ["10.0.0.1", "10.0.0.2"]
}
```

### POST /api/v1/routing/decap/inline

Add an inline decapsulation destination.

**Request Body:**
```json
{
  "dst": "10.0.0.1"
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/routing/decap/inline

Delete an inline decapsulation destination.

**Request Body:**
```json
{
  "dst": "10.0.0.1"
}
```

**Response:**
```json
{
  "success": true
}
```

---

## Healthcheck Management

### GET /api/v1/healthcheck/dsts

Get all healthcheck destination mappings.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "somark": 100,
      "dst": "192.168.1.1"
    }
  ]
}
```

### POST /api/v1/healthcheck/dsts

Add a healthcheck destination mapping.

**Request Body:**
```json
{
  "somark": 100,
  "dst": "192.168.1.1"
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/healthcheck/dsts

Delete a healthcheck destination mapping.

**Request Body:**
```json
{
  "somark": 100
}
```

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/healthcheck/keys

Add a healthcheck key.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE /api/v1/healthcheck/keys

Delete a healthcheck key.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true
}
```

---

## VIP Healthcheck Configuration

These endpoints manage per-VIP healthcheck configuration. When a healthcheck config is set for a VIP, the server registers the VIP, its reals, and the healthcheck config with the external healthcheck service. The server then periodically polls the healthcheck service for real health states.

### Healthcheck Config Object

The healthcheck config object is shared across all endpoints in this section:

```json
{
  "type": "http",
  "port": 8080,
  "http": {
    "path": "/healthz",
    "expected_status": 200,
    "host": "example.com"
  },
  "interval_ms": 5000,
  "timeout_ms": 2000,
  "healthy_threshold": 3,
  "unhealthy_threshold": 3
}
```

**Healthcheck types:**

| Type | Description | Sub-object |
|------|-------------|------------|
| `http` | HTTP GET check | `http`: `path` (required), `expected_status` (default 200), `host` (optional) |
| `https` | HTTPS GET check | `https`: same as `http` + `skip_tls_verify` (default false) |
| `tcp` | TCP connect check | None |
| `dummy` | Always healthy, no checking | None |

**Common fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `port` | int | VIP port | Port to check on each real server |
| `interval_ms` | int | 5000 | Check interval in milliseconds. Ignored for `dummy` |
| `timeout_ms` | int | 2000 | Timeout in milliseconds. Must be less than `interval_ms`. Ignored for `dummy` |
| `healthy_threshold` | int | 3 | Consecutive successful checks before marking a real as healthy. Ignored for `dummy` |
| `unhealthy_threshold` | int | 3 | Consecutive failed checks before marking a real as unhealthy. Ignored for `dummy` |

**Constraints:**
- Non-dummy types require `healthchecker_endpoint` to be configured in the server config, otherwise `FEATURE_DISABLED` is returned.
- The `dummy` type works without a healthcheck service and marks all reals as healthy immediately.
- Setting a healthcheck config registers the VIP and its reals with the healthcheck service. Deleting deregisters them.

### PUT /api/v1/vips/healthcheck

Set or update the healthcheck configuration for a VIP.

**Request Body:**
```json
{
  "vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "healthcheck": {
    "type": "http",
    "http": {
      "path": "/healthz",
      "expected_status": 200
    },
    "interval_ms": 5000,
    "timeout_ms": 2000,
    "healthy_threshold": 3,
    "unhealthy_threshold": 3
  }
}
```

**Response:**
```json
{
  "success": true
}
```

**Errors:**
- `NOT_FOUND` — VIP does not exist
- `INVALID_REQUEST` — Invalid healthcheck config (e.g., `timeout_ms` >= `interval_ms`, missing required sub-object fields)
- `FEATURE_DISABLED` — Non-dummy type requested but `healthchecker_endpoint` is not configured
- `HC_SERVICE_UNAVAILABLE` — Could not reach the healthcheck service to register the VIP

### GET /api/v1/vips/healthcheck

Get the healthcheck configuration for a VIP.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "type": "http",
    "port": 8080,
    "http": {
      "path": "/healthz",
      "expected_status": 200,
      "host": "example.com"
    },
    "interval_ms": 5000,
    "timeout_ms": 2000,
    "healthy_threshold": 3,
    "unhealthy_threshold": 3
  }
}
```

Returns `null` for `data` if no healthcheck is configured for the VIP.

**Errors:**
- `NOT_FOUND` — VIP does not exist

### DELETE /api/v1/vips/healthcheck

Remove the healthcheck configuration from a VIP. This deregisters the VIP from the healthcheck service and stops all health checks for its reals.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true
}
```

**Errors:**
- `NOT_FOUND` — VIP does not exist or has no healthcheck configured
- `HC_SERVICE_UNAVAILABLE` — Could not reach the healthcheck service to deregister (deregistration is still applied locally)

### GET /api/v1/vips/healthcheck/status

Get detailed health status for all reals of a VIP, including timestamps and failure counts from the healthcheck service.

**Query Parameters:**
- address (string)
- port (integer)
- proto (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "vip": {
      "address": "10.0.0.1",
      "port": 80,
      "proto": 6
    },
    "reals": [
      {
        "address": "192.168.1.1",
        "healthy": true,
        "last_check_time": "2025-01-15T10:30:00Z",
        "last_status_change": "2025-01-15T09:00:00Z",
        "consecutive_failures": 0
      },
      {
        "address": "192.168.1.2",
        "healthy": false,
        "last_check_time": "2025-01-15T10:30:00Z",
        "last_status_change": "2025-01-15T10:25:00Z",
        "consecutive_failures": 3
      }
    ]
  }
}
```

**Errors:**
- `NOT_FOUND` — VIP does not exist or has no healthcheck configured
- `HC_SERVICE_UNAVAILABLE` — Could not reach the healthcheck service

### YAML Configuration Extension

Healthcheck configuration can be specified per VIP in the YAML config file using the optional `healthcheck` block:

```yaml
vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
    target_group: web-servers
    healthcheck:
      type: "http"
      port: 8080
      http:
        path: "/healthz"
        expected_status: 200
      interval_ms: 5000
      timeout_ms: 2000
      healthy_threshold: 3
      unhealthy_threshold: 3
```

When loading from config, each VIP with a `healthcheck` block will be registered with the healthcheck service at startup. The `dummy` type can be used in config to mark all reals as healthy without requiring a healthcheck service.

---

## Feature Management

### GET /api/v1/features/check

Check if a feature is available.

**Request Body:**
```json
{
  "feature": 1
}
```

| Feature | Value | Description |
|---------|-------|-------------|
| SrcRouting | 1 | Source-based routing |
| InlineDecap | 2 | Inline packet decapsulation |
| Introspection | 4 | Packet introspection/monitoring |
| GUEEncap | 8 | GUE encapsulation |
| DirectHC | 16 | Direct healthcheck encapsulation |
| LocalDeliveryOpt | 32 | Local delivery optimization |
| FlowDebug | 64 | Flow debugging maps |

**Response:**
```json
{
  "success": true,
  "data": {
    "available": true
  }
}
```

### POST /api/v1/features/install

Install a feature by reloading the BPF program.

**Request Body:**
```json
{
  "feature": 1,
  "prog_path": "/path/to/program_with_feature.o"
}
```

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/features/remove

Remove a feature by reloading the BPF program.

**Request Body:**
```json
{
  "feature": 1,
  "prog_path": "/path/to/program_without_feature.o"
}
```

**Response:**
```json
{
  "success": true
}
```

---

## LRU Management

### DELETE /api/v1/lru

Delete an LRU entry for a specific flow.

**Request Body:**
```json
{
  "dst_vip": {
    "address": "10.0.0.1",
    "port": 80,
    "proto": 6
  },
  "src_ip": "192.168.1.100",
  "src_port": 12345
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "maps": ["lru_map_0", "lru_map_1"]
  }
}
```

### DELETE /api/v1/lru/vip

Purge all LRU entries for a VIP.

**Request Body:**
```json
{
  "address": "10.0.0.1",
  "port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "deleted_count": 1000
  }
}
```

---

## Monitor Control

### POST /api/v1/monitor/stop

Stop the packet monitor.

**Response:**
```json
{
  "success": true
}
```

### POST /api/v1/monitor/restart

Restart the packet monitor.

**Request Body:**
```json
{
  "limit": 1000
}
```

**Response:**
```json
{
  "success": true
}
```

---

## Utility Endpoints

### GET /api/v1/utils/mac

Get the current default router MAC address.

**Response:**
```json
{
  "success": true,
  "data": {
    "mac": "aa:bb:cc:dd:ee:ff"
  }
}
```

### PUT /api/v1/utils/mac

Change the default router MAC address.

**Request Body:**
```json
{
  "mac": "aa:bb:cc:dd:ee:ff"
}
```

**Response:**
```json
{
  "success": true
}
```

### GET /api/v1/utils/real-for-flow

Get the real server for a specific flow.

**Request Body:**
```json
{
  "src": "192.168.1.100",
  "dst": "10.0.0.1",
  "src_port": 12345,
  "dst_port": 80,
  "proto": 6
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "address": "192.168.1.1"
  }
}
```

### POST /api/v1/utils/simulate-packet

Simulate a packet through the BPF program.

**Request Body:**
```json
{
  "packet": "base64_encoded_packet_data"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "packet": "base64_encoded_output_packet"
  }
}
```

### GET /api/v1/utils/prog-fd

Get the katran BPF program file descriptor.

**Response:**
```json
{
  "success": true,
  "data": {
    "fd": 5
  }
}
```

### GET /api/v1/utils/hc-prog-fd

Get the healthchecker BPF program file descriptor.

**Response:**
```json
{
  "success": true,
  "data": {
    "fd": 6
  }
}
```

### GET /api/v1/utils/map-fd

Get a BPF map file descriptor by name.

**Request Body:**
```json
{
  "map_name": "vip_map"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "fd": 7
  }
}
```

### GET /api/v1/utils/global-lru-map-fds

Get global LRU map file descriptors.

**Response:**
```json
{
  "success": true,
  "data": {
    "fds": [8, 9, 10, 11]
  }
}
```

### POST /api/v1/utils/src-ip-encap

Add a source IP for packet encapsulation.

**Request Body:**
```json
{
  "src": "10.0.0.1"
}
```

**Response:**
```json
{
  "success": true
}
```

---

## Configuration Export

### GET /api/v1/config/export

Export the current running configuration as a YAML file. This can be used to save the current state of the load balancer for backup or replication purposes.

**Response:**

Returns a YAML file with Content-Type `application/x-yaml`.

```yaml
# Katran Server Configuration
server:
  host: ""
  port: 8080
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120
  enable_cors: false
  allowed_origins: []
  enable_logging: true
  enable_recovery: true
  static_dir: ""
  bpf_prog_dir: "/path/to/bpf"
  auth:
    enabled: true
    database_path: "/var/lib/katran/auth.db"
    allow_localhost: false
    session_timeout: 24
    bcrypt_cost: 12
    exempt_paths:
      - /metrics

lb:
  interfaces:
    main: "eth0"
    healthcheck: "eth0"
    v4_tunnel: ""
    v6_tunnel: ""
  programs:
    balancer: "/path/to/balancer.bpf.o"
    healthcheck: "/path/to/healthchecking.bpf.o"
  root_map:
    enabled: false
    path: ""
    position: 2
  mac:
    default: "aa:bb:cc:dd:ee:ff"
    local: ""
  capacity:
    max_vips: 512
    max_reals: 4096
    ch_ring_size: 65537
    lru_size: 8000000
    global_lru_size: 100000
    max_lpm_src: 3000000
    max_decap_dst: 6
  cpu:
    forwarding_cores: []
    numa_nodes: []
  xdp:
    attach_flags: 0
    priority: 2307
  encapsulation:
    src_v4: ""
    src_v6: ""
  features:
    enable_healthcheck: true
    tunnel_based_hc_encap: true
    flow_debug: false
    enable_cid_v3: false
    memlock_unlimited: true
    cleanup_on_shutdown: true
    testing: false
    # healthchecker_endpoint: "http://healthchecker.example.com/api"
    # bgp_endpoint: "http://localhost:9100"
    # bgp_min_healthy_reals: 1
  hash_function: "maglev"

target_groups:
  group-0:
    - address: "10.0.0.1"
      weight: 100
      flags: 0
    - address: "10.0.0.2"
      weight: 100

vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
    target_group: group-0
    flags: 0
    healthcheck:
      type: "http"
      port: 8080
      http:
        path: "/healthz"
        expected_status: 200
      interval_ms: 5000
      timeout_ms: 2000
      healthy_threshold: 3
      unhealthy_threshold: 3
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "LB_NOT_INITIALIZED",
    "message": "Load balancer is not initialized"
  }
}
```

### GET /api/v1/config/export/json

Export the current running configuration as JSON. Same data as the YAML export but in JSON format.

**Response:**
```json
{
  "success": true,
  "data": {
    "server": {
      "host": "",
      "port": 8080,
      "read_timeout": 30,
      "write_timeout": 30,
      "idle_timeout": 120,
      "enable_cors": false,
      "allowed_origins": [],
      "enable_logging": true,
      "enable_recovery": true,
      "static_dir": "",
      "bpf_prog_dir": "/path/to/bpf",
      "auth": {
        "enabled": true,
        "database_path": "/var/lib/katran/auth.db",
        "allow_localhost": false,
        "session_timeout": 24,
        "bcrypt_cost": 12,
        "exempt_paths": ["/metrics"]
      }
    },
    "lb": {
      "interfaces": {
        "main": "eth0",
        "healthcheck": "eth0",
        "v4_tunnel": "",
        "v6_tunnel": ""
      },
      "programs": {
        "balancer": "/path/to/balancer.bpf.o",
        "healthcheck": "/path/to/healthchecking.bpf.o"
      },
      "root_map": {
        "enabled": false,
        "path": "",
        "position": 2
      },
      "mac": {
        "default": "aa:bb:cc:dd:ee:ff",
        "local": ""
      },
      "capacity": {
        "max_vips": 512,
        "max_reals": 4096,
        "ch_ring_size": 65537,
        "lru_size": 8000000,
        "global_lru_size": 100000,
        "max_lpm_src": 3000000,
        "max_decap_dst": 6
      },
      "cpu": {
        "forwarding_cores": [],
        "numa_nodes": []
      },
      "xdp": {
        "attach_flags": 0,
        "priority": 2307
      },
      "encapsulation": {
        "src_v4": "",
        "src_v6": ""
      },
      "features": {
        "enable_healthcheck": true,
        "tunnel_based_hc_encap": true,
        "flow_debug": false,
        "enable_cid_v3": false,
        "memlock_unlimited": true,
        "cleanup_on_shutdown": true,
        "testing": false,
        "healthchecker_endpoint": "",
        "bgp_endpoint": "",
        "bgp_min_healthy_reals": 0
      },
      "hash_function": "maglev"
    },
    "target_groups": {
      "group-0": [
        {
          "address": "10.0.0.1",
          "weight": 100,
          "flags": 0
        }
      ]
    },
    "vips": [
      {
        "address": "192.168.1.100",
        "port": 80,
        "proto": "tcp",
        "target_group": "group-0",
        "flags": 0,
        "healthcheck": {
          "type": "http",
          "port": 8080,
          "http": {
            "path": "/healthz",
            "expected_status": 200
          },
          "interval_ms": 5000,
          "timeout_ms": 2000,
          "healthy_threshold": 3,
          "unhealthy_threshold": 3
        }
      }
    ]
  }
}
```

---

## Running the Server

```bash
# Basic usage
./katran-server -port 8080

# With YAML config file (recommended for production)
./katran-server -config /path/to/config.yaml

# With TLS
./katran-server -port 443 -tls-cert /path/to/cert.pem -tls-key /path/to/key.pem

# With CORS enabled
./katran-server -port 8080 -cors -cors-origins "http://localhost:3000"

# Full options (without config file)
./katran-server \
  -host 0.0.0.0 \
  -port 8080 \
  -tls-cert /path/to/cert.pem \
  -tls-key /path/to/key.pem \
  -tls-client-ca /path/to/ca.pem \
  -cors \
  -cors-origins "http://localhost:3000,https://example.com" \
  -read-timeout 30 \
  -write-timeout 30 \
  -idle-timeout 120 \
  -bpf-prog-dir /path/to/bpf
```

### Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | "" | Path to YAML config file (overrides other flags) |
| `-host` | "" | Host to bind to |
| `-port` | 8080 | Port to listen on |
| `-tls-cert` | "" | Path to TLS certificate file |
| `-tls-key` | "" | Path to TLS private key file |
| `-tls-client-ca` | "" | Path to client CA file for mTLS |
| `-cors` | false | Enable CORS |
| `-cors-origins` | "*" | Comma-separated list of allowed CORS origins |
| `-read-timeout` | 30 | Read timeout in seconds |
| `-write-timeout` | 30 | Write timeout in seconds |
| `-idle-timeout` | 120 | Idle timeout in seconds |
| `-no-logging` | false | Disable request logging |
| `-no-recovery` | false | Disable panic recovery |
| `-static-dir` | "" | Path to static files directory for SPA |
| `-bpf-prog-dir` | "" | Base directory for BPF program files |

### Using YAML Configuration

When using the `-config` flag, the server will:
1. Load the configuration from the YAML file
2. Initialize the load balancer with the specified settings
3. Create all VIPs and add backends from target groups
4. Start the HTTP server

Example config file structure:
```yaml
server:
  host: ""
  port: 8080
  bpf_prog_dir: "/path/to/bpf"

lb:
  interfaces:
    main: "eth0"
  programs:
    balancer: "balancer.bpf.o"

target_groups:
  web-servers:
    - address: "10.0.0.1"
      weight: 100

vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
    target_group: web-servers
```

See `config_example.yaml` for a full configuration example with all available options.
