package tsdb

import (
	"fmt"

	"github.com/spf13/viper"
)

// RootConfig 时序存储根配置结构体，统一聚合激活的驱动类型与各数据库独立专属的 Config 结构体。
type RootConfig struct {
	// Type 当前激活的时序驱动类型 (tdengine | sqlite | postgresql | timescaledb | clickhouse | influxdb | iotdb | mock)
	Type DriverType `mapstructure:"type" yaml:"type" json:"type"`

	// 各数据库独立专属配置结构体（定义在各自的驱动源文件中）
	TDengine    TDengineConfig    `mapstructure:"tdengine" yaml:"tdengine" json:"tdengine"`
	SQLite      SQLiteConfig      `mapstructure:"sqlite" yaml:"sqlite" json:"sqlite"`
	PostgreSQL  PostgreSQLConfig  `mapstructure:"postgresql" yaml:"postgresql" json:"postgresql"`
	TimescaleDB TimescaleDBConfig `mapstructure:"timescaledb" yaml:"timescaledb" json:"timescaledb"`
	ClickHouse  ClickHouseConfig  `mapstructure:"clickhouse" yaml:"clickhouse" json:"clickhouse"`
	InfluxDB    InfluxDBConfig    `mapstructure:"influxdb" yaml:"influxdb" json:"influxdb"`
	IoTDB       IoTDBConfig       `mapstructure:"iotdb" yaml:"iotdb" json:"iotdb"`
}

// SetDefaults 自动补齐各数据库独立配置的默认值
func (r *RootConfig) SetDefaults() {
	if r.Type == "" {
		r.Type = DriverTypeMock
	}

	// TDengine
	if r.TDengine.Database == "" {
		r.TDengine.Database = "things_tsdb"
	}
	if r.TDengine.Table == "" {
		r.TDengine.Table = "device_properties"
	}

	// SQLite
	if r.SQLite.Path == "" {
		r.SQLite.Path = "data/things_tsdb.db"
	}
	if r.SQLite.Table == "" {
		r.SQLite.Table = "device_properties"
	}

	// PostgreSQL / TimescaleDB
	if r.PostgreSQL.Table == "" {
		r.PostgreSQL.Table = "device_properties"
	}
	if r.TimescaleDB.Table == "" {
		r.TimescaleDB.Table = "device_properties"
	}

	// ClickHouse
	if r.ClickHouse.Database == "" {
		r.ClickHouse.Database = "things_tsdb"
	}
	if r.ClickHouse.Table == "" {
		r.ClickHouse.Table = "device_properties"
	}

	// InfluxDB
	if r.InfluxDB.URL == "" {
		r.InfluxDB.URL = "http://127.0.0.1:8086"
	}
	if r.InfluxDB.Bucket == "" {
		r.InfluxDB.Bucket = "things_tsdb"
	}

	// IoTDB
	if r.IoTDB.Host == "" {
		r.IoTDB.Host = "127.0.0.1"
	}
	if r.IoTDB.Port == "" {
		r.IoTDB.Port = "6667"
	}
	if r.IoTDB.StorageGroup == "" {
		r.IoTDB.StorageGroup = "root.things"
	}
}

// LoadRootConfigFromViper 从 Viper 实例中强类型加载独立多 DB 根配置
func LoadRootConfigFromViper(v *viper.Viper) (*RootConfig, error) {
	cfg := &RootConfig{}

	// 1. 结构化独立子块解析
	if v.IsSet("tsdb") {
		sub := v.Sub("tsdb")
		if sub != nil {
			if err := sub.Unmarshal(cfg); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tsdb root config: %w", err)
			}
		}
	} else {
		_ = v.Unmarshal(cfg)
	}

	// 2. 兼容传统的顶层扁平配置（向后兼容）
	if cfg.Type == "" {
		cfg.Type = DriverType(v.GetString("tsdb.type"))
	}
	flatDSN := v.GetString("tsdb.dsn")
	if flatDSN != "" {
		if cfg.TDengine.DSN == "" {
			cfg.TDengine.DSN = flatDSN
		}
		if cfg.PostgreSQL.DSN == "" {
			cfg.PostgreSQL.DSN = flatDSN
		}
		if cfg.TimescaleDB.DSN == "" {
			cfg.TimescaleDB.DSN = flatDSN
		}
		if cfg.ClickHouse.DSN == "" {
			cfg.ClickHouse.DSN = flatDSN
		}
	}
	if flatPath := v.GetString("tsdb.path"); flatPath != "" && cfg.SQLite.Path == "" {
		cfg.SQLite.Path = flatPath
	}
	if flatURL := v.GetString("tsdb.url"); flatURL != "" && cfg.InfluxDB.URL == "" {
		cfg.InfluxDB.URL = flatURL
	}
	if flatToken := v.GetString("tsdb.token"); flatToken != "" && cfg.InfluxDB.Token == "" {
		cfg.InfluxDB.Token = flatToken
	}
	if flatOrg := v.GetString("tsdb.org"); flatOrg != "" && cfg.InfluxDB.Org == "" {
		cfg.InfluxDB.Org = flatOrg
	}
	if flatBucket := v.GetString("tsdb.bucket"); flatBucket != "" && cfg.InfluxDB.Bucket == "" {
		cfg.InfluxDB.Bucket = flatBucket
	}

	cfg.SetDefaults()
	return cfg, nil
}
