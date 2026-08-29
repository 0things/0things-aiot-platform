package mqtttransport

import "testing"

func TestDeviceKeyFromTopic(t *testing.T) {
	for _, test := range []struct{ topic, want string }{
		{"/sys/P001/D003/thing/event/property/post", "D003"},
		{"v1/devices/D001/telemetry", "D001"},
		{"devices/D002/commands/reply", "D002"},
		{"v1/gateway/telemetry", ""},
	} {
		if got := deviceKeyFromTopic(test.topic); got != test.want {
			t.Fatalf("deviceKeyFromTopic(%q) = %q, want %q", test.topic, got, test.want)
		}
	}
}
