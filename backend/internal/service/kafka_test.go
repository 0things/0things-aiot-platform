package service

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func newTestLogger(t *testing.T) *log.Logger {
	t.Helper()
	conf := viper.New()
	conf.Set("log.log_file_name", filepath.Join(t.TempDir(), "test.log"))
	return log.NewLog(conf)
}

func TestKafkaService_DisabledWhenNoBrokers(t *testing.T) {
	logger := newTestLogger(t)
	conf := viper.New()

	svc, cleanup, err := NewKafkaService(conf, logger)
	require.NoError(t, err)
	require.NotNil(t, svc)
	require.False(t, svc.enabled)
	require.Nil(t, svc.client)

	ctx := context.Background()

	// All produce calls should return ErrKafkaNotInitialized
	err = svc.Produce(ctx, "test-topic", []byte("key"), []byte("value"))
	require.ErrorIs(t, err, ErrKafkaNotInitialized)

	err = svc.ProduceJSON(ctx, "test-topic", "key", map[string]string{"foo": "bar"})
	require.ErrorIs(t, err, ErrKafkaNotInitialized)

	var wg sync.WaitGroup
	wg.Add(1)
	svc.ProduceAsync(ctx, "test-topic", []byte("key"), []byte("value"), func(err error) {
		require.ErrorIs(t, err, ErrKafkaNotInitialized)
		wg.Done()
	})
	wg.Wait()

	wg.Add(1)
	svc.ProduceJSONAsync(ctx, "test-topic", "key", map[string]string{"foo": "bar"}, func(err error) {
		require.ErrorIs(t, err, ErrKafkaNotInitialized)
		wg.Done()
	})
	wg.Wait()

	err = svc.Flush(ctx)
	require.ErrorIs(t, err, ErrKafkaNotInitialized)

	// Close / cleanup should not panic
	svc.Close()
	cleanup()
}

func TestKafkaService_InvalidJSONSerialization(t *testing.T) {
	logger := newTestLogger(t)
	conf := viper.New()
	conf.Set("data.kafka.device.brokers", []string{"127.0.0.1:9092"})

	svc, cleanup, err := NewKafkaService(conf, logger)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()
	unmarshallableData := map[string]any{
		"invalid": make(chan int),
	}

	// ProduceJSON with unmarshallable data should return a json marshal error
	err = svc.ProduceJSON(ctx, "test-topic", "key", unmarshallableData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to marshal json")

	var wg sync.WaitGroup
	wg.Add(1)
	svc.ProduceJSONAsync(ctx, "test-topic", "key", unmarshallableData, func(err error) {
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to marshal json")
		wg.Done()
	})
	wg.Wait()
}

func TestKafkaService_EnabledInitialization(t *testing.T) {
	logger := newTestLogger(t)
	conf := viper.New()
	conf.Set("data.kafka.device.brokers", []string{"127.0.0.1:9092"})
	conf.Set("data.kafka.client_id", "test-client")

	svc, cleanup, err := NewKafkaService(conf, logger)
	require.NoError(t, err)
	require.NotNil(t, svc)
	require.True(t, svc.enabled)
	require.NotNil(t, svc.client)
	require.Equal(t, []string{"127.0.0.1:9092"}, svc.brokers)

	defer cleanup()
}

