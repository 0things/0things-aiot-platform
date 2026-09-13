#include <app/services.h>
#include <string.h>
#include <zephyr/sys/byteorder.h>
#if defined(CONFIG_HWINFO)
#include <zephyr/drivers/hwinfo.h>
#endif
#if defined(CONFIG_NETWORKING)
#include <zephyr/net/net_if.h>
#endif
#if defined(CONFIG_APP_SETTINGS)
#include <zephyr/settings/settings.h>
#endif
static char identifier[APP_NAME_SIZE];
#if defined(CONFIG_APP_SETTINGS)
static int load_identity(const char *key, size_t len, settings_read_cb read, void *arg, void *param)
{
	ARG_UNUSED(param);
	if ((key && *key) || !len || len >= sizeof(identifier)) {
		return -EINVAL;
	}
	int rc = read(arg, identifier, len);
	if (rc < 0) {
		return rc;
	}
	identifier[len] = 0;
	return app_name_valid(identifier) ? 0 : -EINVAL;
}
#endif
int device_identity_init(void)
{
	uint8_t bytes[10];
	size_t length = 0;
#if defined(CONFIG_HWINFO)
	ssize_t n = hwinfo_get_device_id(bytes, sizeof(bytes));
	if (n > 0) {
		length = n;
	}
#endif
#if defined(CONFIG_NETWORKING)
	if (!length) {
		struct net_if *iface = net_if_get_default();
		struct net_linkaddr *mac = iface ? net_if_get_link_addr(iface) : NULL;
		if (mac && mac->len && mac->len <= sizeof(bytes)) {
			length = mac->len;
			memcpy(bytes, mac->addr, length);
		}
	}
#endif
	if (length) {
		static const char hex[] = "0123456789abcdef";
		for (size_t i = 0; i < length; i++) {
			identifier[2 * i] = hex[bytes[i] >> 4];
			identifier[2 * i + 1] = hex[bytes[i] & 15];
		}
		return 0;
	}
#if defined(CONFIG_APP_SETTINGS)
	int rc = settings_load_subtree_direct("identity/id", load_identity, NULL);
	if (!rc && app_name_valid(identifier)) {
		return 0;
	}
#endif
	if (!app_name_valid(CONFIG_APP_DEVICE_ID)) {
		return -EINVAL;
	}
	strcpy(identifier, CONFIG_APP_DEVICE_ID);
	return 0;
}
const char *device_identity_get(void)
{
	return identifier;
}
