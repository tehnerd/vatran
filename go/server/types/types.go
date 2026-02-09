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
	Address  string
	Port     uint16
	Proto    uint8
	Flags    uint32
	Backends []BackendConfig
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
