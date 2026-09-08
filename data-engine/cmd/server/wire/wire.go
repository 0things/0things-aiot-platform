//go:build wireinject
// +build wireinject

package wire

import (
	"0things/pkg/event"
	"0things/pkg/tsdb"
	"data-engine/internal/consumer"
	"data-engine/internal/repository"
	"data-engine/internal/server"
	"data-engine/internal/service"
	"data-engine/pkg/app"
	"data-engine/pkg/log"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func provideZapLogger(logger *log.Logger) *zap.Logger {
	return logger.Logger
}

func provideTSDBClient(conf *viper.Viper, logger *log.Logger) (tsdb.Client, func(), error) {
	client := tsdb.NewClient(conf, logger.Logger)
	cleanup := func() {
		client.Close()
	}
	return client, cleanup, nil
}

func provideDB(conf *viper.Viper, logger *log.Logger) (*gorm.DB, error) {
	return repository.NewDB(conf, logger.Logger)
}

type eventBusHolder struct {
	consumer event.Consumer
}

func provideEventBus(conf *viper.Viper, logger *log.Logger) (*eventBusHolder, func(), error) {
	_, sub, err := event.NewBus(conf, logger.Logger)
	if err != nil {
		return nil, nil, err
	}
	c := event.NewConsumer(sub, logger.Logger)
	cleanup := func() {
		_ = c.Close()
	}
	return &eventBusHolder{consumer: c}, cleanup, nil
}

func provideEventConsumer(holder *eventBusHolder) event.Consumer {
	return holder.consumer
}

var repositorySet = wire.NewSet(
	provideDB,
	repository.NewRepository,
	repository.NewOTARepository,
	repository.NewShadowRepository,
)

var serviceSet = wire.NewSet(
	service.NewTelemetryService,
	service.NewOTAService,
	service.NewEventService,
)

var consumerSet = wire.NewSet(
	consumer.NewTelemetryConsumer,
	consumer.NewEventConsumer,
	consumer.NewOTAProgressConsumer,
	consumer.NewManager,
)

var serverSet = wire.NewSet(
	server.NewDataEngineServer,
)

func newApp(
	server *server.DataEngineServer,
) *app.App {
	return app.NewApp(
		app.WithServer(server),
		app.WithName("0things-data-engine"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		provideZapLogger,
		provideTSDBClient,
		provideEventBus,
		provideEventConsumer,
		repositorySet,
		serviceSet,
		consumerSet,
		serverSet,
		newApp,
	))
}
