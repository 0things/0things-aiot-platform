package router

import (
	"http-transport/internal/handler"
	"http-transport/pkg/jwt"
	"http-transport/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
