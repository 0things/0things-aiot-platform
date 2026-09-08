package router

import (
	"transport-http/internal/handler"
	"transport-http/pkg/jwt"
	"transport-http/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
