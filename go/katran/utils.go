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

// ChangeMAC changes the default router MAC address.
//
// Updates the MAC address used as the destination for forwarded packets.
// This is typically the MAC of the default gateway or next-hop router.
//
// Parameters:
//   - mac: 6-byte MAC address.
func (lb *LoadBalancer) ChangeMAC(mac []byte) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	if len(mac) != 6 {
		return &KatranError{Code: ErrInvalidArgument, Message: "MAC address must be 6 bytes"}
	}

	ret := C.katran_lb_change_mac(lb.handle, (*C.uint8_t)(unsafe.Pointer(&mac[0])))
	return errorFromCode(ret)
}

// GetMAC retrieves the current default router MAC address.
//
// Returns the MAC address currently configured for packet forwarding.
func (lb *LoadBalancer) GetMAC() ([]byte, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	mac := make([]byte, 6)
	ret := C.katran_lb_get_mac(lb.handle, (*C.uint8_t)(unsafe.Pointer(&mac[0])))
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	return mac, nil
}

// GetRealForFlow determines which real server a flow would be routed to.
//
// Simulates the routing decision for a given 5-tuple without actually
// processing a packet. Useful for debugging and verification.
//
// Parameters:
//   - flow: Flow 5-tuple to simulate.
//
// Returns the real server address, or ErrNotFound if the flow doesn't match any VIP.
func (lb *LoadBalancer) GetRealForFlow(flow Flow) (string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return "", err
	}

	cSrc := C.CString(flow.Src)
	defer C.free(unsafe.Pointer(cSrc))

	cDst := C.CString(flow.Dst)
	defer C.free(unsafe.Pointer(cDst))

	cFlow := C.katran_flow_t{
		src:      cSrc,
		dst:      cDst,
		src_port: C.uint16_t(flow.SrcPort),
		dst_port: C.uint16_t(flow.DstPort),
		proto:    C.uint8_t(flow.Proto),
	}

	// Allocate buffer for real address (INET6_ADDRSTRLEN = 46)
	bufSize := C.size_t(46)
	buf := make([]byte, bufSize)

	ret := C.katran_lb_get_real_for_flow(lb.handle, &cFlow, (*C.char)(unsafe.Pointer(&buf[0])), bufSize)
	if ret != C.KATRAN_OK {
		return "", errorFromCode(ret)
	}

	// Find null terminator
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i]), nil
		}
	}
	return string(buf), nil
}

// SimulatePacket processes a raw packet through the BPF program.
//
// Processes a raw packet through the Katran BPF program and returns
// the resulting packet. Note: This affects BPF state (maps, stats).
//
// Parameters:
//   - inPacket: Input packet data (starting with Ethernet header).
//
// Returns the output packet data.
func (lb *LoadBalancer) SimulatePacket(inPacket []byte) ([]byte, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	if len(inPacket) == 0 {
		return nil, &KatranError{Code: ErrInvalidArgument, Message: "input packet cannot be empty"}
	}

	// Phase 1: Get output size
	var outSize C.size_t
	ret := C.katran_lb_simulate_packet(
		lb.handle,
		(*C.uint8_t)(unsafe.Pointer(&inPacket[0])),
		C.size_t(len(inPacket)),
		nil,
		&outSize,
	)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if outSize == 0 {
		return nil, nil
	}

	// Phase 2: Get output packet
	outPacket := make([]byte, outSize)
	ret = C.katran_lb_simulate_packet(
		lb.handle,
		(*C.uint8_t)(unsafe.Pointer(&inPacket[0])),
		C.size_t(len(inPacket)),
		(*C.uint8_t)(unsafe.Pointer(&outPacket[0])),
		&outSize,
	)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	return outPacket[:outSize], nil
}

// GetKatranProgFD returns the file descriptor of the balancer BPF program.
func (lb *LoadBalancer) GetKatranProgFD() (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return -1, err
	}

	var fd C.int
	ret := C.katran_lb_get_katran_prog_fd(lb.handle, &fd)
	if ret != C.KATRAN_OK {
		return -1, errorFromCode(ret)
	}

	return int(fd), nil
}

// GetHealthcheckerProgFD returns the file descriptor of the healthcheck BPF program.
//
// Returns ErrFeatureDisabled if healthchecking is not enabled.
func (lb *LoadBalancer) GetHealthcheckerProgFD() (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return -1, err
	}

	var fd C.int
	ret := C.katran_lb_get_healthchecker_prog_fd(lb.handle, &fd)
	if ret != C.KATRAN_OK {
		return -1, errorFromCode(ret)
	}

	return int(fd), nil
}

// GetBPFMapFDByName returns a BPF map file descriptor by name.
//
// Parameters:
//   - mapName: Name of the BPF map.
//
// Returns ErrNotFound if the map does not exist.
func (lb *LoadBalancer) GetBPFMapFDByName(mapName string) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return -1, err
	}

	cName := C.CString(mapName)
	defer C.free(unsafe.Pointer(cName))

	var fd C.int
	ret := C.katran_lb_get_bpf_map_fd_by_name(lb.handle, cName, &fd)
	if ret != C.KATRAN_OK {
		return -1, errorFromCode(ret)
	}

	return int(fd), nil
}

// GetGlobalLRUMapsFDs returns file descriptors for global LRU maps.
//
// Returns FDs for all per-CPU global LRU maps.
func (lb *LoadBalancer) GetGlobalLRUMapsFDs() ([]int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_global_lru_maps_fds(lb.handle, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get FDs
	fds := make([]C.int, count)
	ret = C.katran_lb_get_global_lru_maps_fds(lb.handle, &fds[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]int, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = int(fds[i])
	}

	return result, nil
}

// AddSrcIPForPcktEncap adds a source IP for packet encapsulation.
//
// Sets the source IP address to use when Katran encapsulates packets
// (GUE or IPIP). Replaces any existing source of the same address family.
//
// Parameters:
//   - src: Source IP address (IPv4 or IPv6).
func (lb *LoadBalancer) AddSrcIPForPcktEncap(src string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cSrc := C.CString(src)
	defer C.free(unsafe.Pointer(cSrc))

	ret := C.katran_lb_add_src_ip_for_pckt_encap(lb.handle, cSrc)
	return errorFromCode(ret)
}
