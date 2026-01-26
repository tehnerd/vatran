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

// DeleteLRU deletes an LRU entry for a specific flow.
//
// Removes the connection tracking entry for the specified flow from
// all per-CPU and fallback LRU maps.
//
// Parameters:
//   - dstVip: Destination VIP.
//   - srcIP: Source IP address.
//   - srcPort: Source port.
//
// Returns a slice of map names where the entry was deleted.
func (lb *LoadBalancer) DeleteLRU(dstVip VIPKey, srcIP string, srcPort uint16) ([]string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	cVipAddr := C.CString(dstVip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(dstVip.Port),
		proto:   C.uint8_t(dstVip.Proto),
	}

	cSrcIP := C.CString(srcIP)
	defer C.free(unsafe.Pointer(cSrcIP))

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_delete_lru(lb.handle, &cVip, cSrcIP, C.uint16_t(srcPort), nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get map names
	maps := make([]*C.char, count)
	ret = C.katran_lb_delete_lru(lb.handle, &cVip, cSrcIP, C.uint16_t(srcPort), &maps[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]string, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = C.GoString(maps[i])
	}

	// Free C-allocated strings
	C.katran_free_strings(&maps[0], count)
	return result, nil
}

// PurgeVIPLRU purges all LRU entries for a VIP.
//
// Removes all connection tracking entries for the specified VIP
// from all LRU maps. Useful when removing a VIP or for cache invalidation.
//
// Parameters:
//   - dstVip: VIP to purge entries for.
//
// Returns the count of deleted entries.
func (lb *LoadBalancer) PurgeVIPLRU(dstVip VIPKey) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return 0, err
	}

	cVipAddr := C.CString(dstVip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(dstVip.Port),
		proto:   C.uint8_t(dstVip.Proto),
	}

	var deletedCount C.int
	ret := C.katran_lb_purge_vip_lru(lb.handle, &cVip, &deletedCount)
	if ret != C.KATRAN_OK {
		return 0, errorFromCode(ret)
	}

	return int(deletedCount), nil
}
