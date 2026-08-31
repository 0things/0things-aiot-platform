package tsdb

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

// MockClient 内存与平滑模拟时序客户端。
type MockClient struct {
	mu     sync.RWMutex
	data   map[string][]Record
	logger *zap.Logger
}

func NewMockClient(logger *zap.Logger) *MockClient {
	logger.Info("TSDB initialized in Mock / In-Memory mode")
	return &MockClient{
		data:   make(map[string][]Record),
		logger: logger,
	}
}

func (m *MockClient) WriteBatch(ctx context.Context, records []Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, rec := range records {
		key := rec.DeviceKey + ":" + rec.Metric
		m.data[key] = append(m.data[key], rec)
	}

	m.logger.Debug("TSDB mock persisted batch", zap.Int("count", len(records)))
	return nil
}

func (m *MockClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	endTime := filter.EndTime
	if endTime <= 0 {
		endTime = time.Now().UnixMilli()
	}

	startTime := filter.StartTime
	if startTime <= 0 {
		startTime = endTime - 24*3600*1000
	}

	// 查内存中缓存的点
	key := filter.DeviceKey + ":" + filter.Metric
	records := m.data[key]

	points := make([]Point, 0)
	for _, rec := range records {
		ts := rec.Timestamp.UnixMilli()
		if ts >= startTime && ts <= endTime {
			points = append(points, Point{
				Timestamp: ts,
				Metric:    rec.Metric,
				Value:     rec.Value,
			})
			if len(points) >= limit {
				break
			}
		}
	}

	// 若内存暂无真实点，平滑返回拟合曲线供前端预览
	if len(points) == 0 {
		step := (endTime - startTime) / int64(10)
		if step <= 0 {
			step = 60 * 1000
		}
		baseVal := 25.0
		for t := startTime; t <= endTime && len(points) < limit; t += step {
			val := baseVal + float64(t%100)/20.0 + (rand.Float64()-0.5)*2.0
			points = append(points, Point{
				Timestamp: t,
				Metric:    filter.Metric,
				Value:     float64(int(val*100)) / 100.0,
			})
		}
	}

	return points, nil
}

func (m *MockClient) Close() error {
	return nil
}
