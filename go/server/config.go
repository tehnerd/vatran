package server

import (
	"crypto/tls"
	"fmt"
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

// Config contains the server configuration.
type Config struct {
	// Host is the host to bind to (default: "").
	Host string
	// Port is the port to listen on (default: 8080).
	Port int
	// TLS contains TLS configuration. If nil, HTTP is used.
	TLS *TLSConfig
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
