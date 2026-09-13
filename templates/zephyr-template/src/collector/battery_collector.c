#include <app/collector.h>
#include <zephyr/drivers/sensor.h>
#if defined(CONFIG_APP_BATTERY)
#if !defined(CONFIG_APP_MOCK_SENSOR)
static const struct device *const sensor = DEVICE_DT_GET(DT_ALIAS(battery_sensor));
#endif

/**
 * @brief Initialize battery fuel gauge sensor device.
 */
static int init(void)
{
#if defined(CONFIG_APP_MOCK_SENSOR)
	return 0;
#else
	return device_is_ready(sensor) ? 0 : -ENODEV;
#endif
}

/**
 * @brief Sample battery voltage and state of charge (SoC).
 *
 * @param r Output telemetry record populated with voltage (mV) and battery percent.
 */
static int collect(struct telemetry_record *r)
{
	int64_t voltage = 3850, percent = 72;
#if !defined(CONFIG_APP_MOCK_SENSOR)
	struct sensor_value v, p;
	int rc = sensor_sample_fetch(sensor);
	if (!rc) {
		rc = sensor_channel_get(sensor, SENSOR_CHAN_GAUGE_VOLTAGE, &v);
	}
	if (!rc) {
		rc = sensor_channel_get(sensor, SENSOR_CHAN_GAUGE_STATE_OF_CHARGE, &p);
	}
	if (rc) {
		return rc;
	}
	voltage = sensor_value_to_micro(&v) / 1000;
	percent = sensor_value_to_micro(&p) / 1000000;
#endif
	r->field_count = 2;
	r->fields[0] = (struct telemetry_field){"voltage_mv", voltage, 0};
	r->fields[1] = (struct telemetry_field){"battery_percent", percent, 0};
	return 0;
}

static const struct collector_ops ops = {.init = init, .collect = collect};
COLLECTOR_DEFINE(battery, &ops, CONFIG_APP_BATTERY_COLLECTION_INTERVAL,
		 CONFIG_APP_BATTERY_UPLOAD_INTERVAL);
#endif
