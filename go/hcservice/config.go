package hcservice

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig contains HTTP server configuration.
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

// KatranConfig contains the katran server connection settings.
type KatranConfig struct {
	// ServerURL is the base URL of the katran server (e.g., "http://localhost:8080").
	ServerURL string `yaml:"server_url"`
	// Timeout is the HTTP client timeout in seconds.
	Timeout int `yaml:"timeout"`
}

// SomarkConfig contains somark allocator settings.
type SomarkConfig struct {
	// BaseSomark is the starting somark value for allocation.
	BaseSomark uint32 `yaml:"base_somark"`
	// MaxReals is the maximum number of unique reals that can be tracked.
	MaxReals uint32 `yaml:"max_reals"`
}

// SchedulerConfig contains scheduler settings.
type SchedulerConfig struct {
	// SpreadIntervalMs is the window in milliseconds over which to stagger new health checks.
	SpreadIntervalMs int `yaml:"spread_interval_ms"`
	// WorkerCount is the maximum number of concurrent health check workers.
	WorkerCount int `yaml:"worker_count"`
	// TickIntervalMs is the scheduler sweep granularity in milliseconds.
	TickIntervalMs int `yaml:"tick_interval_ms"`
}

// Config is the top-level configuration for the healthcheck service.
type Config struct {
	// Server contains HTTP server settings.
	Server ServerConfig `yaml:"server"`
	// Katran contains katran server connection settings.
	Katran KatranConfig `yaml:"katran"`
	// Somark contains somark allocator settings.
	Somark SomarkConfig `yaml:"somark"`
	// Scheduler contains scheduler settings.
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "127.0.0.1",
			Port:         9000,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Katran: KatranConfig{
			ServerURL: "http://localhost:8080",
			Timeout:   10,
		},
		Somark: SomarkConfig{
			BaseSomark: 10000,
			MaxReals:   4096,
		},
		Scheduler: SchedulerConfig{
			SpreadIntervalMs: 3000,
			WorkerCount:      64,
			TickIntervalMs:   100,
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
	if c.Katran.ServerURL == "" {
		return fmt.Errorf("katran.server_url is required")
	}
	if c.Somark.BaseSomark == 0 {
		return fmt.Errorf("somark.base_somark must be > 0")
	}
	if c.Somark.MaxReals == 0 {
		return fmt.Errorf("somark.max_reals must be > 0")
	}
	if c.Scheduler.WorkerCount <= 0 {
		return fmt.Errorf("scheduler.worker_count must be > 0")
	}
	if c.Scheduler.TickIntervalMs <= 0 {
		return fmt.Errorf("scheduler.tick_interval_ms must be > 0")
	}
	if c.Scheduler.SpreadIntervalMs <= 0 {
		return fmt.Errorf("scheduler.spread_interval_ms must be > 0")
	}
	return nil
}

// Addr returns the server listen address as "host:port".
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
