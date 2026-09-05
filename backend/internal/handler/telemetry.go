package handler

import (
	"errors"
	"net/http"

	apiV1 "aiot-backend/api/v1"
	"aiot-backend/internal/dto"
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

// QueryHistory returns historical telemetry data points for a device property.
// @Summary Query device telemetry history
// @Description Returns telemetry data points within the requested time range.
// @Tags Devices
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param property query string true "Property identifier"
// @Param start_time query int false "Start timestamp"
// @Param end_time query int false "End timestamp"
// @Param limit query int false "Maximum result count"
// @Success 200 {object} apiV1.ApiResponse[[]apiV1.TelemetryPoint] "Successful response"
// @Router /devices/{deviceKey}/telemetry/history [get]
func (h *TelemetryHandler) QueryHistory(c *gin.Context) {
	deviceKey := c.Param("deviceKey")
	if deviceKey == "" {
		apiV1.HandleError(c, http.StatusBadRequest, errors.New("deviceKey is required"), nil)
		return
	}

	var req dto.TelemetryQueryReq
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
