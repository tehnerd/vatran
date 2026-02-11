package main

import (
	"flag"
	"log"

	"github.com/tehnerd/vatran/go/hcservice"
)

func main() {
	configFile := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	cfg, err := hcservice.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := hcservice.New(cfg)

	log.Printf("HC service starting...")
	if err := srv.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("HC service stopped")
}
