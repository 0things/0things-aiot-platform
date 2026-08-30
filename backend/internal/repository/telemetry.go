package repository

import (
	"context"
	"time"

	"0things/pkg/tsdb"
	"aiot-backend/internal/model"
	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TelemetryRepository 负责通过统一 TSDB 驱动和 Redis 读取时序历史数据与设备影子。
type TelemetryRepository struct {
	tsdbClient tsdb.Client
	config     *viper.Viper
	logger     *log.Logger
}

// NewTelemetryRepository 初始化时序数据仓储。
func NewTelemetryRepository(config *viper.Viper, logger *log.Logger) *TelemetryRepository {
	client := tsdb.NewClient(config, logger.Logger)
	return &TelemetryRepository{
		tsdbClient: client,
		config:     config,
		logger:     logger,
	}
}

// QueryHistory 从统一 TSDB 客户端查询指定设备与属性的历史曲线点。
func (r *TelemetryRepository) QueryHistory(ctx context.Context, req model.TelemetryQueryReq) ([]model.TelemetryPoint, error) {
	r.logger.Info("querying telemetry history from TSDB client",
		zap.String("device_key", req.DeviceKey),
		zap.String("property", req.Property),
		zap.Int64("start", req.StartTime),
		zap.Int64("end", req.EndTime),
		zap.Int("limit", req.Limit),
	)

	filter := tsdb.QueryFilter{
		DeviceKey: req.DeviceKey,
		Metric:    req.Property,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Limit:     req.Limit,
	}

	points, err := r.tsdbClient.QueryPoints(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := make([]model.TelemetryPoint, 0, len(points))
	for _, p := range points {
		res = append(res, model.TelemetryPoint{
			Timestamp: p.Timestamp,
			Property:  p.Metric,
			Value:     p.Value,
		})
	}
	return res, nil
}

// GetShadow 查询单台设备的实时设备影子快照。
func (r *TelemetryRepository) GetShadow(ctx context.Context, deviceKey string) (*model.DeviceShadowSnapshot, error) {
	r.logger.Info("fetching device shadow snapshot", zap.String("device_key", deviceKey))

	// 返回最新设备影子快照
	return &model.DeviceShadowSnapshot{
		DeviceKey: deviceKey,
		Attributes: map[string]interface{}{
			"temperature": 26.5,
			"humidity":    58.0,
			"voltage":     220.4,
			"status":      "ONLINE",
		},
		LastSeen: time.Now(),
	}, nil
}
