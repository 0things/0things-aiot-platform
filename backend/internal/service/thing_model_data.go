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
	ListProperties(ctx context.Context, deviceKey string) ([]model.ThingModelProperty, error)
	ListServiceInvocations(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error)
}

type thingModelDeviceRepository interface {
	FindByKey(ctx context.Context, key string) (*model.Device, error)
}

type thingModelTSLRepository interface {
	FindByProductID(ctx context.Context, productID int64) (*model.ProductTSL, error)
}

type thingModelTelemetryRepository interface {
	QueryLatest(ctx context.Context, deviceKey string, identifiers []string) (map[string]model.TelemetryPoint, error)
}

type thingModelInvocationRepository interface {
	List(ctx context.Context, query dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error)
}

type ThingModelDataService struct {
	invocations thingModelInvocationRepository
	devices     thingModelDeviceRepository
	tsls        thingModelTSLRepository
	telemetry   thingModelTelemetryRepository
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

type productTSLDocument struct {
	Properties []productTSLProperty `json:"properties"`
}

type productTSLProperty struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	AccessMode string `json:"accessMode"`
	DataType   struct {
		Type  string `json:"type"`
		Specs struct {
			Unit string `json:"unit"`
		} `json:"specs"`
	} `json:"dataType"`
}

func (s *ThingModelDataService) ListProperties(ctx context.Context, deviceKey string) ([]model.ThingModelProperty, error) {
	device, err := s.devices.FindByKey(ctx, deviceKey)
	if err != nil {
		return nil, err
	}
	tsl, err := s.tsls.FindByProductID(ctx, device.ProductID)
	if err != nil {
		return nil, err
	}

	var document productTSLDocument
	if err := json.Unmarshal([]byte(tsl.TSL), &document); err != nil || document.Properties == nil {
		return nil, ErrInvalidThingModel
	}
	identifiers := make([]string, 0, len(document.Properties))
	for _, property := range document.Properties {
		if property.Identifier == "" || property.Name == "" || property.DataType.Type == "" {
			return nil, ErrInvalidThingModel
		}
		identifiers = append(identifiers, property.Identifier)
	}
	latest, err := s.telemetry.QueryLatest(ctx, deviceKey, identifiers)
	if err != nil {
		return nil, err
	}

	properties := make([]model.ThingModelProperty, 0, len(document.Properties))
	for _, property := range document.Properties {
		item := model.ThingModelProperty{Identifier: property.Identifier, Name: property.Name, DataType: property.DataType.Type, Unit: property.DataType.Specs.Unit, AccessMode: property.AccessMode}
		if point, found := latest[property.Identifier]; found {
			reportedAt := time.UnixMilli(point.Timestamp)
			item.Value = point.Value
			item.ReportedAt = &reportedAt
		}
		properties = append(properties, item)
	}
	return properties, nil
}
