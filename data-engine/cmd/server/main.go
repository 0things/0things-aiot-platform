package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"0things/pkg/event"
	"0things/pkg/tsdb"
	"data-engine/internal/consumer"
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

	// 1. Initialize configuration and logger
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)
	logger.Info("starting 0things data-engine service...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize storage layer (pluggable TSDB client & device shadow)
	tsdbClient := tsdb.NewClient(conf, logger.Logger)
	defer tsdbClient.Close()
	shadowStore := storage.NewShadowStore(conf, logger.Logger)

	// 3. Initialize rule processor, OTA state store and event processor
	ruleProcessor := engine.NewProcessor(conf, logger.Logger, tsdbClient, shadowStore)
	otaStore, err := storage.NewOTAStore(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize OTA state store", zap.Error(err))
	}
	otaProcessor := service.NewOTAProcessor(otaStore, logger.Logger)
	eventProcessor := service.NewEventProcessor(conf, logger.Logger)

	// 4. Initialize event bus and consumer manager
	_, sub, err := event.NewBus(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize event bus subscriber", zap.Error(err))
	}
	eventConsumer := event.NewConsumer(sub, logger.Logger)
	defer eventConsumer.Close()

	telemetryConsumer := consumer.NewTelemetryConsumer(ruleProcessor, logger.Logger)
	deviceEventConsumer := consumer.NewEventConsumer(eventProcessor, logger.Logger)
	otaProgressConsumer := consumer.NewOTAProgressConsumer(otaProcessor, logger.Logger)

	consumerManager := consumer.NewManager(
		eventConsumer,
		telemetryConsumer,
		deviceEventConsumer,
		otaProgressConsumer,
		logger.Logger,
	)

	// 5. Register and start domain event subscriptions
	if err := consumerManager.Start(ctx); err != nil {
		logger.Fatal("failed to start event consumer manager", zap.Error(err))
	}

	logger.Info("0things data-engine initialized successfully with domain event consumers")

	// 6. Wait for termination signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down data-engine gracefully...")
	cancel()
	logger.Info("data-engine stopped successfully")
}
