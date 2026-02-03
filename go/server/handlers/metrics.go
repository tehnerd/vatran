package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tehnerd/vatran/go/server/metrics"
)

// MetricsHandler handles Prometheus metrics endpoint.
type MetricsHandler struct {
	handler http.Handler
}

// NewMetricsHandler creates a new MetricsHandler with a custom Prometheus registry.
//
// The handler uses a custom registry containing:
//   - KatranCollector for Katran-specific metrics
//   - Optional Go and Process collectors can be added by modifying this function
//
// Returns a new MetricsHandler instance.
func NewMetricsHandler() *MetricsHandler {
	// Create a custom registry to avoid polluting with default Go metrics
	registry := prometheus.NewRegistry()

	// Register the Katran collector
	collector := metrics.NewKatranCollector()
	registry.MustRegister(collector)

	// Optionally register Go and Process collectors for debugging
	// registry.MustRegister(prometheus.NewGoCollector())
	// registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	return &MetricsHandler{
		handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}),
	}
}

// HandleMetrics handles GET /metrics - returns Prometheus metrics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *MetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.handler.ServeHTTP(w, r)
}
