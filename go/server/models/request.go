package models

// CreateLBRequest represents a request to create a load balancer.
type CreateLBRequest struct {
	// MainInterface is the name of the main network interface to attach XDP program.
	MainInterface string `json:"main_interface"`
	// V4TunInterface is the name of the IPv4 tunnel interface for healthcheck encapsulation.
	V4TunInterface string `json:"v4_tun_interface,omitempty"`
	// V6TunInterface is the name of the IPv6 tunnel interface for healthcheck encapsulation.
	V6TunInterface string `json:"v6_tun_interface,omitempty"`
	// HCInterface is the interface for attaching healthcheck BPF program.
	HCInterface string `json:"hc_interface,omitempty"`
	// BalancerProgPath is the path to the compiled balancer BPF program (.o file).
	BalancerProgPath string `json:"balancer_prog_path"`
	// HealthcheckingProgPath is the path to the compiled healthcheck BPF program (.o file).
	HealthcheckingProgPath string `json:"healthchecking_prog_path,omitempty"`
	// DefaultMAC is the MAC address of the default gateway/router (hex string, e.g., "aa:bb:cc:dd:ee:ff").
	DefaultMAC string `json:"default_mac,omitempty"`
	// LocalMAC is the MAC address of the local server (hex string).
	LocalMAC string `json:"local_mac,omitempty"`
	// RootMapPath is the path to pinned BPF map from root XDP program.
	RootMapPath string `json:"root_map_path,omitempty"`
	// RootMapPos is the position (index) in the root map for Katran's program FD.
	RootMapPos uint32 `json:"root_map_pos,omitempty"`
	// UseRootMap indicates whether to use the root map (true) or standalone mode (false).
	UseRootMap *bool `json:"use_root_map,omitempty"`
	// MaxVIPs is the maximum number of VIPs that can be configured.
	MaxVIPs uint32 `json:"max_vips,omitempty"`
	// MaxReals is the maximum number of real servers that can be configured.
	MaxReals uint32 `json:"max_reals,omitempty"`
	// CHRingSize is the size of the consistent hashing ring per VIP.
	CHRingSize uint32 `json:"ch_ring_size,omitempty"`
	// LRUSize is the size of the per-CPU LRU connection tracking table.
	LRUSize uint64 `json:"lru_size,omitempty"`
	// MaxLPMSrcSize is the maximum number of source routing LPM entries.
	MaxLPMSrcSize uint32 `json:"max_lpm_src_size,omitempty"`
	// MaxDecapDst is the maximum number of inline decapsulation destinations.
	MaxDecapDst uint32 `json:"max_decap_dst,omitempty"`
	// GlobalLRUSize is the size of per-CPU global LRU maps.
	GlobalLRUSize uint32 `json:"global_lru_size,omitempty"`
	// EnableHC enables healthcheck encapsulation program.
	EnableHC *bool `json:"enable_hc,omitempty"`
	// TunnelBasedHCEncap uses tunnel interfaces for healthcheck encapsulation.
	TunnelBasedHCEncap *bool `json:"tunnel_based_hc_encap,omitempty"`
	// Testing enables testing mode.
	Testing bool `json:"testing,omitempty"`
	// MemlockUnlimited sets RLIMIT_MEMLOCK to unlimited.
	MemlockUnlimited *bool `json:"memlock_unlimited,omitempty"`
	// FlowDebug enables flow debugging maps.
	FlowDebug bool `json:"flow_debug,omitempty"`
	// EnableCIDV3 enables QUIC CID version 3 support.
	EnableCIDV3 bool `json:"enable_cid_v3,omitempty"`
	// CleanupOnShutdown cleans up BPF resources on shutdown.
	CleanupOnShutdown *bool `json:"cleanup_on_shutdown,omitempty"`
	// ForwardingCores is the array of CPU core IDs responsible for packet forwarding.
	ForwardingCores []int32 `json:"forwarding_cores,omitempty"`
	// NUMANodes maps forwarding cores to NUMA nodes.
	NUMANodes []int32 `json:"numa_nodes,omitempty"`
	// XDPAttachFlags are the XDP attach flags.
	XDPAttachFlags uint32 `json:"xdp_attach_flags,omitempty"`
	// Priority is the TC priority for healthcheck program attachment.
	Priority uint32 `json:"priority,omitempty"`
	// MainInterfaceIndex is the interface index for main_interface.
	MainInterfaceIndex uint32 `json:"main_interface_index,omitempty"`
	// HCInterfaceIndex is the interface index for hc_interface.
	HCInterfaceIndex uint32 `json:"hc_interface_index,omitempty"`
	// KatranSrcV4 is the IPv4 source address for GUE-encapsulated packets.
	KatranSrcV4 string `json:"katran_src_v4,omitempty"`
	// KatranSrcV6 is the IPv6 source address for GUE-encapsulated packets.
	KatranSrcV6 string `json:"katran_src_v6,omitempty"`
	// HashFunction is the hash function algorithm for consistent hashing ("maglev" or "maglev_v2").
	// Defaults to "maglev_v2" when omitted.
	HashFunction string `json:"hash_function,omitempty"`
}

// VIPRequest represents a VIP key in API requests.
type VIPRequest struct {
	// Address is the IP address of the VIP (IPv4 or IPv6).
	Address string `json:"address"`
	// Port is the port number in host byte order.
	Port uint16 `json:"port"`
	// Proto is the IP protocol number (e.g., 6 for TCP, 17 for UDP).
	Proto uint8 `json:"proto"`
}

// AddVIPRequest represents a request to add a VIP.
type AddVIPRequest struct {
	VIPRequest
	// Flags are the VIP flags.
	Flags uint32 `json:"flags,omitempty"`
}

// ModifyVIPFlagsRequest represents a request to modify VIP flags.
type ModifyVIPFlagsRequest struct {
	VIPRequest
	// Flag is the flag bits to modify.
	Flag uint32 `json:"flag"`
	// Set is true to set the flags, false to clear them.
	Set bool `json:"set"`
}

// ChangeHashFunctionRequest represents a request to change the hash function for a VIP.
type ChangeHashFunctionRequest struct {
	VIPRequest
	// HashFunction is the hash function to use (0=maglev, 1=maglevV2).
	HashFunction int `json:"hash_function"`
}

// RealRequest represents a real server in API requests.
type RealRequest struct {
	// Address is the IP address of the real server (IPv4 or IPv6).
	Address string `json:"address"`
	// Weight is the weight for consistent hashing (higher = more traffic).
	Weight uint32 `json:"weight"`
	// Flags contains real-specific flags.
	Flags uint8 `json:"flags,omitempty"`
}

// AddRealRequest represents a request to add a real server to a VIP.
type AddRealRequest struct {
	// VIP is the target VIP.
	VIP VIPRequest `json:"vip"`
	// Real is the real server to add.
	Real RealRequest `json:"real"`
}

// DelRealRequest represents a request to delete a real server from a VIP.
type DelRealRequest struct {
	// VIP is the target VIP.
	VIP VIPRequest `json:"vip"`
	// Real is the real server to delete.
	Real RealRequest `json:"real"`
}

// ModifyRealsRequest represents a request to batch modify reals for a VIP.
type ModifyRealsRequest struct {
	// VIP is the target VIP.
	VIP VIPRequest `json:"vip"`
	// Action is the action to perform (0=add, 1=delete).
	Action int `json:"action"`
	// Reals is the list of real servers to modify.
	Reals []RealRequest `json:"reals"`
}

// ModifyRealFlagsRequest represents a request to modify a real server's flags.
type ModifyRealFlagsRequest struct {
	// Address is the real server IP address.
	Address string `json:"address"`
	// Flags is the flag bits to modify.
	Flags uint8 `json:"flags"`
	// Set is true to set the flags, false to clear them.
	Set bool `json:"set"`
}

// GetRealIndexRequest represents a request to get a real server's index.
type GetRealIndexRequest struct {
	// Address is the real server IP address.
	Address string `json:"address"`
}

// GetRealStatsRequest represents a request to get real server statistics.
type GetRealStatsRequest struct {
	// Index is the real server index.
	Index uint32 `json:"index"`
}

// GetVIPStatsRequest represents a request to get VIP statistics.
type GetVIPStatsRequest struct {
	VIPRequest
}

// GetBPFMapStatsRequest represents a request to get BPF map statistics.
type GetBPFMapStatsRequest struct {
	// MapName is the name of the BPF map.
	MapName string `json:"map_name"`
}

// QuicRealRequest represents a QUIC real server mapping in API requests.
type QuicRealRequest struct {
	// Address is the IP address of the real server.
	Address string `json:"address"`
	// ID is the QUIC host ID portion embedded in connection IDs.
	ID uint32 `json:"id"`
}

// ModifyQuicRealsRequest represents a request to modify QUIC real mappings.
type ModifyQuicRealsRequest struct {
	// Action is the action to perform (0=add, 1=delete).
	Action int `json:"action"`
	// Reals is the list of QUIC real mappings to modify.
	Reals []QuicRealRequest `json:"reals"`
}

// AddSrcRoutingRuleRequest represents a request to add source routing rules.
type AddSrcRoutingRuleRequest struct {
	// SrcPrefixes is an array of source IP prefixes (CIDR notation).
	SrcPrefixes []string `json:"src_prefixes"`
	// Dst is the destination address for matching traffic.
	Dst string `json:"dst"`
}

// DelSrcRoutingRuleRequest represents a request to delete source routing rules.
type DelSrcRoutingRuleRequest struct {
	// SrcPrefixes is an array of source IP prefixes to remove.
	SrcPrefixes []string `json:"src_prefixes"`
}

// InlineDecapDstRequest represents a request to add/delete inline decap destination.
type InlineDecapDstRequest struct {
	// Dst is the destination IP address.
	Dst string `json:"dst"`
}

// HealthcheckerDstRequest represents a request to add a healthcheck destination.
type HealthcheckerDstRequest struct {
	// Somark is the socket mark value to match.
	Somark uint32 `json:"somark"`
	// Dst is the destination address for encapsulated packets.
	Dst string `json:"dst"`
}

// DelHealthcheckerDstRequest represents a request to delete a healthcheck destination.
type DelHealthcheckerDstRequest struct {
	// Somark is the socket mark value to remove.
	Somark uint32 `json:"somark"`
}

// HCKeyRequest represents a healthcheck key request.
type HCKeyRequest struct {
	VIPRequest
}

// FeatureRequest represents a feature request.
type FeatureRequest struct {
	// Feature is the feature flag.
	Feature int `json:"feature"`
	// ProgPath is the path to BPF program, empty to use current.
	ProgPath string `json:"prog_path,omitempty"`
}

// DeleteLRURequest represents a request to delete an LRU entry.
type DeleteLRURequest struct {
	// DstVIP is the destination VIP.
	DstVIP VIPRequest `json:"dst_vip"`
	// SrcIP is the source IP address.
	SrcIP string `json:"src_ip"`
	// SrcPort is the source port.
	SrcPort uint16 `json:"src_port"`
}

// PurgeVIPLRURequest represents a request to purge all LRU entries for a VIP.
type PurgeVIPLRURequest struct {
	VIPRequest
}

// RestartMonitorRequest represents a request to restart the monitor.
type RestartMonitorRequest struct {
	// Limit is the maximum number of packets to capture (0 = unlimited).
	Limit uint32 `json:"limit"`
}

// ChangeMACRequest represents a request to change the MAC address.
type ChangeMACRequest struct {
	// MAC is the 6-byte MAC address as hex string (e.g., "aa:bb:cc:dd:ee:ff").
	MAC string `json:"mac"`
}

// FlowRequest represents a flow for routing simulation.
type FlowRequest struct {
	// Src is the source IP address.
	Src string `json:"src"`
	// Dst is the destination IP address.
	Dst string `json:"dst"`
	// SrcPort is the source port number.
	SrcPort uint16 `json:"src_port"`
	// DstPort is the destination port number.
	DstPort uint16 `json:"dst_port"`
	// Proto is the IP protocol number.
	Proto uint8 `json:"proto"`
}

// SimulatePacketRequest represents a request to simulate a packet.
type SimulatePacketRequest struct {
	// Packet is the raw packet data as base64-encoded string.
	Packet string `json:"packet"`
}

// ReloadBalancerProgRequest represents a request to reload the balancer program.
type ReloadBalancerProgRequest struct {
	// Path is the path to the new BPF program file.
	Path string `json:"path"`
	// Config is optional new configuration.
	Config *CreateLBRequest `json:"config,omitempty"`
}

// AddSrcIPForPcktEncapRequest represents a request to add source IP for packet encapsulation.
type AddSrcIPForPcktEncapRequest struct {
	// Src is the source IP address (IPv4 or IPv6).
	Src string `json:"src"`
}
