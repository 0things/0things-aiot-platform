package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestMemoryShadowRepository(t *testing.T) {
	logger := zaptest.NewLogger(t)
	repo := NewShadowRepository(nil, logger)
	ctx := context.Background()

	// 1. Get non-existent shadow
	shadow, err := repo.GetShadow(ctx, "dev-none")
	require.NoError(t, err)
	assert.Nil(t, shadow)

	// 2. Update shadow for new device
	ts1 := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	err = repo.UpdateShadow(ctx, "dev-001", map[string]interface{}{
		"temp": 25.5,
		"mode": "auto",
	}, ts1)
	require.NoError(t, err)

	shadow, err = repo.GetShadow(ctx, "dev-001")
	require.NoError(t, err)
	require.NotNil(t, shadow)
	assert.Equal(t, "dev-001", shadow.DeviceKey)
	assert.Equal(t, 25.5, shadow.Attributes["temp"])
	assert.Equal(t, "auto", shadow.Attributes["mode"])
	assert.Equal(t, ts1, shadow.LastSeen)

	// 3. Update existing device shadow with new attributes
	ts2 := ts1.Add(time.Minute)
	err = repo.UpdateShadow(ctx, "dev-001", map[string]interface{}{
		"temp":     26.0,
		"humidity": 60,
	}, ts2)
	require.NoError(t, err)

	shadow, err = repo.GetShadow(ctx, "dev-001")
	require.NoError(t, err)
	require.NotNil(t, shadow)
	assert.Equal(t, 26.0, shadow.Attributes["temp"])
	assert.Equal(t, "auto", shadow.Attributes["mode"])
	assert.Equal(t, 60, shadow.Attributes["humidity"])
	assert.Equal(t, ts2, shadow.LastSeen)
}
