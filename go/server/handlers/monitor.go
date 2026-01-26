package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// MonitorHandler handles monitor control operations.
type MonitorHandler struct {
	manager *lb.Manager
}

// NewMonitorHandler creates a new MonitorHandler.
//
// Returns a new MonitorHandler instance.
func NewMonitorHandler() *MonitorHandler {
	return &MonitorHandler{
		manager: lb.GetManager(),
	}
}

// HandleStopMonitor handles POST /monitor/stop - stops the monitor.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *MonitorHandler) HandleStopMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	if err := lbInstance.StopMonitor(); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleRestartMonitor handles POST /monitor/restart - restarts the monitor.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *MonitorHandler) HandleRestartMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.RestartMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	if err := lbInstance.RestartMonitor(req.Limit); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}
