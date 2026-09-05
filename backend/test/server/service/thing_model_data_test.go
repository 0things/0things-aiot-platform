package service_test

import (
	"testing"
	"time"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/test/server/testutil"

	"github.com/stretchr/testify/require"
)

func TestDeviceServiceInvocationServiceList(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	output := `{"ok":true}`
	now := time.Now().Truncate(time.Second)
	require.NoError(t, db.Create(&model.DeviceServiceInvocation{UUID: "call-1", DeviceID: 1, ServiceIdentifier: "reboot", ServiceName: "Reboot", InputParams: `{}`, OutputParams: &output, InvokedAt: now}).Error)
	require.NoError(t, db.Create(&model.DeviceServiceInvocation{UUID: "call-2", DeviceID: 1, ServiceIdentifier: "set", ServiceName: "Set", InputParams: `{"on":true}`, InvokedAt: now.Add(-time.Hour)}).Error)

	items, total, err := testutil.NewTestThingModelDataService(db).ListServiceInvocations(testutil.ContextWithTenant(t.Context(), 1), dto.ListDeviceServiceInvocationsQuery{
		DeviceKey: "D001", ServiceIdentifier: "reboot", StartAt: ptrTime(now.Add(-time.Minute)), EndAt: ptrTime(now.Add(time.Minute)), Page: 1, PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, items, 1)
	require.Equal(t, "call-1", items[0].UUID)
	require.NotNil(t, items[0].OutputParams)
}

func ptrTime(value time.Time) *time.Time { return &value }
