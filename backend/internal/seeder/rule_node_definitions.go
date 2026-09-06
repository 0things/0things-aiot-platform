package seeder

import (
	"context"

	"aiot-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	inputPort       = `[{"key":"input","label":"输入"}]`
	successOutput   = `[{"key":"success","label":"Success"},{"key":"failure","label":"Failure"}]`
	conditionOutput = `[{"key":"true","label":"True"},{"key":"false","label":"False"},{"key":"failure","label":"Failure"}]`
	alertOutput     = `[{"key":"created","label":"Created"},{"key":"updated","label":"Updated"},{"key":"failure","label":"Failure"}]`
	emptySchema     = `{"type":"object","properties":{}}`
	emptyConfig     = `{}`
)

type ruleNodeDefinitionSeed struct {
	Key           string
	Category      string
	Name          string
	Description   string
	Icon          string
	ConfigSchema  string
	DefaultConfig string
	InputPorts    string
	OutputPorts   string
	ExecutorKey   string
	IsSystem      bool
	SortOrder     int
}

// seedRuleNodeDefinitions keeps operator-managed names and descriptions intact,
// while synchronizing the functional contract used by the graph and executor.
func seedRuleNodeDefinitions(ctx context.Context, db *gorm.DB) error {
	for _, seed := range defaultRuleNodeDefinitionSeeds() {
		definition := model.RuleNodeDefinition{
			UUID:          uuid.NewSHA1(uuid.NameSpaceURL, []byte("0things/rule-node-definition/"+seed.Key)).String(),
			Key:           seed.Key,
			Version:       1,
			Category:      seed.Category,
			Name:          seed.Name,
			Description:   seed.Description,
			Icon:          seed.Icon,
			ConfigSchema:  seed.ConfigSchema,
			DefaultConfig: seed.DefaultConfig,
			InputPorts:    seed.InputPorts,
			OutputPorts:   seed.OutputPorts,
			ExecutorKey:   seed.ExecutorKey,
			Enabled:       true,
			IsSystem:      seed.IsSystem,
			SortOrder:     seed.SortOrder,
		}

		var existing model.RuleNodeDefinition
		err := db.WithContext(ctx).Where("key = ?", seed.Key).First(&existing).Error
		if err == nil {
			if err := db.WithContext(ctx).Model(&existing).Updates(map[string]any{
				"version":        definition.Version,
				"category":       definition.Category,
				"config_schema":  definition.ConfigSchema,
				"default_config": definition.DefaultConfig,
				"input_ports":    definition.InputPorts,
				"output_ports":   definition.OutputPorts,
				"executor_key":   definition.ExecutorKey,
				"enabled":        definition.Enabled,
				"is_system":      definition.IsSystem,
				"sort_order":     definition.SortOrder,
			}).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err := db.WithContext(ctx).Create(&definition).Error; err != nil {
			return err
		}
	}
	return nil
}

func defaultRuleNodeDefinitionSeeds() []ruleNodeDefinitionSeed {
	definitions := []ruleNodeDefinitionSeed{
		{Key: "system.entry", Category: "System", Name: "消息入口", Description: "规则链的系统消息入口", Icon: "Play", ConfigSchema: emptySchema, DefaultConfig: emptyConfig, InputPorts: `[]`, OutputPorts: `[{"key":"success","label":"Success"}]`, ExecutorKey: "system-entry-v1", IsSystem: true},

		{Key: "filter.message-type", Category: "Filter", Name: "消息类型", Description: "按规则消息类型分流", Icon: "ListFilter", ConfigSchema: `{"type":"object","required":["messageTypes"],"properties":{"messageTypes":{"type":"array","minItems":1,"items":{"type":"string"},"x-ui":"multi-select"}}}`, DefaultConfig: `{"messageTypes":[]}`, InputPorts: inputPort, OutputPorts: conditionOutput, ExecutorKey: "filter-message-type-v1", SortOrder: 10},
		{Key: "filter.originator-type", Category: "Filter", Name: "来源实体", Description: "按消息来源实体类型筛选", Icon: "GitFork", ConfigSchema: `{"type":"object","properties":{"originatorTypes":{"type":"array"}}}`, DefaultConfig: `{"originatorTypes":["device"]}`, InputPorts: inputPort, OutputPorts: conditionOutput, ExecutorKey: "filter-originator-type-v1", SortOrder: 20},
		{Key: "filter.property", Category: "Filter", Name: "属性条件", Description: "按属性值判断消息是否通过", Icon: "Filter", ConfigSchema: `{"type":"object","required":["property","operator"],"properties":{"property":{"type":"string","minLength":1,"x-ui":"property-picker"},"operator":{"type":"string","enum":["eq","ne","gt","gte","lt","lte","contains","exists"],"x-ui":"select"},"value":{"x-ui":"json-value"}}}`, DefaultConfig: `{"property":"","operator":"eq","value":""}`, InputPorts: inputPort, OutputPorts: conditionOutput, ExecutorKey: "filter-property-v1", SortOrder: 30},
		{Key: "filter.value-condition", Category: "Filter", Name: "数值条件", Description: "按数值阈值决定 True 或 False 分支", Icon: "GitBranch", ConfigSchema: `{"type":"object","required":["path","operator","value"],"properties":{"path":{"type":"string","minLength":1,"x-ui":"property-picker"},"operator":{"type":"string","enum":["gt","gte","lt","lte","eq","ne"],"x-ui":"select"},"value":{"type":"number"}}}`, DefaultConfig: `{"path":"","operator":"gt","value":0}`, InputPorts: inputPort, OutputPorts: conditionOutput, ExecutorKey: "filter-value-condition-v1", SortOrder: 40},

		{Key: "enrichment.originator-attributes", Category: "Enrichment", Name: "来源属性", Description: "读取来源实体属性并补充到消息上下文", Icon: "Database", ConfigSchema: `{"type":"object","properties":{"keys":{"type":"array"}}}`, DefaultConfig: `{"keys":[]}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "enrichment-originator-attributes-v1", SortOrder: 10},
		{Key: "enrichment.related-entity", Category: "Enrichment", Name: "关联实体", Description: "按关系查询关联实体并补充上下文", Icon: "Network", ConfigSchema: `{"type":"object","properties":{"relationType":{"type":"string"}}}`, DefaultConfig: `{"relationType":""}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "enrichment-related-entity-v1", SortOrder: 20},
		{Key: "enrichment.latest-telemetry", Category: "Enrichment", Name: "最新遥测", Description: "读取来源实体的最新遥测", Icon: "ChartLine", ConfigSchema: `{"type":"object","properties":{"keys":{"type":"array"}}}`, DefaultConfig: `{"keys":[]}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "enrichment-latest-telemetry-v1", SortOrder: 30},

		{Key: "transformation.map", Category: "Transformation", Name: "字段映射", Description: "映射消息字段和元数据", Icon: "Waypoints", ConfigSchema: `{"type":"object","required":["mappings"],"properties":{"mappings":{"type":"array","items":{"type":"object","required":["from","to"],"properties":{"from":{"type":"string","minLength":1,"x-ui":"path-picker"},"to":{"type":"string","minLength":1,"x-ui":"path-picker"}}},"x-ui":"mapping-list"}}}`, DefaultConfig: `{"mappings":[]}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "transformation-map-v1", SortOrder: 10},
		{Key: "transformation.change-originator", Category: "Transformation", Name: "变更来源实体", Description: "把消息来源切换为关联实体", Icon: "Repeat2", ConfigSchema: `{"type":"object","properties":{"target":{"type":"string"}}}`, DefaultConfig: `{"target":""}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "transformation-change-originator-v1", SortOrder: 20},

		{Key: "action.save-attributes", Category: "Action", Name: "保存属性", Description: "将字段保存为实体属性", Icon: "Save", ConfigSchema: `{"type":"object","required":["scope","attributes"],"properties":{"scope":{"type":"string","enum":["server","shared","client"],"x-ui":"select"},"attributes":{"type":"object","minProperties":1,"additionalProperties":true,"x-ui":"key-value-editor"}}}`, DefaultConfig: `{"scope":"server","attributes":{}}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "action-save-attributes-v1", SortOrder: 10},
		{Key: "action.save-telemetry", Category: "Action", Name: "保存遥测", Description: "将字段保存为遥测数据", Icon: "Activity", ConfigSchema: `{"type":"object","required":["telemetry"],"properties":{"telemetry":{"type":"object","minProperties":1,"additionalProperties":true,"x-ui":"key-value-editor"}}}`, DefaultConfig: `{"telemetry":{}}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "action-save-telemetry-v1", SortOrder: 20},
		{Key: "action.create-alert", Category: "Action", Name: "创建告警", Description: "创建或聚合活动告警", Icon: "Siren", ConfigSchema: `{"type":"object","properties":{"alertType":{"type":"string"},"severity":{"type":"string"}}}`, DefaultConfig: `{"alertType":"","severity":"warning"}`, InputPorts: inputPort, OutputPorts: alertOutput, ExecutorKey: "action-create-alert-v1", SortOrder: 30},
		{Key: "action.clear-alert", Category: "Action", Name: "清除告警", Description: "恢复匹配的活动告警", Icon: "CircleCheck", ConfigSchema: `{"type":"object","properties":{"alertType":{"type":"string"}}}`, DefaultConfig: `{"alertType":""}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "action-clear-alert-v1", SortOrder: 40},
		{Key: "action.device-command", Category: "Action", Name: "设备命令", Description: "向设备下发服务调用或 RPC", Icon: "Send", ConfigSchema: `{"type":"object","properties":{"command":{"type":"string"},"params":{"type":"object"}}}`, DefaultConfig: `{"command":"","params":{}}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "action-device-command-v1", SortOrder: 50},
		{Key: "action.rpc-reply", Category: "Action", Name: "RPC 回复", Description: "回复设备发起的 RPC 请求", Icon: "Reply", ConfigSchema: `{"type":"object","properties":{"response":{"type":"object"}}}`, DefaultConfig: `{"response":{}}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "action-rpc-reply-v1", SortOrder: 60},

		{Key: "external.webhook", Category: "External", Name: "调用 Webhook", Description: "向 HTTPS 公网地址发送 HTTP 请求", Icon: "Webhook", ConfigSchema: `{"type":"object","required":["url","method","timeoutMs","retryCount"],"properties":{"url":{"type":"string","format":"uri","pattern":"^https://","x-ui":"url"},"method":{"type":"string","enum":["POST","PUT","PATCH"],"x-ui":"select"},"headers":{"type":"object","additionalProperties":{"type":"string"},"x-ui":"key-value-editor"},"bodyTemplate":{"type":"object","additionalProperties":true,"x-ui":"json"},"timeoutMs":{"type":"integer","minimum":100,"maximum":30000,"x-ui":"number"},"retryCount":{"type":"integer","minimum":0,"maximum":5,"x-ui":"number"}}}`, DefaultConfig: `{"url":"","method":"POST","headers":{},"bodyTemplate":{},"timeoutMs":5000,"retryCount":3}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "external-webhook-v1", SortOrder: 10},

		{Key: "flow.delay", Category: "Flow", Name: "延迟", Description: "延迟后继续执行", Icon: "Clock3", ConfigSchema: `{"type":"object","required":["durationSeconds"],"properties":{"durationSeconds":{"type":"integer","minimum":1,"maximum":604800,"x-ui":"duration-seconds"}}}`, DefaultConfig: `{"durationSeconds":60}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "flow-delay-v1", SortOrder: 10},
		{Key: "flow.rate-limit", Category: "Flow", Name: "限流", Description: "限制单位时间内的消息数量", Icon: "Gauge", ConfigSchema: `{"type":"object","properties":{"maxCount":{"type":"number"},"windowSeconds":{"type":"number"}}}`, DefaultConfig: `{"maxCount":1,"windowSeconds":60}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "flow-rate-limit-v1", SortOrder: 20},
		{Key: "flow.deduplicate", Category: "Flow", Name: "去重", Description: "按指定键忽略重复消息", Icon: "CopyX", ConfigSchema: `{"type":"object","properties":{"keyTemplate":{"type":"string"},"windowSeconds":{"type":"number"}}}`, DefaultConfig: `{"keyTemplate":"{{message.id}}","windowSeconds":60}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "flow-deduplicate-v1", SortOrder: 30},
		{Key: "flow.message-count", Category: "Flow", Name: "消息计数", Description: "达到指定消息数量后继续", Icon: "Hash", ConfigSchema: `{"type":"object","properties":{"threshold":{"type":"number"}}}`, DefaultConfig: `{"threshold":1}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "flow-message-count-v1", SortOrder: 40},
		{Key: "flow.sub-chain", Category: "Flow", Name: "子规则链", Description: "调用同组织的已发布规则链", Icon: "Workflow", ConfigSchema: `{"type":"object","properties":{"ruleChainUuid":{"type":"string"}}}`, DefaultConfig: `{"ruleChainUuid":""}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "flow-sub-chain-v1", SortOrder: 50},
		{Key: "flow.end", Category: "Flow", Name: "结束", Description: "结束当前消息的规则执行", Icon: "CircleStop", ConfigSchema: emptySchema, DefaultConfig: emptyConfig, InputPorts: inputPort, OutputPorts: `[]`, ExecutorKey: "flow-end-v1", SortOrder: 60},

		{Key: "analytics.aggregate", Category: "Analytics", Name: "时间窗口聚合", Description: "对窗口内数值执行聚合", Icon: "ChartNoAxesCombined", ConfigSchema: `{"type":"object","properties":{"path":{"type":"string"},"operation":{"type":"string"},"windowSeconds":{"type":"number"}}}`, DefaultConfig: `{"path":"","operation":"avg","windowSeconds":300}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "analytics-aggregate-v1", SortOrder: 10},
		{Key: "analytics.calculate", Category: "Analytics", Name: "聚合值计算", Description: "对已聚合字段执行受限算术计算", Icon: "Calculator", ConfigSchema: `{"type":"object","properties":{"leftPath":{"type":"string"},"operator":{"type":"string"},"rightValue":{"type":"number"}}}`, DefaultConfig: `{"leftPath":"","operator":"add","rightValue":0}`, InputPorts: inputPort, OutputPorts: successOutput, ExecutorKey: "analytics-calculate-v1", SortOrder: 20},
	}

	return append(definitions, thingsBoardReferenceSeeds()...)
}

const thingsBoardReferenceSchema = `{"type":"object","additionalProperties":true,"x-ui":"json"}`

type thingsBoardNodeSeed struct {
	Key  string
	Name string
}

// thingsBoardReferenceSeeds mirrors every current open-source ThingsBoard
// @RuleNode descriptor. Their configuration stays as a JSON editor until the
// matching Go executor supplies a stricter schema.
func thingsBoardReferenceSeeds() []ruleNodeDefinitionSeed {
	groups := []struct {
		category string
		nodes    []thingsBoardNodeSeed
	}{
		{category: "Filter", nodes: []thingsBoardNodeSeed{
			{Key: "tb.filter.alarm-status-filter", Name: "告警状态筛选"},
			{Key: "tb.filter.asset-profile-switch", Name: "资产配置切换"},
			{Key: "tb.filter.check-fields-presence", Name: "检查字段存在"},
			{Key: "tb.filter.check-relation-presence", Name: "检查关系存在"},
			{Key: "tb.filter.device-profile-switch", Name: "设备配置切换"},
			{Key: "tb.filter.entity-type-filter", Name: "实体类型筛选"},
			{Key: "tb.filter.entity-type-switch", Name: "实体类型切换"},
			{Key: "tb.filter.gps-geofencing-filter", Name: "GPS 地理围栏筛选"},
			{Key: "tb.filter.message-type-filter", Name: "消息类型筛选"},
			{Key: "tb.filter.message-type-switch", Name: "消息类型切换"},
			{Key: "tb.filter.script", Name: "脚本筛选"},
			{Key: "tb.filter.switch", Name: "条件切换"},
		}},
		{category: "Enrichment", nodes: []thingsBoardNodeSeed{
			{Key: "tb.enrichment.calculate-delta", Name: "计算增量"},
			{Key: "tb.enrichment.customer-attributes", Name: "客户属性"},
			{Key: "tb.enrichment.customer-details", Name: "客户详情"},
			{Key: "tb.enrichment.fetch-device-credentials", Name: "获取设备凭据"},
			{Key: "tb.enrichment.originator-attributes", Name: "来源实体属性"},
			{Key: "tb.enrichment.originator-fields", Name: "来源实体字段"},
			{Key: "tb.enrichment.originator-telemetry", Name: "来源实体遥测"},
			{Key: "tb.enrichment.related-device-attributes", Name: "关联设备属性"},
			{Key: "tb.enrichment.related-entity-data", Name: "关联实体数据"},
			{Key: "tb.enrichment.tenant-attributes", Name: "租户属性"},
			{Key: "tb.enrichment.tenant-details", Name: "租户详情"},
		}},
		{category: "Transformation", nodes: []thingsBoardNodeSeed{
			{Key: "tb.transformation.change-originator", Name: "变更来源实体"},
			{Key: "tb.transformation.copy-key-value-pairs", Name: "复制键值对"},
			{Key: "tb.transformation.deduplication", Name: "消息去重"},
			{Key: "tb.transformation.delete-key-value-pairs", Name: "删除键值对"},
			{Key: "tb.transformation.json-path", Name: "JSON Path"},
			{Key: "tb.transformation.rename-keys", Name: "重命名键"},
			{Key: "tb.transformation.script", Name: "脚本转换"},
			{Key: "tb.transformation.split-array-message", Name: "拆分数组消息"},
			{Key: "tb.transformation.to-email", Name: "转换为邮件"},
		}},
		{category: "Action", nodes: []thingsBoardNodeSeed{
			{Key: "tb.action.assign-to-customer", Name: "分配给客户"},
			{Key: "tb.action.calculated-fields-and-alarm-rules", Name: "计算字段和告警规则"},
			{Key: "tb.action.clear-alarm", Name: "清除告警"},
			{Key: "tb.action.copy-to-view", Name: "复制到实体视图"},
			{Key: "tb.action.create-alarm", Name: "创建告警"},
			{Key: "tb.action.create-relation", Name: "创建关系"},
			{Key: "tb.action.delay-deprecated", Name: "延迟（已弃用）"},
			{Key: "tb.action.delete-attributes", Name: "删除属性"},
			{Key: "tb.action.delete-relation", Name: "删除关系"},
			{Key: "tb.action.device-profile-deprecated", Name: "设备配置（已弃用）"},
			{Key: "tb.action.device-state", Name: "设备状态"},
			{Key: "tb.action.generator", Name: "消息生成器"},
			{Key: "tb.action.gps-geofencing-events", Name: "GPS 地理围栏事件"},
			{Key: "tb.action.log", Name: "日志"},
			{Key: "tb.action.math-function", Name: "数学函数"},
			{Key: "tb.action.message-count", Name: "消息计数"},
			{Key: "tb.action.push-to-cloud", Name: "推送到云端"},
			{Key: "tb.action.push-to-edge", Name: "推送到边缘"},
			{Key: "tb.action.rest-call-reply", Name: "REST 调用回复"},
			{Key: "tb.action.rpc-call-reply", Name: "RPC 调用回复"},
			{Key: "tb.action.rpc-call-request", Name: "RPC 调用请求"},
			{Key: "tb.action.save-attributes", Name: "保存属性"},
			{Key: "tb.action.save-time-series", Name: "保存时序数据"},
			{Key: "tb.action.save-to-custom-table", Name: "保存到自定义表"},
			{Key: "tb.action.synchronization-end", Name: "同步结束"},
			{Key: "tb.action.synchronization-start", Name: "同步开始"},
			{Key: "tb.action.unassign-from-customer", Name: "取消客户分配"},
		}},
		{category: "External", nodes: []thingsBoardNodeSeed{
			{Key: "tb.external.ai-request", Name: "AI 请求"},
			{Key: "tb.external.aws-lambda", Name: "AWS Lambda"},
			{Key: "tb.external.aws-sns", Name: "AWS SNS"},
			{Key: "tb.external.aws-sqs", Name: "AWS SQS"},
			{Key: "tb.external.azure-iot-hub", Name: "Azure IoT Hub"},
			{Key: "tb.external.gcp-pubsub", Name: "GCP Pub/Sub"},
			{Key: "tb.external.kafka", Name: "Kafka"},
			{Key: "tb.external.mqtt", Name: "MQTT"},
			{Key: "tb.external.rabbitmq", Name: "RabbitMQ"},
			{Key: "tb.external.rest-api-call", Name: "REST API 调用"},
			{Key: "tb.external.send-email", Name: "发送邮件"},
			{Key: "tb.external.send-notification", Name: "发送通知"},
			{Key: "tb.external.send-sms", Name: "发送短信"},
			{Key: "tb.external.send-to-slack", Name: "发送到 Slack"},
		}},
		{category: "Flow", nodes: []thingsBoardNodeSeed{
			{Key: "tb.flow.acknowledge", Name: "确认消息"},
			{Key: "tb.flow.checkpoint", Name: "检查点"},
			{Key: "tb.flow.output", Name: "规则链输出"},
			{Key: "tb.flow.rule-chain", Name: "规则链"},
		}},
	}

	var seeds []ruleNodeDefinitionSeed
	for _, group := range groups {
		for index, node := range group.nodes {
			outputPorts := successOutput
			if group.category == "Filter" {
				outputPorts = conditionOutput
			}

			seeds = append(seeds, ruleNodeDefinitionSeed{
				Key:           node.Key,
				Category:      group.category,
				Name:          node.Name,
				Description:   "ThingsBoard 对应节点：" + node.Key,
				Icon:          "Blocks",
				ConfigSchema:  thingsBoardReferenceSchema,
				DefaultConfig: emptyConfig,
				InputPorts:    inputPort,
				OutputPorts:   outputPorts,
				ExecutorKey:   "pending-" + node.Key + "-v1",
				SortOrder:     1000 + index,
			})
		}
	}
	return seeds
}
