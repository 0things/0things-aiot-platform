package service_test

import (
	"context"
	"testing"

	"aiot-backend/internal/service"
	"aiot-backend/test/server/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTagServiceWithSeed(t *testing.T) (*service.DeviceService, context.Context) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)
	ctx := testutil.ContextWithTenant(context.Background(), 1)
	return svc, ctx
}

func TestDeviceService_SetTags_RejectsNumericKey(t *testing.T) {
	svc, ctx := newTagServiceWithSeed(t)

	_, err := svc.SetTags(ctx, "1", map[string]string{"12345": "x"}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "purely numeric")
}

func TestDeviceService_SetTags_RejectsEmptyKey(t *testing.T) {
	svc, ctx := newTagServiceWithSeed(t)

	_, err := svc.SetTags(ctx, "1", map[string]string{"": "x"}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag key")
}

func TestDeviceService_SetTags_RejectsOverlongKey(t *testing.T) {
	svc, ctx := newTagServiceWithSeed(t)

	overlong := make([]byte, 129)
	for i := range overlong {
		overlong[i] = 'a'
	}
	_, err := svc.SetTags(ctx, "1", map[string]string{string(overlong): "x"}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag key")
}

func TestDeviceService_SetTags_AcceptsValidKey(t *testing.T) {
	svc, ctx := newTagServiceWithSeed(t)

	_, err := svc.SetTags(ctx, "1", map[string]string{"room_no": "A1"}, false)
	require.NoError(t, err)

	tags, err := svc.Tags(ctx, "1")
	require.NoError(t, err)
	found := false
	for _, tg := range tags {
		if tg.Key == "room_no" && tg.Value == "A1" {
			found = true
		}
	}
	assert.True(t, found, "expected room_no tag to be persisted")
}

func TestDeviceService_SetTags_InvalidIDParam(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)
	ctx := context.Background()

	_, err := svc.SetTags(ctx, "not-a-number", map[string]string{"k": "v"}, false)
	require.Error(t, err)
}
