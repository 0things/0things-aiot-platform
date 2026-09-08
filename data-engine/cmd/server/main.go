package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"0things/pkg/tsdb"
	"data-engine/internal/engine"
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

	// 2. 初始化持久化存储层 (统一可插拔 TSDB 客户端 & 设备影子)
	tsdbClient := tsdb.NewClient(conf, logger.Logger)
	defer tsdbClient.Close()
	shadowStore := storage.NewShadowStore(conf, logger.Logger)

	// 3. 初始化核心计算、OTA 状态存储与业务处理器
	ruleProcessor := engine.NewProcessor(conf, logger.Logger, tsdbClient, shadowStore)
	_ = ruleProcessor

	otaStore, err := storage.NewOTAStore(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize OTA state store", zap.Error(err))
	}
	otaProcessor := service.NewOTAProcessor(otaStore, logger.Logger)
	_ = otaProcessor

	eventProcessor := service.NewEventProcessor(conf, logger.Logger)
	_ = eventProcessor

	logger.Info("0things data-engine initialized successfully")

	// 4. 监听系统优雅关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down data-engine gracefully...")
	logger.Info("data-engine stopped successfully")
}
