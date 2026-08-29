package job

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// OTADeviceUpgradeReportEvent 定义经桥接后投递到 Kafka 的标准化设备升级上报事件结构体。
type OTADeviceUpgradeReportEvent struct {
	BatchID   string `json:"batch_id"`
	DeviceKey string `json:"device_key"`
	Status    string `json:"status"`
	Version   string `json:"version"`
	Progress  int32  `json:"progress"`
	Desc      string `json:"desc"`
}

// OTAMQTTReportBridge 负责监听设备端 MQTT 升级进度和版本上报主题，
// 并将其标准化清洗后桥接转发至 Kafka 回报队列，实现协议接入层与状态处理层的解耦。
type OTAMQTTReportBridge struct {
	mqtt       service.MQTTServiceInterface
	kafka      service.KafkaServiceInterface
	logger     *log.Logger
	subscribed []string
}

func NewOTAMQTTReportBridge(mqtt service.MQTTServiceInterface, kafka service.KafkaServiceInterface, logger *log.Logger) *OTAMQTTReportBridge {
	return &OTAMQTTReportBridge{
		mqtt:   mqtt,
		kafka:  kafka,
		logger: logger,
	}
}

func (c *OTAMQTTReportBridge) Start(ctx context.Context) error {
	if c.mqtt == nil || c.kafka == nil {
		c.logger.Warn("OTA MQTT report bridge disabled: dependency unavailable")
		return nil
	}

	handler := func(_ mqtt.Client, msg mqtt.Message) {
		var report struct {
			Params struct {
				BatchID   string          `json:"batch_id"`
				DeviceKey string          `json:"device_key"`
				Status    string          `json:"status"`
				Step      json.RawMessage `json:"step"`
				Version   string          `json:"version"`
				Desc      string          `json:"desc"`
			} `json:"params"`
		}

		if err := json.Unmarshal(msg.Payload(), &report); err != nil {
			c.logger.Error("invalid OTA MQTT report payload", zap.Error(err), zap.ByteString("payload", msg.Payload()))
			return
		}

		if report.Params.BatchID == "" || report.Params.DeviceKey == "" {
			c.logger.Warn("OTA MQTT report missing required batch_id or device_key")
			return
		}

		status := report.Params.Status
		if status == "" {
			status = enum.OTAStatusInProgress
		}

		// 兼容设备端 step 字段可能以数字 (85) 或字符串 ("85") 上报的格式
		var step int32
		if err := json.Unmarshal(report.Params.Step, &step); err != nil {
			var text string
			if json.Unmarshal(report.Params.Step, &text) == nil {
				value, _ := strconv.ParseInt(text, 10, 32)
				step = int32(value)
			}
		}
		// 阿里云 OTA 约定负数 step 表示下载、校验、烧写或重启失败，统一映射为 failed。
		if step < 0 {
			status = enum.OTAStatusFailed
			if report.Params.Desc == "" {
				report.Params.Desc = "OTA device reported failure at step " + strconv.FormatInt(int64(step), 10)
			}
		}

		event := OTADeviceUpgradeReportEvent{
			BatchID:   report.Params.BatchID,
			DeviceKey: report.Params.DeviceKey,
			Status:    status,
			Version:   report.Params.Version,
			Progress:  step,
			Desc:      report.Params.Desc,
		}

		// 使用独立超时的 Context 进行 Kafka 生产，防范服务关闭过程中的并发 Context 取消导致在途上报丢失
		produceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := event.BatchID + ":" + event.DeviceKey
		if err := c.kafka.ProduceJSON(produceCtx, enum.KafkaTopicOTAUpgradeReportV1, key, event); err != nil {
			c.logger.Error("failed to publish OTA MQTT report to Kafka",
				zap.Error(err),
				zap.String("batch_id", event.BatchID),
				zap.String("device_key", event.DeviceKey),
				zap.String("desc", event.Desc),
			)
		}
	}

	topics := []string{enum.MQTTTopicOTADeviceProgress, enum.MQTTTopicOTADeviceInform}
	c.subscribed = topics

	for _, topic := range topics {
		if err := c.mqtt.Subscribe(ctx, topic, 1, handler); err != nil {
			c.logger.Warn("failed to subscribe OTA MQTT topic", zap.String("topic", topic), zap.Error(err))
		}
	}
	return nil
}

// Stop 优雅注销所订阅的 MQTT 主题。
func (c *OTAMQTTReportBridge) Stop() {
	if c.mqtt != nil && len(c.subscribed) > 0 {
		_ = c.mqtt.Unsubscribe(c.subscribed...)
	}
}
