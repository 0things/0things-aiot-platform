#include <app/network.h>
int mqtt_manager_init(void)
{
	return 0;
}
int mqtt_manager_start(void)
{
	return 0;
}
void mqtt_manager_stop(void)
{
}
bool mqtt_manager_is_connected(void)
{
	return false;
}
uint32_t mqtt_manager_reconnects(void)
{
	return 0;
}
int mqtt_manager_publish(const char *suffix, const char *payload, size_t size, int qos, bool retain)
{
	ARG_UNUSED(suffix);
	ARG_UNUSED(payload);
	ARG_UNUSED(size);
	ARG_UNUSED(qos);
	ARG_UNUSED(retain);
	return -ENOTCONN;
}
int mqtt_manager_subscribe(const char *suffix)
{
	ARG_UNUSED(suffix);
	return -ENOTCONN;
}
