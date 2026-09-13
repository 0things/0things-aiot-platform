#pragma once
#include <zephyr/kernel.h>
enum network_state {
	NETWORK_DISCONNECTED,
	NETWORK_CONNECTING,
	NETWORK_CONNECTED,
	NETWORK_RECONNECTING,
	NETWORK_FAILED
};
int wifi_manager_init(void);
int wifi_manager_connect(void);
int wifi_manager_disconnect(void);
bool wifi_manager_is_connected(void);
int wifi_manager_wait_connected(k_timeout_t timeout);
enum network_state wifi_manager_state(void);
uint32_t wifi_manager_reconnects(void);
/* Boot init/start once. Publish/subscribe/stop/state are thread-safe, thread context.
 * Publish copies topic/payload; returns 0 only after socket send (QoS0) / PUBACK (QoS1).
 * A timeout is ambiguous delivery; callers retain data and may retry (at least once).
 * Never call blocking publish/subscribe from MQTT callbacks. */
int mqtt_manager_init(void);
int mqtt_manager_start(void);
void mqtt_manager_stop(void);
bool mqtt_manager_is_connected(void);
int mqtt_manager_publish(const char *suffix, const char *payload, size_t length, int qos,
			 bool retain);
int mqtt_manager_subscribe(const char *suffix);
uint32_t mqtt_manager_reconnects(void);
