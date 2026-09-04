package handler

import (
	protocolV1 "aiot-backend/api/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// Keep this response type reference so swag can resolve the generic response annotation.
var _ = protocolV1.DeviceEndpoints{}

type ProtocolHandler struct {
	*Handler
	svc service.ProtocolServiceInterface
}

func NewProtocolHandler(h *Handler, svc service.ProtocolServiceInterface) *ProtocolHandler {
	return &ProtocolHandler{Handler: h, svc: svc}
}

// ListDeviceEndpoints returns connection endpoints generated from the device product protocol.
// @Summary Get device connection endpoints
// @Description Returns protocol-specific connection endpoints for a device.
// @Tags Device connections
// @Produce json
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[protocolV1.DeviceEndpoints] "Successful response"
// @Router /devices/{deviceKey}/endpoints [get]
func (h *ProtocolHandler) ListDeviceEndpoints(c *gin.Context) {
	items, err := h.svc.ListDeviceEndpoints(c, c.Param("deviceKey"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, items)
}
