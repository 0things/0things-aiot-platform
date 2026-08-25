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

func newOTATestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(
		&model.OTAPackage{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{},
		&model.Product{}, &model.Device{}, &model.DeviceState{},
	))
	return db
}

func newOTAServiceForTest(t *testing.T) (*OTAService, *gorm.DB) {
	db := newOTATestDB(t)
	iotDB := &repository.IoTDB{DB: db}
	repo := repository.NewOTARepository(iotDB)
	productRepo := repository.NewProductRepository(iotDB)
	deviceRepo := repository.NewDeviceRepository(iotDB, &repository.IoTRedis{})
	return NewOTAService(repo, productRepo, deviceRepo), db
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

func TestOTAService_DeployAndDispatch(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))

	// Deploy creates pending records and marks package deploying.
	count, err := svc.Deploy(ctx, pkg.UUID, keys)
	require.NoError(t, err)
	require.Equal(t, 2, count)
	got, _ := svc.Get(ctx, pkg.UUID)
	require.Equal(t, "deploying", got.Status)

	// PendingPackageIDs should include this package.
	pending, err := svc.repo.PendingPackageIDs(ctx)
	require.NoError(t, err)
	require.Contains(t, pending, pkg.ID)

	// Dispatch moves pending -> in_progress.
	affected, err := svc.Dispatch(ctx, pkg.UUID)
	require.NoError(t, err)
	require.EqualValues(t, 2, affected)

	stats, err := svc.repo.Statistics(ctx, pkg.ID)
	require.NoError(t, err)
	require.EqualValues(t, 0, stats.Pending)
	require.EqualValues(t, 2, stats.InProgress)
}

func TestOTAService_ReportStatusAggregation(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg := &model.OTAPackage{PackageName: "fw-2", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg, "P001"))
	_, err := svc.Deploy(ctx, pkg.UUID, keys)
	require.NoError(t, err)
	_, err = svc.Dispatch(ctx, pkg.UUID)
	require.NoError(t, err)

	// One success, one still in progress -> package stays deploying.
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
	_, err := svc.Deploy(ctx, pkg.UUID, keys)
	require.NoError(t, err)
	_, err = svc.Dispatch(ctx, pkg.UUID)
	require.NoError(t, err)

	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D1", "success"))
	require.NoError(t, svc.ReportStatus(ctx, pkg.UUID, "D2", "failed"))
	got, _ := svc.Get(ctx, pkg.UUID)
	require.Equal(t, "partial", got.Status)
}

func TestOTAService_DispatchAll(t *testing.T) {
	ctx := context.Background()
	svc, db := newOTAServiceForTest(t)
	_, keys := seedOTAProductAndDevices(t, db)

	pkg1 := &model.OTAPackage{PackageName: "fw-a", Version: "1.0.0", ProductID: 1, Status: "draft"}
	pkg2 := &model.OTAPackage{PackageName: "fw-b", Version: "1.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, svc.Create(ctx, pkg1, "P001"))
	require.NoError(t, svc.Create(ctx, pkg2, "P001"))
	_, err := svc.Deploy(ctx, pkg1.UUID, keys)
	require.NoError(t, err)
	_, err = svc.Deploy(ctx, pkg2.UUID, keys)
	require.NoError(t, err)

	total, err := svc.DispatchAll(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 4, total)

	stats1, err := svc.repo.Statistics(ctx, pkg1.ID)
	require.NoError(t, err)
	require.EqualValues(t, 2, stats1.InProgress)
	stats2, err := svc.repo.Statistics(ctx, pkg2.ID)
	require.NoError(t, err)
	require.EqualValues(t, 2, stats2.InProgress)
}
