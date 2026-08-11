package model

import "time"

type ProductTSL struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	TSL       string    `gorm:"column:tsl"`
	ProductID *int64    `gorm:"column:product_product_tsl"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ProductTSL) TableName() string { return "product_ts_ls" }
