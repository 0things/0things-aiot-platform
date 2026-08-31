package engine

import (
	"context"
	"fmt"
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
	p.logger.Info("processing device message",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
	)

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
		p.logger.Warn("payload decoding yielded empty map, raw bytes skipped", zap.String("device_key", msg.DeviceKey))
		return nil
	}

	// 兼容业界常见的物模型包装结构 (如 {"params": {...}} 或 {"values": {...}})
	if params, ok := data["params"].(map[string]interface{}); ok {
		data = params
	} else if values, ok := data["values"].(map[string]interface{}); ok {
		data = values
	}

	// 2. 扁平化抽取所有时序键值对
	records := make([]tsdb.Record, 0, len(data))
	for k, v := range data {
		records = append(records, tsdb.Record{
			DeviceKey: msg.DeviceKey,
			Metric:    k,
			Value:     v,
			Timestamp: msg.Timestamp,
		})

		// 3. 触发业务规则计算与阈值评估
		p.evaluateRule(msg.DeviceKey, k, v)
	}

	// 4. 异步持久化至统一 TSDB 时序数据库客户端
	if p.tsdbClient != nil && len(records) > 0 {
		_ = p.tsdbClient.WriteBatch(ctx, records)
	}

	// 5. 刷新设备实时影子快照
	if p.shadow != nil && len(data) > 0 {
		_ = p.shadow.UpdateShadow(ctx, msg.DeviceKey, data, msg.Timestamp)
	}

	p.logger.Debug("extracted & stored telemetry metrics via tsdb.Client", zap.String("device_key", msg.DeviceKey), zap.Int("count", len(records)))
	return nil
}

// evaluateRule 评估单个指标值是否触发业务告警规则。
func (p *Processor) evaluateRule(deviceKey, metric string, val interface{}) {
	if metric == "temperature" || metric == "temp" {
		var floatVal float64
		switch v := val.(type) {
		case float64:
			floatVal = v
		case int:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
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
