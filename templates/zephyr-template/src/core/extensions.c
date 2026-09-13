#include <app/services.h>
#include <errno.h>
K_MUTEX_DEFINE(heartbeat_lock);
static int64_t heartbeats[3];
__weak int ota_manager_request(const char *url)
{
	ARG_UNUSED(url);
	return -ENOTSUP;
}
void watchdog_manager_heartbeat(unsigned int subsystem)
{
	if (subsystem >= ARRAY_SIZE(heartbeats)) {
		return;
	}
	k_mutex_lock(&heartbeat_lock, K_FOREVER);
	heartbeats[subsystem] = k_uptime_get();
	k_mutex_unlock(&heartbeat_lock);
}
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
