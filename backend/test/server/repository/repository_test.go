package repository

import (
	"context"
	"testing"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupDeviceRepository(t *testing.T) (*repository.DeviceRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	deviceRepo := repository.NewDeviceRepository(db, nil)
	return deviceRepo, mock
}

func TestDeviceRepository_DB(t *testing.T) {
	deviceRepo, _ := setupDeviceRepository(t)
	ctx := context.Background()

	db := deviceRepo.DB(ctx)
	assert.NotNil(t, db)
}

func TestDeviceRepository_Create(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `devices`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `device_states`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `device_credentials`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	device := &model.Device{DeviceKey: "D002", Name: "New Device", ProductID: 1, OrganizationID: "org-1"}
	credential := &model.DeviceCredential{DeviceUUID: "device-uuid", CredentialType: "mqtt", Username: "mqtt_device-uuid", Password: "hash", Salt: "salt", Enabled: true}
	err := deviceRepo.Create(ctx, device, credential)
	assert.NoError(t, err)
}

func TestDeviceRepository_Create_RollsBackWhenCredentialInsertFails(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `devices`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `device_states`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `device_credentials`").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := deviceRepo.Create(ctx,
		&model.Device{DeviceKey: "D003", Name: "Rollback Device", ProductID: 1, OrganizationID: "org-1"},
		&model.DeviceCredential{DeviceUUID: "device-uuid-rollback", CredentialType: "mqtt", Username: "mqtt_device-uuid-rollback", Password: "hash", Salt: "salt", Enabled: true},
	)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeviceRepository_Delete(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_credentials` SET `enabled`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE `devices` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	device := &model.Device{ID: 1, DeviceUUID: "device-uuid"}
	err := deviceRepo.Delete(ctx, device)
	assert.NoError(t, err)
}

func TestDeviceRepository_Delete_RollsBackWhenCredentialDisableFails(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_credentials` SET `enabled`").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := deviceRepo.Delete(ctx, &model.Device{ID: 1, DeviceUUID: "device-uuid"})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Skip TestDeviceRepository_SaveEnabled - GORM Save uses upsert which is hard to mock precisely with sqlmock

func setupOTARespository(t *testing.T) (*repository.OTARepository, sqlmock.Sqlmock) {
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

func TestOTARepository_List(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	// Count query
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	// Find query
	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	packages, total, err := otaRepo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, packages, 1)
}

func TestOTARepository_Find(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	pkg, err := otaRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "firmware-1", pkg.PackageName)
}

func TestOTARepository_Create(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `ota_packages`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	pkg := &model.OTAPackage{PackageName: "firmware-1", Version: "1.0.0", OrganizationID: "org-1"}
	err := otaRepo.Create(ctx, pkg)
	assert.NoError(t, err)
}

func TestOTARepository_Delete(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	// Find is called first inside Delete
	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `ota_packages` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := otaRepo.Delete(ctx, 1)
	assert.NoError(t, err)
}
