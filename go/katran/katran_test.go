package katran

import (
	"errors"
	"testing"
)

// TestErrorConstants verifies that error code constants are defined correctly.
func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		err      Error
		expected int
	}{
		{"OK", OK, 0},
		{"ErrInvalidArgument", ErrInvalidArgument, -1},
		{"ErrNotFound", ErrNotFound, -2},
		{"ErrAlreadyExists", ErrAlreadyExists, -3},
		{"ErrSpaceExhausted", ErrSpaceExhausted, -4},
		{"ErrBPFFailed", ErrBPFFailed, -5},
		{"ErrFeatureDisabled", ErrFeatureDisabled, -6},
		{"ErrInternal", ErrInternal, -7},
		{"ErrMemory", ErrMemory, -8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.err) != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, int(tt.err))
			}
		})
	}
}

// TestModifyActionConstants verifies that modify action constants are correct.
func TestModifyActionConstants(t *testing.T) {
	if ActionAdd != 0 {
		t.Errorf("ActionAdd expected 0, got %d", ActionAdd)
	}
	if ActionDel != 1 {
		t.Errorf("ActionDel expected 1, got %d", ActionDel)
	}
}

// TestHashFunctionConstants verifies that hash function constants are correct.
func TestHashFunctionConstants(t *testing.T) {
	if HashMaglev != 0 {
		t.Errorf("HashMaglev expected 0, got %d", HashMaglev)
	}
	if HashMaglevV2 != 1 {
		t.Errorf("HashMaglevV2 expected 1, got %d", HashMaglevV2)
	}
}

// TestFeatureConstants verifies that feature constants are correct bit flags.
func TestFeatureConstants(t *testing.T) {
	tests := []struct {
		name    string
		feature Feature
		bit     int
	}{
		{"FeatureSrcRouting", FeatureSrcRouting, 1 << 0},
		{"FeatureInlineDecap", FeatureInlineDecap, 1 << 1},
		{"FeatureIntrospection", FeatureIntrospection, 1 << 2},
		{"FeatureGUEEncap", FeatureGUEEncap, 1 << 3},
		{"FeatureDirectHC", FeatureDirectHC, 1 << 4},
		{"FeatureLocalDeliveryOpt", FeatureLocalDeliveryOpt, 1 << 5},
		{"FeatureFlowDebug", FeatureFlowDebug, 1 << 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.feature) != tt.bit {
				t.Errorf("expected %d, got %d", tt.bit, int(tt.feature))
			}
		})
	}

	// Verify features can be combined with bitwise OR.
	combined := FeatureSrcRouting | FeatureInlineDecap
	if combined != Feature(0x03) {
		t.Errorf("combined features expected 0x03, got 0x%x", combined)
	}
}

// TestKatranError tests the KatranError type.
func TestKatranError(t *testing.T) {
	err := &KatranError{
		Code:    ErrNotFound,
		Message: "VIP not found",
	}

	// Test Error() method.
	errStr := err.Error()
	if errStr != "katran: VIP not found (code: -2)" {
		t.Errorf("unexpected error string: %s", errStr)
	}

	// Test Is() method with matching error.
	target := &KatranError{Code: ErrNotFound}
	if !errors.Is(err, target) {
		t.Error("expected errors.Is to return true for matching code")
	}

	// Test Is() method with non-matching error.
	target2 := &KatranError{Code: ErrAlreadyExists}
	if errors.Is(err, target2) {
		t.Error("expected errors.Is to return false for different code")
	}

	// Test Is() method with non-KatranError.
	if err.Is(errors.New("other error")) {
		t.Error("expected Is to return false for non-KatranError")
	}
}

// TestIsNotFound tests the IsNotFound helper function.
func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "KatranError with ErrNotFound",
			err:      &KatranError{Code: ErrNotFound, Message: "not found"},
			expected: true,
		},
		{
			name:     "KatranError with different code",
			err:      &KatranError{Code: ErrAlreadyExists, Message: "exists"},
			expected: false,
		},
		{
			name:     "non-KatranError",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsNotFound(tt.err); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIsAlreadyExists tests the IsAlreadyExists helper function.
func TestIsAlreadyExists(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "KatranError with ErrAlreadyExists",
			err:      &KatranError{Code: ErrAlreadyExists, Message: "exists"},
			expected: true,
		},
		{
			name:     "KatranError with different code",
			err:      &KatranError{Code: ErrNotFound, Message: "not found"},
			expected: false,
		},
		{
			name:     "non-KatranError",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsAlreadyExists(tt.err); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIsSpaceExhausted tests the IsSpaceExhausted helper function.
func TestIsSpaceExhausted(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "KatranError with ErrSpaceExhausted",
			err:      &KatranError{Code: ErrSpaceExhausted, Message: "full"},
			expected: true,
		},
		{
			name:     "KatranError with different code",
			err:      &KatranError{Code: ErrNotFound, Message: "not found"},
			expected: false,
		},
		{
			name:     "non-KatranError",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsSpaceExhausted(tt.err); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIsBPFFailed tests the IsBPFFailed helper function.
func TestIsBPFFailed(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "KatranError with ErrBPFFailed",
			err:      &KatranError{Code: ErrBPFFailed, Message: "bpf error"},
			expected: true,
		},
		{
			name:     "KatranError with different code",
			err:      &KatranError{Code: ErrNotFound, Message: "not found"},
			expected: false,
		},
		{
			name:     "non-KatranError",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsBPFFailed(tt.err); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIsFeatureDisabled tests the IsFeatureDisabled helper function.
func TestIsFeatureDisabled(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "KatranError with ErrFeatureDisabled",
			err:      &KatranError{Code: ErrFeatureDisabled, Message: "disabled"},
			expected: true,
		},
		{
			name:     "KatranError with different code",
			err:      &KatranError{Code: ErrNotFound, Message: "not found"},
			expected: false,
		},
		{
			name:     "non-KatranError",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsFeatureDisabled(tt.err); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestVIPKey tests the VIPKey struct.
func TestVIPKey(t *testing.T) {
	vip := VIPKey{
		Address: "10.0.0.1",
		Port:    80,
		Proto:   6, // TCP
	}

	if vip.Address != "10.0.0.1" {
		t.Errorf("expected address 10.0.0.1, got %s", vip.Address)
	}
	if vip.Port != 80 {
		t.Errorf("expected port 80, got %d", vip.Port)
	}
	if vip.Proto != 6 {
		t.Errorf("expected proto 6, got %d", vip.Proto)
	}

	// Test IPv6 VIP.
	vip6 := VIPKey{
		Address: "2001:db8::1",
		Port:    443,
		Proto:   6,
	}
	if vip6.Address != "2001:db8::1" {
		t.Errorf("expected IPv6 address 2001:db8::1, got %s", vip6.Address)
	}
}

// TestReal tests the Real struct.
func TestReal(t *testing.T) {
	real := Real{
		Address: "192.168.1.1",
		Weight:  100,
		Flags:   0,
	}

	if real.Address != "192.168.1.1" {
		t.Errorf("expected address 192.168.1.1, got %s", real.Address)
	}
	if real.Weight != 100 {
		t.Errorf("expected weight 100, got %d", real.Weight)
	}
	if real.Flags != 0 {
		t.Errorf("expected flags 0, got %d", real.Flags)
	}
}

// TestQuicReal tests the QuicReal struct.
func TestQuicReal(t *testing.T) {
	qr := QuicReal{
		Address: "10.0.0.10",
		ID:      12345,
	}

	if qr.Address != "10.0.0.10" {
		t.Errorf("expected address 10.0.0.10, got %s", qr.Address)
	}
	if qr.ID != 12345 {
		t.Errorf("expected ID 12345, got %d", qr.ID)
	}
}

// TestFlow tests the Flow struct.
func TestFlow(t *testing.T) {
	flow := Flow{
		Src:     "192.168.1.100",
		Dst:     "10.0.0.1",
		SrcPort: 54321,
		DstPort: 80,
		Proto:   6,
	}

	if flow.Src != "192.168.1.100" {
		t.Errorf("expected src 192.168.1.100, got %s", flow.Src)
	}
	if flow.Dst != "10.0.0.1" {
		t.Errorf("expected dst 10.0.0.1, got %s", flow.Dst)
	}
	if flow.SrcPort != 54321 {
		t.Errorf("expected src_port 54321, got %d", flow.SrcPort)
	}
	if flow.DstPort != 80 {
		t.Errorf("expected dst_port 80, got %d", flow.DstPort)
	}
	if flow.Proto != 6 {
		t.Errorf("expected proto 6, got %d", flow.Proto)
	}
}

// TestLBStats tests the LBStats struct.
func TestLBStats(t *testing.T) {
	stats := LBStats{
		V1: 1000000,
		V2: 500000000,
	}

	if stats.V1 != 1000000 {
		t.Errorf("expected V1 1000000, got %d", stats.V1)
	}
	if stats.V2 != 500000000 {
		t.Errorf("expected V2 500000000, got %d", stats.V2)
	}
}

// TestQuicPacketsStats tests the QuicPacketsStats struct.
func TestQuicPacketsStats(t *testing.T) {
	stats := QuicPacketsStats{
		CHRouted:                 100,
		CIDInitial:               50,
		CIDInvalidServerID:       5,
		CIDInvalidServerIDSample: 1,
		CIDRouted:                200,
		CIDUnknownRealDropped:    2,
		CIDV0:                    10,
		CIDV1:                    20,
		CIDV2:                    30,
		CIDV3:                    40,
		DstMatchInLRU:            80,
		DstMismatchInLRU:         3,
		DstNotFoundInLRU:         7,
	}

	if stats.CHRouted != 100 {
		t.Errorf("expected CHRouted 100, got %d", stats.CHRouted)
	}
	if stats.CIDRouted != 200 {
		t.Errorf("expected CIDRouted 200, got %d", stats.CIDRouted)
	}
}

// TestTPRPacketsStats tests the TPRPacketsStats struct.
func TestTPRPacketsStats(t *testing.T) {
	stats := TPRPacketsStats{
		CHRouted:         100,
		DstMismatchInLRU: 5,
		SIDRouted:        200,
		TCPSyn:           50,
	}

	if stats.CHRouted != 100 {
		t.Errorf("expected CHRouted 100, got %d", stats.CHRouted)
	}
	if stats.SIDRouted != 200 {
		t.Errorf("expected SIDRouted 200, got %d", stats.SIDRouted)
	}
}

// TestHCStats tests the HCStats struct.
func TestHCStats(t *testing.T) {
	stats := HCStats{
		PacketsProcessed: 10000,
		PacketsDropped:   10,
		PacketsSkipped:   5,
		PacketsTooBig:    2,
	}

	if stats.PacketsProcessed != 10000 {
		t.Errorf("expected PacketsProcessed 10000, got %d", stats.PacketsProcessed)
	}
	if stats.PacketsDropped != 10 {
		t.Errorf("expected PacketsDropped 10, got %d", stats.PacketsDropped)
	}
}

// TestBPFMapStats tests the BPFMapStats struct.
func TestBPFMapStats(t *testing.T) {
	stats := BPFMapStats{
		MaxEntries:     65536,
		CurrentEntries: 1024,
	}

	if stats.MaxEntries != 65536 {
		t.Errorf("expected MaxEntries 65536, got %d", stats.MaxEntries)
	}
	if stats.CurrentEntries != 1024 {
		t.Errorf("expected CurrentEntries 1024, got %d", stats.CurrentEntries)
	}
}

// TestMonitorStats tests the MonitorStats struct.
func TestMonitorStats(t *testing.T) {
	stats := MonitorStats{
		Limit:      1000,
		Amount:     500,
		BufferFull: 3,
	}

	if stats.Limit != 1000 {
		t.Errorf("expected Limit 1000, got %d", stats.Limit)
	}
	if stats.Amount != 500 {
		t.Errorf("expected Amount 500, got %d", stats.Amount)
	}
}

// TestUserspaceStats tests the UserspaceStats struct.
func TestUserspaceStats(t *testing.T) {
	stats := UserspaceStats{
		BPFFailedCalls:       5,
		AddrValidationFailed: 10,
	}

	if stats.BPFFailedCalls != 5 {
		t.Errorf("expected BPFFailedCalls 5, got %d", stats.BPFFailedCalls)
	}
	if stats.AddrValidationFailed != 10 {
		t.Errorf("expected AddrValidationFailed 10, got %d", stats.AddrValidationFailed)
	}
}

// TestSrcRoutingRule tests the SrcRoutingRule struct.
func TestSrcRoutingRule(t *testing.T) {
	rule := SrcRoutingRule{
		Src: "10.0.0.0/8",
		Dst: "192.168.1.1",
	}

	if rule.Src != "10.0.0.0/8" {
		t.Errorf("expected Src 10.0.0.0/8, got %s", rule.Src)
	}
	if rule.Dst != "192.168.1.1" {
		t.Errorf("expected Dst 192.168.1.1, got %s", rule.Dst)
	}
}

// TestHealthcheckerDst tests the HealthcheckerDst struct.
func TestHealthcheckerDst(t *testing.T) {
	hc := HealthcheckerDst{
		Somark: 100,
		Dst:    "10.0.0.5",
	}

	if hc.Somark != 100 {
		t.Errorf("expected Somark 100, got %d", hc.Somark)
	}
	if hc.Dst != "10.0.0.5" {
		t.Errorf("expected Dst 10.0.0.5, got %s", hc.Dst)
	}
}

// TestNewWithNilConfig tests that New returns an error for nil config.
func TestNewWithNilConfig(t *testing.T) {
	lb, err := New(nil)
	if err == nil {
		t.Error("expected error for nil config")
		if lb != nil {
			lb.Close()
		}
		return
	}

	if !IsNotFound(err) {
		katranErr, ok := err.(*KatranError)
		if !ok {
			t.Errorf("expected KatranError, got %T", err)
			return
		}
		if katranErr.Code != ErrInvalidArgument {
			t.Errorf("expected ErrInvalidArgument, got %d", katranErr.Code)
		}
	}
}
