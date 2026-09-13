#pragma once

#include <zephyr/kernel.h>

/**
 * @brief Network connectivity states.
 */
enum network_state {
	NETWORK_DISCONNECTED,   /**< Interface down or idle. */
	NETWORK_CONNECTING,     /**< Association / DHCP handshake in progress. */
	NETWORK_CONNECTED,      /**< IP address acquired and interface up. */
	NETWORK_RECONNECTING,   /**< Link dropped, attempting exponential backoff reconnect. */
	NETWORK_FAILED          /**< Permanent failure / invalid credentials. */
};

/**
 * @brief Initialize Wi-Fi subsystem and register management event callbacks.
 *
 * @return 0 on success, or negative errno on error.
 */
int wifi_manager_init(void);

/**
 * @brief Request Wi-Fi station connection.
 *
 * @return 0 on success, -ENODEV if not initialized.
 */
int wifi_manager_connect(void);

/**
 * @brief Disconnect Wi-Fi and cancel pending reconnection attempts.
 *
 * @return 0 on success, negative errno on error.
 */
int wifi_manager_disconnect(void);

/**
 * @brief Check if Wi-Fi interface is fully connected with an assigned IPv4 address.
 *
 * @return true if connected and valid IP present, false otherwise.
 */
bool wifi_manager_is_connected(void);

/**
 * @brief Block until Wi-Fi connection is established or timeout expires.
 *
 * @param timeout Maximum wait timeout.
 * @return 0 if connected before timeout, -ETIMEDOUT otherwise.
 */
int wifi_manager_wait_connected(k_timeout_t timeout);

/**
 * @brief Retrieve current Wi-Fi state.
 *
 * @return Current network_state enum value.
 */
enum network_state wifi_manager_state(void);

/**
 * @brief Get cumulative Wi-Fi reconnection count.
 *
 * @return Number of reconnect attempts.
 */
uint32_t wifi_manager_reconnects(void);

/**
 * @brief Initialize MQTT manager actor thread, queues, and TLS credentials if enabled.
 *
 * @return 0 on success, negative errno on error.
 */
int mqtt_manager_init(void);

/**
 * @brief Start the MQTT client loop and initiate cloud broker connection.
 *
 * @return 0 on success.
 */
int mqtt_manager_start(void);

/**
 * @brief Stop the MQTT client loop and disconnect from broker.
 */
void mqtt_manager_stop(void);

/**
 * @brief Check if MQTT session is actively connected and authenticated.
 *
 * @return true if MQTT_CONNECTED, false otherwise.
 */
bool mqtt_manager_is_connected(void);

/**
 * @brief Publish a message to a device sub-topic.
 *
 * Thread-safe. Synchronously waits for socket send (QoS0) or PUBACK (QoS1).
 *
 * @param suffix Topic suffix (e.g., "telemetry", "event", "state").
 * @param payload Message payload string buffer.
 * @param length Byte length of payload.
 * @param qos MQTT QoS level (0 or 1).
 * @param retain MQTT retain flag.
 * @return 0 on confirmed delivery, negative errno on failure or timeout.
 */
int mqtt_manager_publish(const char *suffix, const char *payload, size_t length, int qos,
			 bool retain);

/**
 * @brief Subscribe to a device sub-topic (e.g., "config", "command").
 *
 * Thread-safe. Synchronously waits for SUBACK.
 *
 * @param suffix Topic suffix to subscribe to.
 * @return 0 on confirmed subscription, negative errno on failure.
 */
int mqtt_manager_subscribe(const char *suffix);

/**
 * @brief Get cumulative MQTT reconnection count.
 *
 * @return Number of reconnect cycles.
 */
uint32_t mqtt_manager_reconnects(void);

