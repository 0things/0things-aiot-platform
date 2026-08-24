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
}//@name ProductCreateProductRequest
type UpdateProductRequest struct {
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Category           string          `json:"category"`
	Status             string          `json:"status"`
	Metadata           json.RawMessage `json:"metadata"`
	NodeType           string          `json:"nodeType"`
	ConnectivityMethod string          `json:"connectivityMethod"`
	AccessProtocol     string          `json:"accessProtocol"`
}//@name ProductUpdateProductRequest
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
	OrganizationID           int64      `json:"organizationId"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}//@name Product
type CreateProductResponse struct {
	Product Product `json:"product"`
}//@name ProductCreateProductResponse
type GetProductResponse struct {
	Product Product `json:"product"`
}//@name ProductGetProductResponse
type GetProductByKeyResponse struct {
	Product Product `json:"product"`
}//@name ProductGetProductByKeyResponse
type UpdateProductResponse struct {
	Product Product `json:"product"`
}//@name ProductUpdateProductResponse
type RestoreProductResponse struct {
	Product Product `json:"product"`
}//@name ProductRestoreProductResponse
type ListProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}//@name ProductListProductsResponse

type SuccessResponse struct {
	Success bool `json:"success"`
}//@name ProductSuccessResponse
