// Package messageparserv1 owns the product message parser HTTP contract.
package messageparserv1

const LanguageJavaScriptES5 = "javascript-es5"

type ProductMessageParser struct {
	ProductKey string `json:"productKey"`
	Language   string `json:"language"`
	Script     string `json:"script"`
	IsDefault  bool   `json:"isDefault"`
} //@name MessageParserProductMessageParser

type UpsertProductMessageParserRequest struct {
	Language string `json:"language" binding:"required"`
	Script   string `json:"script" binding:"required"`
} //@name MessageParserUpsertProductMessageParserRequest

type ExecuteProductMessageParserRequest struct {
	Mode    string `json:"mode" binding:"required,oneof=device_report device_receive custom"`
	Topic   string `json:"topic"`
	RawData string `json:"rawData"`
	JSONObj string `json:"jsonObj"`
} //@name MessageParserExecuteProductMessageParserRequest

type ExecuteProductMessageParserResponse struct {
	Mode       string `json:"mode"`
	JSONOutput string `json:"jsonOutput,omitempty"`
	RawData    string `json:"rawData,omitempty"`
} //@name MessageParserExecuteProductMessageParserResponse
