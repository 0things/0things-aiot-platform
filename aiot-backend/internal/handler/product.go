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
	svc *service.ProductService
}

func NewProductHandler(h *Handler, svc *service.ProductService) *ProductHandler {
	return &ProductHandler{Handler: h, svc: svc}
}

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

func (h *ProductHandler) GetByKey(c *gin.Context) {
	product, err := h.svc.GetByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, productV1.GetProductByKeyResponse{Product: productJSON(*product, 0)})
}

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
