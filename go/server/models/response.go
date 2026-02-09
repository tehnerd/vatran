package models

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response wrapper.
type Response struct {
	// Success indicates whether the operation was successful.
	Success bool `json:"success"`
	// Data contains the response payload on success.
	Data interface{} `json:"data,omitempty"`
	// Error contains error information on failure.
	Error *APIError `json:"error,omitempty"`
}

// LBStatusResponse represents the load balancer status response.
type LBStatusResponse struct {
	// Initialized indicates if the LB instance has been created.
	Initialized bool `json:"initialized"`
	// Ready indicates if BPF programs are loaded and attached.
	Ready bool `json:"ready"`
}

// VIPResponse represents a VIP in API responses.
type VIPResponse struct {
	// Address is the IP address of the VIP.
	Address string `json:"address"`
	// Port is the port number.
	Port uint16 `json:"port"`
	// Proto is the IP protocol number.
	Proto uint8 `json:"proto"`
}

// VIPFlagsResponse represents VIP flags response.
type VIPFlagsResponse struct {
	// Flags are the current VIP flags.
	Flags uint32 `json:"flags"`
}

// RealResponse represents a real server in API responses.
type RealResponse struct {
	// Address is the IP address of the real server.
	Address string `json:"address"`
	// Weight is the weight for consistent hashing.
	Weight uint32 `json:"weight"`
	// Flags contains real-specific flags.
	Flags uint8 `json:"flags"`
	// Healthy indicates whether the real server is healthy and receiving traffic.
	Healthy bool `json:"healthy"`
}

// RealIndexResponse represents a real server index response.
type RealIndexResponse struct {
	// Index is the internal index of the real server.
	Index int64 `json:"index"`
}

// LBStatsResponse represents basic load balancer statistics.
type LBStatsResponse struct {
	// V1 is the first statistic value (typically packets).
	V1 uint64 `json:"v1"`
	// V2 is the second statistic value (typically bytes).
	V2 uint64 `json:"v2"`
}

// QuicPacketsStatsResponse contains QUIC packet routing statistics.
type QuicPacketsStatsResponse struct {
	// CHRouted is the count of packets routed via consistent hashing.
	CHRouted uint64 `json:"ch_routed"`
	// CIDInitial is the count of initial QUIC packets.
	CIDInitial uint64 `json:"cid_initial"`
	// CIDInvalidServerID is the count of packets with invalid server ID.
	CIDInvalidServerID uint64 `json:"cid_invalid_server_id"`
	// CIDInvalidServerIDSample is a sample of packets with invalid server ID.
	CIDInvalidServerIDSample uint64 `json:"cid_invalid_server_id_sample"`
	// CIDRouted is the count of packets routed via CID.
	CIDRouted uint64 `json:"cid_routed"`
	// CIDUnknownRealDropped is the count of packets dropped due to unknown real.
	CIDUnknownRealDropped uint64 `json:"cid_unknown_real_dropped"`
	// CIDV0 is the count of packets using CID version 0.
	CIDV0 uint64 `json:"cid_v0"`
	// CIDV1 is the count of packets using CID version 1.
	CIDV1 uint64 `json:"cid_v1"`
	// CIDV2 is the count of packets using CID version 2.
	CIDV2 uint64 `json:"cid_v2"`
	// CIDV3 is the count of packets using CID version 3.
	CIDV3 uint64 `json:"cid_v3"`
	// DstMatchInLRU is the count of packets with destination match in LRU.
	DstMatchInLRU uint64 `json:"dst_match_in_lru"`
	// DstMismatchInLRU is the count of packets with destination mismatch in LRU.
	DstMismatchInLRU uint64 `json:"dst_mismatch_in_lru"`
	// DstNotFoundInLRU is the count of packets with destination not found in LRU.
	DstNotFoundInLRU uint64 `json:"dst_not_found_in_lru"`
}

// TPRPacketsStatsResponse contains TCP Passive Routing statistics.
type TPRPacketsStatsResponse struct {
	// CHRouted is the count of packets routed via consistent hashing.
	CHRouted uint64 `json:"ch_routed"`
	// DstMismatchInLRU is the count of packets with destination mismatch in LRU.
	DstMismatchInLRU uint64 `json:"dst_mismatch_in_lru"`
	// SIDRouted is the count of packets routed via server ID.
	SIDRouted uint64 `json:"sid_routed"`
	// TCPSyn is the count of TCP SYN packets processed.
	TCPSyn uint64 `json:"tcp_syn"`
}

// HCStatsResponse contains healthcheck program statistics.
type HCStatsResponse struct {
	// PacketsProcessed is the total packets processed.
	PacketsProcessed uint64 `json:"packets_processed"`
	// PacketsDropped is the packets dropped.
	PacketsDropped uint64 `json:"packets_dropped"`
	// PacketsSkipped is the packets skipped.
	PacketsSkipped uint64 `json:"packets_skipped"`
	// PacketsTooBig is the packets exceeding maximum size.
	PacketsTooBig uint64 `json:"packets_too_big"`
}

// BPFMapStatsResponse contains BPF map statistics.
type BPFMapStatsResponse struct {
	// MaxEntries is the maximum number of entries the map can hold.
	MaxEntries uint32 `json:"max_entries"`
	// CurrentEntries is the current number of entries in the map.
	CurrentEntries uint32 `json:"current_entries"`
}

// UserspaceStatsResponse contains userspace library statistics.
type UserspaceStatsResponse struct {
	// BPFFailedCalls is the number of failed BPF syscalls.
	BPFFailedCalls uint64 `json:"bpf_failed_calls"`
	// AddrValidationFailed is the number of address validation failures.
	AddrValidationFailed uint64 `json:"addr_validation_failed"`
}

// MonitorStatsResponse contains monitoring subsystem statistics.
type MonitorStatsResponse struct {
	// Limit is the maximum number of packets to capture.
	Limit uint32 `json:"limit"`
	// Amount is the number of packets captured so far.
	Amount uint32 `json:"amount"`
	// BufferFull is the number of times the buffer was full.
	BufferFull uint32 `json:"buffer_full"`
}

// PerCorePacketsStatsResponse contains per-core packet statistics.
type PerCorePacketsStatsResponse struct {
	// Counts is a slice of packet counts, one per CPU core.
	Counts []int64 `json:"counts"`
}

// FloodStatusResponse contains flood status.
type FloodStatusResponse struct {
	// UnderFlood indicates if the system is under flood conditions.
	UnderFlood bool `json:"under_flood"`
}

// QuicRealResponse represents a QUIC real server mapping in API responses.
type QuicRealResponse struct {
	// Address is the IP address of the real server.
	Address string `json:"address"`
	// ID is the QUIC host ID.
	ID uint32 `json:"id"`
}

// SrcRoutingRuleResponse represents a source routing rule in API responses.
type SrcRoutingRuleResponse struct {
	// Src is the source IP prefix (CIDR notation).
	Src string `json:"src"`
	// Dst is the destination address.
	Dst string `json:"dst"`
}

// SrcRoutingRuleSizeResponse represents the source routing rule count.
type SrcRoutingRuleSizeResponse struct {
	// Size is the number of source routing rules.
	Size uint32 `json:"size"`
}

// HealthcheckerDstResponse represents a healthcheck destination mapping.
type HealthcheckerDstResponse struct {
	// Somark is the socket mark value.
	Somark uint32 `json:"somark"`
	// Dst is the destination address.
	Dst string `json:"dst"`
}

// FeatureStatusResponse represents feature status.
type FeatureStatusResponse struct {
	// Available indicates if the feature is available.
	Available bool `json:"available"`
}

// DeleteLRUResponse represents the result of deleting LRU entries.
type DeleteLRUResponse struct {
	// Maps is a list of map names where entries were deleted.
	Maps []string `json:"maps"`
}

// PurgeVIPLRUResponse represents the result of purging VIP LRU entries.
type PurgeVIPLRUResponse struct {
	// DeletedCount is the number of entries deleted.
	DeletedCount int `json:"deleted_count"`
}

// MACResponse represents a MAC address response.
type MACResponse struct {
	// MAC is the MAC address as hex string.
	MAC string `json:"mac"`
}

// RealForFlowResponse represents the real server for a flow.
type RealForFlowResponse struct {
	// Address is the real server IP address.
	Address string `json:"address"`
}

// SimulatePacketResponse represents the result of packet simulation.
type SimulatePacketResponse struct {
	// Packet is the output packet data as base64-encoded string.
	Packet string `json:"packet"`
}

// ProgFDResponse represents a BPF program file descriptor response.
type ProgFDResponse struct {
	// FD is the file descriptor.
	FD int `json:"fd"`
}

// MapFDsResponse represents a list of BPF map file descriptors.
type MapFDsResponse struct {
	// FDs is the list of file descriptors.
	FDs []int `json:"fds"`
}

// WriteJSON writes a JSON response to the http.ResponseWriter.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - status: The HTTP status code.
//   - data: The data to encode as JSON.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success response.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - data: The response data.
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// WriteCreated writes a created response.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - data: The response data.
func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// WriteError writes an error response.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - status: The HTTP status code.
//   - err: The API error.
func WriteError(w http.ResponseWriter, status int, err *APIError) {
	WriteJSON(w, status, Response{
		Success: false,
		Error:   err,
	})
}

// WriteKatranError writes an error response from a katran error.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - err: The error to map and write.
func WriteKatranError(w http.ResponseWriter, err error) {
	status, apiErr := MapKatranError(err)
	WriteError(w, status, apiErr)
}
