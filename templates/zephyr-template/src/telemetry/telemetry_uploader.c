#include <app/services.h>
#include <app/network.h>
#include <zephyr/sys/atomic.h>
#include <stdio.h>

static atomic_t pending[CONFIG_APP_MAX_COLLECTORS];
static struct k_work_delayable upload_work;
static struct k_work_q upload_queue;
K_THREAD_STACK_DEFINE(upload_stack, 6144);
static size_t cursor;
static struct k_spinlock schedule_lock;
static int64_t wake_at = INT64_MAX;
static int64_t allowed_at;

/**
 * @brief Schedule or advance the next telemetry upload work execution.
 *
 * @param delay_ms Minimum milliseconds from now before running the upload.
 */
static void schedule_upload(int32_t delay_ms)
{
	k_spinlock_key_t key = k_spin_lock(&schedule_lock);
	int64_t deadline = MAX(k_uptime_get() + delay_ms, allowed_at);
	if (deadline < wake_at) {
		wake_at = deadline;
		k_work_reschedule_for_queue(&upload_queue, &upload_work,
					    K_MSEC(MAX(1, deadline - k_uptime_get())));
	}
	k_spin_unlock(&schedule_lock, key);
}

/**
 * @brief Check if any collector has queued records pending transmission.
 */
static bool has_pending(void)
{
	for (size_t i = 0; i < collector_count(); i++) {
		if (atomic_get(&pending[i])) {
			return true;
		}
	}
	return false;
}

static bool was_connected;
static uint16_t age[CONFIG_APP_MAX_COLLECTORS];

/**
 * @brief Work queue handler performing prioritized, round-robin batch telemetry uploads.
 */
static void upload_handler(struct k_work *work)
{
	ARG_UNUSED(work);
	k_spinlock_key_t key = k_spin_lock(&schedule_lock);
	wake_at = INT64_MAX;
	allowed_at = k_uptime_get() + CONFIG_APP_UPLOAD_GAP_MS;
	k_spin_unlock(&schedule_lock, key);
	watchdog_manager_heartbeat(1);
	if (!mqtt_manager_is_connected()) {
		was_connected = false;
		goto again;
	}
	if (!was_connected) {
		telemetry_uploader_flush();
		was_connected = true;
	}
	int chosen = -1, best = -1;
	/* Round-robin among equal priorities; bounded one-batch turn prevents reconnect floods. */
	for (size_t n = 0; n < collector_count(); n++) {
		size_t i = (cursor + n) % collector_count();
		if (!atomic_get(&pending[i])) {
			continue;
		}
		struct collector_config cfg;
		collector_snapshot(i, &cfg, NULL, NULL);
		if (!cfg.enabled || !telemetry_store_count(&collector_at(i)->store)) {
			atomic_val_t generation = atomic_get(&pending[i]);
			atomic_cas(&pending[i], generation, 0);
			continue;
		}
		int score = cfg.priority + age[i];
		if (score > best) {
			best = score;
			chosen = i;
		}
		if (age[i] < 1000) {
			age[i]++;
		}
	}
	if (chosen >= 0) {
		atomic_val_t generation = atomic_get(&pending[chosen]);
		age[chosen] = 0;
		struct collector_runtime *r = collector_at(chosen);
		struct telemetry_record batch[CONFIG_APP_BATCH_SIZE];
		char payload[CONFIG_APP_PAYLOAD_SIZE];
		size_t n = telemetry_store_peek(&r->store, batch, ARRAY_SIZE(batch));
		int rc = telemetry_serialize(device_identity_get(), r->definition->name, batch, n,
					     payload, sizeof(payload));
		/* Fit a smaller batch instead of permanently blocking a collector on payload size.
		 */
		while (rc == -ENOSPC && n > 1) {
			rc = telemetry_serialize(device_identity_get(), r->definition->name, batch,
						 --n, payload, sizeof(payload));
		}
		if (rc > 0) {
			rc = mqtt_manager_publish("telemetry", payload, rc,
						  CONFIG_APP_TELEMETRY_QOS, false);
		}
		if (!rc) {
			telemetry_store_pop(&r->store, batch[n - 1].sequence);
			if (!telemetry_store_count(&r->store)) {
				atomic_cas(&pending[chosen], generation, 0);
			}
		} else {
			health_upload_error();
		}
		cursor = (chosen + 1) % collector_count();
	}
again:
	/* Slow reconciliation also recovers from a dropped state notification. */
	schedule_upload(mqtt_manager_is_connected() && has_pending() ? CONFIG_APP_UPLOAD_GAP_MS
								     : 30000);
}

void telemetry_uploader_request(size_t index)
{
	if (index < collector_count()) {
		atomic_inc(&pending[index]);
		schedule_upload(0);
	}
}

void telemetry_uploader_flush(void)
{
	for (size_t i = 0; i < collector_count(); i++) {
		telemetry_uploader_request(i);
	}
}

/**
 * @brief Zbus event handler triggering immediate upload flush on MQTT connection up.
 */
static void mqtt_event(const struct zbus_channel *channel)
{
	const struct app_event *e = zbus_chan_const_msg(channel);
	if (e->type == APP_EVENT_MQTT_UP) {
		telemetry_uploader_flush();
	}
}

ZBUS_LISTENER_DEFINE(upload_listener, mqtt_event);
ZBUS_CHAN_ADD_OBS(app_events, upload_listener, 2);

void telemetry_uploader_init(void)
{
	k_work_init_delayable(&upload_work, upload_handler);
	k_work_queue_start(&upload_queue, upload_stack, K_THREAD_STACK_SIZEOF(upload_stack), 8,
			   NULL);
	k_thread_name_set(k_work_queue_thread_get(&upload_queue), "uploader");
	schedule_upload(CONFIG_APP_UPLOAD_GAP_MS);
}
