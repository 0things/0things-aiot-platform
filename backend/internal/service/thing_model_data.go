package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

var ErrInvalidThingModel = errors.New("invalid product TSL")

type ThingModelDataServiceInterface interface {
	ListProperties(ctx context.Context, deviceKey string) ([]dto.ThingModelProperty, error)
	ListServiceInvocations(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error)
}

type ThingModelDataService struct {
	invocations *repository.DeviceServiceInvocationRepository
	devices     *repository.DeviceRepository
	tsls        *repository.ProductTSLRepository
	telemetry   *repository.TelemetryRepository
}

func NewThingModelDataService(invocations *repository.DeviceServiceInvocationRepository, devices *repository.DeviceRepository, tsls *repository.ProductTSLRepository, telemetry *repository.TelemetryRepository) *ThingModelDataService {
	return &ThingModelDataService{invocations: invocations, devices: devices, tsls: tsls, telemetry: telemetry}
}

func (s *ThingModelDataService) ListServiceInvocations(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error) {
	device, err := s.devices.FindByKey(ctx, query.DeviceKey)
	if err != nil {
		return nil, 0, err
	}
	query.DeviceID = device.ID
	return s.invocations.List(ctx, query)
}

func (s *ThingModelDataService) ListProperties(ctx context.Context, deviceKey string) ([]dto.ThingModelProperty, error) {
	device, err := s.devices.FindByKey(ctx, deviceKey)
	if err != nil {
		return nil, err
	}
	tsl, err := s.tsls.FindByProductID(ctx, device.ProductID)
	if err != nil {
		return nil, err
	}

	var document dto.ProductTSLDocument
	if err := json.Unmarshal([]byte(tsl.TSL), &document); err != nil || document.Items == nil {
		return nil, ErrInvalidThingModel
	}
	identifiers := make([]string, 0, len(document.Items))
	for _, property := range document.Items {
		if property.Identifier == "" || property.Name == "" || property.DataType.Type == "" {
			return nil, ErrInvalidThingModel
		}
		identifiers = append(identifiers, property.Identifier)
	}
	var latest map[string]dto.TelemetryPoint
	if s.telemetry != nil {
		latest, err = s.telemetry.QueryLatest(ctx, deviceKey, identifiers)
		if err != nil {
			return nil, err
		}
	}

	properties := make([]dto.ThingModelProperty, 0, len(document.Items))
	for _, property := range document.Items {
		item := dto.ThingModelProperty{Identifier: property.Identifier, Name: property.Name, DataType: property.DataType.Type, Unit: property.DataType.Specs.Unit, AccessMode: property.AccessMode}
		if point, found := latest[property.Identifier]; found {
			reportedAt := time.UnixMilli(point.Timestamp)
			item.Value = point.Value
			item.ReportedAt = &reportedAt
		}
		properties = append(properties, item)
	}
	return properties, nil
}
