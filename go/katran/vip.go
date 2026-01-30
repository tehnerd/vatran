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

// AddVIP adds a new Virtual IP (VIP) to the load balancer.
//
// After adding a VIP, use AddRealForVIP() or ModifyRealsForVIP() to
// configure backend servers.
//
// Parameters:
//   - vip: VIP to add (address, port, protocol).
//   - flags: VIP flags (e.g., NO_PORT, NO_LRU). Pass 0 for defaults.
//
// Returns ErrAlreadyExists if the VIP already exists, or ErrSpaceExhausted
// if the max_vips limit is reached.
func (lb *LoadBalancer) AddVIP(vip VIPKey, flags uint32) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	ret := C.katran_lb_add_vip(lb.handle, &cVip, C.uint32_t(flags))
	return errorFromCode(ret)
}

// DelVIP deletes a VIP from the load balancer.
//
// This removes the VIP and all associated real server configurations.
// Traffic to this VIP will no longer be load balanced.
//
// Parameters:
//   - vip: VIP to delete.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) DelVIP(vip VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	ret := C.katran_lb_del_vip(lb.handle, &cVip)
	return errorFromCode(ret)
}

// GetAllVIPs retrieves all configured VIPs.
//
// Returns a slice of all VIPs currently configured in the load balancer.
func (lb *LoadBalancer) GetAllVIPs() ([]VIPKey, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_all_vips(lb.handle, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	cVips := make([]C.katran_vip_key_t, count)
	ret = C.katran_lb_get_all_vips(lb.handle, &cVips[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]VIPKey, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = VIPKey{
			Address: C.GoString(cVips[i].address),
			Port:    uint16(cVips[i].port),
			Proto:   uint8(cVips[i].proto),
		}
	}

	// Free C-allocated strings
	C.katran_free_vips(&cVips[0], count)
	return result, nil
}

// ModifyVIP modifies a VIP's flags.
//
// Sets or clears specific flags on an existing VIP.
//
// Parameters:
//   - vip: VIP to modify.
//   - flag: Flag bits to modify.
//   - set: true to set the flags, false to clear them.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) ModifyVIP(vip VIPKey, flag uint32, set bool) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	setInt := C.int(0)
	if set {
		setInt = 1
	}

	ret := C.katran_lb_modify_vip(lb.handle, &cVip, C.uint32_t(flag), setInt)
	return errorFromCode(ret)
}

// GetVIPFlags retrieves a VIP's current flags.
//
// Parameters:
//   - vip: VIP to query.
//
// Returns the flag bits currently set on the VIP, or ErrNotFound if the VIP
// does not exist.
func (lb *LoadBalancer) GetVIPFlags(vip VIPKey) (uint32, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return 0, err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	var flags C.uint32_t
	ret := C.katran_lb_get_vip_flags(lb.handle, &cVip, &flags)
	if ret != C.KATRAN_OK {
		return 0, errorFromCode(ret)
	}

	return uint32(flags), nil
}

// ChangeHashFunctionForVIP changes the hash function for a VIP's consistent hashing ring.
//
// This triggers a recalculation of the entire ring.
//
// Parameters:
//   - vip: VIP to modify.
//   - fn: Hash function to use.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) ChangeHashFunctionForVIP(vip VIPKey, fn HashFunction) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cVip := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	ret := C.katran_lb_change_hash_function_for_vip(lb.handle, &cVip, C.katran_hash_function_t(fn))
	return errorFromCode(ret)
}
