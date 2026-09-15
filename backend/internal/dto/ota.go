package dto

import "time"

// DeviceDeployment represents the joined OTA deployment projection used by repository and service layers.
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
