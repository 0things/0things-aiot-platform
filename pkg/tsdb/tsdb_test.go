package tsdb

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTSDB_MockClient(t *testing.T) {
	logger := zap.NewNop()
	v := viper.New()
	v.Set("tsdb.type", "mock")

	client := NewClient(v, logger)
	defer client.Close()

	// 1. 测试写入
	records := []Record{
		{
			DeviceKey: "dev_mock_01",
			Metric:    "temperature",
			Value:     28.5,
			Timestamp: time.Now(),
		},
	}

	if err := client.WriteBatch(context.Background(), records); err != nil {
		t.Fatalf("WriteBatch failed: %v", err)
	}

	// 2. 测试查询
	filter := QueryFilter{
		DeviceKey: "dev_mock_01",
		Metric:    "temperature",
		Limit:     10,
	}

	points, err := client.QueryPoints(context.Background(), filter)
	if err != nil {
		t.Fatalf("QueryPoints failed: %v", err)
	}

	if len(points) == 0 {
		t.Fatalf("expected at least 1 point, got 0")
	}
}

func TestTSDB_TDengineClient(t *testing.T) {
	logger := zap.NewNop()
	v := viper.New()
	v.Set("tsdb.type", "tdengine")
	v.Set("tsdb.database", "things_tsdb")

	client := NewClient(v, logger)
	defer client.Close()

	// 测试批量写入排队
	records := []Record{
		{DeviceKey: "dev_td_01", Metric: "voltage", Value: 220.0, Timestamp: time.Now()},
		{DeviceKey: "dev_td_01", Metric: "status", Value: "NORMAL", Timestamp: time.Now()},
	}

	if err := client.WriteBatch(context.Background(), records); err != nil {
		t.Fatalf("WriteBatch on TDengine failed: %v", err)
	}
}
