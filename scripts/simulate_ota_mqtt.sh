#!/usr/bin/env bash

set -euo pipefail

# 模拟单台设备完整 OTA 流程：进度上报 -> 下载完成 -> 版本确认。
# 用法：./scripts/simulate_ota_mqtt.sh <product_key> <device_key> <batch_id> <target_version>
# 可通过 MQTT_HOST、MQTT_PORT、MQTT_INTERVAL 调整连接地址和上报间隔。

if [[ $# -ne 4 ]]; then
  echo "用法: $0 <product_key> <device_key> <batch_id> <target_version>" >&2
  exit 1
fi

PRODUCT_KEY="$1"
DEVICE_KEY="$2"
BATCH_ID="$3"
TARGET_VERSION="$4"
MQTT_HOST="${MQTT_HOST:-127.0.0.1}"
MQTT_PORT="${MQTT_PORT:-1883}"
MQTT_INTERVAL="${MQTT_INTERVAL:-1}"

PROGRESS_TOPIC="/ota/device/progress/${PRODUCT_KEY}/${DEVICE_KEY}"
INFORM_TOPIC="/ota/device/inform/${PRODUCT_KEY}/${DEVICE_KEY}"

command -v mosquitto_pub >/dev/null 2>&1 || {
  echo "未找到 mosquitto_pub，请先安装 Mosquitto 客户端。" >&2
  exit 1
}

publish() {
  local topic="$1"
  local payload="$2"
  echo "发布 ${topic}: ${payload}"
  mosquitto_pub -h "$MQTT_HOST" -p "$MQTT_PORT" -q 1 -t "$topic" -m "$payload"
}

echo "开始模拟设备 ${DEVICE_KEY} 的 OTA 批次 ${BATCH_ID}"

# 进度上报不能单独判定成功，最终成功必须由 inform 上报目标版本确认。
for step in 0 25 50 75 100; do
  publish "$PROGRESS_TOPIC" "{\"params\":{\"batch_id\":\"${BATCH_ID}\",\"device_key\":\"${DEVICE_KEY}\",\"status\":\"in_progress\",\"step\":${step},\"desc\":\"OTA progress ${step}%\"}}"
  sleep "$MQTT_INTERVAL"
done

# 设备重启后上报目标版本，服务端据此将任务标记为 success。
publish "$INFORM_TOPIC" "{\"params\":{\"batch_id\":\"${BATCH_ID}\",\"device_key\":\"${DEVICE_KEY}\",\"status\":\"success\",\"version\":\"${TARGET_VERSION}\",\"step\":100,\"desc\":\"OTA version confirmed\"}}"

echo "OTA 模拟完成"
