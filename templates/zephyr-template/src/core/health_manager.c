#include <app/services.h>
#include <app/network.h>
#include <zephyr/sys/atomic.h>
#include <stdio.h>
#include <inttypes.h>
#include <zephyr/sys/sys_heap.h>

static atomic_t upload_errors;
static struct k_work_delayable health_work;
static struct k_work_q health_queue;
K_THREAD_STACK_DEFINE(health_stack, 3072);

/**
 * @brief Increment the global upload transmission error counter.
 */
void health_upload_error(void)
{
	atomic_inc(&upload_errors);
}

/**
 * @brief Retrieve current upload error counter value.
 *
 * @return Cumulative upload failures.
 */
uint32_t health_upload_errors(void)
{
	return atomic_get(&upload_errors);
}

#if defined(CONFIG_APP_STACK_METRICS)
/**
 * @brief Thread iterator callback to determine minimum unused stack watermark.
 */
static void measure_stack(const struct k_thread *thread, void *data)
{
	size_t unused;
	if (!k_thread_stack_space_get(thread, &unused)) {
		size_t *minimum = data;
		*minimum = MIN(*minimum, unused);
	}
}
#endif

/**
 * @brief Periodic work handler collecting and publishing device health diagnostics.
 */
static void report(struct k_work *work)
{
	ARG_UNUSED(work);
	size_t count = 0;
	uint64_t dropped = 0, errors = 0;

	/* Aggregate queue statistics across all collector ring buffers */
	for (size_t i = 0; i < collector_count(); i++) {
		count += telemetry_store_count(&collector_at(i)->store);
		dropped += telemetry_store_dropped(&collector_at(i)->store);
		uint64_t e;
		collector_snapshot(i, NULL, NULL, &e);
		errors += e;
	}

	char payload[512];
	char stack_unused[24] = "null";
	char heap_free[24] = "null";

#if defined(CONFIG_APP_HEAP_METRICS)
	struct sys_heap **heaps;
	size_t free_bytes = 0;
	/* Single-core only: prevent allocator updates during the bounded counter snapshot */
	unsigned int irq_key = irq_lock();
	int heap_count = sys_heap_array_get(&heaps);
	for (int i = 0; i < heap_count; i++) {
		struct sys_memory_stats stats;
		if (!sys_heap_runtime_stats_get(heaps[i], &stats)) {
			free_bytes += stats.free_bytes;
		}
	}
	irq_unlock(irq_key);
	if (heap_count >= 0) {
		snprintf(heap_free, sizeof(heap_free), "%u", (unsigned)free_bytes);
	}
#endif

#if defined(CONFIG_APP_STACK_METRICS)
	size_t minimum = SIZE_MAX;
	k_thread_foreach(measure_stack, &minimum);
	if (minimum != SIZE_MAX) {
		snprintf(stack_unused, sizeof(stack_unused), "%u", (unsigned)minimum);
	}
#endif

	unsigned int wifi_retries = 0;
#if defined(CONFIG_APP_WIFI)
	wifi_retries = wifi_manager_reconnects();
#endif

	/* Build JSON state diagnostics payload */
	int n = snprintf(
		payload, sizeof(payload),
		"{\"uptime_ms\":%" PRId64 ",\"mqtt\":\"%s\",\"queue_size\":%u,\"dropped\":%" PRIu64
		",\"collector_errors\":%" PRIu64 ",\"upload_errors\":%u,\"wifi_reconnects\":%u,"
		"\"mqtt_reconnects\":%u,\"free_heap\":%s,\"stack_unused_min\":%s}",
		k_uptime_get(), mqtt_manager_is_connected() ? "connected" : "disconnected",
		(unsigned)count, dropped, errors, health_upload_errors(), wifi_retries,
		mqtt_manager_reconnects(), heap_free, stack_unused);

	if (n > 0 && n < sizeof(payload) && mqtt_manager_is_connected()) {
		mqtt_manager_publish("state", payload, n, 1, true);
	}

	/* Reschedule next periodic health report */
	k_work_reschedule_for_queue(&health_queue, &health_work,
				    K_SECONDS(CONFIG_APP_HEALTH_INTERVAL));
}

/**
 * @brief Initialize and start the health reporting delayed work queue.
 */
void health_manager_start(void)
{
	k_work_init_delayable(&health_work, report);
	k_work_queue_start(&health_queue, health_stack, K_THREAD_STACK_SIZEOF(health_stack), 9,
			   NULL);
	k_work_reschedule_for_queue(&health_queue, &health_work,
				    K_SECONDS(CONFIG_APP_HEALTH_INTERVAL));
}

