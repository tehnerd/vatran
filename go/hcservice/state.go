package hcservice

import (
	"fmt"
	"sync"
	"time"

	"github.com/tehnerd/vatran/go/server/types"
)

// VIPKey identifies a VIP by address, port, and protocol.
type VIPKey struct {
	Address string
	Port    uint16
	Proto   uint8
}

// String returns the VIP key as "address:port:proto".
func (k VIPKey) String() string {
	return fmt.Sprintf("%s:%d:%d", k.Address, k.Port, k.Proto)
}

// VIPKeyFromHC converts a types.HCVIPKey to a VIPKey.
//
// Parameters:
//   - hk: The HC VIP key to convert.
//
// Returns the corresponding VIPKey.
func VIPKeyFromHC(hk types.HCVIPKey) VIPKey {
	return VIPKey{Address: hk.Address, Port: hk.Port, Proto: hk.Proto}
}

// RealHealth tracks the health state and check history for a single real server.
type RealHealth struct {
	Address             string
	Weight              uint32
	Flags               uint8
	Healthy             bool
	ConsecutiveSuccess  int
	ConsecutiveFailures int
	LastCheckTime       time.Time
	LastStatusChange    time.Time
}

// VIPState holds the healthcheck configuration and real health states for a single VIP.
type VIPState struct {
	Key       VIPKey
	Config    types.HealthcheckConfig
	Reals     map[string]*RealHealth // realAddr -> health state
	CheckPort int                    // resolved port (config.Port or VIP port)
}

// CheckTarget represents a single check to be scheduled by the scheduler.
type CheckTarget struct {
	VIPKey    VIPKey
	VIPAddr   string
	RealAddr  string
	CheckPort int
	Config    types.HealthcheckConfig
}

// State is the thread-safe store for all VIP and real health tracking.
type State struct {
	mu   sync.RWMutex
	vips map[string]*VIPState // VIPKey.String() -> VIPState
}

// NewState creates a new empty State.
//
// Returns a new State instance.
func NewState() *State {
	return &State{
		vips: make(map[string]*VIPState),
	}
}

// RegisterVIP adds a new VIP with its healthcheck config and initial reals.
// All reals start as healthy.
//
// Parameters:
//   - key: The VIP key.
//   - config: The healthcheck configuration.
//   - reals: The initial set of real servers.
//
// Returns an error if the VIP is already registered.
func (s *State) RegisterVIP(key VIPKey, config types.HealthcheckConfig, reals []RealEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	if _, exists := s.vips[k]; exists {
		return fmt.Errorf("VIP %s already registered", k)
	}

	checkPort := config.Port
	if checkPort == 0 {
		checkPort = int(key.Port)
	}

	vs := &VIPState{
		Key:       key,
		Config:    config,
		Reals:     make(map[string]*RealHealth, len(reals)),
		CheckPort: checkPort,
	}

	now := time.Now()
	for _, r := range reals {
		vs.Reals[r.Address] = &RealHealth{
			Address:          r.Address,
			Weight:           r.Weight,
			Flags:            r.Flags,
			Healthy:          true,
			LastStatusChange: now,
		}
	}

	s.vips[k] = vs
	return nil
}

// UpdateVIP updates the healthcheck configuration for a registered VIP.
// Optionally replaces the reals list.
//
// Parameters:
//   - key: The VIP key.
//   - config: The new healthcheck configuration.
//   - reals: If non-nil, replaces the entire reals list. If nil, keeps existing reals.
//
// Returns the old real addresses (for somark cleanup) and an error if VIP not found.
func (s *State) UpdateVIP(key VIPKey, config types.HealthcheckConfig, reals []RealEntry) (oldReals []string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return nil, fmt.Errorf("VIP %s not found", k)
	}

	vs.Config = config
	checkPort := config.Port
	if checkPort == 0 {
		checkPort = int(key.Port)
	}
	vs.CheckPort = checkPort

	if reals != nil {
		// Collect old real addresses for cleanup
		for addr := range vs.Reals {
			oldReals = append(oldReals, addr)
		}
		// Replace with new reals
		vs.Reals = make(map[string]*RealHealth, len(reals))
		now := time.Now()
		for _, r := range reals {
			vs.Reals[r.Address] = &RealHealth{
				Address:          r.Address,
				Weight:           r.Weight,
				Flags:            r.Flags,
				Healthy:          true,
				LastStatusChange: now,
			}
		}
	}

	return oldReals, nil
}

// DeregisterVIP removes a VIP and all its reals from tracking.
//
// Parameters:
//   - key: The VIP key to remove.
//
// Returns the list of real addresses that were tracked (for somark cleanup) and an error if not found.
func (s *State) DeregisterVIP(key VIPKey) (realAddrs []string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return nil, fmt.Errorf("VIP %s not found", k)
	}

	for addr := range vs.Reals {
		realAddrs = append(realAddrs, addr)
	}
	delete(s.vips, k)
	return realAddrs, nil
}

// AddReals adds reals to a registered VIP. Already existing reals are skipped.
//
// Parameters:
//   - key: The VIP key.
//   - reals: The reals to add.
//
// Returns the count of added and skipped reals, and an error if VIP not found.
func (s *State) AddReals(key VIPKey, reals []RealEntry) (added, skipped int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return 0, 0, fmt.Errorf("VIP %s not found", k)
	}

	now := time.Now()
	for _, r := range reals {
		if _, exists := vs.Reals[r.Address]; exists {
			skipped++
			continue
		}
		vs.Reals[r.Address] = &RealHealth{
			Address:          r.Address,
			Weight:           r.Weight,
			Flags:            r.Flags,
			Healthy:          true,
			LastStatusChange: now,
		}
		added++
	}
	return added, skipped, nil
}

// RemoveReals removes reals from a registered VIP.
//
// Parameters:
//   - key: The VIP key.
//   - addresses: The real addresses to remove.
//
// Returns the count of removed and not-found reals, and an error if VIP not found.
func (s *State) RemoveReals(key VIPKey, addresses []string) (removed, notFound int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return 0, 0, fmt.Errorf("VIP %s not found", k)
	}

	for _, addr := range addresses {
		if _, exists := vs.Reals[addr]; exists {
			delete(vs.Reals, addr)
			removed++
		} else {
			notFound++
		}
	}
	return removed, notFound, nil
}

// UpdateRealHealth updates the health state for a real after a check completes.
//
// Parameters:
//   - key: The VIP key.
//   - realAddr: The real address.
//   - success: Whether the health check succeeded.
func (s *State) UpdateRealHealth(key VIPKey, realAddr string, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return
	}
	rh, exists := vs.Reals[realAddr]
	if !exists {
		return
	}

	now := time.Now()
	rh.LastCheckTime = now

	if success {
		rh.ConsecutiveFailures = 0
		rh.ConsecutiveSuccess++
		if !rh.Healthy && rh.ConsecutiveSuccess >= vs.Config.HealthyThreshold {
			rh.Healthy = true
			rh.LastStatusChange = now
		}
	} else {
		rh.ConsecutiveSuccess = 0
		rh.ConsecutiveFailures++
		if rh.Healthy && rh.ConsecutiveFailures >= vs.Config.UnhealthyThreshold {
			rh.Healthy = false
			rh.LastStatusChange = now
		}
	}
}

// GetAllCheckTargets returns all non-dummy check targets for scheduling.
//
// Returns a slice of CheckTarget for the scheduler to iterate.
func (s *State) GetAllCheckTargets() []CheckTarget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var targets []CheckTarget
	for _, vs := range s.vips {
		if vs.Config.Type == "dummy" {
			continue
		}
		for _, rh := range vs.Reals {
			targets = append(targets, CheckTarget{
				VIPKey:    vs.Key,
				VIPAddr:   vs.Key.Address,
				RealAddr:  rh.Address,
				CheckPort: vs.CheckPort,
				Config:    vs.Config,
			})
		}
	}
	return targets
}

// GetVIPHealth returns the health snapshot for a single VIP.
//
// Parameters:
//   - key: The VIP key.
//
// Returns the HCVIPHealthResponse or an error if not found.
func (s *State) GetVIPHealth(key VIPKey) (*types.HCVIPHealthResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return nil, fmt.Errorf("VIP %s not found", k)
	}

	resp := &types.HCVIPHealthResponse{
		VIP: types.HCVIPKey{
			Address: vs.Key.Address,
			Port:    vs.Key.Port,
			Proto:   vs.Key.Proto,
		},
		Reals: make([]types.HCRealHealth, 0, len(vs.Reals)),
	}

	for _, rh := range vs.Reals {
		entry := types.HCRealHealth{
			Address:             rh.Address,
			Healthy:             rh.Healthy,
			ConsecutiveFailures: rh.ConsecutiveFailures,
		}
		if !rh.LastCheckTime.IsZero() {
			entry.LastCheckTime = rh.LastCheckTime.UTC().Format(time.RFC3339)
		}
		if !rh.LastStatusChange.IsZero() {
			entry.LastStatusChange = rh.LastStatusChange.UTC().Format(time.RFC3339)
		}
		resp.Reals = append(resp.Reals, entry)
	}

	return resp, nil
}

// GetAllHealth returns the health snapshot for all registered VIPs.
//
// Returns a slice of HCVIPHealthResponse.
func (s *State) GetAllHealth() []types.HCVIPHealthResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]types.HCVIPHealthResponse, 0, len(s.vips))
	for _, vs := range s.vips {
		resp := types.HCVIPHealthResponse{
			VIP: types.HCVIPKey{
				Address: vs.Key.Address,
				Port:    vs.Key.Port,
				Proto:   vs.Key.Proto,
			},
			Reals: make([]types.HCRealHealth, 0, len(vs.Reals)),
		}
		for _, rh := range vs.Reals {
			entry := types.HCRealHealth{
				Address:             rh.Address,
				Healthy:             rh.Healthy,
				ConsecutiveFailures: rh.ConsecutiveFailures,
			}
			if !rh.LastCheckTime.IsZero() {
				entry.LastCheckTime = rh.LastCheckTime.UTC().Format(time.RFC3339)
			}
			if !rh.LastStatusChange.IsZero() {
				entry.LastStatusChange = rh.LastStatusChange.UTC().Format(time.RFC3339)
			}
			resp.Reals = append(resp.Reals, entry)
		}
		result = append(result, resp)
	}
	return result
}

// ListVIPs returns a summary of all registered VIPs.
//
// Returns a list of VIP summaries with config and real count.
func (s *State) ListVIPs() []VIPSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]VIPSummary, 0, len(s.vips))
	for _, vs := range s.vips {
		result = append(result, VIPSummary{
			VIP: types.HCVIPKey{
				Address: vs.Key.Address,
				Port:    vs.Key.Port,
				Proto:   vs.Key.Proto,
			},
			Healthcheck: vs.Config,
			RealCount:   len(vs.Reals),
		})
	}
	return result
}

// GetRealAddresses returns the list of real addresses for a VIP.
//
// Parameters:
//   - key: The VIP key.
//
// Returns the list of real addresses or an error if VIP not found.
func (s *State) GetRealAddresses(key VIPKey) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	k := key.String()
	vs, exists := s.vips[k]
	if !exists {
		return nil, fmt.Errorf("VIP %s not found", k)
	}

	addrs := make([]string, 0, len(vs.Reals))
	for addr := range vs.Reals {
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// VIPSummary is the list-view of a registered VIP.
type VIPSummary struct {
	VIP         types.HCVIPKey          `json:"vip"`
	Healthcheck types.HealthcheckConfig `json:"healthcheck"`
	RealCount   int                     `json:"real_count"`
}

// RealEntry represents a real server in registration requests.
type RealEntry struct {
	Address string `json:"address"`
	Weight  uint32 `json:"weight,omitempty"`
	Flags   uint8  `json:"flags,omitempty"`
}
