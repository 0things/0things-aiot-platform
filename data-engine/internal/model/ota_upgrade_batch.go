package model

import "time"

// UpgradeBatch represents an OTA batch record in database.
type UpgradeBatch struct {
	ID                int64     `gorm:"primaryKey"`
	BatchID           string    `gorm:"column:batch_id;uniqueIndex"`
	OTAPackageID      string    `gorm:"column:ota_package_id"`
	BatchName         string    `gorm:"column:batch_name"`
	UpgradeStrategy   string    `gorm:"column:upgrade_strategy"`
	Status            string    `gorm:"column:status"`
	TargetDeviceCount int32     `gorm:"column:target_device_count"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (UpgradeBatch) TableName() string { return "ota_upgrade_batches" }
