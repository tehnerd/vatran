//go:build integration
// +build integration

package katran

import (
	"testing"
)

// TestSrcRoutingRuleOperations tests source routing rule operations.
func TestSrcRoutingRuleOperations(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Check if feature is available.
	hasFeature, err := lb.HasFeature(FeatureSrcRouting)
	if err != nil {
		t.Fatalf("HasFeature failed: %v", err)
	}
	if !hasFeature {
		t.Skip("Source routing feature not enabled")
	}

	srcPrefixes := []string{"10.0.0.0/8", "172.16.0.0/12"}
	dst := "192.168.1.1"

	// Add source routing rule.
	err = lb.AddSrcRoutingRule(srcPrefixes, dst)
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Source routing feature not enabled")
		}
		t.Fatalf("AddSrcRoutingRule failed: %v", err)
	}

	// Get source routing rules.
	rules, err := lb.GetSrcRoutingRules()
	if err != nil {
		t.Fatalf("GetSrcRoutingRules failed: %v", err)
	}

	t.Logf("Got %d source routing rules", len(rules))
	for _, rule := range rules {
		t.Logf("  Rule: %s -> %s", rule.Src, rule.Dst)
	}

	// Get rule count.
	count, err := lb.GetSrcRoutingRuleSize()
	if err != nil {
		t.Fatalf("GetSrcRoutingRuleSize failed: %v", err)
	}
	t.Logf("Source routing rule count: %d", count)

	// Delete source routing rules.
	err = lb.DelSrcRoutingRule(srcPrefixes)
	if err != nil {
		t.Fatalf("DelSrcRoutingRule failed: %v", err)
	}
}

// TestSrcRoutingEmptyPrefixes tests that empty prefixes are handled correctly.
func TestSrcRoutingEmptyPrefixes(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Empty prefixes should not cause error.
	err = lb.AddSrcRoutingRule([]string{}, "10.0.0.1")
	if err != nil {
		t.Errorf("AddSrcRoutingRule with empty prefixes failed: %v", err)
	}

	err = lb.DelSrcRoutingRule([]string{})
	if err != nil {
		t.Errorf("DelSrcRoutingRule with empty prefixes failed: %v", err)
	}
}

// TestClearAllSrcRoutingRules tests clearing all source routing rules.
func TestClearAllSrcRoutingRules(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Check if feature is available.
	hasFeature, err := lb.HasFeature(FeatureSrcRouting)
	if err != nil {
		t.Fatalf("HasFeature failed: %v", err)
	}
	if !hasFeature {
		t.Skip("Source routing feature not enabled")
	}

	// Add some rules first.
	err = lb.AddSrcRoutingRule([]string{"10.0.0.0/8"}, "192.168.1.1")
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Source routing feature not enabled")
		}
		t.Fatalf("AddSrcRoutingRule failed: %v", err)
	}

	// Clear all rules.
	err = lb.ClearAllSrcRoutingRules()
	if err != nil {
		t.Fatalf("ClearAllSrcRoutingRules failed: %v", err)
	}

	// Verify rules are cleared.
	count, err := lb.GetSrcRoutingRuleSize()
	if err != nil {
		t.Fatalf("GetSrcRoutingRuleSize failed: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 rules after clear, got %d", count)
	}
}

// TestInlineDecapDstOperations tests inline decapsulation destination operations.
func TestInlineDecapDstOperations(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	// Check if feature is available.
	hasFeature, err := lb.HasFeature(FeatureInlineDecap)
	if err != nil {
		t.Fatalf("HasFeature failed: %v", err)
	}
	if !hasFeature {
		t.Skip("Inline decap feature not enabled")
	}

	dst := "10.20.30.40"

	// Add inline decap destination.
	err = lb.AddInlineDecapDst(dst)
	if err != nil {
		if IsFeatureDisabled(err) {
			t.Skip("Inline decap feature not enabled")
		}
		t.Fatalf("AddInlineDecapDst failed: %v", err)
	}

	// Get inline decap destinations.
	dsts, err := lb.GetInlineDecapDsts()
	if err != nil {
		t.Fatalf("GetInlineDecapDsts failed: %v", err)
	}

	found := false
	for _, d := range dsts {
		if d == dst {
			found = true
			break
		}
	}
	if !found {
		t.Error("Inline decap destination not found after adding")
	}

	// Delete inline decap destination.
	err = lb.DelInlineDecapDst(dst)
	if err != nil {
		t.Fatalf("DelInlineDecapDst failed: %v", err)
	}

	// Verify deletion.
	dsts, err = lb.GetInlineDecapDsts()
	if err != nil {
		t.Fatalf("GetInlineDecapDsts after delete failed: %v", err)
	}

	for _, d := range dsts {
		if d == dst {
			t.Error("Inline decap destination still exists after delete")
		}
	}
}

// TestDeleteNonExistentInlineDecapDst tests deleting non-existent inline decap destination.
func TestDeleteNonExistentInlineDecapDst(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	err = lb.DelInlineDecapDst("99.99.99.99")
	if err == nil {
		t.Error("Expected error when deleting non-existent inline decap dst")
	}
	if !IsNotFound(err) && !IsFeatureDisabled(err) {
		t.Errorf("Expected ErrNotFound or ErrFeatureDisabled, got: %v", err)
	}
}

// TestGetSrcRoutingStats tests source routing statistics retrieval.
func TestGetSrcRoutingStats(t *testing.T) {
	cfg := NewConfig()
	cfg.Testing = true
	cfg.MainInterface = "lo"
	cfg.EnableHC = false

	lb, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create LoadBalancer: %v", err)
	}
	defer lb.Close()

	stats, err := lb.GetSrcRoutingStats()
	if err != nil {
		t.Fatalf("GetSrcRoutingStats failed: %v", err)
	}

	t.Logf("Src routing stats: V1=%d (local), V2=%d (remote)", stats.V1, stats.V2)
}
