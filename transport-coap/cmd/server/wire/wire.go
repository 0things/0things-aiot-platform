//go:build wireinject
// +build wireinject

package wire

import (
	"transport-coap/internal/server"
	"transport-coap/pkg/app"
	"transport-coap/pkg/log"

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
		provideEventProducer,
		serverSet,
		newApp,
	))
}
