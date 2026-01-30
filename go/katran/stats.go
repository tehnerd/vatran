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

// GetStatsForVIP retrieves packet/byte statistics for a VIP.
//
// Parameters:
//   - vip: VIP to query.
//
// Returns stats where V1 = packets and V2 = bytes, or ErrNotFound if the VIP
// does not exist.
func (lb *LoadBalancer) GetStatsForVIP(vip VIPKey) (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_stats_for_vip(lb.handle, &cVip, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetDecapStatsForVIP retrieves decapsulation statistics for a VIP.
//
// Parameters:
//   - vip: VIP to query.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) GetDecapStatsForVIP(vip VIPKey) (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_decap_stats_for_vip(lb.handle, &cVip, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetLRUStats retrieves LRU cache statistics.
//
// Returns stats where V1 = total packets and V2 = LRU hits.
func (lb *LoadBalancer) GetLRUStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_lru_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetLRUMissStats retrieves LRU miss statistics.
//
// Returns stats where V1 = TCP SYN misses and V2 = non-SYN misses.
func (lb *LoadBalancer) GetLRUMissStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_lru_miss_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetLRUFallbackStats retrieves LRU fallback statistics.
func (lb *LoadBalancer) GetLRUFallbackStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_lru_fallback_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetICMPTooBigStats retrieves ICMP "too big" statistics.
//
// Returns stats where V1 = ICMPv4 count and V2 = ICMPv6 count.
func (lb *LoadBalancer) GetICMPTooBigStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_icmp_too_big_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetCHDropStats retrieves consistent hash drop statistics.
//
// Returns stats where V1 = real ID out of bounds and V2 = real #0 (unmapped).
func (lb *LoadBalancer) GetCHDropStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_ch_drop_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetSrcRoutingStats retrieves source routing statistics.
//
// Returns stats where V1 = local backend and V2 = remote (LPM matched).
func (lb *LoadBalancer) GetSrcRoutingStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_src_routing_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetInlineDecapStats retrieves inline decapsulation statistics.
func (lb *LoadBalancer) GetInlineDecapStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_inline_decap_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetGlobalLRUStats retrieves global LRU statistics.
//
// Returns stats where V1 = map lookup failures and V2 = global LRU routed.
func (lb *LoadBalancer) GetGlobalLRUStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_global_lru_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetDecapStats retrieves general decapsulation statistics.
//
// Returns stats where V1 = IPv4 decapped and V2 = IPv6 decapped.
func (lb *LoadBalancer) GetDecapStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_decap_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetQuicICMPStats retrieves QUIC ICMP statistics.
func (lb *LoadBalancer) GetQuicICMPStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_quic_icmp_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetRealStats retrieves per-real server statistics.
//
// Parameters:
//   - index: Real server index (from GetIndexForReal()).
//
// Returns stats where V1 = packets and V2 = bytes.
func (lb *LoadBalancer) GetRealStats(index uint32) (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_real_stats(lb.handle, C.uint32_t(index), &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetXDPTotalStats retrieves total XDP statistics.
//
// Returns stats where V1 = packets and V2 = bytes.
func (lb *LoadBalancer) GetXDPTotalStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_xdp_total_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetXDPTXStats retrieves XDP TX statistics.
//
// Returns packets/bytes forwarded (XDP_TX).
func (lb *LoadBalancer) GetXDPTXStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_xdp_tx_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetXDPDropStats retrieves XDP drop statistics.
//
// Returns packets/bytes dropped (XDP_DROP).
func (lb *LoadBalancer) GetXDPDropStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_xdp_drop_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetXDPPassStats retrieves XDP pass statistics.
//
// Returns packets/bytes passed to kernel (XDP_PASS).
func (lb *LoadBalancer) GetXDPPassStats() (LBStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return LBStats{}, err
	}

	var stats C.katran_lb_stats_t
	ret := C.katran_lb_get_xdp_pass_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return LBStats{}, errorFromCode(ret)
	}

	return LBStats{V1: uint64(stats.v1), V2: uint64(stats.v2)}, nil
}

// GetTCPServerIDRoutingStats retrieves TCP server ID routing statistics (TPR).
func (lb *LoadBalancer) GetTCPServerIDRoutingStats() (TPRPacketsStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return TPRPacketsStats{}, err
	}

	var stats C.katran_tpr_packets_stats_t
	ret := C.katran_lb_get_tcp_server_id_routing_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return TPRPacketsStats{}, errorFromCode(ret)
	}

	return TPRPacketsStats{
		CHRouted:         uint64(stats.ch_routed),
		DstMismatchInLRU: uint64(stats.dst_mismatch_in_lru),
		SIDRouted:        uint64(stats.sid_routed),
		TCPSyn:           uint64(stats.tcp_syn),
	}, nil
}

// GetQuicPacketsStats retrieves QUIC packet routing statistics.
func (lb *LoadBalancer) GetQuicPacketsStats() (QuicPacketsStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return QuicPacketsStats{}, err
	}

	var stats C.katran_quic_packets_stats_t
	ret := C.katran_lb_get_quic_packets_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return QuicPacketsStats{}, errorFromCode(ret)
	}

	return QuicPacketsStats{
		CHRouted:                 uint64(stats.ch_routed),
		CIDInitial:               uint64(stats.cid_initial),
		CIDInvalidServerID:       uint64(stats.cid_invalid_server_id),
		CIDInvalidServerIDSample: uint64(stats.cid_invalid_server_id_sample),
		CIDRouted:                uint64(stats.cid_routed),
		CIDUnknownRealDropped:    uint64(stats.cid_unknown_real_dropped),
		CIDV0:                    uint64(stats.cid_v0),
		CIDV1:                    uint64(stats.cid_v1),
		CIDV2:                    uint64(stats.cid_v2),
		CIDV3:                    uint64(stats.cid_v3),
		DstMatchInLRU:            uint64(stats.dst_match_in_lru),
		DstMismatchInLRU:         uint64(stats.dst_mismatch_in_lru),
		DstNotFoundInLRU:         uint64(stats.dst_not_found_in_lru),
	}, nil
}

// GetHCProgStats retrieves healthcheck program statistics.
//
// Returns ErrFeatureDisabled if healthchecking is not enabled.
func (lb *LoadBalancer) GetHCProgStats() (HCStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return HCStats{}, err
	}

	var stats C.katran_hc_stats_t
	ret := C.katran_lb_get_hc_prog_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return HCStats{}, errorFromCode(ret)
	}

	return HCStats{
		PacketsProcessed: uint64(stats.packets_processed),
		PacketsDropped:   uint64(stats.packets_dropped),
		PacketsSkipped:   uint64(stats.packets_skipped),
		PacketsTooBig:    uint64(stats.packets_too_big),
	}, nil
}

// GetBPFMapStats retrieves BPF map statistics.
//
// Parameters:
//   - mapName: Name of the BPF map.
//
// Returns ErrNotFound if the map does not exist.
func (lb *LoadBalancer) GetBPFMapStats(mapName string) (BPFMapStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return BPFMapStats{}, err
	}

	cName := C.CString(mapName)
	defer C.free(unsafe.Pointer(cName))

	var stats C.katran_bpf_map_stats_t
	ret := C.katran_lb_get_bpf_map_stats(lb.handle, cName, &stats)
	if ret != C.KATRAN_OK {
		return BPFMapStats{}, errorFromCode(ret)
	}

	return BPFMapStats{
		MaxEntries:     uint32(stats.max_entries),
		CurrentEntries: uint32(stats.current_entries),
	}, nil
}

// GetUserspaceStats retrieves userspace library statistics.
func (lb *LoadBalancer) GetUserspaceStats() (UserspaceStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return UserspaceStats{}, err
	}

	var stats C.katran_userspace_stats_t
	ret := C.katran_lb_get_userspace_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return UserspaceStats{}, errorFromCode(ret)
	}

	return UserspaceStats{
		BPFFailedCalls:       uint64(stats.bpf_failed_calls),
		AddrValidationFailed: uint64(stats.addr_validation_failed),
	}, nil
}

// GetPerCorePacketsStats retrieves per-core packet statistics.
//
// Returns a slice of packet counts, one per CPU core.
func (lb *LoadBalancer) GetPerCorePacketsStats() ([]int64, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_per_core_packets_stats(lb.handle, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	counts := make([]C.int64_t, count)
	ret = C.katran_lb_get_per_core_packets_stats(lb.handle, &counts[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]int64, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = int64(counts[i])
	}

	return result, nil
}

// IsUnderFlood checks if the system is under flood conditions.
//
// Examines connection rate statistics to determine if the system
// is experiencing a traffic flood.
//
// Returns true if under flood, false otherwise.
func (lb *LoadBalancer) IsUnderFlood() (bool, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return false, err
	}

	var result C.int
	ret := C.katran_lb_is_under_flood(lb.handle, &result)
	if ret != C.KATRAN_OK {
		return false, errorFromCode(ret)
	}

	return result != 0, nil
}
