package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"http-transport/internal/kafka"
	"http-transport/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeviceHandler 提供针对 HTTP 设备/网关的标准 RESTful 数据上报接口。
type DeviceHandler struct {
	producer *kafka.Producer
	logger   *zap.Logger
}

func NewDeviceHandler(producer *kafka.Producer, logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		producer: producer,
		logger:   logger,
	}
}

// PostTelemetry 处理设备时序遥测数据上报
// 接口：POST /api/v1/:deviceKey/telemetry ➔ 投递到 device.telemetry.v1
func (h *DeviceHandler) PostTelemetry(c *gin.Context) {
	h.handleIngress(c, "telemetry")
}

// PostAttributes 处理设备属性更新上报
// 接口：POST /api/v1/:deviceKey/attributes ➔ 投递到 device.telemetry.v1
func (h *DeviceHandler) PostAttributes(c *gin.Context) {
	h.handleIngress(c, "attributes")
}

// PostEvent 处理设备特定事件上报（如告警、故障）
// 接口：POST /api/v1/:deviceKey/events/:eventType ➔ 投递到 device.event.v1
func (h *DeviceHandler) PostEvent(c *gin.Context) {
	eventType := c.Param("eventType")
	if eventType == "" {
		eventType = "event"
	}
	h.handleIngress(c, "event")
}

// PostOtaProgress 处理设备 OTA 升级进度上报
// 接口：POST /api/v1/:deviceKey/ota/progress ➔ 投递到 ota.report.v1
func (h *DeviceHandler) PostOtaProgress(c *gin.Context) {
	h.handleIngress(c, "ota_report")
}

// DeviceIngressLegacy 兼容老网关上报路径
// 接口：POST /v1/device-ingress/:deviceKey
func (h *DeviceHandler) DeviceIngressLegacy(c *gin.Context) {
	msgType := c.GetHeader("X-Device-Message-Type")
	if msgType == "" {
		msgType = "telemetry"
	}
	h.handleIngress(c, msgType)
}

// handleIngress 统一校验请求参数，组装 DeviceMessage 并异步投递 Kafka，快速返回 202 Accepted 避免设备端等待。
func (h *DeviceHandler) handleIngress(c *gin.Context, msgType string) {
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

	msg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  c.GetHeader("X-Product-Key"),
		Transport:   "http",
		MessageType: msgType,
		Payload:     json.RawMessage(body),
		Timestamp:   time.Now().UTC(),
		Headers: map[string]string{
			"content-type": c.ContentType(),
			"remote-ip":    c.ClientIP(),
		},
	}

	if err := h.producer.SendDeviceMessage(c.Request.Context(), msg); err != nil {
		h.logger.Error("failed to produce http message to kafka", zap.String("device_key", deviceKey), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "message": "message queue dispatch failed"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"code":      200,
		"message":   "accepted",
		"timestamp": time.Now().UnixMilli(),
	})
}
