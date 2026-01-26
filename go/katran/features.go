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

// HasFeature checks if a feature is available.
//
// Queries whether a specific optional feature is enabled in the
// current BPF program.
//
// Parameters:
//   - feature: Feature flag to check.
//
// Returns true if the feature is available, false otherwise.
func (lb *LoadBalancer) HasFeature(feature Feature) (bool, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return false, err
	}

	var hasFeature C.int
	ret := C.katran_lb_has_feature(lb.handle, C.katran_feature_t(feature), &hasFeature)
	if ret != C.KATRAN_OK {
		return false, errorFromCode(ret)
	}

	return hasFeature != 0, nil
}

// InstallFeature installs a feature by reloading the BPF program.
//
// If the feature is not currently available, attempts to reload the
// BPF program from the specified path to enable it.
//
// Parameters:
//   - feature: Feature to install.
//   - progPath: Path to BPF program with the feature, or empty string to use
//     current program path.
//
// Returns ErrBPFFailed if reload fails or feature not in new program.
func (lb *LoadBalancer) InstallFeature(feature Feature, progPath string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	var cPath *C.char
	if progPath != "" {
		cPath = C.CString(progPath)
		defer C.free(unsafe.Pointer(cPath))
	}

	ret := C.katran_lb_install_feature(lb.handle, C.katran_feature_t(feature), cPath)
	return errorFromCode(ret)
}

// RemoveFeature removes a feature by reloading the BPF program.
//
// If the feature is currently available, attempts to reload the
// BPF program from the specified path to disable it.
//
// Parameters:
//   - feature: Feature to remove.
//   - progPath: Path to BPF program without the feature, or empty string.
//
// Returns ErrBPFFailed if reload fails or feature still in new program.
func (lb *LoadBalancer) RemoveFeature(feature Feature, progPath string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	var cPath *C.char
	if progPath != "" {
		cPath = C.CString(progPath)
		defer C.free(unsafe.Pointer(cPath))
	}

	ret := C.katran_lb_remove_feature(lb.handle, C.katran_feature_t(feature), cPath)
	return errorFromCode(ret)
}
