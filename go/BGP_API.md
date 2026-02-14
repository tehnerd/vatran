# BGP Service API Documentation

Base URL: `http://<host>:<port>` (default port: 9100)

## VIP Route Management

### Advertise a VIP

Announce a VIP prefix via BGP to all configured peers.

**Request:** `POST /api/v1/routes/advertise`

```json
{
  "vip": "10.0.0.1",
  "prefix_len": 32,
  "communities": ["65000:100"],
  "local_pref": 100
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `vip` | string | yes | VIP IP address to advertise |
| `prefix_len` | int | yes | Prefix length (e.g., 32 for /32) |
| `communities` | []string | no | BGP communities (defaults from config) |
| `local_pref` | int | no | Local preference (defaults from config) |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "vip": "10.0.0.1",
    "prefix_len": 32,
    "advertised": true,
    "was_new": true
  }
}
```

### Withdraw a VIP

Withdraw a VIP prefix from BGP announcements.

**Request:** `POST /api/v1/routes/withdraw`

```json
{
  "vip": "10.0.0.1",
  "prefix_len": 32
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `vip` | string | yes | VIP IP address to withdraw |
| `prefix_len` | int | yes | Prefix length |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "vip": "10.0.0.1",
    "prefix_len": 32,
    "advertised": false,
    "was_advertised": true
  }
}
```

### List All Routes

Get the status of all tracked VIP routes.

**Request:** `GET /api/v1/routes`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "vip": "10.0.0.1",
      "prefix_len": 32,
      "advertised": true,
      "since": "2026-02-13T10:00:00Z",
      "communities": ["65000:100"],
      "local_pref": 100
    }
  ]
}
```

### Get Route Status

Get the status of a specific VIP route.

**Request:** `GET /api/v1/routes/vip?address=10.0.0.1`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `address` | string | yes | VIP IP address |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "vip": "10.0.0.1",
    "prefix_len": 32,
    "advertised": true,
    "since": "2026-02-13T10:00:00Z",
    "communities": ["65000:100"],
    "local_pref": 100
  }
}
```

## BGP Peer Management

### List Peers

Get all BGP peers and their session status.

**Request:** `GET /api/v1/peers`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "address": "10.0.0.254",
      "asn": 65001,
      "state": "established",
      "uptime": "2h30m",
      "prefixes_announced": 5
    }
  ]
}
```

### Add Peer

Add a new BGP peer.

**Request:** `POST /api/v1/peers`

```json
{
  "address": "10.0.0.254",
  "asn": 65001,
  "hold_time": 90,
  "keepalive": 30
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `address` | string | yes | Peer IP address |
| `asn` | int | yes | Peer AS number |
| `hold_time` | int | no | Hold timer in seconds (default: 90) |
| `keepalive` | int | no | Keepalive interval in seconds (default: 30) |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": null
}
```

### Remove Peer

Remove a BGP peer.

**Request:** `DELETE /api/v1/peers`

```json
{
  "address": "10.0.0.254"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `address` | string | yes | Peer IP address to remove |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": null
}
```

## Service Health

### Liveness Check

**Request:** `GET /health`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error description"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Malformed request body or missing required fields |
| `NOT_FOUND` | 404 | Route or peer not found |
| `ALREADY_EXISTS` | 409 | Route or peer already exists |
| `INTERNAL_ERROR` | 500 | Internal server or BGP speaker error |

## Configuration

The BGP service is configured via YAML:

```yaml
server:
  port: 9100
  read_timeout: 30
  write_timeout: 30

bgp:
  asn: 65000
  router_id: "10.0.0.1"
  listen_port: 179
  local_pref: 100
  communities:
    - "65000:100"
  peers:
    - address: "10.0.0.254"
      asn: 65001
      hold_time: 90
      keepalive: 30
```
