package bgpservice

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is the top-level BGP service, tying together all components.
type Server struct {
	config   *Config
	state    *State
	bgp      *BGPSpeaker
	handlers *Handlers
	httpSrv  *http.Server
}

// New creates a new Server with all components wired together.
//
// Parameters:
//   - config: The service configuration.
//
// Returns a new Server instance ready to be started.
func New(config *Config) *Server {
	state := NewState()
	bgp := NewBGPSpeaker(&config.BGP)
	bgpHandlers := NewHandlers(state, bgp, &config.BGP)

	mux := http.NewServeMux()
	RegisterRoutes(mux, bgpHandlers)

	httpSrv := &http.Server{
		Addr:         config.Server.Addr(),
		Handler:      mux,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		config:   config,
		state:    state,
		bgp:      bgp,
		handlers: bgpHandlers,
		httpSrv:  httpSrv,
	}
}

// Start launches the BGP speaker and HTTP server. It blocks until the server stops.
//
// Returns an error if the server fails to start (not including graceful shutdown).
func (s *Server) Start() error {
	// Start BGP speaker first
	if err := s.bgp.Start(); err != nil {
		return err
	}

	log.Printf("BGP service listening on %s", s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully shuts down the HTTP server and BGP speaker.
//
// Parameters:
//   - ctx: Context with deadline for graceful shutdown.
//
// Returns an error if shutdown fails.
func (s *Server) Stop(ctx context.Context) error {
	err := s.httpSrv.Shutdown(ctx)
	s.bgp.Stop()
	return err
}

// RunWithGracefulShutdown starts the server and handles SIGINT/SIGTERM for graceful shutdown.
//
// Returns an error if the server fails.
func (s *Server) RunWithGracefulShutdown() error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		log.Printf("Received signal %v, shutting down...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.Stop(ctx)
	}
}
