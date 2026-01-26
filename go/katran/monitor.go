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

// StopMonitor stops the packet monitor.
//
// Stops packet capture/introspection if running.
//
// Returns ErrFeatureDisabled if introspection is not enabled.
func (lb *LoadBalancer) StopMonitor() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	ret := C.katran_lb_stop_monitor(lb.handle)
	return errorFromCode(ret)
}

// RestartMonitor restarts the packet monitor.
//
// Restarts packet capture with the specified packet limit.
//
// Parameters:
//   - limit: Maximum number of packets to capture (0 = unlimited).
//
// Returns ErrFeatureDisabled if introspection is not enabled.
func (lb *LoadBalancer) RestartMonitor(limit uint32) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return err
	}

	ret := C.katran_lb_restart_monitor(lb.handle, C.uint32_t(limit))
	return errorFromCode(ret)
}

// GetMonitorStats retrieves monitor statistics.
//
// Returns statistics from the packet capture subsystem.
//
// Returns ErrFeatureDisabled if introspection is not enabled.
func (lb *LoadBalancer) GetMonitorStats() (MonitorStats, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err := lb.checkClosed(); err != nil {
		return MonitorStats{}, err
	}

	var stats C.katran_monitor_stats_t
	ret := C.katran_lb_get_monitor_stats(lb.handle, &stats)
	if ret != C.KATRAN_OK {
		return MonitorStats{}, errorFromCode(ret)
	}

	return MonitorStats{
		Limit:      uint32(stats.limit),
		Amount:     uint32(stats.amount),
		BufferFull: uint32(stats.buffer_full),
	}, nil
}
