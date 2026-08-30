package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rule-engine/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Processor 是规则引擎的核心计算处理器，负责 TSL 物模型字段展开、时序指标抽取与告警阈值计算。
type Processor struct {
	logger *zap.Logger
}

func NewProcessor(config *viper.Viper, logger *zap.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

// ProcessMessage 处理单条设备上行消息流。
// 执行阶段：
// 1. 反序列化 JSON 载荷；
// 2. 扁平化提取 params / values / 顶级字段为标准时序指标（TelemetryRecord）；
// 3. 执行告警规则判定（evaluateRule）。
func (p *Processor) ProcessMessage(ctx context.Context, msg model.DeviceMessage) error {
	p.logger.Info("processing device message",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
	)

	// 1. 解析原始 Payload JSON
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payloadMap); err != nil {
		p.logger.Warn("payload is not standard json object, raw bytes stored", zap.String("device_key", msg.DeviceKey))
		return nil
	}

	// 兼容业界常见的物模型包装结构 (如 {"params": {...}} 或 {"values": {...}})
	data := payloadMap
	if params, ok := payloadMap["params"].(map[string]interface{}); ok {
		data = params
	} else if values, ok := payloadMap["values"].(map[string]interface{}); ok {
		data = values
	}

	// 2. 扁平化抽取所有时序键值对
	records := make([]model.TelemetryRecord, 0, len(data))
	for k, v := range data {
		records = append(records, model.TelemetryRecord{
			DeviceKey: msg.DeviceKey,
			Metric:    k,
			Value:     v,
			Timestamp: msg.Timestamp,
		})

		// 3. 触发业务规则计算与阈值评估
		p.evaluateRule(msg.DeviceKey, k, v)
	}

	p.logger.Debug("extracted telemetry metrics", zap.String("device_key", msg.DeviceKey), zap.Int("count", len(records)))
	return nil
}

// evaluateRule 评估单个指标值是否触发业务告警规则。
// 示例业务规则：当设备上报的温度 metric（temperature 或 temp）数值超过 70.0℃ 时，自动生成 CRITICAL 严重告警。
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
