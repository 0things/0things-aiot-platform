#pragma once
#include <zephyr/kernel.h>
#include <stdint.h>
#define TELEMETRY_FIELDS 6
#define APP_NAME_SIZE    24
/* Values are fixed-point integers: value / 10^decimals, avoiding float JSON. */
struct telemetry_field {
	const char *name;
	int64_t value;
	uint8_t decimals;
};
struct telemetry_record {
	uint64_t sequence;
	int64_t timestamp_ms;
	int64_t uptime_ms;
	bool time_synced;
	uint8_t field_count;
	struct telemetry_field fields[TELEMETRY_FIELDS];
};
struct telemetry_store {
	struct k_mutex lock;
	struct telemetry_record records[CONFIG_APP_CACHE_SIZE];
	size_t head, count, capacity;
	uint64_t next_sequence, dropped;
};
/* Initialize once before sharing. All other store calls are thread-safe, thread context.
 * Records copy by value; field name strings must have static lifetime. */
void telemetry_store_init(struct telemetry_store *store, size_t capacity);
int telemetry_store_append(struct telemetry_store *store, const struct telemetry_record *record);
size_t telemetry_store_peek(struct telemetry_store *store, struct telemetry_record *out,
			    size_t max);
/* Remove only sequence <= acknowledged: concurrent overflow/appends cannot delete new data. */
void telemetry_store_pop(struct telemetry_store *store, uint64_t acknowledged);
size_t telemetry_store_count(struct telemetry_store *store);
uint64_t telemetry_store_dropped(struct telemetry_store *store);
void telemetry_store_clear(struct telemetry_store *store);
int telemetry_store_resize(struct telemetry_store *store, size_t capacity);
/* Pure/reentrant. Caller owns buffers. Returns length or negative errno; never truncates. */
int telemetry_serialize(const char *device, const char *collector,
			const struct telemetry_record *records, size_t count, char *out,
			size_t size);
int telemetry_topic(const char *device, const char *suffix, char *out, size_t size);
bool app_name_valid(const char *name);
