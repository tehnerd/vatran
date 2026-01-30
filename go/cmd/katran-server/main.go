package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/tehnerd/vatran/go/server"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "", "Path to YAML config file (overrides other flags if provided)")
	host := flag.String("host", "", "Host to bind to")
	port := flag.Int("port", 8080, "Port to listen on")
	tlsCert := flag.String("tls-cert", "", "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "Path to TLS private key file")
	tlsClientCA := flag.String("tls-client-ca", "", "Path to client CA file for mTLS")
	enableCORS := flag.Bool("cors", false, "Enable CORS")
	corsOrigins := flag.String("cors-origins", "*", "Comma-separated list of allowed CORS origins")
	readTimeout := flag.Int("read-timeout", 30, "Read timeout in seconds")
	writeTimeout := flag.Int("write-timeout", 30, "Write timeout in seconds")
	idleTimeout := flag.Int("idle-timeout", 120, "Idle timeout in seconds")
	disableLogging := flag.Bool("no-logging", false, "Disable request logging")
	disableRecovery := flag.Bool("no-recovery", false, "Disable panic recovery")
	staticDir := flag.String("static-dir", "", "Path to static files directory for SPA (e.g., ui/dist)")
	bpfProgDir := flag.String("bpf-prog-dir", "", "Base directory for BPF program files")

	flag.Parse()

	// If config file is provided, use it
	if *configFile != "" {
		if _, err := os.Stat(*configFile); err != nil {
			if os.IsNotExist(err) {
				log.Printf("config file do not exists")
			} else {
				log.Fatalf("Failed to stat config file: %v", err)
			}
		} else {
			runWithConfigFile(*configFile)
			return
		}
	}

	// Build configuration from flags
	config := server.DefaultConfig()
	config.Host = *host
	config.Port = *port
	config.ReadTimeout = *readTimeout
	config.WriteTimeout = *writeTimeout
	config.IdleTimeout = *idleTimeout
	config.EnableLogging = !*disableLogging
	config.EnableRecovery = !*disableRecovery
	config.EnableCORS = *enableCORS

	if *corsOrigins != "" {
		config.AllowedOrigins = strings.Split(*corsOrigins, ",")
	}

	config.StaticDir = *staticDir
	config.BPFProgDir = *bpfProgDir

	// Configure TLS if cert and key provided
	if *tlsCert != "" && *tlsKey != "" {
		config.TLS = &server.TLSConfig{
			CertFile:     *tlsCert,
			KeyFile:      *tlsKey,
			ClientCAFile: *tlsClientCA,
		}
	}

	// Create and start server
	srv := server.New(config)

	log.Printf("Katran REST API Server starting...")
	if err := srv.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

// runWithConfigFile loads configuration from a YAML file and starts the server.
func runWithConfigFile(configPath string) {
	log.Printf("Loading configuration from %s...", configPath)

	// Load config from file
	cfg, err := server.LoadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Convert to server config
	serverConfig := cfg.Server.ToServerConfig()

	// Create server
	srv := server.New(serverConfig)

	// Initialize LB from config (creates LB, loads/attaches BPF, configures VIPs)
	log.Printf("Initializing load balancer from config...")
	if err := srv.InitFromConfig(cfg); err != nil {
		log.Fatalf("Failed to initialize from config: %v", err)
	}

	// Start server
	log.Printf("Katran REST API Server starting on %s...", serverConfig.Addr())
	if err := srv.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
