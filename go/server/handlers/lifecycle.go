package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
)

// LifecycleHandler handles load balancer lifecycle operations.
type LifecycleHandler struct {
	manager *lb.Manager
}

// NewLifecycleHandler creates a new LifecycleHandler.
//
// Returns a new LifecycleHandler instance.
func NewLifecycleHandler() *LifecycleHandler {
	return &LifecycleHandler{
		manager: lb.GetManager(),
	}
}

// HandleCreate handles POST /lb/create - creates the load balancer.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	cfg := requestToConfig(&req)

	if err := h.manager.Create(cfg); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteCreated(w, nil)
}

// HandleClose handles POST /lb/close - closes the load balancer.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleClose(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.Close(); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleStatus handles GET /lb/status - gets the load balancer status.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	initialized, ready := h.manager.Status()
	models.WriteSuccess(w, models.LBStatusResponse{
		Initialized: initialized,
		Ready:       ready,
	})
}

// HandleLoadBPFProgs handles POST /lb/load-bpf-progs - loads BPF programs.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleLoadBPFProgs(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.LoadBPFProgs(); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleAttachBPFProgs handles POST /lb/attach-bpf-progs - attaches BPF programs.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleAttachBPFProgs(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.AttachBPFProgs(); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// HandleReload handles POST /lb/reload - reloads the balancer program.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *LifecycleHandler) HandleReload(w http.ResponseWriter, r *http.Request) {
	var req models.ReloadBalancerProgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteError(w, http.StatusBadRequest,
			models.NewInvalidRequestError("invalid request body: "+err.Error()))
		return
	}

	var cfg *katran.Config
	if req.Config != nil {
		cfg = requestToConfig(req.Config)
	}

	if err := h.manager.ReloadBalancerProg(req.Path, cfg); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// requestToConfig converts a CreateLBRequest to a katran.Config.
func requestToConfig(req *models.CreateLBRequest) *katran.Config {
	cfg := katran.NewConfig()

	cfg.MainInterface = req.MainInterface
	cfg.V4TunInterface = req.V4TunInterface
	cfg.V6TunInterface = req.V6TunInterface
	cfg.HCInterface = req.HCInterface
	cfg.BalancerProgPath = req.BalancerProgPath
	cfg.HealthcheckingProgPath = req.HealthcheckingProgPath
	cfg.RootMapPath = req.RootMapPath
	cfg.KatranSrcV4 = req.KatranSrcV4
	cfg.KatranSrcV6 = req.KatranSrcV6

	// Parse MAC addresses
	if req.DefaultMAC != "" {
		if mac, err := parseMAC(req.DefaultMAC); err == nil {
			cfg.DefaultMAC = mac
		}
	}
	if req.LocalMAC != "" {
		if mac, err := parseMAC(req.LocalMAC); err == nil {
			cfg.LocalMAC = mac
		}
	}

	// Apply non-zero values
	if req.RootMapPos > 0 {
		cfg.RootMapPos = req.RootMapPos
	}
	if req.UseRootMap != nil {
		cfg.UseRootMap = *req.UseRootMap
	}
	if req.MaxVIPs > 0 {
		cfg.MaxVIPs = req.MaxVIPs
	}
	if req.MaxReals > 0 {
		cfg.MaxReals = req.MaxReals
	}
	if req.CHRingSize > 0 {
		cfg.CHRingSize = req.CHRingSize
	}
	if req.LRUSize > 0 {
		cfg.LRUSize = req.LRUSize
	}
	if req.MaxLPMSrcSize > 0 {
		cfg.MaxLPMSrcSize = req.MaxLPMSrcSize
	}
	if req.MaxDecapDst > 0 {
		cfg.MaxDecapDst = req.MaxDecapDst
	}
	if req.GlobalLRUSize > 0 {
		cfg.GlobalLRUSize = req.GlobalLRUSize
	}
	if req.EnableHC != nil {
		cfg.EnableHC = *req.EnableHC
	}
	if req.TunnelBasedHCEncap != nil {
		cfg.TunnelBasedHCEncap = *req.TunnelBasedHCEncap
	}
	cfg.Testing = req.Testing
	if req.MemlockUnlimited != nil {
		cfg.MemlockUnlimited = *req.MemlockUnlimited
	}
	cfg.FlowDebug = req.FlowDebug
	cfg.EnableCIDV3 = req.EnableCIDV3
	if req.CleanupOnShutdown != nil {
		cfg.CleanupOnShutdown = *req.CleanupOnShutdown
	}
	if len(req.ForwardingCores) > 0 {
		cfg.ForwardingCores = req.ForwardingCores
	}
	if len(req.NUMANodes) > 0 {
		cfg.NUMANodes = req.NUMANodes
	}
	if req.XDPAttachFlags > 0 {
		cfg.XDPAttachFlags = req.XDPAttachFlags
	}
	if req.Priority > 0 {
		cfg.Priority = req.Priority
	}
	if req.MainInterfaceIndex > 0 {
		cfg.MainInterfaceIndex = req.MainInterfaceIndex
	}
	if req.HCInterfaceIndex > 0 {
		cfg.HCInterfaceIndex = req.HCInterfaceIndex
	}
	cfg.HashFunc = katran.HashFunction(req.HashFunc)

	return cfg
}
