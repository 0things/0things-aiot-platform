package javascript

import (
	"context"
	"errors"
	"fmt"
	"time"

	"0things/pkg/protocol"

	"github.com/dop251/goja"
)

var _ protocol.Codec = (*Codec)(nil)

// Codec executes user-defined JavaScript data converter scripts for uplink/downlink message parsing.
type Codec struct {
	program *goja.Program
}

// NewCodec compiles the given JavaScript script into a reusable executable program.
func NewCodec(script string) (*Codec, error) {
	if script == "" {
		return nil, errors.New("script cannot be empty")
	}
	prog, err := goja.Compile("converter.js", script, true)
	if err != nil {
		return nil, fmt.Errorf("compile javascript script: %w", err)
	}
	return &Codec{program: prog}, nil
}

// Name returns the codec identifier name.
func (c *Codec) Name() string {
	return "javascript"
}

// Decode executes decodeUplink(bytes) in a sandboxed JavaScript runtime with timeout protection.
func (c *Codec) Decode(ctx context.Context, payload []byte) (map[string]any, error) {
	vm := goja.New()

	// Timeout protection (50ms) to avoid infinite loops in user scripts
	timer := time.AfterFunc(50*time.Millisecond, func() {
		vm.Interrupt("execution timeout (50ms exceeded)")
	})
	defer timer.Stop()

	if _, err := vm.RunProgram(c.program); err != nil {
		return nil, fmt.Errorf("initialize script program: %w", err)
	}

	decodeFn, ok := goja.AssertFunction(vm.Get("decodeUplink"))
	if !ok {
		return nil, errors.New("function decodeUplink(bytes) not found in script")
	}

	// Pass raw bytes as integer array to JavaScript: [0x01, 0x03, ...]
	byteArgs := make([]any, len(payload))
	for i, b := range payload {
		byteArgs[i] = int(b)
	}
	jsBytes := vm.ToValue(byteArgs)

	result, err := decodeFn(goja.Undefined(), jsBytes)
	if err != nil {
		return nil, fmt.Errorf("execute decodeUplink: %w", err)
	}

	exported := result.Export()
	if exported == nil {
		return nil, errors.New("decodeUplink returned nil")
	}

	if outMap, ok := exported.(map[string]any); ok {
		return outMap, nil
	}
	if genericMap, ok := exported.(map[string]interface{}); ok {
		return genericMap, nil
	}

	return nil, fmt.Errorf("decodeUplink must return an object, got %T", exported)
}

// Encode executes optional encodeDownlink(object) function in JavaScript to convert object back to raw bytes.
func (c *Codec) Encode(ctx context.Context, value map[string]any) ([]byte, error) {
	vm := goja.New()

	timer := time.AfterFunc(50*time.Millisecond, func() {
		vm.Interrupt("execution timeout (50ms exceeded)")
	})
	defer timer.Stop()

	if _, err := vm.RunProgram(c.program); err != nil {
		return nil, fmt.Errorf("initialize script program: %w", err)
	}

	encodeFn, ok := goja.AssertFunction(vm.Get("encodeDownlink"))
	if !ok {
		return nil, errors.New("function encodeDownlink(object) not found in script")
	}

	result, err := encodeFn(goja.Undefined(), vm.ToValue(value))
	if err != nil {
		return nil, fmt.Errorf("execute encodeDownlink: %w", err)
	}

	if bytes, ok := result.Export().([]byte); ok {
		return bytes, nil
	}
	if arr, ok := result.Export().([]any); ok {
		out := make([]byte, len(arr))
		for i, item := range arr {
			if num, ok := item.(int64); ok {
				out[i] = byte(num)
			} else if num, ok := item.(float64); ok {
				out[i] = byte(num)
			} else if num, ok := item.(int); ok {
				out[i] = byte(num)
			}
		}
		return out, nil
	}

	return nil, fmt.Errorf("encodeDownlink must return byte array, got %T", result.Export())
}
