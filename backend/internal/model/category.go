package model

import "gorm.io/gorm"

// Category stores a node in the product category tree.
type Category struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	ParentID  *int64         `gorm:"column:parent_id;index" json:"parentId"`
	Name      string         `gorm:"column:name;not null" json:"name"`
	Sort      int            `gorm:"column:sort;not null;default:0" json:"sort"`
	Enabled   bool           `gorm:"column:enabled;not null;default:true" json:"enabled"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	Children  []Category     `gorm:"-" json:"children,omitempty"`
}

func (Category) TableName() string { return "categories" }
