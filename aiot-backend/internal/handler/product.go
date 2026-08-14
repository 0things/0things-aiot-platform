package handler

import (
	"fmt"

	productV1 "0things-backend/api/product/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func productJSON(product model.Product, count int64) productV1.Product {
	return productV1.Product{
		ID:                 product.ID,
		ProductKey:         product.ProductKey,
		Name:               product.Name,
		Description:        product.Description,
		Category:           product.Category,
		Status:             product.Status,
		Metadata:           fmt.Sprint(raw(product.Metadata)),
		DeviceCount:        int32(count),
		NodeType:           product.NodeType,
		ConnectivityMethod: product.ConnectivityMethod,
		AccessProtocol:     product.AccessProtocol,
		TenantID:           product.TenantID,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
		DeletedAt:          deletedAt(product.DeletedAt),
	}
}

type ProductHandler struct {
	*Handler
	svc service.ProductServiceInterface
}

func NewProductHandler(h *Handler, svc service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{Handler: h, svc: svc}
}

// Create godoc
// @Summary 创建产品
// @Schemes
// @Description 创建一个新的产品
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body productV1.CreateProductRequest true "params"
// @Success 200 {object} productV1.CreateProductResponse
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req productV1.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	product, err := h.svc.Create(c, &model.Product{Name: req.Name, Description: req.Description, Category: req.Category, Status: req.Status, Metadata: req.Metadata, NodeType: req.NodeType, ConnectivityMethod: req.ConnectivityMethod, AccessProtocol: req.AccessProtocol})
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.CreateProductResponse{Product: productJSON(*product, 0)})
}

// Get godoc
// @Summary 通过 ID 获取产品
// @Schemes
// @Description 通过产品 ID 获取产品详情
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Success 200 {object} productV1.GetProductResponse
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	productID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	product, err := h.svc.Get(c, productID)
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.GetProductResponse{Product: productJSON(*product, 0)})
}

// GetByKey godoc
// @Summary 通过 productKey 获取产品
// @Schemes
// @Description 通过产品 Key 获取产品详情
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} productV1.GetProductByKeyResponse
// @Router /products/key/{productKey} [get]
func (h *ProductHandler) GetByKey(c *gin.Context) {
	product, err := h.svc.GetByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.GetProductByKeyResponse{Product: productJSON(*product, 0)})
}

// List godoc
// @Summary 获取产品列表
// @Schemes
// @Description 分页获取产品列表，支持按 category、status、searchText 过滤
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param category query string false "产品分类"
// @Param status query string false "产品状态"
// @Param searchText query string false "搜索关键字"
// @Success 200 {object} productV1.ListProductsResponse
// @Router /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	pageNumber, pageSize := page(c, 10)
	products, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("category"), c.Query("status"), c.Query("searchText"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]productV1.Product, 0, len(products))
	for _, product := range products {
		items = append(items, productJSON(product, 0))
	}
	c.JSON(200, productV1.ListProductsResponse{Products: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

// Update godoc
// @Summary 更新产品
// @Schemes
// @Description 通过 productKey 更新产品
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Param request body productV1.UpdateProductRequest true "params"
// @Success 200 {object} productV1.UpdateProductResponse
// @Router /products/key/{productKey} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	product, err := h.svc.GetByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	var req productV1.UpdateProductRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if len(req.Metadata) > 0 {
		product.Metadata = req.Metadata
	}
	if req.NodeType != "" {
		product.NodeType = req.NodeType
	}
	if req.ConnectivityMethod != "" {
		product.ConnectivityMethod = req.ConnectivityMethod
	}
	if req.AccessProtocol != "" {
		product.AccessProtocol = req.AccessProtocol
	}
	if err = h.svc.Save(c, product); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.UpdateProductResponse{Product: productJSON(*product, 0)})
}

// Delete godoc
// @Summary 删除产品
// @Schemes
// @Description 通过产品 ID 软删除产品
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Success 200 {object} productV1.SuccessResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	productID, err := id(c)
	if err == nil {
		err = h.svc.Delete(c, productID)
	}
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.SuccessResponse{Success: true})
}

// Restore godoc
// @Summary 恢复已删除的产品
// @Schemes
// @Description 通过产品 ID 恢复软删除的产品
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "产品 ID"
// @Success 200 {object} productV1.RestoreProductResponse
// @Router /products/{id}/restore [post]
func (h *ProductHandler) Restore(c *gin.Context) {
	productID, err := id(c)
	if err == nil {
		var product *model.Product
		product, err = h.svc.Restore(c, productID)
		if err == nil {
			c.JSON(200, productV1.RestoreProductResponse{Product: productJSON(*product, 0)})
			return
		}
	}
	deviceError(c, err)
}
