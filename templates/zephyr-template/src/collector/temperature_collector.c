#include <app/collector.h>
#include <zephyr/drivers/sensor.h>
#if defined(CONFIG_APP_TEMPERATURE)
#if !defined(CONFIG_APP_MOCK_SENSOR)
#define SENSOR_NODE DT_ALIAS(temperature_sensor)
static const struct device *const sensor = DEVICE_DT_GET(SENSOR_NODE);
#endif
static int init(void)
{
#if defined(CONFIG_APP_MOCK_SENSOR)
	return 0;
#else
	return device_is_ready(sensor) ? 0 : -ENODEV;
#endif
}
static int collect(struct telemetry_record *r)
{
	int64_t value = 23500;
#if !defined(CONFIG_APP_MOCK_SENSOR)
	struct sensor_value v;
	int rc = sensor_sample_fetch(sensor);
	if (!rc) {
		rc = sensor_channel_get(sensor, SENSOR_CHAN_AMBIENT_TEMP, &v);
	}
	if (rc) {
		return rc;
	}
	value = sensor_value_to_micro(&v) / 1000;
#endif
	r->field_count = 1;
	r->fields[0] = (struct telemetry_field){"temperature_c", value, 3};
	return 0;
}
static const struct collector_ops ops = {.init = init, .collect = collect};
COLLECTOR_DEFINE(temperature, &ops, CONFIG_APP_TEMP_COLLECTION_INTERVAL,
		 CONFIG_APP_TEMP_UPLOAD_INTERVAL);
#endif
