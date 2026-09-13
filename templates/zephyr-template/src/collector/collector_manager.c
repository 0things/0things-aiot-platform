#include <app/collector.h>
#include <app/services.h>
#include <zephyr/logging/log.h>
#include <string.h>
LOG_MODULE_REGISTER(collector_manager, CONFIG_APP_LOG_LEVEL);

static struct collector_runtime runtimes[CONFIG_APP_MAX_COLLECTORS];
static size_t registered;
static struct k_work_q collect_queue;
K_THREAD_STACK_DEFINE(collect_stack, 4096);

/**
 * @brief Calculate next execution deadline to prevent drift and phase alignment.
 *
 * @param previous Previous scheduled timestamp (ms).
 * @param now Current uptime timestamp (ms).
 * @param interval Recurrence interval in seconds.
 * @return Next execution deadline in milliseconds.
 */
int64_t collector_next_deadline(int64_t previous, int64_t now, int32_t interval)
{
	int64_t step = (int64_t)interval * 1000;
	if (interval <= 0) {
		return INT64_MAX;
	}
	return previous > now ? previous : previous + ((now - previous) / step + 1) * step;
}

/**
 * @brief Work queue handler executing periodic data collection for a sensor.
 */
static void collect_handler(struct k_work *work)
{
	struct collector_runtime *r = CONTAINER_OF(k_work_delayable_from_work(work),
						   struct collector_runtime, collect_work);
	k_mutex_lock(&r->lock, K_FOREVER);
	if (!r->config.enabled) {
		if (r->started && r->definition->ops->stop) {
			r->definition->ops->stop();
		}
		r->started = false;
		r->state = COLLECTOR_DISABLED;
		k_mutex_unlock(&r->lock);
		return;
	}
	int rc = 0;
	if (!r->initialized) {
		rc = r->definition->ops->init ? r->definition->ops->init() : 0;
		r->initialized = !rc;
	}
	if (!rc && !r->started) {
		rc = r->definition->ops->start ? r->definition->ops->start() : 0;
		r->started = !rc;
	}
	int64_t now = k_uptime_get();
	if (!rc && now >= r->next_collect) {
		struct telemetry_record record = {0};
		r->state = COLLECTOR_COLLECTING;
		rc = r->definition->ops->collect(&record);
		if (!rc) {
			time_manager_stamp(&record);
			rc = telemetry_store_append(&r->store, &record);
		}
		r->next_collect = collector_next_deadline(r->next_collect, k_uptime_get(),
							  r->config.collection_interval_sec);
	}
	if (rc) {
		r->errors++;
		r->state = COLLECTOR_ERROR;
		LOG_WRN("%s collect failed: %d", r->definition->name, rc);
		if (r->failures++ < r->config.retry) {
			/* Fast retry on transient collection failure */
			r->next_collect = k_uptime_get() + 1000;
		} else {
			/* Exceeded retries, fall back to normal period */
			r->failures = 0;
			r->next_collect =
				k_uptime_get() + (int64_t)r->config.collection_interval_sec * 1000;
		}
	} else {
		r->failures = 0;
		r->state = COLLECTOR_IDLE;
	}
	k_work_reschedule_for_queue(&collect_queue, &r->collect_work,
				    K_MSEC(MAX(1, r->next_collect - k_uptime_get())));
	k_mutex_unlock(&r->lock);
	watchdog_manager_heartbeat(0);
}

/**
 * @brief Work queue handler triggering batch upload requests for a collector.
 */
static void upload_handler(struct k_work *work)
{
	struct collector_runtime *r = CONTAINER_OF(k_work_delayable_from_work(work),
						   struct collector_runtime, upload_work);
	k_mutex_lock(&r->lock, K_FOREVER);
	if (r->config.enabled) {
		/* At a shared deadline, let the due sample finish before requesting its batch. */
		if (r->next_collect <= k_uptime_get()) {
			k_work_reschedule(&r->upload_work, K_MSEC(10));
			k_mutex_unlock(&r->lock);
			return;
		}
		telemetry_uploader_request(r - runtimes);
		r->next_upload = collector_next_deadline(r->next_upload, k_uptime_get(),
							 r->config.upload_interval_sec);
		k_work_reschedule(&r->upload_work, K_MSEC(MAX(1, r->next_upload - k_uptime_get())));
	}
	k_mutex_unlock(&r->lock);
}

int collector_manager_init(void)
{
	STRUCT_SECTION_FOREACH(collector, c) {
		if (registered == ARRAY_SIZE(runtimes)) {
			return -ENOMEM;
		}
		if (!app_name_valid(c->name) || !c->ops || !c->ops->collect ||
		    collector_find(c->name) >= 0) {
			return -EINVAL;
		}
		struct collector_runtime *r = &runtimes[registered++];
		r->definition = c;
		k_mutex_init(&r->lock);
		r->config = (struct collector_config){.enabled = true,
						      .collection_interval_sec =
							      c->collection_interval_sec,
						      .upload_interval_sec = c->upload_interval_sec,
						      .cache_size = CONFIG_APP_CACHE_SIZE,
						      .retry = 2,
						      .priority = 5};
		strcpy(r->config.name, c->name);
		telemetry_store_init(&r->store, r->config.cache_size);
		k_work_init_delayable(&r->collect_work, collect_handler);
		k_work_init_delayable(&r->upload_work, upload_handler);
	}
	k_work_queue_start(&collect_queue, collect_stack, K_THREAD_STACK_SIZEOF(collect_stack), 7,
			   NULL);
	k_thread_name_set(k_work_queue_thread_get(&collect_queue), "collectors");
	return registered ? 0 : -ENODEV;
}

void collector_manager_start(void)
{
	for (size_t i = 0; i < registered; i++) {
		collector_configure(i, &runtimes[i].config);
	}
}

size_t collector_count(void)
{
	return registered;
}

struct collector_runtime *collector_at(size_t index)
{
	return index < registered ? &runtimes[index] : NULL;
}

int collector_find(const char *name)
{
	for (size_t i = 0; i < registered; i++) {
		if (!strcmp(runtimes[i].definition->name, name)) {
			return i;
		}
	}
	return -ENOENT;
}

void collector_snapshot(size_t index, struct collector_config *cfg, enum collector_state *state,
			uint64_t *errors)
{
	struct collector_runtime *r = collector_at(index);
	k_mutex_lock(&r->lock, K_FOREVER);
	if (cfg) {
		*cfg = r->config;
	}
	if (state) {
		*state = r->state;
	}
	if (errors) {
		*errors = r->errors;
	}
	k_mutex_unlock(&r->lock);
}

void collector_configure(size_t index, const struct collector_config *cfg)
{
	struct collector_runtime *r = collector_at(index);
	k_mutex_lock(&r->lock, K_FOREVER);
	r->config = *cfg;
	telemetry_store_resize(&r->store, cfg->cache_size);
	r->next_collect = k_uptime_get() + (int64_t)cfg->collection_interval_sec * 1000;
	r->next_upload = k_uptime_get() + (int64_t)cfg->upload_interval_sec * 1000;
	/* Schedule lifecycle on the same queue as driver calls, including disable/stop. */
	k_work_reschedule_for_queue(&collect_queue, &r->collect_work, K_NO_WAIT);
	if (cfg->enabled) {
		k_work_reschedule(&r->upload_work, K_SECONDS(cfg->upload_interval_sec));
	} else {
		k_work_cancel_delayable(&r->upload_work);
	}
	k_mutex_unlock(&r->lock);
}
