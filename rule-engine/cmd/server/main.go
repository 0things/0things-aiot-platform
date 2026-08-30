package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"rule-engine/internal/engine"
	"rule-engine/internal/kafka"
	"rule-engine/pkg/config"
	"rule-engine/pkg/log"
	"go.uber.org/zap"
)

// main 是 0things 规则引擎计算微服务（Rule Engine）的启动入口。
// 职责：
// 1. 作为纯后台计算 Worker 启动，无对外管理 HTTP 端口；
// 2. 从 Kafka device.message.v1 高并发拉取上行设备报文；
// 3. 执行物模型 TSL 解析、时序指标抽取与业务告警规则评估；
// 4. 监听系统终止信号实现平滑停机。
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things Rule Engine Service...")

	// 1. 初始化规则流式计算处理器（TSL 解析与阈值评估）
	processor := engine.NewProcessor(conf, logger.Logger)

	// 2. 初始化 Kafka 消息流消费者
	consumer, err := kafka.NewConsumer(conf, logger.Logger, processor)
	if err != nil {
		logger.Fatal("failed to initialize kafka consumer", zap.Error(err))
	}
	defer consumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Info("Rule Engine is running and consuming device messages...")

	// 阻塞运行消费循环
	if err := consumer.Start(ctx); err != nil {
		logger.Fatal("rule engine consumer failed", zap.Error(err))
	}

	logger.Info("Rule Engine Service stopped gracefully")
}
