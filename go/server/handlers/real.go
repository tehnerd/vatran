package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// RealHandler handles real server management operations.
type RealHandler struct {
	manager *lb.Manager
}

// NewRealHandler creates a new RealHandler.
//
// Returns a new RealHandler instance.
func NewRealHandler() *RealHandler {
	return &RealHandler{
		manager: lb.GetManager(),
	}
}

// HandleVIPReals handles real server operations for a VIP based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleVIPReals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetReals(w, r)
	case http.MethodPost:
		h.handleAddReal(w, r)
	case http.MethodDelete:
		h.handleDelReal(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetReals handles GET /vips/reals - gets all reals for a VIP.
func (h *RealHandler) handleGetReals(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.VIPRequest
	if err := decodeVIPRequest(r, &req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	reals, err := lbInstance.GetRealsForVIP(vip)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	response := make([]models.RealResponse, len(reals))
	for i, real := range reals {
		response[i] = models.RealResponse{
			Address: real.Address,
			Weight:  real.Weight,
			Flags:   real.Flags,
		}
	}

	models.WriteSuccess(w, response)
}

// handleAddReal handles POST /vips/reals - adds a real to a VIP.
func (h *RealHandler) handleAddReal(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.AddRealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	real := katran.Real{
		Address: req.Real.Address,
		Weight:  req.Real.Weight,
		Flags:   req.Real.Flags,
	}

	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	if err := lbInstance.AddRealForVIP(real, vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelReal handles DELETE /vips/reals - removes a real from a VIP.
func (h *RealHandler) handleDelReal(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.DelRealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	real := katran.Real{
		Address: req.Real.Address,
		Weight:  req.Real.Weight,
		Flags:   req.Real.Flags,
	}

	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	if err := lbInstance.DelRealForVIP(real, vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleBatchModifyReals handles PUT /vips/reals/batch - batch modify reals.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleBatchModifyReals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ModifyRealsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	reals := make([]katran.Real, len(req.Reals))
	for i, r := range req.Reals {
		reals[i] = katran.Real{
			Address: r.Address,
			Weight:  r.Weight,
			Flags:   r.Flags,
		}
	}

	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	if err := lbInstance.ModifyRealsForVIP(katran.ModifyAction(req.Action), reals, vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleRealIndex handles GET /reals/index - gets the index for a real.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleRealIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.GetRealIndexRequest
	if err := decodeRealIndexRequest(r, &req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request: "+err.Error()))
		return
	}

	index, err := lbInstance.GetIndexForReal(req.Address)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.RealIndexResponse{Index: index})
}

// HandleRealFlags handles PUT /reals/flags - modifies real flags.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleRealFlags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ModifyRealFlagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.ModifyReal(req.Address, req.Flags, req.Set); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
