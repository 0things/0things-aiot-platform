package router

import (
	"0things-backend/internal/handler"
	"0things-backend/pkg/jwt"
	"0things-backend/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger            *log.Logger
	Config            *viper.Viper
	JWT               *jwt.JWT
	UserHandler       *handler.UserHandler
	ProductHandler    *handler.ProductHandler
	ProductTSLHandler *handler.ProductTSLHandler
	DeviceHandler     *handler.DeviceHandler
	RuleHandler       *handler.RuleHandler
	AlertHandler      *handler.AlertHandler
	OTAHandler        *handler.OTAHandler
	DeviceEventHandler *handler.DeviceEventHandler
}
