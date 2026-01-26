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

// AddSrcRoutingRule adds source-based routing rules.
//
// Configures routing based on source IP prefixes. Packets from matching
// source addresses are forwarded to the specified destination.
//
// Requires FeatureSrcRouting to be enabled.
//
// Parameters:
//   - srcPrefixes: Array of source IP prefixes (CIDR notation, e.g., "10.0.0.0/8").
//   - dst: Destination address for matching traffic.
//
// Returns ErrFeatureDisabled if source routing is not enabled, or
// ErrSpaceExhausted if max_lpm_src_size is exceeded.
func (lb *LoadBalancer) AddSrcRoutingRule(srcPrefixes []string, dst string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	if len(srcPrefixes) == 0 {
		return nil
	}

	// Allocate C strings for source prefixes
	cPrefixes := make([]*C.char, len(srcPrefixes))
	for i, prefix := range srcPrefixes {
		cPrefixes[i] = C.CString(prefix)
	}
	defer func() {
		for _, p := range cPrefixes {
			C.free(unsafe.Pointer(p))
		}
	}()

	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cDst))

	ret := C.katran_lb_add_src_routing_rule(
		lb.handle,
		&cPrefixes[0],
		C.size_t(len(srcPrefixes)),
		cDst,
	)
	return errorFromCode(ret)
}

// DelSrcRoutingRule deletes source-based routing rules.
//
// Removes routing rules for the specified source prefixes.
//
// Parameters:
//   - srcPrefixes: Array of source IP prefixes to remove.
//
// Returns ErrFeatureDisabled if source routing is not enabled.
func (lb *LoadBalancer) DelSrcRoutingRule(srcPrefixes []string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	if len(srcPrefixes) == 0 {
		return nil
	}

	// Allocate C strings for source prefixes
	cPrefixes := make([]*C.char, len(srcPrefixes))
	for i, prefix := range srcPrefixes {
		cPrefixes[i] = C.CString(prefix)
	}
	defer func() {
		for _, p := range cPrefixes {
			C.free(unsafe.Pointer(p))
		}
	}()

	ret := C.katran_lb_del_src_routing_rule(
		lb.handle,
		&cPrefixes[0],
		C.size_t(len(srcPrefixes)),
	)
	return errorFromCode(ret)
}

// ClearAllSrcRoutingRules removes all source-based routing rules.
//
// Returns ErrFeatureDisabled if source routing is not enabled.
func (lb *LoadBalancer) ClearAllSrcRoutingRules() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	ret := C.katran_lb_clear_all_src_routing_rules(lb.handle)
	return errorFromCode(ret)
}

// SrcRoutingRule represents a source routing rule.
type SrcRoutingRule struct {
	// Src is the source IP prefix (CIDR notation).
	Src string
	// Dst is the destination address.
	Dst string
}

// GetSrcRoutingRules retrieves all source-based routing rules.
//
// Returns all configured source prefix to destination mappings.
// Returns ErrFeatureDisabled if source routing is not enabled.
func (lb *LoadBalancer) GetSrcRoutingRules() ([]SrcRoutingRule, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_src_routing_rule(lb.handle, nil, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	srcs := make([]*C.char, count)
	dsts := make([]*C.char, count)
	ret = C.katran_lb_get_src_routing_rule(lb.handle, &srcs[0], &dsts[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]SrcRoutingRule, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = SrcRoutingRule{
			Src: C.GoString(srcs[i]),
			Dst: C.GoString(dsts[i]),
		}
	}

	// Free C-allocated strings
	C.katran_free_src_routing_rules(&srcs[0], &dsts[0], count)
	return result, nil
}

// GetSrcRoutingRuleSize returns the number of source routing rules.
func (lb *LoadBalancer) GetSrcRoutingRuleSize() (uint32, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return 0, err
	}

	var size C.uint32_t
	ret := C.katran_lb_get_src_routing_rule_size(lb.handle, &size)
	if ret != C.KATRAN_OK {
		return 0, errorFromCode(ret)
	}

	return uint32(size), nil
}

// AddInlineDecapDst adds an inline decapsulation destination.
//
// Configures an IP address for which incoming encapsulated packets
// should be decapsulated in the XDP program.
//
// Requires FeatureInlineDecap to be enabled.
//
// Parameters:
//   - dst: Destination IP address to enable decapsulation for.
//
// Returns ErrFeatureDisabled if inline decap is not enabled, or
// ErrSpaceExhausted if max_decap_dst is exceeded.
func (lb *LoadBalancer) AddInlineDecapDst(dst string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cDst))

	ret := C.katran_lb_add_inline_decap_dst(lb.handle, cDst)
	return errorFromCode(ret)
}

// DelInlineDecapDst removes an inline decapsulation destination.
//
// Parameters:
//   - dst: Destination IP address to remove.
//
// Returns ErrFeatureDisabled if inline decap is not enabled, or
// ErrNotFound if the destination was not configured.
func (lb *LoadBalancer) DelInlineDecapDst(dst string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	cDst := C.CString(dst)
	defer C.free(unsafe.Pointer(cDst))

	ret := C.katran_lb_del_inline_decap_dst(lb.handle, cDst)
	return errorFromCode(ret)
}

// GetInlineDecapDsts retrieves all inline decapsulation destinations.
//
// Returns all configured decapsulation destination addresses.
func (lb *LoadBalancer) GetInlineDecapDsts() ([]string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return nil, err
	}

	// Phase 1: Get count
	var count C.size_t
	ret := C.katran_lb_get_inline_decap_dst(lb.handle, nil, &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}
	if count == 0 {
		return nil, nil
	}

	// Phase 2: Get data
	dsts := make([]*C.char, count)
	ret = C.katran_lb_get_inline_decap_dst(lb.handle, &dsts[0], &count)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	// Convert to Go types
	result := make([]string, count)
	for i := C.size_t(0); i < count; i++ {
		result[i] = C.GoString(dsts[i])
	}

	// Free C-allocated strings
	C.katran_free_strings(&dsts[0], count)
	return result, nil
}
