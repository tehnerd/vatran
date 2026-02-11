# Healthcheck Service REST API

This document describes the REST API for the standalone healthcheck service. The katran server communicates with this service to manage health checks for VIP backends.

## Interaction Model

The katran server is the primary consumer of this API:

1. **Startup:** The server calls `POST /api/v1/targets` for each VIP that has a healthcheck configuration, registering the VIP, its reals, and the healthcheck parameters. The server then starts a background poller.
2. **Polling:** The server periodically calls `GET /api/v1/health` to fetch health states for all registered VIPs and applies state transitions via its existing health update logic.
3. **Real add/remove:** When reals are added to or removed from a VIP, the server calls `POST /api/v1/targets/reals` or `DELETE /api/v1/targets/reals` to keep the healthcheck service in sync.
4. **VIP deletion:** When a VIP is deleted, the server calls `DELETE /api/v1/targets` to stop health checks and clean up.
5. **HC service down:** Reals stay in their last known state. The server retries registrations on reconnect and uses `GET /api/v1/health` for a full resync.

## Base URL

All API endpoints are prefixed with `/api/v1` unless otherwise noted.

## Response Format

All responses follow a standard format consistent with the katran server API:

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
| `INVALID_REQUEST` | 400 | Malformed or invalid request |
| `NOT_FOUND` | 404 | VIP or real not found |
| `ALREADY_EXISTS` | 409 | VIP already registered |
| `INTERNAL_ERROR` | 500 | Internal error |

---

## Service Health

### GET /health

Check if the healthcheck service is running.

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

## Target Management

Targets represent VIPs registered with the healthcheck service, along with their reals and healthcheck configuration.

### POST /api/v1/targets

Register a VIP with its reals and healthcheck configuration. The service begins health checking immediately after registration.

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
      "weight": 100,
      "flags": 0
    },
    {
      "address": "192.168.1.2",
      "weight": 100,
      "flags": 0
    }
  ],
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
- `INVALID_REQUEST` — Missing or invalid fields (e.g., unknown healthcheck type, `timeout_ms` >= `interval_ms`)
- `ALREADY_EXISTS` — VIP is already registered. Use `PUT /api/v1/targets` to update

### PUT /api/v1/targets

Update the healthcheck configuration for an already-registered VIP. Optionally replaces the reals list.

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
    "interval_ms": 10000,
    "timeout_ms": 3000,
    "healthy_threshold": 2,
    "unhealthy_threshold": 5
  },
  "reals": [
    {
      "address": "192.168.1.1",
      "weight": 100,
      "flags": 0
    },
    {
      "address": "192.168.1.3",
      "weight": 100,
      "flags": 0
    }
  ]
}
```

The `reals` field is optional. When provided, it replaces the entire reals list for the VIP. When omitted, the existing reals list is preserved.

**Response:**
```json
{
  "success": true
}
```

**Errors:**
- `NOT_FOUND` — VIP is not registered
- `INVALID_REQUEST` — Invalid healthcheck config

### DELETE /api/v1/targets

Deregister a VIP and stop all health checks for its reals.

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
- `NOT_FOUND` — VIP is not registered

### GET /api/v1/targets

List all registered VIPs with their healthcheck configurations and real counts.

**Response:**
```json
{
  "success": true,
  "data": [
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
      },
      "real_count": 2
    }
  ]
}
```

---

## Real Management

Add or remove individual reals from a registered VIP without replacing the entire list.

### POST /api/v1/targets/reals

Add reals to a registered VIP. Reals that already exist are skipped.

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
      "address": "192.168.1.3",
      "weight": 100,
      "flags": 0
    },
    {
      "address": "192.168.1.4",
      "weight": 100,
      "flags": 0
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "added": 2,
    "skipped": 0
  }
}
```

**Errors:**
- `NOT_FOUND` — VIP is not registered
- `INVALID_REQUEST` — Empty reals list or invalid real addresses

### DELETE /api/v1/targets/reals

Remove reals from a registered VIP. Reals that are not found are counted but not treated as errors.

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
      "address": "192.168.1.3"
    },
    {
      "address": "192.168.1.4"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "removed": 1,
    "not_found": 1
  }
}
```

**Errors:**
- `NOT_FOUND` — VIP is not registered
- `INVALID_REQUEST` — Empty reals list

---

## Health State Queries

### GET /api/v1/health/vip

Get the health state for all reals of a single VIP.

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
- `NOT_FOUND` — VIP is not registered

### GET /api/v1/health

Get the health state for all registered VIPs and their reals. This is the primary endpoint used by the katran server's background poller for bulk health state synchronization.

**Response:**
```json
{
  "success": true,
  "data": [
    {
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
    },
    {
      "vip": {
        "address": "10.0.0.2",
        "port": 443,
        "proto": 6
      },
      "reals": [
        {
          "address": "192.168.2.1",
          "healthy": true,
          "last_check_time": "2025-01-15T10:30:00Z",
          "last_status_change": "2025-01-15T08:00:00Z",
          "consecutive_failures": 0
        }
      ]
    }
  ]
}
```

---

## Healthcheck Config Object Reference

The healthcheck config object is used when registering or updating targets.

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
| `interval_ms` | int | 5000 | Check interval in milliseconds |
| `timeout_ms` | int | 2000 | Timeout in milliseconds. Must be less than `interval_ms` |
| `healthy_threshold` | int | 3 | Consecutive successful checks before marking a real as healthy |
| `unhealthy_threshold` | int | 3 | Consecutive failed checks before marking a real as unhealthy |

For `dummy` type, `interval_ms`, `timeout_ms`, `healthy_threshold`, and `unhealthy_threshold` are ignored.
