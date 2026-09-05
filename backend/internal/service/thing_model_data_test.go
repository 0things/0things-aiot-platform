package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

type fakeThingModelDevices struct {
	device *model.Device
	err    error
}

func (f fakeThingModelDevices) FindByKey(context.Context, string) (*model.Device, error) {
	return f.device, f.err
}

type fakeThingModelTSLs struct {
	tsl *model.ProductTSL
	err error
}

func (f fakeThingModelTSLs) FindByProductID(context.Context, int64) (*model.ProductTSL, error) {
	return f.tsl, f.err
}

type fakeThingModelTelemetry struct {
	points map[string]model.TelemetryPoint
}

func (f fakeThingModelTelemetry) QueryLatest(context.Context, string, []string) (map[string]model.TelemetryPoint, error) {
	return f.points, nil
}

func TestThingModelPropertyService_List(t *testing.T) {
	now := time.Now().UnixMilli()
	service := ThingModelDataService{
		devices: fakeThingModelDevices{device: &model.Device{ProductID: 9}},
		tsls:    fakeThingModelTSLs{tsl: &model.ProductTSL{TSL: `{"properties":[{"identifier":"temperature","name":"Temperature","accessMode":"r","dataType":{"type":"double","specs":{"unit":"°C"}}},{"identifier":"humidity","name":"Humidity","accessMode":"rw","dataType":{"type":"double","specs":{"unit":"%"}}}]}`}},
		telemetry: fakeThingModelTelemetry{points: map[string]model.TelemetryPoint{
			"temperature": {Property: "temperature", Value: 24.5, Timestamp: now},
			"undefined":   {Property: "undefined", Value: "ignored", Timestamp: now},
		}},
	}

	properties, err := service.ListProperties(context.Background(), "device-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(properties) != 2 || properties[0].Identifier != "temperature" || properties[1].Identifier != "humidity" {
		t.Fatalf("expected TSL-ordered definitions, got %#v", properties)
	}
	if properties[0].Value != 24.5 || properties[0].ReportedAt == nil {
		t.Fatalf("expected temperature latest point, got %#v", properties[0])
	}
	if properties[1].Value != nil || properties[1].ReportedAt != nil {
		t.Fatalf("expected unreported property to be null, got %#v", properties[1])
	}
}

func TestThingModelPropertyService_ListErrors(t *testing.T) {
	service := ThingModelDataService{
		devices: fakeThingModelDevices{err: repository.ErrNotFound},
		tsls:    fakeThingModelTSLs{}, telemetry: fakeThingModelTelemetry{},
	}
	if _, err := service.ListProperties(context.Background(), "missing"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected device not found, got %v", err)
	}

	service = ThingModelDataService{
		devices: fakeThingModelDevices{device: &model.Device{ProductID: 9}},
		tsls:    fakeThingModelTSLs{tsl: &model.ProductTSL{TSL: "not-json"}}, telemetry: fakeThingModelTelemetry{},
	}
	if _, err := service.ListProperties(context.Background(), "device-1"); !errors.Is(err, ErrInvalidThingModel) {
		t.Fatalf("expected invalid TSL error, got %v", err)
	}
}
