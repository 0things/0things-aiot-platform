package tsdb

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DriverBuilder 定义时序驱动实例构造函数契约。
type DriverBuilder func(config *viper.Viper, logger *zap.Logger) Client

// driversRegistry 统一驱动注册表，键类型为强类型 DriverType 枚举。
var driversRegistry = map[DriverType]DriverBuilder{
	DriverTypeTDengine: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewTDengineClient(config, logger)
	},
	DriverTypeIoTDB: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewIoTDBClient(config, logger)
	},
	DriverTypeTimescaleDB: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewTimescaleDBClient(config, logger)
	},
	DriverTypePostgreSQL: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewTimescaleDBClient(config, logger)
	},
	DriverTypeInfluxDB: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewInfluxDBClient(config, logger)
	},
	DriverTypeClickHouse: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewClickHouseClient(config, logger)
	},
	DriverTypeMock: func(config *viper.Viper, logger *zap.Logger) Client {
		return NewMockClient(logger)
	},
}

// NewClient 根据配置项 tsdb.type 统一分发并初始化对应的时序数据库客户端。
func NewClient(config *viper.Viper, logger *zap.Logger) Client {
	rawType := config.GetString("tsdb.type")
	driverType, err := ParseDriverType(rawType)
	if err != nil {
		logger.Warn("unrecognized or empty tsdb.type in configuration, falling back to mock driver",
			zap.String("configured_type", rawType),
			zap.String("fallback", DriverTypeMock.String()),
			zap.Error(err),
		)
		driverType = DriverTypeMock
	}

	builder, ok := driversRegistry[driverType]
	if !ok {
		return NewMockClient(logger)
	}

	logger.Info("initializing TSDB client", zap.String("driver", driverType.String()))
	return builder(config, logger)
}
