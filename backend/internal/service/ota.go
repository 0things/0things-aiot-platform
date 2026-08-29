package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"aiot-backend/internal/enum"
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
	BatchUpgrade(ctx context.Context, uuid string, deviceKeys []string) (*model.UpgradeBatch, error)
	ReportStatus(ctx context.Context, uuid string, deviceKey string, status string) error
	ReportBatchDevice(ctx context.Context, batchID, deviceKey, status, version string, progress int32, desc ...string) error
	Statistics(ctx context.Context, uuid string, batchID ...string) (UpgradeStatistics, error)
	Batches(ctx context.Context, uuid string) ([]model.UpgradeBatch, error)
	Deployments(ctx context.Context, uuid string, page, size int, status string, batchID ...string) ([]model.DeviceDeployment, int64, error)
	CancelBatch(ctx context.Context, uuid, batchID string) error
	RetryBatch(ctx context.Context, uuid, batchID string) error
}

func (s *OTAService) batchPackage(ctx context.Context, uuid, batchID string) (*model.OTAPackage, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.FindBatch(ctx, pkg.ID, batchID); err != nil {
		return nil, repository.ErrNotFound
	}
	return pkg, nil
}

func (s *OTAService) CancelBatch(ctx context.Context, uuid, batchID string) error {
	pkg, err := s.batchPackage(ctx, uuid, batchID)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateBatchDevicesStatus(ctx, pkg.ID, batchID, []string{enum.OTAStatusPending, enum.OTAStatusSent, enum.OTAStatusInProgress}, enum.OTAStatusCancelled); err != nil {
		return err
	}
	return s.repo.UpdateBatchStatus(ctx, batchID, enum.OTAStatusCancelled)
}

func (s *OTAService) RetryBatch(ctx context.Context, uuid, batchID string) error {
	if s.kafka == nil {
		return errors.New("kafka service is required for OTA batch retry")
	}
	pkg, err := s.batchPackage(ctx, uuid, batchID)
	if err != nil {
		return err
	}
	retryLimited, err := s.repo.CountRetryLimited(ctx, pkg.ID, batchID)
	if err != nil {
		return err
	}
	if retryLimited > 0 {
		return errors.New("OTA retry limit exceeded")
	}
	if err := s.repo.UpdateBatchDevicesStatus(ctx, pkg.ID, batchID, []string{enum.OTAStatusFailed, enum.OTAStatusTimeout}, enum.OTAStatusPending); err != nil {
		return err
	}
	deployments, _, err := s.repo.Deployments(ctx, pkg.ID, 1, 100000, enum.OTAStatusPending, batchID)
	if err != nil {
		return err
	}
	for _, d := range deployments {
		message := map[string]any{"batch_id": batchID, "package_id": pkg.ID, "product_key": d.ProductKey, "device_key": d.DeviceKey, "device_name": d.DeviceName, "target_version": d.TargetVersion, "module": "default", "url": pkg.FileURL, "size": pkg.FileSize, "checksum": pkg.Checksum, "transport_protocol": "mqtt"}
		if s.protocolRepo != nil {
			endpoint, endpointErr := s.protocolRepo.DeviceOTAEndpoint(ctx, d.DeviceID)
			if endpointErr != nil {
				return endpointErr
			}
			if endpoint != nil {
				message["endpoint_id"] = endpoint.EndpointID
				message["endpoint"] = endpoint.Endpoint
				message["transport_protocol"] = endpoint.TransportProtocol
			}
		}
		if err := s.kafka.ProduceJSON(ctx, enum.KafkaTopicOTAUpgradeCommandV1, batchID+":"+d.DeviceKey, message); err != nil {
			_ = s.repo.MarkDispatchResult(ctx, pkg.ID, d.DeviceID, batchID, enum.OTAStatusFailed, err.Error())
			return err
		}
		if err := s.repo.RecordKafkaDispatch(ctx, pkg.ID, d.DeviceID, batchID); err != nil {
			return err
		}
	}
	return s.repo.UpdateBatchStatus(ctx, batchID, enum.OTAStatusPending)
}

type OTAService struct {
	repo         *repository.OTARepository
	productRepo  *repository.ProductRepository
	deviceRepo   *repository.DeviceRepository
	protocolRepo *repository.ProtocolRepository
	kafka        KafkaServiceInterface
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

func NewOTAService(repo *repository.OTARepository, productRepo *repository.ProductRepository, deviceRepo *repository.DeviceRepository, kafka KafkaServiceInterface) *OTAService {
	return &OTAService{repo: repo, productRepo: productRepo, deviceRepo: deviceRepo, kafka: kafka}
}

// NewOTAServiceWithProtocol 注入协议仓储，使 OTA 命令可按设备协议端点选择传输方式。
func NewOTAServiceWithProtocol(repo *repository.OTARepository, productRepo *repository.ProductRepository, deviceRepo *repository.DeviceRepository, kafka KafkaServiceInterface, protocols *repository.ProtocolRepository) *OTAService {
	return &OTAService{repo: repo, productRepo: productRepo, deviceRepo: deviceRepo, protocolRepo: protocols, kafka: kafka}
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

// BatchUpgrade 为指定升级包创建静态升级批次，并将所选设备加入该批次的升级
// （每条设备记录状态为 pending，并关联 upgrade_batch_id）。升级包状态置为
// deploying。返回创建好的批次。
func (s *OTAService) BatchUpgrade(ctx context.Context, uuid string, deviceKeys []string) (*model.UpgradeBatch, error) {
	if s.kafka == nil {
		return nil, errors.New("kafka service is required for OTA batch upgrade")
	}
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
		Status:            enum.OTAStatusPending,
		TargetDeviceCount: int32(len(deviceIDs)),
	}
	if err := s.repo.CreateBatchWithDeployments(ctx, batch, pkg.ID, deviceIDs, pkg.Version); err != nil {
		return batch, err
	}
	var dispatchErr error
	for _, device := range devices {
		message := map[string]any{
			"batch_id":           batchID,
			"package_id":         pkg.ID,
			"product_key":        pkg.ProductKey,
			"device_key":         device.DeviceKey,
			"device_name":        device.Name,
			"target_version":     pkg.Version,
			"module":             "default",
			"url":                pkg.FileURL,
			"size":               pkg.FileSize,
			"checksum":           pkg.Checksum,
			"transport_protocol": "mqtt",
		}
		if s.protocolRepo != nil {
			endpoint, endpointErr := s.protocolRepo.DeviceOTAEndpoint(ctx, device.ID)
			if endpointErr != nil {
				return batch, endpointErr
			}
			if endpoint != nil {
				message["endpoint_id"] = endpoint.EndpointID
				message["endpoint"] = endpoint.Endpoint
				message["transport_protocol"] = endpoint.TransportProtocol
			}
		}
		if err := s.kafka.ProduceJSON(ctx, enum.KafkaTopicOTAUpgradeCommandV1, batchID+":"+device.DeviceKey, message); err != nil {
			_ = s.repo.MarkDispatchResult(ctx, pkg.ID, device.ID, batchID, enum.OTAStatusFailed, err.Error())
			if dispatchErr == nil {
				dispatchErr = err
			}
			continue
		}
		if err := s.repo.RecordKafkaDispatch(ctx, pkg.ID, device.ID, batchID); err != nil {
			return batch, err
		}
	}
	// 继续尝试其余设备，避免单台 Kafka 失败导致批次只下发一部分。
	pkg, err = s.repo.Find(ctx, pkg.ID)
	if err != nil {
		return batch, err
	}
	pkg.Status = enum.OTAPackageDeploying
	if err := s.repo.Save(ctx, pkg); err != nil {
		return batch, err
	}
	if dispatchErr != nil {
		return batch, dispatchErr
	}
	return batch, nil
}

// ReportStatus 上报某台设备对指定升级包的升级结果（in_progress/success/failed），
// 更新设备升级记录并重新聚合升级包状态。
func (s *OTAService) ReportStatus(ctx context.Context, uuid string, deviceKey string, status string) error {
	if status != enum.OTAStatusInProgress && status != enum.OTAStatusSuccess && status != enum.OTAStatusFailed {
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

// ReportBatchDevice 消费并处理设备升级状态/进度上报，自动进行版本匹配校验与超时判断，
// 并联动更新批次内设备记录、升级包整体状态以及升级批次状态。
func (s *OTAService) ReportBatchDevice(ctx context.Context, batchID, deviceKey, status, version string, progress int32, desc ...string) error {
	if status != enum.OTAStatusInProgress && status != enum.OTAStatusSuccess && status != enum.OTAStatusFailed && status != enum.OTAStatusTimeout {
		return errors.New("invalid upgrade status: " + status)
	}
	device, err := s.deviceRepo.FindByKeyForEvent(ctx, deviceKey)
	if err != nil {
		return err
	}
	var targetVersion string
	var packageID string
	batch, err := s.repo.FindBatchByID(ctx, batchID)
	if err != nil {
		return err
	}
	packageID = batch.OTAPackageID
	task, err := s.repo.FindBatchDevice(ctx, batchID, device.ID)
	if err != nil {
		return err
	}
	targetVersion = task.TargetVersion
	if task.FirstProgressAt != nil && task.TimeoutSeconds > 0 && time.Now().Unix()-*task.FirstProgressAt > int64(task.TimeoutSeconds) {
		status = enum.OTAStatusTimeout
	}
	if status != enum.OTAStatusTimeout && version != "" && version == targetVersion {
		status = enum.OTAStatusSuccess
	}
	var errorDesc string
	if len(desc) > 0 {
		errorDesc = desc[0]
	}
	if err := s.repo.UpdateBatchDeviceStatus(ctx, batchID, device.ID, status, version, progress, errorDesc); err != nil {
		return err
	}
	packageIDValue, err := strconv.ParseInt(packageID, 10, 64)
	if err != nil {
		return err
	}
	if err := s.recomputePackageStatus(ctx, packageIDValue); err != nil {
		return err
	}
	return s.recomputeBatchStatus(ctx, batchID, packageIDValue)
}

// recomputeBatchStatus 根据批次内设备状态更新批次自身状态。
func (s *OTAService) recomputeBatchStatus(ctx context.Context, batchID string, packageID int64) error {
	counts, err := s.repo.Statistics(ctx, packageID, batchID)
	if err != nil {
		return err
	}
	status := enum.OTAStatusSuccess
	if counts.Pending > 0 || counts.InProgress > 0 {
		status = enum.OTAStatusPending
	} else if counts.Failed > 0 && counts.Success > 0 {
		status = enum.OTAPackagePartial
	} else if counts.Failed > 0 {
		status = enum.OTAStatusFailed
	}
	return s.repo.UpdateBatchStatus(ctx, batchID, status)
}

// ShouldDispatchBatchDevice prevents replayed Kafka commands from being sent
// again after a task has already reached the broker.
// ClaimBatchDeviceForMQTT 原子领取 Kafka 命令，确保同一设备任务只下发一次 MQTT。
func (s *OTAService) ClaimBatchDeviceForMQTT(ctx context.Context, batchID, deviceKey string) (bool, error) {
	return s.repo.ClaimBatchDeviceForMQTT(ctx, batchID, deviceKey)
}

// ResetMQTTDispatch 将 MQTT 发布失败的任务恢复为 pending，供显式重试接口处理。
func (s *OTAService) ResetMQTTDispatch(ctx context.Context, batchID, deviceKey, dispatchError string) error {
	return s.repo.ResetMQTTDispatch(ctx, batchID, deviceKey, dispatchError)
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
		pkg.Status = enum.OTAStatusSuccess
	case counts.Success == 0:
		pkg.Status = enum.OTAStatusFailed
	default:
		pkg.Status = enum.OTAPackagePartial
	}
	return s.repo.Save(ctx, pkg)
}

func (s *OTAService) Statistics(ctx context.Context, uuid string, batchID ...string) (UpgradeStatistics, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return UpgradeStatistics{}, err
	}
	counts, err := s.repo.Statistics(ctx, pkg.ID, batchID...)
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

func (s *OTAService) Deployments(ctx context.Context, uuid string, page, size int, status string, batchID ...string) ([]model.DeviceDeployment, int64, error) {
	pkg, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.Deployments(ctx, pkg.ID, page, size, status, batchID...)
}
