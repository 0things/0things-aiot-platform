package router

import (
	"mqtt-transport/internal/handler"
	"mqtt-transport/pkg/jwt"
	"mqtt-transport/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
