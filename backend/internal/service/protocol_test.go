package service

import (
	"context"
	"testing"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/security"
	"aiot-backend/internal/tenant"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestProtocolService_ListDeviceEndpointsUsesDeviceMQTTCredentials(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.ProductProtocol{}, &model.Device{}, &model.DeviceCredential{}))
	require.NoError(t, db.Create(&model.Product{ID: 1, ProductKey: "P001", OrganizationID: "org-1"}).Error)
	require.NoError(t, db.Create(&model.ProductProtocol{ProductID: 1, TransportProtocol: string(enum.TransportMQTT), ApplicationProtocol: string(enum.ApplicationJSON)}).Error)
	require.NoError(t, db.Create(&model.Device{ID: 1, DeviceUUID: "device-uuid-1", DeviceKey: "D001", ProductID: 1, OrganizationID: "org-1"}).Error)
	ciphertext, err := security.EncryptCredentials("mqtt-pass", "test-key")
	require.NoError(t, err)
	require.NoError(t, db.Create(&model.DeviceCredential{DeviceUUID: "device-uuid-1", CredentialType: "mqtt", Username: "mqtt-user", Password: "hash", PasswordCiphertext: ciphertext, Enabled: true}).Error)

	config := viper.New()
	config.Set("device_gateway.mqtt.host", "broker.example")
	config.Set("device_gateway.mqtt.port", "1884")
	config.Set("security.device_credentials_key", "test-key")
	svc := NewProtocolService(repository.NewProtocolRepository(db), config)

	result, err := svc.ListDeviceEndpoints(tenant.WithOrganization(context.Background(), "org-1"), "D001")
	require.NoError(t, err)
	require.NotNil(t, result.MQTT)
	require.Equal(t, "broker.example", result.MQTT.Host)
	require.Equal(t, "1884", result.MQTT.Port)
	require.Equal(t, "mqtt-user", result.MQTT.Username)
	require.Equal(t, "mqtt-pass", result.MQTT.Password)
}
