// Package productv1 owns the product HTTP contract. It mirrors
// device-service/api/product/v1/product.proto.
package productv1

import (
	"encoding/json"
	"time"
)

type CreateProductRequest struct {
	Name               string          `json:"name" binding:"required"`
	Description        string          `json:"description"`
	Category           string          `json:"category"`
	Status             string          `json:"status"`
	Metadata           json.RawMessage `json:"metadata"`
	NodeType           string          `json:"nodeType"`
	ConnectivityMethod string          `json:"connectivityMethod"`
	AccessProtocol     string          `json:"accessProtocol"`
}
type UpdateProductRequest struct {
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Category           string          `json:"category"`
	Status             string          `json:"status"`
	Metadata           json.RawMessage `json:"metadata"`
	NodeType           string          `json:"nodeType"`
	ConnectivityMethod string          `json:"connectivityMethod"`
	AccessProtocol     string          `json:"accessProtocol"`
}
type Product struct {
	ID                 int64      `json:"id"`
	ProductKey         string     `json:"productKey"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Category           string     `json:"category"`
	Status             string     `json:"status"`
	Metadata           string     `json:"metadata"`
	DeviceCount        int32      `json:"deviceCount"`
	NodeType           string     `json:"nodeType"`
	ConnectivityMethod string     `json:"connectivityMethod"`
	AccessProtocol     string     `json:"accessProtocol"`
	TenantID           int64      `json:"tenantId"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}
type CreateProductResponse struct {
	Product Product `json:"product"`
}
type GetProductResponse struct {
	Product Product `json:"product"`
}
type GetProductByKeyResponse struct {
	Product Product `json:"product"`
}
type UpdateProductResponse struct {
	Product Product `json:"product"`
}
type RestoreProductResponse struct {
	Product Product `json:"product"`
}
type ListProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
