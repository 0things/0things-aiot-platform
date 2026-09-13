#include <zephyr/ztest.h>
#include <app/config.h>
#include <app/services.h>
#include <zephyr/sys/atomic.h>
#include <string.h>
#include <stdio.h>

static struct telemetry_store store;
static struct telemetry_record records[CONFIG_APP_CACHE_SIZE];
static atomic_t upload_requests[CONFIG_APP_MAX_COLLECTORS];
static int storage_error;

/**
 * @brief Collector operation mock that simulates hardware collection failure.
 */
static int fail_collect(struct telemetry_record *record)
{
	ARG_UNUSED(record);
	return -EIO;
}

static const struct collector_ops failing_ops = {.collect = fail_collect};
COLLECTOR_DEFINE(broken, &failing_ops, 1, 6);

/**
 * @brief Test mock capturing telemetry uploader upload trigger invocations.
 */
void telemetry_uploader_request(size_t i)
{
	atomic_inc(&upload_requests[i]);
}

/**
 * @brief Test mock allowing controlled failure simulation during config saving.
 */
int config_storage_save(const struct app_config *config)
{
	ARG_UNUSED(config);
	return storage_error;
}

/**
 * @brief Global test suite fixture initialization.
 */
static void *setup(void)
{
	zassert_ok(collector_manager_init());
	zassert_ok(config_manager_init());
	time_manager_init();
	return NULL;
}

/**
 * @brief Per-test setup resetting test state and initializing store.
 */
static void before(void *fixture)
{
	ARG_UNUSED(fixture);
	storage_error = 0;
	telemetry_store_init(&store, 6);
}

/**
 * @test Verify ring buffer overflow eviction behavior and sequence acknowledgement pop semantics.
 */
ZTEST(core, ring_overflow_ack_does_not_delete_new_records)
{
	struct telemetry_record r = {0};
	for (int i = 0; i < 6; i++) {
		zassert_ok(telemetry_store_append(&store, &r));
	}
	zassert_equal(telemetry_store_peek(&store, records, 2), 2);
	uint64_t ack = records[1].sequence;
	for (int i = 0; i < 3; i++) {
		zassert_ok(telemetry_store_append(&store, &r));
	}
	telemetry_store_pop(&store, ack);
	zassert_equal(telemetry_store_count(&store), 6);
	zassert_equal(telemetry_store_dropped(&store), 3);
	zassert_equal(telemetry_store_peek(&store, records, 6), 6);
	zassert_equal(records[0].sequence, 4);
	telemetry_store_pop(&store, records[1].sequence);
	zassert_equal(telemetry_store_count(&store), 4);
}

/**
 * @test Verify ring buffer dynamic resizing, wrap-around index calculation, and clear.
 */
ZTEST(core, ring_resize_and_wrap)
{
	struct telemetry_record r = {0};
	for (int i = 0; i < 20; i++) {
		telemetry_store_append(&store, &r);
	}
	zassert_ok(telemetry_store_resize(&store, 2));
	zassert_equal(telemetry_store_peek(&store, records, 6), 2);
	zassert_equal(records[0].sequence, 19);
	zassert_equal(records[1].sequence, 20);
	zassert_equal(telemetry_store_dropped(&store), 18);
	zassert_equal(telemetry_store_resize(&store, 0), -EINVAL);
	telemetry_store_clear(&store);
	zassert_equal(telemetry_store_count(&store), 0);
}

/**
 * @test Verify fixed-point decimal scaling, negative values, and JSON formatting boundaries.
 */
ZTEST(core, serialization_fixed_point_and_bounds)
{
	struct telemetry_record r = {
		.sequence = 1,
		.timestamp_ms = 123,
		.field_count = 2,
		.fields = {{"temperature_c", -125, 3}, {"limit", INT64_MIN, 0}}};
	char json[1024];
	zassert_true(telemetry_serialize("device-001", "temperature", &r, 1, json, sizeof(json)) >
		     0);
	zassert_not_null(strstr(json, "\"temperature_c\":-0.125"));
	zassert_not_null(strstr(json, "-9223372036854775808"));
	zassert_not_null(strstr(json, "\"time_synced\":false"));
	zassert_equal(telemetry_serialize("device-001", "temperature", &r, 1, json, 8), -ENOSPC);
	r.fields[0].name = "bad\"key";
	zassert_equal(telemetry_serialize("device-001", "temperature", &r, 1, json, sizeof(json)),
		      -EINVAL);
}

/**
 * @test Verify MQTT topic validation and rejection of invalid chars/wildcards.
 */
ZTEST(core, topics_reject_injection)
{
	char topic[80];
	zassert_true(telemetry_topic("device-001", "telemetry", topic, sizeof(topic)) > 0);
	zassert_equal(strcmp(topic, "devices/device-001/telemetry"), 0);
	zassert_equal(telemetry_topic("bad/+", "config", topic, sizeof(topic)), -EINVAL);
	zassert_equal(telemetry_topic("device-001", "config", topic, 3), -ENOSPC);
}

/**
 * @test Verify deadline computation logic for periodic collections and uploads.
 */
ZTEST(core, schedule_five_minutes_thirty_minutes)
{
	int64_t collect = 300000, upload = 1800000;
	unsigned int collected = 0, uploaded = 0;
	for (int64_t now = 1000; now <= 1800000; now += 1000) {
		if (now >= collect) {
			collected++;
			collect = collector_next_deadline(collect, now, 300);
		}
		if (now >= upload) {
			uploaded++;
			upload = collector_next_deadline(upload, now, 1800);
		}
	}
	zassert_equal(collected, 6);
	zassert_equal(uploaded, 1);
	zassert_equal(collector_next_deadline(1000, 10000, 1), 11000);
}

/**
 * @test Verify config structure parameter validation rules.
 */
ZTEST(core, validation_and_upload_faster_than_collection)
{
	struct app_config c;
	config_manager_snapshot(&c);
	c.collectors[0].upload_interval_sec = 1;
	zassert_ok(config_validate(&c));
	c.collectors[0].collection_interval_sec = -1;
	zassert_equal(config_validate(&c), -EINVAL);
	c.collectors[0].collection_interval_sec = 1;
	c.collectors[0].cache_size = CONFIG_APP_CACHE_SIZE + 1;
	zassert_equal(config_validate(&c), -EINVAL);
}

/**
 * @test Verify atomic version checking and rollback on persistence failure.
 */
ZTEST(core, version_and_persistence_failure_are_atomic)
{
	struct app_config c, original, after;
	config_manager_snapshot(&original);
	c = original;
	c.version++;
	c.collectors[0].collection_interval_sec = 99;
	storage_error = -EIO;
	zassert_equal(config_manager_apply(&c), -EIO);
	config_manager_snapshot(&after);
	zassert_mem_equal(&original, &after, sizeof(original));
	storage_error = 0;
	zassert_ok(config_manager_apply(&c));
	zassert_equal(config_manager_apply(&c), -ESTALE);
	config_manager_snapshot(&after);
	zassert_equal(after.version, c.version);
}

/**
 * @brief Helper executing remote configuration apply over JSON string.
 */
static int apply(const char *body)
{
	char json[CONFIG_APP_PAYLOAD_SIZE];
	strcpy(json, body);
	return remote_config_apply(json, strlen(json));
}

/**
 * @test Verify remote config parsing, partial updates, and malformed payload rejections.
 */
ZTEST(core, remote_partial_update_and_invalid_transactions)
{
	struct app_config c;
	config_manager_snapshot(&c);
	char json[512];
	snprintf(json, sizeof(json),
		 "{\"version\":%u,\"collectors\":{\"temperature\":{\"collection_interval_sec\":120,"
		 "\"upload_interval_sec\":600}}}",
		 c.version + 1);
	zassert_ok(apply(json));
	config_manager_snapshot(&c);
	int t = -1;
	for (size_t i = 0; i < c.count; i++) {
		if (!strcmp(c.collectors[i].name, "temperature")) {
			t = i;
		}
	}
	zassert_true(t >= 0);
	zassert_equal(c.collectors[t].collection_interval_sec, 120);
	zassert_equal(c.collectors[t].upload_interval_sec, 600);
	const char *bad[] = {
		"{\"version\":999,\"collectors\":{\"temperature\":{\"enabled\":true}}} garbage",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"enabled\":true}}} {}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"collection_interval_sec\":-1}}"
		"}",
		"{\"version\":999,\"collectors\":{\"unknown\":{\"enabled\":true}}}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"enabled\":1}}}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"cache_size\":0}}}",
		"{\"version\":999,\"version\":1000,\"collectors\":{}}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"retry\":1,\"retry\":2}}}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"collection_interval_sec\":1.5}"
		"}}",
		"{\"version\":999,\"collectors\":{\"temperature\":{\"enabled\":true}}",
	};
	for (size_t i = 0; i < ARRAY_SIZE(bad); i++) {
		zassert_true(apply(bad[i]) < 0, "accepted bad JSON %u", i);
	}
	struct app_config after;
	config_manager_snapshot(&after);
	zassert_mem_equal(&c, &after, sizeof(c));
}

/**
 * @test Verify mock collector sampling types and schema invariants.
 */
ZTEST(core, mock_collectors_are_typed)
{
	for (size_t i = 0; i < collector_count(); i++) {
		const struct collector *c = collector_at(i)->definition;
		struct telemetry_record r = {0};
		if (!strcmp(c->name, "broken")) {
			zassert_equal(c->ops->collect(&r), -EIO);
			continue;
		}
		zassert_ok(c->ops->init());
		zassert_ok(c->ops->collect(&r));
		zassert_true(r.field_count >= 1 && r.field_count <= 6);
		char json[1024];
		zassert_true(telemetry_serialize("test", c->name, &r, 1, json, sizeof(json)) > 0);
	}
}

/**
 * @test Verify workqueue offline sampling execution, error recovery, and dynamic reconfiguration.
 */
ZTEST(core, real_workqueue_offline_collection_and_reschedule)
{
	int i = collector_find("temperature");
	struct collector_config cfg;
	collector_snapshot(i, &cfg, NULL, NULL);
	struct collector_config broken_cfg;
	int broken = collector_find("broken");
	collector_snapshot(broken, &broken_cfg, NULL, NULL);
	collector_configure(broken, &broken_cfg);
	cfg.enabled = true;
	cfg.collection_interval_sec = 1;
	cfg.upload_interval_sec = 6;
	atomic_clear(&upload_requests[i]);
	telemetry_store_clear(&collector_at(i)->store);
	collector_configure(i, &cfg);
	k_sleep(K_MSEC(6300));
	zassert_equal(telemetry_store_count(&collector_at(i)->store), 6);
	zassert_equal(atomic_get(&upload_requests[i]), 1);
	uint64_t errors;
	collector_snapshot(broken, NULL, NULL, &errors);
	zassert_true(errors > 0);
	broken_cfg.enabled = false;
	collector_configure(broken, &broken_cfg);
	cfg.enabled = false;
	collector_configure(i, &cfg);
	k_sleep(K_MSEC(1200));
	zassert_equal(telemetry_store_count(&collector_at(i)->store), 6);
	cfg.enabled = true;
	cfg.collection_interval_sec = 1;
	cfg.upload_interval_sec = 1;
	collector_configure(i, &cfg);
	k_sleep(K_MSEC(1200));
	zassert_equal(telemetry_store_count(&collector_at(i)->store), 7);
	zassert_equal(atomic_get(&upload_requests[i]), 2);
	cfg.enabled = false;
	collector_configure(i, &cfg);
}

ZTEST_SUITE(core, NULL, setup, before, NULL, NULL);
