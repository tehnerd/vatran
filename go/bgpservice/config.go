package bgpservice

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig contains HTTP server configuration for the BGP service.
type ServerConfig struct {
	// Host is the address to bind to (empty string binds to all interfaces).
	Host string `yaml:"host"`
	// Port is the port to listen on.
	Port int `yaml:"port"`
	// ReadTimeout is the HTTP read timeout in seconds.
	ReadTimeout int `yaml:"read_timeout"`
	// WriteTimeout is the HTTP write timeout in seconds.
	WriteTimeout int `yaml:"write_timeout"`
}

// BGPConfig contains BGP speaker configuration.
type BGPConfig struct {
	// ASN is the local autonomous system number (required).
	ASN uint32 `yaml:"asn"`
	// RouterID is the BGP router identifier (required).
	RouterID string `yaml:"router_id"`
	// ListenPort is the port for BGP peer connections.
	ListenPort int `yaml:"listen_port"`
	// LocalPref is the default local preference for route announcements.
	LocalPref uint32 `yaml:"local_pref"`
	// Communities is the default list of BGP communities for announcements.
	Communities []string `yaml:"communities"`
	// Peers is the list of initial BGP peers to configure.
	Peers []PeerConfig `yaml:"peers"`
}

// PeerConfig contains configuration for a single BGP peer.
type PeerConfig struct {
	// Address is the peer IP address (required).
	Address string `yaml:"address"`
	// ASN is the peer autonomous system number (required).
	ASN uint32 `yaml:"asn"`
	// HoldTime is the BGP hold timer in seconds.
	HoldTime int `yaml:"hold_time"`
	// Keepalive is the BGP keepalive interval in seconds.
	Keepalive int `yaml:"keepalive"`
}

// Config is the top-level configuration for the BGP service.
type Config struct {
	// Server contains HTTP server settings.
	Server ServerConfig `yaml:"server"`
	// BGP contains BGP speaker settings.
	BGP BGPConfig `yaml:"bgp"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "127.0.0.1",
			Port:         9100,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		BGP: BGPConfig{
			ASN:        0,
			RouterID:   "",
			ListenPort: 179,
			LocalPref:  100,
		},
	}
}

// LoadConfig reads a YAML config file and returns a Config with defaults applied.
//
// Parameters:
//   - path: Path to the YAML configuration file. If empty, returns DefaultConfig.
//
// Returns the parsed Config or an error if reading/parsing fails.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = "127.0.0.1"
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks the configuration for errors.
//
// Returns an error if any required field is missing or invalid.
func (c *Config) Validate() error {
	if c.BGP.ASN == 0 {
		return fmt.Errorf("bgp.asn is required")
	}
	if c.BGP.RouterID == "" {
		return fmt.Errorf("bgp.router_id is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	for i, peer := range c.BGP.Peers {
		if peer.Address == "" {
			return fmt.Errorf("bgp.peers[%d].address is required", i)
		}
		if peer.ASN == 0 {
			return fmt.Errorf("bgp.peers[%d].asn is required", i)
		}
	}
	return nil
}

// Addr returns the server listen address as "host:port".
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
