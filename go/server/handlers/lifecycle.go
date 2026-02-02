package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
	"github.com/tehnerd/vatran/go/server/types"
)

// ErrPathTraversal is returned when a path escapes the base directory.
var ErrPathTraversal = errors.New("path traversal attempt detected")

// sanitizePath joins basePath with userPath and ensures the result does not escape basePath.
// Returns the sanitized absolute path or an error if traversal is detected.
//
// Parameters:
//   - basePath: The trusted base directory path.
//   - userPath: The untrusted user-provided path component.
//
// Returns the sanitized path or ErrPathTraversal if the path escapes basePath.
func sanitizePath(basePath, userPath string) (string, error) {
	// Clean the base path to ensure consistent comparison
	cleanBase := filepath.Clean(basePath)

	// Join and clean the full path
	fullPath := filepath.Clean(filepath.Join(cleanBase, userPath))

	// Verify the result is within the base directory
	// We add a trailing separator to ensure we match the directory, not a prefix
	if !strings.HasPrefix(fullPath, cleanBase+string(filepath.Separator)) && fullPath != cleanBase {
		return "", ErrPathTraversal
	}

	return fullPath, nil
}

const (
	// bpfFSPath is the standard BPF filesystem path.
	bpfFSPath = "/sys/fs/bpf"
)

// LifecycleHandler handles load balancer lifecycle operations.
type LifecycleHandler struct {
	manager    *lb.Manager
	bpfProgDir string
}

// NewLifecycleHandler creates a new LifecycleHandler.
//
// Parameters:
//   - bpfProgDir: Base directory for BPF program files.
//
// Returns a new LifecycleHandler instance.
func NewLifecycleHandler(bpfProgDir string) *LifecycleHandler {
	return &LifecycleHandler{
		manager:    lb.GetManager(),
		bpfProgDir: bpfProgDir,
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

	cfg := requestToConfig(&req, h.bpfProgDir)

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
		cfg = requestToConfig(req.Config, h.bpfProgDir)
	}

	// Resolve reload path relative to bpfProgDir if not absolute
	reloadPath := req.Path
	if reloadPath != "" && !filepath.IsAbs(reloadPath) && h.bpfProgDir != "" {
		var err error
		reloadPath, err = sanitizePath(h.bpfProgDir, reloadPath)
		if err != nil {
			models.WriteError(w, http.StatusBadRequest,
				models.NewInvalidRequestError("invalid reload path: "+err.Error()))
			return
		}
	}

	if err := h.manager.ReloadBalancerProg(reloadPath, cfg); err != nil {
		models.WriteKatranError(w, err)
		return
	}

	models.WriteSuccess(w, nil)
}

// requestToConfig converts a CreateLBRequest to a katran.Config.
//
// Parameters:
//   - req: The CreateLBRequest to convert.
//   - bpfProgDir: Base directory for BPF program files.
//
// BalancerProgPath and HealthcheckingProgPath are resolved relative to bpfProgDir
// if they are not absolute paths and bpfProgDir is set.
// RootMapPath is resolved relative to /sys/fs/bpf/ if it is not an absolute path.
func requestToConfig(req *models.CreateLBRequest, bpfProgDir string) *katran.Config {
	cfg := katran.NewConfig()

	cfg.MainInterface = req.MainInterface
	cfg.V4TunInterface = req.V4TunInterface
	cfg.V6TunInterface = req.V6TunInterface
	cfg.HCInterface = req.HCInterface

	// Resolve BalancerProgPath relative to bpfProgDir
	cfg.BalancerProgPath = req.BalancerProgPath
	if cfg.BalancerProgPath != "" && !filepath.IsAbs(cfg.BalancerProgPath) && bpfProgDir != "" {
		if sanitized, err := sanitizePath(bpfProgDir, cfg.BalancerProgPath); err == nil {
			cfg.BalancerProgPath = sanitized
		}
		// If sanitization fails, keep the original path - it will fail at file access
	}

	// Resolve HealthcheckingProgPath relative to bpfProgDir
	cfg.HealthcheckingProgPath = req.HealthcheckingProgPath
	if cfg.HealthcheckingProgPath != "" && !filepath.IsAbs(cfg.HealthcheckingProgPath) && bpfProgDir != "" {
		if sanitized, err := sanitizePath(bpfProgDir, cfg.HealthcheckingProgPath); err == nil {
			cfg.HealthcheckingProgPath = sanitized
		}
		// If sanitization fails, keep the original path - it will fail at file access
	}

	// Resolve RootMapPath relative to /sys/fs/bpf/
	cfg.RootMapPath = req.RootMapPath
	if cfg.RootMapPath != "" && !filepath.IsAbs(cfg.RootMapPath) {
		if sanitized, err := sanitizePath(bpfFSPath, cfg.RootMapPath); err == nil {
			cfg.RootMapPath = sanitized
		}
		// If sanitization fails, keep the original path - it will fail at file access
	}
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
	} else {
		// Default: EnableHC is true if and only if HealthcheckingProgPath is non-empty
		cfg.EnableHC = cfg.HealthcheckingProgPath != ""
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
	if strings.TrimSpace(req.HashFunction) != "" {
		cfg.HashFunc = katran.HashFunction(types.HashFunctionToInt(req.HashFunction))
	}

	return cfg
}
