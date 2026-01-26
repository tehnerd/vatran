package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// FeaturesHandler handles feature management operations.
type FeaturesHandler struct {
	manager *lb.Manager
}

// NewFeaturesHandler creates a new FeaturesHandler.
//
// Returns a new FeaturesHandler instance.
func NewFeaturesHandler() *FeaturesHandler {
	return &FeaturesHandler{
		manager: lb.GetManager(),
	}
}

// HandleHasFeature handles GET /features/check - checks if a feature is available.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *FeaturesHandler) HandleHasFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.FeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	hasFeature, err := lbInstance.HasFeature(katran.Feature(req.Feature))
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.FeatureStatusResponse{Available: hasFeature})
}

// HandleInstallFeature handles POST /features/install - installs a feature.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *FeaturesHandler) HandleInstallFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.FeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.InstallFeature(katran.Feature(req.Feature), req.ProgPath); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleRemoveFeature handles POST /features/remove - removes a feature.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *FeaturesHandler) HandleRemoveFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.FeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.RemoveFeature(katran.Feature(req.Feature), req.ProgPath); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
