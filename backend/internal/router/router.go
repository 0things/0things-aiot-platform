package router

import (
	"aiot-backend/internal/handler"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/logto"

	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger                      *log.Logger
	Config                      *viper.Viper
	Logto                       *logto.Verifier
	ProductHandler              *handler.ProductHandler
	CategoryHandler             *handler.CategoryHandler
	ProductTSLHandler           *handler.ProductTSLHandler
	ProductMessageParserHandler *handler.ProductMessageParserHandler
	DeviceHandler               *handler.DeviceHandler
	OTAHandler                  *handler.OTAHandler
	FileHandler                 *handler.FileHandler
	DeviceEventHandler          *handler.DeviceEventHandler
	ThingModelDataHandler       *handler.ThingModelDataHandler
	ProtocolHandler             *handler.ProtocolHandler
	DeviceGroupHandler          *handler.DeviceGroupHandler
	TelemetryHandler            *handler.TelemetryHandler
	RuleNodeDefinitionHandler   *handler.RuleNodeDefinitionHandler
	RuleChainHandler            *handler.RuleChainHandler
	OrganizationHandler         *handler.OrganizationHandler
}
