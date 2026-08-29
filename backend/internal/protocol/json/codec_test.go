package jsoncodec

import (
	"context"
	"testing"
)

func TestCodecRoundTrip(t *testing.T) {
	c := New()
	payload, err := c.Encode(context.Background(), map[string]any{"temperature": 25})
	if err != nil {
		t.Fatal(err)
	}
	value, err := c.Decode(context.Background(), payload)
	if err != nil || value["temperature"] != float64(25) {
		t.Fatalf("unexpected decoded value: %#v, err=%v", value, err)
	}
}
