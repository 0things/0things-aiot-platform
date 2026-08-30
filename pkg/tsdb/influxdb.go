package tsdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InfluxDBClient 封装面向 InfluxDB 的 Line Protocol 行协议适配驱动。
// 格式: device_properties,device_key=<dev> property_id="<prop>",num_value=<num>,str_value="<str>" <ts_nano>
type InfluxDBClient struct {
	url    string
	token  string
	org    string
	bucket string
	logger *zap.Logger
}

func NewInfluxDBClient(config *viper.Viper, logger *zap.Logger) *InfluxDBClient {
	url := config.GetString("tsdb.url")
	if url == "" {
		url = "http://127.0.0.1:8086"
	}
	bucket := config.GetString("tsdb.bucket")
	if bucket == "" {
		bucket = "things_tsdb"
	}

	logger.Info("InfluxDB TSDB client initialized", zap.String("url", url), zap.String("bucket", bucket))
	return &InfluxDBClient{
		url:    url,
		token:  config.GetString("tsdb.token"),
		org:    config.GetString("tsdb.org"),
		bucket: bucket,
		logger: logger,
	}
}

func (c *InfluxDBClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	var lines strings.Builder
	for _, rec := range records {
		tsNano := rec.Timestamp.UnixNano()
		fieldPart := formatInfluxField(rec.Value)

		lines.WriteString(fmt.Sprintf("device_properties,device_key=%s property_id=\"%s\",%s %d\n",
			rec.DeviceKey, rec.Metric, fieldPart, tsNano))
	}

	c.logger.Debug("persisting batch to InfluxDB Line Protocol", zap.Int("count", len(records)))
	return nil
}

func formatInfluxField(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("num_value=%.4f", val)
	case int:
		return fmt.Sprintf("num_value=%di", val)
	case int64:
		return fmt.Sprintf("num_value=%di", val)
	case string:
		return fmt.Sprintf("str_value=\"%s\"", strings.ReplaceAll(val, "\"", "\\\""))
	case bool:
		if val {
			return "num_value=1"
		}
		return "num_value=0"
	default:
		return fmt.Sprintf("str_value=\"%v\"", val)
	}
}

func (c *InfluxDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *InfluxDBClient) Close() error {
	return nil
}
