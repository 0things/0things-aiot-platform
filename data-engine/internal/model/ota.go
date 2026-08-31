package model

import (
	"time"
)

// OTADeviceUpgradeReportEvent 表示单台设备的固件升级进度上报事实。
type OTADeviceUpgradeReportEvent struct {
	BatchID   string    `json:"batch_id,omitempty"` // 升级批次 ID
	DeviceKey string    `json:"device_key"`         // 设备全局唯一 Key
	Status    string    `json:"status"`             // 状态：UPGRADING(升级中) / SUCCESS(成功) / FAILED(失败)
	Version   string    `json:"version,omitempty"`  // 目标固件版本号（如 v2.0.0）
	Progress  int       `json:"progress"`           // 进度百分比 (0 ~ 100，或负数错误码)
	Desc      string    `json:"desc,omitempty"`     // 状态描述或错误原因
	Timestamp time.Time `json:"timestamp"`          // 上报时间
}
