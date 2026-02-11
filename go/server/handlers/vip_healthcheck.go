package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
	"github.com/tehnerd/vatran/go/server/types"
)

// VIPHealthcheckHandler handles per-VIP healthcheck configuration operations.
type VIPHealthcheckHandler struct {
	manager *lb.Manager
}

// NewVIPHealthcheckHandler creates a new VIPHealthcheckHandler.
//
// Returns a new VIPHealthcheckHandler instance.
func NewVIPHealthcheckHandler() *VIPHealthcheckHandler {
	return &VIPHealthcheckHandler{
		manager: lb.GetManager(),
	}
}

// HandleVIPHealthcheck dispatches to the correct handler based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *VIPHealthcheckHandler) HandleVIPHealthcheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.handleSetVIPHealthcheck(w, r)
	case http.MethodGet:
		h.handleGetVIPHealthcheck(w, r)
	case http.MethodDelete:
		h.handleDeleteVIPHealthcheck(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleVIPHealthcheckStatus handles GET /vips/healthcheck/status.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *VIPHealthcheckHandler) HandleVIPHealthcheckStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.handleGetVIPHealthcheckStatus(w, r)
}

// handleSetVIPHealthcheck handles PUT /vips/healthcheck - set or update HC config for a VIP.
func (h *VIPHealthcheckHandler) handleSetVIPHealthcheck(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.SetVIPHealthcheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	// Verify VIP exists
	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}
	if _, err := lbInstance.GetVIPFlags(vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	// Apply defaults and validate
	req.Healthcheck.ApplyDefaults()
	if err := req.Healthcheck.Validate(); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid healthcheck config: "+err.Error()))
		return
	}

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("state store not initialized"))
		return
	}

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)

	if req.Healthcheck.Type != "dummy" {
		// Non-dummy: require HC client
		hcClient := h.manager.GetHCClient()
		if hcClient == nil {
			models.WriteError(w, http.StatusNotImplemented,
				models.NewFeatureDisabledError("healthchecker_endpoint is not configured"))
			return
		}

		// Check if updating existing config
		_, hadExisting := state.GetHCConfig(vipKey)

		hcVIP := types.HCVIPKey{
			Address: req.VIP.Address,
			Port:    req.VIP.Port,
			Proto:   req.VIP.Proto,
		}

		reals := state.GetReals(vipKey)

		ctx := context.Background()
		if hadExisting {
			if err := hcClient.UpdateVIP(ctx, hcVIP, &req.Healthcheck, reals); err != nil {
				models.WriteError(w, http.StatusBadGateway,
					models.NewHCServiceUnavailableError("failed to update HC service: "+err.Error()))
				return
			}
		} else {
			if err := hcClient.RegisterVIP(ctx, hcVIP, reals, &req.Healthcheck); err != nil {
				models.WriteError(w, http.StatusBadGateway,
					models.NewHCServiceUnavailableError("failed to register with HC service: "+err.Error()))
				return
			}
		}
	} else {
		// Dummy type: mark all reals healthy immediately
		reals := state.GetReals(vipKey)
		for _, rs := range reals {
			oldHealthy, found := state.UpdateHealth(vipKey, rs.Address, true)
			if found && !oldHealthy {
				real := katran.Real{
					Address: rs.Address,
					Weight:  rs.Weight,
					Flags:   rs.Flags,
				}
				if err := lbInstance.AddRealForVIP(real, vip); err != nil {
					log.Printf("HC dummy: failed to add real %s to katran for VIP %s: %v", rs.Address, vipKey, err)
				}
			}
		}
	}

	// Store config
	hcCfg := req.Healthcheck
	state.SetHCConfig(vipKey, &hcCfg)

	models.WriteSuccess(w, nil)
}

// handleGetVIPHealthcheck handles GET /vips/healthcheck - get HC config for a VIP.
func (h *VIPHealthcheckHandler) handleGetVIPHealthcheck(w http.ResponseWriter, r *http.Request) {
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

	// Verify VIP exists
	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}
	if _, err := lbInstance.GetVIPFlags(vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteSuccess(w, nil)
		return
	}

	vipKey := lb.VIPKeyString(req.Address, req.Port, req.Proto)
	hcCfg, found := state.GetHCConfig(vipKey)
	if !found {
		models.WriteSuccess(w, nil)
		return
	}

	models.WriteSuccess(w, hcCfg)
}

// handleDeleteVIPHealthcheck handles DELETE /vips/healthcheck - remove HC config from a VIP.
func (h *VIPHealthcheckHandler) handleDeleteVIPHealthcheck(w http.ResponseWriter, r *http.Request) {
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

	// Verify VIP exists
	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}
	if _, err := lbInstance.GetVIPFlags(vip); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("state store not initialized"))
		return
	}

	vipKey := lb.VIPKeyString(req.Address, req.Port, req.Proto)
	oldCfg, found := state.DelHCConfig(vipKey)
	if !found {
		models.WriteError(w, http.StatusNotFound,
			models.NewNotFoundError("no healthcheck configuration found for this VIP"))
		return
	}

	// For non-dummy: deregister from HC service (fire-and-forget)
	if oldCfg.Type != "dummy" {
		if hcClient := h.manager.GetHCClient(); hcClient != nil {
			hcVIP := types.HCVIPKey{
				Address: req.Address,
				Port:    req.Port,
				Proto:   req.Proto,
			}
			if err := hcClient.DeregisterVIP(context.Background(), hcVIP); err != nil {
				log.Printf("HC deregister: failed to deregister VIP %s from HC service: %v", vipKey, err)
			}
		}
	}

	models.WriteSuccess(w, nil)
}

// handleGetVIPHealthcheckStatus handles GET /vips/healthcheck/status - detailed health status.
func (h *VIPHealthcheckHandler) handleGetVIPHealthcheckStatus(w http.ResponseWriter, r *http.Request) {
	_, ok := h.manager.Get()
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

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("state store not initialized"))
		return
	}

	vipKey := lb.VIPKeyString(req.Address, req.Port, req.Proto)
	hcCfg, found := state.GetHCConfig(vipKey)
	if !found {
		models.WriteError(w, http.StatusNotFound,
			models.NewNotFoundError("no healthcheck configuration found for this VIP"))
		return
	}

	if hcCfg.Type == "dummy" {
		// Build response from local state (all healthy)
		reals := state.GetReals(vipKey)
		realStatuses := make([]models.RealHealthStatusResponse, len(reals))
		for i, rs := range reals {
			realStatuses[i] = models.RealHealthStatusResponse{
				Address: rs.Address,
				Healthy: rs.Healthy,
			}
		}
		models.WriteSuccess(w, models.VIPHealthStatusResponse{
			VIP: models.VIPResponse{
				Address: req.Address,
				Port:    req.Port,
				Proto:   req.Proto,
			},
			Reals: realStatuses,
		})
		return
	}

	// Non-dummy: proxy to HC service
	hcClient := h.manager.GetHCClient()
	if hcClient == nil {
		models.WriteError(w, http.StatusNotImplemented,
			models.NewFeatureDisabledError("healthchecker_endpoint is not configured"))
		return
	}

	hcVIP := types.HCVIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	hcResp, err := hcClient.GetVIPHealthStatus(context.Background(), hcVIP)
	if err != nil {
		models.WriteError(w, http.StatusBadGateway,
			models.NewHCServiceUnavailableError("failed to query HC service: "+err.Error()))
		return
	}

	// Convert HC service response to our response format
	realStatuses := make([]models.RealHealthStatusResponse, len(hcResp.Reals))
	for i, rh := range hcResp.Reals {
		realStatuses[i] = models.RealHealthStatusResponse{
			Address:             rh.Address,
			Healthy:             rh.Healthy,
			LastCheckTime:       rh.LastCheckTime,
			LastStatusChange:    rh.LastStatusChange,
			ConsecutiveFailures: rh.ConsecutiveFailures,
		}
	}

	models.WriteSuccess(w, models.VIPHealthStatusResponse{
		VIP: models.VIPResponse{
			Address: req.Address,
			Port:    req.Port,
			Proto:   req.Proto,
		},
		Reals: realStatuses,
	})
}
