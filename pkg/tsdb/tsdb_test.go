package tsdb

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTSDB_AllDrivers(t *testing.T) {
	logger := zap.NewNop()
	drivers := []string{"tdengine", "iotdb", "timescaledb", "influxdb", "clickhouse", "mock"}

	records := []Record{
		{DeviceKey: "dev_all_01", Metric: "temperature", Value: 26.5, Timestamp: time.Now()},
		{DeviceKey: "dev_all_01", Metric: "door_state", Value: "CLOSED", Timestamp: time.Now()},
	}

	filter := QueryFilter{
		DeviceKey: "dev_all_01",
		Metric:    "temperature",
		Limit:     10,
	}

	for _, driverName := range drivers {
		t.Run(driverName, func(t *testing.T) {
			v := viper.New()
			v.Set("tsdb.type", driverName)
			v.Set("tsdb.database", "things_tsdb")

			client := NewClient(v, logger)
			defer client.Close()

			// 1. 验证写入
			if err := client.WriteBatch(context.Background(), records); err != nil {
				t.Fatalf("driver %s WriteBatch failed: %v", driverName, err)
			}

			// 2. 验证查询
			points, err := client.QueryPoints(context.Background(), filter)
			if err != nil {
				t.Fatalf("driver %s QueryPoints failed: %v", driverName, err)
			}
			if len(points) == 0 {
				t.Fatalf("driver %s expected points, got empty", driverName)
			}
		})
	}
}
