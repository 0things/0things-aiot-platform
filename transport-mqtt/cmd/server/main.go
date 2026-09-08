package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"transport-mqtt/internal/mqtt"
	"transport-mqtt/pkg/config"
	"transport-mqtt/pkg/log"

	"0things/pkg/event"

	"go.uber.org/zap"
)

// main is the entrypoint for 0things MQTT Transport Service.
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things MQTT Transport Service...")

	// 1. 初始化事件总线
	pub, _, err := event.NewBus(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize event bus", zap.Error(err))
	}
	eventProducer := event.NewProducer(pub)
	defer eventProducer.Close()

	// 2. 初始化 MQTT 传输服务
	mqttService, err := mqtt.NewService(conf, logger.Logger, eventProducer)
	if err != nil {
		logger.Fatal("failed to initialize MQTT transport", zap.Error(err))
	}

	// 捕获系统终止信号，确保 K8s 滚动更新或本地停止时不丢失数据
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 主协程阻塞运行 MQTT 客户端服务
	if err := mqttService.Start(ctx); err != nil {
		logger.Fatal("mqtt service stopped with error", zap.Error(err))
	}

	logger.Info("MQTT Transport Service stopped gracefully")
}
