package repository

import (
	"context"
	"errors"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"gorm.io/gorm"
)

type ProtocolRepository struct{ db *IoTDB }

func (r *ProtocolRepository) DeviceByKey(ctx context.Context, deviceKey string) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Preload("Product").Where("device_key = ? AND organization_id = ?", deviceKey, tenant.GetOrganizationID(ctx)).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &device, nil
}

type DeviceEndpointProtocol struct {
	EndpointID          int64
	Endpoint            string
	TransportProtocol   string
	ApplicationProtocol string
}

// MarkEndpointSeenForEvent 根据设备上报的 device_key 和传输协议更新端点心跳。
func (r *ProtocolRepository) MarkEndpointSeenForEvent(ctx context.Context, deviceKey, transport string, seenAt time.Time) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE device_endpoints
		SET status = ?, last_seen_at = ?, updated_at = ?
		WHERE device_id IN (SELECT id FROM devices WHERE device_key = ?)
		  AND product_protocol_id IN (SELECT id FROM product_protocols WHERE transport_protocol = ?)
	`, "online", seenAt, seenAt, deviceKey, transport).Error
}

// MarkStaleEndpointsOffline 将超过心跳窗口仍未上报的在线端点标记为离线。
func (r *ProtocolRepository) MarkStaleEndpointsOffline(ctx context.Context, cutoff time.Time) error {
	return r.db.WithContext(ctx).Model(&model.DeviceEndpoint{}).
		Where("status = ? AND last_seen_at IS NOT NULL AND last_seen_at < ?", "online", cutoff).
		Updates(map[string]any{"status": "offline", "updated_at": time.Now().UTC()}).Error
}

func NewProtocolRepository(db *IoTDB) *ProtocolRepository { return &ProtocolRepository{db: db} }

func (r *ProtocolRepository) ProductProtocols(ctx context.Context, productID int64) ([]model.ProductProtocol, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", productID, tenant.GetOrganizationID(ctx)).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var result []model.ProductProtocol
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("id").Find(&result).Error
	return result, err
}

func (r *ProtocolRepository) DeviceEndpoints(ctx context.Context, deviceID int64) ([]model.DeviceEndpoint, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", deviceID, tenant.GetOrganizationID(ctx)).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// 端点按需物化：访问设备接入信息时，才根据产品协议创建缺失端点。
	var protocols []model.ProductProtocol
	if err := r.db.WithContext(ctx).Where("product_id = ?", device.ProductID).Find(&protocols).Error; err != nil {
		return nil, err
	}
	for _, protocol := range protocols {
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.DeviceEndpoint{}).
			Where("device_id = ? AND product_protocol_id = ?", deviceID, protocol.ID).Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			if err := r.db.WithContext(ctx).Create(&model.DeviceEndpoint{DeviceID: deviceID, ProductProtocolID: protocol.ID, Enabled: true, Status: "inactive"}).Error; err != nil {
				return nil, err
			}
		}
	}
	var result []model.DeviceEndpoint
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("id").Find(&result).Error
	return result, err
}

// DeviceOTAEndpoint returns the first enabled endpoint whose product protocol
// declares OTA support. It keeps OTA routing independent from MQTT details.
func (r *ProtocolRepository) DeviceOTAEndpoint(ctx context.Context, deviceID int64) (*DeviceEndpointProtocol, error) {
	var result DeviceEndpointProtocol
	err := r.db.WithContext(ctx).Table("device_endpoints AS de").
		Select("de.id AS endpoint_id, de.endpoint, pp.transport_protocol, pp.application_protocol").
		Joins("JOIN product_protocols AS pp ON pp.id = de.product_protocol_id").
		Joins("JOIN devices AS d ON d.id = de.device_id").
		Where("de.device_id = ? AND de.enabled = ? AND d.organization_id = ?", deviceID, true, tenant.GetOrganizationID(ctx)).
		Order("de.id").First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ProtocolRepository) SaveDeviceEndpoint(ctx context.Context, item *model.DeviceEndpoint) error {
	var device model.Device
	if err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", item.DeviceID, tenant.GetOrganizationID(ctx)).First(&device).Error; err != nil {
		return err
	}
	var protocol model.ProductProtocol
	if err := r.db.WithContext(ctx).Where("id = ? AND product_id = ?", item.ProductProtocolID, device.ProductID).First(&protocol).Error; err != nil {
		return err
	}
	if item.ID != 0 && item.Credentials == "" {
		var existing model.DeviceEndpoint
		if err := r.db.WithContext(ctx).Select("credentials").Where("id = ? AND device_id = ?", item.ID, item.DeviceID).First(&existing).Error; err != nil {
			return err
		}
		item.Credentials = existing.Credentials
	}
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ProtocolRepository) DeleteDeviceEndpoint(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND device_id IN (SELECT id FROM devices WHERE organization_id = ?)", id, tenant.GetOrganizationID(ctx)).
		Delete(&model.DeviceEndpoint{}).Error
}
