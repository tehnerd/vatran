package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// LRUHandler handles LRU management operations.
type LRUHandler struct {
	manager *lb.Manager
}

// NewLRUHandler creates a new LRUHandler.
//
// Returns a new LRUHandler instance.
func NewLRUHandler() *LRUHandler {
	return &LRUHandler{
		manager: lb.GetManager(),
	}
}

// HandleDeleteLRU handles DELETE /lru - deletes an LRU entry.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LRUHandler) HandleDeleteLRU(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.DeleteLRURequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	dstVip := katran.VIPKey{
		Address: req.DstVIP.Address,
		Port:    req.DstVIP.Port,
		Proto:   req.DstVIP.Proto,
	}

	maps, err := lbInstance.DeleteLRU(dstVip, req.SrcIP, req.SrcPort)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.DeleteLRUResponse{Maps: maps})
}

// HandlePurgeVIPLRU handles DELETE /lru/vip - purges all LRU entries for a VIP.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LRUHandler) HandlePurgeVIPLRU(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.PurgeVIPLRURequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	dstVip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	deletedCount, err := lbInstance.PurgeVIPLRU(dstVip)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.PurgeVIPLRUResponse{DeletedCount: deletedCount})
}
