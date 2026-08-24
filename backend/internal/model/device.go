package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID             int64          `gorm:"column:id;primaryKey" json:"id"`
	DeviceKey      string         `gorm:"column:device_key" json:"deviceKey"`
	Name           string         `json:"name"`
	ProductID      int64          `gorm:"column:product_id" json:"productId"`
	OrganizationID int64          `gorm:"column:organization_id;not null;default:1" json:"organizationId"`
	Enabled        bool           `json:"enabled"`
	Metadata       string         `gorm:"type:text" json:"metadata"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	Product        Product        `gorm:"foreignKey:ProductID" json:"-"`
	State          DeviceState    `gorm:"foreignKey:DeviceKey;references:DeviceKey" json:"-"`
}

func (Device) TableName() string { return "devices" }

type DeviceState struct {
	ID              int64  `gorm:"column:id;primaryKey"`
	DeviceKey       string `gorm:"column:device_key"`
	State           string `gorm:"column:state"`
	LastOnlineTime  *int64 `gorm:"column:last_online_time"`
	LastOfflineTime *int64 `gorm:"column:last_offline_time"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (DeviceState) TableName() string { return "device_states" }
