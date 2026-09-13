#pragma once

#include <zephyr/kernel.h>
#include <stdint.h>

#define TELEMETRY_FIELDS 6
#define APP_NAME_SIZE    24

/**
 * @brief Fixed-point numeric telemetry data field.
 *
 * Stored as integer value with fractional decimal digits (value / 10^decimals)
 * to avoid floating point calculations and non-deterministic JSON serialization.
 */
struct telemetry_field {
	const char *name;   /**< Static field name string (e.g. "temperature_c"). */
	int64_t value;      /**< Fixed-point integer value. */
	uint8_t decimals;   /**< Number of fractional decimal digits (0-6). */
};

/**
 * @brief Discrete point-in-time sensor sampling record.
 */
struct telemetry_record {
	uint64_t sequence;                                  /**< Monotonically increasing sequence ID. */
	int64_t timestamp_ms;                               /**< UTC timestamp in ms (or uptime if unsynced). */
	int64_t uptime_ms;                                  /**< System uptime in ms at sampling moment. */
	bool time_synced;                                   /**< true if timestamp is synchronized UTC epoch. */
	uint8_t field_count;                                /**< Number of populated fields (0 to TELEMETRY_FIELDS). */
	struct telemetry_field fields[TELEMETRY_FIELDS];     /**< Array of sampled fields. */
};

/**
 * @brief Thread-safe ring buffer cache store for telemetry records.
 */
struct telemetry_store {
	struct k_mutex lock;                                /**< Store access lock. */
	struct telemetry_record records[CONFIG_APP_CACHE_SIZE]; /**< Statically allocated ring buffer slots. */
	size_t head;                                        /**< Oldest record index. */
	size_t count;                                       /**< Current number of stored records. */
	size_t capacity;                                    /**< Current configured maximum capacity. */
	uint64_t next_sequence;                             /**< Sequence number counter for next append. */
	uint64_t dropped;                                   /**< Cumulative dropped record count on overflow. */
};

/**
 * @brief Initialize a telemetry ring buffer store with given initial capacity.
 *
 * @param store Pointer to telemetry store struct.
 * @param capacity Maximum capacity (clamped to 1..CONFIG_APP_CACHE_SIZE).
 */
void telemetry_store_init(struct telemetry_store *store, size_t capacity);

/**
 * @brief Append a new record to the ring buffer, overwriting the oldest if full.
 *
 * Thread-safe. Assigns monotonic sequence number.
 *
 * @param store Pointer to telemetry store.
 * @param record Pointer to source telemetry record (copied by value).
 * @return 0 on success, -EINVAL if record validation fails.
 */
int telemetry_store_append(struct telemetry_store *store, const struct telemetry_record *record);

/**
 * @brief Copy up to `max` oldest records from store without removing them.
 *
 * Thread-safe.
 *
 * @param store Pointer to telemetry store.
 * @param[out] out Destination record buffer array.
 * @param max Maximum number of records to copy.
 * @return Number of records copied.
 */
size_t telemetry_store_peek(struct telemetry_store *store, struct telemetry_record *out,
			    size_t max);

/**
 * @brief Discard all records up to and including the acknowledged sequence number.
 *
 * Thread-safe. Guarantees that newer records appended concurrently are never dropped.
 *
 * @param store Pointer to telemetry store.
 * @param acknowledged Sequence number of last successfully uploaded record.
 */
void telemetry_store_pop(struct telemetry_store *store, uint64_t acknowledged);

/**
 * @brief Get the current count of records in the store.
 *
 * @param store Pointer to telemetry store.
 * @return Number of unacknowledged records.
 */
size_t telemetry_store_count(struct telemetry_store *store);

/**
 * @brief Get the total number of records dropped due to buffer overflow.
 *
 * @param store Pointer to telemetry store.
 * @return Dropped record count.
 */
uint64_t telemetry_store_dropped(struct telemetry_store *store);

/**
 * @brief Clear all records from the store.
 *
 * @param store Pointer to telemetry store.
 */
void telemetry_store_clear(struct telemetry_store *store);

/**
 * @brief Dynamically resize the store capacity.
 *
 * Drops oldest records if current count exceeds new capacity.
 *
 * @param store Pointer to telemetry store.
 * @param capacity Desired capacity (1..CONFIG_APP_CACHE_SIZE).
 * @return 0 on success, -EINVAL on invalid capacity.
 */
int telemetry_store_resize(struct telemetry_store *store, size_t capacity);

/**
 * @brief Serialize telemetry batch into ThingsBoard-compatible JSON format.
 *
 * Pure, re-entrant function. Returns length or negative errno; never truncates.
 *
 * @param device Device identity string.
 * @param collector Collector name string.
 * @param records Array of telemetry records.
 * @param count Number of records in batch.
 * @param[out] out Output buffer.
 * @param size Capacity of output buffer.
 * @return Length of serialized JSON on success, or -ENOSPC / -EINVAL on failure.
 */
int telemetry_serialize(const char *device, const char *collector,
			const struct telemetry_record *records, size_t count, char *out,
			size_t size);

/**
 * @brief Format a standardized MQTT topic string for the device: "devices/{id}/{suffix}".
 *
 * @param device Device identity string.
 * @param suffix Topic suffix (e.g. "telemetry", "state").
 * @param[out] out Destination output buffer.
 * @param size Output buffer capacity.
 * @return Topic string length on success, or negative errno on error.
 */
int telemetry_topic(const char *device, const char *suffix, char *out, size_t size);

/**
 * @brief Validate that an identifier consists of alphanumeric characters, hyphens, or underscores.
 *
 * @param name Null-terminated name string.
 * @return true if valid and non-empty, false otherwise.
 */
bool app_name_valid(const char *name);

