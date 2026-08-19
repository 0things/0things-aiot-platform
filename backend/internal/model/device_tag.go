package model

import (
	"time"

	"gorm.io/gorm"
)

type DeviceTag struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	DeviceID  int64          `gorm:"column:device_id"`
	Key       string         `gorm:"column:key"`
	Value     string         `gorm:"column:value"`
	Source    string         `gorm:"column:source"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (DeviceTag) TableName() string { return "device_tags" }
