package lb

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tehnerd/vatran/go/katran"
	"github.com/tehnerd/vatran/go/server/types"
)

// defaultReconcileEvery is the number of poll cycles between reconciliation runs.
// With a default 5s poll interval, this gives ~30s between reconciliations.
const defaultReconcileEvery = 6

// HCPoller periodically polls the healthcheck service for health states
// and applies transitions to the local state store and katran.
type HCPoller struct {
	manager        *Manager
	interval       time.Duration
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	reconcileEvery int
	pollCount      int
}

// NewHCPoller creates a new HCPoller.
//
// Parameters:
//   - manager: The Manager singleton.
//   - interval: The polling interval.
//
// Returns a new HCPoller instance.
func NewHCPoller(manager *Manager, interval time.Duration) *HCPoller {
	return &HCPoller{
		manager:        manager,
		interval:       interval,
		reconcileEvery: defaultReconcileEvery,
	}
}

// Start begins the background polling loop.
func (p *HCPoller) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg.Add(1)
	go p.run(ctx)
	log.Printf("HC poller started with interval %s", p.interval)
}

// Stop stops the background polling loop and waits for it to finish.
func (p *HCPoller) Stop() {
	if p.cancel != nil {
		p.cancel()
		p.wg.Wait()
		log.Println("HC poller stopped")
	}
}

// run is the main polling loop.
func (p *HCPoller) run(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

// poll fetches health states from the HC service and applies transitions.
func (p *HCPoller) poll(ctx context.Context) {
	hcClient := p.manager.GetHCClient()
	if hcClient == nil {
		return
	}

	state, stateOK := p.manager.GetState()
	if !stateOK {
		return
	}

	lbInstance, lbOK := p.manager.Get()
	if !lbOK {
		return
	}

	// Run reconciliation on a slower cadence
	p.pollCount++
	if p.pollCount >= p.reconcileEvery {
		p.pollCount = 0
		p.reconcile(ctx, state, hcClient)
	}

	allHealth, err := hcClient.GetAllHealth(ctx)
	if err != nil {
		log.Printf("HC poller: failed to fetch health states: %v", err)
		return
	}

	// Track which VIPs had health transitions for BGP evaluation
	vipsWithTransitions := make(map[string]bool)

	for _, vipHealth := range allHealth {
		vipKey := VIPKeyString(vipHealth.VIP.Address, vipHealth.VIP.Port, vipHealth.VIP.Proto)

		// Skip VIPs with dummy HC config
		hcCfg, hasHC := state.GetHCConfig(vipKey)
		if !hasHC || hcCfg.Type == "dummy" {
			continue
		}

		vip := katran.VIPKey{
			Address: vipHealth.VIP.Address,
			Port:    vipHealth.VIP.Port,
			Proto:   vipHealth.VIP.Proto,
		}

		for _, realHealth := range vipHealth.Reals {
			oldHealthy, found := state.UpdateHealth(vipKey, realHealth.Address, realHealth.Healthy)
			if !found {
				continue
			}
			if oldHealthy == realHealth.Healthy {
				continue
			}

			vipsWithTransitions[vipKey] = true

			// Get real info for katran operation
			reals := state.GetReals(vipKey)
			var real katran.Real
			for _, rs := range reals {
				if rs.Address == realHealth.Address {
					real = katran.Real{
						Address: rs.Address,
						Weight:  rs.Weight,
						Flags:   rs.Flags,
					}
					break
				}
			}

			if realHealth.Healthy {
				// unhealthy -> healthy: add to katran
				if err := lbInstance.AddRealForVIP(real, vip); err != nil {
					log.Printf("HC poller: failed to add real %s to VIP %s: %v", realHealth.Address, vipKey, err)
					state.UpdateHealth(vipKey, realHealth.Address, oldHealthy)
				}
			} else {
				// healthy -> unhealthy: remove from katran
				if err := lbInstance.DelRealForVIP(real, vip); err != nil {
					log.Printf("HC poller: failed to remove real %s from VIP %s: %v", realHealth.Address, vipKey, err)
					state.UpdateHealth(vipKey, realHealth.Address, oldHealthy)
				}
			}
		}
	}

	// Evaluate BGP advertise/withdraw for VIPs that had health transitions
	p.evaluateBGP(ctx, state, vipsWithTransitions)
}

// reconcile checks that the HC service knows about all locally registered VIPs
// and re-registers any that are missing. This handles the case where the HC
// service restarts or starts after the katran server, losing its registrations.
//
// Parameters:
//   - ctx: Context for the reconciliation requests.
//   - state: The local VIP/reals state store.
//   - hcClient: The healthcheck service client.
func (p *HCPoller) reconcile(ctx context.Context, state *VIPRealsState, hcClient *HCClient) {
	localConfigs := state.GetAllHCConfigs()

	// Filter to non-dummy configs only
	nonDummy := make(map[string]*types.HealthcheckConfig)
	for vipKey, cfg := range localConfigs {
		if cfg.Type != "dummy" {
			nonDummy[vipKey] = cfg
		}
	}
	if len(nonDummy) == 0 {
		return
	}

	allHealth, err := hcClient.GetAllHealth(ctx)
	if err != nil {
		log.Printf("HC poller reconcile: failed to fetch health states: %v", err)
		return
	}

	// Build a map of VIP keys known to the HC service, along with their real counts
	type hcVIPInfo struct {
		realCount int
	}
	knownVIPs := make(map[string]hcVIPInfo, len(allHealth))
	for _, vipHealth := range allHealth {
		key := VIPKeyString(vipHealth.VIP.Address, vipHealth.VIP.Port, vipHealth.VIP.Proto)
		knownVIPs[key] = hcVIPInfo{realCount: len(vipHealth.Reals)}
	}

	for vipKey, hcCfg := range nonDummy {
		address, port, proto, err := ParseVIPKey(vipKey)
		if err != nil {
			log.Printf("HC poller reconcile: failed to parse VIP key %q: %v", vipKey, err)
			continue
		}
		hcVIP := types.HCVIPKey{Address: address, Port: port, Proto: proto}
		reals := state.GetReals(vipKey)

		info, registered := knownVIPs[vipKey]
		if !registered {
			// VIP is missing from the HC service — re-register it
			if err := hcClient.RegisterVIP(ctx, hcVIP, reals, hcCfg); err != nil {
				log.Printf("HC poller reconcile: failed to re-register VIP %s: %v", vipKey, err)
			} else {
				log.Printf("HC poller reconcile: re-registered VIP %s with %d reals", vipKey, len(reals))
			}
		} else if len(reals) != info.realCount {
			// VIP is registered but real count differs — update
			if err := hcClient.UpdateVIP(ctx, hcVIP, hcCfg, reals); err != nil {
				log.Printf("HC poller reconcile: failed to update reals for VIP %s: %v", vipKey, err)
			} else {
				log.Printf("HC poller reconcile: updated VIP %s reals (local=%d, remote=%d)", vipKey, len(reals), info.realCount)
			}
		}
	}
}

// evaluateBGP checks VIPs with health transitions and advertises or withdraws
// routes via the BGP service based on the healthy real count threshold.
func (p *HCPoller) evaluateBGP(ctx context.Context, state *VIPRealsState, vipKeys map[string]bool) {
	bgpClient := p.manager.GetBGPClient()
	if bgpClient == nil {
		return
	}

	threshold := p.manager.GetBGPMinHealthyReals()

	for vipKey := range vipKeys {
		healthyCount := state.CountHealthyReals(vipKey)
		vipAddress := GetVIPAddress(vipKey)

		if healthyCount >= threshold {
			if err := bgpClient.Advertise(ctx, vipAddress, 32); err != nil {
				log.Printf("HC poller: failed to advertise VIP %s via BGP: %v", vipAddress, err)
			}
		} else {
			if err := bgpClient.Withdraw(ctx, vipAddress, 32); err != nil {
				log.Printf("HC poller: failed to withdraw VIP %s from BGP: %v", vipAddress, err)
			}
		}
	}
}
