package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/xuri/excelize/v2"
)

type DeviceService struct {
	repo        *repository.DeviceRepository
	products    *repository.ProductRepository
	tags        *repository.DeviceTagRepository
	shadows     *repository.DeviceShadowRepository
	pushRecords *repository.PushRecordRepository
}

type BatchUploadError struct {
	Row        int
	ProductKey string
	DeviceName string
	Error      string
}

type DeviceStatistics struct {
	TotalDevices     int64
	ActivatedDevices int64
	OnlineDevices    int64
	OfflineDevices   int64
	InactiveDevices  int64
}

type MQTTParameters struct {
	ClientID    string
	Username    string
	MQTTHostURL string
	Password    string
	Port        int32
}

func NewDeviceService(
	repo *repository.DeviceRepository,
	products *repository.ProductRepository,
	tags *repository.DeviceTagRepository,
	shadows *repository.DeviceShadowRepository,
	pushRecords *repository.PushRecordRepository,
) *DeviceService {
	return &DeviceService{repo: repo, products: products, tags: tags, shadows: shadows, pushRecords: pushRecords}
}

func normalizeDeviceMetadata(value json.RawMessage) (json.RawMessage, error) {
	if len(value) == 0 {
		return value, nil
	}

	var legacyString string
	if json.Unmarshal(value, &legacyString) == nil {
		if !json.Valid([]byte(legacyString)) {
			return nil, errors.New("invalid metadata")
		}
		return json.RawMessage(legacyString), nil
	}
	if !json.Valid(value) {
		return nil, errors.New("invalid metadata")
	}
	return value, nil
}

func deviceKey() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "D" + strings.ToUpper(hex.EncodeToString(b))
}

func (s *DeviceService) CreateDevice(ctx context.Context, d *model.Device) (*model.Device, error) {
	var err error
	if d.Metadata, err = normalizeDeviceMetadata(d.Metadata); err != nil {
		return nil, err
	}
	if _, err := s.products.Find(ctx, d.ProductID); err != nil {
		return nil, err
	}
	if d.DeviceKey == "" {
		d.DeviceKey = deviceKey()
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return s.repo.Find(ctx, d.ID)
}
func (s *DeviceService) Device(ctx context.Context, id int64) (*model.Device, error) {
	return s.repo.Find(ctx, id)
}
func (s *DeviceService) DeviceByKey(ctx context.Context, key string) (*model.Device, error) {
	return s.repo.FindByKey(ctx, key)
}
func (s *DeviceService) ListDevices(ctx context.Context, page, size int, productID int64, states []string, enabled *bool, search string) ([]model.Device, int64, error) {
	return s.repo.List(ctx, page, size, productID, states, enabled, search)
}
func (s *DeviceService) UpdateDevice(ctx context.Context, id int64, name, state string, metadata json.RawMessage) (*model.Device, error) {
	d, err := s.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		d.Name = name
	}
	if len(metadata) > 0 {
		d.Metadata = metadata
	}
	if state != "" && state != d.State.State {
		valid := map[string][]string{"inactive": {"offline"}, "offline": {"online"}, "online": {"offline"}}
		ok := false
		for _, v := range valid[d.State.State] {
			if v == state {
				ok = true
			}
		}
		if !ok {
			return nil, errors.New("invalid status transition")
		}
		d.State.State = state
	}
	err = s.repo.Save(ctx, d)
	if err != nil {
		return nil, err
	}
	return s.Device(ctx, id)
}
func (s *DeviceService) Activate(ctx context.Context, id int64) (*model.Device, error) {
	d, err := s.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	if d.State.State != "inactive" {
		return nil, errors.New("device already activated")
	}
	return s.UpdateDevice(ctx, id, "", "offline", nil)
}
func (s *DeviceService) SetEnabled(ctx context.Context, id int64, v bool) (*model.Device, error) {
	d, err := s.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	d.Enabled = v
	err = s.repo.SaveEnabled(ctx, d)
	return d, err
}

func (s *DeviceService) Tags(ctx context.Context, key string) ([]model.DeviceTag, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.tags.ListTags(ctx, d.ID)
}
func (s *DeviceService) SetTags(ctx context.Context, key string, tags map[string]string, replace bool) ([]model.DeviceTag, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	for k := range tags {
		if k == "" || len(k) > 128 {
			return nil, errors.New("invalid tag key")
		}
	}
	err = s.tags.SetTags(ctx, d.ID, tags, replace)
	if err != nil {
		return nil, err
	}
	return s.Tags(ctx, key)
}
func (s *DeviceService) RemoveTags(ctx context.Context, key string, keys []string) error {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return err
	}
	return s.tags.DeleteTags(ctx, d.ID, keys)
}
func (s *DeviceService) Shadow(ctx context.Context, key string) (*model.DeviceShadow, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.shadows.GetShadow(ctx, d.ID)
}
func (s *DeviceService) MutateShadow(ctx context.Context, key string, version int64, source string, desired, reported *map[string]any, clear bool) (*model.DeviceShadow, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.shadows.MutateShadow(ctx, d.ID, version, source, desired, reported, clear)
}
func (s *DeviceService) ShadowHistory(ctx context.Context, key string) ([]model.DeviceShadowHistory, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.shadows.ListShadowHistory(ctx, d.ID)
}
func (s *DeviceService) Telemetry(ctx context.Context, key string) (string, error) {
	return s.repo.Telemetry(ctx, key)
}
func (s *DeviceService) Stats(ctx context.Context) (DeviceStatistics, error) {
	counts, err := s.repo.Statistics(ctx)
	if err != nil {
		return DeviceStatistics{}, err
	}
	return DeviceStatistics{
		TotalDevices: counts.TotalDevices, ActivatedDevices: counts.ActivatedDevices,
		OnlineDevices: counts.OnlineDevices, OfflineDevices: counts.OfflineDevices, InactiveDevices: counts.InactiveDevices,
	}, nil
}
func (s *DeviceService) MQTT(ctx context.Context, key string) (MQTTParameters, error) {
	d, err := s.DeviceByKey(ctx, key)
	if err != nil {
		return MQTTParameters{}, err
	}
	return MQTTParameters{
		ClientID: d.Product.ProductKey, Username: d.DeviceKey, MQTTHostURL: "x39c51fe.ala.cn-hangzhou.emqxsl.cn", Password: d.DeviceKey, Port: 8883,
	}, nil
}
func (s *DeviceService) DeleteDevice(ctx context.Context, id int64) error {
	d, err := s.Device(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, d)
}

func (s *DeviceService) SimulatePush(ctx context.Context, deviceKey, payload, createdBy string) (*model.DevicePushRecord, error) {
	device, err := s.DeviceByKey(ctx, deviceKey)
	if err != nil {
		return nil, err
	}
	if createdBy == "" {
		createdBy = "system"
	}
	record := &model.DevicePushRecord{
		DeviceID:      device.ID,
		OperationType: "Property",
		OperationName: "simulate",
		Payload:       payload,
		Status:        "success",
		CreatedBy:     createdBy,
	}
	if err := s.pushRecords.CreatePushRecord(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *DeviceService) ListPushRecords(ctx context.Context, deviceKey string, page, size int, operationType, status string) ([]model.DevicePushRecord, int64, error) {
	device, err := s.DeviceByKey(ctx, deviceKey)
	if err != nil {
		return nil, 0, err
	}
	return s.pushRecords.ListPushRecords(ctx, device.ID, page, size, operationType, status)
}

func (s *DeviceService) PushRecord(ctx context.Context, id int64) (*model.DevicePushRecord, error) {
	return s.pushRecords.FindPushRecord(ctx, id)
}

func (s *DeviceService) ClearPushRecords(ctx context.Context, deviceKey string, before *time.Time) (int64, error) {
	device, err := s.DeviceByKey(ctx, deviceKey)
	if err != nil {
		return 0, err
	}
	return s.pushRecords.DeletePushRecords(ctx, device.ID, before)
}
func (s *DeviceService) BatchTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "设备导入模板"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "productKey")
	f.SetCellValue(sheet, "B1", "deviceName")
	f.SetCellValue(sheet, "A2", "WMS003")
	f.SetCellValue(sheet, "B2", "DEV001")
	var b bytes.Buffer
	err := f.Write(&b)
	return b.Bytes(), err
}
func (s *DeviceService) BatchCreate(ctx context.Context, content []byte) (int, []BatchUploadError, error) {
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return 0, nil, err
	}
	if len(rows) < 2 {
		return 0, nil, errors.New("no valid data rows")
	}
	head := map[string]int{}
	for i, v := range rows[0] {
		head[strings.ToLower(strings.TrimSpace(v))] = i
	}
	pk, ok1 := head["productkey"]
	dn, ok2 := head["devicename"]
	if !ok1 || !ok2 {
		return 0, nil, errors.New("missing required column")
	}
	success := 0
	errs := []BatchUploadError{}
	for i, row := range rows[1:] {
		val := func(n int) string {
			if n < len(row) {
				return strings.TrimSpace(row[n])
			}
			return ""
		}
		key, name := val(pk), val(dn)
		if key == "" && name == "" {
			continue
		}
		p, e := s.products.FindByKey(ctx, key)
		if e == nil {
			_, e = s.CreateDevice(ctx, &model.Device{Name: name, ProductID: p.ID})
		}
		if e != nil {
			errs = append(errs, BatchUploadError{Row: i + 2, ProductKey: key, DeviceName: name, Error: e.Error()})
		} else {
			success++
		}
	}
	return success, errs, nil
}
func (s *DeviceService) MockKafka(ctx context.Context, brokers []string, topic, data string) error {
	if len(brokers) == 0 {
		return errors.New("kafka producer not initialized")
	}
	c, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
	if err != nil {
		return err
	}
	defer c.Close()
	return c.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: []byte(data)}).FirstErr()
}
func (s *DeviceService) RestoreDevice(ctx context.Context, id int64) (*model.Device, error) {
	if err := s.repo.Restore(ctx, id); err != nil {
		return nil, err
	}
	return s.Device(ctx, id)
}
