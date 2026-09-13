#pragma once
#include <app/telemetry.h>
#include <zephyr/sys/iterable_sections.h>
struct collector_config {
	char name[APP_NAME_SIZE];
	int32_t collection_interval_sec, upload_interval_sec;
	int32_t cache_size, retry, priority;
	bool enabled;
};
struct collector;
struct collector_ops {
	int (*init)(void);
	int (*start)(void);
	int (*stop)(void);
	int (*collect)(struct telemetry_record *record);
};
struct collector {
	const char *name;
	const struct collector_ops *ops;
	int32_t collection_interval_sec, upload_interval_sec;
};
#define COLLECTOR_DEFINE(id, operations, collect_sec, upload_sec)                                  \
	STRUCT_SECTION_ITERABLE(collector, id) = {.name = #id,                                     \
						  .ops = operations,                               \
						  .collection_interval_sec = collect_sec,          \
						  .upload_interval_sec = upload_sec}
enum collector_state {
	COLLECTOR_DISABLED,
	COLLECTOR_IDLE,
	COLLECTOR_COLLECTING,
	COLLECTOR_ERROR
};
struct collector_runtime {
	const struct collector *definition;
	struct collector_config config;
	struct telemetry_store store;
	struct k_mutex lock;
	struct k_work_delayable collect_work, upload_work;
	enum collector_state state;
	uint64_t errors;
	int64_t next_collect, next_upload;
	uint8_t failures;
	bool initialized, started;
};
/* Init/start are boot-thread only; registry is immutable after init. */
int collector_manager_init(void);
void collector_manager_start(void);
size_t collector_count(void);
struct collector_runtime *collector_at(size_t index);
int collector_find(const char *name);
/* Thread-safe snapshot/apply. Driver callbacks execute only on collector workqueue. */
void collector_snapshot(size_t index, struct collector_config *config, enum collector_state *state,
			uint64_t *errors);
void collector_configure(size_t index, const struct collector_config *config);
/* Pure deadline helper, skips missed periods rather than producing a burst. */
int64_t collector_next_deadline(int64_t previous, int64_t now, int32_t interval_sec);
