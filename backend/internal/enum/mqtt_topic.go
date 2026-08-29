package enum

const (
	// MQTTTopicDevicePropertyPost 是设备属性上报 Topic 模板。
	MQTTTopicDevicePropertyPost = "/sys/%s/%s/thing/event/property/post"
	// MQTTTopicDevicePropertyPostWildcard 是网关订阅所有设备属性上报的 Topic。
	MQTTTopicDevicePropertyPostWildcard = "/sys/+/+/thing/event/property/post"
	// MQTTTopicDevicePropertySet 是云端向设备下发属性设置 Topic 模板。
	MQTTTopicDevicePropertySet = "/sys/%s/%s/thing/service/property/set"
	// MQTTTopicDeviceEventPost 是设备事件上报 Topic 模板，事件标识由设备物模型决定。
	MQTTTopicDeviceEventPost = "/sys/%s/%s/thing/event/%s/post"
	// MQTTTopicDeviceService 是云端向设备下发服务调用 Topic 模板。
	MQTTTopicDeviceService = "/sys/%s/%s/thing/service/%s"
	// MQTTTopicOTADeviceUpgrade 是后端向指定设备下发 OTA 指令的 Topic 模板。
	MQTTTopicOTADeviceUpgrade = "/ota/device/upgrade/%s/%s"
	// MQTTTopicOTADeviceProgress 是设备上报升级进度的订阅 Topic 模板。
	MQTTTopicOTADeviceProgress = "/ota/device/progress/+/+"
	// MQTTTopicOTADeviceInform 是设备上报升级结果的订阅 Topic 模板。
	MQTTTopicOTADeviceInform = "/ota/device/inform/+/+"
)
