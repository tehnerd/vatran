package handlers

import (
	"fmt"
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/models"
	"github.com/tehnerd/vatran/go/server/types"
)

// ConfigHandler handles configuration export operations.
type ConfigHandler struct {
	manager        *lb.Manager
	serverCfg      types.ServerConfigProvider
	configExporter types.ConfigExporter
}

// NewConfigHandler creates a new ConfigHandler.
//
// Parameters:
//   - serverCfg: The server configuration provider.
//
// Returns a new ConfigHandler instance.
func NewConfigHandler(serverCfg types.ServerConfigProvider) *ConfigHandler {
	return &ConfigHandler{
		manager:   lb.GetManager(),
		serverCfg: serverCfg,
	}
}

// SetConfigExporter sets the config exporter for YAML export.
//
// Parameters:
//   - exporter: The config exporter to use.
func (h *ConfigHandler) SetConfigExporter(exporter types.ConfigExporter) {
	h.configExporter = exporter
}

// HandleExportConfig handles GET /config/export - exports current config as YAML.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *ConfigHandler) HandleExportConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.configExporter == nil {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("config exporter not configured"))
		return
	}

	// Get LB instance
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	// Get katran config from manager
	katranCfg := h.extractKatranConfigFromManager()

	// Get all VIPs with their backends
	vipsWithBackends, err := h.getAllVIPsWithBackends(lbInstance)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	// Export as YAML
	yamlData, err := h.configExporter.ExportAsYAML(katranCfg, vipsWithBackends)
	if err != nil {
		models.WriteError(w, http.StatusInternalServerError,
			models.NewInternalError("failed to export config: "+err.Error()))
		return
	}

	// Write YAML response
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"katran-config.yaml\"")
	w.WriteHeader(http.StatusOK)
	w.Write(yamlData)
}

// HandleExportConfigJSON handles GET /config/export/json - exports current config as JSON.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *ConfigHandler) HandleExportConfigJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get LB instance
	lbInstance, ok := h.manager.Get()
	if !ok {
		models.WriteError(w, http.StatusServiceUnavailable, models.NewLBNotInitializedError())
		return
	}

	// Get katran config from manager
	katranCfg := h.extractKatranConfigFromManager()

	// Get all VIPs with their backends
	vipsWithBackends, err := h.getAllVIPsWithBackends(lbInstance)
	if err != nil {
		models.WriteKatranError(w, err)
		return
	}

	// Build response
	response := ConfigExportResponse{
		Server:       h.buildServerConfigResponse(),
		LB:           h.buildLBConfigResponse(katranCfg),
		TargetGroups: make(map[string][]BackendExportResponse),
		VIPs:         make([]VIPExportResponse, 0, len(vipsWithBackends)),
	}

	// Build target groups and VIPs
	targetGroupIndex := 0
	backendHash := make(map[string]string)

	for _, vip := range vipsWithBackends {
		hashKey := hashBackendsForExport(vip.Backends)

		var targetGroupName string
		if existingName, ok := backendHash[hashKey]; ok {
			targetGroupName = existingName
		} else {
			targetGroupName = fmt.Sprintf("group-%d", targetGroupIndex)
			targetGroupIndex++
			backends := make([]BackendExportResponse, len(vip.Backends))
			for i, b := range vip.Backends {
				backends[i] = BackendExportResponse{
					Address: b.Address,
					Weight:  b.Weight,
					Flags:   b.Flags,
				}
			}
			response.TargetGroups[targetGroupName] = backends
			backendHash[hashKey] = targetGroupName
		}

		response.VIPs = append(response.VIPs, VIPExportResponse{
			Address:     vip.Address,
			Port:        vip.Port,
			Proto:       types.NumberToProto(vip.Proto),
			TargetGroup: targetGroupName,
			Flags:       vip.Flags,
		})
	}

	models.WriteSuccess(w, response)
}

// getAllVIPsWithBackends retrieves all VIPs and their backends.
func (h *ConfigHandler) getAllVIPsWithBackends(lbInstance *katran.LoadBalancer) ([]types.VIPWithBackends, error) {
	vips, err := lbInstance.GetAllVIPs()
	if err != nil {
		return nil, err
	}

	result := make([]types.VIPWithBackends, 0, len(vips))
	for _, vip := range vips {
		vipKey := katran.VIPKey{
			Address: vip.Address,
			Port:    vip.Port,
			Proto:   vip.Proto,
		}

		// Get flags for this VIP
		flags, _ := lbInstance.GetVIPFlags(vipKey)

		// Get backends for this VIP
		reals, err := lbInstance.GetRealsForVIP(vipKey)
		if err != nil {
			return nil, err
		}

		backends := make([]types.BackendConfig, len(reals))
		for i, real := range reals {
			backends[i] = types.BackendConfig{
				Address: real.Address,
				Weight:  real.Weight,
				Flags:   real.Flags,
			}
		}

		result = append(result, types.VIPWithBackends{
			Address:  vip.Address,
			Port:     vip.Port,
			Proto:    vip.Proto,
			Flags:    flags,
			Backends: backends,
		})
	}

	return result, nil
}

// extractKatranConfigFromManager extracts configuration from the manager.
func (h *ConfigHandler) extractKatranConfigFromManager() *types.KatranConfigExport {
	cfg, ok := h.manager.GetConfig()
	if !ok || cfg == nil {
		return nil
	}

	return &types.KatranConfigExport{
		MainInterface:          cfg.MainInterface,
		HCInterface:            cfg.HCInterface,
		V4TunInterface:         cfg.V4TunInterface,
		V6TunInterface:         cfg.V6TunInterface,
		BalancerProgPath:       cfg.BalancerProgPath,
		HealthcheckingProgPath: cfg.HealthcheckingProgPath,
		RootMapPath:            cfg.RootMapPath,
		RootMapPos:             cfg.RootMapPos,
		UseRootMap:             cfg.UseRootMap,
		DefaultMAC:             cfg.DefaultMAC,
		LocalMAC:               cfg.LocalMAC,
		MaxVIPs:                cfg.MaxVIPs,
		MaxReals:               cfg.MaxReals,
		CHRingSize:             cfg.CHRingSize,
		LRUSize:                cfg.LRUSize,
		GlobalLRUSize:          cfg.GlobalLRUSize,
		MaxLPMSrcSize:          cfg.MaxLPMSrcSize,
		MaxDecapDst:            cfg.MaxDecapDst,
		ForwardingCores:        cfg.ForwardingCores,
		NUMANodes:              cfg.NUMANodes,
		XDPAttachFlags:         cfg.XDPAttachFlags,
		Priority:               cfg.Priority,
		KatranSrcV4:            cfg.KatranSrcV4,
		KatranSrcV6:            cfg.KatranSrcV6,
		EnableHC:               cfg.EnableHC,
		TunnelBasedHCEncap:     cfg.TunnelBasedHCEncap,
		FlowDebug:              cfg.FlowDebug,
		EnableCIDV3:            cfg.EnableCIDV3,
		MemlockUnlimited:       cfg.MemlockUnlimited,
		CleanupOnShutdown:      cfg.CleanupOnShutdown,
		Testing:                cfg.Testing,
		HashFunc:               int(cfg.HashFunc),
	}
}

// buildServerConfigResponse builds the server config for JSON export.
func (h *ConfigHandler) buildServerConfigResponse() ServerConfigExportResponse {
	resp := ServerConfigExportResponse{
		Host:           h.serverCfg.GetHost(),
		Port:           h.serverCfg.GetPort(),
		ReadTimeout:    h.serverCfg.GetReadTimeout(),
		WriteTimeout:   h.serverCfg.GetWriteTimeout(),
		IdleTimeout:    h.serverCfg.GetIdleTimeout(),
		EnableCORS:     h.serverCfg.IsEnableCORS(),
		AllowedOrigins: h.serverCfg.GetAllowedOrigins(),
		EnableLogging:  h.serverCfg.IsEnableLogging(),
		EnableRecovery: h.serverCfg.IsEnableRecovery(),
		StaticDir:      h.serverCfg.GetStaticDir(),
		BPFProgDir:     h.serverCfg.GetBPFProgDir(),
	}

	tlsInfo := h.serverCfg.GetTLS()
	if tlsInfo != nil {
		resp.TLS = &TLSConfigExportResponse{
			CertFile:     tlsInfo.CertFile,
			KeyFile:      tlsInfo.KeyFile,
			ClientCAFile: tlsInfo.ClientCAFile,
		}
	}

	return resp
}

// buildLBConfigResponse builds the LB config for JSON export.
func (h *ConfigHandler) buildLBConfigResponse(cfg *types.KatranConfigExport) *LBConfigExportResponse {
	if cfg == nil {
		return nil
	}

	return &LBConfigExportResponse{
		Interfaces: InterfacesExportResponse{
			Main:        cfg.MainInterface,
			Healthcheck: cfg.HCInterface,
			V4Tunnel:    cfg.V4TunInterface,
			V6Tunnel:    cfg.V6TunInterface,
		},
		Programs: ProgramsExportResponse{
			Balancer:    cfg.BalancerProgPath,
			Healthcheck: cfg.HealthcheckingProgPath,
		},
		RootMap: RootMapExportResponse{
			Enabled:  cfg.UseRootMap,
			Path:     cfg.RootMapPath,
			Position: cfg.RootMapPos,
		},
		MAC: MACExportResponse{
			Default: types.FormatMAC(cfg.DefaultMAC),
			Local:   types.FormatMAC(cfg.LocalMAC),
		},
		Capacity: CapacityExportResponse{
			MaxVIPs:       cfg.MaxVIPs,
			MaxReals:      cfg.MaxReals,
			CHRingSize:    cfg.CHRingSize,
			LRUSize:       cfg.LRUSize,
			GlobalLRUSize: cfg.GlobalLRUSize,
			MaxLPMSrc:     cfg.MaxLPMSrcSize,
			MaxDecapDst:   cfg.MaxDecapDst,
		},
		CPU: CPUExportResponse{
			ForwardingCores: cfg.ForwardingCores,
			NUMANodes:       cfg.NUMANodes,
		},
		XDP: XDPExportResponse{
			AttachFlags: cfg.XDPAttachFlags,
			Priority:    cfg.Priority,
		},
		Encapsulation: EncapsulationExportResponse{
			SrcV4: cfg.KatranSrcV4,
			SrcV6: cfg.KatranSrcV6,
		},
		Features: FeaturesExportResponse{
			EnableHealthcheck:  cfg.EnableHC,
			TunnelBasedHCEncap: cfg.TunnelBasedHCEncap,
			FlowDebug:          cfg.FlowDebug,
			EnableCIDV3:        cfg.EnableCIDV3,
			MemlockUnlimited:   cfg.MemlockUnlimited,
			CleanupOnShutdown:  cfg.CleanupOnShutdown,
			Testing:            cfg.Testing,
		},
		HashFunction: types.IntToHashFunction(cfg.HashFunc),
	}
}

// hashBackendsForExport creates a hash key from a list of backends.
func hashBackendsForExport(backends []types.BackendConfig) string {
	if len(backends) == 0 {
		return "empty"
	}
	var parts []string
	for _, b := range backends {
		parts = append(parts, fmt.Sprintf("%s:%d:%d", b.Address, b.Weight, b.Flags))
	}
	return joinStrings(parts, ",")
}

func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}

// ConfigExportResponse is the JSON response for config export.
type ConfigExportResponse struct {
	Server       ServerConfigExportResponse       `json:"server"`
	LB           *LBConfigExportResponse          `json:"lb,omitempty"`
	TargetGroups map[string][]BackendExportResponse `json:"target_groups"`
	VIPs         []VIPExportResponse              `json:"vips"`
}

// ServerConfigExportResponse contains server config for export.
type ServerConfigExportResponse struct {
	Host           string                   `json:"host"`
	Port           int                      `json:"port"`
	TLS            *TLSConfigExportResponse `json:"tls,omitempty"`
	ReadTimeout    int                      `json:"read_timeout"`
	WriteTimeout   int                      `json:"write_timeout"`
	IdleTimeout    int                      `json:"idle_timeout"`
	EnableCORS     bool                     `json:"enable_cors"`
	AllowedOrigins []string                 `json:"allowed_origins"`
	EnableLogging  bool                     `json:"enable_logging"`
	EnableRecovery bool                     `json:"enable_recovery"`
	StaticDir      string                   `json:"static_dir"`
	BPFProgDir     string                   `json:"bpf_prog_dir"`
}

// TLSConfigExportResponse contains TLS config for export.
type TLSConfigExportResponse struct {
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
	ClientCAFile string `json:"client_ca_file,omitempty"`
}

// LBConfigExportResponse contains LB config for export.
type LBConfigExportResponse struct {
	Interfaces    InterfacesExportResponse    `json:"interfaces"`
	Programs      ProgramsExportResponse      `json:"programs"`
	RootMap       RootMapExportResponse       `json:"root_map"`
	MAC           MACExportResponse           `json:"mac"`
	Capacity      CapacityExportResponse      `json:"capacity"`
	CPU           CPUExportResponse           `json:"cpu"`
	XDP           XDPExportResponse           `json:"xdp"`
	Encapsulation EncapsulationExportResponse `json:"encapsulation"`
	Features      FeaturesExportResponse      `json:"features"`
	HashFunction  string                      `json:"hash_function"`
}

// InterfacesExportResponse contains interface config for export.
type InterfacesExportResponse struct {
	Main        string `json:"main"`
	Healthcheck string `json:"healthcheck"`
	V4Tunnel    string `json:"v4_tunnel"`
	V6Tunnel    string `json:"v6_tunnel"`
}

// ProgramsExportResponse contains program paths for export.
type ProgramsExportResponse struct {
	Balancer    string `json:"balancer"`
	Healthcheck string `json:"healthcheck"`
}

// RootMapExportResponse contains root map config for export.
type RootMapExportResponse struct {
	Enabled  bool   `json:"enabled"`
	Path     string `json:"path"`
	Position uint32 `json:"position"`
}

// MACExportResponse contains MAC addresses for export.
type MACExportResponse struct {
	Default string `json:"default"`
	Local   string `json:"local"`
}

// CapacityExportResponse contains capacity config for export.
type CapacityExportResponse struct {
	MaxVIPs       uint32 `json:"max_vips"`
	MaxReals      uint32 `json:"max_reals"`
	CHRingSize    uint32 `json:"ch_ring_size"`
	LRUSize       uint64 `json:"lru_size"`
	GlobalLRUSize uint32 `json:"global_lru_size"`
	MaxLPMSrc     uint32 `json:"max_lpm_src"`
	MaxDecapDst   uint32 `json:"max_decap_dst"`
}

// CPUExportResponse contains CPU config for export.
type CPUExportResponse struct {
	ForwardingCores []int32 `json:"forwarding_cores"`
	NUMANodes       []int32 `json:"numa_nodes"`
}

// XDPExportResponse contains XDP config for export.
type XDPExportResponse struct {
	AttachFlags uint32 `json:"attach_flags"`
	Priority    uint32 `json:"priority"`
}

// EncapsulationExportResponse contains encapsulation config for export.
type EncapsulationExportResponse struct {
	SrcV4 string `json:"src_v4"`
	SrcV6 string `json:"src_v6"`
}

// FeaturesExportResponse contains feature flags for export.
type FeaturesExportResponse struct {
	EnableHealthcheck  bool `json:"enable_healthcheck"`
	TunnelBasedHCEncap bool `json:"tunnel_based_hc_encap"`
	FlowDebug          bool `json:"flow_debug"`
	EnableCIDV3        bool `json:"enable_cid_v3"`
	MemlockUnlimited   bool `json:"memlock_unlimited"`
	CleanupOnShutdown  bool `json:"cleanup_on_shutdown"`
	Testing            bool `json:"testing"`
}

// BackendExportResponse contains backend config for export.
type BackendExportResponse struct {
	Address string `json:"address"`
	Weight  uint32 `json:"weight"`
	Flags   uint8  `json:"flags"`
}

// VIPExportResponse contains VIP config for export.
type VIPExportResponse struct {
	Address     string `json:"address"`
	Port        uint16 `json:"port"`
	Proto       string `json:"proto"`
	TargetGroup string `json:"target_group"`
	Flags       uint32 `json:"flags"`
}
