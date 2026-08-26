package enum

const (
	// KafkaTopicOTAUpgradeCommandV1 是后端发送给 OTA 消费者的升级指令 Topic。
	KafkaTopicOTAUpgradeCommandV1 = "ota.upgrade.command.v1"
	// KafkaTopicOTAUpgradeReportV1 是设备升级进度和结果回报 Topic。
	KafkaTopicOTAUpgradeReportV1 = "ota.upgrade.report.v1"
)
