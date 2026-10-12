package tsdb

import (
	"context"
	"sort"
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
		ts := rec.Timestamp
		if ts >= startTime && ts <= endTime {
			points = append(points, Point{
				Timestamp: ts,
				Metric:    rec.Metric,
				Value:     rec.Value,
			})
		}
	}
	sort.Slice(points, func(i, j int) bool {
		if filter.Descending {
			return points[i].Timestamp > points[j].Timestamp
		}
		return points[i].Timestamp < points[j].Timestamp
	})
	if len(points) > limit {
		points = points[:limit]
	}

	return points, nil
}

func (m *MockClient) Close() error {
	return nil
}
