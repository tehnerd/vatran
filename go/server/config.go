package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tehnerd/vatran/go/server/types"
	"gopkg.in/yaml.v3"
)

// TLSConfig contains TLS configuration for HTTPS.
type TLSConfig struct {
	// CertFile is the path to the TLS certificate file.
	CertFile string
	// KeyFile is the path to the TLS private key file.
	KeyFile string
	// MinVersion is the minimum TLS version (default: TLS 1.2).
	MinVersion uint16
	// ClientAuth specifies the client authentication policy for mTLS.
	ClientAuth tls.ClientAuthType
	// ClientCAFile is the path to the client CA certificate file for mTLS.
	ClientCAFile string
}

// AuthConfig contains authentication configuration.
type AuthConfig struct {
	// Enabled indicates whether authentication is enabled.
	Enabled bool
	// DatabasePath is the path to the SQLite database for users and sessions.
	DatabasePath string
	// AllowLocalhost bypasses authentication for localhost requests.
	AllowLocalhost bool
	// SessionTimeout is the session timeout in hours (default: 24).
	SessionTimeout int
	// BcryptCost is the bcrypt cost factor for password hashing (default: 12).
	BcryptCost int
	// ExemptPaths is a list of paths exempt from authentication.
	ExemptPaths []string
	// BootstrapAdmin contains bootstrap admin user configuration.
	BootstrapAdmin *BootstrapAdminConfig
}

// BootstrapAdminConfig contains bootstrap admin user configuration.
type BootstrapAdminConfig struct {
	// Username is the admin username.
	Username string
	// Password is the admin password (will be hashed on first run).
	Password string
}

// Config contains the server configuration.
type Config struct {
	// Host is the host to bind to (default: "").
	Host string
	// Port is the port to listen on (default: 8080).
	Port int
	// TLS contains TLS configuration. If nil, HTTP is used.
	TLS *TLSConfig
	// Auth contains authentication configuration. If nil, auth is disabled.
	Auth *AuthConfig
	// ReadTimeout is the maximum duration for reading the entire request (default: 30s).
	ReadTimeout int
	// WriteTimeout is the maximum duration before timing out writes of the response (default: 30s).
	WriteTimeout int
	// IdleTimeout is the maximum amount of time to wait for the next request (default: 120s).
	IdleTimeout int
	// EnableCORS enables CORS middleware (default: false).
	EnableCORS bool
	// AllowedOrigins is a list of allowed CORS origins. Use "*" for all origins.
	AllowedOrigins []string
	// EnableLogging enables request logging middleware (default: true).
	EnableLogging bool
	// EnableRecovery enables panic recovery middleware (default: true).
	EnableRecovery bool
	// StaticDir is the path to the directory containing static files for the SPA.
	// If empty, no static files are served.
	StaticDir string
	// BPFProgDir is the base directory for BPF program files.
	// BalancerProgPath and HealthcheckingProgPath in requests are relative to this directory.
	BPFProgDir string
}

// DefaultConfig returns a new Config with default values.
//
// Returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Host:           "",
		Port:           8080,
		TLS:            nil,
		ReadTimeout:    30,
		WriteTimeout:   30,
		IdleTimeout:    120,
		EnableCORS:     false,
		AllowedOrigins: []string{"*"},
		EnableLogging:  true,
		EnableRecovery: true,
	}
}

// Addr returns the address string in the format "host:port".
//
// Returns the address string.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsTLS returns whether TLS is enabled.
//
// Returns true if TLS is configured.
func (c *Config) IsTLS() bool {
	return c.TLS != nil && c.TLS.CertFile != "" && c.TLS.KeyFile != ""
}

// Validate validates the configuration.
//
// Returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.TLS != nil {
		if c.TLS.CertFile == "" {
			return fmt.Errorf("TLS cert file is required when TLS is enabled")
		}
		if c.TLS.KeyFile == "" {
			return fmt.Errorf("TLS key file is required when TLS is enabled")
		}
	}
	return nil
}

// GetHost returns the host.
func (c *Config) GetHost() string { return c.Host }

// GetPort returns the port.
func (c *Config) GetPort() int { return c.Port }

// GetReadTimeout returns the read timeout.
func (c *Config) GetReadTimeout() int { return c.ReadTimeout }

// GetWriteTimeout returns the write timeout.
func (c *Config) GetWriteTimeout() int { return c.WriteTimeout }

// GetIdleTimeout returns the idle timeout.
func (c *Config) GetIdleTimeout() int { return c.IdleTimeout }

// IsEnableCORS returns whether CORS is enabled.
func (c *Config) IsEnableCORS() bool { return c.EnableCORS }

// GetAllowedOrigins returns the allowed origins.
func (c *Config) GetAllowedOrigins() []string { return c.AllowedOrigins }

// IsEnableLogging returns whether logging is enabled.
func (c *Config) IsEnableLogging() bool { return c.EnableLogging }

// IsEnableRecovery returns whether recovery is enabled.
func (c *Config) IsEnableRecovery() bool { return c.EnableRecovery }

// GetStaticDir returns the static directory.
func (c *Config) GetStaticDir() string { return c.StaticDir }

// GetBPFProgDir returns the BPF program directory.
func (c *Config) GetBPFProgDir() string { return c.BPFProgDir }

// GetAuth returns the authentication configuration.
func (c *Config) GetAuth() *AuthConfig { return c.Auth }

// IsAuthEnabled returns whether authentication is enabled.
func (c *Config) IsAuthEnabled() bool { return c.Auth != nil && c.Auth.Enabled }

// GetTLS returns the TLS configuration info.
func (c *Config) GetTLS() *types.TLSConfigInfo {
	if c.TLS == nil {
		return nil
	}
	return &types.TLSConfigInfo{
		CertFile:     c.TLS.CertFile,
		KeyFile:      c.TLS.KeyFile,
		MinVersion:   c.TLS.MinVersion,
		ClientAuth:   c.TLS.ClientAuth,
		ClientCAFile: c.TLS.ClientCAFile,
	}
}

// FullConfig is the top-level configuration structure for YAML config files.
// It contains server settings, load balancer configuration, target groups, and VIPs.
type FullConfig struct {
	// Server contains HTTP/HTTPS server configuration.
	Server ServerYAMLConfig `yaml:"server"`
	// LB contains load balancer configuration.
	LB LBConfig `yaml:"lb"`
	// TargetGroups defines named groups of backend servers.
	TargetGroups map[string][]types.BackendConfig `yaml:"target_groups"`
	// VIPs defines virtual IPs that reference target groups.
	VIPs []VIPConfig `yaml:"vips"`
}

// ServerYAMLConfig contains server configuration from YAML.
type ServerYAMLConfig struct {
	// Host is the host to bind to (empty = all interfaces).
	Host string `yaml:"host"`
	// Port is the port to listen on.
	Port int `yaml:"port"`
	// TLS contains TLS configuration (optional).
	TLS *TLSYAMLConfig `yaml:"tls,omitempty"`
	// Auth contains authentication configuration (optional).
	Auth *AuthYAMLConfig `yaml:"auth,omitempty"`
	// ReadTimeout is the read timeout in seconds.
	ReadTimeout int `yaml:"read_timeout"`
	// WriteTimeout is the write timeout in seconds.
	WriteTimeout int `yaml:"write_timeout"`
	// IdleTimeout is the idle timeout in seconds.
	IdleTimeout int `yaml:"idle_timeout"`
	// EnableCORS enables CORS middleware.
	EnableCORS bool `yaml:"enable_cors"`
	// AllowedOrigins is a list of allowed CORS origins.
	AllowedOrigins []string `yaml:"allowed_origins"`
	// EnableLogging enables request logging.
	EnableLogging *bool `yaml:"enable_logging,omitempty"`
	// EnableRecovery enables panic recovery.
	EnableRecovery *bool `yaml:"enable_recovery,omitempty"`
	// StaticDir is the path to static files directory.
	StaticDir string `yaml:"static_dir"`
	// BPFProgDir is the base directory for BPF programs.
	BPFProgDir string `yaml:"bpf_prog_dir"`
}

// TLSYAMLConfig contains TLS configuration from YAML.
type TLSYAMLConfig struct {
	// CertFile is the path to the TLS certificate.
	CertFile string `yaml:"cert_file"`
	// KeyFile is the path to the TLS private key.
	KeyFile string `yaml:"key_file"`
	// MinVersion is the minimum TLS version ("1.2" or "1.3").
	MinVersion string `yaml:"min_version"`
	// ClientAuth is the client authentication policy ("none", "request", "require", "verify").
	ClientAuth string `yaml:"client_auth"`
	// ClientCAFile is the path to the client CA file for mTLS.
	ClientCAFile string `yaml:"client_ca_file"`
}

// AuthYAMLConfig contains authentication configuration from YAML.
type AuthYAMLConfig struct {
	// Enabled indicates whether authentication is enabled.
	Enabled bool `yaml:"enabled"`
	// DatabasePath is the path to the SQLite database.
	DatabasePath string `yaml:"database_path"`
	// AllowLocalhost bypasses authentication for localhost requests.
	AllowLocalhost bool `yaml:"allow_localhost"`
	// SessionTimeout is the session timeout in hours.
	SessionTimeout int `yaml:"session_timeout"`
	// BcryptCost is the bcrypt cost factor.
	BcryptCost int `yaml:"bcrypt_cost"`
	// ExemptPaths is a list of paths exempt from authentication.
	ExemptPaths []string `yaml:"exempt_paths"`
	// BootstrapAdmin contains bootstrap admin user configuration.
	BootstrapAdmin *BootstrapAdminYAMLConfig `yaml:"bootstrap_admin,omitempty"`
}

// BootstrapAdminYAMLConfig contains bootstrap admin user configuration from YAML.
type BootstrapAdminYAMLConfig struct {
	// Username is the admin username.
	Username string `yaml:"username"`
	// Password is the admin password.
	Password string `yaml:"password"`
}

// LBConfig contains load balancer configuration from YAML.
type LBConfig struct {
	// Interfaces contains network interface configuration.
	Interfaces InterfacesConfig `yaml:"interfaces"`
	// Programs contains BPF program paths.
	Programs ProgramsConfig `yaml:"programs"`
	// RootMap contains root map configuration.
	RootMap RootMapConfig `yaml:"root_map"`
	// MAC contains MAC address configuration.
	MAC MACConfig `yaml:"mac"`
	// Capacity contains capacity limits.
	Capacity CapacityConfig `yaml:"capacity"`
	// CPU contains CPU and NUMA configuration.
	CPU CPUConfig `yaml:"cpu"`
	// XDP contains XDP configuration.
	XDP XDPConfig `yaml:"xdp"`
	// Encapsulation contains GUE encapsulation source addresses.
	Encapsulation EncapsulationConfig `yaml:"encapsulation"`
	// Features contains feature flags.
	Features FeaturesConfig `yaml:"features"`
	// HashFunction is the hash function algorithm ("maglev" or "maglev_v2").
	HashFunction string `yaml:"hash_function"`
}

// InterfacesConfig contains network interface configuration.
type InterfacesConfig struct {
	// Main is the main interface for XDP attachment (required).
	Main string `yaml:"main"`
	// Healthcheck is the interface for healthcheck BPF (defaults to main).
	Healthcheck string `yaml:"healthcheck"`
	// V4Tunnel is the IPv4 tunnel interface for HC encapsulation.
	V4Tunnel string `yaml:"v4_tunnel"`
	// V6Tunnel is the IPv6 tunnel interface for HC encapsulation.
	V6Tunnel string `yaml:"v6_tunnel"`
}

// ProgramsConfig contains BPF program paths.
type ProgramsConfig struct {
	// Balancer is the path to the balancer BPF program.
	Balancer string `yaml:"balancer"`
	// Healthcheck is the path to the healthcheck BPF program.
	Healthcheck string `yaml:"healthcheck"`
}

// RootMapConfig contains root map configuration.
type RootMapConfig struct {
	// Enabled indicates whether to use root map mode.
	Enabled *bool `yaml:"enabled,omitempty"`
	// Path is the path to the pinned root map.
	Path string `yaml:"path"`
	// Position is the position in the root map.
	Position uint32 `yaml:"position"`
}

// MACConfig contains MAC address configuration.
type MACConfig struct {
	// Default is the gateway/router MAC address.
	Default string `yaml:"default"`
	// Local is the local server MAC address.
	Local string `yaml:"local"`
}

// CapacityConfig contains capacity limits.
type CapacityConfig struct {
	// MaxVIPs is the maximum number of VIPs.
	MaxVIPs uint32 `yaml:"max_vips"`
	// MaxReals is the maximum number of real servers.
	MaxReals uint32 `yaml:"max_reals"`
	// CHRingSize is the consistent hashing ring size.
	CHRingSize uint32 `yaml:"ch_ring_size"`
	// LRUSize is the per-CPU LRU table size.
	LRUSize uint64 `yaml:"lru_size"`
	// GlobalLRUSize is the per-CPU global LRU size.
	GlobalLRUSize uint32 `yaml:"global_lru_size"`
	// MaxLPMSrc is the maximum source routing LPM entries.
	MaxLPMSrc uint32 `yaml:"max_lpm_src"`
	// MaxDecapDst is the maximum inline decap destinations.
	MaxDecapDst uint32 `yaml:"max_decap_dst"`
}

// CPUConfig contains CPU and NUMA configuration.
type CPUConfig struct {
	// ForwardingCores is the list of CPU cores for forwarding.
	ForwardingCores []int32 `yaml:"forwarding_cores"`
	// NUMANodes maps forwarding cores to NUMA nodes.
	NUMANodes []int32 `yaml:"numa_nodes"`
}

// XDPConfig contains XDP configuration.
type XDPConfig struct {
	// AttachFlags are the XDP attachment flags.
	AttachFlags uint32 `yaml:"attach_flags"`
	// Priority is the TC priority for healthcheck program.
	Priority uint32 `yaml:"priority"`
}

// EncapsulationConfig contains GUE encapsulation source addresses.
type EncapsulationConfig struct {
	// SrcV4 is the IPv4 source address.
	SrcV4 string `yaml:"src_v4"`
	// SrcV6 is the IPv6 source address.
	SrcV6 string `yaml:"src_v6"`
}

// FeaturesConfig contains feature flags.
type FeaturesConfig struct {
	// EnableHealthcheck enables healthcheck program.
	EnableHealthcheck *bool `yaml:"enable_healthcheck,omitempty"`
	// TunnelBasedHCEncap uses tunnel interfaces for HC encapsulation.
	TunnelBasedHCEncap *bool `yaml:"tunnel_based_hc_encap,omitempty"`
	// FlowDebug enables flow debugging.
	FlowDebug bool `yaml:"flow_debug"`
	// EnableCIDV3 enables QUIC CID v3 support.
	EnableCIDV3 bool `yaml:"enable_cid_v3"`
	// MemlockUnlimited sets RLIMIT_MEMLOCK to unlimited.
	MemlockUnlimited *bool `yaml:"memlock_unlimited,omitempty"`
	// CleanupOnShutdown cleans up BPF resources on shutdown.
	CleanupOnShutdown *bool `yaml:"cleanup_on_shutdown,omitempty"`
	// Testing enables testing mode.
	Testing bool `yaml:"testing"`
}

// BackendConfig is an alias to types.BackendConfig for convenience.
type BackendConfig = types.BackendConfig

// VIPConfig represents a VIP definition.
type VIPConfig struct {
	// Address is the VIP IP address.
	Address string `yaml:"address"`
	// Port is the VIP port number.
	Port uint16 `yaml:"port"`
	// Proto is the protocol ("tcp" or "udp").
	Proto string `yaml:"proto"`
	// TargetGroup is the name of the target group to use.
	TargetGroup string `yaml:"target_group"`
	// Flags are VIP-specific flags.
	Flags uint32 `yaml:"flags"`
}

// LoadConfigFromFile loads and parses a YAML configuration file.
//
// Parameters:
//   - path: Path to the YAML configuration file.
//
// Returns the parsed configuration or an error.
func LoadConfigFromFile(path string) (*FullConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg FullConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the full configuration.
//
// Returns an error if the configuration is invalid.
func (fc *FullConfig) Validate() error {
	// Validate server config
	if fc.Server.Port <= 0 || fc.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", fc.Server.Port)
	}

	// Validate TLS if configured
	if fc.Server.TLS != nil {
		if fc.Server.TLS.CertFile != "" && fc.Server.TLS.KeyFile == "" {
			return fmt.Errorf("TLS key file is required when cert file is specified")
		}
		if fc.Server.TLS.KeyFile != "" && fc.Server.TLS.CertFile == "" {
			return fmt.Errorf("TLS cert file is required when key file is specified")
		}
	}

	// Validate LB config - main interface is required
	if fc.LB.Interfaces.Main == "" {
		return fmt.Errorf("lb.interfaces.main is required")
	}

	// Validate LB config - balancer program is required
	if fc.LB.Programs.Balancer == "" {
		return fmt.Errorf("lb.programs.balancer is required")
	}

	// Validate protocol strings in VIPs
	for i, vip := range fc.VIPs {
		proto := strings.ToLower(vip.Proto)
		if proto != "tcp" && proto != "udp" {
			return fmt.Errorf("vip[%d]: invalid protocol %q (must be 'tcp' or 'udp')", i, vip.Proto)
		}
	}

	// Validate that target groups referenced by VIPs exist
	for i, vip := range fc.VIPs {
		if vip.TargetGroup == "" {
			return fmt.Errorf("vip[%d]: target_group is required", i)
		}
		if _, ok := fc.TargetGroups[vip.TargetGroup]; !ok {
			return fmt.Errorf("vip[%d]: target group %q not found", i, vip.TargetGroup)
		}
	}

	// Check for duplicate VIPs (same address:port:proto)
	seen := make(map[string]bool)
	for i, vip := range fc.VIPs {
		key := fmt.Sprintf("%s:%d:%s", vip.Address, vip.Port, strings.ToLower(vip.Proto))
		if seen[key] {
			return fmt.Errorf("vip[%d]: duplicate VIP %s", i, key)
		}
		seen[key] = true
	}

	return nil
}

// ToServerConfig converts ServerYAMLConfig to a server Config.
//
// Returns a Config struct for the HTTP server.
func (sc *ServerYAMLConfig) ToServerConfig() *Config {
	cfg := DefaultConfig()

	cfg.Host = sc.Host
	if sc.Port > 0 {
		cfg.Port = sc.Port
	}
	if sc.ReadTimeout > 0 {
		cfg.ReadTimeout = sc.ReadTimeout
	}
	if sc.WriteTimeout > 0 {
		cfg.WriteTimeout = sc.WriteTimeout
	}
	if sc.IdleTimeout > 0 {
		cfg.IdleTimeout = sc.IdleTimeout
	}
	cfg.EnableCORS = sc.EnableCORS
	if len(sc.AllowedOrigins) > 0 {
		cfg.AllowedOrigins = sc.AllowedOrigins
	}
	if sc.EnableLogging != nil {
		cfg.EnableLogging = *sc.EnableLogging
	}
	if sc.EnableRecovery != nil {
		cfg.EnableRecovery = *sc.EnableRecovery
	}
	cfg.StaticDir = sc.StaticDir
	cfg.BPFProgDir = sc.BPFProgDir

	// Convert TLS config
	if sc.TLS != nil && sc.TLS.CertFile != "" && sc.TLS.KeyFile != "" {
		cfg.TLS = &TLSConfig{
			CertFile:     sc.TLS.CertFile,
			KeyFile:      sc.TLS.KeyFile,
			ClientCAFile: sc.TLS.ClientCAFile,
		}
		// Parse TLS min version
		switch sc.TLS.MinVersion {
		case "1.3":
			cfg.TLS.MinVersion = tls.VersionTLS13
		default:
			cfg.TLS.MinVersion = tls.VersionTLS12
		}
		// Parse client auth
		switch strings.ToLower(sc.TLS.ClientAuth) {
		case "request":
			cfg.TLS.ClientAuth = tls.RequestClientCert
		case "require":
			cfg.TLS.ClientAuth = tls.RequireAnyClientCert
		case "verify":
			cfg.TLS.ClientAuth = tls.RequireAndVerifyClientCert
		default:
			cfg.TLS.ClientAuth = tls.NoClientCert
		}
	}

	// Convert Auth config
	if sc.Auth != nil && sc.Auth.Enabled {
		cfg.Auth = &AuthConfig{
			Enabled:        sc.Auth.Enabled,
			DatabasePath:   sc.Auth.DatabasePath,
			AllowLocalhost: sc.Auth.AllowLocalhost,
			SessionTimeout: sc.Auth.SessionTimeout,
			BcryptCost:     sc.Auth.BcryptCost,
			ExemptPaths:    sc.Auth.ExemptPaths,
		}
		// Apply defaults
		if cfg.Auth.SessionTimeout <= 0 {
			cfg.Auth.SessionTimeout = 24
		}
		if cfg.Auth.BcryptCost <= 0 {
			cfg.Auth.BcryptCost = 12
		}
		if cfg.Auth.DatabasePath == "" {
			cfg.Auth.DatabasePath = "/var/lib/katran/auth.db"
		}
		// Convert bootstrap admin
		if sc.Auth.BootstrapAdmin != nil {
			cfg.Auth.BootstrapAdmin = &BootstrapAdminConfig{
				Username: sc.Auth.BootstrapAdmin.Username,
				Password: sc.Auth.BootstrapAdmin.Password,
			}
		}
	}

	return cfg
}

// ToCreateLBRequest converts LBConfig to a CreateLBRequest-like structure.
// The returned values can be used to initialize the katran Config.
//
// Parameters:
//   - bpfProgDir: Base directory for BPF programs.
//
// Returns maps of configuration values for initializing katran.
func (lc *LBConfig) ToCreateLBRequest(bpfProgDir string) map[string]interface{} {
	result := make(map[string]interface{})

	// Interfaces
	result["main_interface"] = lc.Interfaces.Main
	if lc.Interfaces.Healthcheck != "" {
		result["hc_interface"] = lc.Interfaces.Healthcheck
	} else {
		result["hc_interface"] = lc.Interfaces.Main
	}
	result["v4_tun_interface"] = lc.Interfaces.V4Tunnel
	result["v6_tun_interface"] = lc.Interfaces.V6Tunnel

	// Programs - resolve relative paths
	balancerPath := lc.Programs.Balancer
	if balancerPath != "" && !filepath.IsAbs(balancerPath) && bpfProgDir != "" {
		balancerPath = filepath.Join(bpfProgDir, balancerPath)
	}
	result["balancer_prog_path"] = balancerPath

	hcPath := lc.Programs.Healthcheck
	if hcPath != "" && !filepath.IsAbs(hcPath) && bpfProgDir != "" {
		hcPath = filepath.Join(bpfProgDir, hcPath)
	}
	result["healthchecking_prog_path"] = hcPath

	// Root map
	if lc.RootMap.Enabled != nil {
		result["use_root_map"] = *lc.RootMap.Enabled
	}
	result["root_map_path"] = lc.RootMap.Path
	if lc.RootMap.Position > 0 {
		result["root_map_pos"] = lc.RootMap.Position
	}

	// MAC addresses
	result["default_mac"] = lc.MAC.Default
	result["local_mac"] = lc.MAC.Local

	// Capacity
	if lc.Capacity.MaxVIPs > 0 {
		result["max_vips"] = lc.Capacity.MaxVIPs
	}
	if lc.Capacity.MaxReals > 0 {
		result["max_reals"] = lc.Capacity.MaxReals
	}
	if lc.Capacity.CHRingSize > 0 {
		result["ch_ring_size"] = lc.Capacity.CHRingSize
	}
	if lc.Capacity.LRUSize > 0 {
		result["lru_size"] = lc.Capacity.LRUSize
	}
	if lc.Capacity.GlobalLRUSize > 0 {
		result["global_lru_size"] = lc.Capacity.GlobalLRUSize
	}
	if lc.Capacity.MaxLPMSrc > 0 {
		result["max_lpm_src_size"] = lc.Capacity.MaxLPMSrc
	}
	if lc.Capacity.MaxDecapDst > 0 {
		result["max_decap_dst"] = lc.Capacity.MaxDecapDst
	}

	// CPU
	if len(lc.CPU.ForwardingCores) > 0 {
		result["forwarding_cores"] = lc.CPU.ForwardingCores
	}
	if len(lc.CPU.NUMANodes) > 0 {
		result["numa_nodes"] = lc.CPU.NUMANodes
	}

	// XDP
	if lc.XDP.AttachFlags > 0 {
		result["xdp_attach_flags"] = lc.XDP.AttachFlags
	}
	if lc.XDP.Priority > 0 {
		result["priority"] = lc.XDP.Priority
	}

	// Encapsulation
	result["katran_src_v4"] = lc.Encapsulation.SrcV4
	result["katran_src_v6"] = lc.Encapsulation.SrcV6

	// Features
	if lc.Features.EnableHealthcheck != nil {
		result["enable_hc"] = *lc.Features.EnableHealthcheck
	}
	if lc.Features.TunnelBasedHCEncap != nil {
		result["tunnel_based_hc_encap"] = *lc.Features.TunnelBasedHCEncap
	}
	result["flow_debug"] = lc.Features.FlowDebug
	result["enable_cid_v3"] = lc.Features.EnableCIDV3
	if lc.Features.MemlockUnlimited != nil {
		result["memlock_unlimited"] = *lc.Features.MemlockUnlimited
	}
	if lc.Features.CleanupOnShutdown != nil {
		result["cleanup_on_shutdown"] = *lc.Features.CleanupOnShutdown
	}
	result["testing"] = lc.Features.Testing

	// Hash function
	result["hash_function"] = lc.HashFunction

	return result
}

// ProtoToNumber converts a protocol string to its IP protocol number.
// This is an alias to types.ProtoToNumber.
var ProtoToNumber = types.ProtoToNumber

// HashFunctionToInt converts a hash function string to its integer value.
// This is an alias to types.HashFunctionToInt.
var HashFunctionToInt = types.HashFunctionToInt

// ParseMAC parses a MAC address string to bytes.
//
// Parameters:
//   - mac: MAC address string (e.g., "aa:bb:cc:dd:ee:ff").
//
// Returns the MAC address as bytes or an error.
func ParseMAC(mac string) ([]byte, error) {
	if mac == "" {
		return nil, nil
	}
	// Remove common separators
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	return hex.DecodeString(mac)
}

// FormatMAC formats a MAC address bytes to string.
// This is an alias to types.FormatMAC.
var FormatMAC = types.FormatMAC

// NumberToProto converts an IP protocol number to a string.
// This is an alias to types.NumberToProto.
var NumberToProto = types.NumberToProto

// IntToHashFunction converts a hash function integer to string.
// This is an alias to types.IntToHashFunction.
var IntToHashFunction = types.IntToHashFunction

// ExportConfigAsYAML exports the current running configuration as YAML.
// This constructs a FullConfig from the current server config and LB state.
//
// Parameters:
//   - serverCfg: The server configuration.
//   - katranCfg: The katran configuration (from LB).
//   - vips: List of VIPs with their backends.
//
// Returns the YAML bytes or an error.
func ExportConfigAsYAML(serverCfg *Config, katranCfg *KatranConfigExport, vips []VIPWithBackends) ([]byte, error) {
	fullCfg := buildFullConfigFromRuntime(serverCfg, katranCfg, vips)
	return yaml.Marshal(fullCfg)
}

// KatranConfigExport is an alias to types.KatranConfigExport.
type KatranConfigExport = types.KatranConfigExport

// VIPWithBackends is an alias to types.VIPWithBackends.
type VIPWithBackends = types.VIPWithBackends

// buildFullConfigFromRuntime builds a FullConfig from runtime configuration.
func buildFullConfigFromRuntime(serverCfg *Config, katranCfg *KatranConfigExport, vips []VIPWithBackends) *FullConfig {
	fc := &FullConfig{
		TargetGroups: make(map[string][]BackendConfig),
		VIPs:         make([]VIPConfig, 0, len(vips)),
	}

	// Build server config
	fc.Server = ServerYAMLConfig{
		Host:           serverCfg.Host,
		Port:           serverCfg.Port,
		ReadTimeout:    serverCfg.ReadTimeout,
		WriteTimeout:   serverCfg.WriteTimeout,
		IdleTimeout:    serverCfg.IdleTimeout,
		EnableCORS:     serverCfg.EnableCORS,
		AllowedOrigins: serverCfg.AllowedOrigins,
		StaticDir:      serverCfg.StaticDir,
		BPFProgDir:     serverCfg.BPFProgDir,
	}
	enableLogging := serverCfg.EnableLogging
	fc.Server.EnableLogging = &enableLogging
	enableRecovery := serverCfg.EnableRecovery
	fc.Server.EnableRecovery = &enableRecovery

	// Build TLS config if present
	if serverCfg.TLS != nil {
		minVersion := "1.2"
		if serverCfg.TLS.MinVersion == tls.VersionTLS13 {
			minVersion = "1.3"
		}
		clientAuth := "none"
		switch serverCfg.TLS.ClientAuth {
		case tls.RequestClientCert:
			clientAuth = "request"
		case tls.RequireAnyClientCert:
			clientAuth = "require"
		case tls.RequireAndVerifyClientCert:
			clientAuth = "verify"
		}
		fc.Server.TLS = &TLSYAMLConfig{
			CertFile:     serverCfg.TLS.CertFile,
			KeyFile:      serverCfg.TLS.KeyFile,
			MinVersion:   minVersion,
			ClientAuth:   clientAuth,
			ClientCAFile: serverCfg.TLS.ClientCAFile,
		}
	}

	// Build LB config if katranCfg is provided
	if katranCfg != nil {
		useRootMap := katranCfg.UseRootMap
		enableHC := katranCfg.EnableHC
		tunnelBasedHCEncap := katranCfg.TunnelBasedHCEncap
		memlockUnlimited := katranCfg.MemlockUnlimited
		cleanupOnShutdown := katranCfg.CleanupOnShutdown

		fc.LB = LBConfig{
			Interfaces: InterfacesConfig{
				Main:        katranCfg.MainInterface,
				Healthcheck: katranCfg.HCInterface,
				V4Tunnel:    katranCfg.V4TunInterface,
				V6Tunnel:    katranCfg.V6TunInterface,
			},
			Programs: ProgramsConfig{
				Balancer:    katranCfg.BalancerProgPath,
				Healthcheck: katranCfg.HealthcheckingProgPath,
			},
			RootMap: RootMapConfig{
				Enabled:  &useRootMap,
				Path:     katranCfg.RootMapPath,
				Position: katranCfg.RootMapPos,
			},
			MAC: MACConfig{
				Default: FormatMAC(katranCfg.DefaultMAC),
				Local:   FormatMAC(katranCfg.LocalMAC),
			},
			Capacity: CapacityConfig{
				MaxVIPs:       katranCfg.MaxVIPs,
				MaxReals:      katranCfg.MaxReals,
				CHRingSize:    katranCfg.CHRingSize,
				LRUSize:       katranCfg.LRUSize,
				GlobalLRUSize: katranCfg.GlobalLRUSize,
				MaxLPMSrc:     katranCfg.MaxLPMSrcSize,
				MaxDecapDst:   katranCfg.MaxDecapDst,
			},
			CPU: CPUConfig{
				ForwardingCores: katranCfg.ForwardingCores,
				NUMANodes:       katranCfg.NUMANodes,
			},
			XDP: XDPConfig{
				AttachFlags: katranCfg.XDPAttachFlags,
				Priority:    katranCfg.Priority,
			},
			Encapsulation: EncapsulationConfig{
				SrcV4: katranCfg.KatranSrcV4,
				SrcV6: katranCfg.KatranSrcV6,
			},
			Features: FeaturesConfig{
				EnableHealthcheck:  &enableHC,
				TunnelBasedHCEncap: &tunnelBasedHCEncap,
				FlowDebug:          katranCfg.FlowDebug,
				EnableCIDV3:        katranCfg.EnableCIDV3,
				MemlockUnlimited:   &memlockUnlimited,
				CleanupOnShutdown:  &cleanupOnShutdown,
				Testing:            katranCfg.Testing,
			},
			HashFunction: IntToHashFunction(katranCfg.HashFunc),
		}
	}

	// Build target groups and VIPs
	// Create a unique target group for each unique set of backends
	targetGroupIndex := 0
	backendHash := make(map[string]string) // hash of backends -> target group name

	for _, vip := range vips {
		// Create hash key from backends
		hashKey := hashBackends(vip.Backends)

		// Check if we already have a target group for these backends
		var targetGroupName string
		if existingName, ok := backendHash[hashKey]; ok {
			targetGroupName = existingName
		} else {
			// Create new target group
			targetGroupName = fmt.Sprintf("group-%d", targetGroupIndex)
			targetGroupIndex++
			fc.TargetGroups[targetGroupName] = vip.Backends
			backendHash[hashKey] = targetGroupName
		}

		// Add VIP
		fc.VIPs = append(fc.VIPs, VIPConfig{
			Address:     vip.Address,
			Port:        vip.Port,
			Proto:       NumberToProto(vip.Proto),
			TargetGroup: targetGroupName,
			Flags:       vip.Flags,
		})
	}

	return fc
}

// hashBackends creates a hash key from a list of backends.
func hashBackends(backends []BackendConfig) string {
	if len(backends) == 0 {
		return "empty"
	}
	var parts []string
	for _, b := range backends {
		parts = append(parts, fmt.Sprintf("%s:%d:%d", b.Address, b.Weight, b.Flags))
	}
	return strings.Join(parts, ",")
}
