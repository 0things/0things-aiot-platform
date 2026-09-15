package repository

import (
	"context"
	"testing"
	"time"

	"data-engine/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newTestDeviceEventRepository(t *testing.T) (DeviceEventRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.DeviceEvent{})
	require.NoError(t, err)

	repo := NewRepository(zap.NewNop(), db)
	return NewDeviceEventRepository(repo), db
}

func TestDeviceEventRepository_Create(t *testing.T) {
	ctx := context.Background()
	deviceEventRepo, db := newTestDeviceEventRepository(t)

	evt := &model.DeviceEvent{
		UUID:            "evt-uuid-001",
		DeviceKey:       "dev_sensor_01",
		ProductKey:      "pk_sensor_01",
		EventIdentifier: "overheat_alarm",
		EventType:       "alert",
		EventAt:         time.Now().UnixMilli(),
		Data:            `{"temperature":88.6,"threshold":80.0}`,
	}

	err := deviceEventRepo.Create(ctx, evt)
	require.NoError(t, err)
	assert.NotZero(t, evt.ID)

	var saved model.DeviceEvent
	err = db.First(&saved, "uuid = ?", "evt-uuid-001").Error
	require.NoError(t, err)
	assert.Equal(t, "dev_sensor_01", saved.DeviceKey)
	assert.Equal(t, "pk_sensor_01", saved.ProductKey)
	assert.Equal(t, "overheat_alarm", saved.EventIdentifier)
	assert.Equal(t, "alert", saved.EventType)
	assert.Equal(t, evt.EventAt, saved.EventAt)
	assert.Equal(t, evt.Data, saved.Data)
}

func TestDeviceEventRepository_Create_NilEvent(t *testing.T) {
	ctx := context.Background()
	deviceEventRepo, _ := newTestDeviceEventRepository(t)

	err := deviceEventRepo.Create(ctx, nil)
	assert.NoError(t, err)
}
