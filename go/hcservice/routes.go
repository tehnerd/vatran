package hcservice

import "net/http"

// RegisterRoutes registers all HC service routes on the given ServeMux.
//
// Parameters:
//   - mux: The HTTP serve mux to register routes on.
//   - handlers: The handlers implementing the HC service API.
func RegisterRoutes(mux *http.ServeMux, handlers *Handlers) {
	mux.HandleFunc("/api/v1/targets", handlers.HandleTargets)
	mux.HandleFunc("/api/v1/targets/reals", handlers.HandleTargetReals)
	mux.HandleFunc("/api/v1/health/vip", handlers.HandleHealthVIP)
	mux.HandleFunc("/api/v1/health", handlers.HandleHealth)
	mux.HandleFunc("/health", handlers.HandleServiceHealth)
}
