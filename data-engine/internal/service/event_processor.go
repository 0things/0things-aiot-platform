package service

import (
	"context"

	"data-engine/internal/model"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// EventProcessor 负责处理设备生命周期事件（上线、离线、心跳）以及自定义业务告警事件。
type EventProcessor struct {
	logger *zap.Logger
}

// NewEventProcessor 初始化事件处理器。
func NewEventProcessor(config *viper.Viper, logger *zap.Logger) *EventProcessor {
	return &EventProcessor{
		logger: logger,
	}
}

// HandleEvent 处理单条设备事件消息。
func (p *EventProcessor) HandleEvent(ctx context.Context, msg model.DeviceMessage) error {
	p.logger.Info("processing device event",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
		zap.Time("timestamp", msg.Timestamp),
	)

	// 业务逻辑处理：
	// 1. 若为上下线事件，更新设备在线状态与最后活跃时间
	// 2. 若为硬件故障/告警事件，记录审计日志或分发联动通知
	return nil
}
