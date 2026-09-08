package model

import "time"

// Device represents the device table entity.
type Device struct {
	ID             int64     `gorm:"primaryKey"`
	DeviceKey      string    `gorm:"column:device_key;uniqueIndex"`
	Name           string    `gorm:"column:name"`
	ProductID      int64     `gorm:"column:product_id"`
	OrganizationID int64     `gorm:"column:organization_id;not null;default:1"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Device) TableName() string { return "devices" }
