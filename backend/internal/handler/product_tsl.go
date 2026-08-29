package handler

import (
	productTSLV1 "aiot-backend/api/product_tsl/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductTSLHandler struct {
	*Handler
	svc service.ProductTSLServiceInterface
}

func NewProductTSLHandler(h *Handler, svc service.ProductTSLServiceInterface) *ProductTSLHandler {
	return &ProductTSLHandler{Handler: h, svc: svc}
}

// Get godoc
// @Summary 获取产品物模型（TSL）
// @Schemes
// @Description 通过产品 Key 获取其物模型定义（TSL）
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} v1.ApiResponse[productTSLV1.GetProductTSLResponse]
// @Router /products/{productKey}/tsl [get]
func (h *ProductTSLHandler) Get(c *gin.Context) {
	tsl, err := h.svc.Get(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	productID := int64(0)
	if tsl.ProductID != nil {
		productID = *tsl.ProductID
	}
	v1.HandleSuccess(c, productTSLV1.GetProductTSLResponse{ProductTSL: productTSLV1.ProductTSL{ID: tsl.ID, ProductID: productID, TSL: tsl.TSL, CreatedAt: tsl.CreatedAt, UpdatedAt: tsl.UpdatedAt}})
}

// Put godoc
// @Summary 上传/更新产品物模型（TSL）
// @Schemes
// @Description 通过产品 Key 上传或更新物模型定义
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Param request body productTSLV1.UpsertProductTSLRequest true "params"
// @Success 200 {object} v1.ApiResponse[productTSLV1.SuccessResponse]
// @Router /products/{productKey}/tsl [post]
// @Router /products/{productKey}/tsl [put]
func (h *ProductTSLHandler) Put(c *gin.Context) {
	var req productTSLV1.UpsertProductTSLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	if err := h.svc.Upsert(c, c.Param("productKey"), req.TSL); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productTSLV1.SuccessResponse{Success: true})
}

// Delete godoc
// @Summary 删除产品物模型（TSL）
// @Schemes
// @Description 通过产品 Key 删除其物模型定义
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} v1.ApiResponse[productTSLV1.SuccessResponse]
// @Router /products/{productKey}/tsl [delete]
func (h *ProductTSLHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("productKey")); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productTSLV1.SuccessResponse{Success: true})
}
