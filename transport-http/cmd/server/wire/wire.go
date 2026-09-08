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

	"github.com/google/wire"
	"github.com/spf13/viper"
)

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
		handlerSet,
		routerSet,
		serverSet,
		newApp,
	))
}
