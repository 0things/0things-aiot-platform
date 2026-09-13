#include <app/config.h>
#include <app/services.h>
#include <zephyr/logging/log.h>
#include <string.h>
#if defined(CONFIG_APP_SETTINGS)
#include <zephyr/settings/settings.h>
#endif
LOG_MODULE_REGISTER(config_manager, CONFIG_APP_LOG_LEVEL);

K_MUTEX_DEFINE(config_lock);
static struct app_config active;

/**
 * @brief Validate application configuration struct boundaries and consistency.
 *
 * @param c Pointer to the configuration struct to validate.
 * @return 0 on success, or -EINVAL on invalid format/ranges.
 */
int config_validate(const struct app_config *c)
{
	if (c->schema != 1 || c->count != collector_count()) {
		return -EINVAL;
	}
	for (size_t i = 0; i < c->count; i++) {
		const struct collector_config *p = &c->collectors[i];
		if (!memchr(p->name, 0, sizeof(p->name)) || !app_name_valid(p->name) ||
		    collector_find(p->name) < 0 || p->collection_interval_sec < 1 ||
		    p->collection_interval_sec > 86400 || p->upload_interval_sec < 1 ||
		    p->upload_interval_sec > 604800 || p->cache_size < 1 ||
		    p->cache_size > CONFIG_APP_CACHE_SIZE || p->retry < 0 || p->retry > 10 ||
		    p->priority < 0 || p->priority > 15) {
			return -EINVAL;
		}
		for (size_t j = 0; j < i; j++) {
			if (!strcmp(c->collectors[j].name, p->name)) {
				return -EINVAL;
			}
		}
	}
	return 0;
}

#if defined(CONFIG_APP_SETTINGS)
/**
 * @brief Direct read callback for Zephyr Settings tree loader.
 */
static int load_config(const char *name, size_t length, settings_read_cb read, void *arg,
		       void *param)
{
	ARG_UNUSED(param);
	if ((name && *name) || length != sizeof(struct app_config)) {
		return -EINVAL;
	}
	struct app_config restored;
	int rc = read(arg, &restored, sizeof(restored));
	if (rc != sizeof(restored)) {
		return rc < 0 ? rc : -EIO;
	}
	rc = config_validate(&restored);
	if (!rc) {
		active = restored;
	}
	return rc;
}
#endif

/**
 * @brief Save configuration structure into non-volatile storage.
 */
__weak int config_storage_save(const struct app_config *config)
{
#if defined(CONFIG_APP_SETTINGS)
	return settings_save_one("app/config", config, sizeof(*config));
#else
	ARG_UNUSED(config);
	return 0;
#endif
}

int config_manager_init(void)
{
	active.schema = 1;
	active.count = collector_count();
	for (size_t i = 0; i < active.count; i++) {
		collector_snapshot(i, &active.collectors[i], NULL, NULL);
	}
#if defined(CONFIG_APP_SETTINGS)
	int rc = settings_subsys_init();
	if (rc) {
		return rc;
	}
	rc = settings_load_subtree_direct("app/config", load_config, NULL);
	if (rc) {
		LOG_WRN("Stored config rejected: %d; using defaults", rc);
	}
#endif
	return config_validate(&active);
}

void config_manager_snapshot(struct app_config *out)
{
	k_mutex_lock(&config_lock, K_FOREVER);
	*out = active;
	k_mutex_unlock(&config_lock);
}

int config_manager_apply(const struct app_config *c)
{
	int rc = config_validate(c);
	if (rc) {
		return rc;
	}
	k_mutex_lock(&config_lock, K_FOREVER);
	if (c->version <= active.version) {
		rc = -ESTALE;
		goto done;
	}
	/* Single settings value makes version and data durable together, before runtime mutation.
	 */
	rc = config_storage_save(c);
	if (rc) {
		goto done;
	}
	active = *c;
	for (size_t i = 0; i < c->count; i++) {
		collector_configure(collector_find(c->collectors[i].name), &c->collectors[i]);
	}
	app_event_emit(APP_EVENT_CONFIG_UPDATED);
done:
	k_mutex_unlock(&config_lock);
	return rc;
}
