package tsdb

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewClient 根据配置文件中的 tsdb.type 动态构建并返回统一时序客户端实例。
// 支持: "tdengine" | "iotdb" | "timescaledb" | "influxdb" | "clickhouse" | "mock"
func NewClient(config *viper.Viper, logger *zap.Logger) Client {
	tsdbType := strings.ToLower(strings.TrimSpace(config.GetString("tsdb.type")))

	switch tsdbType {
	case "tdengine", "taos":
		return NewTDengineClient(config, logger)
	case "iotdb":
		return NewIoTDBClient(config, logger)
	case "timescaledb", "timescale", "postgres", "postgresql":
		return NewTimescaleDBClient(config, logger)
	case "influxdb", "influx":
		return NewInfluxDBClient(config, logger)
	case "clickhouse", "ch":
		return NewClickHouseClient(config, logger)
	default:
		return NewMockClient(logger)
	}
}
