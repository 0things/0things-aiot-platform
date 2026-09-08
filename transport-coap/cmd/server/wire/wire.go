//go:build wireinject
// +build wireinject

package wire

import (
	"transport-coap/internal/server"
	"transport-coap/pkg/app"
	"transport-coap/pkg/log"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var serverSet = wire.NewSet(
	server.NewCoAPServer,
)

func newApp(
	coapServer *server.CoAPServer,
) *app.App {
	return app.NewApp(
		app.WithServer(coapServer),
		app.WithName("0things-transport-coap"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		serverSet,
		newApp,
	))
}
