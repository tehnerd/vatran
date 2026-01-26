package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// StatsHandler handles statistics operations.
type StatsHandler struct {
	manager *lb.Manager
}

// NewStatsHandler creates a new StatsHandler.
//
// Returns a new StatsHandler instance.
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		manager: lb.GetManager(),
	}
}

// HandleVIPStats handles GET /stats/vip - gets VIP statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleVIPStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.VIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	stats, err := lbInstance.GetStatsForVIP(vip)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleVIPDecapStats handles GET /stats/vip/decap - gets VIP decap statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleVIPDecapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.VIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	vip := katran.VIPKey{
		Address: req.Address,
		Port:    req.Port,
		Proto:   req.Proto,
	}

	stats, err := lbInstance.GetDecapStatsForVIP(vip)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleRealStats handles GET /stats/real - gets real server statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleRealStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.GetRealStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	stats, err := lbInstance.GetRealStats(req.Index)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleLRUStats handles GET /stats/lru - gets LRU statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleLRUStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetLRUStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleLRUMissStats handles GET /stats/lru/miss - gets LRU miss statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleLRUMissStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetLRUMissStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleLRUFallbackStats handles GET /stats/lru/fallback - gets LRU fallback statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleLRUFallbackStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetLRUFallbackStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleGlobalLRUStats handles GET /stats/lru/global - gets global LRU statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleGlobalLRUStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetGlobalLRUStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleICMPTooBigStats handles GET /stats/icmp-too-big - gets ICMP too big statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleICMPTooBigStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetICMPTooBigStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleCHDropStats handles GET /stats/ch-drop - gets consistent hash drop statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleCHDropStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetCHDropStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleSrcRoutingStats handles GET /stats/src-routing - gets source routing statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleSrcRoutingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetSrcRoutingStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleInlineDecapStats handles GET /stats/inline-decap - gets inline decap statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleInlineDecapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetInlineDecapStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleDecapStats handles GET /stats/decap - gets decap statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleDecapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetDecapStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleQuicICMPStats handles GET /stats/quic-icmp - gets QUIC ICMP statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleQuicICMPStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetQuicICMPStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleQuicPacketsStats handles GET /stats/quic-packets - gets QUIC packets statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleQuicPacketsStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetQuicPacketsStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.QuicPacketsStatsResponse{
		CHRouted:                 stats.CHRouted,
		CIDInitial:               stats.CIDInitial,
		CIDInvalidServerID:       stats.CIDInvalidServerID,
		CIDInvalidServerIDSample: stats.CIDInvalidServerIDSample,
		CIDRouted:                stats.CIDRouted,
		CIDUnknownRealDropped:    stats.CIDUnknownRealDropped,
		CIDV0:                    stats.CIDV0,
		CIDV1:                    stats.CIDV1,
		CIDV2:                    stats.CIDV2,
		CIDV3:                    stats.CIDV3,
		DstMatchInLRU:            stats.DstMatchInLRU,
		DstMismatchInLRU:         stats.DstMismatchInLRU,
		DstNotFoundInLRU:         stats.DstNotFoundInLRU,
	})
}

// HandleTCPServerIDRoutingStats handles GET /stats/tcp-server-id-routing - gets TPR statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleTCPServerIDRoutingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetTCPServerIDRoutingStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.TPRPacketsStatsResponse{
		CHRouted:         stats.CHRouted,
		DstMismatchInLRU: stats.DstMismatchInLRU,
		SIDRouted:        stats.SIDRouted,
		TCPSyn:           stats.TCPSyn,
	})
}

// HandleXDPTotalStats handles GET /stats/xdp/total - gets XDP total statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleXDPTotalStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetXDPTotalStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleXDPTXStats handles GET /stats/xdp/tx - gets XDP TX statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleXDPTXStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetXDPTXStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleXDPDropStats handles GET /stats/xdp/drop - gets XDP drop statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleXDPDropStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetXDPDropStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleXDPPassStats handles GET /stats/xdp/pass - gets XDP pass statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleXDPPassStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetXDPPassStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.LBStatsResponse{V1: stats.V1, V2: stats.V2})
}

// HandleHCProgStats handles GET /stats/hc-prog - gets healthcheck program statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleHCProgStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetHCProgStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.HCStatsResponse{
		PacketsProcessed: stats.PacketsProcessed,
		PacketsDropped:   stats.PacketsDropped,
		PacketsSkipped:   stats.PacketsSkipped,
		PacketsTooBig:    stats.PacketsTooBig,
	})
}

// HandleBPFMapStats handles GET /stats/bpf-map - gets BPF map statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleBPFMapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	var req models.GetBPFMapStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	stats, err := lbInstance.GetBPFMapStats(req.MapName)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.BPFMapStatsResponse{
		MaxEntries:     stats.MaxEntries,
		CurrentEntries: stats.CurrentEntries,
	})
}

// HandleUserspaceStats handles GET /stats/userspace - gets userspace statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleUserspaceStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetUserspaceStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.UserspaceStatsResponse{
		BPFFailedCalls:       stats.BPFFailedCalls,
		AddrValidationFailed: stats.AddrValidationFailed,
	})
}

// HandlePerCorePacketsStats handles GET /stats/per-core-packets - gets per-core packets statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandlePerCorePacketsStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetPerCorePacketsStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.PerCorePacketsStatsResponse{Counts: stats})
}

// HandleFloodStatus handles GET /stats/flood-status - gets flood status.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleFloodStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	underFlood, err := lbInstance.IsUnderFlood()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.FloodStatusResponse{UnderFlood: underFlood})
}

// HandleMonitorStats handles GET /stats/monitor - gets monitor statistics.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *StatsHandler) HandleMonitorStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	stats, err := lbInstance.GetMonitorStats()
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, models.MonitorStatsResponse{
		Limit:      stats.Limit,
		Amount:     stats.Amount,
		BufferFull: stats.BufferFull,
	})
}
