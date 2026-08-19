package model

import "time"

// ProductMessageParser is the single active message parser configuration of a product.
type ProductMessageParser struct {
	ID        int64  `gorm:"primaryKey"`
	TenantID  int64  `gorm:"column:tenant_id;not null;uniqueIndex:idx_product_message_parser_tenant_product"`
	ProductID int64  `gorm:"column:product_id;not null;uniqueIndex:idx_product_message_parser_tenant_product"`
	Language  string `gorm:"column:language;not null"`
	Script    string `gorm:"column:script;type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ProductMessageParser) TableName() string { return "product_message_parsers" }
