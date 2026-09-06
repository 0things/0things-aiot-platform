package main

import (
	"context"
	"flag"

	"aiot-backend/cmd/mcp/wire"
	"aiot-backend/pkg/config"
	"aiot-backend/pkg/log"

	"go.uber.org/zap"
)

func main() {
	confPath := flag.String("conf", "config/local.yml", "config path, e.g. -conf ./config/local.yml")
	flag.Parse()

	conf := config.NewConfig(*confPath)
	// STDOUT is reserved for MCP JSON-RPC frames when the server uses stdio.
	conf.Set("log.output", "stderr")
	logger := log.NewLog(conf)
	app, cleanup, err := wire.NewWire(conf, logger)
	if err != nil {
		logger.Fatal("MCP server initialization failed", zap.Error(err))
	}
	defer cleanup()

	if err := app.Run(context.Background()); err != nil {
		logger.Fatal("MCP server stopped with an error", zap.Error(err))
	}
}
