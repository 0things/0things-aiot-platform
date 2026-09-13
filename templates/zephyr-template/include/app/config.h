#pragma once
#include <app/collector.h>
struct app_config {
	uint32_t schema, version, count;
	struct collector_config collectors[CONFIG_APP_MAX_COLLECTORS];
};
/* Thread context. Init boot-only; other calls serialized, snapshot copies data.
 * Apply validates all entries and persists a single snapshot BEFORE publishing it.
 * Failure leaves running config/version unchanged. */
int config_manager_init(void);
void config_manager_snapshot(struct app_config *out);
int config_validate(const struct app_config *config);
int config_manager_apply(const struct app_config *config);
/* Mutable JSON, not retained. Atomic partial updates; version must increase. */
int remote_config_apply(char *json, size_t length);
/* Weak persistence adapter, overridden by tests. Returns 0 or negative errno. */
int config_storage_save(const struct app_config *config);
