package service

import (
	"context"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

type DeviceServiceInvocationServiceInterface interface {
	List(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error)
}

type DeviceServiceInvocationService struct {
	repo    *repository.DeviceServiceInvocationRepository
	devices *repository.DeviceRepository
}

func NewDeviceServiceInvocationService(repo *repository.DeviceServiceInvocationRepository, devices *repository.DeviceRepository) *DeviceServiceInvocationService {
	return &DeviceServiceInvocationService{repo: repo, devices: devices}
}

func (s *DeviceServiceInvocationService) List(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error) {
	device, err := s.devices.FindByKey(ctx, query.DeviceKey)
	if err != nil {
		return nil, 0, err
	}
	query.DeviceID = device.ID
	return s.repo.List(ctx, query)
}
