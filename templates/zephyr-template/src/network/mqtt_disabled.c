#include <app/network.h>

/**
 * @brief Stub initialization when MQTT support is disabled.
 */
int mqtt_manager_init(void)
{
	return 0;
}

/**
 * @brief Stub start when MQTT support is disabled.
 */
int mqtt_manager_start(void)
{
	return 0;
}

/**
 * @brief Stub stop when MQTT support is disabled.
 */
void mqtt_manager_stop(void)
{
}

/**
 * @brief Stub connection check when MQTT support is disabled.
 */
bool mqtt_manager_is_connected(void)
{
	return false;
}

/**
 * @brief Stub reconnect count when MQTT support is disabled.
 */
uint32_t mqtt_manager_reconnects(void)
{
	return 0;
}

/**
 * @brief Stub publish returning -ENOTCONN when MQTT support is disabled.
 */
int mqtt_manager_publish(const char *suffix, const char *payload, size_t size, int qos, bool retain)
{
	ARG_UNUSED(suffix);
	ARG_UNUSED(payload);
	ARG_UNUSED(size);
	ARG_UNUSED(qos);
	ARG_UNUSED(retain);
	return -ENOTCONN;
}

/**
 * @brief Stub subscribe returning -ENOTCONN when MQTT support is disabled.
 */
int mqtt_manager_subscribe(const char *suffix)
{
	ARG_UNUSED(suffix);
	return -ENOTCONN;
}
