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
	repository.NewProductTSLRepository,
	repository.NewProductMessageParserRepository,
	repository.NewDeviceRepository,
	repository.NewDeviceTagRepository,
	repository.NewDeviceShadowRepository,
	repository.NewPushRecordRepository,
	repository.NewSceneLinkageRepository,
	repository.NewSceneLinkageDetailRepository,
	repository.NewAlertRepository,
	repository.NewOTARepository,
	repository.NewDeviceEventRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewProductService,
	service.NewProductTSLService,
	service.NewProductMessageParserService,
	service.NewDeviceService,
	service.NewSceneLinkageService,
	service.NewSceneLinkageDetailService,
	service.NewAlertService,
	service.NewOTAService,
	service.NewDeviceEventService,
	wire.Bind(new(service.ProductServiceInterface), new(*service.ProductService)),
	wire.Bind(new(service.ProductTSLServiceInterface), new(*service.ProductTSLService)),
	wire.Bind(new(service.ProductMessageParserServiceInterface), new(*service.ProductMessageParserService)),
	wire.Bind(new(service.DeviceServiceInterface), new(*service.DeviceService)),
	wire.Bind(new(service.SceneLinkageServiceInterface), new(*service.SceneLinkageService)),
	wire.Bind(new(service.SceneLinkageDetailServiceInterface), new(*service.SceneLinkageDetailService)),
	wire.Bind(new(service.AlertServiceInterface), new(*service.AlertService)),
	wire.Bind(new(service.OTAServiceInterface), new(*service.OTAService)),
	wire.Bind(new(service.DeviceEventServiceInterface), new(*service.DeviceEventService)),
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewProductHandler,
	handler.NewProductTSLHandler,
	handler.NewProductMessageParserHandler,
	handler.NewDeviceHandler,
	handler.NewSceneLinkageHandler,
	handler.NewSceneLinkageDetailHandler,
	handler.NewAlertHandler,
	handler.NewOTAHandler,
	handler.NewDeviceEventHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
	job.NewDeviceEventConsumer,
)
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
