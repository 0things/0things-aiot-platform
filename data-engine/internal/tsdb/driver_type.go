package tsdb

import (
	"fmt"
	"strings"
)

// DriverType 定义时序数据库驱动类型枚举。
type DriverType string

const (
	DriverTypeTDengine    DriverType = "tdengine"
	DriverTypeIoTDB       DriverType = "iotdb"
	DriverTypeTimescaleDB DriverType = "timescaledb"
	DriverTypePostgreSQL  DriverType = "postgresql"
	DriverTypeInfluxDB    DriverType = "influxdb"
	DriverTypeClickHouse  DriverType = "clickhouse"
	DriverTypeSQLite      DriverType = "sqlite"
	DriverTypeMock        DriverType = "mock"
)

// allDriverTypes 内部静态枚举数组
var allDriverTypes = []DriverType{
	DriverTypeTDengine,
	DriverTypeIoTDB,
	DriverTypeTimescaleDB,
	DriverTypePostgreSQL,
	DriverTypeInfluxDB,
	DriverTypeClickHouse,
	DriverTypeSQLite,
	DriverTypeMock,
}

// String 返回枚举的字符串表达
func (d DriverType) String() string {
	return string(d)
}

// IsValid 检查枚举值是否合法
func (d DriverType) IsValid() bool {
	for _, item := range allDriverTypes {
		if item == d {
			return true
		}
	}
	return false
}

// AllDriverTypes 返回系统定义的所有时序驱动枚举列表
func AllDriverTypes() []DriverType {
	res := make([]DriverType, len(allDriverTypes))
	copy(res, allDriverTypes)
	return res
}

// ParseDriverType 从字符串安全解析为 DriverType 枚举，若非法返回错误
func ParseDriverType(raw string) (DriverType, error) {
	clean := DriverType(strings.ToLower(strings.TrimSpace(raw)))
	if clean.IsValid() {
		return clean, nil
	}
	return DriverTypeMock, fmt.Errorf("invalid tsdb driver type: %q, supported types: %v", raw, AllDriverTypes())
}
