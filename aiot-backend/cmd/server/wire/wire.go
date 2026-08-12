//go:build wireinject
// +build wireinject

package wire

import (
	"0things-backend/internal/handler"
	"0things-backend/internal/job"
	"0things-backend/internal/repository"
	"0things-backend/internal/router"
	"0things-backend/internal/server"
	"0things-backend/internal/service"
	"0things-backend/pkg/app"
	"0things-backend/pkg/jwt"
	"0things-backend/pkg/log"
	"0things-backend/pkg/server/http"
	"0things-backend/pkg/sid"
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
	repository.NewDeviceRepository,
	repository.NewDeviceTagRepository,
	repository.NewDeviceShadowRepository,
	repository.NewPushRecordRepository,
	repository.NewRuleRepository,
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
	service.NewDeviceService,
	service.NewRuleService,
	service.NewAlertService,
	service.NewOTAService,
	service.NewDeviceEventService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewProductHandler,
	handler.NewProductTSLHandler,
	handler.NewDeviceHandler,
	handler.NewRuleHandler,
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
