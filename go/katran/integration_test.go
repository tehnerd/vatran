//go:build integration
// +build integration

package katran

import (
	"testing"
)

// NOTE: Integration tests require the C library to be built and available.
// Run with: go test -tags=integration -v ./...
//
// These tests require:
// 1. The katran C library to be built in ../../_build/
// 2. BPF programs to be compiled
// 3. Sufficient privileges to load BPF programs (typically root or CAP_BPF)

// createTestLB creates a LoadBalancer in testing mode for unit tests.
// Testing mode doesn't actually program the forwarding plane.
func createTestLB(t *testing.T) *LoadBalancer {
	t.Helper()

	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}

	return lb
}

// TestLoadBalancerCreation tests basic LoadBalancer creation and destruction.
func TestLoadBalancerCreation(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	// Verify lb is not nil.
	if lb == nil {
		t.Fatal("LoadBalancer is nil")
	}
}

// TestLoadBalancerDoubleClose tests that Close can be called multiple times safely.
func TestLoadBalancerDoubleClose(t *testing.T) {
	lb := createTestLB(t)

	// First close should succeed.
	err := lb.Close()
	if err != nil {
		t.Errorf("First Close failed: %v", err)
	}

	// Second close should also succeed (no error).
	err = lb.Close()
	if err != nil {
		t.Errorf("Second Close failed: %v", err)
	}
}

// TestLoadBalancerOperationsAfterClose tests that operations fail after Close.
func TestLoadBalancerOperationsAfterClose(t *testing.T) {
	lb := createTestLB(t)
	lb.Close()

	// All operations should return an error after close.
	vip := VIPKey{Address: "10.0.0.1", Port: 80, Proto: 6}

	_, err := lb.GetAllVIPs()
	if err == nil {
		t.Error("GetAllVIPs should fail after Close")
	}

	err = lb.AddVIP(vip, 0)
	if err == nil {
		t.Error("AddVIP should fail after Close")
	}

	err = lb.DelVIP(vip)
	if err == nil {
		t.Error("DelVIP should fail after Close")
	}
}

// TestAddDeleteVIP tests VIP add and delete operations.
func TestAddDeleteVIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6, // TCP
	}

	// Add VIP.
	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Get all VIPs.
	vips, err := lb.GetAllVIPs()
	if err != nil {
		t.Fatalf("GetAllVIPs failed: %v", err)
	}

	// Should have at least one VIP.
	found := false
	for _, v := range vips {
		if v.Address == vip.Address && v.Port == vip.Port && v.Proto == vip.Proto {
			found = true
			break
		}
	}
	if !found {
		t.Error("VIP not found after adding")
	}

	// Delete VIP.
	err = lb.DelVIP(vip)
	if err != nil {
		t.Fatalf("DelVIP failed: %v", err)
	}

	// Verify VIP is deleted.
	vips, err = lb.GetAllVIPs()
	if err != nil {
		t.Fatalf("GetAllVIPs failed after delete: %v", err)
	}

	for _, v := range vips {
		if v.Address == vip.Address && v.Port == vip.Port && v.Proto == vip.Proto {
			t.Error("VIP still exists after delete")
		}
	}
}

// TestAddDuplicateVIP tests that adding a duplicate VIP returns an error.
func TestAddDuplicateVIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.2",
		Port:    443,
		Proto:   6,
	}

	// Add VIP first time.
	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("First AddVIP failed: %v", err)
	}

	// Add same VIP again should fail.
	err = lb.AddVIP(vip, 0)
	if err == nil {
		t.Error("Second AddVIP should have failed")
	}
	if !IsAlreadyExists(err) {
		t.Errorf("Expected ErrAlreadyExists, got: %v", err)
	}
}

// TestDeleteNonExistentVIP tests that deleting a non-existent VIP returns an error.
func TestDeleteNonExistentVIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.99",
		Port:    80,
		Proto:   6,
	}

	err := lb.DelVIP(vip)
	if err == nil {
		t.Error("DelVIP should have failed for non-existent VIP")
	}
	if !IsNotFound(err) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}

// TestIPv6VIP tests VIP operations with IPv6 addresses.
func TestIPv6VIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "2001:db8::1",
		Port:    80,
		Proto:   6,
	}

	// Add IPv6 VIP.
	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP (IPv6) failed: %v", err)
	}

	// Get all VIPs.
	vips, err := lb.GetAllVIPs()
	if err != nil {
		t.Fatalf("GetAllVIPs failed: %v", err)
	}

	found := false
	for _, v := range vips {
		if v.Address == vip.Address && v.Port == vip.Port {
			found = true
			break
		}
	}
	if !found {
		t.Error("IPv6 VIP not found")
	}

	// Delete VIP.
	err = lb.DelVIP(vip)
	if err != nil {
		t.Fatalf("DelVIP (IPv6) failed: %v", err)
	}
}

// TestVIPFlags tests VIP flag operations.
func TestVIPFlags(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.3",
		Port:    8080,
		Proto:   6,
	}

	// Add VIP with flags.
	initialFlags := uint32(0x01)
	err := lb.AddVIP(vip, initialFlags)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Get VIP flags.
	flags, err := lb.GetVIPFlags(vip)
	if err != nil {
		t.Fatalf("GetVIPFlags failed: %v", err)
	}
	if flags != initialFlags {
		t.Errorf("Expected flags %d, got %d", initialFlags, flags)
	}

	// Modify VIP flags.
	err = lb.ModifyVIP(vip, 0x02, true)
	if err != nil {
		t.Fatalf("ModifyVIP (set) failed: %v", err)
	}

	flags, err = lb.GetVIPFlags(vip)
	if err != nil {
		t.Fatalf("GetVIPFlags after modify failed: %v", err)
	}
	if flags&0x02 == 0 {
		t.Error("Flag 0x02 should be set")
	}

	// Clear flags.
	err = lb.ModifyVIP(vip, 0x02, false)
	if err != nil {
		t.Fatalf("ModifyVIP (clear) failed: %v", err)
	}

	flags, err = lb.GetVIPFlags(vip)
	if err != nil {
		t.Fatalf("GetVIPFlags after clear failed: %v", err)
	}
	if flags&0x02 != 0 {
		t.Error("Flag 0x02 should be cleared")
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestAddDeleteReal tests real server add and delete operations.
func TestAddDeleteReal(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.4",
		Port:    80,
		Proto:   6,
	}

	// Add VIP first.
	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	real := Real{
		Address: "192.168.1.1",
		Weight:  100,
		Flags:   0,
	}

	// Add real to VIP.
	err = lb.AddRealForVIP(real, vip)
	if err != nil {
		t.Fatalf("AddRealForVIP failed: %v", err)
	}

	// Get reals for VIP.
	reals, err := lb.GetRealsForVIP(vip)
	if err != nil {
		t.Fatalf("GetRealsForVIP failed: %v", err)
	}

	found := false
	for _, r := range reals {
		if r.Address == real.Address {
			found = true
			if r.Weight != real.Weight {
				t.Errorf("Expected weight %d, got %d", real.Weight, r.Weight)
			}
			break
		}
	}
	if !found {
		t.Error("Real not found after adding")
	}

	// Delete real from VIP.
	err = lb.DelRealForVIP(real, vip)
	if err != nil {
		t.Fatalf("DelRealForVIP failed: %v", err)
	}

	// Verify real is deleted.
	reals, err = lb.GetRealsForVIP(vip)
	if err != nil {
		t.Fatalf("GetRealsForVIP after delete failed: %v", err)
	}

	for _, r := range reals {
		if r.Address == real.Address {
			t.Error("Real still exists after delete")
		}
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestModifyRealsForVIP tests batch real modification.
func TestModifyRealsForVIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.5",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Batch add reals.
	reals := []Real{
		{Address: "192.168.1.10", Weight: 100, Flags: 0},
		{Address: "192.168.1.11", Weight: 200, Flags: 0},
		{Address: "192.168.1.12", Weight: 50, Flags: 0},
	}

	err = lb.ModifyRealsForVIP(ActionAdd, reals, vip)
	if err != nil {
		t.Fatalf("ModifyRealsForVIP (add) failed: %v", err)
	}

	// Verify all reals were added.
	gotReals, err := lb.GetRealsForVIP(vip)
	if err != nil {
		t.Fatalf("GetRealsForVIP failed: %v", err)
	}

	if len(gotReals) != len(reals) {
		t.Errorf("Expected %d reals, got %d", len(reals), len(gotReals))
	}

	// Batch delete reals.
	err = lb.ModifyRealsForVIP(ActionDel, reals, vip)
	if err != nil {
		t.Fatalf("ModifyRealsForVIP (del) failed: %v", err)
	}

	gotReals, err = lb.GetRealsForVIP(vip)
	if err != nil {
		t.Fatalf("GetRealsForVIP after delete failed: %v", err)
	}

	if len(gotReals) != 0 {
		t.Errorf("Expected 0 reals after delete, got %d", len(gotReals))
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestModifyRealsEmptySlice tests that ModifyRealsForVIP handles empty slice.
func TestModifyRealsEmptySlice(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.6",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Empty slice should not cause an error.
	err = lb.ModifyRealsForVIP(ActionAdd, []Real{}, vip)
	if err != nil {
		t.Errorf("ModifyRealsForVIP with empty slice failed: %v", err)
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestGetIndexForReal tests getting the index for a real server.
func TestGetIndexForReal(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.7",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	real := Real{
		Address: "192.168.2.1",
		Weight:  100,
		Flags:   0,
	}

	err = lb.AddRealForVIP(real, vip)
	if err != nil {
		t.Fatalf("AddRealForVIP failed: %v", err)
	}

	// Get index for real.
	index, err := lb.GetIndexForReal(real.Address)
	if err != nil {
		t.Fatalf("GetIndexForReal failed: %v", err)
	}

	if index < 0 {
		t.Errorf("Expected non-negative index, got %d", index)
	}

	// Cleanup.
	lb.DelRealForVIP(real, vip)
	lb.DelVIP(vip)
}

// TestGetIndexForNonExistentReal tests getting index for a non-existent real.
func TestGetIndexForNonExistentReal(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	index, err := lb.GetIndexForReal("192.168.99.99")
	if err == nil {
		t.Errorf("Expected error for non-existent real, got index %d", index)
	}
}

// TestHashFunctionChange tests changing the hash function for a VIP.
func TestHashFunctionChange(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.8",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Change hash function to MaglevV2.
	err = lb.ChangeHashFunctionForVIP(vip, HashMaglevV2)
	if err != nil {
		t.Fatalf("ChangeHashFunctionForVIP failed: %v", err)
	}

	// Change back to Maglev.
	err = lb.ChangeHashFunctionForVIP(vip, HashMaglev)
	if err != nil {
		t.Fatalf("ChangeHashFunctionForVIP (back) failed: %v", err)
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestMACOperations tests MAC address operations.
func TestMACOperations(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	// Get current MAC.
	originalMAC, err := lb.GetMAC()
	if err != nil {
		t.Fatalf("GetMAC failed: %v", err)
	}

	if len(originalMAC) != 6 {
		t.Errorf("Expected MAC length 6, got %d", len(originalMAC))
	}

	// Change MAC.
	newMAC := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	err = lb.ChangeMAC(newMAC)
	if err != nil {
		t.Fatalf("ChangeMAC failed: %v", err)
	}

	// Verify MAC was changed.
	gotMAC, err := lb.GetMAC()
	if err != nil {
		t.Fatalf("GetMAC after change failed: %v", err)
	}

	for i := 0; i < 6; i++ {
		if gotMAC[i] != newMAC[i] {
			t.Errorf("MAC byte %d: expected %02x, got %02x", i, newMAC[i], gotMAC[i])
		}
	}

	// Test invalid MAC length.
	err = lb.ChangeMAC([]byte{0x00, 0x11})
	if err == nil {
		t.Error("Expected error for invalid MAC length")
	}
}

// TestGetRealForFlow tests flow simulation.
func TestGetRealForFlow(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.9",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	real := Real{
		Address: "192.168.3.1",
		Weight:  100,
		Flags:   0,
	}

	err = lb.AddRealForVIP(real, vip)
	if err != nil {
		t.Fatalf("AddRealForVIP failed: %v", err)
	}

	flow := Flow{
		Src:     "172.16.0.1",
		Dst:     "10.0.0.9",
		SrcPort: 54321,
		DstPort: 80,
		Proto:   6,
	}

	// Get real for flow.
	gotReal, err := lb.GetRealForFlow(flow)
	if err != nil {
		t.Fatalf("GetRealForFlow failed: %v", err)
	}

	if gotReal != real.Address {
		t.Errorf("Expected real %s, got %s", real.Address, gotReal)
	}

	// Cleanup.
	lb.DelRealForVIP(real, vip)
	lb.DelVIP(vip)
}

// TestGetRealForFlowNoMatch tests flow with no matching VIP.
func TestGetRealForFlowNoMatch(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	flow := Flow{
		Src:     "172.16.0.1",
		Dst:     "10.99.99.99", // Non-existent VIP
		SrcPort: 54321,
		DstPort: 80,
		Proto:   6,
	}

	_, err := lb.GetRealForFlow(flow)
	if err == nil {
		t.Error("Expected error for flow with no matching VIP")
	}
	if !IsNotFound(err) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}

// TestStatsOperations tests various statistics retrieval.
func TestStatsOperations(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	// Test LRU stats.
	_, err := lb.GetLRUStats()
	if err != nil {
		t.Errorf("GetLRUStats failed: %v", err)
	}

	// Test LRU miss stats.
	_, err = lb.GetLRUMissStats()
	if err != nil {
		t.Errorf("GetLRUMissStats failed: %v", err)
	}

	// Test LRU fallback stats.
	_, err = lb.GetLRUFallbackStats()
	if err != nil {
		t.Errorf("GetLRUFallbackStats failed: %v", err)
	}

	// Test ICMP too big stats.
	_, err = lb.GetICMPTooBigStats()
	if err != nil {
		t.Errorf("GetICMPTooBigStats failed: %v", err)
	}

	// Test CH drop stats.
	_, err = lb.GetCHDropStats()
	if err != nil {
		t.Errorf("GetCHDropStats failed: %v", err)
	}

	// Test decap stats.
	_, err = lb.GetDecapStats()
	if err != nil {
		t.Errorf("GetDecapStats failed: %v", err)
	}

	// Test inline decap stats.
	_, err = lb.GetInlineDecapStats()
	if err != nil {
		t.Errorf("GetInlineDecapStats failed: %v", err)
	}

	// Test global LRU stats.
	_, err = lb.GetGlobalLRUStats()
	if err != nil {
		t.Errorf("GetGlobalLRUStats failed: %v", err)
	}

	// Test userspace stats.
	_, err = lb.GetUserspaceStats()
	if err != nil {
		t.Errorf("GetUserspaceStats failed: %v", err)
	}

	// Test per-core stats.
	_, err = lb.GetPerCorePacketsStats()
	if err != nil {
		t.Errorf("GetPerCorePacketsStats failed: %v", err)
	}
}

// TestVIPStats tests VIP-specific statistics.
func TestVIPStats(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.10",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Get stats for VIP.
	stats, err := lb.GetStatsForVIP(vip)
	if err != nil {
		t.Fatalf("GetStatsForVIP failed: %v", err)
	}

	// Stats should be zero initially.
	if stats.V1 != 0 || stats.V2 != 0 {
		t.Logf("Initial stats: V1=%d, V2=%d", stats.V1, stats.V2)
	}

	// Get decap stats for VIP.
	_, err = lb.GetDecapStatsForVIP(vip)
	if err != nil {
		t.Errorf("GetDecapStatsForVIP failed: %v", err)
	}

	// Cleanup.
	lb.DelVIP(vip)
}

// TestStatsForNonExistentVIP tests getting stats for a non-existent VIP.
func TestStatsForNonExistentVIP(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.99.99.99",
		Port:    80,
		Proto:   6,
	}

	_, err := lb.GetStatsForVIP(vip)
	if err == nil {
		t.Error("Expected error for non-existent VIP")
	}
	if !IsNotFound(err) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}

// TestQuicRealsMapping tests QUIC real server mapping operations.
func TestQuicRealsMapping(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	reals := []QuicReal{
		{Address: "192.168.10.1", ID: 1001},
		{Address: "192.168.10.2", ID: 1002},
	}

	// Add QUIC reals.
	err := lb.ModifyQuicRealsMapping(ActionAdd, reals)
	if err != nil {
		t.Fatalf("ModifyQuicRealsMapping (add) failed: %v", err)
	}

	// Get QUIC reals.
	gotReals, err := lb.GetQuicRealsMapping()
	if err != nil {
		t.Fatalf("GetQuicRealsMapping failed: %v", err)
	}

	if len(gotReals) < len(reals) {
		t.Errorf("Expected at least %d QUIC reals, got %d", len(reals), len(gotReals))
	}

	// Delete QUIC reals.
	err = lb.ModifyQuicRealsMapping(ActionDel, reals)
	if err != nil {
		t.Fatalf("ModifyQuicRealsMapping (del) failed: %v", err)
	}
}

// TestQuicRealsEmptySlice tests that ModifyQuicRealsMapping handles empty slice.
func TestQuicRealsEmptySlice(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	err := lb.ModifyQuicRealsMapping(ActionAdd, []QuicReal{})
	if err != nil {
		t.Errorf("ModifyQuicRealsMapping with empty slice failed: %v", err)
	}
}

// TestIsUnderFlood tests the flood detection function.
func TestIsUnderFlood(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	isFlood, err := lb.IsUnderFlood()
	if err != nil {
		t.Fatalf("IsUnderFlood failed: %v", err)
	}

	// Should not be under flood initially.
	if isFlood {
		t.Log("System reports being under flood")
	}
}

// TestAddSrcIPForPcktEncap tests setting the source IP for packet encapsulation.
func TestAddSrcIPForPcktEncap(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	// Add IPv4 source.
	err := lb.AddSrcIPForPcktEncap("10.1.1.1")
	if err != nil {
		t.Errorf("AddSrcIPForPcktEncap (IPv4) failed: %v", err)
	}

	// Add IPv6 source.
	err = lb.AddSrcIPForPcktEncap("2001:db8::100")
	if err != nil {
		t.Errorf("AddSrcIPForPcktEncap (IPv6) failed: %v", err)
	}
}

// TestLRUOperations tests LRU management operations.
func TestLRUOperations(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.11",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Try to delete LRU for a flow (may return empty if no entries).
	maps, err := lb.DeleteLRU(vip, "192.168.1.100", 12345)
	if err != nil {
		// This might fail if the flow doesn't exist, which is expected.
		t.Logf("DeleteLRU returned: %v", err)
	} else {
		t.Logf("DeleteLRU deleted from maps: %v", maps)
	}

	// Purge VIP LRU.
	count, err := lb.PurgeVIPLRU(vip)
	if err != nil {
		t.Errorf("PurgeVIPLRU failed: %v", err)
	}
	t.Logf("PurgeVIPLRU deleted %d entries", count)

	// Cleanup.
	lb.DelVIP(vip)
}

// TestMultipleLoadBalancers tests creating multiple LoadBalancer instances.
func TestMultipleLoadBalancers(t *testing.T) {
	lb1 := createTestLB(t)
	defer lb1.Close()

	lb2 := createTestLB(t)
	defer lb2.Close()

	// Both should work independently.
	vip1 := VIPKey{Address: "10.1.0.1", Port: 80, Proto: 6}
	vip2 := VIPKey{Address: "10.2.0.1", Port: 80, Proto: 6}

	err := lb1.AddVIP(vip1, 0)
	if err != nil {
		t.Fatalf("lb1.AddVIP failed: %v", err)
	}

	err = lb2.AddVIP(vip2, 0)
	if err != nil {
		t.Fatalf("lb2.AddVIP failed: %v", err)
	}

	// VIPs should be in their respective instances.
	vips1, _ := lb1.GetAllVIPs()
	vips2, _ := lb2.GetAllVIPs()

	t.Logf("lb1 has %d VIPs, lb2 has %d VIPs", len(vips1), len(vips2))
}

// TestConcurrentAccess tests that concurrent access is safe.
func TestConcurrentAccess(t *testing.T) {
	lb := createTestLB(t)
	defer lb.Close()

	vip := VIPKey{
		Address: "10.0.0.12",
		Port:    80,
		Proto:   6,
	}

	err := lb.AddVIP(vip, 0)
	if err != nil {
		t.Fatalf("AddVIP failed: %v", err)
	}

	// Run concurrent operations.
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 100; j++ {
				lb.GetAllVIPs()
				lb.GetVIPFlags(vip)
			}
		}()
	}

	// Wait for all goroutines.
	for i := 0; i < 10; i++ {
		<-done
	}

	// Cleanup.
	lb.DelVIP(vip)
}
