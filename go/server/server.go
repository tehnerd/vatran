package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/lb"
	"github.com/tehnerd/vatran/go/server/middleware"
)

// Server represents the HTTP server for the Katran API.
type Server struct {
	config        *Config
	httpServer    *http.Server
	mux           *http.ServeMux
	authenticator middleware.Authenticator
}

// New creates a new Server with the provided configuration.
//
// Parameters:
//   - config: Server configuration. If nil, default config is used.
//
// Returns a new Server instance.
func New(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	mux := http.NewServeMux()

	return &Server{
		config:        config,
		mux:           mux,
		authenticator: middleware.NewNoOpAuthenticator(),
	}
}

// SetAuthenticator sets the authenticator for the server.
//
// Parameters:
//   - auth: The authenticator to use.
func (s *Server) SetAuthenticator(auth middleware.Authenticator) {
	s.authenticator = auth
}

// Start starts the HTTP server.
//
// This method blocks until the server is stopped. Use StartAsync for non-blocking start.
//
// Returns an error if the server fails to start.
func (s *Server) Start() error {
	if err := s.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Register routes
	RegisterRoutes(s.mux, s.config)

	// Build handler chain with middleware
	var handler http.Handler = s.mux

	// Apply middleware in reverse order (last applied runs first)
	handler = middleware.Auth(s.authenticator)(handler)

	if s.config.EnableCORS {
		corsConfig := middleware.DefaultCORSConfig()
		corsConfig.AllowedOrigins = s.config.AllowedOrigins
		handler = middleware.CORS(corsConfig)(handler)
	}

	if s.config.EnableLogging {
		handler = middleware.Logging()(handler)
	}

	if s.config.EnableRecovery {
		handler = middleware.Recovery()(handler)
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.config.Addr(),
		Handler:      handler,
		ReadTimeout:  time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.IdleTimeout) * time.Second,
	}

	// Configure TLS if enabled
	if s.config.IsTLS() {
		tlsConfig, err := s.buildTLSConfig()
		if err != nil {
			return fmt.Errorf("failed to configure TLS: %w", err)
		}
		s.httpServer.TLSConfig = tlsConfig
	}

	// Start server
	if s.config.IsTLS() {
		log.Printf("Starting HTTPS server on %s", s.config.Addr())
		return s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
	}

	log.Printf("Starting HTTP server on %s", s.config.Addr())
	return s.httpServer.ListenAndServe()
}

// StartAsync starts the HTTP server asynchronously.
//
// Returns a channel that receives an error when the server stops.
func (s *Server) StartAsync() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Start()
	}()
	return errChan
}

// Stop gracefully stops the HTTP server.
//
// Parameters:
//   - ctx: Context for the shutdown. Can be used for timeout.
//
// Returns an error if shutdown fails.
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// RunWithGracefulShutdown starts the server and handles graceful shutdown on SIGINT/SIGTERM.
//
// Returns an error if the server fails.
func (s *Server) RunWithGracefulShutdown() error {
	// Start server asynchronously
	errChan := s.StartAsync()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-quit:
		log.Println("Received shutdown signal")
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.Stop(ctx)
}

// InitFromConfig initializes the load balancer from a FullConfig.
// This creates the LB, loads and attaches BPF programs, and configures all VIPs with their backends.
//
// Parameters:
//   - cfg: The full configuration loaded from YAML.
//
// Returns an error if initialization fails.
func (s *Server) InitFromConfig(cfg *FullConfig) error {
	manager := lb.GetManager()

	// Build katran config from YAML config
	katranCfg := s.buildKatranConfig(cfg)

	// Create LB
	log.Println("Creating load balancer...")
	if err := manager.Create(katranCfg); err != nil {
		return fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Load BPF programs
	log.Println("Loading BPF programs...")
	if err := manager.LoadBPFProgs(); err != nil {
		return fmt.Errorf("failed to load BPF programs: %w", err)
	}

	// Attach BPF programs
	log.Println("Attaching BPF programs...")
	if err := manager.AttachBPFProgs(); err != nil {
		return fmt.Errorf("failed to attach BPF programs: %w", err)
	}

	// Get LB instance for VIP/Real operations
	lbInstance, ok := manager.Get()
	if !ok {
		return fmt.Errorf("load balancer not initialized after creation")
	}

	// Create VIPs and add backends from target groups
	for i, vipCfg := range cfg.VIPs {
		vip := katran.VIPKey{
			Address: vipCfg.Address,
			Port:    vipCfg.Port,
			Proto:   ProtoToNumber(vipCfg.Proto),
		}

		log.Printf("Adding VIP %s:%d/%s...", vipCfg.Address, vipCfg.Port, vipCfg.Proto)
		if err := lbInstance.AddVIP(vip, vipCfg.Flags); err != nil {
			return fmt.Errorf("failed to add VIP[%d] %s:%d: %w", i, vipCfg.Address, vipCfg.Port, err)
		}

		// Add backends from target group
		backends := cfg.TargetGroups[vipCfg.TargetGroup]
		if len(backends) > 0 {
			reals := make([]katran.Real, len(backends))
			for j, backend := range backends {
				reals[j] = katran.Real{
					Address: backend.Address,
					Weight:  backend.Weight,
					Flags:   backend.Flags,
				}
			}

			log.Printf("  Adding %d backends from target group %q...", len(reals), vipCfg.TargetGroup)
			if err := lbInstance.ModifyRealsForVIP(katran.ActionAdd, reals, vip); err != nil {
				return fmt.Errorf("failed to add backends for VIP[%d]: %w", i, err)
			}
		}
	}

	log.Printf("Initialization complete: %d VIPs configured", len(cfg.VIPs))
	return nil
}

// buildKatranConfig converts FullConfig to a katran.Config.
func (s *Server) buildKatranConfig(cfg *FullConfig) *katran.Config {
	katranCfg := katran.NewConfig()

	lc := &cfg.LB
	bpfProgDir := cfg.Server.BPFProgDir

	// Interfaces
	katranCfg.MainInterface = lc.Interfaces.Main
	if lc.Interfaces.Healthcheck != "" {
		katranCfg.HCInterface = lc.Interfaces.Healthcheck
	} else {
		katranCfg.HCInterface = lc.Interfaces.Main
	}
	katranCfg.V4TunInterface = lc.Interfaces.V4Tunnel
	katranCfg.V6TunInterface = lc.Interfaces.V6Tunnel

	// Programs - resolve relative paths
	katranCfg.BalancerProgPath = lc.Programs.Balancer
	if katranCfg.BalancerProgPath != "" && !filepath.IsAbs(katranCfg.BalancerProgPath) && bpfProgDir != "" {
		katranCfg.BalancerProgPath = filepath.Join(bpfProgDir, katranCfg.BalancerProgPath)
	}
	katranCfg.HealthcheckingProgPath = lc.Programs.Healthcheck
	if katranCfg.HealthcheckingProgPath != "" && !filepath.IsAbs(katranCfg.HealthcheckingProgPath) && bpfProgDir != "" {
		katranCfg.HealthcheckingProgPath = filepath.Join(bpfProgDir, katranCfg.HealthcheckingProgPath)
	}

	// Root map
	if lc.RootMap.Enabled != nil {
		katranCfg.UseRootMap = *lc.RootMap.Enabled
	}
	katranCfg.RootMapPath = lc.RootMap.Path
	if lc.RootMap.Position > 0 {
		katranCfg.RootMapPos = lc.RootMap.Position
	}

	// MAC addresses
	if mac, err := ParseMAC(lc.MAC.Default); err == nil && len(mac) == 6 {
		katranCfg.DefaultMAC = mac
	}
	if mac, err := ParseMAC(lc.MAC.Local); err == nil && len(mac) == 6 {
		katranCfg.LocalMAC = mac
	}

	// Capacity
	if lc.Capacity.MaxVIPs > 0 {
		katranCfg.MaxVIPs = lc.Capacity.MaxVIPs
	}
	if lc.Capacity.MaxReals > 0 {
		katranCfg.MaxReals = lc.Capacity.MaxReals
	}
	if lc.Capacity.CHRingSize > 0 {
		katranCfg.CHRingSize = lc.Capacity.CHRingSize
	}
	if lc.Capacity.LRUSize > 0 {
		katranCfg.LRUSize = lc.Capacity.LRUSize
	}
	if lc.Capacity.GlobalLRUSize > 0 {
		katranCfg.GlobalLRUSize = lc.Capacity.GlobalLRUSize
	}
	if lc.Capacity.MaxLPMSrc > 0 {
		katranCfg.MaxLPMSrcSize = lc.Capacity.MaxLPMSrc
	}
	if lc.Capacity.MaxDecapDst > 0 {
		katranCfg.MaxDecapDst = lc.Capacity.MaxDecapDst
	}

	// CPU
	if len(lc.CPU.ForwardingCores) > 0 {
		katranCfg.ForwardingCores = lc.CPU.ForwardingCores
	}
	if len(lc.CPU.NUMANodes) > 0 {
		katranCfg.NUMANodes = lc.CPU.NUMANodes
	}

	// XDP
	if lc.XDP.AttachFlags > 0 {
		katranCfg.XDPAttachFlags = lc.XDP.AttachFlags
	}
	if lc.XDP.Priority > 0 {
		katranCfg.Priority = lc.XDP.Priority
	}

	// Encapsulation
	katranCfg.KatranSrcV4 = lc.Encapsulation.SrcV4
	katranCfg.KatranSrcV6 = lc.Encapsulation.SrcV6

	// Features
	if lc.Features.EnableHealthcheck != nil {
		katranCfg.EnableHC = *lc.Features.EnableHealthcheck
	}
	if lc.Features.TunnelBasedHCEncap != nil {
		katranCfg.TunnelBasedHCEncap = *lc.Features.TunnelBasedHCEncap
	}
	katranCfg.FlowDebug = lc.Features.FlowDebug
	katranCfg.EnableCIDV3 = lc.Features.EnableCIDV3
	if lc.Features.MemlockUnlimited != nil {
		katranCfg.MemlockUnlimited = *lc.Features.MemlockUnlimited
	}
	if lc.Features.CleanupOnShutdown != nil {
		katranCfg.CleanupOnShutdown = *lc.Features.CleanupOnShutdown
	}
	katranCfg.Testing = lc.Features.Testing

	// Hash function
	katranCfg.HashFunc = katran.HashFunction(HashFunctionToInt(lc.HashFunction))

	return katranCfg
}

// buildTLSConfig builds the TLS configuration.
func (s *Server) buildTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if s.config.TLS.MinVersion != 0 {
		tlsConfig.MinVersion = s.config.TLS.MinVersion
	}

	tlsConfig.ClientAuth = s.config.TLS.ClientAuth

	// Load client CA if specified (for mTLS)
	if s.config.TLS.ClientCAFile != "" {
		caCert, err := os.ReadFile(s.config.TLS.ClientCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read client CA file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse client CA certificate")
		}

		tlsConfig.ClientCAs = caCertPool
	}

	return tlsConfig, nil
}
