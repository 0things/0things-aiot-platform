//go:build wireinject
// +build wireinject

package wire

import (
	"transport-mqtt/internal/server"
	"transport-mqtt/pkg/app"
	"transport-mqtt/pkg/log"

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
	server.NewMQTTServer,
)

func newApp(
	mqttServer *server.MQTTServer,
) *app.App {
	return app.NewApp(
		app.WithServer(mqttServer),
		app.WithName("0things-transport-mqtt"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		provideEventProducer,
		serverSet,
		newApp,
	))
}
