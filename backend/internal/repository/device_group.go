package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

type DeviceGroupRepository struct{ db *gorm.DB }

func NewDeviceGroupRepository(db *gorm.DB) *DeviceGroupRepository {
	return &DeviceGroupRepository{db: db}
}

func (r *DeviceGroupRepository) Create(ctx context.Context, group *model.DeviceGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *DeviceGroupRepository) FindByUUID(ctx context.Context, groupUUID string) (*model.DeviceGroup, error) {
	var group model.DeviceGroup
	err := r.db.WithContext(ctx).Where("group_uuid = ? AND organization_id = ?", groupUUID, tenant.GetOrganizationID(ctx)).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &group, err
}

func (r *DeviceGroupRepository) List(ctx context.Context, page, size int, search, groupType string) ([]model.DeviceGroup, int64, error) {
	query := r.db.WithContext(ctx).Where("organization_id = ?", tenant.GetOrganizationID(ctx)).Order("created_at DESC")
	if strings.TrimSpace(search) != "" {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}
	if strings.TrimSpace(groupType) != "" {
		query = query.Where("type = ?", strings.TrimSpace(groupType))
	}
	var total int64
	if err := query.Model(&model.DeviceGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var groups []model.DeviceGroup
	err := query.Offset((page - 1) * size).Limit(size).Find(&groups).Error
	return groups, total, err
}

func (r *DeviceGroupRepository) Save(ctx context.Context, group *model.DeviceGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *DeviceGroupRepository) Delete(ctx context.Context, group *model.DeviceGroup) error {
	return r.db.WithContext(ctx).Delete(group).Error
}

func (r *DeviceGroupRepository) NameExists(ctx context.Context, name, excludeUUID string) (bool, error) {
	query := r.db.WithContext(ctx).Model(&model.DeviceGroup{}).Where("organization_id = ? AND name = ?", tenant.GetOrganizationID(ctx), name)
	if excludeUUID != "" {
		query = query.Where("group_uuid <> ?", excludeUUID)
	}
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *DeviceGroupRepository) AddDevices(ctx context.Context, groupID int64, deviceIDs []int64) error {
	for _, deviceID := range deviceIDs {
		if err := r.db.WithContext(ctx).Where("group_id = ? AND device_id = ?", groupID, deviceID).FirstOrCreate(&model.DeviceGroupMember{GroupID: groupID, DeviceID: deviceID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceGroupRepository) RemoveDevices(ctx context.Context, groupID int64, deviceIDs []int64) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND device_id IN ?", groupID, deviceIDs).Delete(&model.DeviceGroupMember{}).Error
}

func (r *DeviceGroupRepository) DeviceIDs(ctx context.Context, groupID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Model(&model.DeviceGroupMember{}).Where("group_id = ?", groupID).Pluck("device_id", &ids).Error
	return ids, err
}

func (r *DeviceGroupRepository) DB(ctx context.Context) *gorm.DB { return r.db.WithContext(ctx) }

var groupRuleCondition = regexp.MustCompile(`(?i)^\s*([a-zA-Z_][a-zA-Z0-9_.-]*)\s*(=|!=|like|in)\s*(.+?)\s*$`)

var groupRuleColumns = map[string]string{
	"device_key":  "devices.device_key",
	"name":        "devices.name",
	"product_key": "products.product_key",
	"enabled":     "devices.enabled",
	"state":       "device_states.state",
}

// applyGroupRule 将用户配置的动态规则安全解析并映射为 GORM 查询条件。
// 规则语法支持 AND/OR 组合及 IN/=/!=/LIKE 等操作符，并通过严格的白名单字段校验防止 SQL 注入。
func applyGroupRule(query *gorm.DB, rule string) (*gorm.DB, error) {
	var branches []string
	var allArgs []any
	for _, orPart := range strings.Split(rule, " OR ") {
		var conditions []string
		var branchArgs []any
		for _, raw := range strings.Split(orPart, " AND ") {
			match := groupRuleCondition.FindStringSubmatch(strings.TrimSpace(raw))
			if match == nil {
				return nil, fmt.Errorf("invalid group rule condition: %s", strings.TrimSpace(raw))
			}
			field := strings.ToLower(match[1])
			column, supported := groupRuleColumns[field]
			if !supported && !strings.HasPrefix(field, "tag.") {
				return nil, fmt.Errorf("unsupported group rule field: %s", field)
			}
			if field == "tag." {
				return nil, fmt.Errorf("unsupported group rule field: %s", field)
			}
			operator := strings.ToUpper(match[2])
			value := strings.Trim(strings.TrimSpace(match[3]), "'\"")
			condition := column + " " + operator + " ?"
			args := []any{value}
			if strings.HasPrefix(field, "tag.") {
				condition = "EXISTS (SELECT 1 FROM device_tags dt WHERE dt.device_id = devices.id AND dt.key = ? AND dt.value " + operator + " ?)"
				args = []any{strings.TrimPrefix(field, "tag."), value}
			}
			if operator == "IN" {
				values := strings.Split(value, ",")
				for i := range values {
					values[i] = strings.Trim(strings.TrimSpace(values[i]), "'\"")
				}
				if strings.HasPrefix(field, "tag.") {
					condition, args = "EXISTS (SELECT 1 FROM device_tags dt WHERE dt.device_id = devices.id AND dt.key = ? AND dt.value IN ?)", []any{strings.TrimPrefix(field, "tag."), values}
				} else {
					condition, args = column+" IN ?", []any{values}
				}
			}
			conditions = append(conditions, condition)
			branchArgs = append(branchArgs, args...)
		}
		branches = append(branches, "("+strings.Join(conditions, " AND ")+")")
		allArgs = append(allArgs, branchArgs...)
	}
	return query.Where(strings.Join(branches, " OR "), allArgs...), nil
}

func (r *DeviceGroupRepository) Devices(ctx context.Context, group *model.DeviceGroup) ([]model.Device, int64, error) {
	return r.devices(ctx, group, 0, 0, "", "")
}

func (r *DeviceGroupRepository) DevicesPage(ctx context.Context, group *model.DeviceGroup, page, size int, productKey, search string) ([]model.Device, int64, error) {
	return r.devices(ctx, group, page, size, productKey, search)
}

// devices 统一执行手动分组与动态分组的设备成员分页及过滤查询。
// 严格遵守当前组织隔离 (organization_id) 并过滤已软删除的设备 (deleted_at IS NULL)。
func (r *DeviceGroupRepository) devices(ctx context.Context, group *model.DeviceGroup, page, size int, productKey, search string) ([]model.Device, int64, error) {
	base := r.db.WithContext(ctx).
		Table("devices").
		Select("devices.*").
		Joins("LEFT JOIN products ON products.id = devices.product_id").
		Joins("LEFT JOIN device_states ON device_states.device_key = devices.device_key").
		Where("devices.organization_id = ? AND devices.deleted_at IS NULL", tenant.GetOrganizationID(ctx))

	if strings.TrimSpace(productKey) != "" {
		base = base.Where("products.product_key = ?", strings.TrimSpace(productKey))
	}
	if strings.TrimSpace(search) != "" {
		value := "%" + strings.TrimSpace(search) + "%"
		base = base.Where("(devices.device_key LIKE ? OR devices.name LIKE ?)", value, value)
	}

	var query *gorm.DB
	if group.Type == model.DeviceGroupTypeManual {
		ids, err := r.DeviceIDs(ctx, group.ID)
		if err != nil {
			return nil, 0, err
		}
		if len(ids) == 0 {
			return []model.Device{}, 0, nil
		}
		query = base.Where("devices.id IN ?", ids)
	} else {
		var err error
		query, err = applyGroupRule(base, group.Rule)
		if err != nil {
			return nil, 0, err
		}
	}

	var total int64
	var countResult struct {
		Total int64 `gorm:"column:total"`
	}
	if err := query.Session(&gorm.Session{}).Select("COUNT(*) AS total").Scan(&countResult).Error; err != nil {
		return nil, 0, err
	}
	total = countResult.Total

	var devices []model.Device
	findQuery := query.Session(&gorm.Session{}).Select("devices.*").Order("devices.created_at DESC")
	if page > 0 && size > 0 {
		findQuery = findQuery.Offset((page - 1) * size).Limit(size)
	}
	err := findQuery.Find(&devices).Error
	return devices, total, err
}
