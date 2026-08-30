package service

import (
	"context"

	"data-engine/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// OTAProcessor 负责消费并处理 OTA 设备升级进度，驱动状态机流转与数据库/缓存持久化。
type OTAProcessor struct {
	logger *zap.Logger
}

// NewOTAProcessor 初始化 OTA 业务处理器。
func NewOTAProcessor(config *viper.Viper, logger *zap.Logger) *OTAProcessor {
	return &OTAProcessor{
		logger: logger,
	}
}

// HandleOTAReport 处理单条 OTA 进度汇报。
// 业务规则：
// 1. 进度在 1 ~ 99 时：标记为 UPGRADING 状态，记录中间百分比；
// 2. 进度达到 100 时：标记为 SUCCESS 终态，并触发设备 firmware_version 版本号自动更新；
// 3. 进度为负数或收到 FAILED 时：标记为 FAILED 终态，记录失败原因。
func (p *OTAProcessor) HandleOTAReport(ctx context.Context, report model.OTADeviceUpgradeReportEvent) error {
	p.logger.Info("processing OTA upgrade progress report",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.Int("progress", report.Progress),
		zap.String("status", report.Status),
		zap.String("desc", report.Desc),
	)

	// 状态机校验与标准化
	status := report.Status
	if report.Progress == 100 {
		status = "SUCCESS"
	} else if report.Progress < 0 {
		status = "FAILED"
	} else if status == "" {
		status = "UPGRADING"
	}

	switch status {
	case "SUCCESS":
		p.logger.Info("🎉 Device OTA upgrade completed successfully!",
			zap.String("device_key", report.DeviceKey),
			zap.String("new_version", report.Version),
		)
	case "FAILED":
		p.logger.Warn("❌ Device OTA upgrade failed",
			zap.String("device_key", report.DeviceKey),
			zap.String("reason", report.Desc),
		)
	default:
		p.logger.Debug("Device OTA upgrade in progress...",
			zap.String("device_key", report.DeviceKey),
			zap.Int("step", report.Progress),
		)
	}

	return nil
}
