package lb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BGPClient is an HTTP client for communicating with the external BGP service.
type BGPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewBGPClient creates a new BGPClient.
//
// Parameters:
//   - baseURL: The base URL of the BGP service (e.g., "http://localhost:9100").
//
// Returns a new BGPClient instance.
func NewBGPClient(baseURL string) *BGPClient {
	return &BGPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// bgpResponse is the standard response wrapper from the BGP service.
type bgpResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *bgpError       `json:"error,omitempty"`
}

// bgpError is the error response from the BGP service.
type bgpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// bgpAdvertiseRequest is the request body for POST /api/v1/routes/advertise.
type bgpAdvertiseRequest struct {
	VIP       string `json:"vip"`
	PrefixLen uint8  `json:"prefix_len"`
}

// bgpWithdrawRequest is the request body for POST /api/v1/routes/withdraw.
type bgpWithdrawRequest struct {
	VIP       string `json:"vip"`
	PrefixLen uint8  `json:"prefix_len"`
}

// Advertise sends a route advertise request to the BGP service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP IP address to advertise.
//   - prefixLen: The prefix length (e.g., 32 for /32).
//
// Returns an error if the request fails.
func (c *BGPClient) Advertise(ctx context.Context, vip string, prefixLen uint8) error {
	body := bgpAdvertiseRequest{
		VIP:       vip,
		PrefixLen: prefixLen,
	}
	return c.doRequest(ctx, http.MethodPost, "/api/v1/routes/advertise", body)
}

// Withdraw sends a route withdraw request to the BGP service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP IP address to withdraw.
//   - prefixLen: The prefix length.
//
// Returns an error if the request fails.
func (c *BGPClient) Withdraw(ctx context.Context, vip string, prefixLen uint8) error {
	body := bgpWithdrawRequest{
		VIP:       vip,
		PrefixLen: prefixLen,
	}
	return c.doRequest(ctx, http.MethodPost, "/api/v1/routes/withdraw", body)
}

// doRequest performs an HTTP request to the BGP service.
func (c *BGPClient) doRequest(ctx context.Context, method, path string, body interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("BGP service request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read BGP service response: %w", err)
	}

	var bgpResp bgpResponse
	if err := json.Unmarshal(respBody, &bgpResp); err != nil {
		return fmt.Errorf("failed to parse BGP service response: %w", err)
	}

	if !bgpResp.Success {
		msg := "unknown error"
		if bgpResp.Error != nil {
			msg = bgpResp.Error.Message
		}
		return fmt.Errorf("BGP service error: %s", msg)
	}

	return nil
}
