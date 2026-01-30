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

// Error represents a Katran error code.
type Error int

// Error codes returned by Katran API functions.
const (
	// OK indicates the operation completed successfully.
	OK Error = 0
	// ErrInvalidArgument indicates an invalid argument was passed.
	ErrInvalidArgument Error = -1
	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound Error = -2
	// ErrAlreadyExists indicates the resource already exists.
	ErrAlreadyExists Error = -3
	// ErrSpaceExhausted indicates maximum capacity was reached.
	ErrSpaceExhausted Error = -4
	// ErrBPFFailed indicates a BPF operation failed.
	ErrBPFFailed Error = -5
	// ErrFeatureDisabled indicates the requested feature is not enabled.
	ErrFeatureDisabled Error = -6
	// ErrInternal indicates an internal error occurred.
	ErrInternal Error = -7
	// ErrMemory indicates memory allocation failed.
	ErrMemory Error = -8
)

// ModifyAction specifies the action to perform when modifying VIP-real associations.
type ModifyAction int

const (
	// ActionAdd adds a real to a VIP.
	ActionAdd ModifyAction = 0
	// ActionDel removes a real from a VIP.
	ActionDel ModifyAction = 1
)

// HashFunction specifies the hash algorithm for consistent hashing.
type HashFunction int

const (
	// HashMaglev uses the original Maglev consistent hashing algorithm.
	HashMaglev HashFunction = 0
	// HashMaglevV2 uses the improved Maglev V2 algorithm.
	HashMaglevV2 HashFunction = 1
)

// Feature represents optional Katran capabilities.
type Feature int

const (
	// FeatureSrcRouting enables source-based routing support.
	FeatureSrcRouting Feature = 1 << 0
	// FeatureInlineDecap enables inline packet decapsulation.
	FeatureInlineDecap Feature = 1 << 1
	// FeatureIntrospection enables packet introspection/monitoring.
	FeatureIntrospection Feature = 1 << 2
	// FeatureGUEEncap enables GUE (Generic UDP Encapsulation) instead of IPIP.
	FeatureGUEEncap Feature = 1 << 3
	// FeatureDirectHC enables direct healthcheck encapsulation.
	FeatureDirectHC Feature = 1 << 4
	// FeatureLocalDeliveryOpt enables local delivery optimization.
	FeatureLocalDeliveryOpt Feature = 1 << 5
	// FeatureFlowDebug enables flow debugging maps.
	FeatureFlowDebug Feature = 1 << 6
)

// VIPKey uniquely identifies a Virtual IP (VIP).
type VIPKey struct {
	// Address is the IP address of the VIP (IPv4 or IPv6).
	Address string
	// Port is the port number in host byte order.
	Port uint16
	// Proto is the IP protocol number (e.g., 6 for TCP, 17 for UDP).
	Proto uint8
}

// Real represents a real server (backend) definition.
type Real struct {
	// Address is the IP address of the real server (IPv4 or IPv6).
	Address string
	// Weight is the weight for consistent hashing (higher = more traffic).
	Weight uint32
	// Flags contains real-specific flags.
	Flags uint8
}

// QuicReal represents a QUIC connection ID to real server mapping.
type QuicReal struct {
	// Address is the IP address of the real server.
	Address string
	// ID is the QUIC host ID portion embedded in connection IDs.
	ID uint32
}

// Flow represents a 5-tuple flow identifier.
type Flow struct {
	// Src is the source IP address (IPv4 or IPv6).
	Src string
	// Dst is the destination IP address (IPv4 or IPv6).
	Dst string
	// SrcPort is the source port number in host byte order.
	SrcPort uint16
	// DstPort is the destination port number in host byte order.
	DstPort uint16
	// Proto is the IP protocol number.
	Proto uint8
}

// LBStats contains generic statistics counters.
type LBStats struct {
	// V1 is the first statistic value (typically packets).
	V1 uint64
	// V2 is the second statistic value (typically bytes).
	V2 uint64
}

// QuicPacketsStats contains QUIC packet routing statistics.
type QuicPacketsStats struct {
	// CHRouted is the count of packets routed via consistent hashing.
	CHRouted uint64
	// CIDInitial is the count of initial QUIC packets.
	CIDInitial uint64
	// CIDInvalidServerID is the count of packets with invalid server ID.
	CIDInvalidServerID uint64
	// CIDInvalidServerIDSample is a sample of packets with invalid server ID.
	CIDInvalidServerIDSample uint64
	// CIDRouted is the count of packets routed via CID.
	CIDRouted uint64
	// CIDUnknownRealDropped is the count of packets dropped due to unknown real.
	CIDUnknownRealDropped uint64
	// CIDV0 is the count of packets using CID version 0.
	CIDV0 uint64
	// CIDV1 is the count of packets using CID version 1.
	CIDV1 uint64
	// CIDV2 is the count of packets using CID version 2.
	CIDV2 uint64
	// CIDV3 is the count of packets using CID version 3.
	CIDV3 uint64
	// DstMatchInLRU is the count of packets with destination match in LRU.
	DstMatchInLRU uint64
	// DstMismatchInLRU is the count of packets with destination mismatch in LRU.
	DstMismatchInLRU uint64
	// DstNotFoundInLRU is the count of packets with destination not found in LRU.
	DstNotFoundInLRU uint64
}

// TPRPacketsStats contains TCP Passive Routing statistics.
type TPRPacketsStats struct {
	// CHRouted is the count of packets routed via consistent hashing.
	CHRouted uint64
	// DstMismatchInLRU is the count of packets with destination mismatch in LRU.
	DstMismatchInLRU uint64
	// SIDRouted is the count of packets routed via server ID.
	SIDRouted uint64
	// TCPSyn is the count of TCP SYN packets processed.
	TCPSyn uint64
}

// HCStats contains healthcheck program statistics.
type HCStats struct {
	// PacketsProcessed is the total packets processed.
	PacketsProcessed uint64
	// PacketsDropped is the packets dropped.
	PacketsDropped uint64
	// PacketsSkipped is the packets skipped.
	PacketsSkipped uint64
	// PacketsTooBig is the packets exceeding maximum size.
	PacketsTooBig uint64
}

// BPFMapStats contains BPF map statistics.
type BPFMapStats struct {
	// MaxEntries is the maximum number of entries the map can hold.
	MaxEntries uint32
	// CurrentEntries is the current number of entries in the map.
	CurrentEntries uint32
}

// MonitorStats contains monitoring subsystem statistics.
type MonitorStats struct {
	// Limit is the maximum number of packets to capture.
	Limit uint32
	// Amount is the number of packets captured so far.
	Amount uint32
	// BufferFull is the number of times the buffer was full.
	BufferFull uint32
}

// UserspaceStats contains userspace library statistics.
type UserspaceStats struct {
	// BPFFailedCalls is the number of failed BPF syscalls.
	BPFFailedCalls uint64
	// AddrValidationFailed is the number of address validation failures.
	AddrValidationFailed uint64
}
