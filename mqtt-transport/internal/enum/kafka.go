package enum

// 统一管理 0things 消息队列主题与消费组命名契约
const (
	KafkaTopicDeviceTelemetry = "device.telemetry.v1"
	KafkaTopicOTAReport       = "ota.report.v1"
	KafkaTopicDeviceEvent     = "device.event.v1"
	KafkaTopicDeviceCommand   = "device.command.v1"

	ConsumerGroupMqttDownlink = "mqtt-transport-downlink-group"
)
