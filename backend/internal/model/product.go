package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                 int64          `gorm:"column:id;primaryKey" json:"id"`
	ProductKey         string         `gorm:"column:product_key" json:"productKey"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	Category           string         `json:"category"`
	Status             string         `json:"status"`
	Metadata           string         `gorm:"type:text" json:"metadata"`
	NodeType           string          `gorm:"column:node_type" json:"nodeType"`
	ConnectivityMethod string          `gorm:"column:connectivity_method" json:"connectivityMethod"`
	AccessProtocol     string          `gorm:"column:access_protocol" json:"accessProtocol"`
	TenantID           int64           `gorm:"column:tenant_id;not null;default:1" json:"tenantId"`
	DeletedAt          gorm.DeletedAt  `gorm:"column:deleted_at" json:"-"`
	CreatedAt          time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

func (Product) TableName() string { return "products" }
