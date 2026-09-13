package topic

import "strings"

func ExtractDeviceKey(value string) string {
	parts := strings.Split(value, "/")
	if len(parts) >= 4 && parts[1] == "sys" {
		return parts[3]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[5]
	}
	return ""
}

func ExtractProductKey(value string) string {
	parts := strings.Split(value, "/")
	if len(parts) >= 3 && parts[1] == "sys" {
		return parts[2]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[4]
	}
	return ""
}
