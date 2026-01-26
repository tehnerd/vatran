//go:build integration
// +build integration

package katran

import (
	"testing"
)

// TestMonitorOperations tests monitor start/stop operations.
func TestMonitorOperations(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Check if introspection feature is available.
	hasFeature, err := lb.HasFeature(FeatureIntrospection)
	if err != nil {
		t.Fatalf("HasFeature failed: %v", err)
	}
	if !hasFeature {
		t.Skip("Introspection feature not enabled")
	}

	// Start monitor.
	err = lb.RestartMonitor(1000)
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Introspection feature not enabled")
		}
		t.Fatalf("RestartMonitor failed: %v", err)
	}

	// Get monitor stats.
	stats, err := lb.GetMonitorStats()
	if err != nil {
		t.Fatalf("GetMonitorStats failed: %v", err)
	}
	t.Logf("Monitor stats: limit=%d, amount=%d, buffer_full=%d",
		stats.Limit, stats.Amount, stats.BufferFull)

	// Stop monitor.
	err = lb.StopMonitor()
	if err != nil {
		t.Fatalf("StopMonitor failed: %v", err)
	}
}

// TestRestartMonitorWithZeroLimit tests restarting monitor with zero limit (unlimited).
func TestRestartMonitorWithZeroLimit(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Check if introspection feature is available.
	hasFeature, err := lb.HasFeature(FeatureIntrospection)
	if err != nil {
		t.Fatalf("HasFeature failed: %v", err)
	}
	if !hasFeature {
		t.Skip("Introspection feature not enabled")
	}

	// Start monitor with zero limit (unlimited).
	err = lb.RestartMonitor(0)
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Introspection feature not enabled")
		}
		t.Fatalf("RestartMonitor with zero limit failed: %v", err)
	}

	// Stop monitor.
	lb.StopMonitor()
}

// TestMonitorAfterClose tests that monitor operations fail after Close.
func TestMonitorAfterClose(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	lb.Close()

	err = lb.RestartMonitor(100)
	if err == nil {
		t.Error("RestartMonitor should fail after Close")
	}

	err = lb.StopMonitor()
	if err == nil {
		t.Error("StopMonitor should fail after Close")
	}

	_, err = lb.GetMonitorStats()
	if err == nil {
		t.Error("GetMonitorStats should fail after Close")
	}
}

// TestGetMonitorStatsStruct tests the MonitorStats struct fields.
func TestGetMonitorStatsStruct(t *testing.T) {
	// Test struct initialization.
	stats := MonitorStats{
		Limit:      1000,
		Amount:     500,
		BufferFull: 3,
	}

	if stats.Limit != 1000 {
		t.Errorf("Expected Limit 1000, got %d", stats.Limit)
	}
	if stats.Amount != 500 {
		t.Errorf("Expected Amount 500, got %d", stats.Amount)
	}
	if stats.BufferFull != 3 {
		t.Errorf("Expected BufferFull 3, got %d", stats.BufferFull)
	}
}
