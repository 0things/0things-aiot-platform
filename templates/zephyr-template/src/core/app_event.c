#include <app/services.h>
#include <zephyr/logging/log.h>

LOG_MODULE_REGISTER(app_event, CONFIG_APP_LOG_LEVEL);

/* System-wide asynchronous event channel definition */
ZBUS_CHAN_DEFINE(app_events, struct app_event, NULL, NULL, ZBUS_OBSERVERS_EMPTY,
		 ZBUS_MSG_INIT(.type = APP_EVENT_NETWORK_DOWN));

/**
 * @brief Emit an event message to all registered zbus channel listeners.
 *
 * Thread-safe and non-blocking; consumers reconcile their state asynchronously.
 *
 * @param type Event type enumeration value.
 */
void app_event_emit(enum app_event_type type)
{
	struct app_event event = {.type = type};
	int rc = zbus_chan_pub(&app_events, &event, K_NO_WAIT);
	if (rc) {
		LOG_WRN("event %d not delivered: %d", type, rc);
	}
}

