// Package client provides an HTTP client for the Katran Load Balancer REST API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an HTTP client for the Katran REST API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option is a function that configures the Client.
type Option func(*Client)

// WithTimeout sets the HTTP client timeout.
//
// Parameters:
//   - timeout: The timeout duration for HTTP requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
//
// Parameters:
//   - httpClient: The custom HTTP client to use.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// New creates a new Katran API client.
//
// Parameters:
//   - baseURL: The base URL of the Katran server (e.g., "http://localhost:8080").
//   - opts: Optional configuration options.
//
// Returns:
//   - A new Client instance.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Response is the standard API response wrapper.
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *APIError       `json:"error,omitempty"`
}

// APIError represents an API error.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// doRequest performs an HTTP request and decodes the response.
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return apiResp.Error
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	if result != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to decode data: %w", err)
		}
	}

	return nil
}

// VIP represents a VIP configuration.
type VIP struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	Proto   uint8  `json:"proto"`
	Flags   uint32 `json:"flags,omitempty"`
}

// Real represents a real server.
type Real struct {
	Address string `json:"address"`
	Weight  uint32 `json:"weight"`
	Flags   uint8  `json:"flags,omitempty"`
}

// LBStats represents basic load balancer statistics.
type LBStats struct {
	V1 uint64 `json:"v1"`
	V2 uint64 `json:"v2"`
}

// HealthcheckerDst represents a healthcheck destination mapping.
type HealthcheckerDst struct {
	Somark uint32 `json:"somark"`
	Dst    string `json:"dst"`
}

// ListVIPs returns all configured VIPs.
//
// Returns:
//   - A slice of VIP configurations.
//   - An error if the request fails.
func (c *Client) ListVIPs() ([]VIP, error) {
	var vips []VIP
	if err := c.doRequest(http.MethodGet, "/api/v1/vips", nil, &vips); err != nil {
		return nil, err
	}
	return vips, nil
}

// AddVIP adds a new VIP.
//
// Parameters:
//   - vip: The VIP configuration to add.
//
// Returns:
//   - An error if the request fails.
func (c *Client) AddVIP(vip VIP) error {
	return c.doRequest(http.MethodPost, "/api/v1/vips", vip, nil)
}

// DeleteVIP removes a VIP.
//
// Parameters:
//   - address: The IP address of the VIP.
//   - port: The port number.
//   - proto: The IP protocol number.
//
// Returns:
//   - An error if the request fails.
func (c *Client) DeleteVIP(address string, port uint16, proto uint8) error {
	req := VIP{Address: address, Port: port, Proto: proto}
	return c.doRequest(http.MethodDelete, "/api/v1/vips", req, nil)
}

// GetVIPReals returns all real servers for a VIP.
//
// Parameters:
//   - address: The IP address of the VIP.
//   - port: The port number.
//   - proto: The IP protocol number.
//
// Returns:
//   - A slice of Real servers.
//   - An error if the request fails.
func (c *Client) GetVIPReals(address string, port uint16, proto uint8) ([]Real, error) {
	path := fmt.Sprintf("/api/v1/vips/reals?address=%s&port=%d&proto=%d",
		url.QueryEscape(address), port, proto)
	var reals []Real
	if err := c.doRequest(http.MethodGet, path, nil, &reals); err != nil {
		return nil, err
	}
	return reals, nil
}

// AddReal adds a real server to a VIP.
//
// Parameters:
//   - vipAddr: The IP address of the VIP.
//   - vipPort: The port number of the VIP.
//   - vipProto: The IP protocol number of the VIP.
//   - real: The real server to add.
//
// Returns:
//   - An error if the request fails.
func (c *Client) AddReal(vipAddr string, vipPort uint16, vipProto uint8, real Real) error {
	req := struct {
		VIP  VIP  `json:"vip"`
		Real Real `json:"real"`
	}{
		VIP:  VIP{Address: vipAddr, Port: vipPort, Proto: vipProto},
		Real: real,
	}
	return c.doRequest(http.MethodPost, "/api/v1/vips/reals", req, nil)
}

// DeleteReal removes a real server from a VIP.
//
// Parameters:
//   - vipAddr: The IP address of the VIP.
//   - vipPort: The port number of the VIP.
//   - vipProto: The IP protocol number of the VIP.
//   - realAddr: The IP address of the real server to remove.
//
// Returns:
//   - An error if the request fails.
func (c *Client) DeleteReal(vipAddr string, vipPort uint16, vipProto uint8, realAddr string) error {
	req := struct {
		VIP  VIP  `json:"vip"`
		Real Real `json:"real"`
	}{
		VIP:  VIP{Address: vipAddr, Port: vipPort, Proto: vipProto},
		Real: Real{Address: realAddr},
	}
	return c.doRequest(http.MethodDelete, "/api/v1/vips/reals", req, nil)
}

// UpdateReals batch updates real servers for a VIP.
//
// Parameters:
//   - vipAddr: The IP address of the VIP.
//   - vipPort: The port number of the VIP.
//   - vipProto: The IP protocol number of the VIP.
//   - action: 0 for add, 1 for delete.
//   - reals: The list of real servers to modify.
//
// Returns:
//   - An error if the request fails.
func (c *Client) UpdateReals(vipAddr string, vipPort uint16, vipProto uint8, action int, reals []Real) error {
	req := struct {
		VIP    VIP    `json:"vip"`
		Action int    `json:"action"`
		Reals  []Real `json:"reals"`
	}{
		VIP:    VIP{Address: vipAddr, Port: vipPort, Proto: vipProto},
		Action: action,
		Reals:  reals,
	}
	return c.doRequest(http.MethodPut, "/api/v1/vips/reals/batch", req, nil)
}

// GetVIPStats returns statistics for a VIP.
//
// Parameters:
//   - address: The IP address of the VIP.
//   - port: The port number.
//   - proto: The IP protocol number.
//
// Returns:
//   - LBStats containing packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetVIPStats(address string, port uint16, proto uint8) (*LBStats, error) {
	path := fmt.Sprintf("/api/v1/stats/vip?address=%s&port=%d&proto=%d",
		url.QueryEscape(address), port, proto)
	var stats LBStats
	if err := c.doRequest(http.MethodGet, path, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetLRUStats returns LRU cache statistics.
//
// Returns:
//   - LBStats containing LRU statistics.
//   - An error if the request fails.
func (c *Client) GetLRUStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/lru", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetLRUMissStats returns LRU miss statistics.
//
// Returns:
//   - LBStats containing LRU miss statistics.
//   - An error if the request fails.
func (c *Client) GetLRUMissStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/lru/miss", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetLRUFallbackStats returns LRU fallback statistics.
//
// Returns:
//   - LBStats containing LRU fallback statistics.
//   - An error if the request fails.
func (c *Client) GetLRUFallbackStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/lru/fallback", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetXDPTotalStats returns XDP total statistics.
//
// Returns:
//   - LBStats containing total packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetXDPTotalStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/xdp/total", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetXDPTxStats returns XDP TX statistics.
//
// Returns:
//   - LBStats containing TX packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetXDPTxStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/xdp/tx", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetXDPDropStats returns XDP drop statistics.
//
// Returns:
//   - LBStats containing drop packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetXDPDropStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/xdp/drop", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetXDPPassStats returns XDP pass statistics.
//
// Returns:
//   - LBStats containing pass packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetXDPPassStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/xdp/pass", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetDecapStats returns decapsulation statistics.
//
// Returns:
//   - LBStats containing decap packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetDecapStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/decap", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetInlineDecapStats returns inline decapsulation statistics.
//
// Returns:
//   - LBStats containing inline decap packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetInlineDecapStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/inline-decap", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetICMPTooBigStats returns ICMP too big statistics.
//
// Returns:
//   - LBStats containing ICMP too big packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetICMPTooBigStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/icmp-too-big", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetCHDropStats returns consistent hash drop statistics.
//
// Returns:
//   - LBStats containing CH drop packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetCHDropStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/ch-drop", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetSrcRoutingStats returns source routing statistics.
//
// Returns:
//   - LBStats containing source routing packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetSrcRoutingStats() (*LBStats, error) {
	var stats LBStats
	if err := c.doRequest(http.MethodGet, "/api/v1/stats/src-routing", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetMAC returns the current default router MAC address.
//
// Returns:
//   - The MAC address as a string.
//   - An error if the request fails.
func (c *Client) GetMAC() (string, error) {
	var resp struct {
		MAC string `json:"mac"`
	}
	if err := c.doRequest(http.MethodGet, "/api/v1/utils/mac", nil, &resp); err != nil {
		return "", err
	}
	return resp.MAC, nil
}

// SetMAC changes the default router MAC address.
//
// Parameters:
//   - mac: The new MAC address in hex format (e.g., "aa:bb:cc:dd:ee:ff").
//
// Returns:
//   - An error if the request fails.
func (c *Client) SetMAC(mac string) error {
	req := struct {
		MAC string `json:"mac"`
	}{MAC: mac}
	return c.doRequest(http.MethodPut, "/api/v1/utils/mac", req, nil)
}

// ListHealthcheckerDsts returns all healthcheck destination mappings.
//
// Returns:
//   - A slice of HealthcheckerDst mappings.
//   - An error if the request fails.
func (c *Client) ListHealthcheckerDsts() ([]HealthcheckerDst, error) {
	var dsts []HealthcheckerDst
	if err := c.doRequest(http.MethodGet, "/api/v1/healthcheck/dsts", nil, &dsts); err != nil {
		return nil, err
	}
	return dsts, nil
}

// AddHealthcheckerDst adds a healthcheck destination mapping.
//
// Parameters:
//   - somark: The socket mark value.
//   - dst: The destination IP address.
//
// Returns:
//   - An error if the request fails.
func (c *Client) AddHealthcheckerDst(somark uint32, dst string) error {
	req := HealthcheckerDst{Somark: somark, Dst: dst}
	return c.doRequest(http.MethodPost, "/api/v1/healthcheck/dsts", req, nil)
}

// DeleteHealthcheckerDst removes a healthcheck destination mapping.
//
// Parameters:
//   - somark: The socket mark value to remove.
//
// Returns:
//   - An error if the request fails.
func (c *Client) DeleteHealthcheckerDst(somark uint32) error {
	req := struct {
		Somark uint32 `json:"somark"`
	}{Somark: somark}
	return c.doRequest(http.MethodDelete, "/api/v1/healthcheck/dsts", req, nil)
}

// GetRealIndex returns the internal index for a real server.
//
// Parameters:
//   - address: The IP address of the real server.
//
// Returns:
//   - The internal index.
//   - An error if the request fails.
func (c *Client) GetRealIndex(address string) (int64, error) {
	path := fmt.Sprintf("/api/v1/reals/index?address=%s", url.QueryEscape(address))
	var resp struct {
		Index int64 `json:"index"`
	}
	if err := c.doRequest(http.MethodGet, path, nil, &resp); err != nil {
		return 0, err
	}
	return resp.Index, nil
}

// GetRealStats returns statistics for a real server by index.
//
// Parameters:
//   - index: The internal index of the real server.
//
// Returns:
//   - LBStats containing packet and byte counts.
//   - An error if the request fails.
func (c *Client) GetRealStats(index int64) (*LBStats, error) {
	path := fmt.Sprintf("/api/v1/stats/real?index=%d", index)
	var stats LBStats
	if err := c.doRequest(http.MethodGet, path, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// Health checks if the server is running.
//
// Returns:
//   - An error if the server is not healthy.
func (c *Client) Health() error {
	var resp struct {
		Status string `json:"status"`
	}
	return c.doRequest(http.MethodGet, "/health", nil, &resp)
}
