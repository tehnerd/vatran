//go:build integration
// +build integration

package katran

import (
	"testing"
)

// TestHealthcheckerDstOperations tests healthcheck destination operations.
func TestHealthcheckerDstOperations(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = true

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Add healthchecker destination.
	somark := uint32(100)
	dst := "192.168.100.1"

	err = lb.AddHealthcheckerDst(somark, dst)
	if err != nil {
		// HC might not be fully enabled in testing mode.
		if IsFeatureDisabled(err) {
			t.Skip("Healthcheck feature not enabled")
		}
		t.Fatalf("AddHealthcheckerDst failed: %v", err)
	}

	// Get healthchecker destinations.
	dsts, err := lb.GetHealthcheckersDst()
	if err != nil {
		t.Fatalf("GetHealthcheckersDst failed: %v", err)
	}

	found := false
	for _, d := range dsts {
		if d.Somark == somark && d.Dst == dst {
			found = true
			break
		}
	}
	if !found {
		t.Error("Healthchecker destination not found after adding")
	}

	// Delete healthchecker destination.
	err = lb.DelHealthcheckerDst(somark)
	if err != nil {
		t.Fatalf("DelHealthcheckerDst failed: %v", err)
	}

	// Verify deletion.
	dsts, err = lb.GetHealthcheckersDst()
	if err != nil {
		t.Fatalf("GetHealthcheckersDst after delete failed: %v", err)
	}

	for _, d := range dsts {
		if d.Somark == somark {
			t.Error("Healthchecker destination still exists after delete")
		}
	}
}

// TestDeleteNonExistentHealthcheckerDst tests deleting non-existent healthchecker destination.
func TestDeleteNonExistentHealthcheckerDst(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = true

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	err = lb.DelHealthcheckerDst(99999)
	if err == nil {
		t.Error("Expected error when deleting non-existent healthchecker dst")
	}
	if !IsNotFound(err) && !IsFeatureDisabled(err) {
		t.Errorf("Expected ErrNotFound or ErrFeatureDisabled, got: %v", err)
	}
}

// TestHCKeyOperations tests healthcheck key operations.
func TestHCKeyOperations(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = true

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	hcKey := VIPKey{
		Address: "10.10.10.1",
		Port:    8080,
		Proto:   6,
	}

	// Add HC key.
	err = lb.AddHCKey(hcKey)
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Healthcheck feature not enabled")
		}
		t.Fatalf("AddHCKey failed: %v", err)
	}

	// Delete HC key.
	err = lb.DelHCKey(hcKey)
	if err != nil {
		t.Fatalf("DelHCKey failed: %v", err)
	}
}

// TestGetHCProgStats tests healthcheck program statistics retrieval.
func TestGetHCProgStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = true

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	stats, err := lb.GetHCProgStats()
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Healthcheck feature not enabled")
		}
		t.Fatalf("GetHCProgStats failed: %v", err)
	}

	t.Logf("HC Stats: processed=%d, dropped=%d, skipped=%d, too_big=%d",
		stats.PacketsProcessed, stats.PacketsDropped,
		stats.PacketsSkipped, stats.PacketsTooBig)
}
