package service

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type testKafkaService struct{}

func (testKafkaService) Produce(context.Context, string, []byte, []byte) error              { return nil }
func (testKafkaService) ProduceJSON(context.Context, string, string, any) error             { return nil }
func (testKafkaService) ProduceAsync(context.Context, string, []byte, []byte, func(error))  {}
func (testKafkaService) ProduceJSONAsync(context.Context, string, string, any, func(error)) {}
func (testKafkaService) Flush(context.Context) error                                        { return nil }
func (testKafkaService) Close()                                                             {}

func newOTATestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(
		&model.OTAPackage{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{},
		&model.Product{}, &model.ProductProtocol{}, &model.Device{}, &model.DeviceState{},
	))
	return db
}

func newOTAServiceForTest(t *testing.T) (*OTAService, *gorm.DB) {
	db := newOTATestDB(t)
	iotDB := &repository.IoTDB{DB: db}
	repo := repository.NewOTARepository(iotDB)
	productRepo := repository.NewProductRepository(iotDB)
	deviceRepo := repository.NewDeviceRepository(iotDB, &repository.IoTRedis{})
	return NewOTAService(repo, productRepo, deviceRepo, testKafkaService{}), db
}

func seedOTAProductAndDevices(t *testing.T, db *gorm.DB) (int64, []string) {
	t.Helper()
	ctx := context.Background()
	product := &model.Product{ProductKey: "P001", Name: "Sensor", OrganizationID: 1}
	require.NoError(t, db.WithContext(ctx).Create(product).Error)
	dev1 := &model.Device{DeviceKey: "D1", Name: "d1", ProductID: product.ID, OrganizationID: 1}
	dev2 := &model.Device{DeviceKey: "D2", Name: "d2", ProductID: product.ID, OrganizationID: 1}
	require.NoError(t, db.WithContext(ctx).Create(dev1).Error)
	require.NoError(t, db.WithContext(ctx).Create(dev2).Error)
	return product.ID, []string{"D1", "D2"}
}

func TestOTAService_BatchUpgradeAndReportStatus(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))

	// BatchUpgrade creates batch and records, and marks package deploying.
	batch, err := svc.BatchUpgrade(ctx, pkg.UUID, keys)
	require.NoError(t, err)
	require.EqualValues(t, 2, batch.TargetDeviceCount)
	got, _ := svc.Get(ctx, pkg.UUID)
	require.Equal(t, "deploying", got.Status)

	// PendingPackageIDs should include this package.
	pending, err := svc.repo.PendingPackageIDs(ctx)
	require.NoError(t, err)
	require.Contains(t, pending, pkg.ID)
}

func TestOTAService_ReportStatusAggregation(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg := &model.OTAPackage{PackageName: "fw-2", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))
	_, err := svc.BatchUpgrade(ctx, pkg.UUID, keys)
	require.NoError(t, err)

	// One success, one still pending -> package stays deploying.
	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D1", "success"))
	got, _ := svc.Get(ctx, pkg.UUID)
	require.Equal(t, "deploying", got.Status)

	// All success -> package becomes success.
	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D2", "success"))
	got, _ = svc.Get(ctx, pkg.UUID)
	require.Equal(t, "success", got.Status)
}

func TestOTAService_ReportStatusPartial(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg := &model.OTAPackage{PackageName: "fw-3", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))
	_, err := svc.BatchUpgrade(ctx, pkg.UUID, keys)
	require.NoError(t, err)

	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D1", "success"))
	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D2", "failed"))
	got, _ := svc.Get(ctx, pkg.UUID)
	require.Equal(t, "partial", got.Status)
}

func TestOTAService_CancelBatch(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)
	pkg := &model.OTAPackage{PackageName: "fw-cancel", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))
	batch, err := svc.BatchUpgrade(ctx, pkg.UUID, keys)
	require.NoError(t, err)
	require.NoError(t, svc.CancelBatch(ctx, pkg.UUID, batch.BatchID))
	var count int64
	require.NoError(t, db.Model(&model.DeviceUpgradeStatus{}).Where("upgrade_batch_id = ? AND status = ?", batch.BatchID, "cancelled").Count(&count).Error)
	require.EqualValues(t, len(keys), count)
}
