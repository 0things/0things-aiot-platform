package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"aiot-backend/internal/enum"
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
		Joins("JOIN products p ON p.id = ota_packages.product_id AND p.organization_id = ?", tenant.GetOrganizationID(ctx))
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

func (r *OTARepository) FindByUUID(ctx context.Context, uuid string) (*model.OTAPackage, error) {
	var pkg model.OTAPackage
	if err := r.selectWithProduct(ctx, r.DB(ctx).Model(&model.OTAPackage{})).Where("ota_packages.uuid = ?", uuid).First(&pkg).Error; err != nil {
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

func (r *OTARepository) Statistics(ctx context.Context, packageID int64, batchID ...string) (UpgradeStatistics, error) {
	baseQuery := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).Where("ota_package_id = ?", strconv.FormatInt(packageID, 10))
	if len(batchID) > 0 && batchID[0] != "" {
		// 只有传入批次号时才限定范围，兼容旧的全包统计调用。
		baseQuery = baseQuery.Where("upgrade_batch_id = ?", batchID[0])
	}
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return UpgradeStatistics{}, err
	}
	counts := UpgradeStatistics{Total: total}
	countStatus := func(status string) (int64, error) {
		var count int64
		query := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
			Where("ota_package_id = ?", strconv.FormatInt(packageID, 10))
		if status == enum.OTAStatusFailed {
			query = query.Where("status IN ?", []string{enum.OTAStatusFailed, enum.OTAStatusTimeout})
		} else if status == enum.OTAStatusPending {
			query = query.Where("status IN ?", []string{enum.OTAStatusPending, enum.OTAStatusSent})
		} else {
			query = query.Where("status = ?", status)
		}
		if len(batchID) > 0 && batchID[0] != "" {
			query = query.Where("upgrade_batch_id = ?", batchID[0])
		}
		if err := query.
			Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
	var err error
	if counts.Success, err = countStatus(enum.OTAStatusSuccess); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Failed, err = countStatus(enum.OTAStatusFailed); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Cancelled, err = countStatus(enum.OTAStatusCancelled); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.Pending, err = countStatus(enum.OTAStatusPending); err != nil {
		return UpgradeStatistics{}, err
	}
	if counts.InProgress, err = countStatus(enum.OTAStatusInProgress); err != nil {
		return UpgradeStatistics{}, err
	}
	return counts, nil
}

func (r *OTARepository) Batches(ctx context.Context, packageID int64) ([]model.UpgradeBatch, error) {
	var batches []model.UpgradeBatch
	err := r.DB(ctx).Where("ota_package_id = ?", strconv.FormatInt(packageID, 10)).Find(&batches).Error
	return batches, err
}

// FindBatch 按 OTA 包和批次号查询批次，确保批次归属于当前升级包。
func (r *OTARepository) FindBatch(ctx context.Context, packageID int64, batchID string) (*model.UpgradeBatch, error) {
	var batch model.UpgradeBatch
	err := r.DB(ctx).Where("batch_id = ? AND ota_package_id = ?", batchID, strconv.FormatInt(packageID, 10)).First(&batch).Error
	return &batch, err
}

// FindBatchByID 按批次号查询批次，供设备回报场景使用。
func (r *OTARepository) FindBatchByID(ctx context.Context, batchID string) (*model.UpgradeBatch, error) {
	var batch model.UpgradeBatch
	err := r.DB(ctx).Where("batch_id = ?", batchID).First(&batch).Error
	return &batch, err
}

// CountRetryLimited 统计已达到最大重试次数的失败任务。
func (r *OTARepository) CountRetryLimited(ctx context.Context, packageID int64, batchID string) (int64, error) {
	var count int64
	err := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("ota_package_id = ? AND upgrade_batch_id = ? AND status IN ? AND max_retries > 0 AND dispatch_attempts >= max_retries", strconv.FormatInt(packageID, 10), batchID, []string{enum.OTAStatusFailed, enum.OTAStatusTimeout}).
		Count(&count).Error
	return count, err
}

// UpdateBatchStatus 更新批次状态。
func (r *OTARepository) UpdateBatchStatus(ctx context.Context, batchID, status string) error {
	return r.DB(ctx).Model(&model.UpgradeBatch{}).Where("batch_id = ?", batchID).Update("status", status).Error
}

// UpdateBatchDevicesStatus 批量更新指定批次中处于给定状态的设备任务。
func (r *OTARepository) UpdateBatchDevicesStatus(ctx context.Context, packageID int64, batchID string, from []string, status string) error {
	return r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("ota_package_id = ? AND upgrade_batch_id = ? AND status IN ?", strconv.FormatInt(packageID, 10), batchID, from).
		Updates(map[string]any{"status": status, "last_status_change_ts": time.Now().Unix()}).Error
}

// FindBatchDevice 查询批次中的设备任务记录。
func (r *OTARepository) FindBatchDevice(ctx context.Context, batchID string, deviceID int64) (*model.DeviceUpgradeStatus, error) {
	var task model.DeviceUpgradeStatus
	err := r.DB(ctx).Where("upgrade_batch_id = ? AND device_id = ?", batchID, deviceID).First(&task).Error
	return &task, err
}

// ShouldDispatchBatchDevice 判断待下发任务是否仍处于 pending，避免 Kafka 重复消息重复下发。
func (r *OTARepository) ShouldDispatchBatchDevice(ctx context.Context, batchID, deviceKey string) (bool, error) {
	var count int64
	err := r.DB(ctx).Table("ota_device_upgrade_status dus").
		Joins("JOIN devices d ON d.id = dus.device_id").
		Where("dus.upgrade_batch_id = ? AND d.device_key = ? AND dus.status = ?", batchID, deviceKey, enum.OTAStatusPending).
		Count(&count).Error
	return count > 0, err
}

// CreateBatch 持久化一个新的升级批次。
func (r *OTARepository) CreateBatch(ctx context.Context, batch *model.UpgradeBatch) error {
	return r.DB(ctx).Create(batch).Error
}

func (r *OTARepository) CreateBatchWithDeployments(ctx context.Context, batch *model.UpgradeBatch, packageID int64, deviceIDs []int64, targetVersion string) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		seen := make(map[int64]struct{}, len(deviceIDs))
		rows := make([]model.DeviceUpgradeStatus, 0, len(deviceIDs))
		for _, deviceID := range deviceIDs {
			if _, ok := seen[deviceID]; ok {
				continue
			}
			seen[deviceID] = struct{}{}
			now := time.Now().Unix()
			rows = append(rows, model.DeviceUpgradeStatus{
				DeviceID: deviceID, OTAPackageID: strconv.FormatInt(packageID, 10), UpgradeBatchID: batch.BatchID,
				Status: enum.OTAStatusPending, TargetVersion: targetVersion, TimeoutSeconds: 1800, MaxRetries: 3, LastStatusChangeTime: &now,
			})
		}
		if len(rows) > 0 {
			if err := tx.Create(&rows).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateBatchDeployments 为静态升级批次创建设备升级记录。每个设备在每个批次
// 中都有独立记录，状态初始为 pending。返回受影响的设备数量。
func (r *OTARepository) CreateBatchDeployments(ctx context.Context, packageID int64, batchID string, deviceIDs []int64) (int, error) {
	if len(deviceIDs) == 0 {
		return 0, nil
	}
	pkgID := strconv.FormatInt(packageID, 10)
	db := r.DB(ctx)

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

	toInsert := make([]model.DeviceUpgradeStatus, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		now := time.Now().Unix()
		toInsert = append(toInsert, model.DeviceUpgradeStatus{
			DeviceID:             deviceID,
			OTAPackageID:         pkgID,
			UpgradeBatchID:       batchID,
			Status:               enum.OTAStatusPending,
			TargetVersion:        "",
			LastStatusChangeTime: &now,
		})
	}

	if err := db.Create(&toInsert).Error; err != nil {
		return 0, err
	}
	return len(deviceIDs), nil
}

func (r *OTARepository) MarkDispatchResult(ctx context.Context, packageID, deviceID int64, batchID, status, dispatchError string) error {
	updates := map[string]interface{}{
		"status":                status,
		"dispatch_attempts":     gorm.Expr("dispatch_attempts + 1"),
		"last_dispatch_error":   dispatchError,
		"last_status_change_ts": time.Now().Unix(),
	}
	return r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("ota_package_id = ? AND upgrade_batch_id = ? AND device_id = ? AND status = ?", strconv.FormatInt(packageID, 10), batchID, deviceID, enum.OTAStatusPending).
		Updates(updates).Error
}

func (r *OTARepository) Deployments(ctx context.Context, packageID int64, page, size int, status string, batchID ...string) ([]model.DeviceDeployment, int64, error) {
	query := r.DB(ctx).Table("ota_device_upgrade_status dus").
		Select("dus.device_id, d.device_key, d.name as device_name, d.product_id, p.product_key, dus.current_version, dus.target_version, dus.progress, dus.upgrade_batch_id, dus.status, dus.last_status_change_ts, dus.created_at").
		Joins("JOIN devices d ON d.id = dus.device_id").
		Joins("JOIN products p ON p.id = d.product_id").
		Where("dus.ota_package_id = ?", strconv.FormatInt(packageID, 10))
	if status != "" {
		query = query.Where("dus.status = ?", status)
	}
	if len(batchID) > 0 && batchID[0] != "" {
		query = query.Where("dus.upgrade_batch_id = ?", batchID[0])
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var deployments []model.DeviceDeployment
	if err := query.Order("dus.created_at DESC, dus.id DESC").
		Offset((page - 1) * size).Limit(size).Find(&deployments).Error; err != nil {
		return nil, 0, err
	}
	return deployments, total, nil
}

// UpdateDeviceStatus 更新单台设备的升级状态（及可选的当前版本）。
func (r *OTARepository) UpdateDeviceStatus(ctx context.Context, packageID, deviceID int64, status, currentVersion string) error {
	var deployment model.DeviceUpgradeStatus
	if err := r.DB(ctx).
		Where("ota_package_id = ? AND device_id = ?", strconv.FormatInt(packageID, 10), deviceID).
		Order("id DESC").First(&deployment).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"status":                status,
		"last_status_change_ts": time.Now().Unix(),
	}
	if currentVersion != "" {
		updates["current_version"] = currentVersion
	}
	return r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("id = ?", deployment.ID).
		Updates(updates).Error
}

func (r *OTARepository) UpdateBatchDeviceStatus(ctx context.Context, batchID string, deviceID int64, status, currentVersion string, progress int32) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status":                status,
		"progress":              progress,
		"last_status_change_ts": now,
		"last_report_at":        now,
		"first_progress_at":     gorm.Expr("COALESCE(first_progress_at, ?)", now),
	}
	if currentVersion != "" {
		updates["current_version"] = currentVersion
	}
	return r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("upgrade_batch_id = ? AND device_id = ?", batchID, deviceID).
		Updates(updates).Error
}

// PendingPackageIDs 返回当前仍有 pending 设备升级记录的升级包 ID 列表。
// 保留给直接操作仓储的定时任务使用；业务层统一以 UUID 进行筛选。
func (r *OTARepository) PendingPackageIDs(ctx context.Context) ([]int64, error) {
	var ids []string
	if err := r.DB(ctx).Model(&model.DeviceUpgradeStatus{}).
		Where("status IN ?", []string{enum.OTAStatusPending, enum.OTAStatusSent}).
		Distinct("ota_package_id").
		Pluck("ota_package_id", &ids).Error; err != nil {
		return nil, err
	}
	result := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, id)
	}
	return result, nil
}
