package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// QuicHandler handles QUIC management operations.
type QuicHandler struct {
	manager *lb.Manager
}

// NewQuicHandler creates a new QuicHandler.
//
// Returns a new QuicHandler instance.
func NewQuicHandler() *QuicHandler {
	return &QuicHandler{
		manager: lb.GetManager(),
	}
}

// HandleQuicReals handles QUIC real mappings operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *QuicHandler) HandleQuicReals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetQuicReals(w, r)
	case http.MethodPut:
		h.handleModifyQuicReals(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetQuicReals handles GET /quic/reals - gets all QUIC real mappings.
func (h *QuicHandler) handleGetQuicReals(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	reals, err := lbInstance.GetQuicRealsMapping()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	response := make([]models.QuicRealResponse, len(reals))
	for i, real := range reals {
		response[i] = models.QuicRealResponse{
			Address: real.Address,
			ID:      real.ID,
		}
	}

	models.WriteSuccess(w, response)
}

// handleModifyQuicReals handles PUT /quic/reals - modifies QUIC real mappings.
func (h *QuicHandler) handleModifyQuicReals(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.ModifyQuicRealsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	reals := make([]katran.QuicReal, len(req.Reals))
	for i, r := range req.Reals {
		reals[i] = katran.QuicReal{
			Address: r.Address,
			ID:      r.ID,
		}
	}

	if err := lbInstance.ModifyQuicRealsMapping(katran.ModifyAction(req.Action), reals); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
