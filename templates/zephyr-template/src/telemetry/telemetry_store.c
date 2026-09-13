#include <app/telemetry.h>
#include <errno.h>
#include <string.h>

void telemetry_store_init(struct telemetry_store *s, size_t capacity)
{
	memset(s, 0, sizeof(*s));
	k_mutex_init(&s->lock);
	s->capacity = CLAMP(capacity, 1, CONFIG_APP_CACHE_SIZE);
	s->next_sequence = 1;
}

int telemetry_store_append(struct telemetry_store *s, const struct telemetry_record *r)
{
	if (!r || r->field_count > TELEMETRY_FIELDS) {
		return -EINVAL;
	}
	for (size_t i = 0; i < r->field_count; i++) {
		if (!app_name_valid(r->fields[i].name) || r->fields[i].decimals > 6) {
			return -EINVAL;
		}
	}
	k_mutex_lock(&s->lock, K_FOREVER);
	if (s->count == s->capacity) {
		s->head = (s->head + 1) % CONFIG_APP_CACHE_SIZE;
		s->count--;
		s->dropped++;
	}
	size_t tail = (s->head + s->count++) % CONFIG_APP_CACHE_SIZE;
	s->records[tail] = *r;
	s->records[tail].sequence = s->next_sequence++;
	k_mutex_unlock(&s->lock);
	return 0;
}

size_t telemetry_store_peek(struct telemetry_store *s, struct telemetry_record *out, size_t max)
{
	k_mutex_lock(&s->lock, K_FOREVER);
	size_t n = MIN(s->count, max);
	for (size_t i = 0; i < n; i++) {
		out[i] = s->records[(s->head + i) % CONFIG_APP_CACHE_SIZE];
	}
	k_mutex_unlock(&s->lock);
	return n;
}

void telemetry_store_pop(struct telemetry_store *s, uint64_t ack)
{
	k_mutex_lock(&s->lock, K_FOREVER);
	while (s->count && s->records[s->head].sequence <= ack) {
		s->head = (s->head + 1) % CONFIG_APP_CACHE_SIZE;
		s->count--;
	}
	k_mutex_unlock(&s->lock);
}

size_t telemetry_store_count(struct telemetry_store *s)
{
	k_mutex_lock(&s->lock, K_FOREVER);
	size_t n = s->count;
	k_mutex_unlock(&s->lock);
	return n;
}

uint64_t telemetry_store_dropped(struct telemetry_store *s)
{
	k_mutex_lock(&s->lock, K_FOREVER);
	uint64_t n = s->dropped;
	k_mutex_unlock(&s->lock);
	return n;
}

void telemetry_store_clear(struct telemetry_store *s)
{
	k_mutex_lock(&s->lock, K_FOREVER);
	s->count = 0;
	s->head = 0;
	k_mutex_unlock(&s->lock);
}

int telemetry_store_resize(struct telemetry_store *s, size_t capacity)
{
	if (!capacity || capacity > CONFIG_APP_CACHE_SIZE) {
		return -EINVAL;
	}
	k_mutex_lock(&s->lock, K_FOREVER);
	while (s->count > capacity) {
		s->head = (s->head + 1) % CONFIG_APP_CACHE_SIZE;
		s->count--;
		s->dropped++;
	}
	s->capacity = capacity;
	k_mutex_unlock(&s->lock);
	return 0;
}
