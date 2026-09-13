#include <app/network.h>
#include <app/services.h>
#include <zephyr/net/wifi_mgmt.h>
#include <zephyr/net/net_event.h>
#include <zephyr/net/net_if.h>
#include <zephyr/net/dhcpv4.h>
#include <zephyr/sys/atomic.h>
#include <zephyr/random/random.h>
#include <zephyr/logging/log.h>
#include <string.h>
LOG_MODULE_REGISTER(wifi_manager, CONFIG_APP_LOG_LEVEL);

static struct net_mgmt_event_callback callback;
static struct k_work_delayable connect_work;
static struct net_if *interface;
static atomic_t state, wanted, retries;
static bool initialized;
K_MUTEX_DEFINE(control_lock);
static atomic_t backoff = ATOMIC_INIT(1);
static struct k_work_q wifi_queue;
K_THREAD_STACK_DEFINE(wifi_stack, 3072);
K_EVENT_DEFINE(network_ready);

/**
 * @brief Check whether the Wi-Fi interface is up with an active preferred IPv4 address.
 */
static bool has_address(void)
{
	return interface && net_if_is_up(interface) &&
	       net_if_ipv4_get_global_addr(interface, NET_ADDR_PREFERRED);
}

/**
 * @brief Handle network disconnection and schedule next reconnect with exponential backoff & jitter.
 */
static void retry(void)
{
	atomic_set(&state, NETWORK_RECONNECTING);
	k_event_clear(&network_ready, 1);
	app_event_emit(APP_EVENT_NETWORK_DOWN);
	if (atomic_get(&wanted)) {
		k_work_reschedule_for_queue(
			&wifi_queue, &connect_work,
			K_MSEC(atomic_get(&backoff) * 1000 + sys_rand32_get() % 500));
		atomic_set(&backoff, MIN(atomic_get(&backoff) * 2, 60));
	}
}

/**
 * @brief Work queue handler executing connection attempts against configured SSID.
 */
static void connect_handler(struct k_work *work)
{
	ARG_UNUSED(work);
	if (!atomic_get(&wanted)) {
		return;
	}
	if (has_address()) {
		atomic_set(&state, NETWORK_CONNECTED);
		atomic_set(&backoff, 1);
		k_event_post(&network_ready, 1);
		app_event_emit(APP_EVENT_NETWORK_UP);
		return;
	}
	struct wifi_connect_req_params params = {
		.ssid = (const uint8_t *)CONFIG_APP_WIFI_SSID,
		.ssid_length = strlen(CONFIG_APP_WIFI_SSID),
		.psk = (const uint8_t *)CONFIG_APP_WIFI_PASSWORD,
		.psk_length = strlen(CONFIG_APP_WIFI_PASSWORD),
		.security = strlen(CONFIG_APP_WIFI_PASSWORD) ? WIFI_SECURITY_TYPE_PSK
							     : WIFI_SECURITY_TYPE_NONE,
		.channel = WIFI_CHANNEL_ANY,
		.mfp = WIFI_MFP_OPTIONAL,
		.timeout = 20,
	};
	atomic_set(&state, NETWORK_CONNECTING);
	atomic_inc(&retries);
	int rc = net_mgmt(NET_REQUEST_WIFI_CONNECT, interface, &params, sizeof(params));
	if (rc && rc != -EALREADY) {
		retry();
	} else {
		k_work_reschedule_for_queue(&wifi_queue, &connect_work, K_SECONDS(30));
	}
}

/**
 * @brief Net management callback reacting to Wi-Fi connection and IPv4 events.
 */
static void event_handler(struct net_mgmt_event_callback *cb, uint64_t event, struct net_if *iface)
{
	if (iface != interface) {
		return;
	}
	if (!atomic_get(&wanted)) {
		return;
	}
	if (event == NET_EVENT_IPV4_ADDR_ADD && has_address()) {
		atomic_set(&state, NETWORK_CONNECTED);
		atomic_set(&backoff, 1);
		k_event_post(&network_ready, 1);
		k_work_cancel_delayable(&connect_work);
		app_event_emit(APP_EVENT_NETWORK_UP);
	} else if (event == NET_EVENT_WIFI_CONNECT_RESULT) {
		const struct wifi_status *status = cb->info;
		if (!status || status->status) {
			retry();
		} else {
			net_dhcpv4_start(interface);
		}
	} else if (event == NET_EVENT_WIFI_DISCONNECT_RESULT || event == NET_EVENT_IPV4_ADDR_DEL ||
		   event == NET_EVENT_IF_DOWN) {
		retry();
	}
}

int wifi_manager_init(void)
{
	interface = net_if_get_first_wifi();
	if (!interface) {
		return -ENODEV;
	}
	if (!strlen(CONFIG_APP_WIFI_SSID) || strlen(CONFIG_APP_WIFI_SSID) > 32 ||
	    (strlen(CONFIG_APP_WIFI_PASSWORD) &&
	     (strlen(CONFIG_APP_WIFI_PASSWORD) < 8 || strlen(CONFIG_APP_WIFI_PASSWORD) > 64))) {
		atomic_set(&state, NETWORK_FAILED);
		return -EINVAL;
	}
	k_work_init_delayable(&connect_work, connect_handler);
	net_mgmt_init_event_callback(
		&callback, event_handler,
		NET_EVENT_WIFI_CONNECT_RESULT | NET_EVENT_WIFI_DISCONNECT_RESULT |
			NET_EVENT_IPV4_ADDR_ADD | NET_EVENT_IPV4_ADDR_DEL | NET_EVENT_IF_DOWN);
	k_work_queue_start(&wifi_queue, wifi_stack, K_THREAD_STACK_SIZEOF(wifi_stack), 6, NULL);
	net_mgmt_add_event_callback(&callback);
	initialized = true;
	return 0;
}

int wifi_manager_connect(void)
{
	if (!initialized) {
		return -ENODEV;
	}
	k_mutex_lock(&control_lock, K_FOREVER);
	atomic_set(&wanted, 1);
	k_work_reschedule_for_queue(&wifi_queue, &connect_work, K_NO_WAIT);
	k_mutex_unlock(&control_lock);
	return 0;
}

int wifi_manager_disconnect(void)
{
	if (!initialized) {
		return -ENODEV;
	}
	k_mutex_lock(&control_lock, K_FOREVER);
	atomic_clear(&wanted);
	struct k_work_sync sync;
	k_work_cancel_delayable_sync(&connect_work, &sync);
	k_event_clear(&network_ready, 1);
	atomic_set(&state, NETWORK_DISCONNECTED);
	app_event_emit(APP_EVENT_NETWORK_DOWN);
	int rc = net_mgmt(NET_REQUEST_WIFI_DISCONNECT, interface, NULL, 0);
	k_mutex_unlock(&control_lock);
	return rc;
}

bool wifi_manager_is_connected(void)
{
	return atomic_get(&state) == NETWORK_CONNECTED && has_address();
}

int wifi_manager_wait_connected(k_timeout_t timeout)
{
	return k_event_wait(&network_ready, 1, false, timeout) ? 0 : -ETIMEDOUT;
}

enum network_state wifi_manager_state(void)
{
	return atomic_get(&state);
}

uint32_t wifi_manager_reconnects(void)
{
	return MAX(atomic_get(&retries) - 1, 0);
}
