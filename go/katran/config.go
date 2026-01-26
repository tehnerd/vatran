// Copyright (C) 2018-present, Facebook, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; version 2 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package katran

/*
#include "katran_capi.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Config contains the configuration for a Katran load balancer instance.
type Config struct {
	// MainInterface is the name of the main network interface to attach XDP program.
	MainInterface string

	// V4TunInterface is the name of the IPv4 tunnel interface for healthcheck encapsulation.
	V4TunInterface string

	// V6TunInterface is the name of the IPv6 tunnel interface for healthcheck encapsulation.
	V6TunInterface string

	// HCInterface is the interface for attaching healthcheck BPF program.
	HCInterface string

	// BalancerProgPath is the path to the compiled balancer BPF program (.o file).
	BalancerProgPath string

	// HealthcheckingProgPath is the path to the compiled healthcheck BPF program (.o file).
	HealthcheckingProgPath string

	// DefaultMAC is the MAC address of the default gateway/router (6 bytes).
	DefaultMAC []byte

	// LocalMAC is the MAC address of the local server (6 bytes).
	LocalMAC []byte

	// RootMapPath is the path to pinned BPF map from root XDP program.
	RootMapPath string

	// RootMapPos is the position (index) in the root map for Katran's program FD.
	RootMapPos uint32

	// UseRootMap indicates whether to use the root map (true) or standalone mode (false).
	UseRootMap bool

	// MaxVIPs is the maximum number of VIPs that can be configured.
	MaxVIPs uint32

	// MaxReals is the maximum number of real servers that can be configured.
	MaxReals uint32

	// CHRingSize is the size of the consistent hashing ring per VIP.
	CHRingSize uint32

	// LRUSize is the size of the per-CPU LRU connection tracking table.
	LRUSize uint64

	// MaxLPMSrcSize is the maximum number of source routing LPM entries.
	MaxLPMSrcSize uint32

	// MaxDecapDst is the maximum number of inline decapsulation destinations.
	MaxDecapDst uint32

	// GlobalLRUSize is the size of per-CPU global LRU maps.
	GlobalLRUSize uint32

	// EnableHC enables healthcheck encapsulation program.
	EnableHC bool

	// TunnelBasedHCEncap uses tunnel interfaces for healthcheck encapsulation.
	TunnelBasedHCEncap bool

	// Testing enables testing mode - don't actually program the forwarding plane.
	Testing bool

	// MemlockUnlimited sets RLIMIT_MEMLOCK to unlimited.
	MemlockUnlimited bool

	// FlowDebug enables flow debugging maps.
	FlowDebug bool

	// EnableCIDV3 enables QUIC CID version 3 support.
	EnableCIDV3 bool

	// CleanupOnShutdown cleans up BPF resources on shutdown.
	CleanupOnShutdown bool

	// ForwardingCores is the array of CPU core IDs responsible for packet forwarding.
	ForwardingCores []int32

	// NUMANodes maps forwarding cores to NUMA nodes.
	NUMANodes []int32

	// XDPAttachFlags are the XDP attach flags.
	XDPAttachFlags uint32

	// Priority is the TC priority for healthcheck program attachment.
	Priority uint32

	// MainInterfaceIndex is the interface index for main_interface.
	MainInterfaceIndex uint32

	// HCInterfaceIndex is the interface index for hc_interface.
	HCInterfaceIndex uint32

	// KatranSrcV4 is the IPv4 source address for GUE-encapsulated packets.
	KatranSrcV4 string

	// KatranSrcV6 is the IPv6 source address for GUE-encapsulated packets.
	KatranSrcV6 string

	// HashFunc is the hash function algorithm for consistent hashing.
	HashFunc HashFunction
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		RootMapPos:         2,
		UseRootMap:         true,
		MaxVIPs:            512,
		MaxReals:           4096,
		CHRingSize:         65537,
		LRUSize:            8000000,
		MaxLPMSrcSize:      3000000,
		MaxDecapDst:        6,
		GlobalLRUSize:      100000,
		EnableHC:           true,
		TunnelBasedHCEncap: true,
		MemlockUnlimited:   true,
		CleanupOnShutdown:  true,
		Priority:           2307,
		HashFunc:           HashMaglev,
	}
}

// cConfig holds C-allocated resources for a config.
type cConfig struct {
	config           C.katran_config_t
	mainInterface    *C.char
	v4TunInterface   *C.char
	v6TunInterface   *C.char
	hcInterface      *C.char
	balancerProgPath *C.char
	hcProgPath       *C.char
	rootMapPath      *C.char
	katranSrcV4      *C.char
	katranSrcV6      *C.char
	forwardingCores  *C.int32_t
	numaNodes        *C.int32_t
}

// toC converts a Go Config to a C katran_config_t.
// The caller must call free() on the returned cConfig.
func (cfg *Config) toC() *cConfig {
	cc := &cConfig{}

	// Initialize with defaults
	C.katran_config_init(&cc.config)

	// String fields
	if cfg.MainInterface != "" {
		cc.mainInterface = C.CString(cfg.MainInterface)
		cc.config.main_interface = cc.mainInterface
	}
	if cfg.V4TunInterface != "" {
		cc.v4TunInterface = C.CString(cfg.V4TunInterface)
		cc.config.v4_tun_interface = cc.v4TunInterface
	}
	if cfg.V6TunInterface != "" {
		cc.v6TunInterface = C.CString(cfg.V6TunInterface)
		cc.config.v6_tun_interface = cc.v6TunInterface
	}
	if cfg.HCInterface != "" {
		cc.hcInterface = C.CString(cfg.HCInterface)
		cc.config.hc_interface = cc.hcInterface
	}
	if cfg.BalancerProgPath != "" {
		cc.balancerProgPath = C.CString(cfg.BalancerProgPath)
		cc.config.balancer_prog_path = cc.balancerProgPath
	}
	if cfg.HealthcheckingProgPath != "" {
		cc.hcProgPath = C.CString(cfg.HealthcheckingProgPath)
		cc.config.healthchecking_prog_path = cc.hcProgPath
	}
	if cfg.RootMapPath != "" {
		cc.rootMapPath = C.CString(cfg.RootMapPath)
		cc.config.root_map_path = cc.rootMapPath
	}
	if cfg.KatranSrcV4 != "" {
		cc.katranSrcV4 = C.CString(cfg.KatranSrcV4)
		cc.config.katran_src_v4 = cc.katranSrcV4
	}
	if cfg.KatranSrcV6 != "" {
		cc.katranSrcV6 = C.CString(cfg.KatranSrcV6)
		cc.config.katran_src_v6 = cc.katranSrcV6
	}

	// MAC addresses
	if len(cfg.DefaultMAC) == 6 {
		cc.config.default_mac = (*C.uint8_t)(unsafe.Pointer(&cfg.DefaultMAC[0]))
	}
	if len(cfg.LocalMAC) == 6 {
		cc.config.local_mac = (*C.uint8_t)(unsafe.Pointer(&cfg.LocalMAC[0]))
	}

	// Numeric fields
	cc.config.root_map_pos = C.uint32_t(cfg.RootMapPos)
	cc.config.use_root_map = boolToInt(cfg.UseRootMap)
	cc.config.max_vips = C.uint32_t(cfg.MaxVIPs)
	cc.config.max_reals = C.uint32_t(cfg.MaxReals)
	cc.config.ch_ring_size = C.uint32_t(cfg.CHRingSize)
	cc.config.lru_size = C.uint64_t(cfg.LRUSize)
	cc.config.max_lpm_src_size = C.uint32_t(cfg.MaxLPMSrcSize)
	cc.config.max_decap_dst = C.uint32_t(cfg.MaxDecapDst)
	cc.config.global_lru_size = C.uint32_t(cfg.GlobalLRUSize)
	cc.config.enable_hc = boolToInt(cfg.EnableHC)
	cc.config.tunnel_based_hc_encap = boolToInt(cfg.TunnelBasedHCEncap)
	cc.config.testing = boolToInt(cfg.Testing)
	cc.config.memlock_unlimited = boolToInt(cfg.MemlockUnlimited)
	cc.config.flow_debug = boolToInt(cfg.FlowDebug)
	cc.config.enable_cid_v3 = boolToInt(cfg.EnableCIDV3)
	cc.config.cleanup_on_shutdown = boolToInt(cfg.CleanupOnShutdown)
	cc.config.xdp_attach_flags = C.uint32_t(cfg.XDPAttachFlags)
	cc.config.priority = C.uint32_t(cfg.Priority)
	cc.config.main_interface_index = C.uint32_t(cfg.MainInterfaceIndex)
	cc.config.hc_interface_index = C.uint32_t(cfg.HCInterfaceIndex)
	cc.config.hash_function = C.katran_hash_function_t(cfg.HashFunc)

	// Forwarding cores
	if len(cfg.ForwardingCores) > 0 {
		cc.forwardingCores = (*C.int32_t)(C.malloc(C.size_t(len(cfg.ForwardingCores)) * C.size_t(unsafe.Sizeof(C.int32_t(0)))))
		coreSlice := unsafe.Slice(cc.forwardingCores, len(cfg.ForwardingCores))
		for i, core := range cfg.ForwardingCores {
			coreSlice[i] = C.int32_t(core)
		}
		cc.config.forwarding_cores = cc.forwardingCores
		cc.config.forwarding_cores_count = C.size_t(len(cfg.ForwardingCores))
	}

	// NUMA nodes
	if len(cfg.NUMANodes) > 0 {
		cc.numaNodes = (*C.int32_t)(C.malloc(C.size_t(len(cfg.NUMANodes)) * C.size_t(unsafe.Sizeof(C.int32_t(0)))))
		numaSlice := unsafe.Slice(cc.numaNodes, len(cfg.NUMANodes))
		for i, node := range cfg.NUMANodes {
			numaSlice[i] = C.int32_t(node)
		}
		cc.config.numa_nodes = cc.numaNodes
		cc.config.numa_nodes_count = C.size_t(len(cfg.NUMANodes))
	}

	return cc
}

// free releases all C-allocated resources.
func (cc *cConfig) free() {
	if cc.mainInterface != nil {
		C.free(unsafe.Pointer(cc.mainInterface))
	}
	if cc.v4TunInterface != nil {
		C.free(unsafe.Pointer(cc.v4TunInterface))
	}
	if cc.v6TunInterface != nil {
		C.free(unsafe.Pointer(cc.v6TunInterface))
	}
	if cc.hcInterface != nil {
		C.free(unsafe.Pointer(cc.hcInterface))
	}
	if cc.balancerProgPath != nil {
		C.free(unsafe.Pointer(cc.balancerProgPath))
	}
	if cc.hcProgPath != nil {
		C.free(unsafe.Pointer(cc.hcProgPath))
	}
	if cc.rootMapPath != nil {
		C.free(unsafe.Pointer(cc.rootMapPath))
	}
	if cc.katranSrcV4 != nil {
		C.free(unsafe.Pointer(cc.katranSrcV4))
	}
	if cc.katranSrcV6 != nil {
		C.free(unsafe.Pointer(cc.katranSrcV6))
	}
	if cc.forwardingCores != nil {
		C.free(unsafe.Pointer(cc.forwardingCores))
	}
	if cc.numaNodes != nil {
		C.free(unsafe.Pointer(cc.numaNodes))
	}
}

// boolToInt converts a Go bool to a C int.
func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
