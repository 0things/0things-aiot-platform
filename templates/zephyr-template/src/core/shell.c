#include <app/services.h>
#include <app/config.h>
#include <app/network.h>
#include <zephyr/shell/shell.h>
#include <inttypes.h>

/**
 * @brief Print runtime statistics and configuration for a specific collector.
 */
static void show_collector(const struct shell *sh, size_t i)
{
	struct collector_config cfg;
	enum collector_state state;
	uint64_t errors;
	collector_snapshot(i, &cfg, &state, &errors);
	shell_print(sh, "%s enabled=%d state=%d collect=%ds upload=%ds records=%u errors=%" PRIu64,
		    cfg.name, cfg.enabled, state, cfg.collection_interval_sec,
		    cfg.upload_interval_sec,
		    (unsigned)telemetry_store_count(&collector_at(i)->store), errors);
}

/**
 * @brief Shell command handler for general device status.
 */
static int status(const struct shell *sh, size_t argc, char **argv)
{
	ARG_UNUSED(argc);
	ARG_UNUSED(argv);
	shell_print(sh, "device=%s mqtt=%s time=%s uptime=%" PRId64 "ms", device_identity_get(),
		    mqtt_manager_is_connected() ? "connected" : "disconnected",
		    time_manager_is_synced() ? "synced" : "uptime", k_uptime_get());
#if defined(CONFIG_APP_WIFI)
	shell_print(sh, "wifi_state=%d", wifi_manager_state());
#endif
	return 0;
}

/**
 * @brief Shell command handler to list all active collectors.
 */
static int collectors(const struct shell *sh, size_t argc, char **argv)
{
	ARG_UNUSED(argc);
	ARG_UNUSED(argv);
	for (size_t i = 0; i < collector_count(); i++) {
		show_collector(sh, i);
	}
	return 0;
}

/**
 * @brief Shell command handler to inspect a single collector by name.
 */
static int one(const struct shell *sh, size_t argc, char **argv)
{
	ARG_UNUSED(argc);
	int i = collector_find(argv[1]);
	if (i < 0) {
		shell_error(sh, "Unknown collector");
		return i;
	}
	show_collector(sh, i);
	return 0;
}

/**
 * @brief Shell command handler to display the active configuration metadata.
 */
static int config(const struct shell *sh, size_t argc, char **argv)
{
	struct app_config cfg;
	config_manager_snapshot(&cfg);
	shell_print(sh, "version=%u persistence=%s", cfg.version,
		    IS_ENABLED(CONFIG_APP_SETTINGS) ? "settings" : "RAM");
	return collectors(sh, argc, argv);
}

/* Application shell subcommand tree */
SHELL_STATIC_SUBCMD_SET_CREATE(app_commands, SHELL_CMD(status, NULL, "Device status", status),
			       SHELL_CMD(collectors, NULL, "List collectors", collectors),
			       SHELL_CMD_ARG(collector, NULL, "Show collector", one, 2, 0),
			       SHELL_CMD(config, NULL, "Show active config (no secrets)", config),
			       SHELL_SUBCMD_SET_END);
SHELL_CMD_REGISTER(app, &app_commands, "IoT scaffold", NULL);

/* Telemetry cache inspection command */
SHELL_STATIC_SUBCMD_SET_CREATE(telemetry_commands,
			       SHELL_CMD(status, NULL, "Cache status", collectors),
			       SHELL_SUBCMD_SET_END);
SHELL_CMD_REGISTER(telemetry, &telemetry_commands, "Telemetry", NULL);

/* MQTT status inspection command */
SHELL_STATIC_SUBCMD_SET_CREATE(mqtt_commands, SHELL_CMD(status, NULL, "MQTT status", status),
			       SHELL_SUBCMD_SET_END);
SHELL_CMD_REGISTER(mqtt, &mqtt_commands, "MQTT", NULL);

#if !defined(CONFIG_NET_L2_WIFI_SHELL)
/* Standalone Wi-Fi command if L2 shell is disabled */
SHELL_STATIC_SUBCMD_SET_CREATE(wifi_commands, SHELL_CMD(status, NULL, "Network status", status),
			       SHELL_SUBCMD_SET_END);
SHELL_CMD_REGISTER(wifi, &wifi_commands, "Wi-Fi", NULL);
#endif

