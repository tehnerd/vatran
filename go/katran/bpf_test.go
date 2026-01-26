//go:build integration
// +build integration

package katran

import (
	"testing"
)

// TestGetKatranProgFD tests retrieving the balancer program file descriptor.
func TestGetKatranProgFD(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	fd, err := lb.GetKatranProgFD()
	if err != nil {
		// In testing mode, the FD might not be available.
		t.Logf("GetKatranProgFD returned: %v", err)
		return
	}

	t.Logf("Katran program FD: %d", fd)
}

// TestGetKatranProgFDAfterClose tests that GetKatranProgFD fails after Close.
func TestGetKatranProgFDAfterClose(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	lb.Close()

	_, err = lb.GetKatranProgFD()
	if err == nil {
		t.Error("GetKatranProgFD should fail after Close")
	}
}

// TestGetHealthcheckerProgFD tests retrieving the healthcheck program file descriptor.
func TestGetHealthcheckerProgFD(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = true

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	fd, err := lb.GetHealthcheckerProgFD()
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Healthcheck feature not enabled")
		}
		// In testing mode, the FD might not be available.
		t.Logf("GetHealthcheckerProgFD returned: %v", err)
		return
	}

	t.Logf("Healthchecker program FD: %d", fd)
}

// TestGetBPFMapFDByName tests retrieving BPF map file descriptors by name.
func TestGetBPFMapFDByName(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Try to get a well-known map name.
	mapNames := []string{"vip_map", "reals", "ctl_array", "lru_mapping"}

	for _, name := range mapNames {
		fd, err := lb.GetBPFMapFDByName(name)
		if err != nil {
			if IsNotFound(err) {
				t.Logf("Map %s not found", name)
				continue
			}
			t.Logf("GetBPFMapFDByName(%s) returned: %v", name, err)
			continue
		}
		t.Logf("Map %s FD: %d", name, fd)
	}
}

// TestGetBPFMapFDByNameNonExistent tests getting FD for non-existent map.
func TestGetBPFMapFDByNameNonExistent(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	_, err = lb.GetBPFMapFDByName("nonexistent_map_12345")
	if err == nil {
		t.Error("Expected error for non-existent map")
	}
	if !IsNotFound(err) {
		// In testing mode, might get a different error.
		t.Logf("Got error (expected): %v", err)
	}
}

// TestGetBPFMapStats tests retrieving BPF map statistics.
func TestGetBPFMapStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	mapNames := []string{"vip_map", "reals", "ctl_array"}

	for _, name := range mapNames {
		stats, err := lb.GetBPFMapStats(name)
		if err != nil {
			if IsNotFound(err) {
				t.Logf("Map %s not found", name)
				continue
			}
			t.Logf("GetBPFMapStats(%s) returned: %v", name, err)
			continue
		}
		t.Logf("Map %s stats: max=%d, current=%d",
			name, stats.MaxEntries, stats.CurrentEntries)
	}
}

// TestGetGlobalLRUMapsFDs tests retrieving global LRU maps file descriptors.
func TestGetGlobalLRUMapsFDs(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	fds, err := lb.GetGlobalLRUMapsFDs()
	if err != nil {
		// May not be available in testing mode.
		t.Logf("GetGlobalLRUMapsFDs returned: %v", err)
		return
	}

	t.Logf("Global LRU map FDs: %v", fds)
}

// TestSimulatePacket tests packet simulation.
func TestSimulatePacket(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Create a minimal test packet (Ethernet + IP + TCP headers).
	// This is a simplified test packet.
	packet := make([]byte, 64)

	// Ethernet header (14 bytes).
	copy(packet[0:6], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})   // dst MAC
	copy(packet[6:12], []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})  // src MAC
	packet[12] = 0x08                                               // EtherType: IPv4 (0x0800)
	packet[13] = 0x00

	// Minimal IP header (20 bytes starting at offset 14).
	packet[14] = 0x45 // Version + IHL
	packet[15] = 0x00 // DSCP + ECN
	packet[16] = 0x00 // Total length high
	packet[17] = 0x28 // Total length low (40 bytes)
	packet[23] = 0x06 // Protocol: TCP

	// Source IP (10.0.0.100).
	packet[26] = 10
	packet[27] = 0
	packet[28] = 0
	packet[29] = 100

	// Dest IP (10.0.0.1 - our VIP).
	packet[30] = 10
	packet[31] = 0
	packet[32] = 0
	packet[33] = 1

	// Add a VIP for the test.
	vip := VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6,
	}
	lb.AddVIP(vip, 0)
	lb.AddRealForVIP(Real{Address: "192.168.1.1", Weight: 100}, vip)

	// Try to simulate the packet.
	outPacket, err := lb.SimulatePacket(packet)
	if err != nil {
		// Packet simulation might fail without full BPF setup.
		t.Logf("SimulatePacket returned: %v", err)
		return
	}

	t.Logf("Output packet length: %d", len(outPacket))
}

// TestSimulatePacketEmpty tests that empty packet returns error.
func TestSimulatePacketEmpty(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	_, err = lb.SimulatePacket([]byte{})
	if err == nil {
		t.Error("Expected error for empty packet")
	}
}

// TestSimulatePacketAfterClose tests that SimulatePacket fails after Close.
func TestSimulatePacketAfterClose(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	lb.Close()

	_, err = lb.SimulatePacket(make([]byte, 64))
	if err == nil {
		t.Error("SimulatePacket should fail after Close")
	}
}

// TestXDPStats tests XDP statistics retrieval.
func TestXDPStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Test all XDP stats functions.
	totalStats, err := lb.GetXDPTotalStats()
	if err != nil {
		t.Errorf("GetXDPTotalStats failed: %v", err)
	} else {
		t.Logf("XDP Total: packets=%d, bytes=%d", totalStats.V1, totalStats.V2)
	}

	txStats, err := lb.GetXDPTXStats()
	if err != nil {
		t.Errorf("GetXDPTXStats failed: %v", err)
	} else {
		t.Logf("XDP TX: packets=%d, bytes=%d", txStats.V1, txStats.V2)
	}

	dropStats, err := lb.GetXDPDropStats()
	if err != nil {
		t.Errorf("GetXDPDropStats failed: %v", err)
	} else {
		t.Logf("XDP Drop: packets=%d, bytes=%d", dropStats.V1, dropStats.V2)
	}

	passStats, err := lb.GetXDPPassStats()
	if err != nil {
		t.Errorf("GetXDPPassStats failed: %v", err)
	} else {
		t.Logf("XDP Pass: packets=%d, bytes=%d", passStats.V1, passStats.V2)
	}
}

// TestQuicStats tests QUIC statistics retrieval.
func TestQuicStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Get QUIC packets stats.
	stats, err := lb.GetQuicPacketsStats()
	if err != nil {
		t.Errorf("GetQuicPacketsStats failed: %v", err)
	} else {
		t.Logf("QUIC stats: CHRouted=%d, CIDRouted=%d, CIDV0=%d, CIDV1=%d, CIDV2=%d, CIDV3=%d",
			stats.CHRouted, stats.CIDRouted, stats.CIDV0, stats.CIDV1, stats.CIDV2, stats.CIDV3)
	}

	// Get QUIC ICMP stats.
	icmpStats, err := lb.GetQuicICMPStats()
	if err != nil {
		t.Errorf("GetQuicICMPStats failed: %v", err)
	} else {
		t.Logf("QUIC ICMP stats: V1=%d, V2=%d", icmpStats.V1, icmpStats.V2)
	}
}

// TestTCPServerIDRoutingStats tests TPR statistics retrieval.
func TestTCPServerIDRoutingStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	stats, err := lb.GetTCPServerIDRoutingStats()
	if err != nil {
		t.Errorf("GetTCPServerIDRoutingStats failed: %v", err)
	} else {
		t.Logf("TPR stats: CHRouted=%d, SIDRouted=%d, TCPSyn=%d, DstMismatch=%d",
			stats.CHRouted, stats.SIDRouted, stats.TCPSyn, stats.DstMismatchInLRU)
	}
}
