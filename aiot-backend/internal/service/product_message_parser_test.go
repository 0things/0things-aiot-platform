package service

import (
	"strings"
	"testing"

	messageParserV1 "0things-backend/api/message_parser/v1"
	"github.com/stretchr/testify/require"
)

const executableMessageParserScript = `
function transformPayload(topic, rawData) { return { topic: topic, first: rawData[0] }; }
function rawDataToProtocol(rawData) { return { first: rawData[0], length: rawData.length }; }
function protocolToRawData(jsonObj) { return [jsonObj.value, 2]; }
`

func TestExecuteMessageParser(t *testing.T) {
	t.Run("device report", func(t *testing.T) {
		result, err := executeMessageParser(executableMessageParserScript, messageParserV1.ExecuteProductMessageParserRequest{
			Mode:    "device_report",
			RawData: "0x0A0B",
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"first":10,"length":2}`, result.JSONOutput)
	})

	t.Run("device receive", func(t *testing.T) {
		result, err := executeMessageParser(executableMessageParserScript, messageParserV1.ExecuteProductMessageParserRequest{
			Mode:    "device_receive",
			RawData: "0A02",
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"first":10,"length":2}`, result.JSONOutput)
	})

	t.Run("custom", func(t *testing.T) {
		result, err := executeMessageParser(executableMessageParserScript, messageParserV1.ExecuteProductMessageParserRequest{
			Mode:    "custom",
			Topic:   "/devices/demo/custom",
			RawData: "01",
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"topic":"/devices/demo/custom","first":1}`, result.JSONOutput)
	})
}

func TestMessageParserValidation(t *testing.T) {
	require.Error(t, validateMessageParserScript("function rawDataToProtocol() {}"))

	_, err := executeMessageParser(executableMessageParserScript, messageParserV1.ExecuteProductMessageParserRequest{
		Mode:    "device_report",
		RawData: "not hex",
	})
	require.Error(t, err)

	_, err = executeMessageParser(executableMessageParserScript, messageParserV1.ExecuteProductMessageParserRequest{
		Mode:    "custom",
		RawData: "00",
	})
	require.Error(t, err)

	_, err = executeMessageParser(`
function transformPayload(topic, rawData) { return {}; }
function rawDataToProtocol(rawData) { return {}; }
function protocolToRawData(jsonObj) { return []; }
while (true) {}
`, messageParserV1.ExecuteProductMessageParserRequest{Mode: "device_report", RawData: "00"})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "timed out"))
}
