package model

import (
	"time"

	"gorm.io/gorm"
)

type OTAPackage struct {
	ID             int64  `gorm:"primaryKey"`
	UUID           string `gorm:"column:uuid;type:varchar(64);uniqueIndex"`
	PackageName    string `gorm:"column:package_name"`
	Version        string
	ProductID      int64  `gorm:"column:product_id"`
	OrganizationID int64  `gorm:"column:organization_id;not null;default:1"`
	ProductKey     string `gorm:"column:product_key;->" json:"productKey,omitempty"`
	ProductName    string `gorm:"column:product_name;->" json:"productName,omitempty"`
	PackageType    string `gorm:"column:package_type"`
	Status         string
	UploadType     string `gorm:"column:upload_type"`
	FileURL        string `gorm:"column:file_url"`
	FileSize       int64  `gorm:"column:file_size"`
	Checksum       string
	Description    string
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ReleasedAt     *time.Time `gorm:"column:released_at"`
}

func (OTAPackage) TableName() string { return "ota_packages" }

type UpgradeBatch struct {
	ID                int64  `gorm:"primaryKey"`
	BatchID           string `gorm:"column:batch_id"`
	OTAPackageID      string `gorm:"column:ota_package_id"`
	BatchName         string `gorm:"column:batch_name"`
	UpgradeStrategy   string `gorm:"column:upgrade_strategy"`
	Status            string
	TargetDeviceCount int32 `gorm:"column:target_device_count"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (UpgradeBatch) TableName() string { return "ota_upgrade_batches" }

type DeviceUpgradeStatus struct {
	ID                   int64  `gorm:"primaryKey"`
	DeviceID             int64  `gorm:"column:device_id;uniqueIndex:ux_ota_batch_device"`
	OTAPackageID         string `gorm:"column:ota_package_id"`
	UpgradeBatchID       string `gorm:"column:upgrade_batch_id;uniqueIndex:ux_ota_batch_device;index:idx_ota_batch_status"`
	Status               string `gorm:"index:idx_ota_batch_status"`
	TargetVersion        string `gorm:"column:target_version"`
	Progress             int32  `gorm:"column:progress"`
	DispatchAttempts     int32  `gorm:"column:dispatch_attempts"`
	LastDispatchError    string `gorm:"column:last_dispatch_error"`
	FirstProgressAt      *int64 `gorm:"column:first_progress_at"`
	LastReportAt         *int64 `gorm:"column:last_report_at"`
	TimeoutSeconds       int32  `gorm:"column:timeout_seconds"`
	MaxRetries           int32  `gorm:"column:max_retries"`
	CurrentVersion       string `gorm:"column:current_version"`
	LastStatusChangeTime *int64 `gorm:"column:last_status_change_ts"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (DeviceUpgradeStatus) TableName() string { return "ota_device_upgrade_status" }

// DeviceDeployment is the joined OTA deployment projection used by the
// repository and service layers. It is not a persisted table model.
type DeviceDeployment struct {
	DeviceID             int64     `gorm:"column:device_id"`
	DeviceKey            string    `gorm:"column:device_key"`
	DeviceName           string    `gorm:"column:device_name"`
	ProductID            int64     `gorm:"column:product_id"`
	ProductKey           string    `gorm:"column:product_key"`
	CurrentVersion       string    `gorm:"column:current_version"`
	TargetVersion        string    `gorm:"column:target_version"`
	Progress             int32     `gorm:"column:progress"`
	UpgradeBatchID       string    `gorm:"column:upgrade_batch_id"`
	Status               string    `gorm:"column:status"`
	LastStatusChangeTime int64     `gorm:"column:last_status_change_ts"`
	CreatedAt            time.Time `gorm:"column:created_at"`
}
