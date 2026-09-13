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

/* Global buffer holding the immutable device identity string */
static char identifier[APP_NAME_SIZE];

#if defined(CONFIG_APP_SETTINGS)
/**
 * @brief Settings callback to load persisted device identity.
 */
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

/**
 * @brief Initialize device identifier with hardware fallback sequence:
 * 1. Hardware unique ID (hwinfo)
 * 2. Network MAC address
 * 3. Persisted settings entry (identity/id)
 * 4. Kconfig default (CONFIG_APP_DEVICE_ID)
 *
 * @return 0 on success, -EINVAL on invalid identifier.
 */
int device_identity_init(void)
{
	uint8_t bytes[10];
	size_t length = 0;

#if defined(CONFIG_HWINFO)
	/* 1. Try hardware unique device ID */
	ssize_t n = hwinfo_get_device_id(bytes, sizeof(bytes));
	if (n > 0) {
		length = n;
	}
#endif

#if defined(CONFIG_NETWORKING)
	/* 2. Fallback to default network interface MAC address */
	if (!length) {
		struct net_if *iface = net_if_get_default();
		struct net_linkaddr *mac = iface ? net_if_get_link_addr(iface) : NULL;
		if (mac && mac->len && mac->len <= sizeof(bytes)) {
			length = mac->len;
			memcpy(bytes, mac->addr, length);
		}
	}
#endif

	/* Convert binary ID to hex string if found */
	if (length) {
		static const char hex[] = "0123456789abcdef";
		for (size_t i = 0; i < length; i++) {
			identifier[2 * i] = hex[bytes[i] >> 4];
			identifier[2 * i + 1] = hex[bytes[i] & 15];
		}
		return 0;
	}

#if defined(CONFIG_APP_SETTINGS)
	/* 3. Check persistent settings */
	int rc = settings_load_subtree_direct("identity/id", load_identity, NULL);
	if (!rc && app_name_valid(identifier)) {
		return 0;
	}
#endif

	/* 4. Fallback to static Kconfig device ID */
	if (!app_name_valid(CONFIG_APP_DEVICE_ID)) {
		return -EINVAL;
	}
	strcpy(identifier, CONFIG_APP_DEVICE_ID);
	return 0;
}

/**
 * @brief Retrieve the global device identity string.
 *
 * @return Pointer to null-terminated identifier string.
 */
const char *device_identity_get(void)
{
	return identifier;
}
