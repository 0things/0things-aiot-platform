package job

import (
	"context"
	"testing"
	"time"

	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewOTACommandConsumer_NoBrokers(t *testing.T) {
	config := viper.New()
	logger := &log.Logger{Logger: zap.NewNop()}

	consumer, err := NewOTACommandConsumer(config, nil, nil, logger)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = consumer.Start(ctx)
	assert.NoError(t, err)

	consumer.Stop()
}
