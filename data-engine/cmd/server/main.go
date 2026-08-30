package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"data-engine/internal/engine"
	"data-engine/internal/kafka"
	"data-engine/internal/service"
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

	// 2. 初始化核心计算与业务处理器
	ruleProcessor := engine.NewProcessor(conf, logger.Logger)
	otaProcessor := service.NewOTAProcessor(conf, logger.Logger)
	eventProcessor := service.NewEventProcessor(conf, logger.Logger)

	// 3. 初始化遥测流消费者 (消费 device.telemetry.v1)
	telemetryConsumer, err := kafka.NewTelemetryConsumer(conf, logger.Logger, ruleProcessor)
	if err != nil {
		logger.Fatal("failed to initialize telemetry consumer", zap.Error(err))
	}
	defer telemetryConsumer.Close()

	// 4. 初始化 OTA 任务消费者 (消费 ota.report.v1)
	otaConsumer, err := kafka.NewOTAConsumer(conf, logger.Logger, otaProcessor)
	if err != nil {
		logger.Fatal("failed to initialize ota consumer", zap.Error(err))
	}
	defer otaConsumer.Close()

	// 5. 初始化设备生命周期与业务事件消费者 (消费 device.event.v1)
	eventConsumer, err := kafka.NewEventConsumer(conf, logger.Logger, eventProcessor)
	if err != nil {
		logger.Fatal("failed to initialize event consumer", zap.Error(err))
	}
	defer eventConsumer.Close()

	// 6. 并发启动三大核心计算与任务协程
	go func() {
		if err := telemetryConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("telemetry consumer stopped with error", zap.Error(err))
		}
	}()

	go func() {
		if err := otaConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("ota consumer stopped with error", zap.Error(err))
		}
	}()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("event consumer stopped with error", zap.Error(err))
		}
	}()

	logger.Info("0things data-engine running with telemetry, ota & event consumers started successfully")

	// 7. 监听系统优雅关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down data-engine gracefully...")
	cancel()
	logger.Info("data-engine stopped successfully")
}
