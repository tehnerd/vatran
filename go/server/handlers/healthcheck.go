package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// HealthcheckHandler handles healthcheck management operations.
type HealthcheckHandler struct {
	manager *lb.Manager
}

// NewHealthcheckHandler creates a new HealthcheckHandler.
//
// Returns a new HealthcheckHandler instance.
func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{
		manager: lb.GetManager(),
	}
}

// HandleHealthcheckerDsts handles healthcheck destination operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *HealthcheckHandler) HandleHealthcheckerDsts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetHealthcheckerDsts(w, r)
	case http.MethodPost:
		h.handleAddHealthcheckerDst(w, r)
	case http.MethodDelete:
		h.handleDelHealthcheckerDst(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetHealthcheckerDsts handles GET /healthcheck/dsts - gets all healthcheck destinations.
func (h *HealthcheckHandler) handleGetHealthcheckerDsts(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	dsts, err := lbInstance.GetHealthcheckersDst()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	response := make([]models.HealthcheckerDstResponse, len(dsts))
	for i, dst := range dsts {
		response[i] = models.HealthcheckerDstResponse{
			Somark: dst.Somark,
			Dst:    dst.Dst,
		}
	}

	models.WriteSuccess(w, response)
}

// handleAddHealthcheckerDst handles POST /healthcheck/dsts - adds a healthcheck destination.
func (h *HealthcheckHandler) handleAddHealthcheckerDst(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.HealthcheckerDstRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.AddHealthcheckerDst(req.Somark, req.Dst); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelHealthcheckerDst handles DELETE /healthcheck/dsts - deletes a healthcheck destination.
func (h *HealthcheckHandler) handleDelHealthcheckerDst(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.DelHealthcheckerDstRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.DelHealthcheckerDst(req.Somark); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleHCKeys handles healthcheck key operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *HealthcheckHandler) HandleHCKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleAddHCKey(w, r)
	case http.MethodDelete:
		h.handleDelHCKey(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleAddHCKey handles POST /healthcheck/keys - adds a healthcheck key.
func (h *HealthcheckHandler) handleAddHCKey(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.HCKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	hcKey := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.AddHCKey(hcKey); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelHCKey handles DELETE /healthcheck/keys - deletes a healthcheck key.
func (h *HealthcheckHandler) handleDelHCKey(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.HCKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	hcKey := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	if err := lbInstance.DelHCKey(hcKey); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
