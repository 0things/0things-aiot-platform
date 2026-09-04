package service

import (
	"context"
	"errors"
	"strings"

	devicegroupv1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"

	"github.com/google/uuid"
)

var ErrInvalidDeviceGroupType = errors.New("invalid device group type")

type DeviceGroupServiceInterface interface {
	Create(context.Context, *devicegroupv1.CreateDeviceGroupRequest) (*model.DeviceGroup, error)
	List(context.Context, int, int, string, string) ([]model.DeviceGroup, int64, error)
	Get(context.Context, string) (*model.DeviceGroup, error)
	Update(context.Context, string, *devicegroupv1.UpdateDeviceGroupRequest) (*model.DeviceGroup, error)
	Delete(context.Context, string) error
	AddDevices(context.Context, string, []string) error
	RemoveDevices(context.Context, string, []string) error
	Devices(context.Context, string, int, int, string, string) ([]model.Device, int64, error)
	Preview(context.Context, string) ([]model.Device, int64, error)
}

type DeviceGroupService struct {
	repo    *repository.DeviceGroupRepository
	devices *repository.DeviceRepository
}

func NewDeviceGroupService(repo *repository.DeviceGroupRepository, devices *repository.DeviceRepository) *DeviceGroupService {
	return &DeviceGroupService{repo: repo, devices: devices}
}

func (s *DeviceGroupService) Create(ctx context.Context, req *devicegroupv1.CreateDeviceGroupRequest) (*model.DeviceGroup, error) {
	if req.Type != model.DeviceGroupTypeManual && req.Type != model.DeviceGroupTypeDynamic {
		return nil, ErrInvalidDeviceGroupType
	}
	if req.Type == model.DeviceGroupTypeDynamic && strings.TrimSpace(req.Rule) == "" {
		return nil, errors.New("rule is required for dynamic group")
	}
	if req.Type == model.DeviceGroupTypeDynamic {
		if _, _, err := s.Preview(ctx, req.Rule); err != nil {
			return nil, err
		}
	}
	if req.Type == model.DeviceGroupTypeManual {
		req.Rule = ""
	}
	exists, err := s.repo.NameExists(ctx, req.Name, "")
	if err != nil || exists {
		if exists {
			return nil, errors.New("device group name already exists")
		}
		return nil, err
	}
	group := &model.DeviceGroup{GroupUUID: uuid.NewString(), OrganizationID: tenant.GetOrganizationID(ctx), Name: req.Name, Type: req.Type, Description: req.Description, Rule: req.Rule}
	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *DeviceGroupService) List(ctx context.Context, page, size int, search, groupType string) ([]model.DeviceGroup, int64, error) {
	return s.repo.List(ctx, page, size, search, groupType)
}

func (s *DeviceGroupService) Get(ctx context.Context, groupUUID string) (*model.DeviceGroup, error) {
	return s.repo.FindByUUID(ctx, groupUUID)
}

func (s *DeviceGroupService) Update(ctx context.Context, groupUUID string, req *devicegroupv1.UpdateDeviceGroupRequest) (*model.DeviceGroup, error) {
	group, err := s.Get(ctx, groupUUID)
	if err != nil {
		return nil, err
	}
	if req.Name != group.Name {
		exists, err := s.repo.NameExists(ctx, req.Name, groupUUID)
		if err != nil || exists {
			if exists {
				return nil, errors.New("device group name already exists")
			}
			return nil, err
		}
	}
	if group.Type == model.DeviceGroupTypeDynamic && strings.TrimSpace(req.Rule) == "" {
		return nil, errors.New("rule is required for dynamic group")
	}
	if group.Type == model.DeviceGroupTypeDynamic {
		if _, _, err := s.Preview(ctx, req.Rule); err != nil {
			return nil, err
		}
	}
	group.Name, group.Description, group.Rule = req.Name, req.Description, req.Rule
	if group.Type == model.DeviceGroupTypeManual {
		group.Rule = ""
	}
	if err := s.repo.Save(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *DeviceGroupService) Delete(ctx context.Context, groupUUID string) error {
	group, err := s.Get(ctx, groupUUID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, group)
}

func (s *DeviceGroupService) AddDevices(ctx context.Context, groupUUID string, keys []string) error {
	group, err := s.Get(ctx, groupUUID)
	if err != nil {
		return err
	}
	if group.Type != model.DeviceGroupTypeManual {
		return errors.New("dynamic group members are managed by rule")
	}
	devices, err := s.devices.FindByKeys(ctx, keys)
	if err != nil {
		return err
	}
	if len(devices) != len(keys) {
		return repository.ErrNotFound
	}
	ids := make([]int64, len(devices))
	for i, device := range devices {
		ids[i] = device.ID
	}
	return s.repo.AddDevices(ctx, group.ID, ids)
}

func (s *DeviceGroupService) RemoveDevices(ctx context.Context, groupUUID string, keys []string) error {
	group, err := s.Get(ctx, groupUUID)
	if err != nil {
		return err
	}
	if group.Type != model.DeviceGroupTypeManual {
		return errors.New("dynamic group members are managed by rule")
	}
	devices, err := s.devices.FindByKeys(ctx, keys)
	if err != nil {
		return err
	}
	ids := make([]int64, len(devices))
	for i, device := range devices {
		ids[i] = device.ID
	}
	return s.repo.RemoveDevices(ctx, group.ID, ids)
}

func (s *DeviceGroupService) Devices(ctx context.Context, groupUUID string, page, size int, productKey, search string) ([]model.Device, int64, error) {
	group, err := s.Get(ctx, groupUUID)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.DevicesPage(ctx, group, page, size, productKey, search)
}

func (s *DeviceGroupService) Preview(ctx context.Context, rule string) ([]model.Device, int64, error) {
	return s.repo.Devices(ctx, &model.DeviceGroup{Type: model.DeviceGroupTypeDynamic, Rule: rule})
}
