package model

import "time"

// DeviceGroupMember 保存手动分组与设备的内部关联。
type DeviceGroupMember struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"-"`
	GroupID   int64     `gorm:"column:group_id;not null;uniqueIndex:idx_group_device" json:"-"`
	DeviceID  int64     `gorm:"column:device_id;not null;uniqueIndex:idx_group_device" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeviceGroupMember) TableName() string { return "device_group_members" }
