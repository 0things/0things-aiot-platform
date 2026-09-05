package dto

import "aiot-backend/internal/model"

// ProductListItem 是产品列表查询结果，包含关联分类名称，不属于领域模型。
type ProductListItem struct {
	model.Product
	CategoryName string
}

// ProductOption represents lightweight product summary for dropdown selects.
type ProductOption struct {
	ID         int64
	ProductKey string
	Name       string
	NodeType   string
}
