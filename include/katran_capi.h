/* Copyright (C) 2018-present, Facebook, Inc.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

/**
 * @file katran_capi.h
 * @brief C API for the Katran XDP-based L4 Load Balancer.
 *
 * This header provides a C-compatible interface to the Katran load balancer,
 * enabling integration with languages like Go (via CGO), Rust (via FFI),
 * and other languages that can call C functions.
 *
 * ## Usage Pattern
 *
 * 1. Initialize configuration with katran_config_init()
 * 2. Create instance with katran_lb_create()
 * 3. Load BPF programs with katran_lb_load_bpf_progs()
 * 4. Attach to interfaces with katran_lb_attach_bpf_progs()
 * 5. Configure VIPs and reals using the management functions
 * 6. Clean up with katran_lb_destroy()
 *
 * ## Error Handling
 *
 * All functions return katran_error_t. On error, call katran_lb_get_last_error()
 * to retrieve a human-readable error message (thread-local).
 *
 * ## Memory Management
 *
 * - Input strings: Borrowed (caller owns, must remain valid during call)
 * - Output arrays: Allocated by API, freed via katran_free_*() functions
 * - For array outputs: Pass NULL to get count, then allocate and call again
 *
 * ## Thread Safety
 *
 * The underlying KatranLb class is NOT thread-safe. Callers must provide
 * external synchronization if accessing from multiple threads.
 */

#ifndef KATRAN_CAPI_H
#define KATRAN_CAPI_H

#include "katran_capi_types.h"

#ifdef __cplusplus
extern "C" {
#endif

/* ============================================================================
 * INITIALIZATION AND LIFECYCLE
 * ============================================================================ */

/**
 * @brief Initialize a configuration structure with default values.
 *
 * This function sets all fields of the configuration structure to sensible
 * defaults. Callers should call this before setting custom values to ensure
 * all fields are properly initialized.
 *
 * @param[out] config Pointer to configuration structure to initialize.
 *                    Must not be NULL.
 *
 * @return KATRAN_OK on success, KATRAN_ERR_INVALID_ARGUMENT if config is NULL.
 */
katran_error_t katran_config_init(katran_config_t* config);

/**
 * @brief Create a new Katran load balancer instance.
 *
 * Allocates and initializes a new KatranLb instance with the provided
 * configuration. The returned handle must be destroyed with
 * katran_lb_destroy() when no longer needed.
 *
 * After creation, call katran_lb_load_bpf_progs() and
 * katran_lb_attach_bpf_progs() to start load balancing.
 *
 * @param[in]  config Pointer to configuration structure. Must not be NULL.
 *                    String fields must remain valid only during this call.
 * @param[out] handle Pointer to receive the created handle. Must not be NULL.
 *                    On success, *handle will be set to a valid handle.
 *                    On failure, *handle will be set to NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if config or handle is NULL, or if
 *         required configuration fields are missing/invalid.
 * @return KATRAN_ERR_MEMORY if memory allocation fails.
 * @return KATRAN_ERR_INTERNAL for other initialization failures.
 */
katran_error_t katran_lb_create(
    const katran_config_t* config,
    katran_lb_t* handle);

/**
 * @brief Destroy a Katran load balancer instance.
 *
 * Releases all resources associated with the handle, including BPF programs
 * if cleanup_on_shutdown was set in the configuration.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param[in] handle Handle to destroy. May be NULL (no-op).
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INTERNAL if cleanup fails.
 */
katran_error_t katran_lb_destroy(katran_lb_t handle);

/**
 * @brief Load BPF programs into the kernel.
 *
 * Loads the balancer and (optionally) healthcheck BPF programs specified
 * in the configuration. This must be called before katran_lb_attach_bpf_progs().
 *
 * @param[in] handle Valid Katran handle from katran_lb_create().
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_BPF_FAILED if BPF program loading fails.
 */
katran_error_t katran_lb_load_bpf_progs(katran_lb_t handle);

/**
 * @brief Attach loaded BPF programs to network interfaces.
 *
 * Attaches the balancer program to the main interface (XDP) and optionally
 * the healthcheck program to the HC interface (TC). Must be called after
 * katran_lb_load_bpf_progs().
 *
 * @param[in] handle Valid Katran handle with loaded BPF programs.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_BPF_FAILED if attachment fails.
 */
katran_error_t katran_lb_attach_bpf_progs(katran_lb_t handle);

/**
 * @brief Reload the balancer BPF program at runtime.
 *
 * Allows hot-reloading of the balancer program without service interruption.
 * The new program takes effect immediately after successful reload.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] path   Path to the new BPF program file. Must not be NULL.
 * @param[in] config Optional new configuration. Pass NULL to keep current config.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or path is NULL.
 * @return KATRAN_ERR_BPF_FAILED if reload fails.
 */
katran_error_t katran_lb_reload_balancer_prog(
    katran_lb_t handle,
    const char* path,
    const katran_config_t* config);

/* ============================================================================
 * MAC ADDRESS MANAGEMENT
 * ============================================================================ */

/**
 * @brief Change the default router MAC address.
 *
 * Updates the MAC address used as the destination for forwarded packets.
 * This is typically the MAC of the default gateway or next-hop router.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] mac    Pointer to 6-byte MAC address. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or mac is NULL.
 * @return KATRAN_ERR_BPF_FAILED if updating the BPF map fails.
 */
katran_error_t katran_lb_change_mac(katran_lb_t handle, const uint8_t* mac);

/**
 * @brief Get the current default router MAC address.
 *
 * Retrieves the MAC address currently configured for packet forwarding.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] mac    Pointer to 6-byte buffer to receive MAC address.
 *                    Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or mac is NULL.
 */
katran_error_t katran_lb_get_mac(katran_lb_t handle, uint8_t* mac);

/* ============================================================================
 * VIP MANAGEMENT
 * ============================================================================ */

/**
 * @brief Add a new Virtual IP (VIP) to the load balancer.
 *
 * Creates a new VIP that can receive traffic. After adding a VIP, use
 * katran_lb_add_real_for_vip() or katran_lb_modify_reals_for_vip() to
 * configure backend servers.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] vip    VIP to add (address, port, protocol). Must not be NULL.
 *                   The address string must be a valid IPv4 or IPv6 address.
 * @param[in] flags  VIP flags (e.g., NO_PORT, NO_LRU, etc.). Pass 0 for defaults.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or vip is NULL, or if the
 *         address cannot be parsed.
 * @return KATRAN_ERR_ALREADY_EXISTS if the VIP already exists.
 * @return KATRAN_ERR_SPACE_EXHAUSTED if max_vips limit is reached.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t flags);

/**
 * @brief Delete a VIP from the load balancer.
 *
 * Removes the VIP and all associated real server configurations.
 * Traffic to this VIP will no longer be load balanced.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] vip    VIP to delete. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or vip is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_vip(katran_lb_t handle, const katran_vip_key_t* vip);

/**
 * @brief Get all configured VIPs.
 *
 * Retrieves a list of all VIPs currently configured in the load balancer.
 *
 * Usage pattern (two-call):
 *   1. Call with vips=NULL to get count
 *   2. Allocate array of katran_vip_key_t[count]
 *   3. Call again with allocated array
 *   4. Free with katran_free_vips() when done
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    vips   Array to receive VIPs, or NULL to query count.
 * @param[in,out] count  On input with non-NULL vips: size of vips array.
 *                       On output: number of VIPs returned/available.
 *                       Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_all_vips(
    katran_lb_t handle,
    katran_vip_key_t* vips,
    size_t* count);

/**
 * @brief Modify a VIP's flags.
 *
 * Sets or clears specific flags on an existing VIP. Flags control
 * forwarding behavior (e.g., bypass LRU, ignore source port in hash).
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] vip    VIP to modify. Must not be NULL.
 * @param[in] flag   Flag bits to modify.
 * @param[in] set    1 to set the flags, 0 to clear them.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or vip is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_modify_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t flag,
    int set);

/**
 * @brief Get a VIP's current flags.
 *
 * Retrieves the flag bits currently set on a VIP.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[in]  vip    VIP to query. Must not be NULL.
 * @param[out] flags  Pointer to receive the flag value. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle, vip, or flags is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 */
katran_error_t katran_lb_get_vip_flags(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t* flags);

/**
 * @brief Change the hash function for a VIP's consistent hashing ring.
 *
 * Updates the algorithm used to generate the consistent hash ring.
 * This triggers a recalculation of the entire ring.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] vip    VIP to modify. Must not be NULL.
 * @param[in] func   Hash function to use.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or vip is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_change_hash_function_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_hash_function_t func);

/* ============================================================================
 * REAL SERVER MANAGEMENT
 * ============================================================================ */

/**
 * @brief Add a real server to a VIP.
 *
 * Adds a backend server to the VIP's consistent hash ring with the
 * specified weight. Higher weights result in more traffic.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] real   Real server to add (address, weight, flags). Must not be NULL.
 * @param[in] vip    VIP to add the real to. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL or address is invalid.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_SPACE_EXHAUSTED if max_reals limit is reached.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_real_for_vip(
    katran_lb_t handle,
    const katran_new_real_t* real,
    const katran_vip_key_t* vip);

/**
 * @brief Remove a real server from a VIP.
 *
 * Removes the backend server from the VIP's consistent hash ring.
 * The real's weight field is ignored for deletion.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] real   Real server to remove. Must not be NULL.
 * @param[in] vip    VIP to remove the real from. Must not be NULL.
 *
 * @return KATRAN_OK on success (including if real wasn't associated with VIP).
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_real_for_vip(
    katran_lb_t handle,
    const katran_new_real_t* real,
    const katran_vip_key_t* vip);

/**
 * @brief Get all real servers for a VIP.
 *
 * Retrieves the list of backend servers configured for a VIP along
 * with their weights and flags.
 *
 * Usage pattern (two-call):
 *   1. Call with reals=NULL to get count
 *   2. Allocate array of katran_new_real_t[count]
 *   3. Call again with allocated array
 *   4. Free with katran_free_reals() when done
 *
 * @param[in]     handle Valid Katran handle.
 * @param[in]     vip    VIP to query. Must not be NULL.
 * @param[out]    reals  Array to receive reals, or NULL to query count.
 * @param[in,out] count  On input: size of reals array (if non-NULL).
 *                       On output: number of reals returned/available.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle, vip, or count is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_reals_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_new_real_t* reals,
    size_t* count);

/**
 * @brief Batch modify real servers for a VIP.
 *
 * Adds or removes multiple real servers in a single operation.
 * More efficient than individual add/del calls for bulk updates.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] action KATRAN_ACTION_ADD or KATRAN_ACTION_DEL.
 * @param[in] reals  Array of real servers to modify. Must not be NULL if count > 0.
 * @param[in] count  Number of elements in reals array.
 * @param[in] vip    VIP to modify. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any required pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 * @return KATRAN_ERR_SPACE_EXHAUSTED if adding would exceed max_reals.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_modify_reals_for_vip(
    katran_lb_t handle,
    katran_modify_action_t action,
    const katran_new_real_t* reals,
    size_t count,
    const katran_vip_key_t* vip);

/**
 * @brief Get the internal index for a real server address.
 *
 * Retrieves the internal numeric index assigned to a real server.
 * This index can be used with functions like katran_lb_get_real_stats().
 *
 * @param[in]  handle  Valid Katran handle.
 * @param[in]  address Real server IP address as string. Must not be NULL.
 * @param[out] index   Pointer to receive the index. Must not be NULL.
 *                     Set to -1 if the real is not found.
 *
 * @return KATRAN_OK on success (even if real not found, check *index).
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 */
katran_error_t katran_lb_get_index_for_real(
    katran_lb_t handle,
    const char* address,
    int64_t* index);

/**
 * @brief Modify flags on a real server.
 *
 * Sets or clears flags on a real server globally (affects all VIPs
 * using this real).
 *
 * @param[in] handle  Valid Katran handle.
 * @param[in] address Real server IP address. Must not be NULL.
 * @param[in] flags   Flag bits to modify.
 * @param[in] set     1 to set the flags, 0 to clear them.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or address is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the real does not exist.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_modify_real(
    katran_lb_t handle,
    const char* address,
    uint8_t flags,
    int set);

/* ============================================================================
 * QUIC MAPPING
 * ============================================================================ */

/**
 * @brief Modify QUIC connection ID to real server mappings.
 *
 * Adds or removes mappings between QUIC host IDs (embedded in connection IDs)
 * and real server addresses. This enables stateful QUIC routing.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] action KATRAN_ACTION_ADD or KATRAN_ACTION_DEL.
 * @param[in] reals  Array of QUIC real mappings. Must not be NULL if count > 0.
 * @param[in] count  Number of elements in reals array.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL or reals is NULL with count > 0.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_modify_quic_reals_mapping(
    katran_lb_t handle,
    katran_modify_action_t action,
    const katran_quic_real_t* reals,
    size_t count);

/**
 * @brief Get all QUIC connection ID mappings.
 *
 * Retrieves all configured QUIC host ID to real server mappings.
 *
 * Usage pattern (two-call):
 *   1. Call with reals=NULL to get count
 *   2. Allocate array of katran_quic_real_t[count]
 *   3. Call again with allocated array
 *   4. Free with katran_free_quic_reals() when done
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    reals  Array to receive mappings, or NULL to query count.
 * @param[in,out] count  On input: size of reals array (if non-NULL).
 *                       On output: number of mappings returned/available.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_quic_reals_mapping(
    katran_lb_t handle,
    katran_quic_real_t* reals,
    size_t* count);

/* ============================================================================
 * SOURCE ROUTING
 * ============================================================================ */

/**
 * @brief Add source-based routing rules.
 *
 * Configures routing based on source IP prefixes. Packets from matching
 * source addresses are forwarded to the specified destination.
 *
 * Requires KATRAN_FEATURE_SRC_ROUTING to be enabled.
 *
 * @param[in] handle       Valid Katran handle.
 * @param[in] src_prefixes Array of source IP prefixes (CIDR notation, e.g., "10.0.0.0/8").
 *                         Must not be NULL if count > 0.
 * @param[in] count        Number of source prefixes.
 * @param[in] dst          Destination address for matching traffic. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if required pointers are NULL or addresses invalid.
 * @return KATRAN_ERR_FEATURE_DISABLED if source routing is not enabled.
 * @return KATRAN_ERR_SPACE_EXHAUSTED if max_lpm_src_size is exceeded.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_src_routing_rule(
    katran_lb_t handle,
    const char** src_prefixes,
    size_t count,
    const char* dst);

/**
 * @brief Delete source-based routing rules.
 *
 * Removes routing rules for the specified source prefixes.
 *
 * @param[in] handle       Valid Katran handle.
 * @param[in] src_prefixes Array of source IP prefixes to remove.
 *                         Must not be NULL if count > 0.
 * @param[in] count        Number of source prefixes.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if source routing is not enabled.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_src_routing_rule(
    katran_lb_t handle,
    const char** src_prefixes,
    size_t count);

/**
 * @brief Clear all source-based routing rules.
 *
 * Removes all configured source routing rules.
 *
 * @param[in] handle Valid Katran handle.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if source routing is not enabled.
 * @return KATRAN_ERR_BPF_FAILED if clearing BPF maps fails.
 */
katran_error_t katran_lb_clear_all_src_routing_rules(katran_lb_t handle);

/**
 * @brief Get all source-based routing rules.
 *
 * Retrieves all configured source prefix to destination mappings.
 *
 * Usage pattern (two-call):
 *   1. Call with srcs=NULL and dsts=NULL to get count
 *   2. Allocate arrays of char*[count] for both srcs and dsts
 *   3. Call again with allocated arrays
 *   4. Free with katran_free_src_routing_rules() when done
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    srcs   Array to receive source prefixes, or NULL to query count.
 * @param[out]    dsts   Array to receive destinations, or NULL to query count.
 * @param[in,out] count  On input: size of arrays (if non-NULL).
 *                       On output: number of rules returned/available.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if source routing is not enabled.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_src_routing_rule(
    katran_lb_t handle,
    char** srcs,
    char** dsts,
    size_t* count);

/**
 * @brief Get the number of source routing rules.
 *
 * Returns the count of currently configured source routing rules.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] size   Pointer to receive the count. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or size is NULL.
 */
katran_error_t katran_lb_get_src_routing_rule_size(
    katran_lb_t handle,
    uint32_t* size);

/* ============================================================================
 * INLINE DECAPSULATION
 * ============================================================================ */

/**
 * @brief Add an inline decapsulation destination.
 *
 * Configures an IP address for which incoming encapsulated packets
 * should be decapsulated in the XDP program.
 *
 * Requires KATRAN_FEATURE_INLINE_DECAP to be enabled.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] dst    Destination IP address to enable decapsulation for.
 *                   Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or dst is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if inline decap is not enabled.
 * @return KATRAN_ERR_SPACE_EXHAUSTED if max_decap_dst is exceeded.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_inline_decap_dst(
    katran_lb_t handle,
    const char* dst);

/**
 * @brief Remove an inline decapsulation destination.
 *
 * Removes the destination from the decapsulation list.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] dst    Destination IP address to remove. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or dst is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if inline decap is not enabled.
 * @return KATRAN_ERR_NOT_FOUND if the destination was not configured.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_inline_decap_dst(
    katran_lb_t handle,
    const char* dst);

/**
 * @brief Get all inline decapsulation destinations.
 *
 * Retrieves all configured decapsulation destination addresses.
 *
 * Usage pattern (two-call):
 *   1. Call with dsts=NULL to get count
 *   2. Allocate array of char*[count]
 *   3. Call again with allocated array
 *   4. Free with katran_free_strings() when done
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    dsts   Array to receive destination addresses, or NULL to query count.
 * @param[in,out] count  On input: size of dsts array (if non-NULL).
 *                       On output: number of destinations returned/available.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_inline_decap_dst(
    katran_lb_t handle,
    char** dsts,
    size_t* count);

/* ============================================================================
 * HEALTHCHECKING
 * ============================================================================ */

/**
 * @brief Add a healthcheck destination mapping.
 *
 * Maps a socket mark (SO_MARK) to a destination address for healthcheck
 * packet encapsulation. Packets with the specified mark will be
 * encapsulated and sent to the destination.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] somark Socket mark value to match.
 * @param[in] dst    Destination address for encapsulated packets. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or dst is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if healthchecking is not enabled.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_healthchecker_dst(
    katran_lb_t handle,
    uint32_t somark,
    const char* dst);

/**
 * @brief Remove a healthcheck destination mapping.
 *
 * Removes the socket mark to destination mapping.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] somark Socket mark value to remove.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the somark was not configured.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_healthchecker_dst(
    katran_lb_t handle,
    uint32_t somark);

/**
 * @brief Get all healthcheck destination mappings.
 *
 * Retrieves all configured socket mark to destination mappings.
 *
 * Usage pattern (two-call):
 *   1. Call with somarks=NULL and dsts=NULL to get count
 *   2. Allocate arrays
 *   3. Call again with allocated arrays
 *   4. Free dsts with katran_free_strings() when done
 *
 * @param[in]     handle  Valid Katran handle.
 * @param[out]    somarks Array to receive socket marks, or NULL to query count.
 * @param[out]    dsts    Array to receive destinations, or NULL to query count.
 * @param[in,out] count   On input: size of arrays (if non-NULL).
 *                        On output: number of mappings returned/available.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 * @return KATRAN_ERR_MEMORY if internal allocation fails.
 */
katran_error_t katran_lb_get_healthcheckers_dst(
    katran_lb_t handle,
    uint32_t* somarks,
    char** dsts,
    size_t* count);

/**
 * @brief Add a healthcheck key for per-key statistics.
 *
 * Registers a VIP-like key for healthcheck packet tracking.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] hc_key Healthcheck key to add. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or hc_key is NULL.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_hc_key(
    katran_lb_t handle,
    const katran_vip_key_t* hc_key);

/**
 * @brief Remove a healthcheck key.
 *
 * Unregisters a healthcheck key.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] hc_key Healthcheck key to remove. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or hc_key is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the key was not registered.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_del_hc_key(
    katran_lb_t handle,
    const katran_vip_key_t* hc_key);

/* ============================================================================
 * STATISTICS
 * ============================================================================ */

/**
 * @brief Get packet/byte statistics for a VIP.
 *
 * Returns the total packets and bytes forwarded to a VIP's backends.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[in]  vip    VIP to query. Must not be NULL.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = packets, v2 = bytes.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 */
katran_error_t katran_lb_get_stats_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_lb_stats_t* stats);

/**
 * @brief Get decapsulation statistics for a VIP.
 *
 * Returns decapsulation packet counts for a VIP.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[in]  vip    VIP to query. Must not be NULL.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the VIP does not exist.
 */
katran_error_t katran_lb_get_decap_stats_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_lb_stats_t* stats);

/**
 * @brief Get LRU cache statistics.
 *
 * Returns total packets processed and LRU hits/misses.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = total packets, v2 = LRU hits.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_lru_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get LRU miss statistics.
 *
 * Returns breakdown of LRU misses by cause.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = TCP SYN misses, v2 = non-SYN misses.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_lru_miss_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get LRU fallback statistics.
 *
 * Returns count of fallback LRU cache hits.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_lru_fallback_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get ICMP "too big" statistics.
 *
 * Returns count of ICMP Packet Too Big messages generated.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = ICMPv4 count, v2 = ICMPv6 count.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_icmp_too_big_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get consistent hash drop statistics.
 *
 * Returns count of packets dropped during consistent hashing.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = real ID out of bounds, v2 = real #0 (unmapped).
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_ch_drop_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get source routing statistics.
 *
 * Returns local vs remote routing statistics.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = local backend, v2 = remote (LPM matched).
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_src_routing_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get inline decapsulation statistics.
 *
 * Returns count of inline-decapsulated packets.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_inline_decap_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get global LRU statistics.
 *
 * Returns global LRU cache statistics.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = map lookup failures, v2 = global LRU routed.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_global_lru_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get general decapsulation statistics.
 *
 * Returns v4/v6 decapsulation counts.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = IPv4 decapped, v2 = IPv6 decapped.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_decap_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get QUIC ICMP statistics.
 *
 * Returns QUIC-related ICMP message statistics.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_quic_icmp_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get per-real server statistics.
 *
 * Returns packet/byte statistics for a specific real server by index.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[in]  index  Real server index (from katran_lb_get_index_for_real()).
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = packets, v2 = bytes.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_real_stats(
    katran_lb_t handle,
    uint32_t index,
    katran_lb_stats_t* stats);

/**
 * @brief Get total XDP statistics.
 *
 * Returns total packets and bytes processed by XDP program.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *                    v1 = packets, v2 = bytes.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_xdp_total_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get XDP TX statistics.
 *
 * Returns packets/bytes forwarded (XDP_TX).
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_xdp_tx_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get XDP drop statistics.
 *
 * Returns packets/bytes dropped (XDP_DROP).
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_xdp_drop_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get XDP pass statistics.
 *
 * Returns packets/bytes passed to kernel (XDP_PASS).
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_xdp_pass_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats);

/**
 * @brief Get TCP server ID routing statistics.
 *
 * Returns statistics for TCP Passive Routing.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_tcp_server_id_routing_stats(
    katran_lb_t handle,
    katran_tpr_packets_stats_t* stats);

/**
 * @brief Get QUIC packet routing statistics.
 *
 * Returns detailed QUIC routing statistics.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_quic_packets_stats(
    katran_lb_t handle,
    katran_quic_packets_stats_t* stats);

/**
 * @brief Get healthcheck program statistics.
 *
 * Returns packet counters from the healthcheck BPF program.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if healthchecking is not enabled.
 */
katran_error_t katran_lb_get_hc_prog_stats(
    katran_lb_t handle,
    katran_hc_stats_t* stats);

/**
 * @brief Get BPF map statistics.
 *
 * Returns capacity and current usage of a BPF map.
 *
 * @param[in]  handle   Valid Katran handle.
 * @param[in]  map_name Name of the BPF map. Must not be NULL.
 * @param[out] stats    Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the map does not exist.
 */
katran_error_t katran_lb_get_bpf_map_stats(
    katran_lb_t handle,
    const char* map_name,
    katran_bpf_map_stats_t* stats);

/**
 * @brief Get userspace library statistics.
 *
 * Returns statistics about the userspace component.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 */
katran_error_t katran_lb_get_userspace_stats(
    katran_lb_t handle,
    katran_userspace_stats_t* stats);

/**
 * @brief Get per-core packet statistics.
 *
 * Returns packet counts processed on each CPU core.
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    counts Array to receive per-core counts, or NULL to query count.
 * @param[in,out] count  On input: size of counts array (if non-NULL).
 *                       On output: number of cores.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 */
katran_error_t katran_lb_get_per_core_packets_stats(
    katran_lb_t handle,
    int64_t* counts,
    size_t* count);

/**
 * @brief Check if the system is under flood conditions.
 *
 * Examines connection rate statistics to determine if the system
 * is experiencing a traffic flood.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] result Pointer to receive result. Must not be NULL.
 *                    Set to 1 if under flood, 0 otherwise.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or result is NULL.
 */
katran_error_t katran_lb_is_under_flood(
    katran_lb_t handle,
    int* result);

/* ============================================================================
 * FLOW SIMULATION
 * ============================================================================ */

/**
 * @brief Determine which real server a flow would be routed to.
 *
 * Simulates the routing decision for a given 5-tuple without actually
 * processing a packet. Useful for debugging and verification.
 *
 * @param[in]  handle       Valid Katran handle.
 * @param[in]  flow         Flow 5-tuple to simulate. Must not be NULL.
 * @param[out] real_address Buffer to receive the real server address.
 *                          Must not be NULL.
 * @param[in]  buffer_size  Size of real_address buffer.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL or buffer too small.
 * @return KATRAN_ERR_NOT_FOUND if the flow doesn't match any VIP.
 */
katran_error_t katran_lb_get_real_for_flow(
    katran_lb_t handle,
    const katran_flow_t* flow,
    char* real_address,
    size_t buffer_size);

/**
 * @brief Simulate packet processing through the BPF program.
 *
 * Processes a raw packet through the Katran BPF program and returns
 * the resulting packet. Note: This affects BPF state (maps, stats).
 *
 * @param[in]  handle     Valid Katran handle.
 * @param[in]  in_packet  Input packet data (starting with Ethernet header).
 *                        Must not be NULL.
 * @param[in]  in_size    Size of input packet in bytes.
 * @param[out] out_packet Buffer to receive output packet, or NULL to query size.
 * @param[in,out] out_size On input: size of out_packet buffer.
 *                         On output: size of output packet.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if required pointers are NULL.
 * @return KATRAN_ERR_BPF_FAILED if simulation fails.
 */
katran_error_t katran_lb_simulate_packet(
    katran_lb_t handle,
    const uint8_t* in_packet,
    size_t in_size,
    uint8_t* out_packet,
    size_t* out_size);

/* ============================================================================
 * FEATURE MANAGEMENT
 * ============================================================================ */

/**
 * @brief Check if a feature is available.
 *
 * Queries whether a specific optional feature is enabled in the
 * current BPF program.
 *
 * @param[in]  handle      Valid Katran handle.
 * @param[in]  feature     Feature flag to check.
 * @param[out] has_feature Pointer to receive result. Must not be NULL.
 *                         Set to 1 if feature is available, 0 otherwise.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or has_feature is NULL.
 */
katran_error_t katran_lb_has_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    int* has_feature);

/**
 * @brief Install a feature by reloading the BPF program.
 *
 * If the feature is not currently available, attempts to reload the
 * BPF program from the specified path to enable it.
 *
 * @param[in] handle    Valid Katran handle.
 * @param[in] feature   Feature to install.
 * @param[in] prog_path Path to BPF program with the feature, or NULL to use
 *                      current program path.
 *
 * @return KATRAN_OK on success (feature now available).
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_BPF_FAILED if reload fails or feature not in new program.
 */
katran_error_t katran_lb_install_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    const char* prog_path);

/**
 * @brief Remove a feature by reloading the BPF program.
 *
 * If the feature is currently available, attempts to reload the
 * BPF program from the specified path to disable it.
 *
 * @param[in] handle    Valid Katran handle.
 * @param[in] feature   Feature to remove.
 * @param[in] prog_path Path to BPF program without the feature, or NULL.
 *
 * @return KATRAN_OK on success (feature now unavailable).
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_BPF_FAILED if reload fails or feature still in new program.
 */
katran_error_t katran_lb_remove_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    const char* prog_path);

/* ============================================================================
 * LRU OPERATIONS
 * ============================================================================ */

/**
 * @brief Delete an LRU entry for a specific flow.
 *
 * Removes the connection tracking entry for the specified flow from
 * all per-CPU and fallback LRU maps.
 *
 * @param[in]     handle   Valid Katran handle.
 * @param[in]     dst_vip  Destination VIP. Must not be NULL.
 * @param[in]     src_ip   Source IP address. Must not be NULL.
 * @param[in]     src_port Source port.
 * @param[out]    maps     Array to receive names of maps where entry was deleted,
 *                         or NULL to skip.
 * @param[in,out] count    On input: size of maps array (if non-NULL).
 *                         On output: number of maps where entry was deleted.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if required pointers are NULL.
 */
katran_error_t katran_lb_delete_lru(
    katran_lb_t handle,
    const katran_vip_key_t* dst_vip,
    const char* src_ip,
    uint16_t src_port,
    char** maps,
    size_t* count);

/**
 * @brief Purge all LRU entries for a VIP.
 *
 * Removes all connection tracking entries for the specified VIP
 * from all LRU maps. Useful when removing a VIP or for cache invalidation.
 *
 * @param[in]  handle        Valid Katran handle.
 * @param[in]  dst_vip       VIP to purge entries for. Must not be NULL.
 * @param[out] deleted_count Pointer to receive count of deleted entries.
 *                           May be NULL if count not needed.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or dst_vip is NULL.
 */
katran_error_t katran_lb_purge_vip_lru(
    katran_lb_t handle,
    const katran_vip_key_t* dst_vip,
    int* deleted_count);

/* ============================================================================
 * MONITORING
 * ============================================================================ */

/**
 * @brief Stop the packet monitor.
 *
 * Stops packet capture/introspection if running.
 *
 * @param[in] handle Valid Katran handle.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if introspection is not enabled.
 */
katran_error_t katran_lb_stop_monitor(katran_lb_t handle);

/**
 * @brief Restart the packet monitor.
 *
 * Restarts packet capture with the specified packet limit.
 *
 * @param[in] handle  Valid Katran handle.
 * @param[in] limit   Maximum number of packets to capture (0 = unlimited).
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if introspection is not enabled.
 */
katran_error_t katran_lb_restart_monitor(
    katran_lb_t handle,
    uint32_t limit);

/**
 * @brief Get monitor statistics.
 *
 * Returns statistics from the packet capture subsystem.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] stats  Pointer to receive statistics. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or stats is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if introspection is not enabled.
 */
katran_error_t katran_lb_get_monitor_stats(
    katran_lb_t handle,
    katran_monitor_stats_t* stats);

/* ============================================================================
 * UTILITY FUNCTIONS
 * ============================================================================ */

/**
 * @brief Get the file descriptor of the balancer BPF program.
 *
 * Returns the FD for the loaded balancer XDP program.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] fd     Pointer to receive the file descriptor. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or fd is NULL.
 */
katran_error_t katran_lb_get_katran_prog_fd(katran_lb_t handle, int* fd);

/**
 * @brief Get the file descriptor of the healthcheck BPF program.
 *
 * Returns the FD for the loaded healthcheck TC program.
 *
 * @param[in]  handle Valid Katran handle.
 * @param[out] fd     Pointer to receive the file descriptor. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or fd is NULL.
 * @return KATRAN_ERR_FEATURE_DISABLED if healthchecking is not enabled.
 */
katran_error_t katran_lb_get_healthchecker_prog_fd(katran_lb_t handle, int* fd);

/**
 * @brief Get a BPF map file descriptor by name.
 *
 * Returns the FD for a named BPF map.
 *
 * @param[in]  handle   Valid Katran handle.
 * @param[in]  map_name Name of the BPF map. Must not be NULL.
 * @param[out] fd       Pointer to receive the file descriptor. Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if any pointer is NULL.
 * @return KATRAN_ERR_NOT_FOUND if the map does not exist.
 */
katran_error_t katran_lb_get_bpf_map_fd_by_name(
    katran_lb_t handle,
    const char* map_name,
    int* fd);

/**
 * @brief Get file descriptors for global LRU maps.
 *
 * Returns FDs for all per-CPU global LRU maps.
 *
 * @param[in]     handle Valid Katran handle.
 * @param[out]    fds    Array to receive file descriptors, or NULL to query count.
 * @param[in,out] count  On input: size of fds array (if non-NULL).
 *                       On output: number of FDs.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or count is NULL.
 */
katran_error_t katran_lb_get_global_lru_maps_fds(
    katran_lb_t handle,
    int* fds,
    size_t* count);

/**
 * @brief Add a source IP for packet encapsulation.
 *
 * Sets the source IP address to use when Katran encapsulates packets
 * (GUE or IPIP). Replaces any existing source of the same address family.
 *
 * @param[in] handle Valid Katran handle.
 * @param[in] src    Source IP address (IPv4 or IPv6). Must not be NULL.
 *
 * @return KATRAN_OK on success.
 * @return KATRAN_ERR_INVALID_ARGUMENT if handle or src is NULL, or address invalid.
 * @return KATRAN_ERR_BPF_FAILED if updating BPF maps fails.
 */
katran_error_t katran_lb_add_src_ip_for_pckt_encap(
    katran_lb_t handle,
    const char* src);

/**
 * @brief Get the last error message.
 *
 * Returns a human-readable error message for the last error that occurred
 * on the current thread. The returned string is thread-local and remains
 * valid until the next API call on the same thread.
 *
 * @return Pointer to error message string, or empty string if no error.
 *         Never returns NULL.
 */
const char* katran_lb_get_last_error(void);

/* ============================================================================
 * MEMORY MANAGEMENT
 * ============================================================================ */

/**
 * @brief Free a VIP array allocated by the API.
 *
 * Frees memory allocated by katran_lb_get_all_vips().
 *
 * @param[in] vips  Array to free. May be NULL (no-op).
 * @param[in] count Number of elements in the array.
 */
void katran_free_vips(katran_vip_key_t* vips, size_t count);

/**
 * @brief Free a reals array allocated by the API.
 *
 * Frees memory allocated by katran_lb_get_reals_for_vip().
 *
 * @param[in] reals Array to free. May be NULL (no-op).
 * @param[in] count Number of elements in the array.
 */
void katran_free_reals(katran_new_real_t* reals, size_t count);

/**
 * @brief Free a QUIC reals array allocated by the API.
 *
 * Frees memory allocated by katran_lb_get_quic_reals_mapping().
 *
 * @param[in] reals Array to free. May be NULL (no-op).
 * @param[in] count Number of elements in the array.
 */
void katran_free_quic_reals(katran_quic_real_t* reals, size_t count);

/**
 * @brief Free a string array allocated by the API.
 *
 * Frees memory allocated by functions that return string arrays.
 *
 * @param[in] strings Array to free. May be NULL (no-op).
 * @param[in] count   Number of elements in the array.
 */
void katran_free_strings(char** strings, size_t count);

/**
 * @brief Free source routing rule arrays allocated by the API.
 *
 * Frees memory allocated by katran_lb_get_src_routing_rule().
 *
 * @param[in] srcs  Source prefixes array to free. May be NULL.
 * @param[in] dsts  Destinations array to free. May be NULL.
 * @param[in] count Number of elements in each array.
 */
void katran_free_src_routing_rules(char** srcs, char** dsts, size_t count);

#ifdef __cplusplus
}
#endif

#endif /* KATRAN_CAPI_H */
