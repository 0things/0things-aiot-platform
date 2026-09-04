package v1

import "time"

type UpsertProductTSLRequest struct {
	TSL string `json:"tsl"`
} //@name ProductTslUpsertProductTSLRequest

type ProductTSLSuccessResponse struct {
	Success bool `json:"success"`
} //@name ProductTslSuccessResponse

type ProductTSL struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"productId"`
	TSL       string    `json:"tsl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} //@name ProductTslProductTSL

type GetProductTSLResponse struct {
	ProductTSL ProductTSL `json:"productTsl"`
} //@name ProductTslGetProductTSLResponse
