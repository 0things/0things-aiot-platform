package handler

import (
	protocolV1 "aiot-backend/api/protocol/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// 保留协议响应类型引用，供 Swagger 泛型注释解析。
var _ = protocolV1.DeviceEndpoints{}

type ProtocolHandler struct {
	*Handler
	svc service.ProtocolServiceInterface
}

func NewProtocolHandler(h *Handler, svc service.ProtocolServiceInterface) *ProtocolHandler {
	return &ProtocolHandler{Handler: h, svc: svc}
}

// ListDeviceEndpoints 返回设备按产品协议生成的分协议接入参数。
// @Summary 获取设备连接数据
// @Tags devices
// @Produce json
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[protocolV1.DeviceEndpoints]
// @Router /devices/{deviceKey}/endpoints [get]
func (h *ProtocolHandler) ListDeviceEndpoints(c *gin.Context) {
	items, err := h.svc.ListDeviceEndpoints(c, c.Param("deviceKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, items)
}
