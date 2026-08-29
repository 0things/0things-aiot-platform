package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	// DeviceGroupTypeManual 手动维护成员的分组类型
	DeviceGroupTypeManual = "manual"
	// DeviceGroupTypeDynamic 基于规则实时匹配设备的分组类型
	DeviceGroupTypeDynamic = "dynamic"
)

// DeviceGroup 保存组织内的手动或动态设备分组。
type DeviceGroup struct {
	ID             int64          `gorm:"column:id;primaryKey" json:"-"`
	GroupUUID      string         `gorm:"column:group_uuid;size:36;not null;uniqueIndex" json:"groupUuid"`
	OrganizationID int64          `gorm:"column:organization_id;not null;index" json:"-"`
	Name           string         `gorm:"column:name;size:128;not null" json:"name"`
	Type           string         `gorm:"column:type;size:16;not null" json:"type"`
	Description    string         `gorm:"column:description;type:text" json:"description"`
	Rule           string         `gorm:"column:rule;type:text" json:"rule,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeviceGroup) TableName() string { return "device_groups" }
