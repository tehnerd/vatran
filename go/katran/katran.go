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
#cgo CFLAGS: -I${SRCDIR}/../../include
#cgo LDFLAGS: -Wl,--start-group  -lfolly -lglog -lgflags -lpthread -ldl -lbpf -lstdc++ -lm -lelf -lz -levent -liberty -ldouble-conversion -lzstd -lmnl -lfmt -lunwind
#cgo LDFLAGS: -L${SRCDIR}/../../_build_go -lkatran_capi_static -lkatranlb -lpcapwriter -lbpfadapter -liphelpers -lchhelpers -lkatransimulator -lmurmur3 -Wl,--end-group
#include "katran_capi.h"
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

// LoadBalancer represents a Katran load balancer instance.
// It provides thread-safe access to all load balancer operations.
type LoadBalancer struct {
	handle C.katran_lb_t
	mu     sync.Mutex
	closed bool
}

// New creates a new Katran load balancer instance with the provided configuration.
//
// After creation, call LoadBPFProgs() and AttachBPFProgs() to start load balancing.
// When done, call Close() to release resources.
//
// Parameters:
//   - cfg: Configuration for the load balancer. Must not be nil.
//
// Returns an error if creation fails.
func New(cfg *Config) (*LoadBalancer, error) {
	if cfg == nil {
		return nil, &KatranError{Code: ErrInvalidArgument, Message: "config cannot be nil"}
	}

	cc := cfg.toC()
	defer cc.free()

	var handle C.katran_lb_t
	ret := C.katran_lb_create(&cc.config, &handle)
	if ret != C.KATRAN_OK {
		return nil, errorFromCode(ret)
	}

	lb := &LoadBalancer{
		handle: handle,
	}

	// Set finalizer to clean up if user forgets to call Close()
	runtime.SetFinalizer(lb, func(lb *LoadBalancer) {
		lb.Close()
	})

	return lb, nil
}

// Close destroys the load balancer instance and releases all resources.
//
// After calling Close(), the LoadBalancer must not be used.
// It is safe to call Close() multiple times.
func (lb *LoadBalancer) Close() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.closed {
		return nil
	}

	lb.closed = true
	runtime.SetFinalizer(lb, nil)

	ret := C.katran_lb_destroy(lb.handle)
	lb.handle = nil

	return errorFromCode(ret)
}

// LoadBPFProgs loads BPF programs into the kernel.
//
// This must be called before AttachBPFProgs(). The programs loaded are
// specified in the configuration (balancer_prog_path and healthchecking_prog_path).
func (lb *LoadBalancer) LoadBPFProgs() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.closed {
		return &KatranError{Code: ErrInvalidArgument, Message: "load balancer is closed"}
	}

	ret := C.katran_lb_load_bpf_progs(lb.handle)
	return errorFromCode(ret)
}

// AttachBPFProgs attaches loaded BPF programs to network interfaces.
//
// This must be called after LoadBPFProgs(). The programs are attached to
// the interfaces specified in the configuration.
func (lb *LoadBalancer) AttachBPFProgs() error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.closed {
		return &KatranError{Code: ErrInvalidArgument, Message: "load balancer is closed"}
	}

	ret := C.katran_lb_attach_bpf_progs(lb.handle)
	return errorFromCode(ret)
}

// ReloadBalancerProg reloads the balancer BPF program at runtime.
//
// This allows hot-reloading of the balancer program without service interruption.
//
// Parameters:
//   - path: Path to the new BPF program file.
//   - cfg: Optional new configuration. Pass nil to keep current config.
func (lb *LoadBalancer) ReloadBalancerProg(path string, cfg *Config) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.closed {
		return &KatranError{Code: ErrInvalidArgument, Message: "load balancer is closed"}
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var ret C.katran_error_t
	if cfg != nil {
		cc := cfg.toC()
		defer cc.free()
		ret = C.katran_lb_reload_balancer_prog(lb.handle, cPath, &cc.config)
	} else {
		ret = C.katran_lb_reload_balancer_prog(lb.handle, cPath, nil)
	}

	return errorFromCode(ret)
}

// checkClosed returns an error if the load balancer is closed.
// Must be called with lb.mu held.
func (lb *LoadBalancer) checkClosed() error {
	if lb.closed {
		return &KatranError{Code: ErrInvalidArgument, Message: "load balancer is closed"}
	}
	return nil
}
