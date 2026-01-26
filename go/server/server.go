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
	"syscall"
	"time"

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
	RegisterRoutes(s.mux)

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
