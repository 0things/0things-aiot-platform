//go:build wireinject
// +build wireinject

package wire

import (
	"aiot-backend/internal/handler"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/server"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/app"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	mcptransport "aiot-backend/pkg/server/mcp"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewDeviceRepository,
	repository.NewProductRepository,
	repository.NewDeviceTagRepository,
	repository.NewDeviceShadowRepository,
	repository.NewPushRecordRepository,
	repository.NewProductTSLRepository,
	repository.NewDeviceServiceInvocationRepository,
	repository.NewTelemetryRepository,
)

var serviceSet = wire.NewSet(
	service.NewDeviceService,
	service.NewThingModelDataService,
	service.NewTelemetryService,
	wire.Bind(new(handler.MCPDeviceService), new(*service.DeviceService)),
	wire.Bind(new(handler.MCPThingModelService), new(*service.ThingModelDataService)),
	wire.Bind(new(handler.MCPTelemetryService), new(*service.TelemetryService)),
)

func newApp(mcpServer *mcptransport.Server) *app.App {
	return app.NewApp(
		app.WithServer(mcpServer),
		app.WithName("0things-mcp"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		jwt.NewJwt,
		handler.NewMCPHandler,
		server.NewMCPServer,
		server.NewMCPTransportServer,
		newApp,
	))
}
