package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"0things/pkg/tsdb"
	"data-engine/internal/engine"
	"data-engine/internal/kafka"
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

	// 4. 初始化遥测流消费者 (消费 device.telemetry.v1)
	telemetryConsumer, err := kafka.NewTelemetryConsumer(conf, logger.Logger, ruleProcessor)
	if err != nil {
		logger.Fatal("failed to initialize telemetry consumer", zap.Error(err))
	}
	defer telemetryConsumer.Close()

	// 5. data-engine 独立消费 OTA 上报与指令；下发结果由命令 consumer 直接写库。
	otaProgressReportConsumer, err := kafka.NewOTAProgressReportConsumer(conf, logger.Logger, otaProcessor)
	if err != nil {
		logger.Fatal("failed to initialize ota consumer", zap.Error(err))
	}
	defer otaProgressReportConsumer.Close()
	otaDispatcher, err := service.NewOTADispatcher(conf)
	if err != nil {
		logger.Fatal("failed to initialize OTA MQTT dispatcher", zap.Error(err))
	}
	defer otaDispatcher.Close()
	otaCommandConsumer, err := kafka.NewOTACommandConsumer(conf, logger.Logger, otaDispatcher, otaStore)
	if err != nil {
		logger.Fatal("failed to initialize OTA command consumer", zap.Error(err))
	}
	defer otaCommandConsumer.Close()

	// 6. 初始化设备生命周期与业务事件消费者 (消费 device.event.v1)
	eventConsumer, err := kafka.NewEventConsumer(conf, logger.Logger, eventProcessor)
	if err != nil {
		logger.Fatal("failed to initialize event consumer", zap.Error(err))
	}
	defer eventConsumer.Close()

	// 7. 并发启动流计算与各专职 OTA 协程。
	go func() {
		if err := telemetryConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("telemetry consumer stopped with error", zap.Error(err))
		}
	}()

	go func() {
		if err := otaProgressReportConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("ota report consumer stopped with error", zap.Error(err))
		}
	}()
	go func() {
		if err := otaCommandConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("ota command consumer stopped with error", zap.Error(err))
		}
	}()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("event consumer stopped with error", zap.Error(err))
		}
	}()

	logger.Info("0things data-engine running with telemetry, OTA and event consumers started successfully")

	// 8. 监听系统优雅关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down data-engine gracefully...")
	cancel()
	logger.Info("data-engine stopped successfully")
}
