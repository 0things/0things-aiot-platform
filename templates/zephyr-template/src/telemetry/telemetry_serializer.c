#include <app/telemetry.h>
#include <errno.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <inttypes.h>
bool app_name_valid(const char *s)
{
	if (!s || !*s || strlen(s) >= APP_NAME_SIZE) {
		return false;
	}
	for (; *s; s++) {
		if (!((*s >= 'a' && *s <= 'z') || (*s >= 'A' && *s <= 'Z') ||
		      (*s >= '0' && *s <= '9') || *s == '_' || *s == '-')) {
			return false;
		}
	}
	return true;
}
static int add(char **p, size_t *left, const char *fmt, ...)
{
	va_list args;
	va_start(args, fmt);
	int n = vsnprintf(*p, *left, fmt, args);
	va_end(args);
	if (n < 0 || (size_t)n >= *left) {
		return -ENOSPC;
	}
	*p += n;
	*left -= n;
	return 0;
}
int telemetry_topic(const char *id, const char *suffix, char *out, size_t size)
{
	if (!app_name_valid(id) || !app_name_valid(suffix)) {
		return -EINVAL;
	}
	int n = snprintf(out, size, "devices/%s/%s", id, suffix);
	return n < 0 || (size_t)n >= size ? -ENOSPC : n;
}
int telemetry_serialize(const char *id, const char *name, const struct telemetry_record *records,
			size_t count, char *out, size_t size)
{
	if (!app_name_valid(id) || !app_name_valid(name) || !records || !count || !size) {
		return -EINVAL;
	}
	char *p = out;
	size_t left = size;
	int rc;
#define ADD(...)                                                                                   \
	do {                                                                                       \
		rc = add(&p, &left, __VA_ARGS__);                                                  \
		if (rc) {                                                                          \
			return rc;                                                                 \
		}                                                                                  \
	} while (0)
	ADD("{\"device_id\":\"%s\",\"collector\":\"%s\",\"timestamp\":%" PRId64 ",\"records\":[",
	    id, name, records[count - 1].timestamp_ms);
	for (size_t i = 0; i < count; i++) {
		const struct telemetry_record *r = &records[i];
		if (r->field_count > TELEMETRY_FIELDS) {
			return -EINVAL;
		}
		ADD("%s{\"sequence\":%" PRIu64 ",\"timestamp\":%" PRId64 ",\"uptime_ms\":%" PRId64
		    ",\"time_synced\":%s",
		    i ? "," : "", r->sequence, r->timestamp_ms, r->uptime_ms,
		    r->time_synced ? "true" : "false");
		for (size_t j = 0; j < r->field_count; j++) {
			const struct telemetry_field *f = &r->fields[j];
			if (!app_name_valid(f->name) || f->decimals > 6) {
				return -EINVAL;
			}
			uint64_t scale = 1;
			for (int d = 0; d < f->decimals; d++) {
				scale *= 10;
			}
			uint64_t magnitude =
				f->value < 0 ? (uint64_t)(-(f->value + 1)) + 1 : (uint64_t)f->value;
			ADD(",\"%s\":%s%" PRIu64, f->name, f->value < 0 ? "-" : "",
			    magnitude / scale);
			if (f->decimals) {
				ADD(".%0*" PRIu64, f->decimals, magnitude % scale);
			}
		}
		ADD("}");
	}
	ADD("]}");
#undef ADD
	return (int)(p - out);
}
