package handler

import (
	productTSLV1 "0things-backend/api/product_tsl/v1"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductTSLHandler struct {
	*Handler
	svc *service.ProductTSLService
}

func NewProductTSLHandler(h *Handler, svc *service.ProductTSLService) *ProductTSLHandler {
	return &ProductTSLHandler{Handler: h, svc: svc}
}

// Get godoc
// @Summary 获取产品物模型（TSL）
// @Schemes
// @Description 通过产品 ID 获取其物模型定义（TSL）
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Success 200 {object} productTSLV1.GetProductTSLResponse
// @Router /products/{id}/tsl [get]
func (h *ProductTSLHandler) Get(c *gin.Context) {
	tsl, err := h.svc.Get(c, c.Param("id"))
	if err != nil {
		deviceError(c, err)
		return
	}
	productID := int64(0)
	if tsl.ProductID != nil {
		productID = *tsl.ProductID
	}
	c.JSON(200, productTSLV1.GetProductTSLResponse{ProductTSL: productTSLV1.ProductTSL{ID: tsl.ID, ProductID: productID, TSL: tsl.TSL, CreatedAt: tsl.CreatedAt, UpdatedAt: tsl.UpdatedAt}})
}

// Put godoc
// @Summary 上传/更新产品物模型（TSL）
// @Schemes
// @Description 通过产品 ID 上传或更新物模型定义
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Param request body productTSLV1.UpsertProductTSLRequest true "params"
// @Success 200 {object} productTSLV1.SuccessResponse
// @Router /products/{id}/tsl [post]
// @Router /products/{id}/tsl [put]
func (h *ProductTSLHandler) Put(c *gin.Context) {
	var req productTSLV1.UpsertProductTSLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	if err := h.svc.Upsert(c, c.Param("id"), req.TSL); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productTSLV1.SuccessResponse{Success: true})
}

// Delete godoc
// @Summary 删除产品物模型（TSL）
// @Schemes
// @Description 通过产品 ID 删除其物模型定义
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Success 200 {object} productTSLV1.SuccessResponse
// @Router /products/{id}/tsl [delete]
func (h *ProductTSLHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("id")); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productTSLV1.SuccessResponse{Success: true})
}
