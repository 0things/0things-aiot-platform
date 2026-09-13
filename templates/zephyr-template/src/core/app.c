#include <app/app.h>
#include <app/services.h>
#include <app/config.h>
#include <app/network.h>
#include <zephyr/logging/log.h>

LOG_MODULE_REGISTER(app, CONFIG_APP_LOG_LEVEL);

/**
 * @brief Initialize all subsystem modules and run the application services.
 *
 * Boot-thread only. Initializes collector manager, config manager, device identity,
 * time manager, telemetry uploader, MQTT client, and optional Wi-Fi interface.
 *
 * @return 0 on success, negative errno on error.
 */
int app_run(void)
{
	/* 1. Discover and initialize registered sensor collectors */
	int rc = collector_manager_init();
	if (rc) {
		return rc;
	}

	/* 2. Load persisted configuration and apply default validation */
	rc = config_manager_init();
	if (rc) {
		return rc;
	}

	/* 3. Determine unique device identity (HWID, MAC, settings, or Kconfig) */
	rc = device_identity_init();
	if (rc) {
		return rc;
	}

	/* 4. Start local timekeeping and telemetry batch uploader worker */
	time_manager_init();
	telemetry_uploader_init();

	/* 5. Apply restored configuration to active collector instances */
	struct app_config restored;
	config_manager_snapshot(&restored);
	for (size_t i = 0; i < restored.count; i++) {
		collector_configure(collector_find(restored.collectors[i].name),
				    &restored.collectors[i]);
	}

	/* 6. Initialize MQTT communication manager */
	rc = mqtt_manager_init();
	if (rc) {
		return rc;
	}

#if defined(CONFIG_APP_WIFI)
	/* 7. Initialize and start Wi-Fi station connection if enabled */
	rc = wifi_manager_init();
	if (rc) {
		LOG_ERR("Wi-Fi init failed: %d", rc);
	} else {
		wifi_manager_connect();
	}
#endif

	/* 8. Start MQTT connection loop and periodic health reporter */
	mqtt_manager_start();
	health_manager_start();

	LOG_INF("IoT scaffold ready: %s (%u collectors)", device_identity_get(),
		(unsigned)collector_count());
	return 0;
}

