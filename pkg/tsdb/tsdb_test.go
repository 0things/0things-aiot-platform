package tsdb

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTSDB_DedicatedConfigs(t *testing.T) {
	yamlConfig := []byte(`
tsdb:
  type: tdengine
  tdengine:
    dsn: "root:taosdata@ws(127.0.0.1:6041)/things_tsdb"
    database: "things_tsdb"
    table: "device_properties"
  sqlite:
    path: "data/things_tsdb.db"
    table: "device_properties"
  influxdb:
    url: "http://127.0.0.1:8086"
    token: "my-token"
    org: "my-org"
    bucket: "my-bucket"
  iotdb:
    host: "127.0.0.1"
    port: "6667"
    storage_group: "root.custom"
`)

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlConfig)); err != nil {
		t.Fatalf("failed to read yaml config: %v", err)
	}

	rootCfg, err := LoadRootConfigFromViper(v)
	if err != nil {
		t.Fatalf("LoadRootConfigFromViper failed: %v", err)
	}

	if rootCfg.Type != DriverTypeTDengine {
		t.Fatalf("expected DriverTypeTDengine, got %v", rootCfg.Type)
	}
	if rootCfg.TDengine.DSN != "root:taosdata@ws(127.0.0.1:6041)/things_tsdb" {
		t.Fatalf("unexpected TDengine DSN: %s", rootCfg.TDengine.DSN)
	}
	if rootCfg.SQLite.Path != "data/things_tsdb.db" {
		t.Fatalf("unexpected SQLite path: %s", rootCfg.SQLite.Path)
	}
	if rootCfg.InfluxDB.URL != "http://127.0.0.1:8086" || rootCfg.InfluxDB.Token != "my-token" || rootCfg.InfluxDB.Bucket != "my-bucket" {
		t.Fatalf("unexpected InfluxDB params: %s / %s / %s", rootCfg.InfluxDB.URL, rootCfg.InfluxDB.Token, rootCfg.InfluxDB.Bucket)
	}
	if rootCfg.IoTDB.Host != "127.0.0.1" || rootCfg.IoTDB.StorageGroup != "root.custom" {
		t.Fatalf("unexpected IoTDB host/group: %s / %s", rootCfg.IoTDB.Host, rootCfg.IoTDB.StorageGroup)
	}
}

func TestTSDB_EnumAndAllDrivers(t *testing.T) {
	logger := zap.NewNop()
	drivers := AllDriverTypes()

	if len(drivers) == 0 {
		t.Fatal("expected non-empty driver types")
	}

	for _, dt := range drivers {
		if !dt.IsValid() {
			t.Fatalf("driver type %s should be valid", dt)
		}
		parsed, err := ParseDriverType(dt.String())
		if err != nil || parsed != dt {
			t.Fatalf("ParseDriverType(%s) failed: %v", dt, err)
		}
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

	for _, driverType := range drivers {
		driverName := driverType.String()
		t.Run(driverName, func(t *testing.T) {
			v := viper.New()
			v.Set("tsdb.type", driverName)
			v.Set("tsdb.database", "things_tsdb")

			client := NewClient(v, logger)
			defer client.Close()

			if err := client.WriteBatch(context.Background(), records); err != nil {
				t.Fatalf("driver %s WriteBatch failed: %v", driverName, err)
			}

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
