package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// RuleChain stores a versioned visual rule graph owned by one organization.
type RuleChain struct {
	ID             int64           `gorm:"column:id;primaryKey" json:"id"`
	UUID           string          `gorm:"column:uuid;type:varchar(36);not null;uniqueIndex:uk_rule_chains_uuid" json:"uuid"`
	OrganizationID int64           `gorm:"column:organization_id;not null;index:idx_rule_chains_org" json:"organizationId"`
	Name           string          `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description    string          `gorm:"column:description;type:text" json:"description"`
	Status         string          `gorm:"column:status;type:varchar(16);not null;default:draft;index:idx_rule_chains_org_status" json:"status"`
	Version        int             `gorm:"column:version;not null;default:1" json:"version"`
	Graph          json.RawMessage `gorm:"column:graph;type:json;not null" json:"graph"`
	CreatedAt      time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"column:deleted_at" json:"-"`
}

func (RuleChain) TableName() string { return "rule_chains" }
