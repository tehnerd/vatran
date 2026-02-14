package types

import (
	"crypto/tls"
	"fmt"
	"strings"
)

// TLSConfigInfo provides TLS configuration information for export.
type TLSConfigInfo struct {
	CertFile     string
	KeyFile      string
	MinVersion   uint16
	ClientAuth   tls.ClientAuthType
	ClientCAFile string
}

// AuthConfigInfo provides authentication configuration information for export.
type AuthConfigInfo struct {
	Enabled        bool
	DatabasePath   string
	AllowLocalhost bool
	SessionTimeout int
	BcryptCost     int
	ExemptPaths    []string
}

// ServerConfigProvider provides server configuration for handlers.
// This interface breaks the import cycle between server and handlers packages.
type ServerConfigProvider interface {
	GetHost() string
	GetPort() int
	GetReadTimeout() int
	GetWriteTimeout() int
	GetIdleTimeout() int
	IsEnableCORS() bool
	GetAllowedOrigins() []string
	IsEnableLogging() bool
	IsEnableRecovery() bool
	GetStaticDir() string
	GetBPFProgDir() string
	GetTLS() *TLSConfigInfo
	GetAuthInfo() *AuthConfigInfo
}

// ConfigExporter exports configuration to YAML format.
// This interface is implemented by the server package.
type ConfigExporter interface {
	ExportAsYAML(katranCfg *KatranConfigExport, vips []VIPWithBackends) ([]byte, error)
}

// ConfigExporterFunc is a function type that implements ConfigExporter.
type ConfigExporterFunc func(katranCfg *KatranConfigExport, vips []VIPWithBackends) ([]byte, error)

// ExportAsYAML implements ConfigExporter.
func (f ConfigExporterFunc) ExportAsYAML(katranCfg *KatranConfigExport, vips []VIPWithBackends) ([]byte, error) {
	return f(katranCfg, vips)
}

// BackendConfig represents a backend server in a target group.
type BackendConfig struct {
	// Address is the IP address of the backend.
	Address string `yaml:"address" json:"address"`
	// Weight is the weight for consistent hashing.
	Weight uint32 `yaml:"weight" json:"weight"`
	// Flags are backend-specific flags.
	Flags uint8 `yaml:"flags" json:"flags"`
	// Healthy indicates whether the backend is healthy and receiving traffic.
	Healthy bool `yaml:"healthy,omitempty" json:"healthy"`
}

// VIPWithBackends contains a VIP and its backends.
type VIPWithBackends struct {
	Address     string
	Port        uint16
	Proto       uint8
	Flags       uint32
	Backends    []BackendConfig
	Healthcheck *HealthcheckConfig
}

// HealthcheckHTTPConfig contains HTTP-specific healthcheck configuration.
type HealthcheckHTTPConfig struct {
	// Path is the HTTP path to check (e.g., "/healthz").
	Path string `yaml:"path" json:"path"`
	// ExpectedStatus is the expected HTTP status code (default: 200).
	ExpectedStatus int `yaml:"expected_status,omitempty" json:"expected_status,omitempty"`
	// Host is the optional Host header value.
	Host string `yaml:"host,omitempty" json:"host,omitempty"`
}

// HealthcheckHTTPSConfig contains HTTPS-specific healthcheck configuration.
type HealthcheckHTTPSConfig struct {
	// Path is the HTTP path to check (e.g., "/healthz").
	Path string `yaml:"path" json:"path"`
	// ExpectedStatus is the expected HTTP status code (default: 200).
	ExpectedStatus int `yaml:"expected_status,omitempty" json:"expected_status,omitempty"`
	// Host is the optional Host header value.
	Host string `yaml:"host,omitempty" json:"host,omitempty"`
	// SkipTLSVerify skips TLS certificate verification (default: false).
	SkipTLSVerify bool `yaml:"skip_tls_verify,omitempty" json:"skip_tls_verify,omitempty"`
}

// HealthcheckConfig contains per-VIP healthcheck configuration.
type HealthcheckConfig struct {
	// Type is the healthcheck type ("http", "https", "tcp", "dummy").
	Type string `yaml:"type" json:"type"`
	// Port is the port to check on each real server (default: VIP port).
	Port int `yaml:"port,omitempty" json:"port,omitempty"`
	// HTTP contains HTTP-specific healthcheck settings.
	HTTP *HealthcheckHTTPConfig `yaml:"http,omitempty" json:"http,omitempty"`
	// HTTPS contains HTTPS-specific healthcheck settings.
	HTTPS *HealthcheckHTTPSConfig `yaml:"https,omitempty" json:"https,omitempty"`
	// IntervalMs is the check interval in milliseconds (default: 5000).
	IntervalMs int `yaml:"interval_ms,omitempty" json:"interval_ms,omitempty"`
	// TimeoutMs is the check timeout in milliseconds (default: 2000).
	TimeoutMs int `yaml:"timeout_ms,omitempty" json:"timeout_ms,omitempty"`
	// HealthyThreshold is the number of consecutive successes before marking healthy (default: 3).
	HealthyThreshold int `yaml:"healthy_threshold,omitempty" json:"healthy_threshold,omitempty"`
	// UnhealthyThreshold is the number of consecutive failures before marking unhealthy (default: 3).
	UnhealthyThreshold int `yaml:"unhealthy_threshold,omitempty" json:"unhealthy_threshold,omitempty"`
}

// ApplyDefaults fills in default values for zero-valued fields.
func (hc *HealthcheckConfig) ApplyDefaults() {
	if hc.IntervalMs <= 0 {
		hc.IntervalMs = 5000
	}
	if hc.TimeoutMs <= 0 {
		hc.TimeoutMs = 2000
	}
	if hc.HealthyThreshold <= 0 {
		hc.HealthyThreshold = 3
	}
	if hc.UnhealthyThreshold <= 0 {
		hc.UnhealthyThreshold = 3
	}
	if hc.HTTP != nil && hc.HTTP.ExpectedStatus <= 0 {
		hc.HTTP.ExpectedStatus = 200
	}
	if hc.HTTPS != nil && hc.HTTPS.ExpectedStatus <= 0 {
		hc.HTTPS.ExpectedStatus = 200
	}
}

// Validate checks the healthcheck configuration for errors.
//
// Returns an error if the configuration is invalid.
func (hc *HealthcheckConfig) Validate() error {
	validTypes := map[string]bool{"http": true, "https": true, "tcp": true, "dummy": true}
	if !validTypes[hc.Type] {
		return fmt.Errorf("invalid healthcheck type %q (must be http, https, tcp, or dummy)", hc.Type)
	}
	if hc.Type == "dummy" {
		return nil
	}
	if hc.Type == "http" && hc.HTTP == nil {
		return fmt.Errorf("http healthcheck requires 'http' configuration")
	}
	if hc.Type == "http" && hc.HTTP.Path == "" {
		return fmt.Errorf("http healthcheck requires 'path' in http configuration")
	}
	if hc.Type == "https" && hc.HTTPS == nil {
		return fmt.Errorf("https healthcheck requires 'https' configuration")
	}
	if hc.Type == "https" && hc.HTTPS.Path == "" {
		return fmt.Errorf("https healthcheck requires 'path' in https configuration")
	}
	if hc.TimeoutMs >= hc.IntervalMs {
		return fmt.Errorf("timeout_ms (%d) must be less than interval_ms (%d)", hc.TimeoutMs, hc.IntervalMs)
	}
	if hc.Port < 0 || hc.Port > 65535 {
		return fmt.Errorf("invalid healthcheck port: %d", hc.Port)
	}
	return nil
}

// HCRealHealth represents the health state of a single real from the HC service.
type HCRealHealth struct {
	// Address is the real server IP address.
	Address string `json:"address"`
	// Healthy indicates whether the real is healthy.
	Healthy bool `json:"healthy"`
	// LastCheckTime is the timestamp of the last health check.
	LastCheckTime string `json:"last_check_time,omitempty"`
	// LastStatusChange is the timestamp of the last status change.
	LastStatusChange string `json:"last_status_change,omitempty"`
	// ConsecutiveFailures is the number of consecutive failed checks.
	ConsecutiveFailures int `json:"consecutive_failures"`
}

// HCVIPHealthResponse represents the health response for a single VIP from the HC service.
type HCVIPHealthResponse struct {
	// VIP identifies the virtual IP.
	VIP HCVIPKey `json:"vip"`
	// Reals contains health states for all reals of this VIP.
	Reals []HCRealHealth `json:"reals"`
}

// HCVIPKey is the VIP identifier used in HC service responses.
type HCVIPKey struct {
	// Address is the IP address of the VIP.
	Address string `json:"address"`
	// Port is the port number.
	Port uint16 `json:"port"`
	// Proto is the IP protocol number.
	Proto uint8 `json:"proto"`
}

// KatranConfigExport contains the exported katran configuration.
type KatranConfigExport struct {
	MainInterface          string
	HCInterface            string
	V4TunInterface         string
	V6TunInterface         string
	BalancerProgPath       string
	HealthcheckingProgPath string
	RootMapPath            string
	RootMapPos             uint32
	UseRootMap             bool
	DefaultMAC             []byte
	LocalMAC               []byte
	MaxVIPs                uint32
	MaxReals               uint32
	CHRingSize             uint32
	LRUSize                uint64
	GlobalLRUSize          uint32
	MaxLPMSrcSize          uint32
	MaxDecapDst            uint32
	ForwardingCores        []int32
	NUMANodes              []int32
	XDPAttachFlags         uint32
	Priority               uint32
	KatranSrcV4            string
	KatranSrcV6            string
	EnableHC               bool
	TunnelBasedHCEncap     bool
	FlowDebug              bool
	EnableCIDV3            bool
	MemlockUnlimited       bool
	CleanupOnShutdown      bool
	Testing               bool
	HashFunc              int
	HealthcheckerEndpoint string
	BGPEndpoint           string
	BGPMinHealthyReals    int
}

// NumberToProto converts an IP protocol number to a string.
//
// Parameters:
//   - proto: Protocol number (6 for TCP, 17 for UDP).
//
// Returns "tcp" or "udp".
func NumberToProto(proto uint8) string {
	switch proto {
	case 6:
		return "tcp"
	case 17:
		return "udp"
	default:
		return fmt.Sprintf("%d", proto)
	}
}

// ProtoToNumber converts a protocol string to its IP protocol number.
//
// Parameters:
//   - proto: Protocol string ("tcp" or "udp").
//
// Returns the protocol number (6 for TCP, 17 for UDP).
func ProtoToNumber(proto string) uint8 {
	switch strings.ToLower(proto) {
	case "tcp":
		return 6
	case "udp":
		return 17
	default:
		return 0
	}
}

// FormatMAC formats a MAC address bytes to string.
//
// Parameters:
//   - mac: MAC address as bytes (6 bytes).
//
// Returns the MAC address as string (e.g., "aa:bb:cc:dd:ee:ff").
func FormatMAC(mac []byte) string {
	if len(mac) != 6 {
		return ""
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

// IntToHashFunction converts a hash function integer to string.
//
// Parameters:
//   - hashFunc: Hash function integer (0 for maglev, 1 for maglev_v2).
//
// Returns "maglev" or "maglev_v2".
func IntToHashFunction(hashFunc int) string {
	switch hashFunc {
	case 1:
		return "maglev_v2"
	default:
		return "maglev"
	}
}

// HashFunctionToInt converts a hash function string to its integer value.
//
// Parameters:
//   - hashFunc: Hash function string ("maglev" or "maglev_v2").
//
// Returns 0 for maglev, 1 for maglev_v2.
func HashFunctionToInt(hashFunc string) int {
	switch strings.ToLower(hashFunc) {
	case "maglev_v2":
		return 1
	default:
		return 0
	}
}
