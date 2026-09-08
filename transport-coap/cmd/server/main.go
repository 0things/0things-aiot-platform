package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"transport-coap/internal/coap"
	"transport-coap/pkg/config"
	"transport-coap/pkg/log"

	"go.uber.org/zap"
)

// main is the entrypoint for 0things CoAP Transport Service.
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things CoAP Transport Service...")

	// 1. 初始化并启动 CoAP UDP 服务
	coapService := coap.NewService(conf, logger.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := coapService.Start(ctx); err != nil {
		logger.Fatal("coap service failed", zap.Error(err))
	}

	logger.Info("CoAP Transport Service stopped gracefully")
}
