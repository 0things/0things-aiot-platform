package enum

// 统一管理 0things 消息队列主题与消费组命名契约
const (
	// KafkaTopicDeviceTelemetry 设备时序遥测与属性上报主题 (由 rule-engine 消费)
	KafkaTopicDeviceTelemetry = "device.telemetry.v1"

	// ConsumerGroupRuleEngine rule-engine 专属消费组
	ConsumerGroupRuleEngine = "rule-engine-consumer-group"
)
