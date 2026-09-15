package javascript

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJavaScriptCodec_Decode(t *testing.T) {
	script := `
function decodeUplink(bytes) {
    var rawTemp = (bytes[0] << 8) | bytes[1];
    return {
        "temperature": rawTemp / 10.0,
        "humidity": bytes[2]
    };
}
`
	codec, err := NewCodec(script)
	require.NoError(t, err)
	assert.Equal(t, "javascript", codec.Name())

	// Test decoding raw bytes: [0x00, 0xFF, 0x32] -> temperature: 25.5, humidity: 50
	payload := []byte{0x00, 0xFF, 0x32}
	result, err := codec.Decode(context.Background(), payload)
	require.NoError(t, err)
	assert.Equal(t, 25.5, result["temperature"])
	assert.Equal(t, int64(50), result["humidity"])
}

func TestJavaScriptCodec_Encode(t *testing.T) {
	script := `
function decodeUplink(bytes) { return {}; }
function encodeDownlink(obj) {
    return [0x01, 0x05, obj.switch ? 0xFF : 0x00];
}
`
	codec, err := NewCodec(script)
	require.NoError(t, err)

	bytes, err := codec.Encode(context.Background(), map[string]any{"switch": true})
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x05, 0xFF}, bytes)
}

func TestJavaScriptCodec_TimeoutProtection(t *testing.T) {
	infiniteLoopScript := `
function decodeUplink(bytes) {
    while(true) {}
}
`
	codec, err := NewCodec(infiniteLoopScript)
	require.NoError(t, err)

	_, err = codec.Decode(context.Background(), []byte{0x01})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}
