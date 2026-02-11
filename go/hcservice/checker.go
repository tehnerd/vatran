package hcservice

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/tehnerd/vatran/go/server/types"
)

// CheckResult represents the outcome of a single health check.
type CheckResult struct {
	// Success indicates whether the check passed.
	Success bool
	// Error is the error message if the check failed.
	Error string
}

// Checker performs a health check against a VIP address using a specific somark.
type Checker interface {
	// Check performs a health check.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout.
	//   - vipAddr: The VIP address to dial (katran routes via somark).
	//   - checkPort: The port to check on.
	//   - somark: The SO_MARK value for routing to the real.
	//   - config: The healthcheck configuration.
	//
	// Returns the check result.
	Check(ctx context.Context, vipAddr string, checkPort int, somark int, config *types.HealthcheckConfig) CheckResult
}

// TCPChecker performs TCP connection health checks.
type TCPChecker struct {
	dialer *SomarkDialer
}

// NewTCPChecker creates a new TCPChecker.
//
// Parameters:
//   - dialer: The somark dialer for creating marked connections.
//
// Returns a new TCPChecker instance.
func NewTCPChecker(dialer *SomarkDialer) *TCPChecker {
	return &TCPChecker{dialer: dialer}
}

// Check performs a TCP connection check. A successful connection means the real is healthy.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - vipAddr: The VIP address to dial.
//   - checkPort: The port to check on.
//   - somark: The SO_MARK value.
//   - config: The healthcheck configuration (unused for TCP).
//
// Returns the check result.
func (c *TCPChecker) Check(ctx context.Context, vipAddr string, checkPort int, somark int, config *types.HealthcheckConfig) CheckResult {
	addr := net.JoinHostPort(vipAddr, fmt.Sprintf("%d", checkPort))
	conn, err := c.dialer.DialContext(ctx, "tcp", addr, somark)
	if err != nil {
		return CheckResult{Success: false, Error: err.Error()}
	}
	conn.Close()
	return CheckResult{Success: true}
}

// HTTPChecker performs HTTP GET health checks.
type HTTPChecker struct {
	dialer *SomarkDialer
}

// NewHTTPChecker creates a new HTTPChecker.
//
// Parameters:
//   - dialer: The somark dialer for creating marked connections.
//
// Returns a new HTTPChecker instance.
func NewHTTPChecker(dialer *SomarkDialer) *HTTPChecker {
	return &HTTPChecker{dialer: dialer}
}

// Check performs an HTTP GET check against the VIP address.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - vipAddr: The VIP address to dial.
//   - checkPort: The port to check on.
//   - somark: The SO_MARK value.
//   - config: The healthcheck configuration containing HTTP settings.
//
// Returns the check result.
func (c *HTTPChecker) Check(ctx context.Context, vipAddr string, checkPort int, somark int, config *types.HealthcheckConfig) CheckResult {
	transport := &http.Transport{
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return c.dialer.DialContext(ctx, network, addr, somark)
		},
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("http://%s:%d%s", vipAddr, checkPort, config.HTTP.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CheckResult{Success: false, Error: err.Error()}
	}
	if config.HTTP.Host != "" {
		req.Host = config.HTTP.Host
	}

	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{Success: false, Error: err.Error()}
	}
	resp.Body.Close()

	if resp.StatusCode != config.HTTP.ExpectedStatus {
		return CheckResult{
			Success: false,
			Error:   fmt.Sprintf("unexpected status %d (expected %d)", resp.StatusCode, config.HTTP.ExpectedStatus),
		}
	}
	return CheckResult{Success: true}
}

// HTTPSChecker performs HTTPS GET health checks.
type HTTPSChecker struct {
	dialer *SomarkDialer
}

// NewHTTPSChecker creates a new HTTPSChecker.
//
// Parameters:
//   - dialer: The somark dialer for creating marked connections.
//
// Returns a new HTTPSChecker instance.
func NewHTTPSChecker(dialer *SomarkDialer) *HTTPSChecker {
	return &HTTPSChecker{dialer: dialer}
}

// Check performs an HTTPS GET check against the VIP address.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - vipAddr: The VIP address to dial.
//   - checkPort: The port to check on.
//   - somark: The SO_MARK value.
//   - config: The healthcheck configuration containing HTTPS settings.
//
// Returns the check result.
func (c *HTTPSChecker) Check(ctx context.Context, vipAddr string, checkPort int, somark int, config *types.HealthcheckConfig) CheckResult {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: config.HTTPS.SkipTLSVerify,
	}
	if config.HTTPS.Host != "" {
		tlsCfg.ServerName = config.HTTPS.Host
	}

	transport := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   tlsCfg,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return c.dialer.DialContext(ctx, network, addr, somark)
		},
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("https://%s:%d%s", vipAddr, checkPort, config.HTTPS.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CheckResult{Success: false, Error: err.Error()}
	}
	if config.HTTPS.Host != "" {
		req.Host = config.HTTPS.Host
	}

	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{Success: false, Error: err.Error()}
	}
	resp.Body.Close()

	if resp.StatusCode != config.HTTPS.ExpectedStatus {
		return CheckResult{
			Success: false,
			Error:   fmt.Sprintf("unexpected status %d (expected %d)", resp.StatusCode, config.HTTPS.ExpectedStatus),
		}
	}
	return CheckResult{Success: true}
}

// DummyChecker always returns success without performing any check.
type DummyChecker struct{}

// NewDummyChecker creates a new DummyChecker.
//
// Returns a new DummyChecker instance.
func NewDummyChecker() *DummyChecker {
	return &DummyChecker{}
}

// Check always returns a successful result.
//
// Parameters:
//   - ctx: Context (unused).
//   - vipAddr: The VIP address (unused).
//   - checkPort: The port (unused).
//   - somark: The SO_MARK value (unused).
//   - config: The healthcheck configuration (unused).
//
// Returns a successful check result.
func (c *DummyChecker) Check(ctx context.Context, vipAddr string, checkPort int, somark int, config *types.HealthcheckConfig) CheckResult {
	return CheckResult{Success: true}
}

// NewChecker creates the appropriate Checker for a given healthcheck type.
//
// Parameters:
//   - hcType: The healthcheck type ("tcp", "http", "https", "dummy").
//   - dialer: The somark dialer for creating marked connections.
//
// Returns the appropriate Checker implementation.
func NewChecker(hcType string, dialer *SomarkDialer) Checker {
	switch hcType {
	case "tcp":
		return NewTCPChecker(dialer)
	case "http":
		return NewHTTPChecker(dialer)
	case "https":
		return NewHTTPSChecker(dialer)
	default:
		return NewDummyChecker()
	}
}

// checkerTimeout returns the check timeout as a time.Duration.
func checkerTimeout(config *types.HealthcheckConfig) time.Duration {
	return time.Duration(config.TimeoutMs) * time.Millisecond
}
