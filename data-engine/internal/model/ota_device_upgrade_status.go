package model

import "time"

// DeviceUpgradeStatus represents the per-device OTA upgrade progress status record.
type DeviceUpgradeStatus struct {
	ID                   int64     `gorm:"primaryKey"`
	DeviceID             int64     `gorm:"column:device_id;uniqueIndex:ux_ota_batch_device"`
	OTAPackageID         string    `gorm:"column:ota_package_id"`
	UpgradeBatchID       string    `gorm:"column:upgrade_batch_id;uniqueIndex:ux_ota_batch_device;index:idx_ota_batch_status"`
	Status               string    `gorm:"column:status;index:idx_ota_batch_status"`
	Module               string    `gorm:"column:module;not null;default:default"`
	TargetVersion        string    `gorm:"column:target_version"`
	Progress             int32     `gorm:"column:progress"`
	DispatchAttempts     int32     `gorm:"column:dispatch_attempts"`
	LastDispatchError    string    `gorm:"column:last_dispatch_error"`
	FirstProgressAt      *int64    `gorm:"column:first_progress_at"`
	LastReportAt         *int64    `gorm:"column:last_report_at"`
	TimeoutSeconds       int32     `gorm:"column:timeout_seconds"`
	MaxRetries           int32     `gorm:"column:max_retries"`
	CurrentVersion       string    `gorm:"column:current_version"`
	LastStatusChangeTime *int64    `gorm:"column:last_status_change_ts"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (DeviceUpgradeStatus) TableName() string { return "ota_device_upgrade_status" }
