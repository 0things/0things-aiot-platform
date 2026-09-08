package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"http-transport/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeviceHandler provides RESTful ingress endpoints for HTTP devices and gateways.
type DeviceHandler struct {
	logger *zap.Logger
}

func NewDeviceHandler(logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		logger: logger,
	}
}

// PostTelemetry handles telemetry and property updates.
func (h *DeviceHandler) PostTelemetry(c *gin.Context) {
	h.handleIngress(c, "telemetry", nil)
}

// PostAttributes handles attributes report.
func (h *DeviceHandler) PostAttributes(c *gin.Context) {
	h.handleIngress(c, "attributes", nil)
}

// PostEvent handles custom device event/alarm reports.
func (h *DeviceHandler) PostEvent(c *gin.Context) {
	eventType := c.Param("eventType")
	if eventType == "" {
		eventType = "event"
	}
	h.handleIngress(c, "event", map[string]string{"event_type": eventType})
}

// PostOtaProgress handles device OTA progress reports.
func (h *DeviceHandler) PostOtaProgress(c *gin.Context) {
	h.handleIngress(c, "ota_report", nil)
}

// DeviceIngressLegacy handles legacy ingress path.
func (h *DeviceHandler) DeviceIngressLegacy(c *gin.Context) {
	msgType := c.GetHeader("X-Device-Message-Type")
	if msgType == "" {
		msgType = "telemetry"
	}
	h.handleIngress(c, msgType, nil)
}

func (h *DeviceHandler) handleIngress(c *gin.Context, msgType string, extraHeaders map[string]string) {
	deviceKey := strings.TrimSpace(c.Param("deviceKey"))
	if deviceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "deviceKey is required"})
		return
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 4<<20))
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request body"})
		return
	}

	headers := map[string]string{
		"content-type": c.ContentType(),
		"remote-ip":    c.ClientIP(),
	}
	for k, v := range extraHeaders {
		headers[k] = v
	}

	msg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  c.GetHeader("X-Product-Key"),
		Transport:   "http",
		MessageType: msgType,
		Payload:     json.RawMessage(body),
		Timestamp:   time.Now().UTC(),
		Headers:     headers,
	}

	h.logger.Info("received HTTP ingress message",
		zap.String("device_key", msg.DeviceKey),
		zap.String("product_key", msg.ProductKey),
		zap.String("msg_type", msg.MessageType),
		zap.Int("payload_bytes", len(body)),
	)

	c.JSON(http.StatusAccepted, gin.H{
		"code":      200,
		"message":   "accepted",
		"timestamp": time.Now().UnixMilli(),
	})
}
