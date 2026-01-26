#include <stdio.h>
#include <stdlib.h>

#include "katran_capi.h"

static void die_on_error(katran_error_t rc, const char* what) {
  if (rc == KATRAN_OK) {
    return;
  }
  const char* msg = katran_lb_get_last_error();
  fprintf(stderr, "%s failed: %s (rc=%d)\n", what, msg ? msg : "<no error>", rc);
  exit(1);
}

int main(void) {
  katran_config_t config;
  die_on_error(katran_config_init(&config), "katran_config_init");

  // Run in testing mode so we don't need to load/attach BPF programs.
  config.testing = 1;
  config.enable_hc = 0;
  config.use_root_map = 0;

  katran_lb_t lb = NULL;
  die_on_error(katran_lb_create(&config, &lb), "katran_lb_create");

  katran_vip_key_t vip = {
      .address = "10.0.0.1",
      .port = 80,
      .proto = 6, // TCP
  };
  die_on_error(katran_lb_add_vip(lb, &vip, 0), "katran_lb_add_vip");

  katran_new_real_t real = {
      .address = "10.0.0.2",
      .weight = 1,
      .flags = 0,
  };
  die_on_error(katran_lb_add_real_for_vip(lb, &real, &vip),
               "katran_lb_add_real_for_vip");

  size_t count = 0;
  die_on_error(katran_lb_get_reals_for_vip(lb, &vip, NULL, &count),
               "katran_lb_get_reals_for_vip(count)");

  if (count > 0) {
    katran_new_real_t* reals = calloc(count, sizeof(*reals));
    if (!reals) {
      fprintf(stderr, "calloc failed\n");
      katran_lb_destroy(lb);
      return 1;
    }

    die_on_error(katran_lb_get_reals_for_vip(lb, &vip, reals, &count),
                 "katran_lb_get_reals_for_vip(list)");

    for (size_t i = 0; i < count; i++) {
      printf("real[%zu]: %s weight=%u flags=%u\n",
             i,
             reals[i].address ? reals[i].address : "<null>",
             reals[i].weight,
             reals[i].flags);
    }

    katran_free_reals(reals, count);
  }

  die_on_error(katran_lb_destroy(lb), "katran_lb_destroy");
  printf("katran capi smoke test OK\n");
  return 0;
}
