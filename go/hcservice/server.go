package hcservice

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is the top-level healthcheck service, tying together all components.
type Server struct {
	config    *Config
	state     *State
	somarks   *SomarkAllocator
	katran    *KatranClient
	dialer    *SomarkDialer
	scheduler *Scheduler
	handlers  *Handlers
	httpSrv   *http.Server
}

// New creates a new Server with all components wired together.
//
// Parameters:
//   - config: The service configuration.
//
// Returns a new Server instance ready to be started.
func New(config *Config) *Server {
	state := NewState()
	somarks := NewSomarkAllocator(config.Somark.BaseSomark, config.Somark.MaxReals)
	katranClient := NewKatranClient(config.Katran.ServerURL, config.Katran.Timeout)
	dialer := NewSomarkDialer(time.Duration(config.Katran.Timeout) * time.Second)

	// Build checkers map
	checkers := map[string]Checker{
		"tcp":   NewTCPChecker(dialer),
		"http":  NewHTTPChecker(dialer),
		"https": NewHTTPSChecker(dialer),
		"dummy": NewDummyChecker(),
	}

	scheduler := NewScheduler(state, somarks, checkers, config.Scheduler)
	handlers := NewHandlers(state, somarks, katranClient, scheduler)

	mux := http.NewServeMux()
	RegisterRoutes(mux, handlers)

	httpSrv := &http.Server{
		Addr:         config.Server.Addr(),
		Handler:      mux,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		config:    config,
		state:     state,
		somarks:   somarks,
		katran:    katranClient,
		dialer:    dialer,
		scheduler: scheduler,
		handlers:  handlers,
		httpSrv:   httpSrv,
	}
}

// Start launches the scheduler and HTTP server. It blocks until the server stops.
//
// Returns an error if the server fails to start (not including graceful shutdown).
func (s *Server) Start() error {
	s.scheduler.Start()
	log.Printf("HC service listening on %s", s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully shuts down the HTTP server and scheduler.
//
// Parameters:
//   - ctx: Context with deadline for graceful shutdown.
//
// Returns an error if shutdown fails.
func (s *Server) Stop(ctx context.Context) error {
	s.scheduler.Stop()
	return s.httpSrv.Shutdown(ctx)
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
