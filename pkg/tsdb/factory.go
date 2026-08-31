package tsdb

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewClientFromRootConfig 根据包含各个 DB 专属子配置的 RootConfig 统一分发并初始化目标客户端
func NewClientFromRootConfig(rootCfg *RootConfig, logger *zap.Logger) Client {
	if rootCfg == nil {
		return NewMockClient(logger)
	}

	driverType, err := ParseDriverType(string(rootCfg.Type))
	if err != nil {
		logger.Warn("unrecognized tsdb driver type, falling back to mock",
			zap.String("configured_type", string(rootCfg.Type)),
			zap.String("fallback", DriverTypeMock.String()),
			zap.Error(err),
		)
		driverType = DriverTypeMock
	}

	logger.Info("initializing TSDB client with dedicated DB config", zap.String("driver", driverType.String()))

	switch driverType {
	case DriverTypeTDengine:
		return NewTDengineClient(rootCfg.TDengine, logger)
	case DriverTypeSQLite:
		return NewSQLiteClient(rootCfg.SQLite, logger)
	case DriverTypeTimescaleDB:
		return NewTimescaleDBClient(rootCfg.TimescaleDB, logger)
	case DriverTypePostgreSQL:
		return NewPostgreSQLClient(rootCfg.PostgreSQL, logger)
	case DriverTypeClickHouse:
		return NewClickHouseClient(rootCfg.ClickHouse, logger)
	case DriverTypeInfluxDB:
		return NewInfluxDBClient(rootCfg.InfluxDB, logger)
	case DriverTypeIoTDB:
		return NewIoTDBClient(rootCfg.IoTDB, logger)
	default:
		return NewMockClient(logger)
	}
}

// NewClient 从 Viper 实例中统一加载多 DB 配置并构建目标 TSDB 客户端
func NewClient(v *viper.Viper, logger *zap.Logger) Client {
	rootCfg, err := LoadRootConfigFromViper(v)
	if err != nil {
		logger.Warn("failed to parse tsdb root config from viper, falling back to mock", zap.Error(err))
		return NewMockClient(logger)
	}
	return NewClientFromRootConfig(rootCfg, logger)
}
