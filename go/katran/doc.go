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

// Package katran provides Go bindings for the Katran XDP-based L4 load balancer.
//
// Katran is a high-performance layer 4 load balancer based on XDP (eXpress Data Path)
// and BPF (Berkeley Packet Filter). This package provides a Go-native interface to
// the Katran C API via CGO.
//
// # Basic Usage
//
// Create a load balancer instance with configuration:
//
//	cfg := katran.NewConfig()
//	cfg.MainInterface = "eth0"
//	cfg.BalancerProgPath = "/path/to/balancer.o"
//	cfg.DefaultMAC = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
//
//	lb, err := katran.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer lb.Close()
//
//	// Load and attach BPF programs
//	if err := lb.LoadBPFProgs(); err != nil {
//	    log.Fatal(err)
//	}
//	if err := lb.AttachBPFProgs(); err != nil {
//	    log.Fatal(err)
//	}
//
// # VIP and Real Server Management
//
// Add a VIP and configure backend servers:
//
//	vip := katran.VIPKey{
//	    Address: "10.0.0.1",
//	    Port:    80,
//	    Proto:   6, // TCP
//	}
//
//	if err := lb.AddVIP(vip, 0); err != nil {
//	    log.Fatal(err)
//	}
//
//	real := katran.Real{
//	    Address: "192.168.1.10",
//	    Weight:  100,
//	}
//
//	if err := lb.AddRealForVIP(real, vip); err != nil {
//	    log.Fatal(err)
//	}
//
// # Statistics
//
// Query load balancer statistics:
//
//	stats, err := lb.GetStatsForVIP(vip)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Packets: %d, Bytes: %d\n", stats.V1, stats.V2)
//
// # Thread Safety
//
// The LoadBalancer type provides internal synchronization via a mutex.
// It is safe to call methods from multiple goroutines concurrently.
//
// # Memory Management
//
// The Go garbage collector will automatically clean up the load balancer
// when it is no longer referenced. A finalizer is registered to call Close()
// if the user forgets. However, it is recommended to explicitly call Close()
// when done to ensure timely cleanup of BPF resources.
package katran
