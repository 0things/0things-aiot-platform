package modbuscodec

import (
	"context"
	"testing"
)

func TestCodecDecodeRegisters(t *testing.T) {
	c := New()
	value, err := c.Decode(context.Background(), []byte{0, 1, 0, 0, 0, 7, 1, 3, 4, 0, 25, 0, 26})
	if err != nil {
		t.Fatal(err)
	}
	registers := value["registers"].([]uint16)
	if len(registers) != 2 || registers[0] != 25 || registers[1] != 26 {
		t.Fatalf("unexpected registers: %#v", registers)
	}
}
