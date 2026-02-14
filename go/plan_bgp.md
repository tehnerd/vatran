# BGP Routing Service Plan

## Context

The katran load balancer needs BGP integration to advertise healthy VIPs to the network. When a VIP has enough healthy backends (reals), it should be announced via BGP so routers direct traffic to it. When health drops below threshold, the VIP should be withdrawn. This is a standalone microservice following the same architecture as the existing healthcheck service (`go/hcservice/`).

## Deliverables

1. **`go/BGP_API.md`** - API documentation
2. **`go/bgpservice/`** - BGP service package (standalone microservice)
3. **`go/cmd/bgp-service/main.go`** - Entry point
4. **`go/server/lb/bgpclient.go`** - Client in katran server to push updates to BGP service
5. **Integration in `go/server/lb/poller.go`** - Trigger advertise/withdraw after health transitions

---

## Part 1: BGP Service API Design (`go/BGP_API.md`)

Base path: `/api/v1`

### VIP Route Management (katran pushes to these)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/routes/advertise` | Advertise a VIP via BGP |
| `POST` | `/api/v1/routes/withdraw` | Withdraw a VIP from BGP |
| `GET` | `/api/v1/routes` | List all tracked VIPs and their advertise status |
| `GET` | `/api/v1/routes/vip?address=...` | Get status of a specific VIP |

### BGP Peer Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/peers` | List all BGP peers and their session status |
| `POST` | `/api/v1/peers` | Add a BGP peer |
| `DELETE` | `/api/v1/peers` | Remove a BGP peer |

### Service Health

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check |

### Request/Response Formats

**Advertise request** (`POST /api/v1/routes/advertise`):
```json
{
  "vip": "10.0.0.1",
  "prefix_len": 32,
  "communities": ["65000:100"],
  "local_pref": 100
}
```

**Withdraw request** (`POST /api/v1/routes/withdraw`):
```json
{
  "vip": "10.0.0.1",
  "prefix_len": 32
}
```

**Add peer** (`POST /api/v1/peers`):
```json
{
  "address": "10.0.0.254",
  "asn": 65001,
  "hold_time": 90,
  "keepalive": 30
}
```

**All responses** follow the existing pattern:
```json
{
  "success": true,
  "data": { ... },
  "error": { "code": "...", "message": "..." }
}
```

**Route status response** (from `GET /api/v1/routes`):
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

**Peer status response** (from `GET /api/v1/peers`):
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

---

## Part 2: BGP Service Package (`go/bgpservice/`)

### File Structure

```
go/bgpservice/
  config.go     - YAML config, defaults, validation
  server.go     - Server struct, New(), Start(), Stop(), RunWithGracefulShutdown()
  routes.go     - RegisterRoutes()
  handlers.go   - HTTP handlers
  state.go      - Thread-safe VIP advertisement tracking
  bgp.go        - GoBGP speaker wrapper (advertise/withdraw/peer management)
```

### config.go

```go
type Config struct {
    Server  ServerConfig  `yaml:"server"`
    BGP     BGPConfig     `yaml:"bgp"`
}

type ServerConfig struct {
    Host         string `yaml:"host"`
    Port         int    `yaml:"port"`          // default: 9100
    ReadTimeout  int    `yaml:"read_timeout"`  // default: 30
    WriteTimeout int    `yaml:"write_timeout"` // default: 30
}

type BGPConfig struct {
    ASN         uint32      `yaml:"asn"`          // local AS number (required)
    RouterID    string      `yaml:"router_id"`    // BGP router ID (required)
    ListenPort  int         `yaml:"listen_port"`  // default: 179
    LocalPref   uint32      `yaml:"local_pref"`   // default local pref for announcements
    Communities []string    `yaml:"communities"`  // default communities
    Peers       []PeerConfig `yaml:"peers"`       // initial peers
}

type PeerConfig struct {
    Address   string `yaml:"address"`    // peer IP (required)
    ASN       uint32 `yaml:"asn"`        // peer AS (required)
    HoldTime  int    `yaml:"hold_time"`  // default: 90
    Keepalive int    `yaml:"keepalive"`  // default: 30
}
```

Sample YAML:
```yaml
server:
  port: 9100

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
```

Pattern: follows `go/hcservice/config.go` — `DefaultConfig()`, `LoadConfig(path)`, `Validate()`, `Addr()`.

### server.go

Follows `go/hcservice/server.go` exactly:
```go
type Server struct {
    config   *Config
    state    *State
    bgp      *BGPSpeaker
    handlers *Handlers
    httpSrv  *http.Server
}
```

- `New(config)` — creates State, BGPSpeaker, Handlers, wires HTTP mux
- `Start()` — starts BGP speaker, then HTTP server (blocking)
- `Stop(ctx)` — stops HTTP server, then BGP speaker
- `RunWithGracefulShutdown()` — signal handling (same pattern as hcservice)

### state.go

Thread-safe tracking of advertised VIPs:
```go
type RouteState struct {
    VIP        string    // VIP address
    PrefixLen  uint8     // prefix length (typically 32 or 128)
    Advertised bool      // currently announced?
    Since      time.Time // when last state change happened
    Communities []string
    LocalPref  uint32
}

type State struct {
    mu     sync.RWMutex
    routes map[string]*RouteState  // key: "vip/prefixlen" e.g. "10.0.0.1/32"
}
```

Methods:
- `Advertise(vip, prefixLen, communities, localPref) (bool, error)` — returns true if newly advertised
- `Withdraw(vip, prefixLen) (bool, error)` — returns true if was advertised
- `GetRoute(vip) (*RouteState, bool)` — get single route status
- `GetAllRoutes() []RouteState` — list all tracked routes

### bgp.go — GoBGP Integration

Wrapper around the GoBGP server library (`github.com/osrg/gobgp/v3`):

```go
type BGPSpeaker struct {
    server    *gobgpserver.BgpServer
    config    *BGPConfig
}
```

Methods:
- `NewBGPSpeaker(config) *BGPSpeaker`
- `Start() error` — initialize GoBGP server, set global config, add initial peers
- `Stop()` — stop GoBGP server
- `AnnounceRoute(vip string, prefixLen uint8, communities []string, localPref uint32) error`
- `WithdrawRoute(vip string, prefixLen uint8) error`
- `AddPeer(cfg PeerConfig) error`
- `RemovePeer(address string) error`
- `ListPeers() ([]PeerStatus, error)` — query GoBGP for peer states

Uses GoBGP's `api.AddPathRequest` / `api.DeletePathRequest` for route manipulation and `api.AddPeerRequest` / `api.DeletePeerRequest` for peer management.

### handlers.go

Follows `go/hcservice/handlers.go` pattern — same `apiResponse`/`apiError` wrappers, same `writeSuccessResponse`/`writeErrorResponse` helpers.

```go
type Handlers struct {
    state *State
    bgp   *BGPSpeaker
    config *BGPConfig  // for defaults (local_pref, communities)
}
```

Handler methods dispatch by HTTP method like hcservice.

### routes.go

```go
func RegisterRoutes(mux *http.ServeMux, handlers *Handlers) {
    mux.HandleFunc("/api/v1/routes/advertise", handlers.HandleAdvertise)
    mux.HandleFunc("/api/v1/routes/withdraw", handlers.HandleWithdraw)
    mux.HandleFunc("/api/v1/routes", handlers.HandleRoutes)
    mux.HandleFunc("/api/v1/routes/vip", handlers.HandleRouteVIP)
    mux.HandleFunc("/api/v1/peers", handlers.HandlePeers)
    mux.HandleFunc("/health", handlers.HandleServiceHealth)
}
```

---

## Part 3: Entry Point (`go/cmd/bgp-service/main.go`)

Follows `go/cmd/hc-service/main.go` pattern:
```go
func main() {
    configFile := flag.String("config", "", "Path to YAML config file")
    flag.Parse()
    cfg, err := bgpservice.LoadConfig(*configFile)
    // ...
    srv := bgpservice.New(cfg)
    srv.RunWithGracefulShutdown()
}
```

---

## Part 4: Katran Server Integration

### New file: `go/server/lb/bgpclient.go`

HTTP client for pushing advertise/withdraw to the BGP service. Follows `go/server/lb/hcclient.go` pattern exactly:

```go
type BGPClient struct {
    baseURL    string
    httpClient *http.Client
}
```

Methods:
- `NewBGPClient(baseURL string) *BGPClient`
- `Advertise(ctx, vip string, prefixLen uint8) error` — `POST /api/v1/routes/advertise`
- `Withdraw(ctx, vip string, prefixLen uint8) error` — `POST /api/v1/routes/withdraw`

Uses same `doRequest` / response-parsing pattern as `HCClient`.

### Modify: `go/server/lb/manager.go`

Add BGP client to Manager (similar to how HC client is managed):
- New field: `bgpClient *BGPClient`
- New method: `SetBGPEndpoint(url string)` — called during init if configured
- New method: `GetBGPClient() *BGPClient`

### Modify: `go/server/lb/poller.go`

After processing health transitions for a VIP, evaluate whether to advertise or withdraw:

```go
// After processing all real health changes for a VIP:
healthyCount := state.CountHealthyReals(vipKey)
bgpClient := p.manager.GetBGPClient()
if bgpClient == nil {
    continue
}

// threshold from config (e.g., stored in state or manager)
if healthyCount >= threshold && !isAdvertised(vipKey) {
    bgpClient.Advertise(ctx, vipAddress, 32)
}
if healthyCount < threshold && isAdvertised(vipKey) {
    bgpClient.Withdraw(ctx, vipAddress, 32)
}
```

### Modify: `go/server/lb/state.go`

Add helper method:
- `CountHealthyReals(vipKey string) int` — counts reals where `Healthy == true`
- `GetVIPAddress(vipKey string) string` — extracts address from key

### Modify: `go/server/config.go`

Add BGP endpoint to features config:
```go
type FeaturesConfig struct {
    // ... existing fields ...
    BGPEndpoint string `yaml:"bgp_endpoint"` // e.g., "http://localhost:9100"
    BGPMinHealthyReals int `yaml:"bgp_min_healthy_reals"` // default: 1
}
```

### Modify: `go/server/server.go`

In `InitFromConfig()`, if `bgp_endpoint` is configured:
- Call `manager.SetBGPEndpoint(endpoint)`
- Store min healthy reals threshold for poller use

---

## Part 5: Data Flow

```
Health check detects real goes down
    ↓
HC Poller polls HC service, gets updated health
    ↓
Poller updates state store (UpdateHealth)
    ↓
Poller applies katran transition (add/remove real)
    ↓
Poller counts remaining healthy reals for VIP
    ↓
If count < threshold AND VIP is advertised:
    → BGPClient.Withdraw(vip) → BGP Service → GoBGP withdraws from peers
    ↓
If count >= threshold AND VIP is NOT advertised:
    → BGPClient.Advertise(vip) → BGP Service → GoBGP announces to peers
```

---

## Implementation Order

1. `go/BGP_API.md` — API documentation
2. `go/bgpservice/config.go` — config structs, defaults, validation, YAML loading
3. `go/bgpservice/state.go` — thread-safe route state tracking
4. `go/bgpservice/bgp.go` — GoBGP speaker wrapper
5. `go/bgpservice/handlers.go` — HTTP API handlers
6. `go/bgpservice/routes.go` — route registration
7. `go/bgpservice/server.go` — server wiring and lifecycle
8. `go/cmd/bgp-service/main.go` — entry point
9. `go/server/lb/bgpclient.go` — BGP client in katran server
10. `go/server/lb/state.go` — add `CountHealthyReals` helper
11. `go/server/lb/manager.go` — add BGP client management
12. `go/server/lb/poller.go` — add advertise/withdraw logic after health transitions
13. `go/server/config.go` — add `bgp_endpoint` and `bgp_min_healthy_reals` to config
14. `go/server/server.go` — wire BGP client during `InitFromConfig()`
15. Update `go.mod` — add `github.com/osrg/gobgp/v3` dependency

## Files Modified (existing)

- `go/server/lb/manager.go` — add BGP client field + getter/setter
- `go/server/lb/state.go` — add `CountHealthyReals()` method
- `go/server/lb/poller.go` — add post-transition BGP advertise/withdraw logic
- `go/server/config.go` — add BGP fields to FeaturesConfig
- `go/server/server.go` — wire BGP client in InitFromConfig
- `go/go.mod` / `go/go.sum` — new GoBGP dependency

## Files Created (new)

- `go/BGP_API.md`
- `go/bgpservice/config.go`
- `go/bgpservice/state.go`
- `go/bgpservice/bgp.go`
- `go/bgpservice/handlers.go`
- `go/bgpservice/routes.go`
- `go/bgpservice/server.go`
- `go/cmd/bgp-service/main.go`
- `go/server/lb/bgpclient.go`

## Verification

1. **Build**: `cd go && go build ./...` — all packages compile
2. **Unit tests**: Write tests for `bgpservice/state.go` and `bgpservice/handlers.go`
3. **Config validation**: Test `LoadConfig` with valid and invalid YAML
4. **Integration smoke test**: Start bgp-service with config, verify `/health` returns OK, verify `GET /api/v1/peers` returns configured peers
5. **Existing tests**: `cd go && go test ./server/...` — ensure no regressions
