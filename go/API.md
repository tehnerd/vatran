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
  "hash_func": 0
}
```

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

Get all real servers for a VIP.

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
  "data": [
    {
      "address": "192.168.1.1",
      "weight": 100,
      "flags": 0
    }
  ]
}
```

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

**Request Body:**
```json
{
  "address": "192.168.1.1"
}
```

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

---

## Statistics

### GET /api/v1/stats/vip

Get VIP statistics.

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
    "v1": 1000000,
    "v2": 500000000
  }
}
```

### GET /api/v1/stats/vip/decap

Get VIP decapsulation statistics.

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
    "v1": 100,
    "v2": 50000
  }
}
```

### GET /api/v1/stats/real

Get real server statistics.

**Request Body:**
```json
{
  "index": 1
}
```

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

### GET /api/v1/stats/lru/global

Get global LRU statistics.

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

### GET /api/v1/stats/decap

Get decapsulation statistics.

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

## Running the Server

```bash
# Basic usage
./katran-server -port 8080

# With TLS
./katran-server -port 443 -tls-cert /path/to/cert.pem -tls-key /path/to/key.pem

# With CORS enabled
./katran-server -port 8080 -cors -cors-origins "http://localhost:3000"

# Full options
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
  -idle-timeout 120
```

### Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
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
