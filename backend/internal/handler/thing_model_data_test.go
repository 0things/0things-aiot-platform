package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type fakeThingModelPropertyService struct {
	properties []model.ThingModelProperty
	err        error
}

func (f fakeThingModelPropertyService) ListProperties(context.Context, string) ([]model.ThingModelProperty, error) {
	return f.properties, f.err
}

func (f fakeThingModelPropertyService) ListServiceInvocations(context.Context, dto.ListDeviceServiceInvocationsQuery) ([]model.DeviceServiceInvocation, int64, error) {
	return nil, 0, nil
}

func TestThingModelPropertyHandler_List(t *testing.T) {
	reportedAt := time.Now()
	h := NewThingModelDataHandler(baseHandler(t), fakeThingModelPropertyService{properties: []model.ThingModelProperty{{
		Identifier: "temperature", Name: "Temperature", DataType: "double", Unit: "°C", AccessMode: "r", Value: 25.5, ReportedAt: &reportedAt,
	}}})
	c, w := hctx(http.MethodGet, "/devices/device-1/thing-model/properties", nil, nil, gin.Params{{Key: "deviceKey", Value: "device-1"}})
	h.ListProperties(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestThingModelPropertyHandler_ListError(t *testing.T) {
	h := NewThingModelDataHandler(baseHandler(t), fakeThingModelPropertyService{err: repository.ErrNotFound})
	c, w := hctx(http.MethodGet, "/devices/missing/thing-model/properties", nil, nil, gin.Params{{Key: "deviceKey", Value: "missing"}})
	h.ListProperties(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
