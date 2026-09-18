#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
用法:
  ./scripts/mqtt-publish.sh <host> <port> <username> <password> <productKey> <deviceKey> [payload.json]

示例:
  ./scripts/mqtt-publish.sh 127.0.0.1 1883 username password product_key device_key

可选参数:
  第 7 个参数为 JSON payload 文件路径；不传时使用内置温度示例
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" || $# -lt 6 ]]; then
  usage
  [[ $# -lt 6 ]] && exit 1 || exit 0
fi

if ! command -v mosquitto_pub >/dev/null 2>&1; then
  echo "错误：未找到 mosquitto_pub，请先安装 mosquitto-clients。" >&2
  exit 1
fi

MQTT_HOST="$1"
MQTT_PORT="$2"
MQTT_USERNAME="$3"
MQTT_PASSWORD="$4"
PRODUCT_KEY="$5"
DEVICE_KEY="$6"
MQTT_CLIENT_ID="0things-mqtt-publisher-$(date +%s)"
TOPIC="/sys/thing/property/post/${PRODUCT_KEY}/${DEVICE_KEY}"

if [[ -n "${7:-}" ]]; then
  PAYLOAD_FILE="$7"
  if [[ ! -f "$PAYLOAD_FILE" ]]; then
    echo "错误：找不到 payload 文件：$PAYLOAD_FILE" >&2
    exit 1
  fi
  PAYLOAD="$(<"$PAYLOAD_FILE")"
else
  NOW_MS="$(($(date +%s) * 1000))"
  PAYLOAD="{\"version\":\"1.0\",\"params\":{\"temperature\":{\"value\":25.5,\"time\":${NOW_MS}}}}"
fi

mosquitto_pub \
  -h "$MQTT_HOST" \
  -p "$MQTT_PORT" \
  -u "$MQTT_USERNAME" \
  -P "$MQTT_PASSWORD" \
  -i "$MQTT_CLIENT_ID" \
  -q 1 \
  -t "$TOPIC" \
  -m "$PAYLOAD"

echo "已推送到 $MQTT_HOST:$MQTT_PORT $TOPIC"
