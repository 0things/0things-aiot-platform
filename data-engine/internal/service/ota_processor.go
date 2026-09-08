package service

import (
	"context"
	"fmt"

	"data-engine/internal/model"
	"go.uber.org/zap"
)

type OTAReportStore interface {
	RecordReport(context.Context, model.OTAUpgradeReport) error
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
func (p *OTAProcessor) HandleOTAReport(ctx context.Context, report model.OTAUpgradeReport) error {
	if report.EventType != "progress" && report.EventType != "inform" {
		return fmt.Errorf("unsupported OTA event_type %q", report.EventType)
	}
	p.logger.Info("processing OTA device report",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.String("event_type", report.EventType),
	)
	return p.store.RecordReport(ctx, report)
}
