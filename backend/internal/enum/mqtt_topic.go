package enum

const (
	// MQTTTopicOTADeviceUpgrade 是后端向指定设备下发 OTA 指令的 Topic 模板。
	MQTTTopicOTADeviceUpgrade = "/ota/device/upgrade/%s/%s"
	// MQTTTopicOTADeviceProgress 是设备上报升级进度的订阅 Topic 模板。
	MQTTTopicOTADeviceProgress = "/ota/device/progress/+/+"
	// MQTTTopicOTADeviceInform 是设备上报升级结果的订阅 Topic 模板。
	MQTTTopicOTADeviceInform = "/ota/device/inform/+/+"
)
