package v1

type RuleNodeDefinition struct {
	UUID          string `json:"uuid"`
	Key           string `json:"key"`
	Category      string `json:"category"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Icon          string `json:"icon"`
	ConfigSchema  string `json:"configSchema"`
	DefaultConfig string `json:"defaultConfig"`
	InputPorts    string `json:"inputPorts"`
	OutputPorts   string `json:"outputPorts"`
	ExecutorKey   string `json:"executorKey"`
	IsSystem      bool   `json:"isSystem"`
}

type ListRuleNodeDefinitionsResponse struct {
	Items []RuleNodeDefinition `json:"items"`
}
