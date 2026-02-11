package hcservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tehnerd/vatran/go/server/types"
)

// apiResponse is the standard HC service response wrapper.
type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *apiError   `json:"error,omitempty"`
}

// apiError is the error payload in an HC service response.
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// registerTargetRequest is the request body for POST /api/v1/targets.
type registerTargetRequest struct {
	VIP         types.HCVIPKey          `json:"vip"`
	Reals       []RealEntry             `json:"reals"`
	Healthcheck types.HealthcheckConfig `json:"healthcheck"`
}

// updateTargetRequest is the request body for PUT /api/v1/targets.
type updateTargetRequest struct {
	VIP         types.HCVIPKey          `json:"vip"`
	Healthcheck types.HealthcheckConfig `json:"healthcheck"`
	Reals       []RealEntry             `json:"reals,omitempty"`
}

// realsRequest is the request body for POST/DELETE /api/v1/targets/reals.
type realsRequest struct {
	VIP   types.HCVIPKey `json:"vip"`
	Reals []RealEntry    `json:"reals"`
}

// realsAddedResponse is the response data for POST /api/v1/targets/reals.
type realsAddedResponse struct {
	Added   int `json:"added"`
	Skipped int `json:"skipped"`
}

// realsRemovedResponse is the response data for DELETE /api/v1/targets/reals.
type realsRemovedResponse struct {
	Removed  int `json:"removed"`
	NotFound int `json:"not_found"`
}

// Handlers implements the HTTP handlers for the HC service API.
type Handlers struct {
	state     *State
	somarks   *SomarkAllocator
	katran    *KatranClient
	scheduler *Scheduler
}

// NewHandlers creates a new Handlers.
//
// Parameters:
//   - state: The shared health state store.
//   - somarks: The somark allocator.
//   - katran: The katran client for somark registration.
//   - scheduler: The check scheduler.
//
// Returns a new Handlers instance.
func NewHandlers(state *State, somarks *SomarkAllocator, katran *KatranClient, scheduler *Scheduler) *Handlers {
	return &Handlers{
		state:     state,
		somarks:   somarks,
		katran:    katran,
		scheduler: scheduler,
	}
}

// HandleTargets dispatches /api/v1/targets requests by HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleTargets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleRegisterTarget(w, r)
	case http.MethodPut:
		h.handleUpdateTarget(w, r)
	case http.MethodDelete:
		h.handleDeleteTarget(w, r)
	case http.MethodGet:
		h.handleListTargets(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleTargetReals dispatches /api/v1/targets/reals requests by HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleTargetReals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleAddReals(w, r)
	case http.MethodDelete:
		h.handleRemoveReals(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleHealthVIP handles GET /api/v1/health/vip.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleHealthVIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key, err := parseVIPQuery(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	resp, err := h.state.GetVIPHealth(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	writeSuccessResponse(w, resp)
}

// HandleHealth handles GET /api/v1/health.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := h.state.GetAllHealth()
	writeSuccessResponse(w, resp)
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

// handleRegisterTarget handles POST /api/v1/targets.
func (h *Handlers) handleRegisterTarget(w http.ResponseWriter, r *http.Request) {
	var req registerTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	// Validate and apply defaults to the healthcheck config
	req.Healthcheck.ApplyDefaults()
	if err := req.Healthcheck.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	key := VIPKeyFromHC(req.VIP)

	// Register VIP in state
	if err := h.state.RegisterVIP(key, req.Healthcheck, req.Reals); err != nil {
		writeErrorResponse(w, http.StatusConflict, "ALREADY_EXISTS", err.Error())
		return
	}

	// Acquire somarks and register with katran
	var registeredReals []string
	ctx := r.Context()
	for _, real := range req.Reals {
		somark, isNew, err := h.somarks.Acquire(real.Address)
		if err != nil {
			h.rollbackRegistration(ctx, key, registeredReals)
			writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		if isNew {
			if err := h.katran.RegisterDst(ctx, somark, real.Address); err != nil {
				// Release the somark we just acquired
				h.somarks.Release(real.Address)
				h.rollbackRegistration(ctx, key, registeredReals)
				writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR",
					fmt.Sprintf("failed to register somark with katran: %v", err))
				return
			}
		}
		registeredReals = append(registeredReals, real.Address)
	}

	// Notify scheduler
	h.scheduler.NotifyVIPRegistered(key, registeredReals)

	writeSuccessResponse(w, nil)
}

// handleUpdateTarget handles PUT /api/v1/targets.
func (h *Handlers) handleUpdateTarget(w http.ResponseWriter, r *http.Request) {
	var req updateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	req.Healthcheck.ApplyDefaults()
	if err := req.Healthcheck.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	key := VIPKeyFromHC(req.VIP)
	ctx := r.Context()

	oldReals, err := h.state.UpdateVIP(key, req.Healthcheck, req.Reals)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	if req.Reals != nil {
		// Release somarks for old reals
		for _, addr := range oldReals {
			somark, isLast, err := h.somarks.Release(addr)
			if err != nil {
				log.Printf("warning: failed to release somark for %s: %v", addr, err)
				continue
			}
			if isLast {
				if err := h.katran.DeregisterDst(ctx, somark); err != nil {
					log.Printf("warning: failed to deregister somark %d from katran: %v", somark, err)
				}
			}
		}

		// Acquire somarks for new reals
		var newRealAddrs []string
		for _, real := range req.Reals {
			somark, isNew, err := h.somarks.Acquire(real.Address)
			if err != nil {
				log.Printf("warning: failed to acquire somark for %s: %v", real.Address, err)
				continue
			}
			if isNew {
				if err := h.katran.RegisterDst(ctx, somark, real.Address); err != nil {
					log.Printf("warning: failed to register somark with katran for %s: %v", real.Address, err)
					h.somarks.Release(real.Address)
					continue
				}
			}
			newRealAddrs = append(newRealAddrs, real.Address)
		}

		h.scheduler.NotifyVIPUpdated(key, oldReals, newRealAddrs)
	}

	writeSuccessResponse(w, nil)
}

// handleDeleteTarget handles DELETE /api/v1/targets.
func (h *Handlers) handleDeleteTarget(w http.ResponseWriter, r *http.Request) {
	var vipKey types.HCVIPKey
	if err := json.NewDecoder(r.Body).Decode(&vipKey); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}

	key := VIPKeyFromHC(vipKey)
	ctx := r.Context()

	realAddrs, err := h.state.DeregisterVIP(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	// Release somarks
	for _, addr := range realAddrs {
		somark, isLast, err := h.somarks.Release(addr)
		if err != nil {
			log.Printf("warning: failed to release somark for %s: %v", addr, err)
			continue
		}
		if isLast {
			if err := h.katran.DeregisterDst(ctx, somark); err != nil {
				log.Printf("warning: failed to deregister somark %d from katran: %v", somark, err)
			}
		}
	}

	h.scheduler.NotifyVIPDeregistered(key, realAddrs)
	writeSuccessResponse(w, nil)
}

// handleListTargets handles GET /api/v1/targets.
func (h *Handlers) handleListTargets(w http.ResponseWriter, r *http.Request) {
	vips := h.state.ListVIPs()
	writeSuccessResponse(w, vips)
}

// handleAddReals handles POST /api/v1/targets/reals.
func (h *Handlers) handleAddReals(w http.ResponseWriter, r *http.Request) {
	var req realsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}
	if len(req.Reals) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "reals list is empty")
		return
	}

	key := VIPKeyFromHC(req.VIP)
	ctx := r.Context()

	added, skipped, err := h.state.AddReals(key, req.Reals)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	// Acquire somarks for the newly added reals
	var addedAddrs []string
	for _, real := range req.Reals {
		// Only acquire for reals that were actually added (not skipped)
		somark, isNew, err := h.somarks.Acquire(real.Address)
		if err != nil {
			log.Printf("warning: failed to acquire somark for %s: %v", real.Address, err)
			continue
		}
		if isNew {
			if err := h.katran.RegisterDst(ctx, somark, real.Address); err != nil {
				log.Printf("warning: failed to register somark with katran for %s: %v", real.Address, err)
				h.somarks.Release(real.Address)
				continue
			}
		}
		addedAddrs = append(addedAddrs, real.Address)
	}

	h.scheduler.NotifyRealsAdded(key, addedAddrs)

	writeSuccessResponse(w, realsAddedResponse{Added: added, Skipped: skipped})
}

// handleRemoveReals handles DELETE /api/v1/targets/reals.
func (h *Handlers) handleRemoveReals(w http.ResponseWriter, r *http.Request) {
	var req realsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body: "+err.Error())
		return
	}
	if len(req.Reals) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "reals list is empty")
		return
	}

	key := VIPKeyFromHC(req.VIP)
	ctx := r.Context()

	// Collect addresses
	addresses := make([]string, len(req.Reals))
	for i, r := range req.Reals {
		addresses[i] = r.Address
	}

	removed, notFound, err := h.state.RemoveReals(key, addresses)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	// Release somarks for removed reals
	var removedAddrs []string
	for _, addr := range addresses {
		somark, isLast, err := h.somarks.Release(addr)
		if err != nil {
			// Not tracked means it wasn't found in the first place
			continue
		}
		if isLast {
			if err := h.katran.DeregisterDst(ctx, somark); err != nil {
				log.Printf("warning: failed to deregister somark %d from katran: %v", somark, err)
			}
		}
		removedAddrs = append(removedAddrs, addr)
	}

	h.scheduler.NotifyRealsRemoved(key, removedAddrs)

	writeSuccessResponse(w, realsRemovedResponse{Removed: removed, NotFound: notFound})
}

// rollbackRegistration cleans up a partially registered VIP on failure.
func (h *Handlers) rollbackRegistration(ctx context.Context, key VIPKey, registeredReals []string) {
	for _, addr := range registeredReals {
		somark, isLast, err := h.somarks.Release(addr)
		if err != nil {
			continue
		}
		if isLast {
			h.katran.DeregisterDst(ctx, somark)
		}
	}
	h.state.DeregisterVIP(key)
}

// parseVIPQuery extracts a VIPKey from query parameters.
func parseVIPQuery(r *http.Request) (VIPKey, error) {
	addr := r.URL.Query().Get("address")
	if addr == "" {
		return VIPKey{}, fmt.Errorf("missing 'address' query parameter")
	}

	portStr := r.URL.Query().Get("port")
	if portStr == "" {
		return VIPKey{}, fmt.Errorf("missing 'port' query parameter")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return VIPKey{}, fmt.Errorf("invalid 'port' query parameter")
	}

	protoStr := r.URL.Query().Get("proto")
	if protoStr == "" {
		return VIPKey{}, fmt.Errorf("missing 'proto' query parameter")
	}
	proto, err := strconv.Atoi(protoStr)
	if err != nil || proto < 0 || proto > 255 {
		return VIPKey{}, fmt.Errorf("invalid 'proto' query parameter")
	}

	return VIPKey{Address: addr, Port: uint16(port), Proto: uint8(proto)}, nil
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
