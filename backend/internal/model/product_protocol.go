package model

import "time"

// ProductProtocol 表示产品支持的传输协议与应用协议组合。
type ProductProtocol struct {
	ID                  int64     `gorm:"column:id;primaryKey" json:"id"`
	ProductID           int64     `gorm:"column:product_id;not null;index;uniqueIndex:idx_product_protocol" json:"productId"`
	TransportProtocol   string    `gorm:"column:transport_protocol;size:32;not null;uniqueIndex:idx_product_protocol" json:"transportProtocol"`
	ApplicationProtocol string    `gorm:"column:application_protocol;size:32;not null;uniqueIndex:idx_product_protocol" json:"applicationProtocol"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ProductProtocol) TableName() string { return "product_protocols" }
