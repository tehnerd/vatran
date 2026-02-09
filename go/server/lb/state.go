package lb

import (
	"fmt"
	"sync"
)

// RealState represents the tracked state of a real server for a VIP.
type RealState struct {
	// Address is the IP address of the real server.
	Address string
	// Weight is the weight for consistent hashing.
	Weight uint32
	// Flags contains real-specific flags.
	Flags uint8
	// Healthy indicates whether the real is currently healthy and receiving traffic.
	Healthy bool
}

// VIPRealsState is a thread-safe state store that tracks health state per real per VIP.
// Katran's C library only knows about healthy reals, while this store tracks all reals.
type VIPRealsState struct {
	mu                    sync.RWMutex
	vips                  map[string]map[string]*RealState // vipKey -> realAddr -> state
	healthcheckerEndpoint string
}

// VIPKeyString builds a canonical string key for a VIP from its address, port, and protocol.
//
// Parameters:
//   - address: The IP address of the VIP.
//   - port: The port number.
//   - proto: The IP protocol number.
//
// Returns a canonical key string like "10.0.0.1:80:6".
func VIPKeyString(address string, port uint16, proto uint8) string {
	return fmt.Sprintf("%s:%d:%d", address, port, proto)
}

// NewVIPRealsState creates a new VIPRealsState with the given healthchecker endpoint.
//
// Parameters:
//   - healthcheckerEndpoint: The URL of the healthchecker service. If empty,
//     newly added reals default to healthy.
//
// Returns a new VIPRealsState instance.
func NewVIPRealsState(healthcheckerEndpoint string) *VIPRealsState {
	return &VIPRealsState{
		vips:                  make(map[string]map[string]*RealState),
		healthcheckerEndpoint: healthcheckerEndpoint,
	}
}

// DefaultHealthy returns the default health state for newly added reals.
// When no healthchecker endpoint is configured, reals default to healthy.
// When a healthchecker endpoint is configured, reals default to unhealthy
// until the healthchecker reports them as healthy.
//
// Returns true if reals should default to healthy.
func (s *VIPRealsState) DefaultHealthy() bool {
	return len(s.healthcheckerEndpoint) == 0
}

// InitVIP initializes tracking for a VIP. If the VIP is already tracked,
// this is a no-op.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
func (s *VIPRealsState) InitVIP(vipKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.vips[vipKey]; !exists {
		s.vips[vipKey] = make(map[string]*RealState)
	}
}

// CleanVIP removes all state for a VIP and returns the list of previously
// healthy reals (which were in katran).
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns a slice of RealState for reals that were healthy.
func (s *VIPRealsState) CleanVIP(vipKey string) []RealState {
	s.mu.Lock()
	defer s.mu.Unlock()
	reals, exists := s.vips[vipKey]
	if !exists {
		return nil
	}
	var healthy []RealState
	for _, rs := range reals {
		if rs.Healthy {
			healthy = append(healthy, *rs)
		}
	}
	delete(s.vips, vipKey)
	return healthy
}

// AddReal adds a real to a VIP with the default health state.
// If the VIP is not yet tracked, it auto-initializes the VIP entry.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//   - address: The real server IP address.
//   - weight: The weight for consistent hashing.
//   - flags: Real-specific flags.
//
// Returns the new RealState.
func (s *VIPRealsState) AddReal(vipKey, address string, weight uint32, flags uint8) *RealState {
	return s.AddRealWithHealth(vipKey, address, weight, flags, s.DefaultHealthy())
}

// AddRealWithHealth adds a real to a VIP with an explicit health state.
// If the VIP is not yet tracked, it auto-initializes the VIP entry.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//   - address: The real server IP address.
//   - weight: The weight for consistent hashing.
//   - flags: Real-specific flags.
//   - healthy: Whether the real is initially healthy.
//
// Returns the new RealState.
func (s *VIPRealsState) AddRealWithHealth(vipKey, address string, weight uint32, flags uint8, healthy bool) *RealState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.vips[vipKey]; !exists {
		s.vips[vipKey] = make(map[string]*RealState)
	}
	rs := &RealState{
		Address: address,
		Weight:  weight,
		Flags:   flags,
		Healthy: healthy,
	}
	s.vips[vipKey][address] = rs
	return rs
}

// DelReal removes a real from a VIP's state.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//   - address: The real server IP address.
//
// Returns the removed RealState and true if found, or nil and false if not found.
func (s *VIPRealsState) DelReal(vipKey, address string) (*RealState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	reals, exists := s.vips[vipKey]
	if !exists {
		return nil, false
	}
	rs, found := reals[address]
	if !found {
		return nil, false
	}
	delete(reals, address)
	return rs, true
}

// GetReals returns all reals (healthy and unhealthy) for a VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns a slice of all RealState entries, or nil if VIP is not tracked.
func (s *VIPRealsState) GetReals(vipKey string) []RealState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	reals, exists := s.vips[vipKey]
	if !exists {
		return nil
	}
	result := make([]RealState, 0, len(reals))
	for _, rs := range reals {
		result = append(result, *rs)
	}
	return result
}

// UpdateHealth updates the health state of a real for a VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//   - address: The real server IP address.
//   - healthy: The new health state.
//
// Returns the old health state and whether the real was found.
func (s *VIPRealsState) UpdateHealth(vipKey, address string, healthy bool) (oldHealthy bool, found bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	reals, exists := s.vips[vipKey]
	if !exists {
		return false, false
	}
	rs, ok := reals[address]
	if !ok {
		return false, false
	}
	oldHealthy = rs.Healthy
	rs.Healthy = healthy
	return oldHealthy, true
}

// Clear removes all tracked state. Called on LB close.
func (s *VIPRealsState) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vips = make(map[string]map[string]*RealState)
}

// GetHealthcheckerEndpoint returns the configured healthchecker endpoint URL.
//
// Returns the healthchecker endpoint string.
func (s *VIPRealsState) GetHealthcheckerEndpoint() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthcheckerEndpoint
}
