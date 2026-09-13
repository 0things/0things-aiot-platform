#include <app/network.h>
#include <app/services.h>
#include <app/config.h>
#include <zephyr/net/mqtt.h>
#include <zephyr/net/net_if.h>
#include <zephyr/net/socket.h>
#include <zephyr/net/tls_credentials.h>
#include <zephyr/data/json.h>
#include <zephyr/random/random.h>
#include <zephyr/sys/atomic.h>
#include <zephyr/logging/log.h>
#include <string.h>
#include <stdio.h>
LOG_MODULE_REGISTER(mqtt_manager, CONFIG_APP_LOG_LEVEL);
enum mqtt_state {
	MQTT_DISCONNECTED,
	MQTT_CONNECTING,
	MQTT_CONNECTED,
	MQTT_RECONNECTING
};
static atomic_t state, running, reconnects;
static struct mqtt_client client;
static struct net_sockaddr_storage broker;
static uint8_t rx[CONFIG_APP_PAYLOAD_SIZE + 128], tx[CONFIG_APP_PAYLOAD_SIZE + 128];
static char presence_topic[80], config_topic[80], command_topic[80];
static struct mqtt_topic will_topic;
static struct mqtt_utf8 will_message = {.utf8 = (uint8_t *)"{\"online\":false}", .size = 16};
static struct mqtt_utf8 username, password;
static bool socket_open;
static uint16_t message_id;
static uint16_t waiting_id;
static bool acknowledged;
static int ack_result;
static struct k_thread mqtt_thread, control_thread;
K_THREAD_STACK_DEFINE(mqtt_stack, 6144);
K_THREAD_STACK_DEFINE(control_stack, 6144);
K_SEM_DEFINE(wakeup, 0, 1);
K_SEM_DEFINE(completed, 0, 1);
K_MUTEX_DEFINE(api_lock);
struct request {
	char topic[80];
	char payload[CONFIG_APP_PAYLOAD_SIZE];
	size_t length;
	int qos;
	bool retain, subscribe;
};
static struct request request;
static atomic_t request_pending;
static int request_result;
struct incoming {
	bool config;
	size_t length;
	char payload[CONFIG_APP_PAYLOAD_SIZE];
};
K_MSGQ_DEFINE(incoming_queue, sizeof(struct incoming), 3, 4);
static struct incoming receive_buffer;
#if defined(CONFIG_APP_TLS)
static const unsigned char ca_certificate[] = {
#include "app_ca.inc"
	0};
static const sec_tag_t tags[] = {42};
#endif
static int socket_fd(void)
{
#if defined(CONFIG_APP_TLS)
	return client.transport.tls.sock;
#else
	return client.transport.tcp.sock;
#endif
}
static bool network_ready(void)
{
	struct net_if *iface = net_if_get_default();
	return iface && net_if_is_up(iface) &&
	       net_if_ipv4_get_global_addr(iface, NET_ADDR_PREFERRED);
}
static uint16_t next_id(void)
{
	if (!++message_id) {
		message_id++;
	}
	return message_id;
}
static int read_payload(uint8_t *buffer, size_t length)
{
	int64_t deadline = k_uptime_get() + 5000;
	size_t offset = 0;
	while (offset < length) {
		int rc = mqtt_read_publish_payload(&client, buffer + offset, length - offset);
		if (rc > 0) {
			offset += rc;
			continue;
		}
		if (rc != -EAGAIN) {
			return rc ? rc : -ECONNRESET;
		}
		if (k_uptime_get() >= deadline) {
			return -ETIMEDOUT;
		}
		struct zsock_pollfd fd = {.fd = socket_fd(), .events = ZSOCK_POLLIN};
		rc = zsock_poll(&fd, 1, 100);
		if (rc < 0) {
			return -errno;
		}
	}
	return 0;
}
static void mqtt_callback(struct mqtt_client *c, const struct mqtt_evt *event)
{
	switch (event->type) {
	case MQTT_EVT_CONNACK:
		atomic_set(&state, event->result ? MQTT_DISCONNECTED : MQTT_CONNECTED);
		break;
	case MQTT_EVT_DISCONNECT:
		atomic_set(&state, MQTT_DISCONNECTED);
		break;
	case MQTT_EVT_PUBACK:
		if (event->param.puback.message_id == waiting_id) {
			acknowledged = true;
			ack_result = event->result;
		}
		break;
	case MQTT_EVT_SUBACK:
		if (event->param.suback.message_id == waiting_id) {
			acknowledged = true;
			ack_result = event->result;
			if (event->param.suback.return_codes.len != 1 ||
			    event->param.suback.return_codes.data[0] > 1) {
				ack_result = -EACCES;
			}
		}
		break;
	case MQTT_EVT_PUBLISH: {
		const struct mqtt_publish_param *p = &event->param.publish;
		size_t n = p->message.payload.len;
		if (n >= sizeof(receive_buffer.payload)) {
			/* Abort oversized input: never leave unread bytes desynchronizing the MQTT
			 * stream. */
			atomic_set(&state, MQTT_DISCONNECTED);
			break;
		}
		int rc = read_payload((uint8_t *)receive_buffer.payload, n);
		if (rc) {
			atomic_set(&state, MQTT_DISCONNECTED);
			break;
		}
		receive_buffer.payload[n] = 0;
		receive_buffer.length = n;
		struct mqtt_utf8 topic = p->message.topic.topic;
		bool is_config = topic.size == strlen(config_topic) &&
				 !memcmp(topic.utf8, config_topic, topic.size);
		bool is_command = topic.size == strlen(command_topic) &&
				  !memcmp(topic.utf8, command_topic, topic.size);
		if (is_config || is_command) {
			receive_buffer.config = is_config;
			rc = k_msgq_put(&incoming_queue, &receive_buffer, K_NO_WAIT);
			if (rc) {
				atomic_set(&state, MQTT_DISCONNECTED);
				break;
			}
		}
		if (p->message.topic.qos == MQTT_QOS_1_AT_LEAST_ONCE) {
			struct mqtt_puback_param ack = {.message_id = p->message_id};
			if (mqtt_publish_qos1_ack(c, &ack)) {
				atomic_set(&state, MQTT_DISCONNECTED);
			}
		}
		break;
	}
	default:
		break;
	}
}
static int pump(int timeout)
{
	struct zsock_pollfd fd = {.fd = socket_fd(), .events = ZSOCK_POLLIN};
	int rc = zsock_poll(&fd, 1, timeout);
	if (rc < 0) {
		return -errno;
	}
	if (fd.revents & (ZSOCK_POLLERR | ZSOCK_POLLHUP | ZSOCK_POLLNVAL)) {
		return -ECONNRESET;
	}
	if (fd.revents & ZSOCK_POLLIN) {
		rc = mqtt_input(&client);
		if (rc && rc != -EAGAIN) {
			return rc;
		}
	}
	if (atomic_get(&state) == MQTT_DISCONNECTED) {
		return -ENOTCONN;
	}
	rc = mqtt_live(&client);
	return rc == -EAGAIN ? 0 : rc;
}
static int wait_ack(void)
{
	int64_t deadline = k_uptime_get() + 10000;
	while (!acknowledged && atomic_get(&running)) {
		if (k_uptime_get() >= deadline) {
			return -ETIMEDOUT;
		}
		int rc = pump(200);
		if (rc) {
			return rc;
		}
	}
	return acknowledged ? ack_result : -ECANCELED;
}
static int publish(const char *topic, const char *payload, size_t length, int qos, bool retain)
{
	struct mqtt_publish_param p = {0};
	p.message.topic.topic = (struct mqtt_utf8){.utf8 = (uint8_t *)topic, .size = strlen(topic)};
	p.message.topic.qos = qos;
	p.message.payload.data = (uint8_t *)payload;
	p.message.payload.len = length;
	p.message_id = next_id();
	p.retain_flag = retain;
	waiting_id = p.message_id;
	acknowledged = false;
	ack_result = 0;
	int rc = mqtt_publish(&client, &p);
	return rc || !qos ? rc : wait_ack();
}
static int subscribe(const char *topic)
{
	struct mqtt_topic t = {.topic = {.utf8 = (uint8_t *)topic, .size = strlen(topic)},
			       .qos = MQTT_QOS_1_AT_LEAST_ONCE};
	struct mqtt_subscription_list list = {.list = &t, .list_count = 1, .message_id = next_id()};
	waiting_id = list.message_id;
	acknowledged = false;
	ack_result = 0;
	int rc = mqtt_subscribe(&client, &list);
	return rc ? rc : wait_ack();
}
static int connect_broker(void)
{
	socket_open = false;
	struct zsock_addrinfo hints = {.ai_family = NET_AF_INET, .ai_socktype = NET_SOCK_STREAM};
	struct zsock_addrinfo *result;
	char port[8];
	snprintf(port, sizeof(port), "%d", CONFIG_APP_BROKER_PORT);
	int rc = zsock_getaddrinfo(CONFIG_APP_BROKER, port, &hints, &result);
	if (rc) {
		return -EHOSTUNREACH;
	}
	memcpy(&broker, result->ai_addr, result->ai_addrlen);
	zsock_freeaddrinfo(result);
	mqtt_client_init(&client);
	client.broker = &broker;
	client.evt_cb = mqtt_callback;
	client.client_id = (struct mqtt_utf8){.utf8 = (uint8_t *)device_identity_get(),
					      .size = strlen(device_identity_get())};
	client.protocol_version = MQTT_VERSION_3_1_1;
	client.rx_buf = rx;
	client.rx_buf_size = sizeof(rx);
	client.tx_buf = tx;
	client.tx_buf_size = sizeof(tx);
	client.clean_session = 1;
	client.will_topic = &will_topic;
	client.will_message = &will_message;
	client.will_retain = 1;
	if (username.size) {
		client.user_name = &username;
		client.password = &password;
	}
#if defined(CONFIG_APP_TLS)
	client.transport.type = MQTT_TRANSPORT_SECURE;
	client.transport.tls.config =
		(struct mqtt_sec_config){.peer_verify = TLS_PEER_VERIFY_REQUIRED,
					 .sec_tag_list = tags,
					 .sec_tag_count = ARRAY_SIZE(tags),
					 .hostname = CONFIG_APP_BROKER};
#else
	client.transport.type = MQTT_TRANSPORT_NON_SECURE;
#endif
	atomic_set(&state, MQTT_CONNECTING);
	rc = mqtt_connect(&client);
	if (rc) {
		return rc;
	}
	socket_open = true;
	struct timeval timeout = {.tv_sec = 5};
	zsock_setsockopt(socket_fd(), SOL_SOCKET, SO_SNDTIMEO, &timeout, sizeof(timeout));
	int64_t deadline = k_uptime_get() + 10000;
	while (atomic_get(&state) == MQTT_CONNECTING && k_uptime_get() < deadline) {
		rc = pump(200);
		if (rc) {
			return rc;
		}
	}
	if (!mqtt_manager_is_connected()) {
		return -ETIMEDOUT;
	}
	rc = subscribe(config_topic);
	if (!rc) {
		rc = subscribe(command_topic);
	}
	if (!rc) {
		rc = publish(presence_topic, "{\"online\":true}", 15, 1, true);
	}
	return rc;
}
static void complete_request(int result)
{
	request_result = result;
	atomic_clear(&request_pending);
	k_sem_give(&completed);
}
static void mqtt_entry(void *a, void *b, void *c)
{
	ARG_UNUSED(a);
	ARG_UNUSED(b);
	ARG_UNUSED(c);
	uint32_t backoff = 1;
	for (;;) {
		watchdog_manager_heartbeat(2);
		if (!atomic_get(&running) || !network_ready() ||
		    (IS_ENABLED(CONFIG_APP_TLS) && !time_manager_is_synced())) {
			if (atomic_get(&request_pending)) {
				complete_request(-ENOTCONN);
			}
			k_sem_take(&wakeup, K_SECONDS(1));
			continue;
		}
		atomic_inc(&reconnects);
		int rc = connect_broker();
		if (!rc) {
			backoff = 1;
			app_event_emit(APP_EVENT_MQTT_UP);
			while (atomic_get(&running) && network_ready() &&
			       mqtt_manager_is_connected()) {
				watchdog_manager_heartbeat(2);
				if (atomic_get(&request_pending)) {
					rc = request.subscribe
						     ? subscribe(request.topic)
						     : publish(request.topic, request.payload,
							       request.length, request.qos,
							       request.retain);
					complete_request(rc);
					if (rc) {
						break;
					}
				}
				rc = pump(200);
				if (rc) {
					break;
				}
			}
		}
		if (socket_open) {
			mqtt_abort(&client);
			socket_open = false;
		}
		atomic_set(&state, MQTT_RECONNECTING);
		app_event_emit(APP_EVENT_MQTT_DOWN);
		if (atomic_get(&request_pending)) {
			complete_request(-ENOTCONN);
		}
		LOG_WRN("MQTT disconnected: %d; retry in %us", rc, backoff);
		k_sem_take(&wakeup, K_MSEC(backoff * 1000 + sys_rand32_get() % 500));
		backoff = MIN(backoff * 2, 60U);
	}
}
struct command_data {
	char *command;
	char *url;
};
static void control_entry(void *a, void *b, void *c)
{
	ARG_UNUSED(a);
	ARG_UNUSED(b);
	ARG_UNUSED(c);
	static const struct json_obj_descr command_descr[] = {
		JSON_OBJ_DESCR_PRIM_NAMED(struct command_data, "command", command, JSON_TOK_STRING),
		JSON_OBJ_DESCR_PRIM_NAMED(struct command_data, "url", url, JSON_TOK_STRING),
	};
	struct incoming in;
	for (;;) {
		k_msgq_get(&incoming_queue, &in, K_FOREVER);
		int rc;
		if (in.config) {
			rc = remote_config_apply(in.payload, in.length);
		} else {
			struct command_data command = {0};
			int64_t fields = json_obj_parse(in.payload, in.length, command_descr,
							ARRAY_SIZE(command_descr), &command);
			if (fields < 0 || !(fields & 1)) {
				rc = -EINVAL;
			} else if (!strcmp(command.command, "flush")) {
				telemetry_uploader_flush();
				rc = 0;
			} else if (!strcmp(command.command, "firmware_update") && (fields & 2)) {
				rc = ota_manager_request(command.url);
			} else {
				rc = -ENOTSUP;
			}
		}
		const char *reason = rc == -ESTALE    ? "stale version"
				     : rc == -EINVAL  ? "invalid configuration or command"
				     : rc == -ENOTSUP ? "unsupported command"
				     : rc             ? "operation failed"
						      : "accepted";
		char response[192];
		int n = snprintf(response, sizeof(response),
				 "{\"type\":\"%s_%s\",\"errno\":%d,\"reason\":\"%s\"}",
				 in.config ? "config" : "command", rc ? "rejected" : "applied", rc,
				 reason);
		mqtt_manager_publish("event", response, n, 1, false);
	}
}
static void network_event(const struct zbus_channel *channel)
{
	const struct app_event *e = zbus_chan_const_msg(channel);
	if (e->type == APP_EVENT_NETWORK_UP || e->type == APP_EVENT_NETWORK_DOWN) {
		k_sem_give(&wakeup);
	}
}
ZBUS_LISTENER_DEFINE(mqtt_listener, network_event);
ZBUS_CHAN_ADD_OBS(app_events, mqtt_listener, 3);
int mqtt_manager_init(void)
{
	telemetry_topic(device_identity_get(), "presence", presence_topic, sizeof(presence_topic));
	telemetry_topic(device_identity_get(), "config", config_topic, sizeof(config_topic));
	telemetry_topic(device_identity_get(), "command", command_topic, sizeof(command_topic));
	will_topic = (struct mqtt_topic){
		.topic = {.utf8 = (uint8_t *)presence_topic, .size = strlen(presence_topic)},
		.qos = MQTT_QOS_1_AT_LEAST_ONCE};
	username = (struct mqtt_utf8){.utf8 = (uint8_t *)CONFIG_APP_MQTT_USERNAME,
				      .size = strlen(CONFIG_APP_MQTT_USERNAME)};
	password = (struct mqtt_utf8){.utf8 = (uint8_t *)CONFIG_APP_MQTT_PASSWORD,
				      .size = strlen(CONFIG_APP_MQTT_PASSWORD)};
#if defined(CONFIG_APP_TLS)
	int rc = tls_credential_add(tags[0], TLS_CREDENTIAL_CA_CERTIFICATE, ca_certificate,
				    sizeof(ca_certificate));
	if (rc) {
		return rc;
	}
#endif
	k_thread_create(&mqtt_thread, mqtt_stack, K_THREAD_STACK_SIZEOF(mqtt_stack), mqtt_entry,
			NULL, NULL, NULL, 6, 0, K_NO_WAIT);
	k_thread_create(&control_thread, control_stack, K_THREAD_STACK_SIZEOF(control_stack),
			control_entry, NULL, NULL, NULL, 8, 0, K_NO_WAIT);
	k_thread_name_set(&mqtt_thread, "mqtt");
	k_thread_name_set(&control_thread, "remote_config");
	return 0;
}
int mqtt_manager_start(void)
{
	atomic_set(&running, 1);
	k_sem_give(&wakeup);
	return 0;
}
void mqtt_manager_stop(void)
{
	atomic_clear(&running);
	k_sem_give(&wakeup);
}
bool mqtt_manager_is_connected(void)
{
	return atomic_get(&state) == MQTT_CONNECTED;
}
uint32_t mqtt_manager_reconnects(void)
{
	return MAX(atomic_get(&reconnects) - 1, 0);
}
static int submit(const char *suffix, const char *payload, size_t length, int qos, bool retain,
		  bool sub)
{
	if (length >= sizeof(request.payload) || qos < 0 || qos > 1 || (!payload && length)) {
		return -EINVAL;
	}
	int rc = k_mutex_lock(&api_lock, K_SECONDS(15));
	if (rc) {
		return rc;
	}
	if (!mqtt_manager_is_connected()) {
		rc = -ENOTCONN;
		goto out;
	}
	rc = telemetry_topic(device_identity_get(), suffix, request.topic, sizeof(request.topic));
	if (rc < 0) {
		goto out;
	}
	if (length) {
		memcpy(request.payload, payload, length);
	}
	request.length = length;
	request.qos = qos;
	request.retain = retain;
	request.subscribe = sub;
	atomic_set(&request_pending, 1);
	/* The actor always completes the owned slot, including disconnect and ACK timeout. */
	k_sem_take(&completed, K_FOREVER);
	rc = request_result;
out:
	k_mutex_unlock(&api_lock);
	return rc;
}
int mqtt_manager_publish(const char *suffix, const char *payload, size_t length, int qos,
			 bool retain)
{
	return submit(suffix, payload, length, qos, retain, false);
}
int mqtt_manager_subscribe(const char *suffix)
{
	return submit(suffix, NULL, 0, 1, false, true);
}
