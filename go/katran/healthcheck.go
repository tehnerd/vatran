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

// AddHealthcheckerDst adds a healthcheck destination mapping.
//
// Maps a socket mark (SO_MARK) to a destination address for healthcheck
// packet encapsulation. Packets with the specified mark will be
// encapsulated and sent to the destination.
//
// Parameters:
//   - somark: Socket mark value to match.
//   - dst: Destination address for encapsulated packets.
//
// Returns ErrFeatureDisabled if healthchecking is not enabled.
func (lb *LoadBalancer) AddHealthcheckerDst(somark uint32, dst string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cDst))

	ret := C.katran_lb_add_healthchecker_dst(lb.handle, C.uint32_t(somark), cDst)
	return errorFromCode(ret)
}

// DelHealthcheckerDst removes a healthcheck destination mapping.
//
// Parameters:
//   - somark: Socket mark value to remove.
//
// Returns ErrNotFound if the somark was not configured.
func (lb *LoadBalancer) DelHealthcheckerDst(somark uint32) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	ret := C.katran_lb_del_healthchecker_dst(lb.handle, C.uint32_t(somark))
	return errorFromCode(ret)
}

// HealthcheckerDst represents a healthcheck destination mapping.
type HealthcheckerDst struct {
	// Somark is the socket mark value.
	Somark uint32
	// Dst is the destination address.
	Dst string
}

// GetHealthcheckersDst retrieves all healthcheck destination mappings.
//
// Returns all configured socket mark to destination mappings.
func (lb *LoadBalancer) GetHealthcheckersDst() ([]HealthcheckerDst, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_healthcheckers_dst(lb.handle, nil, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	somarks := make([]C.uint32_t, count)
	dsts := make([]*C.char, count)
	ret = C.katran_lb_get_healthcheckers_dst(lb.handle, &somarks[0], &dsts[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]HealthcheckerDst, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = HealthcheckerDst{
			Somark: uint32(somarks[i]),
			Dst:    C.GoString(dsts[i]),
		}
	}

	// Free C-allocated strings
	C.katran_free_strings(&dsts[0], count)
	return result, nil
}

// AddHCKey adds a healthcheck key for per-key statistics.
//
// Registers a VIP-like key for healthcheck packet tracking.
//
// Parameters:
//   - hcKey: Healthcheck key to add.
func (lb *LoadBalancer) AddHCKey(hcKey VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(hcKey.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cKey := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(hcKey.Port),
		proto:   C.uint8_t(hcKey.Proto),
	}

	ret := C.katran_lb_add_hc_key(lb.handle, &cKey)
	return errorFromCode(ret)
}

// DelHCKey removes a healthcheck key.
//
// Parameters:
//   - hcKey: Healthcheck key to remove.
//
// Returns ErrNotFound if the key was not registered.
func (lb *LoadBalancer) DelHCKey(hcKey VIPKey) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cAddr := C.CString(hcKey.Address)
	defer C.free(unsafe.Pointer(cAddr))

	cKey := C.katran_vip_key_t{
		address: cAddr,
		port:    C.uint16_t(hcKey.Port),
		proto:   C.uint8_t(hcKey.Proto),
	}

	ret := C.katran_lb_del_hc_key(lb.handle, &cKey)
	return errorFromCode(ret)
}
