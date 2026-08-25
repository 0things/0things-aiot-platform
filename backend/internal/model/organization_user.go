package model

import (
	"time"
)

type OrganizationUser struct {
	Id             int64      `gorm:"primaryKey;autoIncrement"`
	OrganizationID int64      `gorm:"not null;uniqueIndex:idx_org_user"`
	UserID         string     `gorm:"not null;uniqueIndex:idx_org_user"`
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (OrganizationUser) TableName() string {
	return "organization_users"
}
