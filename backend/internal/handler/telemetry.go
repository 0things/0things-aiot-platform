package handler

import (
	"errors"
	"net/http"

	apiV1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/log"
	"github.com/gin-gonic/gin"
)

type TelemetryHandler struct {
	telemetryService service.TelemetryServiceInterface
	logger           *log.Logger
}

func NewTelemetryHandler(telemetryService service.TelemetryServiceInterface, logger *log.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryService: telemetryService,
		logger:           logger,
	}
}

// QueryHistory 查询指定设备的时序历史折线数据点
// @Summary      查询设备时序历史数据
// @Tags         时序数据 (Telemetry)
// @Produce      json
// @Param        deviceKey  path      string  true  "设备Key"
// @Param        property   query     string  true  "物模型属性标识 (如 temperature)"
// @Param        start_time query     int     false "起始时间戳(ms)"
// @Param        end_time   query     int     false "结束时间戳(ms)"
// @Param        limit      query     int     false "返回点数上限"
// @Success      200        {object}  apiV1.Response
// @Router       /v1/devices/{deviceKey}/telemetry/history [get]
func (h *TelemetryHandler) QueryHistory(c *gin.Context) {
	deviceKey := c.Param("deviceKey")
	if deviceKey == "" {
		apiV1.HandleError(c, http.StatusBadRequest, errors.New("deviceKey is required"), nil)
		return
	}

	var req model.TelemetryQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiV1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	req.DeviceKey = deviceKey

	points, err := h.telemetryService.QueryHistory(c.Request.Context(), req)
	if err != nil {
		apiV1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	apiV1.HandleSuccess(c, points)
}
