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

func (h *ProductTSLHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("id")); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productTSLV1.SuccessResponse{Success: true})
}
