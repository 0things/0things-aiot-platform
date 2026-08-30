package enum

// 统一管理 0things MQTT 传输层主题契约常量与格式模板
const (
	// --- 上行订阅通配符主题 (Uplink Wildcard Topics) ---

	// MQTTSubTelemetry 遥测与属性上报订阅主题: /sys/{deviceKey}/thing/event/property/post
	MQTTSubTelemetry = "/sys/+/thing/event/property/post"

	// MQTTSubOTAProgress OTA 固件升级进度上报订阅主题: /sys/{deviceKey}/ota/device/progress
	MQTTSubOTAProgress = "/sys/+/ota/device/progress"

	// MQTTSubEvent 自定义业务事件上报订阅主题: /sys/{deviceKey}/thing/event/{eventName}/post
	MQTTSubEvent = "/sys/+/thing/event/+/post"

	// --- 下行控制发布模板 (Downlink Topic Templates) ---

	// MQTTTplPropertySet 属性设置下发模板，格式化参数: deviceKey
	MQTTTplPropertySet = "/sys/%s/thing/service/property/set"

	// MQTTTplOTAUpgrade OTA 升级通知下发模板，格式化参数: deviceKey
	MQTTTplOTAUpgrade = "/sys/%s/ota/device/upgrade"
)
