#include <app/services.h>
#include <zephyr/sys/clock.h>
#include <time.h>
#if defined(CONFIG_APP_RTC)
#include <zephyr/drivers/rtc.h>
#include <zephyr/sys/timeutil.h>
#endif
#include <zephyr/logging/log.h>
#if defined(CONFIG_APP_MQTT)
#include <zephyr/net/sntp.h>
#include <zephyr/net/net_if.h>
#endif
LOG_MODULE_REGISTER(time_manager, CONFIG_APP_LOG_LEVEL);
K_MUTEX_DEFINE(time_lock);
static int64_t offset_ms;
static bool synced;
static struct k_work_delayable sync_work;
static struct k_work_q time_queue;
K_THREAD_STACK_DEFINE(time_stack, 3072);
static void sync_handler(struct k_work *work)
{
	ARG_UNUSED(work);
#if defined(CONFIG_APP_MQTT)
	struct sntp_time t;
	int rc = sntp_simple(CONFIG_APP_SNTP_SERVER, 5000, &t);
	if (!rc && t.seconds > 1577836800) {
		k_mutex_lock(&time_lock, K_FOREVER);
		struct timespec wall = {.tv_sec = t.seconds,
					.tv_nsec = ((uint64_t)t.fraction * 1000000000) >> 32};
		sys_clock_settime(SYS_CLOCK_REALTIME, &wall);
		offset_ms = (int64_t)t.seconds * 1000 + (((uint64_t)t.fraction * 1000) >> 32) -
			    k_uptime_get();
		synced = true;
		k_mutex_unlock(&time_lock);
		app_event_emit(APP_EVENT_TIME_SYNCED);
	}
	k_work_reschedule_for_queue(&time_queue, &sync_work, K_SECONDS(rc ? 60 : 21600));
#endif
}
static void network_event(const struct zbus_channel *channel)
{
	const struct app_event *e = zbus_chan_const_msg(channel);
	if (e->type == APP_EVENT_NETWORK_UP) {
		k_work_reschedule_for_queue(&time_queue, &sync_work, K_NO_WAIT);
	}
}
ZBUS_LISTENER_DEFINE(time_listener, network_event);
ZBUS_CHAN_ADD_OBS(app_events, time_listener, 1);
void time_manager_init(void)
{
	/* A product RTC backend may initialize CLOCK_REALTIME before application startup. */
	struct timespec ts;
#if defined(CONFIG_APP_RTC)
	const struct device *rtc = DEVICE_DT_GET(DT_ALIAS(app_rtc));
	struct rtc_time rt;
	if (device_is_ready(rtc) && !rtc_get_time(rtc, &rt)) {
		struct tm calendar = {.tm_sec = rt.tm_sec,
				      .tm_min = rt.tm_min,
				      .tm_hour = rt.tm_hour,
				      .tm_mday = rt.tm_mday,
				      .tm_mon = rt.tm_mon,
				      .tm_year = rt.tm_year};
		ts.tv_sec = timeutil_timegm64(&calendar);
		ts.tv_nsec = 0;
		if (ts.tv_sec > 1577836800) {
			sys_clock_settime(SYS_CLOCK_REALTIME, &ts);
		}
	}
#endif
	if (!sys_clock_gettime(SYS_CLOCK_REALTIME, &ts) && ts.tv_sec > 1577836800) {
		offset_ms = (int64_t)ts.tv_sec * 1000 + ts.tv_nsec / 1000000 - k_uptime_get();
		synced = true;
	}
	k_work_init_delayable(&sync_work, sync_handler);
	k_work_queue_start(&time_queue, time_stack, K_THREAD_STACK_SIZEOF(time_stack), 8, NULL);
#if defined(CONFIG_APP_MQTT)
	k_work_reschedule_for_queue(&time_queue, &sync_work, K_SECONDS(1));
#endif
}
bool time_manager_is_synced(void)
{
	k_mutex_lock(&time_lock, K_FOREVER);
	bool valid = synced;
	k_mutex_unlock(&time_lock);
	return valid;
}
int64_t time_manager_now_ms(void)
{
	k_mutex_lock(&time_lock, K_FOREVER);
	int64_t result = k_uptime_get() + (synced ? offset_ms : 0);
	k_mutex_unlock(&time_lock);
	return result;
}
void time_manager_stamp(struct telemetry_record *r)
{
	k_mutex_lock(&time_lock, K_FOREVER);
	r->uptime_ms = k_uptime_get();
	r->time_synced = synced;
	r->timestamp_ms = r->uptime_ms + (synced ? offset_ms : 0);
	k_mutex_unlock(&time_lock);
}
