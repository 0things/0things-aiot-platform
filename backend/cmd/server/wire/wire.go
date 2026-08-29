//go:build wireinject
// +build wireinject

package wire

import (
	"aiot-backend/internal/handler"
	"aiot-backend/internal/job"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/router"
	"aiot-backend/internal/server"
	"aiot-backend/internal/service"
	"aiot-backend/internal/transport"
	coaptransport "aiot-backend/internal/transport/coap"
	"aiot-backend/pkg/app"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/server/http"
	"aiot-backend/pkg/sid"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	//repository.NewMongo,
	repository.NewRepository,
	repository.NewIoTDB,
	repository.NewIoTRedis,
	repository.NewProductRepository,
	repository.NewCategoryRepository,
	repository.NewProductTSLRepository,
	repository.NewProductMessageParserRepository,
	repository.NewDeviceRepository,
	repository.NewDeviceGroupRepository,
	repository.NewDeviceTagRepository,
	repository.NewDeviceShadowRepository,
	repository.NewPushRecordRepository,
	repository.NewSceneLinkageRepository,
	repository.NewSceneLinkageDetailRepository,
	repository.NewOTARepository,
	repository.NewDeviceEventRepository,
	repository.NewProtocolRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewOrganizationRepository,
	repository.NewOrganizationUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewProductService,
	service.NewCategoryService,
	service.NewProductTSLService,
	service.NewProductMessageParserService,
	service.NewDeviceService,
	service.NewDeviceGroupService,
	service.NewKafkaService,
	service.NewMQTTService,
	service.NewSceneLinkageService,
	service.NewSceneLinkageDetailService,
	provideOTAService,
	service.NewFileService,
	service.NewDeviceEventService,
	provideProtocolService,
	wire.Bind(new(service.ProductServiceInterface), new(*service.ProductService)),
	wire.Bind(new(service.CategoryServiceInterface), new(*service.CategoryService)),
	wire.Bind(new(service.ProductTSLServiceInterface), new(*service.ProductTSLService)),
	wire.Bind(new(service.ProductMessageParserServiceInterface), new(*service.ProductMessageParserService)),
	wire.Bind(new(service.DeviceServiceInterface), new(*service.DeviceService)),
	wire.Bind(new(service.DeviceGroupServiceInterface), new(*service.DeviceGroupService)),
	wire.Bind(new(service.KafkaServiceInterface), new(*service.KafkaService)),
	wire.Bind(new(service.MQTTServiceInterface), new(*service.MQTTService)),
	wire.Bind(new(service.SceneLinkageServiceInterface), new(*service.SceneLinkageService)),
	wire.Bind(new(service.SceneLinkageDetailServiceInterface), new(*service.SceneLinkageDetailService)),
	wire.Bind(new(service.OTAServiceInterface), new(*service.OTAService)),
	wire.Bind(new(service.FileServiceInterface), new(*service.FileService)),
	wire.Bind(new(service.DeviceEventServiceInterface), new(*service.DeviceEventService)),
	wire.Bind(new(service.ProtocolServiceInterface), new(*service.ProtocolService)),
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewProductHandler,
	handler.NewCategoryHandler,
	handler.NewProductTSLHandler,
	handler.NewProductMessageParserHandler,
	handler.NewDeviceHandler,
	handler.NewDeviceGroupHandler,
	handler.NewSceneLinkageHandler,
	handler.NewSceneLinkageDetailHandler,
	handler.NewOTAHandler,
	handler.NewFileHandler,
	handler.NewDeviceEventHandler,
	handler.NewProtocolHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
	job.NewDeviceEventConsumer,
	job.NewDeviceMessageConsumer,
	job.NewMQTTTransportAdapterForWire,
	provideTransportRegistry,
	provideOTACommandConsumer,
	job.NewOTAMQTTReportBridge,
	job.NewOTAReportConsumer,
)

func provideOTACommandConsumer(config *viper.Viper, mqtt service.MQTTServiceInterface, ota *service.OTAService, logger *log.Logger, registry *transport.Registry) (*job.OTACommandConsumer, error) {
	return job.NewOTACommandConsumer(config, mqtt, ota, logger, registry)
}

func provideOTAService(repo *repository.OTARepository, productRepo *repository.ProductRepository, deviceRepo *repository.DeviceRepository, kafka service.KafkaServiceInterface, protocols *repository.ProtocolRepository) *service.OTAService {
	return service.NewOTAServiceWithProtocol(repo, productRepo, deviceRepo, kafka, protocols)
}

func provideProtocolService(repo *repository.ProtocolRepository, config *viper.Viper) *service.ProtocolService {
	return service.NewProtocolService(repo, config)
}

func provideTransportRegistry(adapter transport.Adapter) (*transport.Registry, error) {
	// 管理服务只负责消费 OTA 命令，协议连接由对应适配器执行；CoAP 适配器无需在此监听端口。
	return transport.NewRegistry(adapter, coaptransport.NewAdapter(""))
}

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		wire.Struct(new(router.RouterDeps), "*"),
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
