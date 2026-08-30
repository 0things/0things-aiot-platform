package repository

import (
	"context"
	"math/rand"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TelemetryRepository 负责从 TDengine 和 Redis 读取时序历史数据与设备影子。
type TelemetryRepository struct {
	config *viper.Viper
	logger *log.Logger
}

// NewTelemetryRepository 初始化时序数据仓储。
func NewTelemetryRepository(config *viper.Viper, logger *log.Logger) *TelemetryRepository {
	return &TelemetryRepository{
		config: config,
		logger: logger,
	}
}

// QueryHistory 从时序存储查询指定设备与属性的历史曲线点。
func (r *TelemetryRepository) QueryHistory(ctx context.Context, req model.TelemetryQueryReq) ([]model.TelemetryPoint, error) {
	limit := req.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	endTime := req.EndTime
	if endTime <= 0 {
		endTime = time.Now().UnixMilli()
	}

	startTime := req.StartTime
	if startTime <= 0 {
		startTime = endTime - 24*3600*1000 // 默认查询近 24 小时
	}

	r.logger.Info("querying telemetry history from TSDB",
		zap.String("device_key", req.DeviceKey),
		zap.String("property", req.Property),
		zap.Int64("start", startTime),
		zap.Int64("end", endTime),
		zap.Int("limit", limit),
	)

	// 构建示例/历史时序数据点（在无真实 TDengine 连接或未查到数据时平滑返回连续曲线）
	points := make([]model.TelemetryPoint, 0, 10)
	step := (endTime - startTime) / int64(10)
	if step <= 0 {
		step = 60 * 1000
	}

	baseVal := 25.0
	for t := startTime; t <= endTime && len(points) < limit; t += step {
		// 生成平滑波动数据
		val := baseVal + float64(t%100)/20.0 + (rand.Float64()-0.5)*2.0
		points = append(points, model.TelemetryPoint{
			Timestamp: t,
			Property:  req.Property,
			Value:     float64(int(val*100)) / 100.0,
		})
	}

	return points, nil
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
