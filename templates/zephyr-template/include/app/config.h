#pragma once

#include <app/collector.h>

/**
 * @brief Top-level application and collector configuration image.
 */
struct app_config {
	uint32_t schema;                                                /**< Schema version identifier (1). */
	uint32_t version;                                               /**< Monotonically increasing config version. */
	uint32_t count;                                                 /**< Number of active collector configs. */
	struct collector_config collectors[CONFIG_APP_MAX_COLLECTORS];  /**< Collector parameter array. */
};

/**
 * @brief Initialize configuration manager, load persisted configuration or defaults.
 *
 * Boot-thread only. Validates active config before startup.
 *
 * @return 0 on success, or negative errno on validation failure.
 */
int config_manager_init(void);

/**
 * @brief Retrieve a thread-safe copy of the active system configuration.
 *
 * @param[out] out Pointer to buffer receiving current configuration snapshot.
 */
void config_manager_snapshot(struct app_config *out);

/**
 * @brief Validate integrity and parameter bounds of a configuration image.
 *
 * Checks schema version, valid collector names, non-zero intervals, and bounds.
 *
 * @param config Pointer to configuration struct to validate.
 * @return 0 if valid, -EINVAL otherwise.
 */
int config_validate(const struct app_config *config);

/**
 * @brief Atomically apply and persist a new configuration image.
 *
 * Rejects stale versions (version <= current). Persists to non-volatile storage
 * before reconfiguring collectors and emitting APP_EVENT_CONFIG_UPDATED.
 *
 * @param config Pointer to desired new configuration.
 * @return 0 on success, -ESTALE if version is outdated, or negative errno on error.
 */
int config_manager_apply(const struct app_config *config);

/**
 * @brief Parse and apply a JSON remote configuration payload received from cloud.
 *
 * Validates version monotonicity, parses partial or complete collector configurations,
 * and delegates to config_manager_apply.
 *
 * @param json Null-terminated JSON string buffer (may be modified during parsing).
 * @param length Byte length of JSON payload.
 * @return 0 on success, negative errno on parse or validation error.
 */
int remote_config_apply(char *json, size_t length);

/**
 * @brief Weak persistence adapter function for non-volatile storage.
 *
 * Default implementation saves to Zephyr Settings subsystem; can be overridden in tests.
 *
 * @param config Configuration snapshot to save.
 * @return 0 on success, negative errno on flash write error.
 */
int config_storage_save(const struct app_config *config);

