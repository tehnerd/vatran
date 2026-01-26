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

#include "katran_capi.h"

#include <cstring>
#include <memory>
#include <stdexcept>
#include <string>
#include <vector>

#include <folly/IPAddress.h>

#include "katran/lib/BpfAdapter.h"
#include "katran/lib/KatranLb.h"
#include "katran/lib/KatranLbStructs.h"
#include "katran/lib/KatranSimulatorUtils.h"

// Thread-local storage for the last error message
thread_local std::string g_last_error;

/**
 * Internal handle structure containing the KatranLb instance.
 */
struct katran_lb_handle {
  std::unique_ptr<katran::KatranLb> lb;
};

/**
 * Helper to set the thread-local error message and return an error code.
 */
static katran_error_t setError(katran_error_t code, const std::string& msg) {
  g_last_error = msg;
  return code;
}

/**
 * Helper to clear the error state.
 */
static void clearError() {
  g_last_error.clear();
}

/**
 * Helper to convert C VipKey to C++ VipKey.
 */
static katran::VipKey toVipKey(const katran_vip_key_t* vip) {
  katran::VipKey result;
  result.address = vip->address ? vip->address : "";
  result.port = vip->port;
  result.proto = vip->proto;
  return result;
}

/**
 * Helper to convert C NewReal to C++ NewReal.
 */
static katran::NewReal toNewReal(const katran_new_real_t* real) {
  katran::NewReal result;
  result.address = real->address ? real->address : "";
  result.weight = real->weight;
  result.flags = real->flags;
  return result;
}

/**
 * Helper to convert C QuicReal to C++ QuicReal.
 */
static katran::QuicReal toQuicReal(const katran_quic_real_t* real) {
  katran::QuicReal result;
  result.address = real->address ? real->address : "";
  result.id = real->id;
  return result;
}

/**
 * Helper to convert C KatranFlow to C++ KatranFlow.
 */
static katran::KatranFlow toKatranFlow(const katran_flow_t* flow) {
  katran::KatranFlow result;
  result.src = flow->src ? flow->src : "";
  result.dst = flow->dst ? flow->dst : "";
  result.srcPort = flow->src_port;
  result.dstPort = flow->dst_port;
  result.proto = flow->proto;
  return result;
}

/**
 * Helper to convert C HashFunction to C++ HashFunction.
 */
static katran::HashFunction toHashFunction(katran_hash_function_t func) {
  switch (func) {
    case KATRAN_HASH_MAGLEV_V2:
      return katran::HashFunction::MaglevV2;
    case KATRAN_HASH_MAGLEV:
    default:
      return katran::HashFunction::Maglev;
  }
}

/**
 * Helper to convert C feature enum to C++ feature enum.
 */
static katran::KatranFeatureEnum toFeatureEnum(katran_feature_t feature) {
  switch (feature) {
    case KATRAN_FEATURE_SRC_ROUTING:
      return katran::KatranFeatureEnum::SrcRouting;
    case KATRAN_FEATURE_INLINE_DECAP:
      return katran::KatranFeatureEnum::InlineDecap;
    case KATRAN_FEATURE_INTROSPECTION:
      return katran::KatranFeatureEnum::Introspection;
    case KATRAN_FEATURE_GUE_ENCAP:
      return katran::KatranFeatureEnum::GueEncap;
    case KATRAN_FEATURE_DIRECT_HC:
      return katran::KatranFeatureEnum::DirectHealthchecking;
    case KATRAN_FEATURE_LOCAL_DELIVERY_OPT:
      return katran::KatranFeatureEnum::LocalDeliveryOptimization;
    case KATRAN_FEATURE_FLOW_DEBUG:
      return katran::KatranFeatureEnum::FlowDebug;
    default:
      return katran::KatranFeatureEnum::SrcRouting;
  }
}

/**
 * Helper to convert C ModifyAction to C++ ModifyAction.
 */
static katran::ModifyAction toModifyAction(katran_modify_action_t action) {
  return action == KATRAN_ACTION_ADD ? katran::ModifyAction::ADD
                                     : katran::ModifyAction::DEL;
}

/**
 * Helper to duplicate a string (caller must free).
 */
static char* duplicateString(const std::string& str) {
  char* result = static_cast<char*>(malloc(str.size() + 1));
  if (result) {
    memcpy(result, str.c_str(), str.size() + 1);
  }
  return result;
}

extern "C" {

/* ============================================================================
 * INITIALIZATION AND LIFECYCLE
 * ============================================================================ */

katran_error_t katran_config_init(katran_config_t* config) {
  if (!config) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL config pointer");
  }

  memset(config, 0, sizeof(katran_config_t));

  // Set defaults matching KatranConfig
  config->root_map_pos = 2;
  config->use_root_map = 1;
  config->max_vips = 512;
  config->max_reals = 4096;
  config->ch_ring_size = 65537;
  config->lru_size = 8000000;
  config->max_lpm_src_size = 3000000;
  config->max_decap_dst = 6;
  config->global_lru_size = 100000;
  config->enable_hc = 1;
  config->tunnel_based_hc_encap = 1;
  config->memlock_unlimited = 1;
  config->cleanup_on_shutdown = 1;
  config->priority = 2307;
  config->hash_function = KATRAN_HASH_MAGLEV;

  clearError();
  return KATRAN_OK;
}

katran_error_t katran_lb_create(
    const katran_config_t* config,
    katran_lb_t* handle) {
  if (!config || !handle) {
    return setError(
        KATRAN_ERR_INVALID_ARGUMENT, "NULL config or handle pointer");
  }

  *handle = nullptr;

  try {
    // Build C++ config from C config
    katran::KatranConfig cppConfig;

    cppConfig.mainInterface =
        config->main_interface ? config->main_interface : "";
    cppConfig.v4TunInterface =
        config->v4_tun_interface ? config->v4_tun_interface : "";
    cppConfig.v6TunInterface =
        config->v6_tun_interface ? config->v6_tun_interface : "";
    cppConfig.hcInterface = config->hc_interface ? config->hc_interface : "";
    cppConfig.balancerProgPath =
        config->balancer_prog_path ? config->balancer_prog_path : "";
    cppConfig.healthcheckingProgPath =
        config->healthchecking_prog_path ? config->healthchecking_prog_path
                                         : "";

    if (config->default_mac) {
      cppConfig.defaultMac.assign(
          config->default_mac, config->default_mac + 6);
    }
    if (config->local_mac) {
      cppConfig.localMac.assign(config->local_mac, config->local_mac + 6);
    }

    cppConfig.rootMapPath = config->root_map_path ? config->root_map_path : "";
    cppConfig.rootMapPos = config->root_map_pos;
    cppConfig.useRootMap = config->use_root_map != 0;

    cppConfig.maxVips = config->max_vips;
    cppConfig.maxReals = config->max_reals;
    cppConfig.chRingSize = config->ch_ring_size;
    cppConfig.LruSize = config->lru_size;
    cppConfig.maxLpmSrcSize = config->max_lpm_src_size;
    cppConfig.maxDecapDst = config->max_decap_dst;
    cppConfig.globalLruSize = config->global_lru_size;

    cppConfig.enableHc = config->enable_hc != 0;
    cppConfig.tunnelBasedHCEncap = config->tunnel_based_hc_encap != 0;
    cppConfig.testing = config->testing != 0;
    cppConfig.memlockUnlimited = config->memlock_unlimited != 0;
    cppConfig.flowDebug = config->flow_debug != 0;
    cppConfig.enableCidV3 = config->enable_cid_v3 != 0;
    cppConfig.cleanupOnShutdown = config->cleanup_on_shutdown != 0;

    if (config->forwarding_cores && config->forwarding_cores_count > 0) {
      cppConfig.forwardingCores.assign(
          config->forwarding_cores,
          config->forwarding_cores + config->forwarding_cores_count);
    }
    if (config->numa_nodes && config->numa_nodes_count > 0) {
      cppConfig.numaNodes.assign(
          config->numa_nodes,
          config->numa_nodes + config->numa_nodes_count);
    }

    cppConfig.xdpAttachFlags = config->xdp_attach_flags;
    cppConfig.priority = config->priority;
    cppConfig.mainInterfaceIndex = config->main_interface_index;
    cppConfig.hcInterfaceIndex = config->hc_interface_index;

    cppConfig.katranSrcV4 =
        config->katran_src_v4 ? config->katran_src_v4 : "";
    cppConfig.katranSrcV6 =
        config->katran_src_v6 ? config->katran_src_v6 : "";

    cppConfig.hashFunction = toHashFunction(config->hash_function);

    // Create BpfAdapter and KatranLb
    auto bpfAdapter = std::make_unique<katran::BpfAdapter>(true);
    auto internal = new katran_lb_handle();
    internal->lb =
        std::make_unique<katran::KatranLb>(cppConfig, std::move(bpfAdapter));

    *handle = internal;
    clearError();
    return KATRAN_OK;

  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_destroy(katran_lb_t handle) {
  if (!handle) {
    clearError();
    return KATRAN_OK;
  }

  try {
    delete handle;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_load_bpf_progs(katran_lb_t handle) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    handle->lb->loadBpfProgs();
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

katran_error_t katran_lb_attach_bpf_progs(katran_lb_t handle) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    handle->lb->attachBpfProgs();
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

katran_error_t katran_lb_reload_balancer_prog(
    katran_lb_t handle,
    const char* path,
    const katran_config_t* config) {
  if (!handle || !path) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or path");
  }

  try {
    std::optional<katran::KatranConfig> cppConfig;
    // Note: If config is provided, we would convert it here.
    // For simplicity, we use nullopt for now.
    (void)config; // Silence unused parameter warning

    if (!handle->lb->reloadBalancerProg(path, cppConfig)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to reload balancer prog");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

/* ============================================================================
 * MAC ADDRESS MANAGEMENT
 * ============================================================================ */

katran_error_t katran_lb_change_mac(katran_lb_t handle, const uint8_t* mac) {
  if (!handle || !mac) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or mac");
  }

  try {
    std::vector<uint8_t> macVec(mac, mac + 6);
    if (!handle->lb->changeMac(macVec)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to change MAC");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_mac(katran_lb_t handle, uint8_t* mac) {
  if (!handle || !mac) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or mac");
  }

  try {
    auto macVec = handle->lb->getMac();
    if (macVec.size() >= 6) {
      memcpy(mac, macVec.data(), 6);
    } else {
      memset(mac, 0, 6);
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * VIP MANAGEMENT
 * ============================================================================ */

katran_error_t katran_lb_add_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t flags) {
  if (!handle || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or vip");
  }

  try {
    auto cppVip = toVipKey(vip);
    if (!handle->lb->addVip(cppVip, flags)) {
      return setError(KATRAN_ERR_SPACE_EXHAUSTED, "Failed to add VIP");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip) {
  if (!handle || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or vip");
  }

  try {
    auto cppVip = toVipKey(vip);
    if (!handle->lb->delVip(cppVip)) {
      return setError(KATRAN_ERR_NOT_FOUND, "VIP not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_all_vips(
    katran_lb_t handle,
    katran_vip_key_t* vips,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto allVips = handle->lb->getAllVips();

    if (!vips) {
      // Query count only
      *count = allVips.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, allVips.size());
    for (size_t i = 0; i < toCopy; i++) {
      vips[i].address = duplicateString(allVips[i].address);
      vips[i].port = allVips[i].port;
      vips[i].proto = allVips[i].proto;
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_modify_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t flag,
    int set) {
  if (!handle || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or vip");
  }

  try {
    auto cppVip = toVipKey(vip);
    if (!handle->lb->modifyVip(cppVip, flag, set != 0)) {
      return setError(KATRAN_ERR_NOT_FOUND, "VIP not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_vip_flags(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    uint32_t* flags) {
  if (!handle || !vip || !vip->address || !flags) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppVip = toVipKey(vip);
    *flags = handle->lb->getVipFlags(cppVip);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_NOT_FOUND, e.what());
  }
}

katran_error_t katran_lb_change_hash_function_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_hash_function_t func) {
  if (!handle || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or vip");
  }

  try {
    auto cppVip = toVipKey(vip);
    if (!handle->lb->changeHashFunctionForVip(cppVip, toHashFunction(func))) {
      return setError(KATRAN_ERR_NOT_FOUND, "VIP not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * REAL SERVER MANAGEMENT
 * ============================================================================ */

katran_error_t katran_lb_add_real_for_vip(
    katran_lb_t handle,
    const katran_new_real_t* real,
    const katran_vip_key_t* vip) {
  if (!handle || !real || !real->address || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppReal = toNewReal(real);
    auto cppVip = toVipKey(vip);
    if (!handle->lb->addRealForVip(cppReal, cppVip)) {
      return setError(KATRAN_ERR_NOT_FOUND, "VIP not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_real_for_vip(
    katran_lb_t handle,
    const katran_new_real_t* real,
    const katran_vip_key_t* vip) {
  if (!handle || !real || !real->address || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppReal = toNewReal(real);
    auto cppVip = toVipKey(vip);
    handle->lb->delRealForVip(cppReal, cppVip);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_reals_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_new_real_t* reals,
    size_t* count) {
  if (!handle || !vip || !vip->address || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppVip = toVipKey(vip);
    auto allReals = handle->lb->getRealsForVip(cppVip);

    if (!reals) {
      *count = allReals.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, allReals.size());
    for (size_t i = 0; i < toCopy; i++) {
      reals[i].address = duplicateString(allReals[i].address);
      reals[i].weight = allReals[i].weight;
      reals[i].flags = allReals[i].flags;
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_NOT_FOUND, e.what());
  }
}

katran_error_t katran_lb_modify_reals_for_vip(
    katran_lb_t handle,
    katran_modify_action_t action,
    const katran_new_real_t* reals,
    size_t count,
    const katran_vip_key_t* vip) {
  if (!handle || !vip || !vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or vip");
  }
  if (count > 0 && !reals) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL reals with count > 0");
  }

  try {
    std::vector<katran::NewReal> cppReals;
    cppReals.reserve(count);
    for (size_t i = 0; i < count; i++) {
      cppReals.push_back(toNewReal(&reals[i]));
    }

    auto cppVip = toVipKey(vip);
    if (!handle->lb->modifyRealsForVip(
            toModifyAction(action), cppReals, cppVip)) {
      return setError(KATRAN_ERR_NOT_FOUND, "VIP not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::invalid_argument& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_index_for_real(
    katran_lb_t handle,
    const char* address,
    int64_t* index) {
  if (!handle || !address || !index) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    *index = handle->lb->getIndexForReal(address);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_modify_real(
    katran_lb_t handle,
    const char* address,
    uint8_t flags,
    int set) {
  if (!handle || !address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or address");
  }

  try {
    if (!handle->lb->modifyReal(address, flags, set != 0)) {
      return setError(KATRAN_ERR_NOT_FOUND, "Real not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * QUIC MAPPING
 * ============================================================================ */

katran_error_t katran_lb_modify_quic_reals_mapping(
    katran_lb_t handle,
    katran_modify_action_t action,
    const katran_quic_real_t* reals,
    size_t count) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }
  if (count > 0 && !reals) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL reals with count > 0");
  }

  try {
    std::vector<katran::QuicReal> cppReals;
    cppReals.reserve(count);
    for (size_t i = 0; i < count; i++) {
      cppReals.push_back(toQuicReal(&reals[i]));
    }

    handle->lb->modifyQuicRealsMapping(toModifyAction(action), cppReals);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_quic_reals_mapping(
    katran_lb_t handle,
    katran_quic_real_t* reals,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto mapping = handle->lb->getQuicRealsMapping();

    if (!reals) {
      *count = mapping.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, mapping.size());
    for (size_t i = 0; i < toCopy; i++) {
      reals[i].address = duplicateString(mapping[i].address);
      reals[i].id = mapping[i].id;
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * SOURCE ROUTING
 * ============================================================================ */

katran_error_t katran_lb_add_src_routing_rule(
    katran_lb_t handle,
    const char** src_prefixes,
    size_t count,
    const char* dst) {
  if (!handle || !dst) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or dst");
  }
  if (count > 0 && !src_prefixes) {
    return setError(
        KATRAN_ERR_INVALID_ARGUMENT, "NULL src_prefixes with count > 0");
  }

  try {
    std::vector<std::string> srcs;
    srcs.reserve(count);
    for (size_t i = 0; i < count; i++) {
      if (src_prefixes[i]) {
        srcs.push_back(src_prefixes[i]);
      }
    }

    int errors = handle->lb->addSrcRoutingRule(srcs, dst);
    if (errors > 0) {
      return setError(
          KATRAN_ERR_BPF_FAILED,
          "Failed to add some source routing rules");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_src_routing_rule(
    katran_lb_t handle,
    const char** src_prefixes,
    size_t count) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }
  if (count > 0 && !src_prefixes) {
    return setError(
        KATRAN_ERR_INVALID_ARGUMENT, "NULL src_prefixes with count > 0");
  }

  try {
    std::vector<std::string> srcs;
    srcs.reserve(count);
    for (size_t i = 0; i < count; i++) {
      if (src_prefixes[i]) {
        srcs.push_back(src_prefixes[i]);
      }
    }

    handle->lb->delSrcRoutingRule(srcs);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_clear_all_src_routing_rules(katran_lb_t handle) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    handle->lb->clearAllSrcRoutingRules();
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_src_routing_rule(
    katran_lb_t handle,
    char** srcs,
    char** dsts,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto rules = handle->lb->getSrcRoutingRule();

    if (!srcs || !dsts) {
      *count = rules.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, rules.size());
    size_t i = 0;
    for (const auto& [src, dst] : rules) {
      if (i >= toCopy)
        break;
      srcs[i] = duplicateString(src);
      dsts[i] = duplicateString(dst);
      i++;
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_src_routing_rule_size(
    katran_lb_t handle,
    uint32_t* size) {
  if (!handle || !size) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or size");
  }

  try {
    *size = handle->lb->getSrcRoutingRuleSize();
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * INLINE DECAPSULATION
 * ============================================================================ */

katran_error_t katran_lb_add_inline_decap_dst(
    katran_lb_t handle,
    const char* dst) {
  if (!handle || !dst) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or dst");
  }

  try {
    if (!handle->lb->addInlineDecapDst(dst)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to add decap destination");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_inline_decap_dst(
    katran_lb_t handle,
    const char* dst) {
  if (!handle || !dst) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or dst");
  }

  try {
    if (!handle->lb->delInlineDecapDst(dst)) {
      return setError(KATRAN_ERR_NOT_FOUND, "Decap destination not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_inline_decap_dst(
    katran_lb_t handle,
    char** dsts,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto decapDsts = handle->lb->getInlineDecapDst();

    if (!dsts) {
      *count = decapDsts.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, decapDsts.size());
    for (size_t i = 0; i < toCopy; i++) {
      dsts[i] = duplicateString(decapDsts[i]);
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * HEALTHCHECKING
 * ============================================================================ */

katran_error_t katran_lb_add_healthchecker_dst(
    katran_lb_t handle,
    uint32_t somark,
    const char* dst) {
  if (!handle || !dst) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or dst");
  }

  try {
    if (!handle->lb->addHealthcheckerDst(somark, dst)) {
      return setError(
          KATRAN_ERR_BPF_FAILED, "Failed to add healthchecker destination");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_healthchecker_dst(
    katran_lb_t handle,
    uint32_t somark) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    if (!handle->lb->delHealthcheckerDst(somark)) {
      return setError(KATRAN_ERR_NOT_FOUND, "Healthchecker somark not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_healthcheckers_dst(
    katran_lb_t handle,
    uint32_t* somarks,
    char** dsts,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto hcDsts = handle->lb->getHealthcheckersDst();

    if (!somarks || !dsts) {
      *count = hcDsts.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, hcDsts.size());
    size_t i = 0;
    for (const auto& [mark, dst] : hcDsts) {
      if (i >= toCopy)
        break;
      somarks[i] = mark;
      dsts[i] = duplicateString(dst);
      i++;
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::bad_alloc& e) {
    return setError(KATRAN_ERR_MEMORY, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_add_hc_key(
    katran_lb_t handle,
    const katran_vip_key_t* hc_key) {
  if (!handle || !hc_key || !hc_key->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or hc_key");
  }

  try {
    auto cppHcKey = toVipKey(hc_key);
    if (!handle->lb->addHcKey(cppHcKey)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to add HC key");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_del_hc_key(
    katran_lb_t handle,
    const katran_vip_key_t* hc_key) {
  if (!handle || !hc_key || !hc_key->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or hc_key");
  }

  try {
    auto cppHcKey = toVipKey(hc_key);
    if (!handle->lb->delHcKey(cppHcKey)) {
      return setError(KATRAN_ERR_NOT_FOUND, "HC key not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * STATISTICS
 * ============================================================================ */

katran_error_t katran_lb_get_stats_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_lb_stats_t* stats) {
  if (!handle || !vip || !vip->address || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppVip = toVipKey(vip);
    auto cppStats = handle->lb->getStatsForVip(cppVip);
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_NOT_FOUND, e.what());
  }
}

katran_error_t katran_lb_get_decap_stats_for_vip(
    katran_lb_t handle,
    const katran_vip_key_t* vip,
    katran_lb_stats_t* stats) {
  if (!handle || !vip || !vip->address || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppVip = toVipKey(vip);
    auto cppStats = handle->lb->getDecapStatsForVip(cppVip);
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_NOT_FOUND, e.what());
  }
}

katran_error_t katran_lb_get_lru_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getLruStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_lru_miss_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getLruMissStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_lru_fallback_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getLruFallbackStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_icmp_too_big_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getIcmpTooBigStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_ch_drop_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getChDropStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_src_routing_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getSrcRoutingStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_inline_decap_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getInlineDecapStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_global_lru_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getGlobalLruStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_decap_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getDecapStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_quic_icmp_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getQuicIcmpStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_real_stats(
    katran_lb_t handle,
    uint32_t index,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getRealStats(index);
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_xdp_total_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getXdpTotalStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_xdp_tx_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getXdpTxStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_xdp_drop_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getXdpDropStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_xdp_pass_stats(
    katran_lb_t handle,
    katran_lb_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getXdpPassStats();
    stats->v1 = cppStats.v1;
    stats->v2 = cppStats.v2;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_tcp_server_id_routing_stats(
    katran_lb_t handle,
    katran_tpr_packets_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getTcpServerIdRoutingStats();
    stats->ch_routed = cppStats.ch_routed;
    stats->dst_mismatch_in_lru = cppStats.dst_mismatch_in_lru;
    stats->sid_routed = cppStats.sid_routed;
    stats->tcp_syn = cppStats.tcp_syn;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_quic_packets_stats(
    katran_lb_t handle,
    katran_quic_packets_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getLbQuicPacketsStats();
    stats->ch_routed = cppStats.ch_routed;
    stats->cid_initial = cppStats.cid_initial;
    stats->cid_invalid_server_id = cppStats.cid_invalid_server_id;
    stats->cid_invalid_server_id_sample = cppStats.cid_invalid_server_id_sample;
    stats->cid_routed = cppStats.cid_routed;
    stats->cid_unknown_real_dropped = cppStats.cid_unknown_real_dropped;
    stats->cid_v0 = cppStats.cid_v0;
    stats->cid_v1 = cppStats.cid_v1;
    stats->cid_v2 = cppStats.cid_v2;
    stats->cid_v3 = cppStats.cid_v3;
    stats->dst_match_in_lru = cppStats.dst_match_in_lru;
    stats->dst_mismatch_in_lru = cppStats.dst_mismatch_in_lru;
    stats->dst_not_found_in_lru = cppStats.dst_not_found_in_lru;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_hc_prog_stats(
    katran_lb_t handle,
    katran_hc_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getStatsForHealthCheckProgram();
    stats->packets_processed = cppStats.packetsProcessed;
    stats->packets_dropped = cppStats.packetsDropped;
    stats->packets_skipped = cppStats.packetsSkipped;
    stats->packets_too_big = cppStats.packetsTooBig;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_bpf_map_stats(
    katran_lb_t handle,
    const char* map_name,
    katran_bpf_map_stats_t* stats) {
  if (!handle || !map_name || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppStats = handle->lb->getBpfMapStats(map_name);
    stats->max_entries = cppStats.maxEntries;
    stats->current_entries = cppStats.currentEntries;
    clearError();
    return KATRAN_OK;
  } catch (const std::runtime_error& e) {
    return setError(KATRAN_ERR_NOT_FOUND, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_userspace_stats(
    katran_lb_t handle,
    katran_userspace_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getKatranLbStats();
    stats->bpf_failed_calls = cppStats.bpfFailedCalls;
    stats->addr_validation_failed = cppStats.addrValidationFailed;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_per_core_packets_stats(
    katran_lb_t handle,
    int64_t* counts,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto perCore = handle->lb->getPerCorePacketsStats();

    if (!counts) {
      *count = perCore.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, perCore.size());
    for (size_t i = 0; i < toCopy; i++) {
      counts[i] = perCore[i];
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_is_under_flood(katran_lb_t handle, int* result) {
  if (!handle || !result) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or result");
  }

  try {
    *result = handle->lb->isUnderFlood() ? 1 : 0;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * FLOW SIMULATION
 * ============================================================================ */

katran_error_t katran_lb_get_real_for_flow(
    katran_lb_t handle,
    const katran_flow_t* flow,
    char* real_address,
    size_t buffer_size) {
  if (!handle || !flow || !real_address || buffer_size == 0) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument or zero buffer");
  }

  try {
    auto cppFlow = toKatranFlow(flow);
    auto result = handle->lb->getRealForFlow(cppFlow);

    if (result.empty()) {
      return setError(KATRAN_ERR_NOT_FOUND, "No real found for flow");
    }

    if (result.size() + 1 > buffer_size) {
      return setError(KATRAN_ERR_INVALID_ARGUMENT, "Buffer too small");
    }

    memcpy(real_address, result.c_str(), result.size() + 1);
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_simulate_packet(
    katran_lb_t handle,
    const uint8_t* in_packet,
    size_t in_size,
    uint8_t* out_packet,
    size_t* out_size) {
  if (!handle || !in_packet || !out_size) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    std::string inPkt(reinterpret_cast<const char*>(in_packet), in_size);
    auto result = handle->lb->simulatePacket(inPkt);

    if (!out_packet) {
      *out_size = result.size();
      clearError();
      return KATRAN_OK;
    }

    if (result.size() > *out_size) {
      return setError(KATRAN_ERR_INVALID_ARGUMENT, "Output buffer too small");
    }

    memcpy(out_packet, result.data(), result.size());
    *out_size = result.size();

    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

/* ============================================================================
 * FEATURE MANAGEMENT
 * ============================================================================ */

katran_error_t katran_lb_has_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    int* has_feature) {
  if (!handle || !has_feature) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or has_feature");
  }

  try {
    *has_feature = handle->lb->hasFeature(toFeatureEnum(feature)) ? 1 : 0;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_install_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    const char* prog_path) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    std::string path = prog_path ? prog_path : "";
    if (!handle->lb->installFeature(toFeatureEnum(feature), path)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to install feature");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

katran_error_t katran_lb_remove_feature(
    katran_lb_t handle,
    katran_feature_t feature,
    const char* prog_path) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    std::string path = prog_path ? prog_path : "";
    if (!handle->lb->removeFeature(toFeatureEnum(feature), path)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to remove feature");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_BPF_FAILED, e.what());
  }
}

/* ============================================================================
 * LRU OPERATIONS
 * ============================================================================ */

katran_error_t katran_lb_delete_lru(
    katran_lb_t handle,
    const katran_vip_key_t* dst_vip,
    const char* src_ip,
    uint16_t src_port,
    char** maps,
    size_t* count) {
  if (!handle || !dst_vip || !dst_vip->address || !src_ip) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    auto cppVip = toVipKey(dst_vip);
    auto deletedMaps = handle->lb->deleteLru(cppVip, src_ip, src_port);

    if (count) {
      if (!maps) {
        *count = deletedMaps.size();
      } else {
        size_t toCopy = std::min(*count, deletedMaps.size());
        for (size_t i = 0; i < toCopy; i++) {
          maps[i] = duplicateString(deletedMaps[i]);
        }
        *count = toCopy;
      }
    }

    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_purge_vip_lru(
    katran_lb_t handle,
    const katran_vip_key_t* dst_vip,
    int* deleted_count) {
  if (!handle || !dst_vip || !dst_vip->address) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or dst_vip");
  }

  try {
    auto cppVip = toVipKey(dst_vip);
    auto result = handle->lb->purgeVipLru(cppVip);

    if (!result.error.empty()) {
      return setError(KATRAN_ERR_INTERNAL, result.error);
    }

    if (deleted_count) {
      *deleted_count = result.deletedCount;
    }

    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * MONITORING
 * ============================================================================ */

katran_error_t katran_lb_stop_monitor(katran_lb_t handle) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    if (!handle->lb->stopKatranMonitor()) {
      return setError(
          KATRAN_ERR_FEATURE_DISABLED, "Introspection not enabled");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_restart_monitor(katran_lb_t handle, uint32_t limit) {
  if (!handle) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle");
  }

  try {
    if (!handle->lb->restartKatranMonitor(limit)) {
      return setError(
          KATRAN_ERR_FEATURE_DISABLED, "Introspection not enabled");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_monitor_stats(
    katran_lb_t handle,
    katran_monitor_stats_t* stats) {
  if (!handle || !stats) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or stats");
  }

  try {
    auto cppStats = handle->lb->getKatranMonitorStats();
    stats->limit = cppStats.limit;
    stats->amount = cppStats.amount;
    stats->buffer_full = cppStats.bufferFull;
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

/* ============================================================================
 * UTILITY FUNCTIONS
 * ============================================================================ */

katran_error_t katran_lb_get_katran_prog_fd(katran_lb_t handle, int* fd) {
  if (!handle || !fd) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or fd");
  }

  try {
    *fd = handle->lb->getKatranProgFd();
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_healthchecker_prog_fd(
    katran_lb_t handle,
    int* fd) {
  if (!handle || !fd) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or fd");
  }

  try {
    *fd = handle->lb->getHealthcheckerProgFd();
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_bpf_map_fd_by_name(
    katran_lb_t handle,
    const char* map_name,
    int* fd) {
  if (!handle || !map_name || !fd) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL argument");
  }

  try {
    *fd = handle->lb->getBpfMapFdByName(map_name);
    if (*fd < 0) {
      return setError(KATRAN_ERR_NOT_FOUND, "Map not found");
    }
    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_get_global_lru_maps_fds(
    katran_lb_t handle,
    int* fds,
    size_t* count) {
  if (!handle || !count) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or count");
  }

  try {
    auto mapFds = handle->lb->getGlobalLruMapsFds();

    if (!fds) {
      *count = mapFds.size();
      clearError();
      return KATRAN_OK;
    }

    size_t toCopy = std::min(*count, mapFds.size());
    for (size_t i = 0; i < toCopy; i++) {
      fds[i] = mapFds[i];
    }
    *count = toCopy;

    clearError();
    return KATRAN_OK;
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

katran_error_t katran_lb_add_src_ip_for_pckt_encap(
    katran_lb_t handle,
    const char* src) {
  if (!handle || !src) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, "NULL handle or src");
  }

  try {
    auto ipAddr = folly::IPAddress(src);
    if (!handle->lb->addSrcIpForPcktEncap(ipAddr)) {
      return setError(KATRAN_ERR_BPF_FAILED, "Failed to add source IP");
    }
    clearError();
    return KATRAN_OK;
  } catch (const folly::IPAddressFormatException& e) {
    return setError(KATRAN_ERR_INVALID_ARGUMENT, e.what());
  } catch (const std::exception& e) {
    return setError(KATRAN_ERR_INTERNAL, e.what());
  }
}

const char* katran_lb_get_last_error(void) {
  return g_last_error.c_str();
}

/* ============================================================================
 * MEMORY MANAGEMENT
 * ============================================================================ */

void katran_free_vips(katran_vip_key_t* vips, size_t count) {
  if (!vips)
    return;
  for (size_t i = 0; i < count; i++) {
    free(const_cast<char*>(vips[i].address));
  }
}

void katran_free_reals(katran_new_real_t* reals, size_t count) {
  if (!reals)
    return;
  for (size_t i = 0; i < count; i++) {
    free(const_cast<char*>(reals[i].address));
  }
}

void katran_free_quic_reals(katran_quic_real_t* reals, size_t count) {
  if (!reals)
    return;
  for (size_t i = 0; i < count; i++) {
    free(const_cast<char*>(reals[i].address));
  }
}

void katran_free_strings(char** strings, size_t count) {
  if (!strings)
    return;
  for (size_t i = 0; i < count; i++) {
    free(strings[i]);
  }
}

void katran_free_src_routing_rules(char** srcs, char** dsts, size_t count) {
  katran_free_strings(srcs, count);
  katran_free_strings(dsts, count);
}

} // extern "C"
