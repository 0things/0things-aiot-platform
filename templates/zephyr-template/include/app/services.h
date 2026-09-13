#pragma once

#include <app/collector.h>
#include <zephyr/zbus/zbus.h>

/**
 * @brief System-wide broadcast event types.
 */
enum app_event_type {
	APP_EVENT_NETWORK_UP,       /**< Network interface acquired IP. */
	APP_EVENT_NETWORK_DOWN,     /**< Network interface disconnected. */
	APP_EVENT_MQTT_UP,          /**< MQTT session connected and authenticated. */
	APP_EVENT_MQTT_DOWN,        /**< MQTT session disconnected. */
	APP_EVENT_TIME_SYNCED,      /**< Wall clock synchronized via SNTP/RTC. */
	APP_EVENT_CONFIG_UPDATED    /**< Runtime configuration changed and persisted. */
};

/**
 * @brief Event message payload published over zbus.
 */
struct app_event {
	enum app_event_type type;
};

ZBUS_CHAN_DECLARE(app_events);

/**
 * @brief Emit a system-wide application event asynchronously over zbus channel.
 *
 * Thread-safe and non-blocking.
 *
 * @param type Event type enum.
 */
void app_event_emit(enum app_event_type type);

/**
 * @brief Initialize device identity from hardware ID, MAC address, Settings, or Kconfig.
 *
 * Boot-thread only. Identity is immutable after initialization.
 *
 * @return 0 on success, negative errno on failure.
 */
int device_identity_init(void);

/**
 * @brief Retrieve the null-terminated unique device identity string.
 *
 * @return Pointer to static device identity string.
 */
const char *device_identity_get(void);

/**
 * @brief Initialize the time manager subsystem, SNTP sync queue, and RTC if present.
 */
void time_manager_init(void);

/**
 * @brief Check if real-time clock is synchronized with a valid wall-clock time source.
 *
 * @return true if wall-clock synced, false if running on monotonic uptime only.
 */
bool time_manager_is_synced(void);

/**
 * @brief Get current timestamp in milliseconds (UTC unix epoch if synced, uptime otherwise).
 *
 * @return Timestamp in milliseconds.
 */
int64_t time_manager_now_ms(void);

/**
 * @brief Stamp telemetry record with uptime, sync flag, and current UTC epoch timestamp.
 *
 * @param[out] record Telemetry record to populate.
 */
void time_manager_stamp(struct telemetry_record *record);

/**
 * @brief Initialize background telemetry batch uploader worker thread and queue.
 */
void telemetry_uploader_init(void);

/**
 * @brief Request an upload cycle for a specific collector index.
 *
 * Thread-safe with request coalescing.
 *
 * @param index 0-based collector index.
 */
void telemetry_uploader_request(size_t index);

/**
 * @brief Request an immediate upload flush for all collectors.
 */
void telemetry_uploader_flush(void);

/**
 * @brief Start periodic health status report worker.
 */
void health_manager_start(void);

/**
 * @brief Record an upload transmission failure.
 */
void health_upload_error(void);

/**
 * @brief Retrieve the total count of upload transmission failures.
 *
 * @return Failure count.
 */
uint32_t health_upload_errors(void);

/**
 * @brief Request firmware update download from specified URL (weak extension hook).
 *
 * @param url Remote firmware binary URL.
 * @return 0 if accepted, -ENOTSUP if unhandled.
 */
int ota_manager_request(const char *url);

/**
 * @brief Refresh software watchdog heartbeat timestamp for a given subsystem.
 *
 * @param subsystem Subsystem index (0: collector, 1: uploader, 2: mqtt).
 */
void watchdog_manager_heartbeat(unsigned int subsystem);

/**
 * @brief Verify that all monitored subsystems refreshed heartbeats within the maximum age.
 *
 * @param maximum_age_ms Maximum tolerated inactivity time in milliseconds.
 * @return true if all subsystems are healthy, false otherwise.
 */
bool watchdog_manager_healthy(int64_t maximum_age_ms);

