#pragma once
#include <app/collector.h>
#include <zephyr/zbus/zbus.h>
enum app_event_type {
	APP_EVENT_NETWORK_UP,
	APP_EVENT_NETWORK_DOWN,
	APP_EVENT_MQTT_UP,
	APP_EVENT_MQTT_DOWN,
	APP_EVENT_TIME_SYNCED,
	APP_EVENT_CONFIG_UPDATED
};
struct app_event {
	enum app_event_type type;
};
ZBUS_CHAN_DECLARE(app_events);
/* Thread-safe, nonblocking notification. State consumers reconcile independently. */
void app_event_emit(enum app_event_type type);
/* Boot-only init; immutable identifier afterwards. */
int device_identity_init(void);
const char *device_identity_get(void);
void time_manager_init(void);
bool time_manager_is_synced(void);
int64_t time_manager_now_ms(void);
void time_manager_stamp(struct telemetry_record *record);
/* Thread-safe request coalescing; no caller-owned data retained. */
void telemetry_uploader_init(void);
void telemetry_uploader_request(size_t index);
void telemetry_uploader_flush(void);
void health_manager_start(void);
void health_upload_error(void);
uint32_t health_upload_errors(void);
/* Extension hooks return -ENOTSUP until a product backend is provided. */
int ota_manager_request(const char *url);
void watchdog_manager_heartbeat(unsigned int subsystem);
bool watchdog_manager_healthy(int64_t maximum_age_ms);
