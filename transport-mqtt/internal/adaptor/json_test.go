package adaptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonMqttAdaptor_ConvertToTelemetryPayload(t *testing.T) {
	adaptor := NewJsonMqttAdaptor()

	// Standard物模型属性上报格式
	rawJSON := []byte(`{
		"version": "1.0",
		"sys": {
			"ack": 0
		},
		"params": {
			"Power": {
				"value": "on",
				"time": 1524448722000
			},
			"WF": {
				"value": 23.6,
				"time": 1524448722000
			}
		}
	}`)

	res, err := adaptor.ConvertToTelemetryPayload(rawJSON)
	require.NoError(t, err)
	require.Len(t, res, 1)

	assert.Equal(t, int64(1524448722000), res[0].Timestamp)
	assert.Equal(t, "on", res[0].Values["Power"])
	assert.Equal(t, 23.6, res[0].Values["WF"])

	// 异常空参数
	_, err = adaptor.ConvertToTelemetryPayload([]byte(`{"params":{}}`))
	require.Error(t, err)

	// 异常：time 缺失或 <= 0
	_, err = adaptor.ConvertToTelemetryPayload([]byte(`{
		"params": {
			"Power": { "value": "on", "time": 0 }
		}
	}`))
	require.Error(t, err)

}
