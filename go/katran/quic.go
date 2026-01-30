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

// ModifyQuicRealsMapping modifies QUIC connection ID to real server mappings.
//
// Adds or removes mappings between QUIC host IDs (embedded in connection IDs)
// and real server addresses. This enables stateful QUIC routing.
//
// Parameters:
//   - action: ActionAdd or ActionDel.
//   - reals: Array of QUIC real mappings.
func (lb *LoadBalancer) ModifyQuicRealsMapping(action ModifyAction, reals []QuicReal) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	if len(reals) == 0 {
		return nil
	}

	// Allocate C strings for real addresses
	cReals := make([]C.katran_quic_real_t, len(reals))
	cRealAddrs := make([]*C.char, len(reals))
	for i, real := range reals {
		cRealAddrs[i] = C.CString(real.Address)
		cReals[i] = C.katran_quic_real_t{
			address: cRealAddrs[i],
			id:      C.uint32_t(real.ID),
		}
	}
	defer func() {
		for _, addr := range cRealAddrs {
			C.free(unsafe.Pointer(addr))
		}
	}()

	ret := C.katran_lb_modify_quic_reals_mapping(
		lb.handle,
		C.katran_modify_action_t(action),
		&cReals[0],
		C.size_t(len(reals)),
	)
	return errorFromCode(ret)
}

// GetQuicRealsMapping retrieves all QUIC connection ID mappings.
//
// Returns all configured QUIC host ID to real server mappings.
func (lb *LoadBalancer) GetQuicRealsMapping() ([]QuicReal, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_quic_reals_mapping(lb.handle, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	cReals := make([]C.katran_quic_real_t, count)
	ret = C.katran_lb_get_quic_reals_mapping(lb.handle, &cReals[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]QuicReal, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = QuicReal{
			Address: C.GoString(cReals[i].address),
			ID:      uint32(cReals[i].id),
		}
	}

	// Free C-allocated strings
	C.katran_free_quic_reals(&cReals[0], count)
	return result, nil
}
