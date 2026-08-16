package model

import (
	"time"
)

type DeviceShadow struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	DeviceID  int64     `gorm:"column:device_id"`
	Desired   string    `gorm:"column:desired;type:text"`
	Reported  string    `gorm:"column:reported;type:text"`
	Metadata  string    `gorm:"column:metadata;type:text"`
	Version   int64     `gorm:"column:version"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (DeviceShadow) TableName() string { return "device_shadows" }

type DeviceShadowHistory struct {
	ID        int64     `gorm:"primaryKey"`
	DeviceID  int64     `gorm:"index;not null"`
	Version   int64     `gorm:"not null"`
	Source    string    `gorm:"size:16;not null"`
	Desired   string    `gorm:"type:text"`
	Reported  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index;not null"`
}

func (DeviceShadowHistory) TableName() string { return "device_shadow_histories" }
