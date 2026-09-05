package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"aiot-backend/internal/dto"
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

var (
	groupRuleCondition = regexp.MustCompile(`(?i)^\s*([a-zA-Z_][a-zA-Z0-9_.-]*)\s*(=|!=|like|in)\s*(.+?)\s*$`)
	groupRuleOrSplit   = regexp.MustCompile(`(?i)\s+OR\s+`)
	groupRuleAndSplit  = regexp.MustCompile(`(?i)\s+AND\s+`)
)

var groupRuleColumns = map[string]string{
	"device_key":  "devices.device_key",
	"name":        "devices.name",
	"product_key": "products.product_key",
	"enabled":     "devices.enabled",
	"state":       "device_states.state",
}

// applyGroupRule securely parses dynamic rule expressions and converts them into GORM query conditions.
// Supports JSON key-value rules as well as AND/OR composite expressions with IN/=/!=/LIKE operators against whitelisted fields.
func applyGroupRule(query *gorm.DB, rule string) (*gorm.DB, error) {
	trimmed := strings.TrimSpace(rule)
	if trimmed == "" {
		return query, nil
	}

	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var jsonMap map[string]any
		if err := json.Unmarshal([]byte(trimmed), &jsonMap); err == nil {
			if len(jsonMap) == 0 {
				return query, nil
			}
			keys := make([]string, 0, len(jsonMap))
			for k := range jsonMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			var conditions []string
			var allArgs []any
			for _, k := range keys {
				v := jsonMap[k]
				field := strings.ToLower(k)
				column, supported := groupRuleColumns[field]
				if !supported && !strings.HasPrefix(field, "tag.") {
					return nil, fmt.Errorf("unsupported group rule field: %s", field)
				}
				if field == "tag." {
					return nil, fmt.Errorf("unsupported group rule field: %s", field)
				}

				switch val := v.(type) {
				case []any:
					strVals := make([]string, len(val))
					for i, item := range val {
						strVals[i] = fmt.Sprintf("%v", item)
					}
					if strings.HasPrefix(field, "tag.") {
						conditions = append(conditions, "EXISTS (SELECT 1 FROM device_tags dt WHERE dt.device_id = devices.id AND dt.key = ? AND dt.value IN ?)")
						allArgs = append(allArgs, strings.TrimPrefix(field, "tag."), strVals)
					} else {
						conditions = append(conditions, column+" IN ?")
						allArgs = append(allArgs, strVals)
					}
				case []string:
					if strings.HasPrefix(field, "tag.") {
						conditions = append(conditions, "EXISTS (SELECT 1 FROM device_tags dt WHERE dt.device_id = devices.id AND dt.key = ? AND dt.value IN ?)")
						allArgs = append(allArgs, strings.TrimPrefix(field, "tag."), val)
					} else {
						conditions = append(conditions, column+" IN ?")
						allArgs = append(allArgs, val)
					}
				default:
					if field == "enabled" {
						bVal := false
						if b, ok := val.(bool); ok {
							bVal = b
						} else if s, ok := val.(string); ok && (strings.ToLower(s) == "true" || s == "1") {
							bVal = true
						}
						conditions = append(conditions, column+" = ?")
						allArgs = append(allArgs, bVal)
					} else if strings.HasPrefix(field, "tag.") {
						conditions = append(conditions, "EXISTS (SELECT 1 FROM device_tags dt WHERE dt.device_id = devices.id AND dt.key = ? AND dt.value = ?)")
						allArgs = append(allArgs, strings.TrimPrefix(field, "tag."), fmt.Sprintf("%v", val))
					} else {
						conditions = append(conditions, column+" = ?")
						allArgs = append(allArgs, val)
					}
				}
			}
			return query.Where("("+strings.Join(conditions, " AND ")+")", allArgs...), nil
		}
	}

	var branches []string
	var allArgs []any
	for _, orPart := range groupRuleOrSplit.Split(trimmed, -1) {
		var conditions []string
		var branchArgs []any
		for _, raw := range groupRuleAndSplit.Split(orPart, -1) {
			rawTrimmed := strings.TrimSpace(raw)
			if rawTrimmed == "" {
				continue
			}
			match := groupRuleCondition.FindStringSubmatch(rawTrimmed)
			if match == nil {
				return nil, fmt.Errorf("invalid group rule condition: %s", rawTrimmed)
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
			var argVal any = value
			if field == "enabled" {
				argVal = strings.ToLower(value) == "true" || value == "1"
			}
			condition := column + " " + operator + " ?"
			args := []any{argVal}
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
		if len(conditions) > 0 {
			branches = append(branches, "("+strings.Join(conditions, " AND ")+")")
			allArgs = append(allArgs, branchArgs...)
		}
	}
	if len(branches) == 0 {
		return query, nil
	}
	return query.Where(strings.Join(branches, " OR "), allArgs...), nil
}

// Devices queries paginated device members for both manual and dynamic device groups.
// Enforces tenant isolation (organization_id) and excludes soft-deleted records.
func (r *DeviceGroupRepository) Devices(ctx context.Context, query dto.ListDeviceGroupDevicesQuery) ([]model.Device, int64, error) {
	base := r.db.WithContext(ctx).
		Table("devices").
		Select("devices.*").
		Joins("LEFT JOIN products ON products.id = devices.product_id").
		Joins("LEFT JOIN device_states ON device_states.device_key = devices.device_key").
		Where("devices.organization_id = ? AND devices.deleted_at IS NULL", tenant.GetOrganizationID(ctx))

	validProductKeys := make([]string, 0, len(query.ProductKeys))
	for _, pk := range query.ProductKeys {
		if trimmed := strings.TrimSpace(pk); trimmed != "" {
			validProductKeys = append(validProductKeys, trimmed)
		}
	}
	if len(validProductKeys) > 0 {
		base = base.Where("products.product_key IN ?", validProductKeys)
	}
	if strings.TrimSpace(query.Search) != "" {
		value := "%" + strings.TrimSpace(query.Search) + "%"
		base = base.Where("(devices.device_key LIKE ? OR devices.name LIKE ? OR products.name LIKE ?)", value, value, value)
	}

	var q *gorm.DB
	if query.GroupType == model.DeviceGroupTypeManual {
		ids, err := r.DeviceIDs(ctx, query.GroupID)
		if err != nil {
			return nil, 0, err
		}
		if len(ids) == 0 {
			return []model.Device{}, 0, nil
		}
		q = base.Where("devices.id IN ?", ids)
	} else {
		var err error
		q, err = applyGroupRule(base, query.Rule)
		if err != nil {
			return nil, 0, err
		}
	}

	var total int64
	var countResult struct {
		Total int64 `gorm:"column:total"`
	}
	if err := q.Session(&gorm.Session{}).Select("COUNT(*) AS total").Scan(&countResult).Error; err != nil {
		return nil, 0, err
	}
	total = countResult.Total

	var devices []model.Device
	findQuery := q.Session(&gorm.Session{}).Select("devices.*").Order("devices.created_at DESC")
	if query.Page > 0 && query.PageSize > 0 {
		findQuery = findQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
	}
	err := findQuery.Preload("Product").Preload("State").Find(&devices).Error
	return devices, total, err
}
