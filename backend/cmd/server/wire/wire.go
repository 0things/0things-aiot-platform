//go:build wireinject
// +build wireinject

package wire

import (
	"aiot-backend/internal/handler"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/router"
	"aiot-backend/internal/server"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/app"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/logto"
	"aiot-backend/pkg/server/http"

	"aiot-backend/internal/event"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewRedis,
	repository.NewProductRepository,
	repository.NewCategoryRepository,
	repository.NewProductTSLRepository,
	repository.NewProductMessageParserRepository,
	repository.NewDeviceRepository,
	repository.NewDeviceGroupRepository,
	repository.NewDeviceTagRepository,
	repository.NewDeviceShadowRepository,
	repository.NewPushRecordRepository,
	repository.NewOTARepository,
	repository.NewDeviceEventRepository,
	repository.NewDeviceServiceInvocationRepository,
	repository.NewProtocolRepository,
	repository.NewTelemetryRepository,
	repository.NewTransaction,
	repository.NewRuleNodeDefinitionRepository,
	repository.NewRuleChainRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewProductService,
	service.NewCategoryService,
	service.NewProductTSLService,
	service.NewProductMessageParserService,
	service.NewDeviceService,
	service.NewDeviceGroupService,
	service.NewTelemetryService,
	service.NewOTAService,
	service.NewFileService,
	service.NewDeviceEventService,
	service.NewThingModelDataService,
	provideProtocolService,
	provideEventProducer,
	service.NewRuleNodeDefinitionService,
	service.NewRuleChainService,
	service.NewOrganizationProvisioningService,
	wire.Bind(new(service.ProductServiceInterface), new(*service.ProductService)),
	wire.Bind(new(service.CategoryServiceInterface), new(*service.CategoryService)),
	wire.Bind(new(service.ProductTSLServiceInterface), new(*service.ProductTSLService)),
	wire.Bind(new(service.ProductMessageParserServiceInterface), new(*service.ProductMessageParserService)),
	wire.Bind(new(service.DeviceServiceInterface), new(*service.DeviceService)),
	wire.Bind(new(service.DeviceGroupServiceInterface), new(*service.DeviceGroupService)),
	wire.Bind(new(service.TelemetryServiceInterface), new(*service.TelemetryService)),
	wire.Bind(new(service.OTAServiceInterface), new(*service.OTAService)),
	wire.Bind(new(service.FileServiceInterface), new(*service.FileService)),
	wire.Bind(new(service.DeviceEventServiceInterface), new(*service.DeviceEventService)),
	wire.Bind(new(service.ThingModelDataServiceInterface), new(*service.ThingModelDataService)),
	wire.Bind(new(service.ProtocolServiceInterface), new(*service.ProtocolService)),
	wire.Bind(new(service.RuleNodeDefinitionServiceInterface), new(*service.RuleNodeDefinitionService)),
	wire.Bind(new(service.RuleChainServiceInterface), new(*service.RuleChainService)),
	wire.Bind(new(service.OrganizationProvisioningServiceInterface), new(*service.OrganizationProvisioningService)),
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewProductHandler,
	handler.NewCategoryHandler,
	handler.NewProductTSLHandler,
	handler.NewProductMessageParserHandler,
	handler.NewDeviceHandler,
	handler.NewDeviceGroupHandler,
	handler.NewOTAHandler,
	handler.NewFileHandler,
	handler.NewDeviceEventHandler,
	handler.NewThingModelDataHandler,
	handler.NewProtocolHandler,
	handler.NewTelemetryHandler,
	handler.NewRuleNodeDefinitionHandler,
	handler.NewRuleChainHandler,
	handler.NewOrganizationHandler,
)

func provideProtocolService(repo *repository.ProtocolRepository, config *viper.Viper) *service.ProtocolService {
	return service.NewProtocolService(repo, config)
}

func provideLogtoVerifier(config *viper.Viper) (*logto.Verifier, error) {
	return logto.NewVerifier(config)
}

func provideLogtoManagementClient(config *viper.Viper) (*logto.ManagementClient, error) {
	return logto.NewManagementClient(config)
}

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
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("0things-backend"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		wire.Struct(new(router.RouterDeps), "*"),
		provideLogtoVerifier,
		provideLogtoManagementClient,
		newApp,
	))
}
