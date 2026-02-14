package bgpservice

import (
	"fmt"
	"sync"
	"time"
)

// RouteState represents the tracked state of an advertised VIP route.
type RouteState struct {
	// VIP is the VIP IP address.
	VIP string `json:"vip"`
	// PrefixLen is the prefix length (e.g., 32 for /32).
	PrefixLen uint8 `json:"prefix_len"`
	// Advertised indicates whether the route is currently announced via BGP.
	Advertised bool `json:"advertised"`
	// Since is the timestamp of the last state change.
	Since time.Time `json:"since"`
	// Communities is the list of BGP communities applied to this route.
	Communities []string `json:"communities"`
	// LocalPref is the local preference value for this route.
	LocalPref uint32 `json:"local_pref"`
}

// State is a thread-safe store tracking VIP advertisement state.
type State struct {
	mu     sync.RWMutex
	routes map[string]*RouteState // key: "vip/prefixlen" e.g. "10.0.0.1/32"
}

// NewState creates a new empty State.
//
// Returns a new State instance.
func NewState() *State {
	return &State{
		routes: make(map[string]*RouteState),
	}
}

// routeKey builds the map key from a VIP address and prefix length.
//
// Parameters:
//   - vip: The VIP IP address.
//   - prefixLen: The prefix length.
//
// Returns the route key string.
func routeKey(vip string, prefixLen uint8) string {
	return fmt.Sprintf("%s/%d", vip, prefixLen)
}

// Advertise marks a VIP route as advertised. If the route is already advertised,
// it updates communities and local_pref but returns false for wasNew.
//
// Parameters:
//   - vip: The VIP IP address.
//   - prefixLen: The prefix length.
//   - communities: BGP communities to apply.
//   - localPref: Local preference value.
//
// Returns true if the route was newly advertised, false if already advertised.
func (s *State) Advertise(vip string, prefixLen uint8, communities []string, localPref uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := routeKey(vip, prefixLen)
	existing, exists := s.routes[key]

	if exists && existing.Advertised {
		existing.Communities = communities
		existing.LocalPref = localPref
		return false
	}

	s.routes[key] = &RouteState{
		VIP:         vip,
		PrefixLen:   prefixLen,
		Advertised:  true,
		Since:       time.Now(),
		Communities: communities,
		LocalPref:   localPref,
	}
	return true
}

// Withdraw marks a VIP route as withdrawn.
//
// Parameters:
//   - vip: The VIP IP address.
//   - prefixLen: The prefix length.
//
// Returns true if the route was previously advertised, false otherwise.
func (s *State) Withdraw(vip string, prefixLen uint8) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := routeKey(vip, prefixLen)
	existing, exists := s.routes[key]
	if !exists || !existing.Advertised {
		return false
	}

	existing.Advertised = false
	existing.Since = time.Now()
	return true
}

// GetRoute returns the state of a specific VIP route.
//
// Parameters:
//   - vip: The VIP IP address.
//
// Returns the RouteState and true if found, or nil and false if not tracked.
func (s *State) GetRoute(vip string) (*RouteState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Search through routes for matching VIP
	for _, rs := range s.routes {
		if rs.VIP == vip {
			copy := *rs
			return &copy, true
		}
	}
	return nil, false
}

// GetAllRoutes returns a copy of all tracked route states.
//
// Returns a slice of all RouteState entries.
func (s *State) GetAllRoutes() []RouteState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]RouteState, 0, len(s.routes))
	for _, rs := range s.routes {
		result = append(result, *rs)
	}
	return result
}

// IsAdvertised checks if a VIP is currently advertised.
//
// Parameters:
//   - vip: The VIP IP address.
//   - prefixLen: The prefix length.
//
// Returns true if the route exists and is currently advertised.
func (s *State) IsAdvertised(vip string, prefixLen uint8) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := routeKey(vip, prefixLen)
	rs, exists := s.routes[key]
	return exists && rs.Advertised
}
