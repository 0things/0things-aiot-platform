package kafka

// Kafka topic constants for the telemetry service.
const (
	// TopicTelemetry is the topic for IoT telemetry data.
	TopicTelemetry = "iotTelemetry"
	// GroupTelemetry is the consumer group for telemetry topic.
	GroupTelemetry = "telemetry-consumer"

	// TopicEvent is the topic for IoT device events.
	TopicEvent = "iotEvent"
	// GroupEvent is the consumer group for event topic.
	GroupEvent = "event-consumer"

	// TopicClientConnectState is the topic for client connection state events.
	TopicClientConnectState = "iotClientConnectStateEvent"
	// GroupClientConnectState is the consumer group for client connect state topic.
	GroupClientConnectState = "client-state-consumer"
)
