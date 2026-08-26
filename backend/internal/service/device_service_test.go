package service

import (
	"bytes"
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func newDeviceSvc(t *testing.T) (*DeviceService, *gorm.DB, context.Context) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Product{}, &model.Device{}, &model.DeviceState{}, &model.DeviceTag{},
		&model.DeviceShadow{}, &model.DeviceShadowHistory{}, &model.DevicePushRecord{},
	))
	require.NoError(t, db.Create(&model.Product{ID: 1, ProductKey: "P001", Name: "Test", OrganizationID: 1}).Error)
	require.NoError(t, db.Create(&model.Device{ID: 1, DeviceKey: "D001", Name: "Test Device", ProductID: 1, OrganizationID: 1, Enabled: true}).Error)
	require.NoError(t, db.Create(&model.DeviceState{ID: 1, DeviceKey: "D001", State: "online"}).Error)

	iotDB := &repository.IoTDB{DB: db}
	kafkaSvc := &KafkaService{enabled: false}
	svc := NewDeviceService(
		repository.NewDeviceRepository(iotDB, &repository.IoTRedis{}),
		repository.NewProductRepository(iotDB),
		repository.NewDeviceTagRepository(iotDB),
		repository.NewDeviceShadowRepository(iotDB),
		repository.NewPushRecordRepository(iotDB),
		kafkaSvc,
	)
	return svc, db, context.Background()
}

func TestDeviceService_CreateDevice(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.CreateDevice(ctx, &model.Device{Name: "dev", ProductID: 1})
	require.NoError(t, err)
	require.NotEmpty(t, d.DeviceKey)
}

func TestDeviceService_CreateDevice_InvalidMetadata(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.CreateDevice(ctx, &model.Device{Name: "dev", ProductID: 1, Metadata: "not-json"})
	require.Error(t, err)
}

func TestDeviceService_CreateDevice_ProductNotFound(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.CreateDevice(ctx, &model.Device{Name: "dev", ProductID: 999})
	require.Error(t, err)
}

func TestDeviceService_Device(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.Device(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "D001", d.DeviceKey)
}

func TestDeviceService_Device_NotFound(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.Device(ctx, 999)
	require.Error(t, err)
}

func TestDeviceService_DeviceByKey(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.DeviceByKey(ctx, "D001")
	require.NoError(t, err)
	require.Equal(t, int64(1), d.ID)
}

func TestDeviceService_ListDevices(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	ds, n, err := svc.ListDevices(ctx, 1, 10, 0, nil, nil, "")
	require.NoError(t, err)
	require.Equal(t, int64(1), n)
	require.Len(t, ds, 1)
}

func TestDeviceService_UpdateDevice(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.UpdateDevice(ctx, 1, "renamed", "", "")
	require.NoError(t, err)
	require.Equal(t, "renamed", d.Name)
}

func TestDeviceService_UpdateDevice_StateTransition(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.UpdateDevice(ctx, 1, "", "offline", "")
	require.NoError(t, err)
	require.Equal(t, "offline", d.State.State)
}

func TestDeviceService_UpdateDevice_InvalidTransition(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.UpdateDevice(ctx, 1, "", "inactive", "")
	require.Error(t, err)
}

func TestDeviceService_Activate(t *testing.T) {
	svc, db, ctx := newDeviceSvc(t)
	require.NoError(t, db.Create(&model.Device{ID: 2, DeviceKey: "D002", Name: "d2", ProductID: 1, OrganizationID: 1, Enabled: true}).Error)
	require.NoError(t, db.Create(&model.DeviceState{ID: 2, DeviceKey: "D002", State: "inactive"}).Error)
	d, err := svc.Activate(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, "offline", d.State.State)
}

func TestDeviceService_Activate_AlreadyActivated(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.Activate(ctx, 1)
	require.Error(t, err)
}

func TestDeviceService_SetEnabled(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	d, err := svc.SetEnabled(ctx, 1, false)
	require.NoError(t, err)
	require.False(t, d.Enabled)
}

func TestDeviceService_Stats(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	s, err := svc.Stats(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), s.TotalDevices)
}

func TestDeviceService_Tags(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.SetTags(ctx, "1", map[string]string{"env": "prod"}, false)
	require.NoError(t, err)
	tags, err := svc.Tags(ctx, "1")
	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.NoError(t, svc.RemoveTags(ctx, "1", []string{"env"}))
}

func TestDeviceService_SetTags_InvalidKey(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.SetTags(ctx, "1", map[string]string{"": "v"}, false)
	require.Error(t, err)
	_, err = svc.SetTags(ctx, "1", map[string]string{"123": "v"}, false)
	require.Error(t, err)
}

func TestDeviceService_Shadow(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	desired := map[string]any{"power": "on"}
	sh, err := svc.MutateShadow(ctx, "D001", 0, "app", &desired, nil, false)
	require.NoError(t, err)
	require.NotZero(t, sh.Version)
	got, err := svc.Shadow(ctx, "D001")
	require.NoError(t, err)
	require.NotNil(t, got)
	hist, err := svc.ShadowHistory(ctx, "D001")
	require.NoError(t, err)
	require.NotNil(t, hist)
}

func TestDeviceService_Telemetry_NoRedis(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, err := svc.Telemetry(ctx, "D001")
	require.Error(t, err)
}

func TestDeviceService_MQTT(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	p, err := svc.MQTT(ctx, "D001")
	require.NoError(t, err)
	require.Equal(t, "D001", p.Username)
}

func TestDeviceService_DeleteAndRestore(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	require.NoError(t, svc.DeleteDevice(ctx, 1))
	_, err := svc.Device(ctx, 1)
	require.Error(t, err)
	d, err := svc.RestoreDevice(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), d.ID)
}

func TestDeviceService_SimulatePush(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	rec, err := svc.SimulatePush(ctx, "D001", "payload", "u1")
	require.NoError(t, err)
	require.Equal(t, "payload", rec.Payload)
	recs, n, err := svc.ListPushRecords(ctx, "D001", 1, 10, "", "")
	require.NoError(t, err)
	require.Equal(t, int64(1), n)
	require.Len(t, recs, 1)
	got, err := svc.PushRecord(ctx, rec.ID)
	require.NoError(t, err)
	require.Equal(t, rec.ID, got.ID)
	deleted, err := svc.ClearPushRecords(ctx, "D001", nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
}

func TestDeviceService_BatchTemplate(t *testing.T) {
	svc, _, _ := newDeviceSvc(t)
	b, err := svc.BatchTemplate()
	require.NoError(t, err)
	require.NotEmpty(t, b)
}

func TestDeviceService_BatchCreate(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "productKey")
	f.SetCellValue("Sheet1", "B1", "deviceName")
	f.SetCellValue("Sheet1", "A2", "P001")
	f.SetCellValue("Sheet1", "B2", "NEWDEV")
	var buf bytes.Buffer
	require.NoError(t, f.Write(&buf))
	n, errs, err := svc.BatchCreate(ctx, buf.Bytes())
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Empty(t, errs)
}

func TestDeviceService_BatchCreate_NoData(t *testing.T) {
	svc, _, ctx := newDeviceSvc(t)
	_, _, err := svc.BatchCreate(ctx, []byte("not-an-xlsx"))
	require.Error(t, err)
}
