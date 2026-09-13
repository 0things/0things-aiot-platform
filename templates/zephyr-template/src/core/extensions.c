#include <app/services.h>
#include <errno.h>

K_MUTEX_DEFINE(heartbeat_lock);

/* Heartbeat timestamps for subsystems (0: collector, 1: uploader, 2: mqtt) */
static int64_t heartbeats[3];

/**
 * @brief Default weak stub for OTA firmware download request.
 *
 * Can be overridden by concrete product implementations.
 *
 * @param url Remote firmware URL string.
 * @return -ENOTSUP by default.
 */
__weak int ota_manager_request(const char *url)
{
	ARG_UNUSED(url);
	return -ENOTSUP;
}

/**
 * @brief Record software watchdog activity heartbeat for a subsystem.
 *
 * @param subsystem Subsystem index (0: collector, 1: uploader, 2: mqtt).
 */
void watchdog_manager_heartbeat(unsigned int subsystem)
{
	if (subsystem >= ARRAY_SIZE(heartbeats)) {
		return;
	}
	k_mutex_lock(&heartbeat_lock, K_FOREVER);
	heartbeats[subsystem] = k_uptime_get();
	k_mutex_unlock(&heartbeat_lock);
}

/**
 * @brief Check if all monitored subsystems updated heartbeats within max age.
 *
 * @param age Maximum allowable elapsed time since last heartbeat in milliseconds.
 * @return true if all subsystems are within deadline, false otherwise.
 */
bool watchdog_manager_healthy(int64_t age)
{
	bool healthy = true;
	k_mutex_lock(&heartbeat_lock, K_FOREVER);
	for (size_t i = 0; i < ARRAY_SIZE(heartbeats); i++) {
		if (k_uptime_get() - heartbeats[i] > age) {
			healthy = false;
		}
	}
	k_mutex_unlock(&heartbeat_lock);
	return healthy;
}

