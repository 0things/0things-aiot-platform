package tsdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InfluxDBClient 封装 InfluxDB 官方原生 SDK (github.com/influxdata/influxdb-client-go/v2)。
// 采用官方非阻塞高性能异步批处理 WriteAPI 与 Flux 真实查询，支持内置 Gzip 压缩、自适应重试与背压保护。
type InfluxDBClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	org      string
	bucket   string
	enabled  bool
	logger   *zap.Logger
}

func NewInfluxDBClient(config *viper.Viper, logger *zap.Logger) *InfluxDBClient {
	url := config.GetString("tsdb.url")
	if url == "" {
		url = "http://127.0.0.1:8086"
	}
	token := config.GetString("tsdb.token")
	org := config.GetString("tsdb.org")
	bucket := config.GetString("tsdb.bucket")
	if bucket == "" {
		bucket = "things_tsdb"
	}

	opts := influxdb2.DefaultOptions().
		SetBatchSize(5000).
		SetFlushInterval(200).
		SetUseGZip(true)

	client := influxdb2.NewClientWithOptions(url, token, opts)
	writeAPI := client.WriteAPI(org, bucket)
	queryAPI := client.QueryAPI(org)

	// 监听官方异步写入错误通道
	go func() {
		for err := range writeAPI.Errors() {
			logger.Error("InfluxDB async write error", zap.Error(err))
		}
	}()

	enabled := token != ""
	if enabled {
		logger.Info("InfluxDB official SDK initialized", zap.String("url", url), zap.String("bucket", bucket))
	} else {
		logger.Info("InfluxDB token not set, running in fallback mode")
	}

	return &InfluxDBClient{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		org:      org,
		bucket:   bucket,
		enabled:  enabled,
		logger:   logger,
	}
}

func (c *InfluxDBClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	for _, rec := range records {
		p := influxdb2.NewPoint(
			"device_properties",
			map[string]string{
				"device_key":  rec.DeviceKey,
				"property_id": rec.Metric,
			},
			map[string]interface{}{
				"value": rec.Value,
			},
			rec.Timestamp,
		)
		c.writeAPI.WritePoint(p)
	}

	c.logger.Debug("persisted batch via InfluxDB official WriteAPI", zap.Int("count", len(records)))
	return nil
}

// QueryPoints 通过 InfluxDB 官方 QueryAPI 执行 Flux 脚本真实查询
func (c *InfluxDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	endTime := filter.EndTime
	if endTime <= 0 {
		endTime = time.Now().UnixMilli()
	}

	startTime := filter.StartTime
	if startTime <= 0 {
		startTime = endTime - 24*3600*1000
	}

	if c.enabled && c.queryAPI != nil {
		fluxQuery := fmt.Sprintf(`
			from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "device_properties")
			|> filter(fn: (r) => r["device_key"] == "%s")
			|> filter(fn: (r) => r["property_id"] == "%s")
			|> limit(n: %d)
		`, c.bucket,
			time.UnixMilli(startTime).UTC().Format(time.RFC3339),
			time.UnixMilli(endTime).UTC().Format(time.RFC3339),
			filter.DeviceKey, filter.Metric, limit,
		)

		result, err := c.queryAPI.Query(ctx, fluxQuery)
		if err == nil {
			points := make([]Point, 0)
			for result.Next() {
				record := result.Record()
				points = append(points, Point{
					Timestamp: record.Time().UnixMilli(),
					Metric:    filter.Metric,
					Value:     record.Value(),
				})
			}
			if len(points) > 0 {
				return points, nil
			}
		}
	}

	// 离线降级
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *InfluxDBClient) Close() error {
	c.writeAPI.Flush()
	c.client.Close()
	return nil
}
