package testutil

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.ProductTSL{},
		&model.Device{},
		&model.DeviceState{},
		&model.DeviceTag{},
		&model.DeviceShadow{},
		&model.DeviceShadowHistory{},
		&model.DeviceEvent{},
		&model.DevicePushRecord{},
		&model.OTAPackage{},
		&model.UpgradeBatch{},
		&model.DeviceUpgradeStatus{},
		&model.SceneLinkage{},
		&model.SceneLinkageDetail{},
		&model.Alert{},
	)
	require.NoError(t, err)

	return db
}

func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()
	product := &model.Product{ID: 1, ProductKey: "P001", Name: "Test Product", TenantID: 1}
	db.Create(product)

	device := &model.Device{ID: 1, DeviceKey: "D001", Name: "Test Device", ProductID: 1, TenantID: 1, Enabled: true}
	db.Create(device)

	db.Create(&model.DeviceState{ID: 1, DeviceKey: "D001", State: "online"})
}

func NewTestRepositories(db *gorm.DB) (*repository.DeviceRepository, *repository.ProductRepository, *repository.DeviceTagRepository, *repository.DeviceShadowRepository, *repository.PushRecordRepository) {
	iotDB := &repository.IoTDB{DB: db}
	iotRedis := &repository.IoTRedis{Client: nil}

	deviceRepo := repository.NewDeviceRepository(iotDB, iotRedis)
	productRepo := repository.NewProductRepository(iotDB)
	tagRepo := repository.NewDeviceTagRepository(iotDB)
	shadowRepo := repository.NewDeviceShadowRepository(iotDB)
	pushRepo := repository.NewPushRecordRepository(iotDB)

	return deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo
}

func NewTestDeviceService(db *gorm.DB) *service.DeviceService {
	deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo := NewTestRepositories(db)
	return service.NewDeviceService(deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo)
}

func NewTestProductService(db *gorm.DB) *service.ProductService {
	productRepo := repository.NewProductRepository(&repository.IoTDB{DB: db})
	return service.NewProductService(productRepo)
}

func NewTestAlertService(db *gorm.DB) *service.AlertService {
	alertRepo := repository.NewAlertRepository(&repository.IoTDB{DB: db})
	return service.NewAlertService(alertRepo)
}

func NewTestOTATotalService(db *gorm.DB) *service.OTAService {
	otaRepo := repository.NewOTARepository(&repository.IoTDB{DB: db})
	productRepo := repository.NewProductRepository(&repository.IoTDB{DB: db})
	return service.NewOTAService(otaRepo, productRepo)
}

func NewTestSceneLinkageService(db *gorm.DB) *service.SceneLinkageService {
	sceneRepo := repository.NewSceneLinkageRepository(&repository.IoTDB{DB: db})
	return service.NewSceneLinkageService(sceneRepo)
}

func NewTestDeviceEventService(db *gorm.DB) *service.DeviceEventService {
	eventRepo := repository.NewDeviceEventRepository(&repository.IoTDB{DB: db})
	deviceRepo := repository.NewDeviceRepository(&repository.IoTDB{DB: db}, &repository.IoTRedis{Client: nil})
	return service.NewDeviceEventService(eventRepo, deviceRepo)
}

func NewTestProductTSLService(db *gorm.DB) *service.ProductTSLService {
	productRepo := repository.NewProductRepository(&repository.IoTDB{DB: db})
	tslRepo := repository.NewProductTSLRepository(&repository.IoTDB{DB: db})
	return service.NewProductTSLService(productRepo, tslRepo)
}

func ContextWithTenant(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, "tenant_id", tenantID)
}
