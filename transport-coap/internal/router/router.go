package router

import (
	"transport-coap/internal/handler"
	"transport-coap/pkg/jwt"
	"transport-coap/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
