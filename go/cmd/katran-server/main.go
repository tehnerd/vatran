package main

import (
	"flag"
	"log"
	"strings"

	"github.com/tehnerd/vatran/go/server"
)

func main() {
	// Parse command line flags
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

	flag.Parse()

	// Build configuration
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
