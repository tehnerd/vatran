package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// VIPHandler handles VIP management operations.
type VIPHandler struct {
	manager *lb.Manager
}

// NewVIPHandler creates a new VIPHandler.
//
// Returns a new VIPHandler instance.
func NewVIPHandler() *VIPHandler {
	return &VIPHandler{
		manager: lb.GetManager(),
	}
}

// HandleVIPs handles VIP operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *VIPHandler) HandleVIPs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetAllVIPs(w, r)
	case http.MethodPost:
		h.handleAddVIP(w, r)
	case http.MethodDelete:
		h.handleDelVIP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetAllVIPs handles GET /vips - retrieves all VIPs.
func (h *VIPHandler) handleGetAllVIPs(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	vips, err := lbInstance.GetAllVIPs()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	response := make([]models.VIPResponse, len(vips))
	for i, vip := range vips {
		response[i] = models.VIPResponse{
			Address: vip.Address,
			Port:    vip.Port,
			Proto:   vip.Proto,
		}
	}

	models.WriteSuccess(w, response)
}

// handleAddVIP handles POST /vips - adds a new VIP.
func (h *VIPHandler) handleAddVIP(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.AddVIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.AddVIP(vip, req.Flags); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelVIP handles DELETE /vips - deletes a VIP.
func (h *VIPHandler) handleDelVIP(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.VIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.DelVIP(vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleVIPFlags handles VIP flags operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *VIPHandler) HandleVIPFlags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetVIPFlags(w, r)
	case http.MethodPut:
		h.handleModifyVIPFlags(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetVIPFlags handles GET /vips/flags - gets VIP flags.
func (h *VIPHandler) handleGetVIPFlags(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.VIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	flags, err := lbInstance.GetVIPFlags(vip)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.VIPFlagsResponse{Flags: flags})
}

// handleModifyVIPFlags handles PUT /vips/flags - modifies VIP flags.
func (h *VIPHandler) handleModifyVIPFlags(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ModifyVIPFlagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.ModifyVIP(vip, req.Flag, req.Set); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleHashFunction handles PUT /vips/hash-function - changes hash function.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *VIPHandler) HandleHashFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ChangeHashFunctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.ChangeHashFunctionForVIP(vip, katran.HashFunction(req.HashFunction)); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
