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

// AddRealForVIP adds a real server to a VIP.
//
// Adds a backend server to the VIP's consistent hash ring with the
// specified weight. Higher weights result in more traffic.
//
// Parameters:
//   - real: Real server to add (address, weight, flags).
//   - vip: VIP to add the real to.
//
// Returns ErrNotFound if the VIP does not exist, or ErrSpaceExhausted
// if the max_reals limit is reached.
func (lb *LoadBalancer) AddRealForVIP(real Real, vip VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cRealAddr := C.CString(real.Address)
	defer C.free(unsafe.Pointer(cRealAddr))

	cVipAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cReal := C.katran_new_real_t{
		address: cRealAddr,
		weight:  C.uint32_t(real.Weight),
		flags:   C.uint8_t(real.Flags),
	}

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	ret := C.katran_lb_add_real_for_vip(lb.handle, &cReal, &cVip)
	return errorFromCode(ret)
}

// DelRealForVIP removes a real server from a VIP.
//
// Removes the backend server from the VIP's consistent hash ring.
// The real's weight field is ignored for deletion.
//
// Parameters:
//   - real: Real server to remove.
//   - vip: VIP to remove the real from.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) DelRealForVIP(real Real, vip VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cRealAddr := C.CString(real.Address)
	defer C.free(unsafe.Pointer(cRealAddr))

	cVipAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cReal := C.katran_new_real_t{
		address: cRealAddr,
		weight:  C.uint32_t(real.Weight),
		flags:   C.uint8_t(real.Flags),
	}

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	ret := C.katran_lb_del_real_for_vip(lb.handle, &cReal, &cVip)
	return errorFromCode(ret)
}

// GetRealsForVIP retrieves all real servers for a VIP.
//
// Returns the list of backend servers configured for a VIP along
// with their weights and flags.
//
// Parameters:
//   - vip: VIP to query.
//
// Returns ErrNotFound if the VIP does not exist.
func (lb *LoadBalancer) GetRealsForVIP(vip VIPKey) ([]Real, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	cVipAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_reals_for_vip(lb.handle, &cVip, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	cReals := make([]C.katran_new_real_t, count)
	ret = C.katran_lb_get_reals_for_vip(lb.handle, &cVip, &cReals[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]Real, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = Real{
			Address: C.GoString(cReals[i].address),
			Weight:  uint32(cReals[i].weight),
			Flags:   uint8(cReals[i].flags),
		}
	}

	// Free C-allocated strings
	C.katran_free_reals(&cReals[0], count)
	return result, nil
}

// ModifyRealsForVIP performs a batch modification of real servers for a VIP.
//
// Adds or removes multiple real servers in a single operation.
// More efficient than individual add/del calls for bulk updates.
//
// Parameters:
//   - action: ActionAdd or ActionDel.
//   - reals: Array of real servers to modify.
//   - vip: VIP to modify.
//
// Returns ErrNotFound if the VIP does not exist, or ErrSpaceExhausted
// if adding would exceed max_reals.
func (lb *LoadBalancer) ModifyRealsForVIP(action ModifyAction, reals []Real, vip VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	if len(reals) == 0 {
		return nil
	}

	cVipAddr := C.CString(vip.Address)
	defer C.free(unsafe.Pointer(cVipAddr))

	cVip := C.katran_vip_key_t{
		address: cVipAddr,
		port:    C.uint16_t(vip.Port),
		proto:   C.uint8_t(vip.Proto),
	}

	// Allocate C strings for real addresses
	cReals := make([]C.katran_new_real_t, len(reals))
	cRealAddrs := make([]*C.char, len(reals))
	for i, real := range reals {
		cRealAddrs[i] = C.CString(real.Address)
		cReals[i] = C.katran_new_real_t{
			address: cRealAddrs[i],
			weight:  C.uint32_t(real.Weight),
			flags:   C.uint8_t(real.Flags),
		}
	}
	defer func() {
		for _, addr := range cRealAddrs {
			C.free(unsafe.Pointer(addr))
		}
	}()

	ret := C.katran_lb_modify_reals_for_vip(
		lb.handle,
		C.katran_modify_action_t(action),
		&cReals[0],
		C.size_t(len(reals)),
		&cVip,
	)
	return errorFromCode(ret)
}

// GetIndexForReal retrieves the internal index for a real server address.
//
// This index can be used with functions like GetRealStats().
//
// Parameters:
//   - address: Real server IP address.
//
// Returns the index, or -1 if the real is not found.
func (lb *LoadBalancer) GetIndexForReal(address string) (int64, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return -1, err
	}

	cAddr := C.CString(address)
	defer C.free(unsafe.Pointer(cAddr))

	var index C.int64_t
	ret := C.katran_lb_get_index_for_real(lb.handle, cAddr, &index)
	if ret != C.KATRAN_OK {
		return -1, errorFromCode(ret)
	}

	return int64(index), nil
}

// ModifyReal modifies flags on a real server.
//
// Sets or clears flags on a real server globally (affects all VIPs
// using this real).
//
// Parameters:
//   - address: Real server IP address.
//   - flags: Flag bits to modify.
//   - set: true to set the flags, false to clear them.
//
// Returns ErrNotFound if the real does not exist.
func (lb *LoadBalancer) ModifyReal(address string, flags uint8, set bool) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(address)
	defer C.free(unsafe.Pointer(cAddr))

	setInt := C.int(0)
	if set {
		setInt = 1
	}

	ret := C.katran_lb_modify_real(lb.handle, cAddr, C.uint8_t(flags), setInt)
	return errorFromCode(ret)
}
