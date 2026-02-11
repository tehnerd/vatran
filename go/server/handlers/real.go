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
// Queries the state store for health info. Falls back to katran (all healthy) if VIP not in state.
func (h *RealHandler) handleGetReals(w http.ResponseWriter, r *http.Request) {
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

	vipKey := lb.VIPKeyString(req.Address, req.Port, req.Proto)

	// Try state store first
	state, stateOK := h.manager.GetState()
	if stateOK {
		stateReals := state.GetReals(vipKey)
		if stateReals != nil {
			response := make([]models.RealResponse, len(stateReals))
			for i, rs := range stateReals {
				response[i] = models.RealResponse{
					Address: rs.Address,
					Weight:  rs.Weight,
					Flags:   rs.Flags,
					Healthy: rs.Healthy,
				}
			}
			models.WriteSuccess(w, response)
			return
		}
	}

	// Fall back to katran (VIP not tracked in state store)
	lbInstance, _ := h.manager.Get()
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
			Healthy: true,
		}
	}

	models.WriteSuccess(w, response)
}

// handleAddReal handles POST /vips/reals - adds a real to a VIP.
// Adds to state store first. Only calls katran AddRealForVIP if the real is healthy.
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

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)

	// Add to state store
	state, stateOK := h.manager.GetState()
	var rs *lb.RealState
	if stateOK {
		rs = state.AddReal(vipKey, req.Real.Address, req.Real.Weight, req.Real.Flags)
	}

	// Only add to katran if healthy (or no state store)
	if !stateOK || (rs != nil && rs.Healthy) {
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
			// Rollback state store on katran error
			if stateOK {
				state.DelReal(vipKey, req.Real.Address)
			}
			models.WriteKatranError(w, err)
			return
		}
	}

	// Notify HC service if VIP has non-dummy HC config (fire-and-forget)
	if stateOK {
		if hcCfg, hasHC := state.GetHCConfig(vipKey); hasHC && hcCfg.Type != "dummy" {
			if hcClient := h.manager.GetHCClient(); hcClient != nil {
				hcVIP := types.HCVIPKey{Address: req.VIP.Address, Port: req.VIP.Port, Proto: req.VIP.Proto}
				newReals := []lb.RealState{{Address: req.Real.Address, Weight: req.Real.Weight, Flags: req.Real.Flags}}
				if err := hcClient.AddReals(context.Background(), hcVIP, newReals); err != nil {
					log.Printf("HC notify: failed to add real %s to HC service for VIP %s: %v", req.Real.Address, vipKey, err)
				}
			}
		}
	}

	models.WriteCreated(w, nil)
}

// handleDelReal handles DELETE /vips/reals - removes a real from a VIP.
// Removes from state store. Only calls katran DelRealForVIP if the real was healthy.
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

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)

	// Remove from state store
	state, stateOK := h.manager.GetState()
	var wasHealthy bool
	if stateOK {
		oldState, found := state.DelReal(vipKey, req.Real.Address)
		if found {
			wasHealthy = oldState.Healthy
		} else {
			// Not in state store, assume it was healthy (in katran)
			wasHealthy = true
		}
	} else {
		wasHealthy = true
	}

	// Only delete from katran if it was healthy (was in katran)
	if wasHealthy {
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
	}

	// Notify HC service if VIP has non-dummy HC config (fire-and-forget)
	if stateOK {
		if hcCfg, hasHC := state.GetHCConfig(vipKey); hasHC && hcCfg.Type != "dummy" {
			if hcClient := h.manager.GetHCClient(); hcClient != nil {
				hcVIP := types.HCVIPKey{Address: req.VIP.Address, Port: req.VIP.Port, Proto: req.VIP.Proto}
				if err := hcClient.RemoveReals(context.Background(), hcVIP, []string{req.Real.Address}); err != nil {
					log.Printf("HC notify: failed to remove real %s from HC service for VIP %s: %v", req.Real.Address, vipKey, err)
				}
			}
		}
	}

	models.WriteSuccess(w, nil)
}

// HandleBatchModifyReals handles PUT /vips/reals/batch - batch modify reals.
// For adds: adds all to state, filters healthy ones for katran.
// For deletes: removes all from state, filters previously-healthy ones for katran.
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

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)
	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	state, stateOK := h.manager.GetState()

	if !stateOK {
		// No state store, pass through directly to katran
		reals := make([]katran.Real, len(req.Reals))
		for i, r := range req.Reals {
			reals[i] = katran.Real{
				Address: r.Address,
				Weight:  r.Weight,
				Flags:   r.Flags,
			}
		}
		if err := lbInstance.ModifyRealsForVIP(katran.ModifyAction(req.Action), reals, vip); err != nil {
			models.WriteKatranError(w, err)
			return
		}
		models.WriteSuccess(w, nil)
		return
	}

	if katran.ModifyAction(req.Action) == katran.ActionAdd {
		// Add action: add all to state, filter healthy for katran
		var healthyReals []katran.Real
		for _, r := range req.Reals {
			rs := state.AddReal(vipKey, r.Address, r.Weight, r.Flags)
			if rs.Healthy {
				healthyReals = append(healthyReals, katran.Real{
					Address: r.Address,
					Weight:  r.Weight,
					Flags:   r.Flags,
				})
			}
		}
		if len(healthyReals) > 0 {
			if err := lbInstance.ModifyRealsForVIP(katran.ActionAdd, healthyReals, vip); err != nil {
				models.WriteKatranError(w, err)
				return
			}
		}
	} else {
		// Delete action: remove all from state, filter previously-healthy for katran
		var healthyReals []katran.Real
		for _, r := range req.Reals {
			oldState, found := state.DelReal(vipKey, r.Address)
			wasHealthy := !found || (found && oldState.Healthy)
			if wasHealthy {
				healthyReals = append(healthyReals, katran.Real{
					Address: r.Address,
					Weight:  r.Weight,
					Flags:   r.Flags,
				})
			}
		}
		if len(healthyReals) > 0 {
			if err := lbInstance.ModifyRealsForVIP(katran.ActionDel, healthyReals, vip); err != nil {
				models.WriteKatranError(w, err)
				return
			}
		}
	}

	models.WriteSuccess(w, nil)
}

// HandleHealthUpdate handles PUT /vips/reals/health - update a single real's health state.
// On state change: healthy->unhealthy calls DelRealForVIP, unhealthy->healthy calls AddRealForVIP.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleHealthUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("state store not initialized"))
		return
	}

	var req models.UpdateRealHealthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)

	oldHealthy, found := state.UpdateHealth(vipKey, req.Address, req.Healthy)
	if !found {
		models.WriteError(w, http.StatusNotFound,
			models.NewNotFoundError("real not found in state store"))
		return
	}

	// No state change, nothing to do in katran
	if oldHealthy == req.Healthy {
		models.WriteSuccess(w, nil)
		return
	}

	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	if req.Healthy {
		// unhealthy -> healthy: add to katran
		// Get current state to get weight/flags
		reals := state.GetReals(vipKey)
		var real katran.Real
		for _, rs := range reals {
			if rs.Address == req.Address {
				real = katran.Real{
					Address: rs.Address,
					Weight:  rs.Weight,
					Flags:   rs.Flags,
				}
				break
			}
		}
		if err := lbInstance.AddRealForVIP(real, vip); err != nil {
			// Rollback state change
			state.UpdateHealth(vipKey, req.Address, oldHealthy)
			models.WriteKatranError(w, err)
			return
		}
	} else {
		// healthy -> unhealthy: remove from katran
		reals := state.GetReals(vipKey)
		var real katran.Real
		for _, rs := range reals {
			if rs.Address == req.Address {
				real = katran.Real{
					Address: rs.Address,
					Weight:  rs.Weight,
					Flags:   rs.Flags,
				}
				break
			}
		}
		if err := lbInstance.DelRealForVIP(real, vip); err != nil {
			// Rollback state change
			state.UpdateHealth(vipKey, req.Address, oldHealthy)
			models.WriteKatranError(w, err)
			return
		}
	}

	models.WriteSuccess(w, nil)
}

// HandleBatchHealthUpdate handles PUT /vips/reals/health/batch - batch update real health states.
// Collects transitions and batch calls ModifyRealsForVIP for adds and deletes.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RealHandler) HandleBatchHealthUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	state, stateOK := h.manager.GetState()
	if !stateOK {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("state store not initialized"))
		return
	}

	var req models.BatchUpdateRealHealthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vipKey := lb.VIPKeyString(req.VIP.Address, req.VIP.Port, req.VIP.Proto)
	vip := katran.VIPKey{
		Address: req.VIP.Address,
		Port:    req.VIP.Port,
		Proto:   req.VIP.Proto,
	}

	// Collect transitions
	type transition struct {
		address    string
		oldHealthy bool
	}
	var toAdd []katran.Real
	var toDel []katran.Real
	var transitions []transition

	// Get current reals for weight/flags lookup
	currentReals := state.GetReals(vipKey)
	realMap := make(map[string]lb.RealState, len(currentReals))
	for _, rs := range currentReals {
		realMap[rs.Address] = rs
	}

	for _, update := range req.Reals {
		oldHealthy, found := state.UpdateHealth(vipKey, update.Address, update.Healthy)
		if !found {
			continue
		}
		if oldHealthy == update.Healthy {
			continue
		}
		transitions = append(transitions, transition{address: update.Address, oldHealthy: oldHealthy})

		rs := realMap[update.Address]
		real := katran.Real{
			Address: rs.Address,
			Weight:  rs.Weight,
			Flags:   rs.Flags,
		}
		if update.Healthy {
			toAdd = append(toAdd, real)
		} else {
			toDel = append(toDel, real)
		}
	}

	// Batch add healthy reals to katran
	if len(toAdd) > 0 {
		if err := lbInstance.ModifyRealsForVIP(katran.ActionAdd, toAdd, vip); err != nil {
			// Rollback all transitions
			for _, t := range transitions {
				state.UpdateHealth(vipKey, t.address, t.oldHealthy)
			}
			models.WriteKatranError(w, err)
			return
		}
	}

	// Batch delete unhealthy reals from katran
	if len(toDel) > 0 {
		if err := lbInstance.ModifyRealsForVIP(katran.ActionDel, toDel, vip); err != nil {
			// Rollback all transitions
			for _, t := range transitions {
				state.UpdateHealth(vipKey, t.address, t.oldHealthy)
			}
			models.WriteKatranError(w, err)
			return
		}
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
