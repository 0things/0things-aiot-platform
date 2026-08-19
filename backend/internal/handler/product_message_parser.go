package handler

import (
	"net/http"

	messageParserV1 "aiot-backend/api/message_parser/v1"
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
// @Summary 获取产品消息解析器
// @Description 获取产品唯一的消息解析脚本；尚未配置时返回默认模板。
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Success 200 {object} v1.ApiResponse[messageParserV1.ProductMessageParser]
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
// @Summary 保存产品消息解析器
// @Description 按产品覆盖保存 JavaScript ES5 消息解析脚本。
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Param request body messageParserV1.UpsertProductMessageParserRequest true "params"
// @Success 200 {object} v1.ApiResponse[messageParserV1.ProductMessageParser]
// @Router /products/key/{productKey}/message-parser [put]
func (h *ProductMessageParserHandler) Put(c *gin.Context) {
	var req messageParserV1.UpsertProductMessageParserRequest
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
// @Summary 模拟执行产品消息解析器
// @Description 仅在受限 JavaScript ES5 运行时中执行保存的产品脚本。
// @Tags 产品模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param productKey path string true "产品 Key"
// @Param request body messageParserV1.ExecuteProductMessageParserRequest true "params"
// @Success 200 {object} v1.ApiResponse[messageParserV1.ExecuteProductMessageParserResponse]
// @Router /products/key/{productKey}/message-parser/execute [post]
func (h *ProductMessageParserHandler) Execute(c *gin.Context) {
	var req messageParserV1.ExecuteProductMessageParserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	result, err := h.svc.Execute(c, c.Param("productKey"), req)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, result)
}

func toProductMessageParserResponse(productKey string, parser *model.ProductMessageParser, isDefault bool) messageParserV1.ProductMessageParser {
	return messageParserV1.ProductMessageParser{
		ProductKey: productKey,
		Language:   parser.Language,
		Script:     parser.Script,
		IsDefault:  isDefault,
	}
}
