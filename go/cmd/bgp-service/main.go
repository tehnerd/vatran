package main

import (
	"flag"
	"log"

	"github.com/tehnerd/vatran/go/bgpservice"
)

func main() {
	configFile := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	cfg, err := bgpservice.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := bgpservice.New(cfg)

	log.Printf("BGP service starting...")
	if err := srv.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("BGP service stopped")
}
