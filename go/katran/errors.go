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
*/
import "C"
import "fmt"

// KatranError represents an error returned by the Katran API.
type KatranError struct {
	// Code is the error code from the C API.
	Code Error
	// Message is a human-readable error description.
	Message string
}

// Error implements the error interface.
func (e *KatranError) Error() string {
	return fmt.Sprintf("katran: %s (code: %d)", e.Message, e.Code)
}

// Is checks if the error matches the target error code.
func (e *KatranError) Is(target error) bool {
	if t, ok := target.(*KatranError); ok {
		return e.Code == t.Code
	}
	return false
}

// errorFromCode creates a KatranError from a C error code.
func errorFromCode(code C.katran_error_t) error {
	if code == C.KATRAN_OK {
		return nil
	}
	msg := C.GoString(C.katran_lb_get_last_error())
	return &KatranError{
		Code:    Error(code),
		Message: msg,
	}
}

// IsNotFound checks if the error indicates a resource was not found.
func IsNotFound(err error) bool {
	if e, ok := err.(*KatranError); ok {
		return e.Code == ErrNotFound
	}
	return false
}

// IsAlreadyExists checks if the error indicates the resource already exists.
func IsAlreadyExists(err error) bool {
	if e, ok := err.(*KatranError); ok {
		return e.Code == ErrAlreadyExists
	}
	return false
}

// IsSpaceExhausted checks if the error indicates maximum capacity was reached.
func IsSpaceExhausted(err error) bool {
	if e, ok := err.(*KatranError); ok {
		return e.Code == ErrSpaceExhausted
	}
	return false
}

// IsBPFFailed checks if the error indicates a BPF operation failed.
func IsBPFFailed(err error) bool {
	if e, ok := err.(*KatranError); ok {
		return e.Code == ErrBPFFailed
	}
	return false
}

// IsFeatureDisabled checks if the error indicates a feature is not enabled.
func IsFeatureDisabled(err error) bool {
	if e, ok := err.(*KatranError); ok {
		return e.Code == ErrFeatureDisabled
	}
	return false
}
