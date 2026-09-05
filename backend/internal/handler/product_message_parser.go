package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductMessageParserHandler struct {
	*Handler
	svc service.ProductMessageParserServiceInterface
}

func NewProductMessageParserHandler(h *Handler, svc service.ProductMessageParserServiceInterface) *ProductMessageParserHandler {
	return &ProductMessageParserHandler{Handler: h, svc: svc}
}

// Get godoc
// @Summary Get product message parser
// @Description Gets the product message parser or its default template.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Success 200 {object} v1.ApiResponse[v1.ProductMessageParser] "Successful response"
// @Router /products/key/{productKey}/message-parser [get]
func (h *ProductMessageParserHandler) Get(c *gin.Context) {
	parser, isDefault, err := h.svc.Get(c, c.Param("productKey"))
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, toProductMessageParserResponse(c.Param("productKey"), parser, isDefault))
}

// Put godoc
// @Summary Save product message parser
// @Description Saves product message parser.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Param request body v1.UpsertProductMessageParserRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.ProductMessageParser] "Successful response"
// @Router /products/key/{productKey}/message-parser [put]
func (h *ProductMessageParserHandler) Put(c *gin.Context) {
	var req v1.UpsertProductMessageParserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	parser, err := h.svc.Save(c, c.Param("productKey"), req.Language, req.Script)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, toProductMessageParserResponse(c.Param("productKey"), parser, false))
}

// Execute godoc
// @Summary Execute product message parser
// @Description Executes product message parser.
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "Product key"
// @Param request body v1.ExecuteProductMessageParserRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.ExecuteProductMessageParserResponse] "Successful response"
// @Router /products/key/{productKey}/message-parser/execute [post]
func (h *ProductMessageParserHandler) Execute(c *gin.Context) {
	var req v1.ExecuteProductMessageParserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	result, err := h.svc.Execute(c, c.Param("productKey"), req)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, result)
}

func toProductMessageParserResponse(productKey string, parser *model.ProductMessageParser, isDefault bool) v1.ProductMessageParser {
	return v1.ProductMessageParser{
		ProductKey: productKey,
		Language:   parser.Language,
		Script:     parser.Script,
		IsDefault:  isDefault,
	}
}
