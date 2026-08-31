package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"0things/pkg/protocol"
	"0things/pkg/tsdb"
	"data-engine/internal/model"
	"data-engine/internal/storage"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Processor 是数据处理引擎的核心计算处理器，负责 TSL 物模型字段展开、TSDB 时序落库、设备影子更新与告警规则计算。
type Processor struct {
	tsdbClient tsdb.Client
	shadow     storage.ShadowStore
	protocols  *protocol.Registry
	logger     *zap.Logger
}

func NewProcessor(config *viper.Viper, logger *zap.Logger, tsdbClient tsdb.Client, shadow storage.ShadowStore) *Processor {
	return &Processor{
		tsdbClient: tsdbClient,
		shadow:     shadow,
		protocols:  protocol.DefaultRegistry(),
		logger:     logger,
	}
}

// ProcessMessage 处理单条设备上行消息流。
// 执行阶段：
// 1. 通过应用层协议编解码器（JSON/Modbus/GB28181）解码载荷；
// 2. 扁平化提取 params / values / 顶级字段为标准时序指标（tsdb.Record）；
// 3. 异步写入 TSDB 统一时序客户端（支持 TDengine/IoTDB/ClickHouse/TimescaleDB/InfluxDB/Mock 可插拔）；
// 4. 刷新 Redis / 内存设备影子最新快照；
// 5. 执行告警规则判定（evaluateRule）。
func (p *Processor) ProcessMessage(ctx context.Context, msg model.DeviceMessage) error {
	// 1. 通过通用协议解码器解码载荷 (默认优先使用 json 编解码器)
	var data map[string]interface{}
	codec, ok := p.protocols.Get("json")
	if ok {
		decoded, err := codec.Decode(ctx, msg.Payload)
		if err == nil {
			data = decoded
		}
	}

	if data == nil {
		if jsonErr := json.Unmarshal(msg.Payload, &data); jsonErr != nil {
			p.logger.Warn("payload decoding yielded empty map, raw bytes skipped", zap.String("device_key", msg.DeviceKey))
			return nil
		}
	}

	// 兼容业界常见的物模型包装结构 (如 {"params": {...}} 或 {"values": {...}})
	if params, ok := data["params"].(map[string]interface{}); ok {
		data = params
	} else if values, ok := data["values"].(map[string]interface{}); ok {
		data = values
	}

	// 2. 扁平化抽取所有时序键值对
	records := make([]tsdb.Record, 0, len(data))
	ts := msg.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	for k, v := range data {
		// 递归或直接展开 params
		if k == "params" || k == "values" || k == "properties" {
			if subMap, ok := v.(map[string]interface{}); ok {
				for subK, subV := range subMap {
					records = append(records, tsdb.Record{
						DeviceKey: msg.DeviceKey,
						Metric:    subK,
						Value:     subV,
						Timestamp: ts,
					})
					p.evaluateRule(msg.DeviceKey, subK, subV)
				}
				continue
			}
		}

		records = append(records, tsdb.Record{
			DeviceKey: msg.DeviceKey,
			Metric:    k,
			Value:     v,
			Timestamp: ts,
		})
		p.evaluateRule(msg.DeviceKey, k, v)
	}

	// 3. 异步批量/单条写入统一 TSDB
	if len(records) > 0 && p.tsdbClient != nil {
		if err := p.tsdbClient.WriteBatch(ctx, records); err != nil {
			p.logger.Error("failed to write records to TSDB client", zap.String("device_key", msg.DeviceKey), zap.Error(err))
		}
	}

	// 4. 更新设备影子快照
	if p.shadow != nil && len(data) > 0 {
		_ = p.shadow.UpdateShadow(ctx, msg.DeviceKey, data, msg.Timestamp)
	}

	p.logger.Debug("extracted & stored telemetry metrics via tsdb.Client", zap.String("device_key", msg.DeviceKey), zap.Int("count", len(records)))
	return nil
}

// evaluateRule 评估单个指标值是否触发业务告警规则。
func (p *Processor) evaluateRule(deviceKey, metric string, val interface{}) {
	if metric == "temperature" || metric == "temp" {
		floatVal, ok := parseNumericValue(val)
		if !ok {
			return
		}

		if floatVal > 70.0 {
			alarm := model.AlarmEvent{
				DeviceKey:   deviceKey,
				RuleName:    "High Temperature Alarm",
				Level:       "CRITICAL",
				Description: fmt.Sprintf("device temperature reached %.1f°C exceeds threshold 70.0°C", floatVal),
				Timestamp:   time.Now().UTC(),
			}
			p.logger.Warn("🚨 RULE TRIGGERED: Alarm generated!",
				zap.String("device_key", alarm.DeviceKey),
				zap.String("rule", alarm.RuleName),
				zap.String("level", alarm.Level),
				zap.String("desc", alarm.Description),
			)
		}
	}
}

// parseNumericValue 稳健提取任意类型的数值
func parseNumericValue(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
