package hcservice

import (
	"context"
	"net"
	"syscall"
	"time"
)

// SomarkDialer creates network connections with SO_MARK set on the socket.
// This allows katran's BPF program to route healthcheck packets to the correct real
// via tunnel encapsulation.
type SomarkDialer struct {
	// Timeout is the connection timeout.
	Timeout time.Duration
}

// NewSomarkDialer creates a new SomarkDialer.
//
// Parameters:
//   - timeout: The dial timeout duration.
//
// Returns a new SomarkDialer instance.
func NewSomarkDialer(timeout time.Duration) *SomarkDialer {
	return &SomarkDialer{Timeout: timeout}
}

// DialContext dials the given address with SO_MARK set on the socket.
//
// Parameters:
//   - ctx: Context for cancellation and deadline.
//   - network: The network type (e.g., "tcp", "tcp4", "tcp6").
//   - addr: The address to connect to (e.g., "10.0.0.1:80").
//   - somark: The SO_MARK value to set on the socket.
//
// Returns the established connection or an error.
func (d *SomarkDialer) DialContext(ctx context.Context, network, addr string, somark int) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: d.Timeout,
		Control: func(network, address string, c syscall.RawConn) error {
			var sErr error
			if err := c.Control(func(fd uintptr) {
				sErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, somark)
			}); err != nil {
				return err
			}
			return sErr
		},
	}
	return dialer.DialContext(ctx, network, addr)
}
