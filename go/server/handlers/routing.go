package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// RoutingHandler handles routing and decapsulation operations.
type RoutingHandler struct {
	manager *lb.Manager
}

// NewRoutingHandler creates a new RoutingHandler.
//
// Returns a new RoutingHandler instance.
func NewRoutingHandler() *RoutingHandler {
	return &RoutingHandler{
		manager: lb.GetManager(),
	}
}

// HandleSrcRoutingRules handles source routing rule operations based on HTTP method.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RoutingHandler) HandleSrcRoutingRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetSrcRoutingRules(w, r)
	case http.MethodPost:
		h.handleAddSrcRoutingRule(w, r)
	case http.MethodDelete:
		h.handleDelSrcRoutingRule(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetSrcRoutingRules handles GET /routing/src-rules - gets all source routing rules.
func (h *RoutingHandler) handleGetSrcRoutingRules(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	rules, err := lbInstance.GetSrcRoutingRules()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	response := make([]models.SrcRoutingRuleResponse, len(rules))
	for i, rule := range rules {
		response[i] = models.SrcRoutingRuleResponse{
			Src: rule.Src,
			Dst: rule.Dst,
		}
	}

	models.WriteSuccess(w, response)
}

// handleAddSrcRoutingRule handles POST /routing/src-rules - adds source routing rules.
func (h *RoutingHandler) handleAddSrcRoutingRule(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.AddSrcRoutingRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.AddSrcRoutingRule(req.SrcPrefixes, req.Dst); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelSrcRoutingRule handles DELETE /routing/src-rules - deletes source routing rules.
func (h *RoutingHandler) handleDelSrcRoutingRule(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.DelSrcRoutingRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.DelSrcRoutingRule(req.SrcPrefixes); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleClearSrcRoutingRules handles DELETE /routing/src-rules/all - clears all source routing rules.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RoutingHandler) HandleClearSrcRoutingRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	if err := lbInstance.ClearAllSrcRoutingRules(); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleSrcRoutingRuleSize handles GET /routing/src-rules/size - gets source routing rule count.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RoutingHandler) HandleSrcRoutingRuleSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	size, err := lbInstance.GetSrcRoutingRuleSize()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.SrcRoutingRuleSizeResponse{Size: size})
}

// HandleInlineDecapDsts handles inline decapsulation destination operations.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *RoutingHandler) HandleInlineDecapDsts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetInlineDecapDsts(w, r)
	case http.MethodPost:
		h.handleAddInlineDecapDst(w, r)
	case http.MethodDelete:
		h.handleDelInlineDecapDst(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetInlineDecapDsts handles GET /routing/decap/inline - gets all inline decap destinations.
func (h *RoutingHandler) handleGetInlineDecapDsts(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	dsts, err := lbInstance.GetInlineDecapDsts()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, dsts)
}

// handleAddInlineDecapDst handles POST /routing/decap/inline - adds an inline decap destination.
func (h *RoutingHandler) handleAddInlineDecapDst(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.InlineDecapDstRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.AddInlineDecapDst(req.Dst); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// handleDelInlineDecapDst handles DELETE /routing/decap/inline - deletes an inline decap destination.
func (h *RoutingHandler) handleDelInlineDecapDst(w http.ResponseWriter, r *http.Request) {
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.InlineDecapDstRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.DelInlineDecapDst(req.Dst); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
