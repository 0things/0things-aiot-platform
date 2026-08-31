package gb28181codec

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
)

// Codec 解析 GB28181 使用的 SIP 文本起始行和头部；SDP body 原样保留。
type Codec struct{}

func New() *Codec           { return &Codec{} }
func (*Codec) Name() string { return "gb28181" }

func (*Codec) Decode(_ context.Context, payload []byte) (map[string]any, error) {
	scanner := bufio.NewScanner(bytes.NewReader(payload))
	result := map[string]any{"headers": map[string]string{}}
	first := true
	headers := result["headers"].(map[string]string)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if first {
			result["start_line"] = line
			first = false
			continue
		}
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if first {
		return nil, fmt.Errorf("empty SIP message")
	}
	return result, nil
}

func (*Codec) Encode(_ context.Context, value map[string]any) ([]byte, error) {
	start, ok := value["start_line"].(string)
	if !ok || strings.TrimSpace(start) == "" {
		return nil, fmt.Errorf("start_line is required")
	}
	var b strings.Builder
	b.WriteString(start + "\r\n")
	if headers, ok := value["headers"].(map[string]string); ok {
		for key, item := range headers {
			fmt.Fprintf(&b, "%s: %s\r\n", key, item)
		}
	}
	b.WriteString("\r\n")
	if body, ok := value["body"].(string); ok {
		b.WriteString(body)
	}
	return []byte(b.String()), nil
}
