package gb28181codec

import (
	"context"
	"testing"
)

func TestCodecDecodeSIP(t *testing.T) {
	value, err := New().Decode(context.Background(), []byte("MESSAGE sip:a SIP/2.0\r\nCall-ID: 1\r\n\r\n"))
	if err != nil || value["start_line"] != "MESSAGE sip:a SIP/2.0" {
		t.Fatalf("unexpected SIP value: %#v, err=%v", value, err)
	}
}
