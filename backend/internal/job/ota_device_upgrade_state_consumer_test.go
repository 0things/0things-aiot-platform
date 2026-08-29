package job

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewOTAReportConsumer_NoBrokers(t *testing.T) {
	config := viper.New()
	logger := &log.Logger{Logger: zap.NewNop()}

	consumer, err := NewOTAReportConsumer(config, nil, logger)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = consumer.Start(ctx)
	assert.NoError(t, err)

	consumer.Stop()
}

func TestOTADeviceUpgradeReportEvent_JSONSerialization(t *testing.T) {
	event := OTADeviceUpgradeReportEvent{
		BatchID:   "batch-123",
		DeviceKey: "dev-key-456",
		Status:    "failed",
		Version:   "1.0.0",
		Progress:  45,
		Desc:      "checksum verification failed",
	}

	data, err := json.Marshal(event)
	assert.NoError(t, err)

	var parsed OTADeviceUpgradeReportEvent
	assert.NoError(t, json.Unmarshal(data, &parsed))

	assert.Equal(t, "batch-123", parsed.BatchID)
	assert.Equal(t, "dev-key-456", parsed.DeviceKey)
	assert.Equal(t, "failed", parsed.Status)
	assert.Equal(t, "1.0.0", parsed.Version)
	assert.Equal(t, int32(45), parsed.Progress)
	assert.Equal(t, "checksum verification failed", parsed.Desc)
}
