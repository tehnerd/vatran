package lb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tehnerd/vatran/go/server/types"
)

// HCClient is an HTTP client for communicating with the external healthcheck service.
type HCClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHCClient creates a new HCClient.
//
// Parameters:
//   - baseURL: The base URL of the healthcheck service (e.g., "http://localhost:9000").
//
// Returns a new HCClient instance.
func NewHCClient(baseURL string) *HCClient {
	return &HCClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// hcResponse is the standard response wrapper from the HC service.
type hcResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *hcError        `json:"error,omitempty"`
}

// hcError is the error response from the HC service.
type hcError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// registerRequest is the request body for POST /api/v1/targets.
type registerRequest struct {
	VIP         types.HCVIPKey         `json:"vip"`
	Reals       []realEntry            `json:"reals"`
	Healthcheck *types.HealthcheckConfig `json:"healthcheck"`
}

// updateRequest is the request body for PUT /api/v1/targets.
type updateRequest struct {
	VIP         types.HCVIPKey         `json:"vip"`
	Healthcheck *types.HealthcheckConfig `json:"healthcheck"`
	Reals       []realEntry            `json:"reals,omitempty"`
}

// realEntry is a real server entry sent to the HC service.
type realEntry struct {
	Address string `json:"address"`
	Weight  uint32 `json:"weight,omitempty"`
	Flags   uint8  `json:"flags,omitempty"`
}

// realsRequest is the request body for POST/DELETE /api/v1/targets/reals.
type realsRequest struct {
	VIP   types.HCVIPKey `json:"vip"`
	Reals []realEntry    `json:"reals"`
}

// RegisterVIP registers a VIP with the healthcheck service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//   - reals: The list of reals to register.
//   - hcConfig: The healthcheck configuration.
//
// Returns an error if the registration fails.
func (c *HCClient) RegisterVIP(ctx context.Context, vip types.HCVIPKey, reals []RealState, hcConfig *types.HealthcheckConfig) error {
	entries := make([]realEntry, len(reals))
	for i, r := range reals {
		entries[i] = realEntry{Address: r.Address, Weight: r.Weight, Flags: r.Flags}
	}
	body := registerRequest{
		VIP:         vip,
		Reals:       entries,
		Healthcheck: hcConfig,
	}
	return c.doRequest(ctx, http.MethodPost, "/api/v1/targets", body, nil)
}

// UpdateVIP updates the healthcheck configuration for a registered VIP.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//   - hcConfig: The new healthcheck configuration.
//   - reals: Optional list of reals to replace. Pass nil to keep existing.
//
// Returns an error if the update fails.
func (c *HCClient) UpdateVIP(ctx context.Context, vip types.HCVIPKey, hcConfig *types.HealthcheckConfig, reals []RealState) error {
	body := updateRequest{
		VIP:         vip,
		Healthcheck: hcConfig,
	}
	if reals != nil {
		entries := make([]realEntry, len(reals))
		for i, r := range reals {
			entries[i] = realEntry{Address: r.Address, Weight: r.Weight, Flags: r.Flags}
		}
		body.Reals = entries
	}
	return c.doRequest(ctx, http.MethodPut, "/api/v1/targets", body, nil)
}

// DeregisterVIP deregisters a VIP from the healthcheck service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//
// Returns an error if the deregistration fails.
func (c *HCClient) DeregisterVIP(ctx context.Context, vip types.HCVIPKey) error {
	return c.doRequest(ctx, http.MethodDelete, "/api/v1/targets", vip, nil)
}

// AddReals adds reals to a registered VIP on the healthcheck service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//   - reals: The list of reals to add.
//
// Returns an error if the request fails.
func (c *HCClient) AddReals(ctx context.Context, vip types.HCVIPKey, reals []RealState) error {
	entries := make([]realEntry, len(reals))
	for i, r := range reals {
		entries[i] = realEntry{Address: r.Address, Weight: r.Weight, Flags: r.Flags}
	}
	body := realsRequest{VIP: vip, Reals: entries}
	return c.doRequest(ctx, http.MethodPost, "/api/v1/targets/reals", body, nil)
}

// RemoveReals removes reals from a registered VIP on the healthcheck service.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//   - addresses: The list of real addresses to remove.
//
// Returns an error if the request fails.
func (c *HCClient) RemoveReals(ctx context.Context, vip types.HCVIPKey, addresses []string) error {
	entries := make([]realEntry, len(addresses))
	for i, addr := range addresses {
		entries[i] = realEntry{Address: addr}
	}
	body := realsRequest{VIP: vip, Reals: entries}
	return c.doRequest(ctx, http.MethodDelete, "/api/v1/targets/reals", body, nil)
}

// GetVIPHealthStatus retrieves the health status for a single VIP.
//
// Parameters:
//   - ctx: Context for the request.
//   - vip: The VIP key.
//
// Returns the health response or an error.
func (c *HCClient) GetVIPHealthStatus(ctx context.Context, vip types.HCVIPKey) (*types.HCVIPHealthResponse, error) {
	url := fmt.Sprintf("/api/v1/health/vip?address=%s&port=%d&proto=%d", vip.Address, vip.Port, vip.Proto)
	var result types.HCVIPHealthResponse
	if err := c.doRequest(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllHealth retrieves health states for all registered VIPs.
//
// Parameters:
//   - ctx: Context for the request.
//
// Returns a slice of VIP health responses or an error.
func (c *HCClient) GetAllHealth(ctx context.Context) ([]types.HCVIPHealthResponse, error) {
	var result []types.HCVIPHealthResponse
	if err := c.doRequest(ctx, http.MethodGet, "/api/v1/health", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// doRequest performs an HTTP request to the HC service.
func (c *HCClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
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
		return fmt.Errorf("HC service request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HC service response: %w", err)
	}

	var hcResp hcResponse
	if err := json.Unmarshal(respBody, &hcResp); err != nil {
		return fmt.Errorf("failed to parse HC service response: %w", err)
	}

	if !hcResp.Success {
		msg := "unknown error"
		if hcResp.Error != nil {
			msg = hcResp.Error.Message
		}
		return fmt.Errorf("HC service error: %s", msg)
	}

	if result != nil && hcResp.Data != nil {
		if err := json.Unmarshal(hcResp.Data, result); err != nil {
			return fmt.Errorf("failed to parse HC service data: %w", err)
		}
	}

	return nil
}
