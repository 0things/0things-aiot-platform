package router

import (
	"coap-transport/internal/handler"
	"coap-transport/pkg/jwt"
	"coap-transport/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
