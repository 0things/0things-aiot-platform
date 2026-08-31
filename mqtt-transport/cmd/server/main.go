package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"mqtt-transport/internal/kafka"
	"mqtt-transport/internal/mqtt"
	"mqtt-transport/pkg/config"
	"mqtt-transport/pkg/log"
	"go.uber.org/zap"
)

// main 是 0things MQTT 传输微服务的启动入口。
// 职责：
// 1. 初始化配置与全局日志；
// 2. 建立 Kafka 上行生产通道与下行消费通道；
// 3. 启动 Paho MQTT 客户端订阅上行报文并桥接到 Kafka；
// 4. 监听操作系统信号（SIGINT/SIGTERM），实现优雅停机。
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things MQTT Transport Service...")

	// 1. 初始化 Kafka Producer（用于将设备 MQTT 上行报文投递到 device.message.v1）
	producer, cleanupProducer, err := kafka.NewProducer(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize kafka producer", zap.Error(err))
	}
	defer cleanupProducer()

	// 2. 初始化 MQTT 传输服务（维持与 Broker 的连接与 Topic 订阅）
	mqttService := mqtt.NewService(conf, logger.Logger, producer)

	// 3. 初始化 Kafka Consumer（用于消费 device.command.v1 下行指令并通过 MQTT 推送到物理设备）
	downlinkConsumer, err := kafka.NewConsumer(conf, logger.Logger, mqttService.HandleDownlinkCommand)
	if err != nil {
		logger.Fatal("failed to initialize kafka downlink consumer", zap.Error(err))
	}
	defer downlinkConsumer.Close()

	// 捕获系统终止信号，确保 K8s 滚动更新或本地停止时不丢失数据
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 后台异步启动下行命令消费循环
	go func() {
		if err := downlinkConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			logger.Error("downlink consumer stopped with error", zap.Error(err))
		}
	}()

	// 主协程阻塞运行 MQTT 客户端服务
	if err := mqttService.Start(ctx); err != nil {
		logger.Fatal("mqtt service stopped with error", zap.Error(err))
	}

	logger.Info("MQTT Transport Service stopped gracefully")
}
