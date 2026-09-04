package handler

import (
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
// @Summary Get product TSL
// @Schemes
// @Description Gets the thing specification language definition for a product.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Success 200 {object} v1.ApiResponse[v1.GetProductTSLResponse] "Successful response"
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
	v1.HandleSuccess(c, v1.GetProductTSLResponse{ProductTSL: v1.ProductTSL{ID: tsl.ID, ProductID: productID, TSL: tsl.TSL, CreatedAt: tsl.CreatedAt, UpdatedAt: tsl.UpdatedAt}})
}

// Put godoc
// @Summary Save product TSL
// @Schemes
// @Description Creates or updates the thing specification language definition for a product.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Param request body v1.UpsertProductTSLRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.ProductTSLSuccessResponse] "Successful response"
// @Router /products/{productKey}/tsl [post]
// @Router /products/{productKey}/tsl [put]
func (h *ProductTSLHandler) Put(c *gin.Context) {
	var req v1.UpsertProductTSLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	if err := h.svc.Upsert(c, c.Param("productKey"), req.TSL); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, v1.ProductTSLSuccessResponse{Success: true})
}

// Delete godoc
// @Summary Delete product TSL
// @Schemes
// @Description Deletes the thing specification language definition for a product.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Success 200 {object} v1.ApiResponse[v1.ProductTSLSuccessResponse] "Successful response"
// @Router /products/{productKey}/tsl [delete]
func (h *ProductTSLHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("productKey")); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, v1.ProductTSLSuccessResponse{Success: true})
}
