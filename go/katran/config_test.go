package katran

import (
	"testing"
)

// TestNewConfig tests that NewConfig creates a config with proper defaults.
func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}

	// Check default values.
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"RootMapPos", cfg.RootMapPos, uint32(2)},
		{"UseRootMap", cfg.UseRootMap, true},
		{"MaxVIPs", cfg.MaxVIPs, uint32(512)},
		{"MaxReals", cfg.MaxReals, uint32(4096)},
		{"CHRingSize", cfg.CHRingSize, uint32(65537)},
		{"LRUSize", cfg.LRUSize, uint64(8000000)},
		{"MaxLPMSrcSize", cfg.MaxLPMSrcSize, uint32(3000000)},
		{"MaxDecapDst", cfg.MaxDecapDst, uint32(6)},
		{"GlobalLRUSize", cfg.GlobalLRUSize, uint32(100000)},
		{"EnableHC", cfg.EnableHC, true},
		{"TunnelBasedHCEncap", cfg.TunnelBasedHCEncap, true},
		{"MemlockUnlimited", cfg.MemlockUnlimited, true},
		{"CleanupOnShutdown", cfg.CleanupOnShutdown, true},
		{"Priority", cfg.Priority, uint32(2307)},
		{"HashFunc", cfg.HashFunc, HashMaglev},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.got)
			}
		})
	}

	// Check that other fields are zero values.
	if cfg.MainInterface != "" {
		t.Errorf("expected MainInterface to be empty, got %s", cfg.MainInterface)
	}
	if cfg.Testing != false {
		t.Error("expected Testing to be false")
	}
	if len(cfg.ForwardingCores) != 0 {
		t.Errorf("expected ForwardingCores to be empty, got %v", cfg.ForwardingCores)
	}
	if len(cfg.NUMANodes) != 0 {
		t.Errorf("expected NUMANodes to be empty, got %v", cfg.NUMANodes)
	}
}

// TestConfigFields tests that all Config fields can be set.
func TestConfigFields(t *testing.T) {
	cfg := &Config{
		MainInterface:          "eth0",
		V4TunInterface:         "ipip0",
		V6TunInterface:         "ip6tnl0",
		HCInterface:            "eth1",
		BalancerProgPath:       "/path/to/balancer.o",
		HealthcheckingProgPath: "/path/to/healthcheck.o",
		DefaultMAC:             []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		LocalMAC:               []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		RootMapPath:            "/sys/fs/bpf/root_map",
		RootMapPos:             3,
		UseRootMap:             true,
		MaxVIPs:                1024,
		MaxReals:               8192,
		CHRingSize:             131072,
		LRUSize:                16000000,
		MaxLPMSrcSize:          6000000,
		MaxDecapDst:            12,
		GlobalLRUSize:          200000,
		EnableHC:               true,
		TunnelBasedHCEncap:     false,
		Testing:                true,
		MemlockUnlimited:       true,
		FlowDebug:              true,
		EnableCIDV3:            true,
		CleanupOnShutdown:      true,
		ForwardingCores:        []int32{0, 1, 2, 3},
		NUMANodes:              []int32{0, 0, 1, 1},
		XDPAttachFlags:         0x01,
		Priority:               1000,
		MainInterfaceIndex:     2,
		HCInterfaceIndex:       3,
		KatranSrcV4:            "10.0.0.1",
		KatranSrcV6:            "2001:db8::1",
		HashFunc:               HashMaglevV2,
	}

	// Verify all fields are set correctly.
	if cfg.MainInterface != "eth0" {
		t.Errorf("expected MainInterface eth0, got %s", cfg.MainInterface)
	}
	if cfg.V4TunInterface != "ipip0" {
		t.Errorf("expected V4TunInterface ipip0, got %s", cfg.V4TunInterface)
	}
	if cfg.V6TunInterface != "ip6tnl0" {
		t.Errorf("expected V6TunInterface ip6tnl0, got %s", cfg.V6TunInterface)
	}
	if cfg.HCInterface != "eth1" {
		t.Errorf("expected HCInterface eth1, got %s", cfg.HCInterface)
	}
	if cfg.BalancerProgPath != "/path/to/balancer.o" {
		t.Errorf("expected BalancerProgPath /path/to/balancer.o, got %s", cfg.BalancerProgPath)
	}
	if cfg.HealthcheckingProgPath != "/path/to/healthcheck.o" {
		t.Errorf("expected HealthcheckingProgPath /path/to/healthcheck.o, got %s", cfg.HealthcheckingProgPath)
	}
	if len(cfg.DefaultMAC) != 6 {
		t.Errorf("expected DefaultMAC length 6, got %d", len(cfg.DefaultMAC))
	}
	if cfg.DefaultMAC[0] != 0x00 || cfg.DefaultMAC[5] != 0x55 {
		t.Error("DefaultMAC values incorrect")
	}
	if len(cfg.LocalMAC) != 6 {
		t.Errorf("expected LocalMAC length 6, got %d", len(cfg.LocalMAC))
	}
	if cfg.RootMapPath != "/sys/fs/bpf/root_map" {
		t.Errorf("expected RootMapPath /sys/fs/bpf/root_map, got %s", cfg.RootMapPath)
	}
	if cfg.RootMapPos != 3 {
		t.Errorf("expected RootMapPos 3, got %d", cfg.RootMapPos)
	}
	if !cfg.UseRootMap {
		t.Error("expected UseRootMap true")
	}
	if cfg.MaxVIPs != 1024 {
		t.Errorf("expected MaxVIPs 1024, got %d", cfg.MaxVIPs)
	}
	if cfg.MaxReals != 8192 {
		t.Errorf("expected MaxReals 8192, got %d", cfg.MaxReals)
	}
	if cfg.CHRingSize != 131072 {
		t.Errorf("expected CHRingSize 131072, got %d", cfg.CHRingSize)
	}
	if cfg.LRUSize != 16000000 {
		t.Errorf("expected LRUSize 16000000, got %d", cfg.LRUSize)
	}
	if cfg.MaxLPMSrcSize != 6000000 {
		t.Errorf("expected MaxLPMSrcSize 6000000, got %d", cfg.MaxLPMSrcSize)
	}
	if cfg.MaxDecapDst != 12 {
		t.Errorf("expected MaxDecapDst 12, got %d", cfg.MaxDecapDst)
	}
	if cfg.GlobalLRUSize != 200000 {
		t.Errorf("expected GlobalLRUSize 200000, got %d", cfg.GlobalLRUSize)
	}
	if !cfg.EnableHC {
		t.Error("expected EnableHC true")
	}
	if cfg.TunnelBasedHCEncap {
		t.Error("expected TunnelBasedHCEncap false")
	}
	if !cfg.Testing {
		t.Error("expected Testing true")
	}
	if !cfg.MemlockUnlimited {
		t.Error("expected MemlockUnlimited true")
	}
	if !cfg.FlowDebug {
		t.Error("expected FlowDebug true")
	}
	if !cfg.EnableCIDV3 {
		t.Error("expected EnableCIDV3 true")
	}
	if !cfg.CleanupOnShutdown {
		t.Error("expected CleanupOnShutdown true")
	}
	if len(cfg.ForwardingCores) != 4 {
		t.Errorf("expected ForwardingCores length 4, got %d", len(cfg.ForwardingCores))
	}
	if cfg.ForwardingCores[0] != 0 || cfg.ForwardingCores[3] != 3 {
		t.Error("ForwardingCores values incorrect")
	}
	if len(cfg.NUMANodes) != 4 {
		t.Errorf("expected NUMANodes length 4, got %d", len(cfg.NUMANodes))
	}
	if cfg.XDPAttachFlags != 0x01 {
		t.Errorf("expected XDPAttachFlags 0x01, got %d", cfg.XDPAttachFlags)
	}
	if cfg.Priority != 1000 {
		t.Errorf("expected Priority 1000, got %d", cfg.Priority)
	}
	if cfg.MainInterfaceIndex != 2 {
		t.Errorf("expected MainInterfaceIndex 2, got %d", cfg.MainInterfaceIndex)
	}
	if cfg.HCInterfaceIndex != 3 {
		t.Errorf("expected HCInterfaceIndex 3, got %d", cfg.HCInterfaceIndex)
	}
	if cfg.KatranSrcV4 != "10.0.0.1" {
		t.Errorf("expected KatranSrcV4 10.0.0.1, got %s", cfg.KatranSrcV4)
	}
	if cfg.KatranSrcV6 != "2001:db8::1" {
		t.Errorf("expected KatranSrcV6 2001:db8::1, got %s", cfg.KatranSrcV6)
	}
	if cfg.HashFunc != HashMaglevV2 {
		t.Errorf("expected HashFunc HashMaglevV2, got %d", cfg.HashFunc)
	}
}

// TestConfigCopy tests that Config can be copied correctly.
func TestConfigCopy(t *testing.T) {
	original := NewConfig()
	original.MainInterface = "eth0"
	original.MaxVIPs = 2048
	original.ForwardingCores = []int32{0, 1, 2, 3}
	original.DefaultMAC = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}

	// Create a copy (shallow copy).
	copied := *original

	// Modify the copy.
	copied.MainInterface = "eth1"
	copied.MaxVIPs = 4096

	// Verify original is unchanged.
	if original.MainInterface != "eth0" {
		t.Error("original MainInterface was modified")
	}
	if original.MaxVIPs != 2048 {
		t.Error("original MaxVIPs was modified")
	}

	// Note: slices are shared in shallow copy.
	// This is expected Go behavior, not a bug in the Config type.
}

// TestConfigWithTestingMode tests that Testing mode can be set.
func TestConfigWithTestingMode(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true

	if !cfg.Testing {
		t.Error("expected Testing to be true")
	}
}

// TestConfigMACAddresses tests MAC address handling.
func TestConfigMACAddresses(t *testing.T) {
	cfg := NewConfig()

	// Valid 6-byte MAC.
	cfg.DefaultMAC = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	if len(cfg.DefaultMAC) != 6 {
		t.Errorf("expected DefaultMAC length 6, got %d", len(cfg.DefaultMAC))
	}

	// Setting to nil should work (will be handled by toC()).
	cfg.DefaultMAC = nil
	if cfg.DefaultMAC != nil {
		t.Error("expected DefaultMAC to be nil")
	}

	// Wrong length should still be stored (validation happens in toC()).
	cfg.DefaultMAC = []byte{0x00, 0x11} // Wrong length
	if len(cfg.DefaultMAC) != 2 {
		t.Errorf("expected DefaultMAC length 2, got %d", len(cfg.DefaultMAC))
	}
}

// TestConfigForwardingCores tests ForwardingCores handling.
func TestConfigForwardingCores(t *testing.T) {
	cfg := NewConfig()

	// Empty slice.
	if len(cfg.ForwardingCores) != 0 {
		t.Errorf("expected empty ForwardingCores, got %d elements", len(cfg.ForwardingCores))
	}

	// Set cores.
	cfg.ForwardingCores = []int32{0, 2, 4, 6}
	if len(cfg.ForwardingCores) != 4 {
		t.Errorf("expected 4 cores, got %d", len(cfg.ForwardingCores))
	}
	if cfg.ForwardingCores[0] != 0 {
		t.Errorf("expected core 0 at index 0, got %d", cfg.ForwardingCores[0])
	}
	if cfg.ForwardingCores[3] != 6 {
		t.Errorf("expected core 6 at index 3, got %d", cfg.ForwardingCores[3])
	}

	// Single core.
	cfg.ForwardingCores = []int32{0}
	if len(cfg.ForwardingCores) != 1 {
		t.Errorf("expected 1 core, got %d", len(cfg.ForwardingCores))
	}
}

// TestConfigNUMANodes tests NUMANodes handling.
func TestConfigNUMANodes(t *testing.T) {
	cfg := NewConfig()

	// Empty slice.
	if len(cfg.NUMANodes) != 0 {
		t.Errorf("expected empty NUMANodes, got %d elements", len(cfg.NUMANodes))
	}

	// Set NUMA nodes (should match ForwardingCores length).
	cfg.ForwardingCores = []int32{0, 1, 2, 3}
	cfg.NUMANodes = []int32{0, 0, 1, 1}
	if len(cfg.NUMANodes) != 4 {
		t.Errorf("expected 4 NUMA nodes, got %d", len(cfg.NUMANodes))
	}
	if cfg.NUMANodes[0] != 0 {
		t.Errorf("expected NUMA node 0 at index 0, got %d", cfg.NUMANodes[0])
	}
	if cfg.NUMANodes[2] != 1 {
		t.Errorf("expected NUMA node 1 at index 2, got %d", cfg.NUMANodes[2])
	}
}

// TestConfigIPAddresses tests IP address string handling.
func TestConfigIPAddresses(t *testing.T) {
	cfg := NewConfig()

	// IPv4 addresses.
	cfg.KatranSrcV4 = "192.168.1.1"
	if cfg.KatranSrcV4 != "192.168.1.1" {
		t.Errorf("expected KatranSrcV4 192.168.1.1, got %s", cfg.KatranSrcV4)
	}

	// IPv6 addresses.
	cfg.KatranSrcV6 = "fe80::1"
	if cfg.KatranSrcV6 != "fe80::1" {
		t.Errorf("expected KatranSrcV6 fe80::1, got %s", cfg.KatranSrcV6)
	}

	// Full IPv6.
	cfg.KatranSrcV6 = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	if cfg.KatranSrcV6 != "2001:0db8:85a3:0000:0000:8a2e:0370:7334" {
		t.Errorf("unexpected KatranSrcV6: %s", cfg.KatranSrcV6)
	}
}

// TestConfigHashFunction tests HashFunction field.
func TestConfigHashFunction(t *testing.T) {
	cfg := NewConfig()

	// Default should be HashMaglev.
	if cfg.HashFunc != HashMaglev {
		t.Errorf("expected default HashFunc HashMaglev, got %d", cfg.HashFunc)
	}

	// Can be changed to HashMaglevV2.
	cfg.HashFunc = HashMaglevV2
	if cfg.HashFunc != HashMaglevV2 {
		t.Errorf("expected HashFunc HashMaglevV2, got %d", cfg.HashFunc)
	}
}

// TestConfigPaths tests file path fields.
func TestConfigPaths(t *testing.T) {
	cfg := NewConfig()

	// Set all path fields.
	cfg.BalancerProgPath = "/usr/lib/katran/balancer.o"
	cfg.HealthcheckingProgPath = "/usr/lib/katran/healthcheck.o"
	cfg.RootMapPath = "/sys/fs/bpf/katran/root_map"

	if cfg.BalancerProgPath != "/usr/lib/katran/balancer.o" {
		t.Errorf("unexpected BalancerProgPath: %s", cfg.BalancerProgPath)
	}
	if cfg.HealthcheckingProgPath != "/usr/lib/katran/healthcheck.o" {
		t.Errorf("unexpected HealthcheckingProgPath: %s", cfg.HealthcheckingProgPath)
	}
	if cfg.RootMapPath != "/sys/fs/bpf/katran/root_map" {
		t.Errorf("unexpected RootMapPath: %s", cfg.RootMapPath)
	}

	// Empty paths are valid.
	cfg.BalancerProgPath = ""
	if cfg.BalancerProgPath != "" {
		t.Error("expected empty BalancerProgPath")
	}
}

// TestConfigInterfaces tests interface name fields.
func TestConfigInterfaces(t *testing.T) {
	cfg := NewConfig()

	// Set interface names.
	cfg.MainInterface = "enp0s3"
	cfg.V4TunInterface = "ipip_v4"
	cfg.V6TunInterface = "ip6_v6"
	cfg.HCInterface = "veth0"

	if cfg.MainInterface != "enp0s3" {
		t.Errorf("unexpected MainInterface: %s", cfg.MainInterface)
	}
	if cfg.V4TunInterface != "ipip_v4" {
		t.Errorf("unexpected V4TunInterface: %s", cfg.V4TunInterface)
	}
	if cfg.V6TunInterface != "ip6_v6" {
		t.Errorf("unexpected V6TunInterface: %s", cfg.V6TunInterface)
	}
	if cfg.HCInterface != "veth0" {
		t.Errorf("unexpected HCInterface: %s", cfg.HCInterface)
	}
}

// TestConfigInterfaceIndices tests interface index fields.
func TestConfigInterfaceIndices(t *testing.T) {
	cfg := NewConfig()

	// Default should be 0.
	if cfg.MainInterfaceIndex != 0 {
		t.Errorf("expected MainInterfaceIndex 0, got %d", cfg.MainInterfaceIndex)
	}
	if cfg.HCInterfaceIndex != 0 {
		t.Errorf("expected HCInterfaceIndex 0, got %d", cfg.HCInterfaceIndex)
	}

	// Set indices.
	cfg.MainInterfaceIndex = 2
	cfg.HCInterfaceIndex = 5

	if cfg.MainInterfaceIndex != 2 {
		t.Errorf("expected MainInterfaceIndex 2, got %d", cfg.MainInterfaceIndex)
	}
	if cfg.HCInterfaceIndex != 5 {
		t.Errorf("expected HCInterfaceIndex 5, got %d", cfg.HCInterfaceIndex)
	}
}

// TestConfigBooleanFields tests all boolean configuration fields.
func TestConfigBooleanFields(t *testing.T) {
	cfg := NewConfig()

	// Test toggling all boolean fields.
	boolTests := []struct {
		name  string
		get   func() bool
		set   func(bool)
		dflt  bool
	}{
		{"UseRootMap", func() bool { return cfg.UseRootMap }, func(v bool) { cfg.UseRootMap = v }, true},
		{"EnableHC", func() bool { return cfg.EnableHC }, func(v bool) { cfg.EnableHC = v }, true},
		{"TunnelBasedHCEncap", func() bool { return cfg.TunnelBasedHCEncap }, func(v bool) { cfg.TunnelBasedHCEncap = v }, true},
		{"Testing", func() bool { return cfg.Testing }, func(v bool) { cfg.Testing = v }, false},
		{"MemlockUnlimited", func() bool { return cfg.MemlockUnlimited }, func(v bool) { cfg.MemlockUnlimited = v }, true},
		{"FlowDebug", func() bool { return cfg.FlowDebug }, func(v bool) { cfg.FlowDebug = v }, false},
		{"EnableCIDV3", func() bool { return cfg.EnableCIDV3 }, func(v bool) { cfg.EnableCIDV3 = v }, false},
		{"CleanupOnShutdown", func() bool { return cfg.CleanupOnShutdown }, func(v bool) { cfg.CleanupOnShutdown = v }, true},
	}

	for _, tt := range boolTests {
		t.Run(tt.name, func(t *testing.T) {
			// Check default.
			if tt.get() != tt.dflt {
				t.Errorf("expected default %v, got %v", tt.dflt, tt.get())
			}

			// Toggle to opposite.
			tt.set(!tt.dflt)
			if tt.get() != !tt.dflt {
				t.Errorf("expected %v after toggle, got %v", !tt.dflt, tt.get())
			}

			// Toggle back.
			tt.set(tt.dflt)
			if tt.get() != tt.dflt {
				t.Errorf("expected %v after toggle back, got %v", tt.dflt, tt.get())
			}
		})
	}
}
