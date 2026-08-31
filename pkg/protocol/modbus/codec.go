package modbuscodec

import (
	"context"
	"encoding/binary"
	"fmt"
)

// Codec 解析 Modbus TCP ADU 的基础读寄存器报文（功能码 03/04）。
type Codec struct{}

func New() *Codec           { return &Codec{} }
func (*Codec) Name() string { return "modbus" }

func (*Codec) Decode(_ context.Context, payload []byte) (map[string]any, error) {
	if len(payload) < 8 {
		return nil, fmt.Errorf("modbus tcp frame is too short")
	}
	length := int(binary.BigEndian.Uint16(payload[4:6]))
	if length < 2 || len(payload) < 6+length {
		return nil, fmt.Errorf("invalid modbus tcp length")
	}
	pdu := payload[7 : 6+length]
	result := map[string]any{"transaction_id": binary.BigEndian.Uint16(payload[:2]), "unit_id": payload[6], "function": pdu[0]}
	if pdu[0] != 3 && pdu[0] != 4 {
		return result, nil
	}
	if len(pdu) < 2 || int(pdu[1])+2 > len(pdu) || pdu[1]%2 != 0 {
		return nil, fmt.Errorf("invalid modbus register payload")
	}
	registers := make([]uint16, int(pdu[1])/2)
	for i := range registers {
		registers[i] = binary.BigEndian.Uint16(pdu[2+i*2:])
	}
	result["registers"] = registers
	return result, nil
}

func (*Codec) Encode(_ context.Context, value map[string]any) ([]byte, error) {
	unit, function, start, quantity := byteValue(value, "unit_id"), byteValue(value, "function"), uint16Value(value, "start_address"), uint16Value(value, "quantity")
	if unit == 0 || (function != 3 && function != 4) || quantity == 0 {
		return nil, fmt.Errorf("unit_id, function(03/04), start_address and quantity are required")
	}
	frame := make([]byte, 12)
	binary.BigEndian.PutUint16(frame[4:6], 6)
	frame[6] = unit
	frame[7] = function
	binary.BigEndian.PutUint16(frame[8:10], start)
	binary.BigEndian.PutUint16(frame[10:12], quantity)
	return frame, nil
}

func byteValue(value map[string]any, key string) byte {
	switch v := value[key].(type) {
	case int:
		return byte(v)
	case uint8:
		return v
	case float64:
		return byte(v)
	default:
		return 0
	}
}

func uint16Value(value map[string]any, key string) uint16 {
	switch v := value[key].(type) {
	case int:
		return uint16(v)
	case uint16:
		return v
	case float64:
		return uint16(v)
	default:
		return 0
	}
}
