//go:build wireinject
// +build wireinject

package wire

import (
	"transport-http/internal/handler"
	"transport-http/internal/router"
	"transport-http/internal/server"
	"transport-http/pkg/app"
	"transport-http/pkg/log"
	pkgHTTP "transport-http/pkg/server/http"

	"0things/pkg/event"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

func provideEventProducer(conf *viper.Viper, logger *log.Logger) (event.Producer, func(), error) {
	pub, _, err := event.NewBus(conf, logger.Logger)
	if err != nil {
		return nil, nil, err
	}
	producer := event.NewProducer(pub)
	cleanup := func() {
		_ = producer.Close()
	}
	return producer, cleanup, nil
}

var handlerSet = wire.NewSet(
	handler.NewDeviceHandler,
)

var routerSet = wire.NewSet(
	router.NewRouter,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
)

func newApp(
	httpServer *pkgHTTP.Server,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("0things-transport-http"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		provideEventProducer,
		handlerSet,
		routerSet,
		serverSet,
		newApp,
	))
}
