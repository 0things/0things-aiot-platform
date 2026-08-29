package router

import (
	"rule-engine/internal/handler"
	"rule-engine/pkg/jwt"
	"rule-engine/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
