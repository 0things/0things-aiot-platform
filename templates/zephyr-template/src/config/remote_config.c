#include <app/config.h>
#include <zephyr/data/json.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
/**
 * @brief Helper to check whether a JSON key matches a literal string.
 */
static bool key_is(const struct json_obj_key_value *kv, const char *name)
{
	return strlen(name) == kv->key_len && !memcmp(kv->key, name, kv->key_len);
}

/**
 * @brief Parse a JSON number token into a signed 64-bit integer.
 *
 * @param t JSON token.
 * @param out Output parsed integer value.
 * @return 0 on success, or -EINVAL on invalid format/overflow.
 */
static int integer(const struct json_token *t, int64_t *out)
{
	char text[24];
	size_t n = t->end - t->start;
	if (t->type != JSON_TOK_NUMBER || !n || n >= sizeof(text)) {
		return -EINVAL;
	}
	memcpy(text, t->start, n);
	text[n] = 0;
	errno = 0;
	char *end;
	*out = strtoll(text, &end, 10);
	return errno || *end ? -EINVAL : 0;
}

/**
 * @brief Parse key-value properties of an individual collector from JSON.
 *
 * @param c Output collector configuration struct to update.
 * @param token JSON object token containing collector properties.
 * @return 0 on success, or negative error code on failure.
 */
static int parse_collector(struct collector_config *c, struct json_token *token)
{
	struct json_obj obj;
	struct json_obj_key_value kv;
	unsigned int seen = 0;
	int rc = json_obj_separate_parse_init(&obj, token->start, token->end - token->start);
	if (rc) {
		return rc;
	}
	while (!(rc = json_obj_next_key_value(&obj, &kv)) && kv.key) {
		static const char *const names[] = {"enabled",
						    "collection_interval_sec",
						    "upload_interval_sec",
						    "cache_size",
						    "retry",
						    "priority"};
		size_t i;
		for (i = 0; i < ARRAY_SIZE(names); i++) {
			if (key_is(&kv, names[i])) {
				break;
			}
		}
		if (i == ARRAY_SIZE(names) || (seen & BIT(i))) {
			return -EINVAL;
		}
		seen |= BIT(i);
		if (!i) {
			if (kv.value.type != JSON_TOK_TRUE && kv.value.type != JSON_TOK_FALSE) {
				return -EINVAL;
			}
			c->enabled = kv.value.type == JSON_TOK_TRUE;
		} else {
			int64_t n;
			if (integer(&kv.value, &n) || n < 0 || n > INT32_MAX) {
				return -EINVAL;
			}
			switch (i) {
			case 1:
				c->collection_interval_sec = n;
				break;
			case 2:
				c->upload_interval_sec = n;
				break;
			case 3:
				c->cache_size = n;
				break;
			case 4:
				c->retry = n;
				break;
			case 5:
				c->priority = n;
				break;
			}
		}
	}
	return rc;
}

int remote_config_apply(char *json, size_t length)
{
	if (!json || !length || length >= CONFIG_APP_PAYLOAD_SIZE) {
		return -EMSGSIZE;
	}
	struct app_config next;
	config_manager_snapshot(&next);
	struct json_obj root;
	struct json_obj_key_value kv;
	unsigned int seen = 0;
	int rc = json_obj_separate_parse_init(&root, json, length);
	if (rc) {
		return rc;
	}
	while (!(rc = json_obj_next_key_value(&root, &kv)) && kv.key) {
		if (key_is(&kv, "version")) {
			int64_t n;
			if ((seen & 1) || integer(&kv.value, &n) || n < 1 || n > UINT32_MAX) {
				return -EINVAL;
			}
			next.version = n;
			seen |= 1;
		} else if (key_is(&kv, "collectors")) {
			if (seen & 2) {
				return -EINVAL;
			}
			seen |= 2;
			struct json_obj members;
			struct json_obj_key_value member;
			unsigned int updated = 0;
			rc = json_obj_separate_parse_init(&members, kv.value.start,
							  kv.value.end - kv.value.start);
			if (rc) {
				return rc;
			}
			while (!(rc = json_obj_next_key_value(&members, &member)) && member.key) {
				size_t i;
				for (i = 0; i < next.count; i++) {
					if (key_is(&member, next.collectors[i].name)) {
						break;
					}
				}
				if (i == next.count || (updated & BIT(i))) {
					return -EINVAL;
				}
				updated |= BIT(i);
				rc = parse_collector(&next.collectors[i], &member.value);
				if (rc) {
					return rc;
				}
			}
			if (rc || !updated) {
				return -EINVAL;
			}
		} else {
			return -EINVAL;
		}
	}
	if (rc || seen != 3) {
		return rc ? rc : -EINVAL;
	}
	/* The incremental decoder stops at the root close; reject a second JSON value. */
	for (char *p = root.lex.pos; p < json + length; p++) {
		if (*p != ' ' && *p != '\t' && *p != '\r' && *p != '\n') {
			return -EINVAL;
		}
	}
	return config_manager_apply(&next);
}
