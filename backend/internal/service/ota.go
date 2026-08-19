package service

import (
	"context"
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
	Statistics(ctx context.Context, packageName string) (UpgradeStatistics, error)
	Batches(ctx context.Context, packageName string) ([]model.UpgradeBatch, error)
	Deployments(ctx context.Context, packageName string, page, size int, status string) ([]model.DeviceDeployment, int64, error)
}

type OTAService struct {
	repo        *repository.OTARepository
	productRepo *repository.ProductRepository
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

func NewOTAService(repo *repository.OTARepository, productRepo *repository.ProductRepository) *OTAService {
	return &OTAService{repo: repo, productRepo: productRepo}
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
