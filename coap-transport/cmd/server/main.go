package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"coap-transport/internal/coap"
	"coap-transport/internal/kafka"
	"coap-transport/pkg/config"
	"coap-transport/pkg/log"
	"go.uber.org/zap"
)

// main 是 0things CoAP 传输微服务的启动入口。
// 职责：
// 1. 监听 UDP 5683 端口接收 NB-IoT / 受限设备报文；
// 2. 将数据封装投递至 Kafka device.message.v1 主题；
// 3. 及时向设备回复确认，支持低功耗休眠唤醒。
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things CoAP Transport Service...")

	// 1. 初始化 Kafka Producer
	producer, cleanupProducer, err := kafka.NewProducer(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize kafka producer", zap.Error(err))
	}
	defer cleanupProducer()

	// 2. 初始化并启动 CoAP UDP 服务
	coapService := coap.NewService(conf, logger.Logger, producer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := coapService.Start(ctx); err != nil {
		logger.Fatal("coap service failed", zap.Error(err))
	}

	logger.Info("CoAP Transport Service stopped gracefully")
}
