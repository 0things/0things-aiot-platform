package handler

import (
	productV1 "aiot-backend/api/product/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func productJSON(product model.Product, count int64) productV1.Product {
	return productV1.Product{
		ID:                 product.ID,
		ProductKey:         product.ProductKey,
		Name:               product.Name,
		Description:        product.Description,
		CategoryID:         product.CategoryID,
		Status:             product.Status,
		DeviceCount:        int32(count),
		NodeType:           product.NodeType,
		ConnectivityMethod: product.ConnectivityMethod,
		AccessProtocol:     product.AccessProtocol,
		Protocols: func() []productV1.ProductProtocolInput {
			items := make([]productV1.ProductProtocolInput, len(product.Protocols))
			for i, protocol := range product.Protocols {
				items[i] = productV1.ProductProtocolInput{TransportProtocol: protocol.TransportProtocol, ApplicationProtocol: protocol.ApplicationProtocol}
			}
			return items
		}(),
		OrganizationID: product.OrganizationID,
		CreatedAt:      product.CreatedAt,
		UpdatedAt:      product.UpdatedAt,
		DeletedAt:      deletedAt(product.DeletedAt),
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
// @Success 200 {object} v1.ApiResponse[productV1.CreateProductResponse]
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req productV1.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	protocols := make([]model.ProductProtocol, len(req.Protocols))
	for i, item := range req.Protocols {
		protocols[i] = model.ProductProtocol{TransportProtocol: item.TransportProtocol, ApplicationProtocol: item.ApplicationProtocol}
	}
	if req.AccessProtocol == "default" {
		// default 是快捷组合，统一展开为 HTTP 和 MQTT 两个通用 JSON 接入端点。
		protocols = []model.ProductProtocol{
			{TransportProtocol: "http", ApplicationProtocol: "json"},
			{TransportProtocol: "mqtt", ApplicationProtocol: "json"},
		}
	}
	product, err := h.svc.Create(c, &model.Product{Name: req.Name, Description: req.Description, CategoryID: req.CategoryID, Status: req.Status, NodeType: req.NodeType, ConnectivityMethod: req.ConnectivityMethod, AccessProtocol: req.AccessProtocol, Protocols: protocols})
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productV1.CreateProductResponse{Product: productJSON(*product, 0)})
}

// Get godoc
// @Summary 通过 productKey 获取产品
// @Schemes
// @Description 通过产品 Key 获取产品详情
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} v1.ApiResponse[productV1.GetProductResponse]
// @Router /products/{productKey} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	product, err := h.svc.GetByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productV1.GetProductResponse{Product: productJSON(*product, 0)})
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
// @Success 200 {object} v1.ApiResponse[productV1.ListProductsResponse]
// @Router /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	pageNumber, pageSize := page(c, 10)
	products, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("category"), c.Query("status"), c.Query("searchText"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]productV1.ProductListItem, 0, len(products))
	for _, product := range products {
		item := productV1.ProductListItem{Product: productJSON(product.Product, 0), CategoryName: product.CategoryName}
		items = append(items, item)
	}
	v1.HandleSuccess(c, productV1.ListProductsResponse{Products: items, Total: total, Page: pageNumber, PageSize: pageSize})
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
// @Success 200 {object} v1.ApiResponse[productV1.UpdateProductResponse]
// @Router /products/{productKey} [put]
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
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		product.Status = req.Status
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
	if req.Protocols != nil {
		product.Protocols = make([]model.ProductProtocol, len(req.Protocols))
		for i, item := range req.Protocols {
			product.Protocols[i] = model.ProductProtocol{TransportProtocol: item.TransportProtocol, ApplicationProtocol: item.ApplicationProtocol}
		}
	}
	if req.AccessProtocol == "default" {
		// default 是快捷组合，统一展开为 HTTP 和 MQTT 两个通用 JSON 接入端点。
		product.Protocols = []model.ProductProtocol{
			{TransportProtocol: "http", ApplicationProtocol: "json"},
			{TransportProtocol: "mqtt", ApplicationProtocol: "json"},
		}
	}
	if err = h.svc.Save(c, product); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productV1.UpdateProductResponse{Product: productJSON(*product, 0)})
}

// Delete godoc
// @Summary 删除产品
// @Schemes
// @Description 通过产品 Key 软删除产品
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} v1.ApiResponse[productV1.SuccessResponse]
// @Router /products/{productKey} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	err := h.svc.DeleteByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productV1.SuccessResponse{Success: true})
}

// Restore 保留旧调用兼容，产品恢复路由不再注册。
func (h *ProductHandler) Restore(c *gin.Context) {
	product, err := h.svc.RestoreByKey(c, c.Param("productKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, productV1.RestoreProductResponse{Product: productJSON(*product, 0)})
}
