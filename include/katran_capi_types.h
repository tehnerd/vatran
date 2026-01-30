/**
 * @file katran_capi_types.h
 * @brief C type definitions for the Katran load balancer C API.
 *
 * This header defines all C-compatible structures and enumerations used
 * by the Katran C API. These types map to their C++ counterparts in the
 * katran namespace but are designed for use from C code (including CGO).
 */

#ifndef KATRAN_CAPI_TYPES_H
#define KATRAN_CAPI_TYPES_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a KatranLb instance.
 *
 * This handle encapsulates all internal state of a Katran load balancer
 * instance. Users should treat this as an opaque pointer and only interact
 * with it through the C API functions.
 *
 * Lifecycle:
 *   1. Create with katran_lb_create()
 *   2. Use with various API functions
 *   3. Destroy with katran_lb_destroy()
 */
typedef struct katran_lb_handle* katran_lb_t;

/**
 * @brief Error codes returned by Katran C API functions.
 *
 * All API functions return one of these error codes. Success is indicated
 * by KATRAN_OK (0). All error codes are negative values.
 */
typedef enum {
  /** Operation completed successfully */
  KATRAN_OK = 0,
  /** Invalid argument passed to function (NULL pointer, invalid value, etc.) */
  KATRAN_ERR_INVALID_ARGUMENT = -1,
  /** Requested resource (VIP, real, etc.) was not found */
  KATRAN_ERR_NOT_FOUND = -2,
  /** Resource already exists (duplicate VIP, real, etc.) */
  KATRAN_ERR_ALREADY_EXISTS = -3,
  /** Maximum capacity reached (too many VIPs, reals, etc.) */
  KATRAN_ERR_SPACE_EXHAUSTED = -4,
  /** BPF operation failed (map update, program load, etc.) */
  KATRAN_ERR_BPF_FAILED = -5,
  /** Requested feature is not enabled in the current configuration */
  KATRAN_ERR_FEATURE_DISABLED = -6,
  /** Internal error (unexpected exception, state corruption, etc.) */
  KATRAN_ERR_INTERNAL = -7,
  /** Memory allocation failed */
  KATRAN_ERR_MEMORY = -8
} katran_error_t;

/**
 * @brief Action to perform when modifying VIP-real associations.
 */
typedef enum {
  /** Add a real to a VIP */
  KATRAN_ACTION_ADD = 0,
  /** Remove a real from a VIP */
  KATRAN_ACTION_DEL = 1
} katran_modify_action_t;

/**
 * @brief Hash function algorithm for consistent hashing ring generation.
 */
typedef enum {
  /** Original Maglev consistent hashing algorithm */
  KATRAN_HASH_MAGLEV = 0,
  /** Improved Maglev V2 algorithm with better distribution */
  KATRAN_HASH_MAGLEV_V2 = 1
} katran_hash_function_t;

/**
 * @brief Feature flags indicating optional Katran capabilities.
 *
 * These flags can be used with katran_lb_has_feature(),
 * katran_lb_install_feature(), and katran_lb_remove_feature() to query
 * and manage optional features at runtime.
 */
typedef enum {
  /** Source-based routing support */
  KATRAN_FEATURE_SRC_ROUTING = 1 << 0,
  /** Inline packet decapsulation */
  KATRAN_FEATURE_INLINE_DECAP = 1 << 1,
  /** Packet introspection/monitoring */
  KATRAN_FEATURE_INTROSPECTION = 1 << 2,
  /** GUE (Generic UDP Encapsulation) instead of IPIP */
  KATRAN_FEATURE_GUE_ENCAP = 1 << 3,
  /** Direct healthcheck encapsulation */
  KATRAN_FEATURE_DIRECT_HC = 1 << 4,
  /** Local delivery optimization (XDP_PASS for local traffic) */
  KATRAN_FEATURE_LOCAL_DELIVERY_OPT = 1 << 5,
  /** Flow debugging maps enabled */
  KATRAN_FEATURE_FLOW_DEBUG = 1 << 6
} katran_feature_t;

/**
 * @brief Virtual IP (VIP) identifier.
 *
 * A VIP is uniquely identified by the combination of its IP address,
 * port number, and protocol. This struct is used to specify which VIP
 * to operate on in most API functions.
 */
typedef struct {
  /**
   * IP address of the VIP as a string.
   * Supports both IPv4 (e.g., "10.0.0.1") and IPv6 (e.g., "fc00::1").
   * The string must remain valid for the duration of the API call.
   */
  const char* address;

  /**
   * Port number in host byte order.
   * Use 0 for port-independent VIPs (requires appropriate VIP flags).
   */
  uint16_t port;

  /**
   * IP protocol number (e.g., 6 for TCP, 17 for UDP).
   */
  uint8_t proto;
} katran_vip_key_t;

/**
 * @brief Real server (backend) definition for VIP association.
 *
 * Represents a backend server that can receive traffic for a VIP.
 * Used when adding or modifying reals for a VIP.
 */
typedef struct {
  /**
   * IP address of the real server as a string.
   * Supports both IPv4 and IPv6 addresses.
   * The string must remain valid for the duration of the API call.
   */
  const char* address;

  /**
   * Weight for consistent hashing (higher = more traffic).
   * A weight of 0 effectively removes the real from the hash ring.
   */
  uint32_t weight;

  /**
   * Real-specific flags (see Katran real flag constants).
   */
  uint8_t flags;
} katran_new_real_t;

/**
 * @brief QUIC connection ID to real server mapping.
 *
 * Used to configure direct routing based on QUIC connection IDs,
 * allowing stateful QUIC connections to be routed to specific backends.
 */
typedef struct {
  /**
   * IP address of the real server as a string.
   * The string must remain valid for the duration of the API call.
   */
  const char* address;

  /**
   * QUIC connection ID (server-generated portion).
   * This is the host ID portion embedded in QUIC connection IDs.
   */
  uint32_t id;
} katran_quic_real_t;

/**
 * @brief 5-tuple flow identifier for packet simulation.
 *
 * Represents a network flow using the standard 5-tuple: source IP,
 * destination IP, source port, destination port, and protocol.
 */
typedef struct {
  /**
   * Source IP address as a string (IPv4 or IPv6).
   * The string must remain valid for the duration of the API call.
   */
  const char* src;

  /**
   * Destination IP address as a string (IPv4 or IPv6).
   * This should typically be a configured VIP address.
   */
  const char* dst;

  /**
   * Source port number in host byte order.
   */
  uint16_t src_port;

  /**
   * Destination port number in host byte order.
   */
  uint16_t dst_port;

  /**
   * IP protocol number (e.g., 6 for TCP, 17 for UDP).
   */
  uint8_t proto;
} katran_flow_t;

/**
 * @brief Generic statistics counters.
 *
 * Used to return various statistics from the load balancer.
 * The interpretation of v1 and v2 depends on the specific stats function.
 */
typedef struct {
  /**
   * First statistic value.
   * Typically represents packets or primary counter.
   */
  uint64_t v1;

  /**
   * Second statistic value.
   * Typically represents bytes or secondary counter.
   */
  uint64_t v2;
} katran_lb_stats_t;

/**
 * @brief QUIC packet processing statistics.
 *
 * Detailed statistics about QUIC packet routing decisions.
 */
typedef struct {
  /** Packets routed via consistent hashing */
  uint64_t ch_routed;
  /** Initial QUIC packets (no CID routing possible) */
  uint64_t cid_initial;
  /** Packets with invalid server ID in CID */
  uint64_t cid_invalid_server_id;
  /** Sample of packets with invalid server ID */
  uint64_t cid_invalid_server_id_sample;
  /** Packets successfully routed via CID */
  uint64_t cid_routed;
  /** Packets dropped due to unknown real for CID */
  uint64_t cid_unknown_real_dropped;
  /** Packets using CID version 0 */
  uint64_t cid_v0;
  /** Packets using CID version 1 */
  uint64_t cid_v1;
  /** Packets using CID version 2 */
  uint64_t cid_v2;
  /** Packets using CID version 3 */
  uint64_t cid_v3;
  /** Packets with destination match in LRU */
  uint64_t dst_match_in_lru;
  /** Packets with destination mismatch in LRU */
  uint64_t dst_mismatch_in_lru;
  /** Packets with destination not found in LRU */
  uint64_t dst_not_found_in_lru;
} katran_quic_packets_stats_t;

/**
 * @brief TCP server ID routing statistics (TPR).
 *
 * Statistics for TCP Passive Routing based on server IDs.
 */
typedef struct {
  /** Packets routed via consistent hashing */
  uint64_t ch_routed;
  /** Packets with destination mismatch in LRU */
  uint64_t dst_mismatch_in_lru;
  /** Packets routed via server ID */
  uint64_t sid_routed;
  /** TCP SYN packets processed */
  uint64_t tcp_syn;
} katran_tpr_packets_stats_t;

/**
 * @brief Healthcheck program statistics.
 *
 * Packet counters for the healthcheck encapsulation program.
 */
typedef struct {
  /** Total packets processed */
  uint64_t packets_processed;
  /** Packets dropped */
  uint64_t packets_dropped;
  /** Packets skipped (no action taken) */
  uint64_t packets_skipped;
  /** Packets exceeding maximum size */
  uint64_t packets_too_big;
} katran_hc_stats_t;

/**
 * @brief BPF map statistics.
 *
 * Information about BPF map capacity and current usage.
 */
typedef struct {
  /** Maximum number of entries the map can hold */
  uint32_t max_entries;
  /** Current number of entries in the map */
  uint32_t current_entries;
} katran_bpf_map_stats_t;

/**
 * @brief Katran monitor (introspection) statistics.
 *
 * Statistics from the packet capture/monitoring subsystem.
 */
typedef struct {
  /** Maximum number of packets to capture */
  uint32_t limit;
  /** Number of packets captured so far */
  uint32_t amount;
  /** Number of times the buffer was full */
  uint32_t buffer_full;
} katran_monitor_stats_t;

/**
 * @brief Userspace library statistics.
 *
 * Statistics about the userspace component of Katran.
 */
typedef struct {
  /** Number of failed BPF syscalls */
  uint64_t bpf_failed_calls;
  /** Number of address validation failures */
  uint64_t addr_validation_failed;
} katran_userspace_stats_t;

/**
 * @brief Katran load balancer configuration.
 *
 * Complete configuration for initializing a KatranLb instance.
 * All string pointers must remain valid only during the katran_lb_create() call.
 */
typedef struct {
  /* ==================== Interface Configuration ==================== */

  /**
   * Name of the main network interface to attach XDP program.
   * Example: "eth0", "ens192".
   * Required field.
   */
  const char* main_interface;

  /**
   * Name of the IPv4 tunnel interface for healthcheck encapsulation.
   * Used when tunnel_based_hc_encap is enabled.
   * Can be NULL if not using tunnel-based healthchecks.
   */
  const char* v4_tun_interface;

  /**
   * Name of the IPv6 tunnel interface for healthcheck encapsulation.
   * Used when tunnel_based_hc_encap is enabled.
   * Can be NULL if not using tunnel-based healthchecks.
   */
  const char* v6_tun_interface;

  /**
   * Interface for attaching healthcheck BPF program.
   * If NULL or empty, main_interface is used.
   */
  const char* hc_interface;

  /* ==================== BPF Program Paths ==================== */

  /**
   * Path to the compiled balancer BPF program (.o file).
   * Required field.
   */
  const char* balancer_prog_path;

  /**
   * Path to the compiled healthcheck BPF program (.o file).
   * Required when enable_hc is true.
   */
  const char* healthchecking_prog_path;

  /* ==================== MAC Addresses ==================== */

  /**
   * MAC address of the default gateway/router (6 bytes).
   * Packets will be L2-forwarded to this address.
   * If NULL, must be set later via katran_lb_change_mac().
   */
  const uint8_t* default_mac;

  /**
   * MAC address of the local server (6 bytes).
   * Used for packet source address when needed.
   * Can be NULL if not required.
   */
  const uint8_t* local_mac;

  /* ==================== Root Map Configuration ==================== */

  /**
   * Path to pinned BPF map from root XDP program.
   * Used in "shared" mode where a root program tail-calls into Katran.
   * Set to NULL or empty string for "standalone" mode.
   */
  const char* root_map_path;

  /**
   * Position (index) in the root map for Katran's program FD.
   * Only used when root_map_path is set.
   * Default: 2
   */
  uint32_t root_map_pos;

  /**
   * Whether to use the root map (1) or standalone mode (0).
   * Default: 1 (true)
   */
  int use_root_map;

  /* ==================== Capacity Limits ==================== */

  /**
   * Maximum number of VIPs that can be configured.
   * Default: 512
   */
  uint32_t max_vips;

  /**
   * Maximum number of real servers that can be configured.
   * Default: 4096
   */
  uint32_t max_reals;

  /**
   * Size of the consistent hashing ring per VIP.
   * Larger values provide better distribution but use more memory.
   * Default: 65537
   */
  uint32_t ch_ring_size;

  /**
   * Size of the per-CPU LRU connection tracking table.
   * Larger values support more concurrent connections.
   * Default: 8000000
   */
  uint64_t lru_size;

  /**
   * Maximum number of source routing LPM entries.
   * Default: 3000000
   */
  uint32_t max_lpm_src_size;

  /**
   * Maximum number of inline decapsulation destinations.
   * Default: 6
   */
  uint32_t max_decap_dst;

  /**
   * Size of per-CPU global LRU maps.
   * Default: 100000
   */
  uint32_t global_lru_size;

  /* ==================== Feature Flags ==================== */

  /**
   * Enable healthcheck encapsulation program (1) or not (0).
   * Default: 1 (true)
   */
  int enable_hc;

  /**
   * Use tunnel interfaces for healthcheck encapsulation (1) or direct (0).
   * Default: 1 (true)
   */
  int tunnel_based_hc_encap;

  /**
   * Testing mode - don't actually program the forwarding plane (1).
   * Useful for unit tests.
   * Default: 0 (false)
   */
  int testing;

  /**
   * Set RLIMIT_MEMLOCK to unlimited (1) or leave unchanged (0).
   * Required for loading BPF programs on some systems.
   * Default: 1 (true)
   */
  int memlock_unlimited;

  /**
   * Enable flow debugging maps (1) or not (0).
   * Provides additional debugging information but uses more memory.
   * Default: 0 (false)
   */
  int flow_debug;

  /**
   * Enable QUIC CID version 3 support (1) or not (0).
   * Default: 0 (false)
   */
  int enable_cid_v3;

  /**
   * Clean up BPF resources on shutdown (1) or leave attached (0).
   * Default: 1 (true)
   */
  int cleanup_on_shutdown;

  /* ==================== CPU/NUMA Configuration ==================== */

  /**
   * Array of CPU core IDs responsible for packet forwarding.
   * These cores will have dedicated per-CPU LRU maps.
   * Can be NULL if using automatic core detection.
   */
  const int32_t* forwarding_cores;

  /**
   * Number of elements in forwarding_cores array.
   */
  size_t forwarding_cores_count;

  /**
   * Array mapping forwarding cores to NUMA nodes.
   * Must have same length as forwarding_cores, or be NULL.
   * Enables NUMA-aware memory allocation for LRU maps.
   */
  const int32_t* numa_nodes;

  /**
   * Number of elements in numa_nodes array.
   */
  size_t numa_nodes_count;

  /* ==================== XDP Configuration ==================== */

  /**
   * XDP attach flags (e.g., XDP_FLAGS_SKB_MODE, XDP_FLAGS_DRV_MODE).
   * Default: 0 (native mode if supported)
   */
  uint32_t xdp_attach_flags;

  /**
   * TC priority for healthcheck program attachment.
   * Default: 2307
   */
  uint32_t priority;

  /**
   * Interface index for main_interface.
   * If 0, will be resolved from main_interface name.
   */
  uint32_t main_interface_index;

  /**
   * Interface index for hc_interface.
   * If 0, will be resolved from hc_interface name.
   */
  uint32_t hc_interface_index;

  /* ==================== GUE Source Addresses ==================== */

  /**
   * IPv4 source address for GUE-encapsulated packets.
   * Required when using GUE encapsulation.
   * Example: "10.0.0.1"
   */
  const char* katran_src_v4;

  /**
   * IPv6 source address for GUE-encapsulated packets.
   * Required when using GUE encapsulation with IPv6.
   * Example: "fc00::1"
   */
  const char* katran_src_v6;

  /* ==================== Hash Function ==================== */

  /**
   * Hash function algorithm for consistent hashing.
   * Default: KATRAN_HASH_MAGLEV
   */
  katran_hash_function_t hash_function;
} katran_config_t;

#ifdef __cplusplus
}
#endif

#endif /* KATRAN_CAPI_TYPES_H */
