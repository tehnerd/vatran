package hcservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// KatranClient is an HTTP client for registering/deregistering somark destinations
// with the katran server's healthcheck BPF program.
type KatranClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewKatranClient creates a new KatranClient.
//
// Parameters:
//   - baseURL: The base URL of the katran server (e.g., "http://localhost:8080").
//   - timeout: The HTTP client timeout in seconds.
//
// Returns a new KatranClient instance.
func NewKatranClient(baseURL string, timeout int) *KatranClient {
	return &KatranClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// katranDstRequest is the request body for somark destination operations.
type katranDstRequest struct {
	Somark uint32 `json:"somark"`
	Dst    string `json:"dst,omitempty"`
}

// katranResponse is the standard response wrapper from the katran server.
type katranResponse struct {
	Success bool            `json:"success"`
	Error   *katranRespErr  `json:"error,omitempty"`
}

// katranRespErr is the error field in a katran response.
type katranRespErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RegisterDst registers a somark-to-destination mapping with katran's BPF healthcheck program.
//
// Parameters:
//   - ctx: Context for the request.
//   - somark: The SO_MARK value.
//   - dst: The tunnel destination address for this real.
//
// Returns an error if the registration fails.
func (c *KatranClient) RegisterDst(ctx context.Context, somark uint32, dst string) error {
	body := katranDstRequest{Somark: somark, Dst: dst}
	return c.doRequest(ctx, http.MethodPost, "/api/v1/healthcheck/dsts", body)
}

// DeregisterDst removes a somark destination mapping from katran's BPF healthcheck program.
//
// Parameters:
//   - ctx: Context for the request.
//   - somark: The SO_MARK value to remove.
//
// Returns an error if the deregistration fails.
func (c *KatranClient) DeregisterDst(ctx context.Context, somark uint32) error {
	body := katranDstRequest{Somark: somark}
	return c.doRequest(ctx, http.MethodDelete, "/api/v1/healthcheck/dsts", body)
}

// doRequest performs an HTTP request to the katran server.
func (c *KatranClient) doRequest(ctx context.Context, method, path string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("katran request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read katran response: %w", err)
	}

	var katranResp katranResponse
	if err := json.Unmarshal(respBody, &katranResp); err != nil {
		return fmt.Errorf("failed to parse katran response: %w", err)
	}

	if !katranResp.Success {
		msg := "unknown error"
		if katranResp.Error != nil {
			msg = katranResp.Error.Message
		}
		return fmt.Errorf("katran error: %s", msg)
	}

	return nil
}
