package lb

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/tehnerd/vatran/go/server/types"
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
	hcConfigs             map[string]*types.HealthcheckConfig // vipKey -> HC config
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
		hcConfigs:             make(map[string]*types.HealthcheckConfig),
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
	delete(s.hcConfigs, vipKey)
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
	s.hcConfigs = make(map[string]*types.HealthcheckConfig)
}

// GetHealthcheckerEndpoint returns the configured healthchecker endpoint URL.
//
// Returns the healthchecker endpoint string.
func (s *VIPRealsState) GetHealthcheckerEndpoint() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthcheckerEndpoint
}

// SetHCConfig stores or updates the healthcheck configuration for a VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//   - config: The healthcheck configuration to store.
func (s *VIPRealsState) SetHCConfig(vipKey string, config *types.HealthcheckConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hcConfigs[vipKey] = config
}

// GetHCConfig retrieves the healthcheck configuration for a VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns the healthcheck config and true if found, or nil and false if not found.
func (s *VIPRealsState) GetHCConfig(vipKey string) (*types.HealthcheckConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.hcConfigs[vipKey]
	return cfg, ok
}

// DelHCConfig removes the healthcheck configuration for a VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns the removed config and true if found, or nil and false if not found.
func (s *VIPRealsState) DelHCConfig(vipKey string) (*types.HealthcheckConfig, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, ok := s.hcConfigs[vipKey]
	if ok {
		delete(s.hcConfigs, vipKey)
	}
	return cfg, ok
}

// GetAllHCConfigs returns a copy of all stored healthcheck configurations.
//
// Returns a map of vipKey to HealthcheckConfig.
func (s *VIPRealsState) GetAllHCConfigs() map[string]*types.HealthcheckConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]*types.HealthcheckConfig, len(s.hcConfigs))
	for k, v := range s.hcConfigs {
		result[k] = v
	}
	return result
}

// CountHealthyReals returns the number of healthy reals for a given VIP.
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns the count of healthy reals.
func (s *VIPRealsState) CountHealthyReals(vipKey string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	reals, exists := s.vips[vipKey]
	if !exists {
		return 0
	}
	count := 0
	for _, rs := range reals {
		if rs.Healthy {
			count++
		}
	}
	return count
}

// GetVIPAddress extracts the address portion from a canonical VIP key string.
// The key format is "address:port:proto" (e.g., "10.0.0.1:80:6").
//
// Parameters:
//   - vipKey: The canonical VIP key string.
//
// Returns the address portion of the key.
func GetVIPAddress(vipKey string) string {
	// Find the last two colons to extract the address
	// Handle IPv6 addresses which contain colons
	lastColon := -1
	secondLastColon := -1
	for i := len(vipKey) - 1; i >= 0; i-- {
		if vipKey[i] == ':' {
			if lastColon == -1 {
				lastColon = i
			} else {
				secondLastColon = i
				break
			}
		}
	}
	if secondLastColon >= 0 {
		return vipKey[:secondLastColon]
	}
	return vipKey
}

// ParseVIPKey parses a canonical VIP key string back into its components.
// The key format is "address:port:proto" (e.g., "10.0.0.1:80:6").
// IPv6 addresses with colons are handled correctly.
//
// Parameters:
//   - vipKey: The canonical VIP key string to parse.
//
// Returns the address, port, protocol, and any parsing error.
func ParseVIPKey(vipKey string) (address string, port uint16, proto uint8, err error) {
	// Find the last two colons (handles IPv6 addresses)
	lastColon := -1
	secondLastColon := -1
	for i := len(vipKey) - 1; i >= 0; i-- {
		if vipKey[i] == ':' {
			if lastColon == -1 {
				lastColon = i
			} else {
				secondLastColon = i
				break
			}
		}
	}
	if secondLastColon < 0 || lastColon < 0 {
		return "", 0, 0, fmt.Errorf("invalid VIP key format: %q", vipKey)
	}

	address = vipKey[:secondLastColon]
	portStr := vipKey[secondLastColon+1 : lastColon]
	protoStr := vipKey[lastColon+1:]

	p, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid port in VIP key %q: %w", vipKey, err)
	}
	pr, err := strconv.ParseUint(protoStr, 10, 8)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid proto in VIP key %q: %w", vipKey, err)
	}

	return address, uint16(p), uint8(pr), nil
}
