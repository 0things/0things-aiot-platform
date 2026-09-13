#include <app/collector.h>
#include <zephyr/drivers/sensor.h>
#if defined(CONFIG_APP_IMU)
#if !defined(CONFIG_APP_MOCK_SENSOR)
static const struct device *const sensor = DEVICE_DT_GET(DT_ALIAS(imu_sensor));
#endif

/**
 * @brief Initialize 6-axis IMU (accelerometer & gyroscope) device.
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
 * @brief Sample 3-axis acceleration and 3-axis angular velocity.
 *
 * @param r Output telemetry record populated with accel (mg) and gyro (mdps).
 */
static int collect(struct telemetry_record *r)
{
	int64_t values[6] = {10, 20, 1001, 100, 200, 300};
	static const char *const names[] = {"accel_x_mg",  "accel_y_mg",  "accel_z_mg",
					    "gyro_x_mdps", "gyro_y_mdps", "gyro_z_mdps"};
#if !defined(CONFIG_APP_MOCK_SENSOR)
	struct sensor_value accel[3], gyro[3];
	int rc = sensor_sample_fetch(sensor);
	if (!rc) {
		rc = sensor_channel_get(sensor, SENSOR_CHAN_ACCEL_XYZ, accel);
	}
	if (!rc) {
		rc = sensor_channel_get(sensor, SENSOR_CHAN_GYRO_XYZ, gyro);
	}
	if (rc) {
		return rc;
	}
	for (int i = 0; i < 3; i++) {
		values[i] = sensor_value_to_micro(&accel[i]) * 1000 / SENSOR_G;
		/* rad/s in micro-units -> millidegrees/s, using Zephyr's micro-radian pi. */
		values[i + 3] = sensor_value_to_micro(&gyro[i]) * 180000 / SENSOR_PI;
	}
#endif
	r->field_count = 6;
	for (int i = 0; i < 6; i++) {
		r->fields[i] = (struct telemetry_field){names[i], values[i], 0};
	}
	return 0;
}

static const struct collector_ops ops = {.init = init, .collect = collect};
COLLECTOR_DEFINE(imu, &ops, CONFIG_APP_IMU_COLLECTION_INTERVAL, CONFIG_APP_IMU_UPLOAD_INTERVAL);
#endif
