package router

import (
	"transport-mqtt/internal/handler"
	"transport-mqtt/pkg/jwt"
	"transport-mqtt/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
