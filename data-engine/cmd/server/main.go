package main

import (
	"context"
	"flag"

	"data-engine/cmd/server/wire"
	"data-engine/pkg/config"
	"data-engine/pkg/log"
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

	logger.Info("starting 0things data-engine service...")
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
