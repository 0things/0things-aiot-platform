package service

import (
	"context"
	"fmt"

	"0things/pkg/event"

	"go.uber.org/zap"
)

type OTAReportStore interface {
	RecordReport(context.Context, event.OTAUpgradeReport) error
}

// OTAProcessor 负责消费并处理 OTA 设备升级进度，驱动状态机流转与数据库持久化。
type OTAProcessor struct {
	store  OTAReportStore
	logger *zap.Logger
}

func NewOTAProcessor(store OTAReportStore, logger *zap.Logger) *OTAProcessor {
	return &OTAProcessor{store: store, logger: logger}
}

// HandleOTAReport persists an OTA progress or final-version device report.
func (p *OTAProcessor) HandleOTAReport(ctx context.Context, report event.OTAUpgradeReport) error {
	eventType := report.EventType
	if eventType == "" {
		eventType = event.OTAEventTypeProgress
		report.EventType = event.OTAEventTypeProgress
	}
	if eventType != event.OTAEventTypeProgress && eventType != event.OTAEventTypeInform {
		return fmt.Errorf("unsupported OTA event_type %q", eventType)
	}
	p.logger.Info("processing OTA device report",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.String("event_type", string(eventType)),
	)
	return p.store.RecordReport(ctx, report)
}
