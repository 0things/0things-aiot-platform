package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"0things/pkg/event"
	"0things/pkg/tsdb"
	"data-engine/internal/engine"
	"data-engine/internal/model"
	"data-engine/internal/service"
	"data-engine/internal/storage"
	"data-engine/pkg/config"
	"data-engine/pkg/log"

	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()

	// 1. 初始化配置与日志
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)
	logger.Info("starting 0things data-engine service...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. 初始化持久化存储层 (统一可插拔 TSDB 客户端 & 设备影子)
	tsdbClient := tsdb.NewClient(conf, logger.Logger)
	defer tsdbClient.Close()
	shadowStore := storage.NewShadowStore(conf, logger.Logger)

	// 3. 初始化核心计算、OTA 状态存储与业务处理器
	ruleProcessor := engine.NewProcessor(conf, logger.Logger, tsdbClient, shadowStore)
	otaStore, err := storage.NewOTAStore(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize OTA state store", zap.Error(err))
	}
	otaProcessor := service.NewOTAProcessor(otaStore, logger.Logger)
	eventProcessor := service.NewEventProcessor(conf, logger.Logger)

	// 4. 初始化 OTA MQTT 调度器
	otaDispatcher, err := service.NewOTADispatcher(conf)
	if err != nil {
		logger.Warn("OTA MQTT dispatcher not connected; command forwarding will be skipped", zap.Error(err))
	} else {
		defer otaDispatcher.Close()
	}

	// 5. 初始化事件总线与消费者
	_, sub, err := event.NewBus(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize event bus subscriber", zap.Error(err))
	}
	consumer := event.NewConsumer(sub, logger.Logger)
	defer consumer.Close()

	// 6. 订阅设备遥测与属性上报
	err = event.Subscribe(ctx, consumer, event.TopicDeviceTelemetryReport, func(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
		modelMsg := model.DeviceMessage{
			DeviceKey:   msg.DeviceKey,
			ProductKey:  msg.ProductKey,
			Transport:   msg.Transport,
			MessageType: msg.MessageType,
			Payload:     msg.Payload,
			Timestamp:   msg.Timestamp,
			Headers:     msg.Headers,
		}
		return ruleProcessor.ProcessMessage(ctx, modelMsg)
	})
	if err != nil {
		logger.Fatal("failed to subscribe to device telemetry reports", zap.Error(err))
	}

	err = event.Subscribe(ctx, consumer, event.TopicDeviceAttributeReport, func(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
		modelMsg := model.DeviceMessage{
			DeviceKey:   msg.DeviceKey,
			ProductKey:  msg.ProductKey,
			Transport:   msg.Transport,
			MessageType: msg.MessageType,
			Payload:     msg.Payload,
			Timestamp:   msg.Timestamp,
			Headers:     msg.Headers,
		}
		return ruleProcessor.ProcessMessage(ctx, modelMsg)
	})
	if err != nil {
		logger.Fatal("failed to subscribe to device attribute reports", zap.Error(err))
	}

	// 7. 订阅设备自定义业务告警/事件上报
	err = event.Subscribe(ctx, consumer, event.TopicDeviceEventReport, func(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
		modelMsg := model.DeviceMessage{
			DeviceKey:   msg.DeviceKey,
			ProductKey:  msg.ProductKey,
			Transport:   msg.Transport,
			MessageType: msg.MessageType,
			Payload:     msg.Payload,
			Timestamp:   msg.Timestamp,
			Headers:     msg.Headers,
		}
		return eventProcessor.HandleEvent(ctx, modelMsg)
	})
	if err != nil {
		logger.Fatal("failed to subscribe to device event reports", zap.Error(err))
	}

	// 8. 订阅 OTA 下发升级指令
	err = event.Subscribe(ctx, consumer, event.TopicOTAUpgradeCommand, func(ctx context.Context, cmd *event.OTAUpgradeCommand, meta map[string]string) error {
		logger.Info("received OTA upgrade command event",
			zap.String("device_key", cmd.DeviceKey),
			zap.String("batch_id", cmd.BatchID),
			zap.String("target_version", cmd.TargetVersion),
		)
		if otaDispatcher != nil {
			modelCmd := model.OTAUpgradeCommand{
				BatchID:       cmd.BatchID,
				PackageID:     cmd.PackageID,
				ProductKey:    cmd.ProductKey,
				DeviceKey:     cmd.DeviceKey,
				DeviceName:    cmd.DeviceName,
				Module:        cmd.Module,
				TargetVersion: cmd.TargetVersion,
				DownloadURL:   cmd.DownloadURL,
				FileSize:      cmd.FileSize,
				SHA256:        cmd.SHA256,
				ExpiresAt:     cmd.ExpiresAt,
			}
			if err := otaDispatcher.Dispatch(ctx, modelCmd); err != nil {
				logger.Error("failed to dispatch OTA command via MQTT", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Fatal("failed to subscribe to OTA upgrade commands", zap.Error(err))
	}

	// 9. 订阅设备 OTA 进度上报
	err = event.Subscribe(ctx, consumer, event.TopicOTAProgressReport, func(ctx context.Context, report *event.OTAUpgradeReport, meta map[string]string) error {
		logger.Info("received OTA progress report event",
			zap.String("device_key", report.DeviceKey),
			zap.String("batch_id", report.BatchID),
			zap.String("status", report.Status),
		)
		modelReport := model.OTAUpgradeReport{
			EventType:       report.EventType,
			BatchID:         report.BatchID,
			ProductKey:      report.ProductKey,
			DeviceKey:       report.DeviceKey,
			Module:          report.Module,
			Progress:        report.Progress,
			Stage:           report.Stage,
			ReportedVersion: report.ReportedVersion,
			ReportedAt:      report.ReportedAt,
		}
		if modelReport.EventType == "" {
			modelReport.EventType = "progress"
		}
		return otaProcessor.HandleOTAReport(ctx, modelReport)
	})
	if err != nil {
		logger.Fatal("failed to subscribe to OTA progress reports", zap.Error(err))
	}

	logger.Info("0things data-engine initialized successfully with telemetry and OTA event consumers")

	// 10. 监听系统优雅关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down data-engine gracefully...")
	cancel()
	logger.Info("data-engine stopped successfully")
}
