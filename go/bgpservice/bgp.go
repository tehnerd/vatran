package bgpservice

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	gobgpapi "github.com/osrg/gobgp/v3/api"
	gobgpserver "github.com/osrg/gobgp/v3/pkg/server"
	"google.golang.org/protobuf/types/known/anypb"
)

// PeerStatus represents the current status of a BGP peer session.
type PeerStatus struct {
	// Address is the peer IP address.
	Address string `json:"address"`
	// ASN is the peer autonomous system number.
	ASN uint32 `json:"asn"`
	// State is the BGP session state (e.g., "established", "active", "idle").
	State string `json:"state"`
	// Uptime is the duration the session has been established.
	Uptime string `json:"uptime"`
	// PrefixesAnnounced is the number of prefixes announced to this peer.
	PrefixesAnnounced int `json:"prefixes_announced"`
}

// BGPSpeaker wraps the GoBGP server to provide route announcement and peer management.
type BGPSpeaker struct {
	server *gobgpserver.BgpServer
	config *BGPConfig
}

// NewBGPSpeaker creates a new BGPSpeaker with the given configuration.
//
// Parameters:
//   - config: BGP configuration including ASN, router ID, and initial peers.
//
// Returns a new BGPSpeaker instance.
func NewBGPSpeaker(config *BGPConfig) *BGPSpeaker {
	return &BGPSpeaker{
		config: config,
	}
}

// Start initializes and starts the GoBGP server, sets the global configuration,
// and adds any initial peers from the configuration.
//
// Returns an error if the BGP server fails to start.
func (b *BGPSpeaker) Start() error {
	b.server = gobgpserver.NewBgpServer()
	go b.server.Serve()

	// Set global BGP config
	listenPort := int32(b.config.ListenPort)
	if listenPort == 0 {
		listenPort = 179
	}

	err := b.server.StartBgp(context.Background(), &gobgpapi.StartBgpRequest{
		Global: &gobgpapi.Global{
			Asn:        b.config.ASN,
			RouterId:   b.config.RouterID,
			ListenPort: listenPort,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start BGP server: %w", err)
	}

	log.Printf("BGP speaker started: ASN=%d RouterID=%s ListenPort=%d",
		b.config.ASN, b.config.RouterID, listenPort)

	// Add initial peers
	for _, peer := range b.config.Peers {
		if err := b.AddPeer(peer); err != nil {
			log.Printf("Warning: failed to add initial peer %s: %v", peer.Address, err)
		}
	}

	return nil
}

// Stop gracefully stops the GoBGP server.
func (b *BGPSpeaker) Stop() {
	if b.server != nil {
		b.server.Stop()
		log.Println("BGP speaker stopped")
	}
}

// AnnounceRoute advertises a VIP prefix via BGP to all peers.
//
// Parameters:
//   - vip: The VIP IP address to announce.
//   - prefixLen: The prefix length.
//   - communities: BGP communities to attach (format: "asn:value").
//   - localPref: Local preference value.
//
// Returns an error if the announcement fails.
func (b *BGPSpeaker) AnnounceRoute(vip string, prefixLen uint8, communities []string, localPref uint32) error {
	path, err := b.buildPath(vip, prefixLen, communities, localPref)
	if err != nil {
		return fmt.Errorf("failed to build BGP path: %w", err)
	}

	_, err = b.server.AddPath(context.Background(), &gobgpapi.AddPathRequest{
		Path: path,
	})
	if err != nil {
		return fmt.Errorf("failed to announce route %s/%d: %w", vip, prefixLen, err)
	}

	log.Printf("BGP: announced route %s/%d", vip, prefixLen)
	return nil
}

// WithdrawRoute withdraws a VIP prefix from BGP announcements.
//
// Parameters:
//   - vip: The VIP IP address to withdraw.
//   - prefixLen: The prefix length.
//
// Returns an error if the withdrawal fails.
func (b *BGPSpeaker) WithdrawRoute(vip string, prefixLen uint8) error {
	path, err := b.buildPath(vip, prefixLen, nil, 0)
	if err != nil {
		return fmt.Errorf("failed to build BGP path: %w", err)
	}

	err = b.server.DeletePath(context.Background(), &gobgpapi.DeletePathRequest{
		Path: path,
	})
	if err != nil {
		return fmt.Errorf("failed to withdraw route %s/%d: %w", vip, prefixLen, err)
	}

	log.Printf("BGP: withdrew route %s/%d", vip, prefixLen)
	return nil
}

// AddPeer adds a new BGP peer to the speaker.
//
// Parameters:
//   - cfg: The peer configuration.
//
// Returns an error if adding the peer fails.
func (b *BGPSpeaker) AddPeer(cfg PeerConfig) error {
	holdTime := uint64(cfg.HoldTime)
	if holdTime == 0 {
		holdTime = 90
	}
	keepalive := uint64(cfg.Keepalive)
	if keepalive == 0 {
		keepalive = 30
	}

	// Determine address family based on peer address
	afiSafis := []*gobgpapi.AfiSafi{
		{
			Config: &gobgpapi.AfiSafiConfig{
				Family: &gobgpapi.Family{
					Afi:  gobgpapi.Family_AFI_IP,
					Safi: gobgpapi.Family_SAFI_UNICAST,
				},
				Enabled: true,
			},
		},
	}
	if net.ParseIP(cfg.Address) != nil && net.ParseIP(cfg.Address).To4() == nil {
		afiSafis = append(afiSafis, &gobgpapi.AfiSafi{
			Config: &gobgpapi.AfiSafiConfig{
				Family: &gobgpapi.Family{
					Afi:  gobgpapi.Family_AFI_IP6,
					Safi: gobgpapi.Family_SAFI_UNICAST,
				},
				Enabled: true,
			},
		})
	}

	err := b.server.AddPeer(context.Background(), &gobgpapi.AddPeerRequest{
		Peer: &gobgpapi.Peer{
			Conf: &gobgpapi.PeerConf{
				NeighborAddress: cfg.Address,
				PeerAsn:         cfg.ASN,
			},
			Timers: &gobgpapi.Timers{
				Config: &gobgpapi.TimersConfig{
					HoldTime:          holdTime,
					KeepaliveInterval: keepalive,
				},
			},
			AfiSafis: afiSafis,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add peer %s (ASN %d): %w", cfg.Address, cfg.ASN, err)
	}

	log.Printf("BGP: added peer %s ASN=%d", cfg.Address, cfg.ASN)
	return nil
}

// RemovePeer removes a BGP peer from the speaker.
//
// Parameters:
//   - address: The peer IP address to remove.
//
// Returns an error if removing the peer fails.
func (b *BGPSpeaker) RemovePeer(address string) error {
	err := b.server.DeletePeer(context.Background(), &gobgpapi.DeletePeerRequest{
		Address: address,
	})
	if err != nil {
		return fmt.Errorf("failed to remove peer %s: %w", address, err)
	}

	log.Printf("BGP: removed peer %s", address)
	return nil
}

// ListPeers returns the current status of all configured BGP peers.
//
// Returns a slice of PeerStatus or an error.
func (b *BGPSpeaker) ListPeers() ([]PeerStatus, error) {
	var peers []PeerStatus

	err := b.server.ListPeer(context.Background(), &gobgpapi.ListPeerRequest{}, func(p *gobgpapi.Peer) {
		status := PeerStatus{
			Address: p.Conf.NeighborAddress,
			ASN:     p.Conf.PeerAsn,
		}

		if p.State != nil {
			status.State = strings.ToLower(p.State.SessionState.String())

			if p.Timers != nil && p.Timers.State != nil && p.Timers.State.Uptime != nil {
				uptime := p.Timers.State.Uptime.AsTime()
				if !uptime.IsZero() {
					status.Uptime = time.Since(uptime).Truncate(time.Second).String()
				}
			}

			// Count announced prefixes from AfiSafi state
			for _, afi := range p.AfiSafis {
				if afi.State != nil {
					status.PrefixesAnnounced += int(afi.State.Advertised)
				}
			}
		}

		peers = append(peers, status)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list peers: %w", err)
	}

	if peers == nil {
		peers = []PeerStatus{}
	}
	return peers, nil
}

// buildPath constructs a GoBGP Path for a given VIP prefix with attributes.
func (b *BGPSpeaker) buildPath(vip string, prefixLen uint8, communities []string, localPref uint32) (*gobgpapi.Path, error) {
	// Build NLRI
	nlri, err := anypb.New(&gobgpapi.IPAddressPrefix{
		PrefixLen: uint32(prefixLen),
		Prefix:    vip,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create NLRI: %w", err)
	}

	// Determine address family
	family := &gobgpapi.Family{
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_UNICAST,
	}
	ip := net.ParseIP(vip)
	if ip != nil && ip.To4() == nil {
		family.Afi = gobgpapi.Family_AFI_IP6
	}

	// Build path attributes
	var pattrs []*anypb.Any

	// Origin: IGP
	origin, err := anypb.New(&gobgpapi.OriginAttribute{
		Origin: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create origin attribute: %w", err)
	}
	pattrs = append(pattrs, origin)

	// Next hop
	nextHop := "0.0.0.0"
	if family.Afi == gobgpapi.Family_AFI_IP6 {
		nextHop = "::"
	}
	nh, err := anypb.New(&gobgpapi.NextHopAttribute{
		NextHop: nextHop,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create next-hop attribute: %w", err)
	}
	pattrs = append(pattrs, nh)

	// Local preference
	if localPref > 0 {
		lp, err := anypb.New(&gobgpapi.LocalPrefAttribute{
			LocalPref: localPref,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create local-pref attribute: %w", err)
		}
		pattrs = append(pattrs, lp)
	}

	// Communities
	if len(communities) > 0 {
		communityValues, err := parseCommunities(communities)
		if err != nil {
			return nil, err
		}
		comm, err := anypb.New(&gobgpapi.CommunitiesAttribute{
			Communities: communityValues,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create communities attribute: %w", err)
		}
		pattrs = append(pattrs, comm)
	}

	return &gobgpapi.Path{
		Family: family,
		Nlri:   nlri,
		Pattrs: pattrs,
	}, nil
}

// parseCommunities converts community strings (e.g., "65000:100") to uint32 values.
func parseCommunities(communities []string) ([]uint32, error) {
	result := make([]uint32, 0, len(communities))
	for _, c := range communities {
		parts := strings.SplitN(c, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid community format %q (expected ASN:VALUE)", c)
		}
		high, err := strconv.ParseUint(parts[0], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid community ASN in %q: %w", c, err)
		}
		low, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid community value in %q: %w", c, err)
		}
		result = append(result, uint32(high)<<16|uint32(low))
	}
	return result, nil
}
