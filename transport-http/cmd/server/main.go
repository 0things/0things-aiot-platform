package main

import (
	"context"
	"flag"

	"transport-http/cmd/server/wire"
	"transport-http/pkg/config"
	"transport-http/pkg/log"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("starting 0things HTTP Transport Service...")
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
