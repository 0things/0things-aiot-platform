package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"mq-service/internal/kafka"
	"mq-service/internal/service"
	"mq-service/pkg/config"
	"mq-service/pkg/log"
	"go.uber.org/zap"
)

// main 是 0things 异步任务消费微服务（mq-service）的启动入口。
// 职责：
// 1. 专职从 Kafka 消费海量异步任务（如 OTA 设备升级进度回报、生命周期事件等）；
// 2. 驱动后台状态机流转与数据库持久化；
// 3. 支持横向水平扩容多个 Worker Pod 并行消费洪峰，彻底解放 Web 后台。
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things MQ Task Consumer Service (mq-service)...")

	// 1. 初始化 OTA 业务处理器
	otaProcessor := service.NewOTAProcessor(conf, logger.Logger)

	// 2. 初始化 OTA 进度 Kafka 消费者
	otaConsumer, err := kafka.NewOTAConsumer(conf, logger.Logger, otaProcessor)
	if err != nil {
		logger.Fatal("failed to initialize OTA Kafka consumer", zap.Error(err))
	}
	defer otaConsumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Info("mq-service is running and listening for Kafka tasks...")

	// 启动 OTA 进度消费循环
	if err := otaConsumer.Start(ctx); err != nil {
		logger.Fatal("OTA consumer stopped with error", zap.Error(err))
	}

	logger.Info("mq-service stopped gracefully")
}
