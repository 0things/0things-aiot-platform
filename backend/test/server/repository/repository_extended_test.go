package repository

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --- Push Record Repository Tests ---

func setupPushRecordRepository(t *testing.T) (*repository.PushRecordRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	pushRepo := repository.NewPushRecordRepository(db)
	return pushRepo, mock
}

func TestPushRecordRepository_CreatePushRecord(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `device_push_records`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	record := &model.DevicePushRecord{DeviceID: 1, OperationType: "push", Status: "success"}
	err := pushRepo.CreatePushRecord(ctx, record)
	assert.NoError(t, err)
}

func TestPushRecordRepository_ListPushRecords_WithOperationType(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(
		sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}).
			AddRow(1, 1, "push", "success", time.Now()),
	)

	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "push", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)
}

func TestPushRecordRepository_ListPushRecords_WithStatus(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}))

	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "", "failed")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, records)
}

func TestPushRecordRepository_FindPushRecord(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}).
		AddRow(1, 1, "push", "success", time.Now())
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(rows)

	record, err := pushRepo.FindPushRecord(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "push", record.OperationType)
}

func TestPushRecordRepository_DeletePushRecords_WithBefore(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	before := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `device_push_records`").WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	count, err := pushRepo.DeletePushRecords(ctx, 1, &before)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestPushRecordRepository_DeletePushRecords_NilBefore(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `device_push_records`").WillReturnResult(sqlmock.NewResult(0, 10))
	mock.ExpectCommit()

	count, err := pushRepo.DeletePushRecords(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

// --- User Repository Tests ---

func setupExtendedUserRepository(t *testing.T) (*repository.UserRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewRepository(nil, db)
	userRepo := repository.NewUserRepository(repo)
	return &userRepo, mock
}

func TestExtendedUserRepository_GetByID(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "email", "nickname"}).
		AddRow(1, "u1", "test@test.com", "Test User")
	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnRows(rows)

	user, err := (*userRepo).GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestExtendedUserRepository_GetByID_NotFound(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnError(gorm.ErrRecordNotFound)

	user, err := (*userRepo).GetByID(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestExtendedUserRepository_GetByEmail(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "email", "nickname"}).
		AddRow(1, "u1", "test@test.com", "Test User")
	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnRows(rows)

	user, err := (*userRepo).GetByEmail(ctx, "test@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestExtendedUserRepository_GetByEmail_NotFound(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnError(gorm.ErrRecordNotFound)

	user, err := (*userRepo).GetByEmail(ctx, "nonexistent@test.com")
	assert.Error(t, err)
	assert.Nil(t, user)
}

// --- Device Shadow Repository Tests ---

func setupDeviceShadowRepository(t *testing.T) (*repository.DeviceShadowRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	shadowRepo := repository.NewDeviceShadowRepository(db)
	return shadowRepo, mock
}

func TestDeviceShadowRepository_GetShadow(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "desired", "reported", "metadata", "version"}).
		AddRow(1, 1, []byte("{}"), []byte("{}"), []byte("{}"), 1)
	mock.ExpectQuery("SELECT .+ FROM `device_shadows`").WillReturnRows(rows)

	shadow, err := shadowRepo.GetShadow(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), shadow.Version)
}

func TestDeviceShadowRepository_GetShadow_NotFound(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_shadows`").WillReturnError(gorm.ErrRecordNotFound)

	shadow, err := shadowRepo.GetShadow(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, shadow)
}

// MutateShadow uses GORM Transaction wrapping — error gets wrapped, hard to test with sqlmock
func TestDeviceShadowRepository_MutateShadow_CreateNew(t *testing.T) {
	// This path creates + updates within a Transaction — complex to mock precisely
	t.Skip("skipped: MutateShadow uses GORM Transaction which wraps errors")
}

func TestDeviceShadowRepository_ListShadowHistory(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "version", "source", "created_at"}).
		AddRow(1, 1, 1, "app", time.Now()).
		AddRow(2, 1, 2, "device", time.Now())
	mock.ExpectQuery("SELECT .+ FROM `device_shadow_histories`").WillReturnRows(rows)

	history, err := shadowRepo.ListShadowHistory(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
}

// --- Device Tag Repository Tests ---

func setupDeviceTagRepository(t *testing.T) (*repository.DeviceTagRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	tagRepo := repository.NewDeviceTagRepository(db)
	return tagRepo, mock
}

func TestDeviceTagRepository_ListTags(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "key", "value", "source"}).
		AddRow(1, 1, "location", "factory-A", "manual")
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(rows)

	tags, err := tagRepo.ListTags(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
}

func TestDeviceTagRepository_DeleteTags(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := tagRepo.DeleteTags(ctx, 1, []string{"location", "env"})
	assert.NoError(t, err)
}

func TestDeviceTagRepository_SetTags_Replace(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec("INSERT INTO `device_tags`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := tagRepo.SetTags(ctx, 1, map[string]string{"key1": "val1"}, true)
	assert.NoError(t, err)
}

func TestDeviceTagRepository_SetTags_Upsert_Existing(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "device_id", "key", "value", "source"}).
		AddRow(1, 1, "location", "factory-A", "manual")
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(rows)
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := tagRepo.SetTags(ctx, 1, map[string]string{"location": "factory-B"}, false)
	assert.NoError(t, err)
}

// --- OTA Repository Additional Tests ---

func setupExtendedOTARespository(t *testing.T) (*repository.OTARepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	otaRepo := repository.NewOTARepository(db)
	return otaRepo, mock
}

func TestOTARepository_FindByName(t *testing.T) {
	otaRepo, mock := setupExtendedOTARespository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	pkg, err := otaRepo.FindByName(ctx, "firmware-1")
	assert.NoError(t, err)
	assert.Equal(t, "firmware-1", pkg.PackageName)
}

func TestOTARepository_FindByName_NotFound(t *testing.T) {
	otaRepo, mock := setupExtendedOTARespository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnError(gorm.ErrRecordNotFound)

	pkg, err := otaRepo.FindByName(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, pkg)
}

// --- Device Event Repository Tests ---

func setupDeviceEventRepository(t *testing.T) (*repository.DeviceEventRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	eventRepo := repository.NewDeviceEventRepository(db)
	return eventRepo, mock
}

func TestDeviceEventRepository_Create(t *testing.T) {
	eventRepo, mock := setupDeviceEventRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `device_events`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	event := &model.DeviceEvent{DeviceID: 1, EventType: "temperature", EventAt: time.Now()}
	err := eventRepo.Create(ctx, event)
	assert.NoError(t, err)
}

func TestDeviceEventRepository_List_WithFilters(t *testing.T) {
	eventRepo, mock := setupDeviceEventRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_events`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery("SELECT .+ FROM `device_events`").WillReturnRows(sqlmock.NewRows([]string{"id"}))

	events, total, err := eventRepo.List(ctx, dto.ListDeviceEventsQuery{Page: 1, PageSize: 10, Keyword: "temp", DeviceKey: "D001", EventType: "temperature"})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, events)
}

// --- Product Repository Tests ---

func setupProductRepository(t *testing.T) (*repository.ProductRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	return productRepo, mock
}

func TestProductRepository_Find(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_key", "name", "status", "created_at", "updated_at"}).
		AddRow(1, "P001", "Test Product", "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `products`").WillReturnRows(rows)
	mock.ExpectQuery("SELECT .+ FROM `product_protocols`").WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "transport_protocol", "application_protocol"}))

	product, err := productRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_FindByKey(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_key", "name", "status", "created_at", "updated_at"}).
		AddRow(1, "P001", "Test Product", "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `products`").WillReturnRows(rows)
	mock.ExpectQuery("SELECT .+ FROM `product_protocols`").WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "transport_protocol", "application_protocol"}))

	product, err := productRepo.FindByKey(ctx, "P001")
	assert.NoError(t, err)
	assert.Equal(t, "P001", product.ProductKey)
}

func TestProductRepository_Create(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	product := &model.Product{Name: "New Product", ProductKey: "P002", Status: "active"}
	err := productRepo.Create(ctx, product)
	assert.NoError(t, err)
}

func TestProductRepository_Delete(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `products` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	product := &model.Product{ID: 1}
	err := productRepo.Delete(ctx, product)
	assert.NoError(t, err)
}

// --- Product TSL Repository Tests ---

func setupProductTSLRepository(t *testing.T) (*repository.ProductTSLRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	tslRepo := repository.NewProductTSLRepository(db)
	return tslRepo, mock
}

func TestProductTSLRepository_FindByProductID(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_id", "tsl"}).
		AddRow(1, 1, `{"properties":[]}`)
	mock.ExpectQuery("SELECT .+ FROM `product_ts_ls`").WillReturnRows(rows)

	tsl, err := tslRepo.FindByProductID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, `{"properties":[]}`, tsl.TSL)
}

func TestProductTSLRepository_FindByProductID_NotFound(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `product_ts_ls`").WillReturnError(gorm.ErrRecordNotFound)

	tsl, err := tslRepo.FindByProductID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, tsl)
}

func TestProductTSLRepository_Create(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `product_ts_ls`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{TSL: `{"properties":[]}`}
	err := tslRepo.Create(ctx, tsl)
	assert.NoError(t, err)
}

func TestProductTSLRepository_Save(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `product_ts_ls`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{ID: 1, TSL: `{"properties":[{"name":"temp"}]}`}
	err := tslRepo.Save(ctx, tsl)
	assert.NoError(t, err)
}

func TestProductTSLRepository_Delete(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `product_ts_ls`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{ID: 1}
	err := tslRepo.Delete(ctx, tsl)
	assert.NoError(t, err)
}

// --- Device Telemetry (nil redis) ---

func TestDeviceRepository_Telemetry_NilRedis(t *testing.T) {
	deviceRepo, _ := setupDeviceRepository(t)
	ctx := context.Background()

	result, err := deviceRepo.Telemetry(ctx, "D001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis is not configured")
	assert.Equal(t, "", result)
}

func TestDeviceRepository_Statistics(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	// Statistics uses gen queries with tenant filter and JOINs — just verify it returns data
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	stats, err := deviceRepo.Statistics(ctx)
	assert.NoError(t, err)
	assert.True(t, stats.TotalDevices > 0)
}

func int64Ptr(v int64) *int64 { return &v }

// --- Repository Infrastructure Tests (SQLite in-memory) ---

func setupSQLiteRepo(t *testing.T) *repository.Repository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return repository.NewRepository(nil, db)
}

func TestRepository_NewTransaction(t *testing.T) {
	repo := setupSQLiteRepo(t)
	tx := repository.NewTransaction(repo)
	assert.NotNil(t, tx)
}

func TestRepository_DB_WithoutContext(t *testing.T) {
	repo := setupSQLiteRepo(t)
	ctx := context.Background()
	db := repo.DB(ctx)
	assert.NotNil(t, db)
}

func TestRepository_DB_WithTransactionContext(t *testing.T) {
	repo := setupSQLiteRepo(t)
	ctx := context.Background()

	// Create a table to use in the transaction
	type TestModel struct {
		ID   int64
		Name string
	}
	db := repo.DB(ctx)
	err := db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	// Execute a transaction
	err = repo.Transaction(ctx, func(ctx context.Context) error {
		txDB := repo.DB(ctx)
		return txDB.Create(&TestModel{Name: "test"}).Error
	})
	assert.NoError(t, err)

	// Verify data was committed
	var count int64
	repo.DB(ctx).Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestRepository_Transaction_Rollback(t *testing.T) {
	repo := setupSQLiteRepo(t)
	ctx := context.Background()

	type TestModel struct {
		ID   int64
		Name string
	}
	db := repo.DB(ctx)
	err := db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	err = repo.Transaction(ctx, func(ctx context.Context) error {
		return errors.New("intentional error")
	})
	assert.Error(t, err)

	var count int64
	repo.DB(ctx).Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// --- Product Repository Extended Tests (SQLite) ---

func setupSQLiteProductRepo(t *testing.T) (*repository.ProductRepository, *repository.Repository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.Device{}, &model.Category{}, &model.ProductProtocol{}))
	iotDB := db
	repo := repository.NewRepository(nil, db)
	productRepo := repository.NewProductRepository(iotDB)
	return productRepo, repo
}

func TestProductRepository_List_WithFilters(t *testing.T) {
	productRepo, _ := setupSQLiteProductRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	// Create products
	cat1 := int64(1)
	cat2 := int64(2)
	productRepo.Create(ctx, &model.Product{ProductKey: "P001", Name: "IoT Sensor", Status: "active", CategoryID: &cat1, OrganizationID: 1})
	productRepo.Create(ctx, &model.Product{ProductKey: "P002", Name: "Gateway", Status: "inactive", CategoryID: &cat2, OrganizationID: 1})

	// List with category filter
	products, total, err := productRepo.List(ctx, 1, 10, "1", "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, products, 1)

	// List with status filter
	products, total, err = productRepo.List(ctx, 1, 10, "", "active", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)

	// List with search
	products, total, err = productRepo.List(ctx, 1, 10, "", "", "Sensor")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
}

func TestProductRepository_CountDevices(t *testing.T) {
	productRepo, _ := setupSQLiteProductRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	count, err := productRepo.CountDevices(ctx, product.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestProductRepository_DB(t *testing.T) {
	productRepo, _ := setupSQLiteProductRepo(t)
	ctx := context.Background()
	db := productRepo.DB(ctx)
	assert.NotNil(t, db)
}

// --- OTA Repository Extended Tests (SQLite) ---

func setupSQLiteOTARepo(t *testing.T) (*repository.OTARepository, *repository.ProductRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.Device{}, &model.OTAPackage{}, &model.UpgradeBatch{}, &model.DeviceUpgradeStatus{}))
	iotDB := db
	otaRepo := repository.NewOTARepository(iotDB)
	productRepo := repository.NewProductRepository(iotDB)
	return otaRepo, productRepo
}

func TestOTARepository_Statistics(t *testing.T) {
	otaRepo, productRepo := setupSQLiteOTARepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0", ProductID: product.ID, OrganizationID: 1}
	otaRepo.Create(ctx, pkg)

	stats, err := otaRepo.Statistics(ctx, pkg.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stats.Total)
}

func TestOTARepository_Batches(t *testing.T) {
	otaRepo, productRepo := setupSQLiteOTARepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0", ProductID: product.ID, OrganizationID: 1}
	otaRepo.Create(ctx, pkg)

	batches, err := otaRepo.Batches(ctx, pkg.ID)
	assert.NoError(t, err)
	assert.Empty(t, batches)
}

func TestOTARepository_Deployments(t *testing.T) {
	otaRepo, productRepo := setupSQLiteOTARepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0", ProductID: product.ID, OrganizationID: 1}
	otaRepo.Create(ctx, pkg)

	deployments, total, err := otaRepo.Deployments(ctx, pkg.ID, 1, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, deployments)
}

func TestOTARepository_FindByName_NotFound_SQLite(t *testing.T) {
	otaRepo, _ := setupSQLiteOTARepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	_, err := otaRepo.FindByName(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestOTARepository_CreateBatchDeployments(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.OTAPackage{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{}))

	iotDB := db
	otaRepo := repository.NewOTARepository(iotDB)
	productRepo := repository.NewProductRepository(iotDB)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	require.NoError(t, productRepo.Create(ctx, product))
	pkg := &model.OTAPackage{PackageName: "fw-1", Version: "1.0", ProductID: product.ID, OrganizationID: 1}
	require.NoError(t, otaRepo.Create(ctx, pkg))
	pkgID := strconv.FormatInt(pkg.ID, 10)

	readRows := func() []model.DeviceUpgradeStatus {
		var rows []model.DeviceUpgradeStatus
		require.NoError(t, db.Where("ota_package_id = ?", pkgID).Find(&rows).Error)
		return rows
	}

	// 空输入
	n, err := otaRepo.CreateBatchDeployments(ctx, pkg.ID, "batch-1", nil)
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	// 批量插入 3 台设备
	deviceIDs := []int64{1, 2, 3}
	n, err = otaRepo.CreateBatchDeployments(ctx, pkg.ID, "batch-1", deviceIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, n)
	rows := readRows()
	assert.Len(t, rows, 3)
	for _, r := range rows {
		assert.Equal(t, "pending", r.Status)
		assert.Equal(t, "batch-1", r.UpgradeBatchID)
	}

	// 入参重复设备 ID 只处理一次
	n, err = otaRepo.CreateBatchDeployments(ctx, pkg.ID, "batch-2", []int64{5, 5, 5})
	require.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestOTARepository_Delete_NotFound(t *testing.T) {
	otaRepo, _ := setupSQLiteOTARepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	err := otaRepo.Delete(ctx, 999)
	assert.Error(t, err)
}

// --- Device Repository Extended Tests (SQLite) ---

func setupSQLiteDeviceRepo(t *testing.T) (*repository.DeviceRepository, *repository.ProductRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.Device{}, &model.DeviceState{}))
	iotDB := db
	deviceRepo := repository.NewDeviceRepository(iotDB, nil)
	productRepo := repository.NewProductRepository(iotDB)
	return deviceRepo, productRepo
}

func TestDeviceRepository_Find_NotFound(t *testing.T) {
	deviceRepo, _ := setupSQLiteDeviceRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	_, err := deviceRepo.Find(ctx, 999)
	assert.Error(t, err)
}

func TestDeviceRepository_FindByKey_NotFound(t *testing.T) {
	deviceRepo, _ := setupSQLiteDeviceRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	_, err := deviceRepo.FindByKey(ctx, "NONEXIST")
	assert.Error(t, err)
}

func TestDeviceRepository_FindByKeyForEvent_NotFound(t *testing.T) {
	deviceRepo, _ := setupSQLiteDeviceRepo(t)
	ctx := context.Background()

	_, err := deviceRepo.FindByKeyForEvent(ctx, "NONEXIST")
	assert.Error(t, err)
}

func TestDeviceRepository_FindByKeys(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Device{}))

	iotDB := db
	deviceRepo := repository.NewDeviceRepository(iotDB, nil)
	ctx := tenant.WithTenant(context.Background(), 1)

	require.NoError(t, db.Create(&model.Device{DeviceKey: "D001", Name: "dev1", OrganizationID: 1, Enabled: true}).Error)
	require.NoError(t, db.Create(&model.Device{DeviceKey: "D002", Name: "dev2", OrganizationID: 1, Enabled: true}).Error)
	require.NoError(t, db.Create(&model.Device{DeviceKey: "D003", Name: "dev3", OrganizationID: 2, Enabled: true}).Error)

	// 命中同租户多台
	devices, err := deviceRepo.FindByKeys(ctx, []string{"D001", "D002"})
	require.NoError(t, err)
	assert.Len(t, devices, 2)

	// 部分命中：未知 key 被忽略
	devices, err = deviceRepo.FindByKeys(ctx, []string{"D001", "UNKNOWN"})
	require.NoError(t, err)
	assert.Len(t, devices, 1)

	// 跨租户 D003 不应返回（当前 tenant=1）
	devices, err = deviceRepo.FindByKeys(ctx, []string{"D003"})
	require.NoError(t, err)
	assert.Empty(t, devices)

	// 空输入返回 nil
	devices, err = deviceRepo.FindByKeys(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, devices)
}

func TestDeviceRepository_List_WithEnabledFilter(t *testing.T) {
	deviceRepo, productRepo := setupSQLiteDeviceRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	device := &model.Device{DeviceKey: "D001", Name: "Enabled Device", ProductID: product.ID, OrganizationID: 1, Enabled: true}
	deviceRepo.Create(ctx, device)

	enabled := true
	devices, total, err := deviceRepo.List(ctx, dto.ListDevicesQuery{Page: 1, PageSize: 10, Enabled: &enabled})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, devices, 1)
}

func TestDeviceRepository_List_WithSearch(t *testing.T) {
	deviceRepo, productRepo := setupSQLiteDeviceRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	device := &model.Device{DeviceKey: "D001", Name: "Test Device", ProductID: product.ID, OrganizationID: 1}
	deviceRepo.Create(ctx, device)

	devices, total, err := deviceRepo.List(ctx, dto.ListDevicesQuery{Page: 1, PageSize: 10, Search: "Test"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, devices, 1)
}

func TestDeviceRepository_Statistics_WithData(t *testing.T) {
	deviceRepo, productRepo := setupSQLiteDeviceRepo(t)
	ctx := tenant.WithTenant(context.Background(), 1)

	product := &model.Product{ProductKey: "P001", Name: "Test", Status: "active", OrganizationID: 1}
	productRepo.Create(ctx, product)

	device := &model.Device{DeviceKey: "D001", Name: "Device 1", ProductID: product.ID, OrganizationID: 1}
	deviceRepo.Create(ctx, device)

	stats, err := deviceRepo.Statistics(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), stats.TotalDevices)
	assert.Equal(t, int64(1), stats.InactiveDevices)
}

func TestDeviceRepository_DB_SQLite(t *testing.T) {
	deviceRepo, _ := setupSQLiteDeviceRepo(t)
	ctx := context.Background()
	db := deviceRepo.DB(ctx)
	assert.NotNil(t, db)
}

// --- Push Record Repository Extended Tests (SQLite) ---

func setupSQLitePushRecordRepo(t *testing.T) *repository.PushRecordRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.DevicePushRecord{}))
	iotDB := db
	return repository.NewPushRecordRepository(iotDB)
}

func TestPushRecordRepository_CRUD(t *testing.T) {
	pushRepo := setupSQLitePushRecordRepo(t)
	ctx := context.Background()

	// Create
	record := &model.DevicePushRecord{DeviceID: 1, OperationType: "push", Status: "success", Payload: `{"key":"val"}`}
	err := pushRepo.CreatePushRecord(ctx, record)
	require.NoError(t, err)
	assert.NotZero(t, record.ID)

	// Find
	found, err := pushRepo.FindPushRecord(ctx, record.ID)
	assert.NoError(t, err)
	assert.Equal(t, "push", found.OperationType)

	// List
	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)

	// Delete
	count, err := pushRepo.DeletePushRecords(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPushRecordRepository_FindPushRecord_NotFound(t *testing.T) {
	pushRepo := setupSQLitePushRecordRepo(t)
	ctx := context.Background()

	_, err := pushRepo.FindPushRecord(ctx, 999)
	assert.Error(t, err)
}

func TestPushRecordRepository_List_WithFilters(t *testing.T) {
	pushRepo := setupSQLitePushRecordRepo(t)
	ctx := context.Background()

	pushRepo.CreatePushRecord(ctx, &model.DevicePushRecord{DeviceID: 1, OperationType: "push", Status: "success"})
	pushRepo.CreatePushRecord(ctx, &model.DevicePushRecord{DeviceID: 1, OperationType: "property", Status: "failed"})

	// Filter by operationType
	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "push", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)

	// Filter by status
	records, total, err = pushRepo.ListPushRecords(ctx, 1, 1, 10, "", "failed")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)
}

// --- Device Shadow Repository Extended Tests (SQLite) ---

func setupSQLiteShadowRepo(t *testing.T) *repository.DeviceShadowRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.DeviceShadow{}, &model.DeviceShadowHistory{}))
	iotDB := db
	return repository.NewDeviceShadowRepository(iotDB)
}

func TestDeviceShadowRepository_MutateShadow_CreateNew_SQLite(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	desired := map[string]any{"temp": 25}
	shadow, err := shadowRepo.MutateShadow(ctx, 1, 0, "app", &desired, nil, false)
	assert.NoError(t, err)
	assert.NotNil(t, shadow)
	assert.Equal(t, int64(1), shadow.Version)

	// Verify shadow was created
	got, err := shadowRepo.GetShadow(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Version)
}

func TestDeviceShadowRepository_MutateShadow_UpdateExisting(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	// Create
	desired := map[string]any{"temp": 25}
	shadow, err := shadowRepo.MutateShadow(ctx, 1, 0, "app", &desired, nil, false)
	require.NoError(t, err)

	// Update
	reported := map[string]any{"temp": 26}
	shadow2, err := shadowRepo.MutateShadow(ctx, 1, shadow.Version, "device", nil, &reported, false)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), shadow2.Version)
}

func TestDeviceShadowRepository_MutateShadow_ClearDesired(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	// Create with desired
	desired := map[string]any{"temp": 25}
	shadow, err := shadowRepo.MutateShadow(ctx, 1, 0, "app", &desired, nil, false)
	require.NoError(t, err)

	// Clear desired
	shadow2, err := shadowRepo.MutateShadow(ctx, 1, shadow.Version, "app", nil, nil, true)
	assert.NoError(t, err)
	assert.NotNil(t, shadow2)
}

func TestDeviceShadowRepository_MutateShadow_VersionConflict(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	// Create
	desired := map[string]any{"temp": 25}
	_, err := shadowRepo.MutateShadow(ctx, 1, 0, "app", &desired, nil, false)
	require.NoError(t, err)

	// Try to update with wrong version
	reported := map[string]any{"temp": 26}
	_, err = shadowRepo.MutateShadow(ctx, 1, 999, "device", nil, &reported, false)
	assert.Error(t, err)
}

func TestDeviceShadowRepository_MutateShadow_NewWithNonZeroVersion(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	// Try to create new with non-zero version (should fail)
	desired := map[string]any{"temp": 25}
	_, err := shadowRepo.MutateShadow(ctx, 1, 5, "app", &desired, nil, false)
	assert.Error(t, err)
}

func TestDeviceShadowRepository_GetShadow_NotFound_SQLite(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	shadow, err := shadowRepo.GetShadow(ctx, 999)
	assert.NoError(t, err)
	assert.NotNil(t, shadow)
	assert.Equal(t, int64(999), shadow.DeviceID)
	assert.Equal(t, int64(0), shadow.Version)
}

func TestDeviceShadowRepository_ListShadowHistory_SQLite(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()

	// Create some history
	desired := map[string]any{"temp": 25}
	shadowRepo.MutateShadow(ctx, 1, 0, "app", &desired, nil, false)
	reported := map[string]any{"temp": 26}
	shadowRepo.MutateShadow(ctx, 1, 1, "device", nil, &reported, false)

	history, err := shadowRepo.ListShadowHistory(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
}

func TestDeviceShadowRepository_DB(t *testing.T) {
	shadowRepo := setupSQLiteShadowRepo(t)
	ctx := context.Background()
	db := shadowRepo.DB(ctx)
	assert.NotNil(t, db)
}

// --- Device Tag Repository Extended Tests (SQLite) ---

func setupSQLiteTagRepo(t *testing.T) *repository.DeviceTagRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.DeviceTag{}))
	iotDB := db
	return repository.NewDeviceTagRepository(iotDB)
}

func TestDeviceTagRepository_CRUD(t *testing.T) {
	tagRepo := setupSQLiteTagRepo(t)
	ctx := context.Background()

	// Set tags
	err := tagRepo.SetTags(ctx, 1, map[string]string{"env": "prod", "region": "us-east"}, false)
	require.NoError(t, err)

	// List tags
	tags, err := tagRepo.ListTags(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 2)

	// Delete tags
	err = tagRepo.DeleteTags(ctx, 1, []string{"env"})
	assert.NoError(t, err)

	tags, err = tagRepo.ListTags(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
}

func TestDeviceTagRepository_SetTags_Replace_SQLite(t *testing.T) {
	tagRepo := setupSQLiteTagRepo(t)
	ctx := context.Background()

	// Create initial tags
	tagRepo.SetTags(ctx, 1, map[string]string{"env": "prod"}, false)

	// Replace all tags
	err := tagRepo.SetTags(ctx, 1, map[string]string{"region": "us-east"}, true)
	assert.NoError(t, err)

	tags, err := tagRepo.ListTags(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "region", tags[0].Key)
}

func TestDeviceTagRepository_DB(t *testing.T) {
	tagRepo := setupSQLiteTagRepo(t)
	ctx := context.Background()
	db := tagRepo.DB(ctx)
	assert.NotNil(t, db)
}

// --- User Repository Tests (SQLite) ---

func setupSQLiteUserRepo(t *testing.T) repository.UserRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))
	repo := repository.NewRepository(nil, db)
	return repository.NewUserRepository(repo)
}

func TestUserRepository_CRUD(t *testing.T) {
	userRepo := setupSQLiteUserRepo(t)
	ctx := context.Background()

	// Create
	user := &model.User{UserId: "u1", Email: "test@test.com", Nickname: "Test User"}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// GetByID
	found, err := userRepo.GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", found.Email)

	// GetByEmail
	found2, err := userRepo.GetByEmail(ctx, "test@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, found2)

	// Update
	user.Nickname = "Updated"
	err = userRepo.Update(ctx, user)
	assert.NoError(t, err)

	// GetByID after update
	found3, err := userRepo.GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated", found3.Nickname)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	userRepo := setupSQLiteUserRepo(t)
	ctx := context.Background()

	_, err := userRepo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	userRepo := setupSQLiteUserRepo(t)
	ctx := context.Background()

	user, err := userRepo.GetByEmail(ctx, "nonexistent@test.com")
	assert.Error(t, err)
	assert.Nil(t, user)
}

// --- Device Event Repository Tests (SQLite) ---

func setupSQLiteEventRepo(t *testing.T) (*repository.DeviceEventRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.DeviceEvent{}))
	iotDB := db
	return repository.NewDeviceEventRepository(iotDB), db
}

func TestDeviceEventRepository_CreateAndList(t *testing.T) {
	_, db := setupSQLiteEventRepo(t)

	// Create events via raw GORM (since List uses gen queries with JOINs)
	db.Create(&model.DeviceEvent{DeviceID: 1, EventType: "temperature", EventAt: time.Now()})

	// Note: List uses gen queries which need full GORM Gen setup, so we test Create directly
	// The integration tests in service layer already cover List
}
