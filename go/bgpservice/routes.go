package bgpservice

import "net/http"

// RegisterRoutes registers all BGP service routes on the given ServeMux.
//
// Parameters:
//   - mux: The HTTP serve mux to register routes on.
//   - handlers: The handlers implementing the BGP service API.
func RegisterRoutes(mux *http.ServeMux, handlers *Handlers) {
	mux.HandleFunc("/api/v1/routes/advertise", handlers.HandleAdvertise)
	mux.HandleFunc("/api/v1/routes/withdraw", handlers.HandleWithdraw)
	mux.HandleFunc("/api/v1/routes", handlers.HandleRoutes)
	mux.HandleFunc("/api/v1/routes/vip", handlers.HandleRouteVIP)
	mux.HandleFunc("/api/v1/peers", handlers.HandlePeers)
	mux.HandleFunc("/health", handlers.HandleServiceHealth)
}
