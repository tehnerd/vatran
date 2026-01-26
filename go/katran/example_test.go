//go:build integration
// +build integration

// Package katran_test contains example tests demonstrating katran library usage.
// These examples require the C library to be built and appropriate privileges.
// Run with: go test -tags=integration -v -run "^Example"
package katran_test

import (
	"fmt"
	"log"

	"github.com/tehnerd/vatran/go/katran"
)

// Example_basicUsage demonstrates basic LoadBalancer creation and VIP management.
func Example_basicUsage() {
	// Create a configuration with defaults.
	cfg := katran.NewConfig()
	cfg.Testing = true          // Enable testing mode (no actual BPF operations)
	cfg.MainInterface = "lo"    // Use loopback for testing
	cfg.EnableHC = false        // Disable healthcheck for this example

	// Create the load balancer instance.
	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	// Define a VIP (Virtual IP).
	vip := katran.VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6, // TCP
	}

	// Add the VIP to the load balancer.
	if err := lb.AddVIP(vip, 0); err != nil {
		log.Fatalf("Failed to add VIP: %v", err)
	}

	// Add real servers (backends) to the VIP.
	reals := []katran.Real{
		{Address: "192.168.1.1", Weight: 100, Flags: 0},
		{Address: "192.168.1.2", Weight: 100, Flags: 0},
		{Address: "192.168.1.3", Weight: 50, Flags: 0},
	}

	for _, real := range reals {
		if err := lb.AddRealForVIP(real, vip); err != nil {
			log.Printf("Failed to add real %s: %v", real.Address, err)
		}
	}

	// Get all configured VIPs.
	vips, err := lb.GetAllVIPs()
	if err != nil {
		log.Fatalf("Failed to get VIPs: %v", err)
	}
	fmt.Printf("Configured VIPs: %d\n", len(vips))

	// Get reals for a specific VIP.
	configuredReals, err := lb.GetRealsForVIP(vip)
	if err != nil {
		log.Fatalf("Failed to get reals: %v", err)
	}
	fmt.Printf("Reals for VIP %s:%d: %d\n", vip.Address, vip.Port, len(configuredReals))

	// Output:
	// Configured VIPs: 1
	// Reals for VIP 10.0.0.1:80: 3
}

// Example_errorHandling demonstrates error handling with Katran errors.
func Example_errorHandling() {
	cfg := katran.NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	vip := katran.VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6,
	}

	// Add VIP.
	if err := lb.AddVIP(vip, 0); err != nil {
		log.Fatalf("Failed to add VIP: %v", err)
	}

	// Try to add the same VIP again - should fail with ErrAlreadyExists.
	err = lb.AddVIP(vip, 0)
	if err != nil {
		if katran.IsAlreadyExists(err) {
			fmt.Println("VIP already exists (expected)")
		} else {
			log.Fatalf("Unexpected error: %v", err)
		}
	}

	// Try to delete a non-existent VIP.
	nonExistentVIP := katran.VIPKey{
		Address: "10.99.99.99",
		Port:    80,
		Proto:   6,
	}
	err = lb.DelVIP(nonExistentVIP)
	if err != nil {
		if katran.IsNotFound(err) {
			fmt.Println("VIP not found (expected)")
		} else {
			log.Fatalf("Unexpected error: %v", err)
		}
	}

	// Output:
	// VIP already exists (expected)
	// VIP not found (expected)
}

// Example_batchOperations demonstrates batch operations for efficiency.
func Example_batchOperations() {
	cfg := katran.NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	vip := katran.VIPKey{
		Address: "10.0.0.1",
		Port:    443,
		Proto:   6,
	}

	if err := lb.AddVIP(vip, 0); err != nil {
		log.Fatalf("Failed to add VIP: %v", err)
	}

	// Batch add multiple reals at once (more efficient than individual adds).
	reals := []katran.Real{
		{Address: "192.168.1.1", Weight: 100, Flags: 0},
		{Address: "192.168.1.2", Weight: 100, Flags: 0},
		{Address: "192.168.1.3", Weight: 100, Flags: 0},
		{Address: "192.168.1.4", Weight: 100, Flags: 0},
	}

	err = lb.ModifyRealsForVIP(katran.ActionAdd, reals, vip)
	if err != nil {
		log.Fatalf("Failed to batch add reals: %v", err)
	}

	// Verify.
	configuredReals, _ := lb.GetRealsForVIP(vip)
	fmt.Printf("Added %d reals in batch\n", len(configuredReals))

	// Batch remove some reals.
	realsToRemove := []katran.Real{
		{Address: "192.168.1.3", Weight: 0, Flags: 0},
		{Address: "192.168.1.4", Weight: 0, Flags: 0},
	}

	err = lb.ModifyRealsForVIP(katran.ActionDel, realsToRemove, vip)
	if err != nil {
		log.Fatalf("Failed to batch remove reals: %v", err)
	}

	// Verify.
	configuredReals, _ = lb.GetRealsForVIP(vip)
	fmt.Printf("Remaining reals after batch remove: %d\n", len(configuredReals))

	// Output:
	// Added 4 reals in batch
	// Remaining reals after batch remove: 2
}

// Example_flowSimulation demonstrates flow simulation for debugging.
func Example_flowSimulation() {
	cfg := katran.NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	// Set up a VIP with reals.
	vip := katran.VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6,
	}
	lb.AddVIP(vip, 0)
	lb.AddRealForVIP(katran.Real{Address: "192.168.1.1", Weight: 100}, vip)

	// Simulate a flow to see which real it would be routed to.
	flow := katran.Flow{
		Src:     "172.16.0.100",
		Dst:     "10.0.0.1",
		SrcPort: 54321,
		DstPort: 80,
		Proto:   6,
	}

	realAddr, err := lb.GetRealForFlow(flow)
	if err != nil {
		log.Fatalf("Failed to get real for flow: %v", err)
	}

	fmt.Printf("Flow %s:%d -> %s:%d would be routed to: %s\n",
		flow.Src, flow.SrcPort, flow.Dst, flow.DstPort, realAddr)

	// Output:
	// Flow 172.16.0.100:54321 -> 10.0.0.1:80 would be routed to: 192.168.1.1
}

// Example_statistics demonstrates statistics retrieval.
func Example_statistics() {
	cfg := katran.NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	vip := katran.VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6,
	}
	lb.AddVIP(vip, 0)

	// Get VIP statistics.
	stats, err := lb.GetStatsForVIP(vip)
	if err != nil {
		log.Fatalf("Failed to get VIP stats: %v", err)
	}
	fmt.Printf("VIP stats - Packets: %d, Bytes: %d\n", stats.V1, stats.V2)

	// Get LRU statistics.
	lruStats, err := lb.GetLRUStats()
	if err != nil {
		log.Fatalf("Failed to get LRU stats: %v", err)
	}
	fmt.Printf("LRU stats - Total: %d, Hits: %d\n", lruStats.V1, lruStats.V2)

	// Get userspace statistics.
	usStats, err := lb.GetUserspaceStats()
	if err != nil {
		log.Fatalf("Failed to get userspace stats: %v", err)
	}
	fmt.Printf("Userspace stats - BPF failures: %d, Addr validation failures: %d\n",
		usStats.BPFFailedCalls, usStats.AddrValidationFailed)

	// Output:
	// VIP stats - Packets: 0, Bytes: 0
	// LRU stats - Total: 0, Hits: 0
	// Userspace stats - BPF failures: 0, Addr validation failures: 0
}

// Example_configCustomization demonstrates customizing the configuration.
func Example_configCustomization() {
	// Start with defaults.
	cfg := katran.NewConfig()

	// Customize for your environment.
	cfg.MainInterface = "eth0"
	cfg.MaxVIPs = 1024
	cfg.MaxReals = 16384
	cfg.CHRingSize = 131072
	cfg.LRUSize = 16000000
	cfg.HashFunc = katran.HashMaglevV2
	cfg.Testing = true // Only for this example

	// Set MAC addresses.
	cfg.DefaultMAC = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	cfg.LocalMAC = []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

	// Set forwarding cores for multi-core systems.
	cfg.ForwardingCores = []int32{0, 1, 2, 3}
	cfg.NUMANodes = []int32{0, 0, 1, 1}

	// Enable specific features.
	cfg.EnableHC = true
	cfg.FlowDebug = true

	fmt.Printf("Config: MaxVIPs=%d, MaxReals=%d, HashFunc=%d\n",
		cfg.MaxVIPs, cfg.MaxReals, cfg.HashFunc)

	// Output:
	// Config: MaxVIPs=1024, MaxReals=16384, HashFunc=1
}

// Example_ipv6Support demonstrates IPv6 VIP and real server support.
func Example_ipv6Support() {
	cfg := katran.NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := katran.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	defer lb.Close()

	// IPv6 VIP.
	vip := katran.VIPKey{
		Address: "2001:db8::1",
		Port:    443,
		Proto:   6,
	}

	if err := lb.AddVIP(vip, 0); err != nil {
		log.Fatalf("Failed to add IPv6 VIP: %v", err)
	}

	// IPv6 real servers.
	reals := []katran.Real{
		{Address: "2001:db8:1::1", Weight: 100, Flags: 0},
		{Address: "2001:db8:1::2", Weight: 100, Flags: 0},
	}

	for _, real := range reals {
		if err := lb.AddRealForVIP(real, vip); err != nil {
			log.Printf("Failed to add IPv6 real: %v", err)
		}
	}

	// Verify.
	configuredReals, _ := lb.GetRealsForVIP(vip)
	fmt.Printf("IPv6 VIP %s has %d reals\n", vip.Address, len(configuredReals))

	// Output:
	// IPv6 VIP 2001:db8::1 has 2 reals
}
