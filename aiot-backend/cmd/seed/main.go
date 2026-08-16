package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/glebarez/sqlite"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	userDB, err := sql.Open("sqlite", "storage/aiot-test.db?_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}
	defer userDB.Close()

	deviceDB, err := sql.Open("sqlite", "storage/aiot-device.db?_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}
	defer deviceDB.Close()

	// Clear existing data
	fmt.Println("Clearing existing data...")
	userDB.Exec("DELETE FROM users")
	deviceDB.Exec("DELETE FROM device_events")
	deviceDB.Exec("DELETE FROM device_tags")
	deviceDB.Exec("DELETE FROM device_shadow_histories")
	deviceDB.Exec("DELETE FROM device_shadows")
	deviceDB.Exec("DELETE FROM device_states")
	deviceDB.Exec("DELETE FROM devices")
	deviceDB.Exec("DELETE FROM products")
	deviceDB.Exec("DELETE FROM product_ts_ls")

	// --- users (50) ---
	fmt.Println("Seeding users...")
	for i := 1; i <= 50; i++ {
		_, err := userDB.Exec(`INSERT OR IGNORE INTO users (user_id, nickname, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("user_%03d", i),
			fmt.Sprintf("用户%d", i),
			"hashed_password_"+fmt.Sprintf("%03d", i),
			fmt.Sprintf("user%d@example.com", i),
			time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("user insert error: %v", err)
		}
	}

	// --- products (50) ---
	fmt.Println("Seeding products...")
	categories := []string{"传感器", "网关", "摄像头", "温控器", "电表", "水表", "气体检测", "门锁", "开关", "插座"}
	statuses := []string{"active", "inactive", "draft"}
	nodeTypes := []string{"device", "gateway", "sub_device"}
	connMethods := []string{"wifi", "ble", "zigbee", "lora", "4g", "ethernet"}
	protocols := []string{"mqtt", "http", "coap", "modbus", "tcp"}

	for i := 1; i <= 50; i++ {
		metadata, _ := json.Marshal(map[string]interface{}{
			"manufacturer": fmt.Sprintf("厂商%d", rand.Intn(10)+1),
			"model":        fmt.Sprintf("MODEL-%c%c%c", 'A'+rand.Intn(26), 'A'+rand.Intn(26), 'A'+rand.Intn(26)),
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO products (product_key, name, description, category, status, metadata, node_type, connectivity_method, access_protocol, tenant_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("pk_product_%03d", i),
			fmt.Sprintf("产品%d", i),
			fmt.Sprintf("这是产品%d的描述信息", i),
			categories[rand.Intn(len(categories))],
			statuses[rand.Intn(len(statuses))],
			string(metadata),
			nodeTypes[rand.Intn(len(nodeTypes))],
			connMethods[rand.Intn(len(connMethods))],
			protocols[rand.Intn(len(protocols))],
			int64(rand.Intn(3) + 1),
			time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product insert error: %v", err)
		}
	}

	// --- devices (50) ---
	fmt.Println("Seeding devices...")
	for i := 1; i <= 50; i++ {
		metadata, _ := json.Marshal(map[string]interface{}{
			"firmware_version": fmt.Sprintf("v%d.%d.%d", rand.Intn(3)+1, rand.Intn(10), rand.Intn(20)),
			"hardware_version": fmt.Sprintf("hw%d.%d", rand.Intn(5)+1, rand.Intn(3)),
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO devices (device_key, name, product_id, tenant_id, enabled, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("dk_device_%03d", i),
			fmt.Sprintf("设备%d", i),
			int64(rand.Intn(50) + 1),
			int64(rand.Intn(3) + 1),
			rand.Intn(2) == 1,
			string(metadata),
			time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device insert error: %v", err)
		}
	}

	// --- device_states (50) ---
	fmt.Println("Seeding device_states...")
	states := []string{"online", "offline", "sleep", "error", "updating"}
	for i := 1; i <= 50; i++ {
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_states (device_key, state, created_at, updated_at) VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("dk_device_%03d", i),
			states[rand.Intn(len(states))],
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_state insert error: %v", err)
		}
	}

	// --- device_shadows (50) ---
	fmt.Println("Seeding device_shadows...")
	for i := 1; i <= 50; i++ {
		desired, _ := json.Marshal(map[string]interface{}{
			"target_temperature": rand.Intn(30) + 15,
			"target_humidity":    rand.Intn(60) + 20,
		})
		reported, _ := json.Marshal(map[string]interface{}{
			"current_temperature": rand.Intn(30) + 15,
			"current_humidity":    rand.Intn(60) + 20,
			"battery":             rand.Intn(100),
		})
		meta, _ := json.Marshal(map[string]interface{}{
			"source": "device_report",
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_shadows (device_id, desired, reported, metadata, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			int64(i),
			string(desired),
			string(reported),
			string(meta),
			int64(rand.Intn(100) + 1),
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_shadow insert error: %v", err)
		}
	}

	// --- device_shadow_histories (50) ---
	fmt.Println("Seeding device_shadow_histories...")
	sources := []string{"device_report", "cloud_set", "rule_engine", "user_action"}
	for i := 1; i <= 50; i++ {
		desired, _ := json.Marshal(map[string]interface{}{"target_temp": rand.Intn(30) + 15})
		reported, _ := json.Marshal(map[string]interface{}{"current_temp": rand.Intn(30) + 15})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_shadow_histories (device_id, version, source, desired, reported, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			int64(rand.Intn(50) + 1),
			int64(i),
			sources[rand.Intn(len(sources))],
			string(desired),
			string(reported),
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		)
		if err != nil {
			log.Printf("device_shadow_history insert error: %v", err)
		}
	}

	// --- device_tags (50) ---
	fmt.Println("Seeding device_tags...")
	tagKeys := []string{"location", "floor", "area", "department", "owner", "priority", "region", "building"}
	tagVals := []string{"1楼", "2楼", "3楼", "A区", "B区", "C区", "办公室", "仓库", "会议室", "大厅"}
	tagSources := []string{"manual", "auto", "import", "rule"}
	for i := 1; i <= 50; i++ {
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_tags (device_id, key, value, source, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			int64(rand.Intn(50) + 1),
			tagKeys[rand.Intn(len(tagKeys))],
			tagVals[rand.Intn(len(tagVals))],
			tagSources[rand.Intn(len(tagSources))],
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_tag insert error: %v", err)
		}
	}

	// --- device_events (50) ---
	fmt.Println("Seeding device_events...")
	eventTypes := []string{"temperature_alert", "humidity_alert", "offline", "online", "firmware_update", "battery_low", "reboot", "data_report"}
	for i := 1; i <= 50; i++ {
		data, _ := json.Marshal(map[string]interface{}{
			"value":     rand.Float64() * 100,
			"unit":      "℃",
			"severity":  []string{"info", "warning", "critical"}[rand.Intn(3)],
			"timestamp": time.Now().Unix(),
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_events (device_id, event_type, event_at, data, created_at) VALUES (?, ?, ?, ?, ?)`,
			int64(rand.Intn(50) + 1),
			eventTypes[rand.Intn(len(eventTypes))],
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			string(data),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_event insert error: %v", err)
		}
	}

	// --- product_ts_ls (50) ---
	fmt.Println("Seeding product_ts_ls...")
	for i := 1; i <= 50; i++ {
		tsl, _ := json.Marshal(map[string]interface{}{
			"schema": "https://iotx-tsl.oss-ap-southeast-1.aliyuncs.com/schema.json",
			"profile": map[string]interface{}{
				"productKey": fmt.Sprintf("pk_product_%03d", i),
			},
			"properties": []interface{}{
				map[string]interface{}{
					"identifier": "temperature",
					"name":       "温度",
					"type":       "float",
				},
			},
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO product_ts_ls (product_key, tsl, created_at, updated_at) VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("pk_product_%03d", i),
			string(tsl),
			time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product_tsl insert error: %v", err)
		}
	}

	fmt.Println("Seed complete!")
}
