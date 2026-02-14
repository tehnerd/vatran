package bgpservice

import (
	"encoding/json"
	"net/http"
)

// apiResponse is the standard BGP service response wrapper.
type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *apiError   `json:"error,omitempty"`
}

// apiError is the error payload in a BGP service response.
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// advertiseRequest is the request body for POST /api/v1/routes/advertise.
type advertiseRequest struct {
	VIP         string   `json:"vip"`
	PrefixLen   uint8    `json:"prefix_len"`
	Communities []string `json:"communities"`
	LocalPref   uint32   `json:"local_pref"`
}

// withdrawRequest is the request body for POST /api/v1/routes/withdraw.
type withdrawRequest struct {
	VIP       string `json:"vip"`
	PrefixLen uint8  `json:"prefix_len"`
}

// advertiseResponse is the response data for a route advertise operation.
type advertiseResponse struct {
	VIP       string `json:"vip"`
	PrefixLen uint8  `json:"prefix_len"`
	Advertised bool  `json:"advertised"`
	WasNew    bool   `json:"was_new"`
}

// withdrawResponse is the response data for a route withdraw operation.
type withdrawResponse struct {
	VIP           string `json:"vip"`
	PrefixLen     uint8  `json:"prefix_len"`
	Advertised    bool   `json:"advertised"`
	WasAdvertised bool   `json:"was_advertised"`
}

// addPeerRequest is the request body for POST /api/v1/peers.
type addPeerRequest struct {
	Address   string `json:"address"`
	ASN       uint32 `json:"asn"`
	HoldTime  int    `json:"hold_time"`
	Keepalive int    `json:"keepalive"`
}

// removePeerRequest is the request body for DELETE /api/v1/peers.
type removePeerRequest struct {
	Address string `json:"address"`
}

// Handlers implements the HTTP handlers for the BGP service API.
type Handlers struct {
	state  *State
	bgp    *BGPSpeaker
	config *BGPConfig
}

// NewHandlers creates a new Handlers instance.
//
// Parameters:
//   - state: The shared route state store.
//   - bgp: The BGP speaker for route announcements.
//   - config: The BGP configuration for default values.
//
// Returns a new Handlers instance.
func NewHandlers(state *State, bgp *BGPSpeaker, config *BGPConfig) *Handlers {
	return &Handlers{
		state:  state,
		bgp:    bgp,
		config: config,
	}
}

// HandleAdvertise handles POST /api/v1/routes/advertise.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleAdvertise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req advertiseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	if req.VIP == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "vip is required")
		return
	}
	if req.PrefixLen == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "prefix_len is required")
		return
	}

	// Apply defaults from config
	communities := req.Communities
	if communities == nil {
		communities = h.config.Communities
	}
	localPref := req.LocalPref
	if localPref == 0 {
		localPref = h.config.LocalPref
	}

	// Announce via BGP
	if err := h.bgp.AnnounceRoute(req.VIP, req.PrefixLen, communities, localPref); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Update state
	wasNew := h.state.Advertise(req.VIP, req.PrefixLen, communities, localPref)

	writeSuccessResponse(w, advertiseResponse{
		VIP:        req.VIP,
		PrefixLen:  req.PrefixLen,
		Advertised: true,
		WasNew:     wasNew,
	})
}

// HandleWithdraw handles POST /api/v1/routes/withdraw.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req withdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	if req.VIP == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "vip is required")
		return
	}
	if req.PrefixLen == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "prefix_len is required")
		return
	}

	// Withdraw from BGP
	if err := h.bgp.WithdrawRoute(req.VIP, req.PrefixLen); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Update state
	wasAdvertised := h.state.Withdraw(req.VIP, req.PrefixLen)

	writeSuccessResponse(w, withdrawResponse{
		VIP:           req.VIP,
		PrefixLen:     req.PrefixLen,
		Advertised:    false,
		WasAdvertised: wasAdvertised,
	})
}

// HandleRoutes handles GET /api/v1/routes.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	routes := h.state.GetAllRoutes()
	writeSuccessResponse(w, routes)
}

// HandleRouteVIP handles GET /api/v1/routes/vip?address=...
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleRouteVIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "missing 'address' query parameter")
		return
	}

	route, found := h.state.GetRoute(address)
	if !found {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "route not found for VIP: "+address)
		return
	}

	writeSuccessResponse(w, route)
}

// HandlePeers dispatches /api/v1/peers requests by HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandlePeers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListPeers(w, r)
	case http.MethodPost:
		h.handleAddPeer(w, r)
	case http.MethodDelete:
		h.handleRemovePeer(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleServiceHealth handles GET /health (liveness check).
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleServiceHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeSuccessResponse(w, map[string]string{"status": "ok"})
}

// handleListPeers handles GET /api/v1/peers.
func (h *Handlers) handleListPeers(w http.ResponseWriter, r *http.Request) {
	peers, err := h.bgp.ListPeers()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeSuccessResponse(w, peers)
}

// handleAddPeer handles POST /api/v1/peers.
func (h *Handlers) handleAddPeer(w http.ResponseWriter, r *http.Request) {
	var req addPeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	if req.Address == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "address is required")
		return
	}
	if req.ASN == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "asn is required")
		return
	}

	peerCfg := PeerConfig{
		Address:   req.Address,
		ASN:       req.ASN,
		HoldTime:  req.HoldTime,
		Keepalive: req.Keepalive,
	}

	if err := h.bgp.AddPeer(peerCfg); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeSuccessResponse(w, nil)
}

// handleRemovePeer handles DELETE /api/v1/peers.
func (h *Handlers) handleRemovePeer(w http.ResponseWriter, r *http.Request) {
	var req removePeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	if req.Address == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "address is required")
		return
	}

	if err := h.bgp.RemovePeer(req.Address); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeSuccessResponse(w, nil)
}

// writeSuccessResponse writes a JSON success response.
func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse{Success: true, Data: data})
}

// writeErrorResponse writes a JSON error response.
func writeErrorResponse(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiResponse{
		Success: false,
		Error:   &apiError{Code: code, Message: message},
	})
}
