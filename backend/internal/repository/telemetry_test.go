package repository

import (
	"context"
	"testing"
	"time"

	"0things/pkg/tsdb"
	"aiot-backend/pkg/log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTelemetryRepository_QueryLatest(t *testing.T) {
	testLogger := &log.Logger{Logger: zap.NewNop()}
	client := tsdb.NewMockClient(testLogger.Logger)
	now := time.Now()
	if err := client.WriteBatch(context.Background(), []tsdb.Record{
		{DeviceKey: "device-1", Metric: "temperature", Value: 21.5, Timestamp: now.Add(-2 * time.Minute)},
		{DeviceKey: "device-1", Metric: "temperature", Value: 24.8, Timestamp: now.Add(-time.Minute)},
		{DeviceKey: "device-1", Metric: "humidity", Value: 55, Timestamp: now.Add(-30 * time.Second)},
	}); err != nil {
		t.Fatal(err)
	}

	repo := NewTelemetryRepositoryWithClient(client, viper.New(), testLogger)
	points, err := repo.QueryLatest(context.Background(), "device-1", []string{"temperature", "humidity", "not-reported"})
	if err != nil {
		t.Fatal(err)
	}
	if points["temperature"].Value != 24.8 {
		t.Fatalf("expected newest temperature, got %#v", points["temperature"])
	}
	if points["humidity"].Value != 55 {
		t.Fatalf("expected humidity point, got %#v", points["humidity"])
	}
	if _, found := points["not-reported"]; found {
		t.Fatal("expected no fabricated point for an unreported property")
	}
}
