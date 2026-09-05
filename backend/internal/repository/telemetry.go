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
	return NewTelemetryRepositoryWithClient(tsdb.NewClient(config, logger.Logger), config, logger)
}

// NewTelemetryRepositoryWithClient allows tests to provide a deterministic TSDB client.
func NewTelemetryRepositoryWithClient(client tsdb.Client, config *viper.Viper, logger *log.Logger) *TelemetryRepository {
	return &TelemetryRepository{
		tsdbClient: client,
		config:     config,
		logger:     logger,
	}
}

// QueryLatest returns the newest real point for each requested property. The
// storage access remains server-side so callers never fan out history requests.
func (r *TelemetryRepository) QueryLatest(ctx context.Context, deviceKey string, identifiers []string) (map[string]model.TelemetryPoint, error) {
	latest := make(map[string]model.TelemetryPoint, len(identifiers))
	for _, identifier := range identifiers {
		points, err := r.tsdbClient.QueryPoints(ctx, tsdb.QueryFilter{
			DeviceKey:           deviceKey,
			Metric:              identifier,
			StartTime:           1,
			EndTime:             time.Now().UnixMilli(),
			Limit:               1,
			Descending:          true,
			DisableMockFallback: true,
		})
		if err != nil {
			return nil, err
		}
		if len(points) == 0 {
			continue
		}
		latest[identifier] = model.TelemetryPoint{Timestamp: points[0].Timestamp, Property: points[0].Metric, Value: points[0].Value}
	}
	return latest, nil
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
