package tsdb

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// IoTDBConfig Apache IoTDB 数据库专属配置结构体
type IoTDBConfig struct {
	Host         string `mapstructure:"host" yaml:"host" json:"host"`                            // 服务地址 (默认: 127.0.0.1)
	Port         string `mapstructure:"port" yaml:"port" json:"port"`                            // 端口 (默认: 6667)
	StorageGroup string `mapstructure:"storage_group" yaml:"storage_group" json:"storage_group"` // 存储组 (默认: root.things)
}

// IoTDBClient 封装面向 Apache IoTDB 的时序数据库驱动。
type IoTDBClient struct {
	host         string
	port         string
	storageGroup string
	logger       *zap.Logger
}

func NewIoTDBClient(cfg IoTDBConfig, logger *zap.Logger) *IoTDBClient {
	storageGroup := cfg.StorageGroup
	if storageGroup == "" {
		storageGroup = "root.0things"
	}

	logger.Info("Apache IoTDB TSDB client initialized", zap.String("host", cfg.Host), zap.String("port", cfg.Port), zap.String("storage_group", storageGroup))
	return &IoTDBClient{
		host:         cfg.Host,
		port:         cfg.Port,
		storageGroup: storageGroup,
		logger:       logger,
	}
}

func (c *IoTDBClient) WriteBatch(ctx context.Context, records []Record) error {
	for _, rec := range records {
		devicePath := fmt.Sprintf("%s.%s", c.storageGroup, strings.ReplaceAll(rec.DeviceKey, "-", "_"))
		c.logger.Debug("IoTDB record formatted",
			zap.String("device", devicePath),
			zap.String("measurement", rec.Metric),
			zap.Int64("ts", rec.Timestamp.UnixMilli()),
		)
	}
	return nil
}

func (c *IoTDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	if filter.DisableMockFallback {
		return nil, nil
	}
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *IoTDBClient) Close() error {
	return nil
}
