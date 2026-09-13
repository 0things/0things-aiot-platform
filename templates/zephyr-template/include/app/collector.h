#pragma once

#include <app/telemetry.h>
#include <zephyr/sys/iterable_sections.h>

/**
 * @brief Collector configuration parameters.
 */
struct collector_config {
	char name[APP_NAME_SIZE];               /**< Unique collector identifier name. */
	int32_t collection_interval_sec;        /**< Sampling period in seconds. */
	int32_t upload_interval_sec;            /**< Batch upload period in seconds. */
	int32_t cache_size;                     /**< Max ring buffer record capacity. */
	int32_t retry;                          /**< Max consecutive retry count on failure. */
	int32_t priority;                       /**< Upload priority ranking (0-15). */
	bool enabled;                           /**< Administrative enable/disable switch. */
};

struct collector;

/**
 * @brief Sensor driver operation callbacks for a collector.
 */
struct collector_ops {
	int (*init)(void);                                  /**< One-time driver hardware initialization. */
	int (*start)(void);                                 /**< Start active sensor acquisition. */
	int (*stop)(void);                                  /**< Stop sensor acquisition to conserve power. */
	int (*collect)(struct telemetry_record *record);    /**< Poll current sample into telemetry record. */
};

/**
 * @brief Static descriptor declared at compile-time via COLLECTOR_DEFINE.
 */
struct collector {
	const char *name;                       /**< Collector name string. */
	const struct collector_ops *ops;        /**< Driver function pointers. */
	int32_t collection_interval_sec;        /**< Default collection period in seconds. */
	int32_t upload_interval_sec;            /**< Default upload period in seconds. */
};

/**
 * @brief Macro to define and register a collector into the iterable linker section.
 */
#define COLLECTOR_DEFINE(id, operations, collect_sec, upload_sec)                                  \
	STRUCT_SECTION_ITERABLE(collector, id) = {.name = #id,                                     \
						  .ops = operations,                               \
						  .collection_interval_sec = collect_sec,          \
						  .upload_interval_sec = upload_sec}

/**
 * @brief Runtime operational state of a collector.
 */
enum collector_state {
	COLLECTOR_DISABLED,  /**< Administratively disabled. */
	COLLECTOR_IDLE,      /**< Waiting for next sample deadline. */
	COLLECTOR_COLLECTING,/**< Executing driver sample acquisition. */
	COLLECTOR_ERROR      /**< Last acquisition failed. */
};

/**
 * @brief Dynamic runtime state associated with an active collector instance.
 */
struct collector_runtime {
	const struct collector *definition;     /**< Pointer to immutable static definition. */
	struct collector_config config;         /**< Current runtime configuration. */
	struct telemetry_store store;           /**< Ring buffer telemetry store. */
	struct k_mutex lock;                    /**< State protection mutex. */
	struct k_work_delayable collect_work;   /**< Scheduled sampling work item. */
	struct k_work_delayable upload_work;    /**< Scheduled upload trigger work item. */
	enum collector_state state;             /**< Current lifecycle state. */
	uint64_t errors;                        /**< Cumulative acquisition error counter. */
	int64_t next_collect;                   /**< Next sampling deadline timestamp (ms). */
	int64_t next_upload;                    /**< Next upload deadline timestamp (ms). */
	uint8_t failures;                       /**< Consecutive failure counter. */
	bool initialized;                       /**< Driver init status. */
	bool started;                           /**< Driver start status. */
};

/**
 * @brief Discover and initialize all collectors declared in the iterable section.
 *
 * Boot-thread only; registry becomes immutable after initialization.
 *
 * @return 0 on success, or negative errno on error.
 */
int collector_manager_init(void);

/**
 * @brief Start background collection and upload timers for all enabled collectors.
 */
void collector_manager_start(void);

/**
 * @brief Get the total number of registered collectors.
 *
 * @return Total collector count.
 */
size_t collector_count(void);

/**
 * @brief Access collector runtime by array index.
 *
 * @param index 0-based collector index.
 * @return Pointer to collector runtime, or NULL if out of range.
 */
struct collector_runtime *collector_at(size_t index);

/**
 * @brief Find collector index by unique name.
 *
 * @param name Collector identifier string.
 * @return 0-based index if found, or -ENOENT if not found.
 */
int collector_find(const char *name);

/**
 * @brief Take a thread-safe snapshot of a collector's status and metrics.
 *
 * @param index 0-based collector index.
 * @param[out] config Optional buffer to receive current configuration.
 * @param[out] state Optional pointer to receive current state.
 * @param[out] errors Optional pointer to receive cumulative error count.
 */
void collector_snapshot(size_t index, struct collector_config *config, enum collector_state *state,
			uint64_t *errors);

/**
 * @brief Reconfigure a collector dynamically at runtime.
 *
 * Thread-safe. Resizes cache store, updates intervals, and reschedules timers.
 *
 * @param index 0-based collector index.
 * @param config Pointer to desired new configuration.
 */
void collector_configure(size_t index, const struct collector_config *config);

/**
 * @brief Calculate the next execution deadline, skipping missed periods on overrun.
 *
 * Pure function. Ensures non-bursting predictable schedule.
 *
 * @param previous Previous scheduled deadline timestamp (ms).
 * @param now Current uptime timestamp (ms).
 * @param interval_sec Periodic interval in seconds.
 * @return Next monotonic deadline timestamp (ms).
 */
int64_t collector_next_deadline(int64_t previous, int64_t now, int32_t interval_sec);

