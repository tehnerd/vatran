//go:build integration
// +build integration

package katran

import (
	"testing"
)

// TestHasFeature tests the HasFeature function.
func TestHasFeature(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Test each feature.
	features := []struct {
		name    string
		feature Feature
	}{
		{"SrcRouting", FeatureSrcRouting},
		{"InlineDecap", FeatureInlineDecap},
		{"Introspection", FeatureIntrospection},
		{"GUEEncap", FeatureGUEEncap},
		{"DirectHC", FeatureDirectHC},
		{"LocalDeliveryOpt", FeatureLocalDeliveryOpt},
		{"FlowDebug", FeatureFlowDebug},
	}

	for _, tt := range features {
		t.Run(tt.name, func(t *testing.T) {
			has, err := lb.HasFeature(tt.feature)
			if err != nil {
				t.Errorf("HasFeature(%s) failed: %v", tt.name, err)
				return
			}
			t.Logf("Feature %s: %v", tt.name, has)
		})
	}
}

// TestHasFeatureAfterClose tests that HasFeature fails after Close.
func TestHasFeatureAfterClose(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	lb.Close()

	_, err = lb.HasFeature(FeatureSrcRouting)
	if err == nil {
		t.Error("HasFeature should fail after Close")
	}
}

// TestInstallFeature tests the InstallFeature function.
// Note: This test may fail if the BPF program path is not available.
func TestInstallFeature(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Try to install a feature (this will likely fail in testing mode
	// due to missing BPF programs, but we test the API).
	err = lb.InstallFeature(FeatureFlowDebug, "")
	if err != nil {
		// Expected to fail without actual BPF program.
		t.Logf("InstallFeature failed as expected: %v", err)
	}
}

// TestRemoveFeature tests the RemoveFeature function.
func TestRemoveFeature(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Try to remove a feature (this will likely fail in testing mode).
	err = lb.RemoveFeature(FeatureFlowDebug, "")
	if err != nil {
		// Expected to fail without actual BPF program.
		t.Logf("RemoveFeature failed as expected: %v", err)
	}
}

// TestFeatureCombinations tests combining multiple features.
func TestFeatureCombinations(t *testing.T) {
	// Test that features can be combined with bitwise operations.
	combined := FeatureSrcRouting | FeatureInlineDecap | FeatureGUEEncap

	if combined&FeatureSrcRouting == 0 {
		t.Error("Combined features should include SrcRouting")
	}
	if combined&FeatureInlineDecap == 0 {
		t.Error("Combined features should include InlineDecap")
	}
	if combined&FeatureGUEEncap == 0 {
		t.Error("Combined features should include GUEEncap")
	}
	if combined&FeatureIntrospection != 0 {
		t.Error("Combined features should not include Introspection")
	}
}
