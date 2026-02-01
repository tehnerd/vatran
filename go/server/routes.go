package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tehnerd/vatran/go/server/handlers"
	"github.com/tehnerd/vatran/go/server/models"
	"github.com/tehnerd/vatran/go/server/types"
)

const (
	// APIBasePath is the base path for all API endpoints.
	APIBasePath = "/api/v1"
)

// RegisterRoutes registers all API routes on the provided mux.
//
// Parameters:
//   - mux: The http.ServeMux to register routes on.
//   - config: The server configuration.
//   - authHandler: Optional auth handler for login/logout routes.
func RegisterRoutes(mux *http.ServeMux, config *Config, authHandler *handlers.AuthHandler) {
	// Initialize handlers
	lifecycleHandler := handlers.NewLifecycleHandler(config.BPFProgDir)
	vipHandler := handlers.NewVIPHandler()
	realHandler := handlers.NewRealHandler()
	statsHandler := handlers.NewStatsHandler()
	quicHandler := handlers.NewQuicHandler()
	routingHandler := handlers.NewRoutingHandler()
	healthcheckHandler := handlers.NewHealthcheckHandler()
	featuresHandler := handlers.NewFeaturesHandler()
	lruHandler := handlers.NewLRUHandler()
	monitorHandler := handlers.NewMonitorHandler()
	utilsHandler := handlers.NewUtilsHandler()
	configHandler := handlers.NewConfigHandler(config)
	// Set up the config exporter for YAML export
	configHandler.SetConfigExporter(createConfigExporter(config))

	// Health check endpoint (not versioned)
	mux.HandleFunc("/health", handleHealth)

	// Lifecycle endpoints
	mux.HandleFunc(APIBasePath+"/lb/create", lifecycleHandler.HandleCreate)
	mux.HandleFunc(APIBasePath+"/lb/close", lifecycleHandler.HandleClose)
	mux.HandleFunc(APIBasePath+"/lb/status", lifecycleHandler.HandleStatus)
	mux.HandleFunc(APIBasePath+"/lb/load-bpf-progs", lifecycleHandler.HandleLoadBPFProgs)
	mux.HandleFunc(APIBasePath+"/lb/attach-bpf-progs", lifecycleHandler.HandleAttachBPFProgs)
	mux.HandleFunc(APIBasePath+"/lb/reload", lifecycleHandler.HandleReload)

	// VIP endpoints
	mux.HandleFunc(APIBasePath+"/vips", vipHandler.HandleVIPs)
	mux.HandleFunc(APIBasePath+"/vips/flags", vipHandler.HandleVIPFlags)
	mux.HandleFunc(APIBasePath+"/vips/hash-function", vipHandler.HandleHashFunction)

	// Real server endpoints
	mux.HandleFunc(APIBasePath+"/vips/reals", realHandler.HandleVIPReals)
	mux.HandleFunc(APIBasePath+"/vips/reals/batch", realHandler.HandleBatchModifyReals)
	mux.HandleFunc(APIBasePath+"/reals/index", realHandler.HandleRealIndex)
	mux.HandleFunc(APIBasePath+"/reals/flags", realHandler.HandleRealFlags)

	// Statistics endpoints
	mux.HandleFunc(APIBasePath+"/stats/vip", statsHandler.HandleVIPStats)
	mux.HandleFunc(APIBasePath+"/stats/vip/decap", statsHandler.HandleVIPDecapStats)
	mux.HandleFunc(APIBasePath+"/stats/real", statsHandler.HandleRealStats)
	mux.HandleFunc(APIBasePath+"/stats/lru", statsHandler.HandleLRUStats)
	mux.HandleFunc(APIBasePath+"/stats/lru/miss", statsHandler.HandleLRUMissStats)
	mux.HandleFunc(APIBasePath+"/stats/lru/fallback", statsHandler.HandleLRUFallbackStats)
	mux.HandleFunc(APIBasePath+"/stats/lru/global", statsHandler.HandleGlobalLRUStats)
	mux.HandleFunc(APIBasePath+"/stats/icmp-too-big", statsHandler.HandleICMPTooBigStats)
	mux.HandleFunc(APIBasePath+"/stats/ch-drop", statsHandler.HandleCHDropStats)
	mux.HandleFunc(APIBasePath+"/stats/src-routing", statsHandler.HandleSrcRoutingStats)
	mux.HandleFunc(APIBasePath+"/stats/inline-decap", statsHandler.HandleInlineDecapStats)
	mux.HandleFunc(APIBasePath+"/stats/decap", statsHandler.HandleDecapStats)
	mux.HandleFunc(APIBasePath+"/stats/quic-icmp", statsHandler.HandleQuicICMPStats)
	mux.HandleFunc(APIBasePath+"/stats/quic-packets", statsHandler.HandleQuicPacketsStats)
	mux.HandleFunc(APIBasePath+"/stats/tcp-server-id-routing", statsHandler.HandleTCPServerIDRoutingStats)
	mux.HandleFunc(APIBasePath+"/stats/xdp/total", statsHandler.HandleXDPTotalStats)
	mux.HandleFunc(APIBasePath+"/stats/xdp/tx", statsHandler.HandleXDPTXStats)
	mux.HandleFunc(APIBasePath+"/stats/xdp/drop", statsHandler.HandleXDPDropStats)
	mux.HandleFunc(APIBasePath+"/stats/xdp/pass", statsHandler.HandleXDPPassStats)
	mux.HandleFunc(APIBasePath+"/stats/hc-prog", statsHandler.HandleHCProgStats)
	mux.HandleFunc(APIBasePath+"/stats/bpf-map", statsHandler.HandleBPFMapStats)
	mux.HandleFunc(APIBasePath+"/stats/userspace", statsHandler.HandleUserspaceStats)
	mux.HandleFunc(APIBasePath+"/stats/per-core-packets", statsHandler.HandlePerCorePacketsStats)
	mux.HandleFunc(APIBasePath+"/stats/flood-status", statsHandler.HandleFloodStatus)
	mux.HandleFunc(APIBasePath+"/stats/monitor", statsHandler.HandleMonitorStats)

	// QUIC endpoints
	mux.HandleFunc(APIBasePath+"/quic/reals", quicHandler.HandleQuicReals)

	// Routing endpoints
	mux.HandleFunc(APIBasePath+"/routing/src-rules", routingHandler.HandleSrcRoutingRules)
	mux.HandleFunc(APIBasePath+"/routing/src-rules/all", routingHandler.HandleClearSrcRoutingRules)
	mux.HandleFunc(APIBasePath+"/routing/src-rules/size", routingHandler.HandleSrcRoutingRuleSize)
	mux.HandleFunc(APIBasePath+"/routing/decap/inline", routingHandler.HandleInlineDecapDsts)

	// Healthcheck endpoints
	mux.HandleFunc(APIBasePath+"/healthcheck/dsts", healthcheckHandler.HandleHealthcheckerDsts)
	mux.HandleFunc(APIBasePath+"/healthcheck/keys", healthcheckHandler.HandleHCKeys)

	// Features endpoints
	mux.HandleFunc(APIBasePath+"/features/check", featuresHandler.HandleHasFeature)
	mux.HandleFunc(APIBasePath+"/features/install", featuresHandler.HandleInstallFeature)
	mux.HandleFunc(APIBasePath+"/features/remove", featuresHandler.HandleRemoveFeature)

	// LRU endpoints
	mux.HandleFunc(APIBasePath+"/lru", lruHandler.HandleDeleteLRU)
	mux.HandleFunc(APIBasePath+"/lru/vip", lruHandler.HandlePurgeVIPLRU)

	// Monitor endpoints
	mux.HandleFunc(APIBasePath+"/monitor/stop", monitorHandler.HandleStopMonitor)
	mux.HandleFunc(APIBasePath+"/monitor/restart", monitorHandler.HandleRestartMonitor)

	// Utility endpoints
	mux.HandleFunc(APIBasePath+"/utils/mac", utilsHandler.HandleMAC)
	mux.HandleFunc(APIBasePath+"/utils/real-for-flow", utilsHandler.HandleRealForFlow)
	mux.HandleFunc(APIBasePath+"/utils/simulate-packet", utilsHandler.HandleSimulatePacket)
	mux.HandleFunc(APIBasePath+"/utils/prog-fd", utilsHandler.HandleKatranProgFD)
	mux.HandleFunc(APIBasePath+"/utils/hc-prog-fd", utilsHandler.HandleHealthcheckerProgFD)
	mux.HandleFunc(APIBasePath+"/utils/map-fd", utilsHandler.HandleBPFMapFD)
	mux.HandleFunc(APIBasePath+"/utils/global-lru-map-fds", utilsHandler.HandleGlobalLRUMapsFDs)
	mux.HandleFunc(APIBasePath+"/utils/src-ip-encap", utilsHandler.HandleAddSrcIPForPcktEncap)

	// Config export endpoints
	mux.HandleFunc(APIBasePath+"/config/export", configHandler.HandleExportConfig)
	mux.HandleFunc(APIBasePath+"/config/export/json", configHandler.HandleExportConfigJSON)

	// Auth endpoints (if auth is enabled)
	if authHandler != nil {
		mux.HandleFunc("/login", authHandler.HandleLogin)
		mux.HandleFunc("/logout", authHandler.HandleLogout)
	}

	// Serve SPA static files if staticDir is configured
	if config.StaticDir != "" {
		mux.Handle("/", newSPAHandler(config.StaticDir))
	}
}

// handleHealth handles the health check endpoint.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	models.WriteSuccess(w, map[string]string{"status": "ok"})
}

// spaHandler serves a Single Page Application (SPA).
// It serves static files from the configured directory and falls back to index.html
// for any path that doesn't match a static file (supporting client-side routing).
type spaHandler struct {
	staticDir string
	fileServer http.Handler
}

// newSPAHandler creates a new SPA handler.
//
// Parameters:
//   - staticDir: The directory containing the SPA static files.
//
// Returns an http.Handler that serves the SPA.
func newSPAHandler(staticDir string) http.Handler {
	return &spaHandler{
		staticDir:  staticDir,
		fileServer: http.FileServer(http.Dir(staticDir)),
	}
}

// ServeHTTP implements the http.Handler interface for spaHandler.
func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path
	path := filepath.Clean(r.URL.Path)

	// Don't serve API routes through the SPA handler
	if strings.HasPrefix(path, APIBasePath) {
		http.NotFound(w, r)
		return
	}

	// Sanitize the path to prevent directory traversal
	cleanBase := filepath.Clean(h.staticDir)
	fullPath := filepath.Clean(filepath.Join(cleanBase, path))

	// Ensure the path is within staticDir
	if !strings.HasPrefix(fullPath, cleanBase+string(filepath.Separator)) && fullPath != cleanBase {
		http.NotFound(w, r)
		return
	}

	info, err := os.Stat(fullPath)

	// If the file exists and is not a directory, serve it
	if err == nil && !info.IsDir() {
		h.fileServer.ServeHTTP(w, r)
		return
	}

	// For thirdparty dependencies, return 404 if file doesn't exist
	// (don't fall back to index.html for these paths)
	if strings.HasPrefix(path, "/thirdparty/") {
		http.NotFound(w, r)
		return
	}

	// For any other path (including directories and non-existent files),
	// serve index.html for client-side routing
	http.ServeFile(w, r, filepath.Join(h.staticDir, "index.html"))
}

// createConfigExporter creates a config exporter function for the given server config.
func createConfigExporter(config *Config) types.ConfigExporter {
	return types.ConfigExporterFunc(func(katranCfg *types.KatranConfigExport, vips []types.VIPWithBackends) ([]byte, error) {
		return ExportConfigAsYAML(config, katranCfg, vips)
	})
}
