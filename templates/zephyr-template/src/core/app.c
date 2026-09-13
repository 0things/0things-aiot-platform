#include <app/app.h>
#include <app/services.h>
#include <app/config.h>
#include <app/network.h>
#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(app, CONFIG_APP_LOG_LEVEL);
int app_run(void)
{
	int rc = collector_manager_init();
	if (rc) {
		return rc;
	}
	rc = config_manager_init();
	if (rc) {
		return rc;
	}
	rc = device_identity_init();
	if (rc) {
		return rc;
	}
	time_manager_init();
	telemetry_uploader_init();
	struct app_config restored;
	config_manager_snapshot(&restored);
	for (size_t i = 0; i < restored.count; i++) {
		collector_configure(collector_find(restored.collectors[i].name),
				    &restored.collectors[i]);
	}
	rc = mqtt_manager_init();
	if (rc) {
		return rc;
	}
#if defined(CONFIG_APP_WIFI)
	rc = wifi_manager_init();
	if (rc) {
		LOG_ERR("Wi-Fi init failed: %d", rc);
	} else {
		wifi_manager_connect();
	}
#endif
	mqtt_manager_start();
	health_manager_start();
	LOG_INF("IoT scaffold ready: %s (%u collectors)", device_identity_get(),
		(unsigned)collector_count());
	return 0;
}
