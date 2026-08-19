package router

import (
	"aiot-backend/internal/handler"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger                      *log.Logger
	Config                      *viper.Viper
	JWT                         *jwt.JWT
	UserHandler                 *handler.UserHandler
	ProductHandler              *handler.ProductHandler
	ProductTSLHandler           *handler.ProductTSLHandler
	ProductMessageParserHandler *handler.ProductMessageParserHandler
	DeviceHandler               *handler.DeviceHandler
	SceneLinkageHandler         *handler.SceneLinkageHandler
	SceneLinkageDetailHandler   *handler.SceneLinkageDetailHandler
	AlertHandler                *handler.AlertHandler
	OTAHandler                  *handler.OTAHandler
	DeviceEventHandler          *handler.DeviceEventHandler
}
