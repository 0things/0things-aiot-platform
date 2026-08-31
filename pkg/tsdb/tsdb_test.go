package tsdb

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTSDB_EnumAndAllDrivers(t *testing.T) {
	logger := zap.NewNop()
	drivers := AllDriverTypes()

	if len(drivers) == 0 {
		t.Fatal("expected non-empty driver types")
	}

	// 1. 验证枚举的解析与合法性判断
	for _, dt := range drivers {
		if !dt.IsValid() {
			t.Fatalf("driver type %s should be valid", dt)
		}
		parsed, err := ParseDriverType(dt.String())
		if err != nil || parsed != dt {
			t.Fatalf("ParseDriverType(%s) failed: %v", dt, err)
		}
	}

	// 验证非法枚举报错
	if _, err := ParseDriverType("unknown_invalid_db"); err == nil {
		t.Fatal("expected error for invalid driver type")
	}

	records := []Record{
		{DeviceKey: "dev_all_01", Metric: "temperature", Value: 26.5, Timestamp: time.Now()},
		{DeviceKey: "dev_all_01", Metric: "door_state", Value: "CLOSED", Timestamp: time.Now()},
	}

	filter := QueryFilter{
		DeviceKey: "dev_all_01",
		Metric:    "temperature",
		Limit:     10,
	}

	// 2. 验证所有枚举驱动的写入与查询能力
	for _, driverType := range drivers {
		driverName := driverType.String()
		t.Run(driverName, func(t *testing.T) {
			v := viper.New()
			v.Set("tsdb.type", driverName)
			v.Set("tsdb.database", "things_tsdb")

			client := NewClient(v, logger)
			defer client.Close()

			// 验证写入
			if err := client.WriteBatch(context.Background(), records); err != nil {
				t.Fatalf("driver %s WriteBatch failed: %v", driverName, err)
			}

			// 验证查询
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
