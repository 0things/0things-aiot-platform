//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"telemetry-service/internal/biz"
	"telemetry-service/internal/conf"
	"telemetry-service/internal/data"
	"telemetry-service/internal/data/kafka"
	"telemetry-service/internal/server"
	"telemetry-service/internal/service"
	"telemetry-service/internal/service/consumer"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		kafka.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		consumer.ProviderSet,
		newApp,
	))
}
