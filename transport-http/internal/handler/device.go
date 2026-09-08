package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"transport-http/pkg/log"

	"0things/pkg/event"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeviceHandler provides RESTful ingress endpoints for HTTP devices and gateways.
type DeviceHandler struct {
	logger        *log.Logger
	eventProducer event.Producer
}

// NewDeviceHandler initializes DeviceHandler.
func NewDeviceHandler(logger *log.Logger, eventProducer event.Producer) *DeviceHandler {
	return &DeviceHandler{
		logger:        logger,
		eventProducer: eventProducer,
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

	h.logger.Info("received HTTP ingress message",
		zap.String("device_key", deviceKey),
		zap.String("product_key", c.GetHeader("X-Product-Key")),
		zap.String("msg_type", msgType),
		zap.Int("payload_bytes", len(body)),
	)

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// 投递至事件总线
	if h.eventProducer != nil {
		if msgType == "ota_report" {
			var report event.OTAUpgradeReport
			if err := json.Unmarshal(body, &report); err == nil {
				if report.DeviceKey == "" {
					report.DeviceKey = deviceKey
				}
				if report.ProductKey == "" {
					report.ProductKey = c.GetHeader("X-Product-Key")
				}
				if report.ReportedAt.IsZero() {
					report.ReportedAt = time.Now().UTC()
				}
				if err := h.eventProducer.Publish(ctx, event.TopicOTAProgressReport, &report); err != nil {
					h.logger.Error("failed to publish OTA progress report event", zap.Error(err))
				}
			}
		} else {
			msg := event.DeviceMessage{
				DeviceKey:   deviceKey,
				ProductKey:  c.GetHeader("X-Product-Key"),
				Transport:   "http",
				MessageType: msgType,
				Payload:     json.RawMessage(body),
				Timestamp:   time.Now().UTC(),
				Headers:     headers,
			}

			var topic event.Topic
			switch msgType {
			case "attributes":
				topic = event.TopicDeviceAttributeReport
			case "event":
				topic = event.TopicDeviceEventReport
			default:
				topic = event.TopicDeviceTelemetryReport
			}

			if err := h.eventProducer.Publish(ctx, topic, &msg); err != nil {
				h.logger.Error("failed to publish device message event", zap.String("topic", topic.String()), zap.Error(err))
			}
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"code":      200,
		"message":   "accepted",
		"timestamp": time.Now().UnixMilli(),
	})
}
