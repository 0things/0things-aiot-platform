package enum

const (
	// OTAStatusPending 表示任务已创建，等待下发。
	OTAStatusPending = "pending"
	// OTAStatusSent 表示升级指令已发送到 Kafka。
	OTAStatusSent = "sent"
	// OTAStatusInProgress 表示设备已开始执行升级。
	OTAStatusInProgress = "in_progress"
	// OTAStatusSuccess 表示设备已升级到目标版本。
	OTAStatusSuccess = "success"
	// OTAStatusFailed 表示设备升级失败。
	OTAStatusFailed = "failed"
	// OTAStatusTimeout 表示设备在规定时间内未完成升级。
	OTAStatusTimeout = "timeout"
	// OTAStatusCancelled 表示升级任务被用户取消。
	OTAStatusCancelled = "cancelled"
	// OTAPackageDeploying 表示升级包存在进行中的设备升级任务。
	OTAPackageDeploying = "deploying"
	// OTAPackagePartial 表示同一批次中部分设备成功、部分设备失败。
	OTAPackagePartial = "partial"
)
