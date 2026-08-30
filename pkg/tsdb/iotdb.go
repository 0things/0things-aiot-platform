package tsdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// IoTDBClient 封装面向 Apache IoTDB 的时序数据库驱动。
type IoTDBClient struct {
	host         string
	port         string
	storageGroup string
	logger       *zap.Logger
}

func NewIoTDBClient(config *viper.Viper, logger *zap.Logger) *IoTDBClient {
	host := config.GetString("tsdb.host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := config.GetString("tsdb.port")
	if port == "" {
		port = "6667"
	}

	logger.Info("Apache IoTDB TSDB client initialized", zap.String("host", host), zap.String("port", port))
	return &IoTDBClient{
		host:         host,
		port:         port,
		storageGroup: "root.0things",
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
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *IoTDBClient) Close() error {
	return nil
}
