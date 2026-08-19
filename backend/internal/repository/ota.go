package repository

import (
	"context"
	"errors"
	"strconv"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"gorm.io/gorm"
)

type OTARepository struct{ db *IoTDB }

type UpgradeStatistics struct {
	Total      int64
	Success    int64
	Failed     int64
	Cancelled  int64
	Pending    int64
	InProgress int64
}

func NewOTARepository(db *IoTDB) *OTARepository          { return &OTARepository{db: db} }
func (r *OTARepository) DB(ctx context.Context) *gorm.DB { return r.db.WithContext(ctx) }

func (r *OTARepository) selectWithProduct(ctx context.Context, query *gorm.DB) *gorm.DB {
	return query.
		Select("ota_packages.*, p.product_key AS product_key, p.name AS product_name").
		Joins("JOIN products p ON p.id = ota_packages.product_id AND p.tenant_id = ?", tenant.GetTenantID(ctx))
}

func (r *OTARepository) Find(ctx context.Context, id int64) (*model.OTAPackage, error) {
	var pkg model.OTAPackage
	if err := r.selectWithProduct(ctx, r.DB(ctx).Model(&model.OTAPackage{})).First(&pkg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pkg, nil
}
func (r *OTARepository) FindByName(ctx context.Context, name string) (*model.OTAPackage, error) {
	var pkg model.OTAPackage
	if err := r.selectWithProduct(ctx, r.DB(ctx).Model(&model.OTAPackage{})).Where("ota_packages.package_name = ?", name).First(&pkg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *OTARepository) List(ctx context.Context, page, size int) ([]model.OTAPackage, int64, error) {
	base := r.selectWithProduct(ctx, r.DB(ctx).Model(&model.OTAPackage{}))
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := r.selectWithProduct(ctx, r.DB(ctx).Model(&model.OTAPackage{}))
	var packages []model.OTAPackage
	if err := query.Order("ota_packages.created_at DESC").Offset((page - 1) * size).Limit(size).Find(&packages).Error; err != nil {
		return nil, 0, err
	}
	return packages, total, nil
}

func (r *OTARepository) Create(ctx context.Context, pkg *model.OTAPackage) error {
	return r.DB(ctx).Create(pkg).Error
}

func (r *OTARepository) Save(ctx context.Context, pkg *model.OTAPackage) error {
	return r.DB(ctx).Save(pkg).Error
}

func (r *OTARepository) Delete(ctx context.Context, id int64) error {
	pkg, err := r.Find(ctx, id)
	if err != nil {
		return err
	}
	return r.DB(ctx).Delete(pkg).Error
}

func (r *OTARepository) Statistics(ctx context.Context, packageID int64) (UpgradeStatistics, error) {
	baseQuery := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).Where("ota_package_id = ?", strconv.FormatInt(packageID, 10))
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return UpgradeStatistics{}, err
	}
	counts := UpgradeStatistics{Total: total}
	countStatus := func(status string) (int64, error) {
		var count int64
		if err := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
			Where("ota_package_id = ? AND status = ?", strconv.FormatInt(packageID, 10), status).
			Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
	var err error
	if counts.Success, err = countStatus("success"); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Failed, err = countStatus("failed"); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Cancelled, err = countStatus("cancelled"); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Pending, err = countStatus("pending"); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.InProgress, err = countStatus("in_progress"); err != nil {
		return UpgradeStatistics{}, err
	}
	return counts, nil
}

func (r *OTARepository) Batches(ctx context.Context, packageID int64) ([]model.UpgradeBatch, error) {
	var batches []model.UpgradeBatch
	err := r.DB(ctx).Where("ota_package_id = ?", strconv.FormatInt(packageID, 10)).Find(&batches).Error
	return batches, err
}

func (r *OTARepository) Deployments(ctx context.Context, packageID int64, page, size int, status string) ([]model.DeviceDeployment, int64, error) {
	query := r.DB(ctx).Table("device_upgrade_status dus").
		Select("dus.device_id, d.device_key, d.name as device_name, d.product_id, p.product_key, dus.current_version, dus.upgrade_batch_id, dus.status, dus.last_status_change_time, dus.created_at").
		Joins("JOIN devices d ON d.id = dus.device_id").
		Joins("JOIN products p ON p.id = d.product_id").
		Where("dus.ota_package_id = ?", strconv.FormatInt(packageID, 10))
	if status != "" {
		query = query.Where("dus.status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var deployments []model.DeviceDeployment
	if err := query.Offset((page - 1) * size).Limit(size).Find(&deployments).Error; err != nil {
		return nil, 0, err
	}
	return deployments, total, nil
}

// CreateDeployments 为给定的 OTA 升级包在每台目标设备上记录一条待升级记录。
// 若同一设备+升级包已存在记录，则更新为 "pending" 而非重复插入。
// 返回受影响的设备数量。
func (r *OTARepository) CreateDeployments(ctx context.Context, packageID int64, deviceIDs []int64) (int, error) {
	if len(deviceIDs) == 0 {
		return 0, nil
	}
	pkgID := strconv.FormatInt(packageID, 10)
	db := r.DB(ctx)

	// 去重，避免同一设备重复处理。
	seen := make(map[int64]struct{}, len(deviceIDs))
	uniq := make([]int64, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	deviceIDs = uniq

	// 一次性查出本包下这批设备已有的升级记录。
	var existing []model.DeviceUpgradeStatus
	if err := db.Where("ota_package_id = ? AND device_id IN ?", pkgID, deviceIDs).
		Find(&existing).Error; err != nil {
		return 0, err
	}
	existingStatus := make(map[int64]string, len(existing))
	for _, e := range existing {
		existingStatus[e.DeviceID] = e.Status
	}

	toInsert := make([]model.DeviceUpgradeStatus, 0, len(deviceIDs))
	toUpdate := make([]int64, 0)
	for _, deviceID := range deviceIDs {
		if status, ok := existingStatus[deviceID]; ok {
			if status != "pending" {
				toUpdate = append(toUpdate, deviceID)
			}
		} else {
			toInsert = append(toInsert, model.DeviceUpgradeStatus{
				DeviceID:     deviceID,
				OTAPackageID: pkgID,
				Status:       "pending",
			})
		}
	}

	if len(toUpdate) > 0 {
		if err := db.Model(&model.DeviceUpgradeStatus{}).
			Where("ota_package_id = ? AND device_id IN ?", pkgID, toUpdate).
			Update("status", "pending").Error; err != nil {
			return 0, err
		}
	}
	if len(toInsert) > 0 {
		if err := db.Create(&toInsert).Error; err != nil {
			return 0, err
		}
	}
	return len(deviceIDs), nil
}
