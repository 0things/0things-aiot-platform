package service

import (
	"context"
	"errors"
	"strconv"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

type OTAServiceInterface interface {
	List(ctx context.Context, page, size int) ([]model.OTAPackage, int64, error)
	Get(ctx context.Context, id int64) (*model.OTAPackage, error)
	Create(ctx context.Context, pkg *model.OTAPackage, productKey string) error
	Update(ctx context.Context, pkg *model.OTAPackage) error
	Delete(ctx context.Context, id int64) error
	Deploy(ctx context.Context, packageID int64, deviceKeys []string) (int, error)
	Dispatch(ctx context.Context, packageID int64) (int64, error)
	DispatchAll(ctx context.Context) (int64, error)
	ReportStatus(ctx context.Context, packageID int64, deviceKey string, status string) error
	Statistics(ctx context.Context, packageName string) (UpgradeStatistics, error)
	Batches(ctx context.Context, packageName string) ([]model.UpgradeBatch, error)
	Deployments(ctx context.Context, packageName string, page, size int, status string) ([]model.DeviceDeployment, int64, error)
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

func (s *OTAService) Get(ctx context.Context, id int64) (*model.OTAPackage, error) {
	return s.repo.Find(ctx, id)
}

func (s *OTAService) Create(ctx context.Context, pkg *model.OTAPackage, productKey string) error {
	product, err := s.productRepo.FindByKey(ctx, productKey)
	if err != nil {
		return err
	}
	pkg.ProductID = product.ID
	pkg.ProductKey = product.ProductKey
	pkg.ProductName = product.Name
	return s.repo.Create(ctx, pkg)
}

func (s *OTAService) Update(ctx context.Context, pkg *model.OTAPackage) error {
	if _, err := s.productRepo.Find(ctx, pkg.ProductID); err != nil {
		return err
	}
	return s.repo.Save(ctx, pkg)
}

func (s *OTAService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *OTAService) Deploy(ctx context.Context, packageID int64, deviceKeys []string) (int, error) {
	if _, err := s.repo.Find(ctx, packageID); err != nil {
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
	count, err := s.repo.CreateDeployments(ctx, packageID, deviceIDs)
	if err != nil {
		return 0, err
	}
	pkg, err := s.repo.Find(ctx, packageID)
	if err != nil {
		return count, err
	}
	pkg.Status = "deploying"
	if err := s.repo.Save(ctx, pkg); err != nil {
		return count, err
	}
	return count, nil
}

// Dispatch 将指定升级包下所有 pending 的设备升级记录推进为 in_progress（下发命令），
// 并确保升级包处于 deploying 状态。返回受影响的设备数量。
func (s *OTAService) Dispatch(ctx context.Context, packageID int64) (int64, error) {
	if _, err := s.repo.Find(ctx, packageID); err != nil {
		return 0, err
	}
	affected, err := s.repo.DispatchPending(ctx, packageID)
	if err != nil {
		return 0, err
	}
	pkg, err := s.repo.Find(ctx, packageID)
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
	pkgIDs, err := s.repo.PendingPackageIDs(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, id := range pkgIDs {
		n, err := s.Dispatch(ctx, id)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

// ReportStatus 上报某台设备对指定升级包的升级结果（in_progress/success/failed），
// 更新设备升级记录并重新聚合升级包状态。
func (s *OTAService) ReportStatus(ctx context.Context, packageID int64, deviceKey string, status string) error {
	if status != "in_progress" && status != "success" && status != "failed" {
		return errors.New("invalid upgrade status: " + status)
	}
	device, err := s.deviceRepo.FindByKey(ctx, deviceKey)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateDeviceStatus(ctx, packageID, device.ID, status, ""); err != nil {
		return err
	}
	return s.recomputePackageStatus(ctx, packageID)
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

func (s *OTAService) Statistics(ctx context.Context, packageName string) (UpgradeStatistics, error) {
	pkg, err := s.repo.FindByName(ctx, packageName)
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

func (s *OTAService) Batches(ctx context.Context, packageName string) ([]model.UpgradeBatch, error) {
	pkg, err := s.repo.FindByName(ctx, packageName)
	if err != nil {
		return nil, err
	}
	return s.repo.Batches(ctx, pkg.ID)
}

func (s *OTAService) Deployments(ctx context.Context, packageName string, page, size int, status string) ([]model.DeviceDeployment, int64, error) {
	pkg, err := s.repo.FindByName(ctx, packageName)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.Deployments(ctx, pkg.ID, page, size, status)
}
