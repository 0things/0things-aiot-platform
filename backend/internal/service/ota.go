package service

import (
	"context"
	"errors"
	"strconv"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	gouuid "github.com/google/uuid"
)

type OTAServiceInterface interface {
	List(ctx context.Context, page, size int) ([]model.OTAPackage, int64, error)
	Get(ctx context.Context, uuid string) (*model.OTAPackage, error)
	Create(ctx context.Context, pkg *model.OTAPackage, productKey string) error
	Update(ctx context.Context, pkg *model.OTAPackage) error
	Delete(ctx context.Context, uuid string) error
	Deploy(ctx context.Context, uuid string, deviceKeys []string) (int, error)
	BatchUpgrade(ctx context.Context, uuid string, deviceKeys []string) (*model.UpgradeBatch, error)
	Dispatch(ctx context.Context, uuid string) (int64, error)
	DispatchAll(ctx context.Context) (int64, error)
	ReportStatus(ctx context.Context, uuid string, deviceKey string, status string) error
	Statistics(ctx context.Context, uuid string) (UpgradeStatistics, error)
	Batches(ctx context.Context, uuid string) ([]model.UpgradeBatch, error)
	Deployments(ctx context.Context, uuid string, page, size int, status string) ([]model.DeviceDeployment, int64, error)
}

type OTAService struct {
	repo        *repository.OTARepository
	productRepo *repository.ProductRepository
	deviceRepo  *repository.DeviceRepository
}

type UpgradeStatistics struct {
	PackageID          string
	TotalTargetDevices int64
	SuccessfulUpgrades int64
	FailedUpgrades     int64
	CancelledUpgrades  int64
	PendingUpgrades    int64
	InProgressUpgrades int64
}

func NewOTAService(repo *repository.OTARepository, productRepo *repository.ProductRepository, deviceRepo *repository.DeviceRepository) *OTAService {
	return &OTAService{repo: repo, productRepo: productRepo, deviceRepo: deviceRepo}
}

func (s *OTAService) List(ctx context.Context, page, size int) ([]model.OTAPackage, int64, error) {
	return s.repo.List(ctx, page, size)
}

func (s *OTAService) Get(ctx context.Context, uuid string) (*model.OTAPackage, error) {
	return s.repo.FindByUUID(ctx, uuid)
}

func (s *OTAService) Create(ctx context.Context, pkg *model.OTAPackage, productKey string) error {
	product, err := s.productRepo.FindByKey(ctx, productKey)
	if err != nil {
		return err
	}
	pkg.ProductID = product.ID
	pkg.ProductKey = product.ProductKey
	pkg.ProductName = product.Name
	if pkg.UUID == "" {
		pkg.UUID = gouuid.NewString()
	}
	return s.repo.Create(ctx, pkg)
}

func (s *OTAService) Update(ctx context.Context, pkg *model.OTAPackage) error {
	if _, err := s.productRepo.Find(ctx, pkg.ProductID); err != nil {
		return err
	}
	return s.repo.Save(ctx, pkg)
}

func (s *OTAService) Delete(ctx context.Context, uuid string) error {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, pkg.ID)
}

func (s *OTAService) Deploy(ctx context.Context, uuid string, deviceKeys []string) (int, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return 0, err
	}
	devices, err := s.deviceRepo.FindByKeys(ctx, deviceKeys)
	if err != nil {
		return 0, err
	}
	deviceIDs := make([]int64, 0, len(devices))
	for _, d := range devices {
		deviceIDs = append(deviceIDs, d.ID)
	}
	count, err := s.repo.CreateDeployments(ctx, pkg.ID, deviceIDs)
	if err != nil {
		return 0, err
	}
	pkg, err = s.repo.Find(ctx, pkg.ID)
	if err != nil {
		return count, err
	}
	pkg.Status = "deploying"
	if err := s.repo.Save(ctx, pkg); err != nil {
		return count, err
	}
	return count, nil
}

// BatchUpgrade 为指定升级包创建静态升级批次，并将所选设备加入该批次的升级
// （每条设备记录状态为 pending，并关联 upgrade_batch_id）。升级包状态置为
// deploying。返回创建好的批次。
func (s *OTAService) BatchUpgrade(ctx context.Context, uuid string, deviceKeys []string) (*model.UpgradeBatch, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	devices, err := s.deviceRepo.FindByKeys(ctx, deviceKeys)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, errors.New("no valid devices found for the given device keys")
	}
	deviceIDs := make([]int64, 0, len(devices))
	for _, d := range devices {
		deviceIDs = append(deviceIDs, d.ID)
	}

	batchID := gouuid.NewString()
	batch := &model.UpgradeBatch{
		BatchID:           batchID,
		OTAPackageID:      strconv.FormatInt(pkg.ID, 10),
		UpgradeStrategy:   "static",
		Status:            "pending",
		TargetDeviceCount: int32(len(deviceIDs)),
	}
	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return nil, err
	}
	if _, err := s.repo.CreateBatchDeployments(ctx, pkg.ID, batchID, deviceIDs); err != nil {
		return nil, err
	}
	pkg, err = s.repo.Find(ctx, pkg.ID)
	if err != nil {
		return batch, err
	}
	pkg.Status = "deploying"
	if err := s.repo.Save(ctx, pkg); err != nil {
		return batch, err
	}
	return batch, nil
}

// Dispatch 将指定升级包下所有 pending 的设备升级记录推进为 in_progress（下发命令），
// 并确保升级包处于 deploying 状态。返回受影响的设备数量。
func (s *OTAService) Dispatch(ctx context.Context, uuid string) (int64, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return 0, err
	}
	affected, err := s.repo.DispatchPending(ctx, pkg.ID)
	if err != nil {
		return 0, err
	}
	pkg, err = s.repo.Find(ctx, pkg.ID)
	if err != nil {
		return affected, err
	}
	if pkg.Status != "deploying" {
		pkg.Status = "deploying"
		if err := s.repo.Save(ctx, pkg); err != nil {
			return affected, err
		}
	}
	return affected, nil
}

// DispatchAll 扫描所有仍有 pending 记录的升级包并批量下发。返回受影响设备总数。
func (s *OTAService) DispatchAll(ctx context.Context) (int64, error) {
	uuids, err := s.repo.PendingPackageUUIDs(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, uuid := range uuids {
		n, err := s.Dispatch(ctx, uuid)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

// ReportStatus 上报某台设备对指定升级包的升级结果（in_progress/success/failed），
// 更新设备升级记录并重新聚合升级包状态。
func (s *OTAService) ReportStatus(ctx context.Context, uuid string, deviceKey string, status string) error {
	if status != "in_progress" && status != "success" && status != "failed" {
		return errors.New("invalid upgrade status: " + status)
	}
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	device, err := s.deviceRepo.FindByKey(ctx, deviceKey)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateDeviceStatus(ctx, pkg.ID, device.ID, status, ""); err != nil {
		return err
	}
	return s.recomputePackageStatus(ctx, pkg.ID)
}

// recomputePackageStatus 根据设备升级记录聚合升级包状态：
// 仍有 pending/in_progress 时保持 deploying；全部成功为 success；
// 全部失败为 failed；部分成功部分失败为 partial。
func (s *OTAService) recomputePackageStatus(ctx context.Context, packageID int64) error {
	counts, err := s.repo.Statistics(ctx, packageID)
	if err != nil {
		return err
	}
	if counts.Pending > 0 || counts.InProgress > 0 {
		return nil
	}
	pkg, err := s.repo.Find(ctx, packageID)
	if err != nil {
		return err
	}
	switch {
	case counts.Failed == 0:
		pkg.Status = "success"
	case counts.Success == 0:
		pkg.Status = "failed"
	default:
		pkg.Status = "partial"
	}
	return s.repo.Save(ctx, pkg)
}

func (s *OTAService) Statistics(ctx context.Context, uuid string) (UpgradeStatistics, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return UpgradeStatistics{}, err
	}
	counts, err := s.repo.Statistics(ctx, pkg.ID)
	if err != nil {
		return UpgradeStatistics{}, err
	}
	return UpgradeStatistics{
		PackageID: strconv.FormatInt(pkg.ID, 10), TotalTargetDevices: counts.Total,
		SuccessfulUpgrades: counts.Success, FailedUpgrades: counts.Failed,
		CancelledUpgrades: counts.Cancelled, PendingUpgrades: counts.Pending,
		InProgressUpgrades: counts.InProgress,
	}, nil
}

func (s *OTAService) Batches(ctx context.Context, uuid string) ([]model.UpgradeBatch, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return s.repo.Batches(ctx, pkg.ID)
}

func (s *OTAService) Deployments(ctx context.Context, uuid string, page, size int, status string) ([]model.DeviceDeployment, int64, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.Deployments(ctx, pkg.ID, page, size, status)
}
