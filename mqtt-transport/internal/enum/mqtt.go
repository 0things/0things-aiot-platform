package enum

// 统一管理 0things MQTT 传输层主题契约常量与格式模板
const (
	// --- 上行订阅通配符主题 (Uplink Wildcard Topics) ---

	// MQTTSubTelemetry 遥测与属性上报订阅主题: /sys/{productKey}/{deviceKey}/thing/event/property/post
	MQTTSubTelemetry = "/sys/+/+/thing/event/property/post"

	// MQTTSubOTAProgress OTA 固件升级进度上报订阅主题: /sys/{productKey}/{deviceKey}/ota/device/progress
	MQTTSubOTAProgress = "/sys/+/+/ota/device/progress"
	// MQTTSubOTAInform OTA 固件版本/升级上报订阅主题: /ota/device/inform/{productKey}/{deviceKey}
	MQTTSubOTAInform = "/ota/device/inform/+/+"
	// MQTTSubOTAProgressV1 OTA 升级进度上报订阅主题: /ota/device/progress/{productKey}/{deviceKey}
	MQTTSubOTAProgressV1 = "/ota/device/progress/+/+"

	// MQTTSubEvent 自定义业务事件上报订阅主题: /sys/{productKey}/{deviceKey}/thing/event/{eventName}/post
	MQTTSubEvent = "/sys/+/+/thing/event/+/post"

	// --- 下行主题模板 (Downlink Topic Templates) ---

	// MQTTTplPropertySet 属性设置下发主题模板: /sys/{productKey}/{deviceKey}/thing/service/property/set
	MQTTTplPropertySet = "/sys/%s/%s/thing/service/property/set"

	// MQTTTplOTAUpgrade OTA 升级指令下发主题模板: /sys/{productKey}/{deviceKey}/ota/device/upgrade
	MQTTTplOTAUpgrade = "/sys/%s/%s/ota/device/upgrade"
	// MQTTTplOTAUpgradeV1 OTA 升级指令下发主题模板: /ota/device/upgrade/{productKey}/{deviceKey}
	MQTTTplOTAUpgradeV1 = "/ota/device/upgrade/%s/%s"
)
