package enum

// 统一管理 0things 消息队列主题与消费组命名契约
const (
	// KafkaTopicDeviceTelemetry 设备时序遥测与属性上报主题 (由 data-engine 消费)
	KafkaTopicDeviceTelemetry = "device.telemetry.v1"

	// KafkaTopicOTAReport 设备 OTA 固件升级进度汇报主题 (由 data-engine 消费)
	KafkaTopicOTAReport = "ota.report.v1"

	// KafkaTopicDeviceEvent 设备生命周期/异常告警事件主题 (由 data-engine 消费)
	KafkaTopicDeviceEvent = "device.event.v1"

	// 消费组默认名称
	ConsumerGroupTelemetry = "data-engine-telemetry-consumer-group"
	ConsumerGroupOTA       = "data-engine-ota-consumer-group"
	ConsumerGroupEvent     = "data-engine-event-consumer-group"
)
