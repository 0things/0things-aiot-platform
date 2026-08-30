package enum

// 统一管理 0things 消息队列主题与消费组命名契约
const (
	// KafkaTopicDeviceTelemetry 设备时序遥测与属性上报主题 (由 rule-engine 消费)
	KafkaTopicDeviceTelemetry = "device.telemetry.v1"

	// KafkaTopicOTAReport 设备 OTA 固件升级进度汇报主题 (由 mq-service 消费)
	KafkaTopicOTAReport = "ota.report.v1"

	// KafkaTopicDeviceEvent 设备生命周期/异常告警事件主题 (由 mq-service 消费)
	KafkaTopicDeviceEvent = "device.event.v1"

	// KafkaTopicDeviceCommand 云端下发控制/OTA 指令主题 (由各 transport 消费)
	KafkaTopicDeviceCommand = "device.command.v1"

	// ConsumerGroupMQServiceOTA mq-service 专职 OTA 任务消费组
	ConsumerGroupMQServiceOTA = "mq-service-ota-consumer-group"
)
