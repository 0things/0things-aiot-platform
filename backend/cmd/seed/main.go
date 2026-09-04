package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"aiot-backend/pkg/config"

	_ "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()

	conf := config.NewConfig(*envConf)
	driver := conf.GetString("data.db.aiot.driver")
	dsn := conf.GetString("data.db.aiot.dsn")
	if driver == "" {
		driver = "sqlite"
	}
	if dsn == "" {
		dsn = "storage/aiot-device.db?_busy_timeout=5000"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userDB := db
	deviceDB := db

	// Clear existing data
	fmt.Println("Clearing existing data...")
	userDB.Exec("DELETE FROM organization_users")
	userDB.Exec("DELETE FROM organizations")
	userDB.Exec("DELETE FROM users")
	deviceDB.Exec("DELETE FROM scene_linkage_detail")
	deviceDB.Exec("DELETE FROM scene_linkage")
	deviceDB.Exec("DELETE FROM device_group_members")
	deviceDB.Exec("DELETE FROM device_groups")
	deviceDB.Exec("DELETE FROM device_push_records")
	deviceDB.Exec("DELETE FROM device_endpoints")
	deviceDB.Exec("DELETE FROM device_service_invocations")
	deviceDB.Exec("DELETE FROM device_events")
	deviceDB.Exec("DELETE FROM device_tags")
	deviceDB.Exec("DELETE FROM device_shadow_histories")
	deviceDB.Exec("DELETE FROM device_shadows")
	deviceDB.Exec("DELETE FROM device_states")
	deviceDB.Exec("DELETE FROM devices")
	deviceDB.Exec("DELETE FROM product_message_parsers")
	deviceDB.Exec("DELETE FROM product_protocols")
	deviceDB.Exec("DELETE FROM product_ts_ls")
	deviceDB.Exec("DELETE FROM products")
	deviceDB.Exec("DELETE FROM categories")
	deviceDB.Exec("DELETE FROM ota_device_upgrade_status")
	deviceDB.Exec("DELETE FROM ota_upgrade_batches")
	deviceDB.Exec("DELETE FROM ota_packages")
	userDB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('users','organizations','organization_users')")
	deviceDB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('device_events','device_service_invocations','device_tags','device_shadow_histories','device_shadows','device_states','devices','products','categories','product_protocols','product_message_parsers','product_ts_ls','device_endpoints','device_push_records','device_groups','device_group_members','scene_linkage','scene_linkage_detail','ota_packages','ota_upgrade_batches','ota_device_upgrade_status')")

	// --- categories ---
	fmt.Println("Seeding categories...")
	defaultCategories := []string{"传感器", "执行器", "网关", "控制器", "显示设备", "摄像头", "其他"}
	for i, name := range defaultCategories {
		_, err := deviceDB.Exec(`INSERT OR REPLACE INTO categories (id, name, sort, enabled) VALUES (?, ?, ?, ?)`,
			i+1, name, i, true,
		)
		if err != nil {
			log.Printf("category insert error: %v", err)
		}
	}

	// --- organizations (3) ---
	fmt.Println("Seeding organizations...")
	orgs := []struct {
		id   int64
		name string
		slug string
	}{
		{id: 1, name: "默认组织", slug: "default-org"},
		{id: 2, name: "研发中心", slug: "rd-center"},
		{id: 3, name: "测试环境", slug: "test-env"},
	}
	for _, org := range orgs {
		_, err := userDB.Exec(`INSERT OR REPLACE INTO organizations (id, name, slug, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			org.id,
			org.name,
			org.slug,
			time.Now().Add(-30*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("organization insert error: %v", err)
		}
	}

	// --- users (50) & organization_users ---
	fmt.Println("Seeding users and memberships...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt generate password error: %v", err)
	}

	for i := 1; i <= 50; i++ {
		userID := fmt.Sprintf("user_%03d", i)
		_, err := userDB.Exec(`INSERT OR IGNORE INTO users (user_id, nickname, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			userID,
			fmt.Sprintf("用户%d", i),
			string(hashedPassword),
			fmt.Sprintf("user%d@example.com", i),
			time.Now().Add(-time.Duration(rand.Intn(365))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("user insert error: %v", err)
		}

		if i == 1 {
			// user_001 加入全部 3 个组织，且 org 3 的 last_login_at 最新
			t1 := time.Now().Add(-24 * time.Hour)
			t2 := time.Now().Add(-2 * time.Hour)
			t3 := time.Now()
			_, _ = userDB.Exec(`INSERT OR REPLACE INTO organization_users (organization_id, user_id, last_login_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, 1, userID, t1, time.Now(), time.Now())
			_, _ = userDB.Exec(`INSERT OR REPLACE INTO organization_users (organization_id, user_id, last_login_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, 2, userID, t2, time.Now(), time.Now())
			_, _ = userDB.Exec(`INSERT OR REPLACE INTO organization_users (organization_id, user_id, last_login_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, 3, userID, t3, time.Now(), time.Now())
		} else {
			// 其它用户随机分配 1~2 个组织
			orgID := int64(rand.Intn(3) + 1)
			var loginTime *time.Time
			if rand.Intn(2) == 1 {
				t := time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour)
				loginTime = &t
			}
			_, _ = userDB.Exec(`INSERT OR REPLACE INTO organization_users (organization_id, user_id, last_login_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, orgID, userID, loginTime, time.Now(), time.Now())
		}
	}

	// --- products (50) ---
	fmt.Println("Seeding products...")
	statuses := []string{"active", "inactive", "draft"}
	nodeTypes := []string{"device", "gateway", "sub_device"}
	connMethods := []string{"wifi", "ble", "zigbee", "lora", "4g", "ethernet"}
	protocols := []string{"mqtt", "http", "coap", "modbus", "tcp"}

	productOrgMap := make(map[int]int64)
	for i := 1; i <= 50; i++ {
		orgID := int64(rand.Intn(3) + 1)
		categoryID := int64(rand.Intn(len(defaultCategories)) + 1)
		productOrgMap[i] = orgID
		metadata, _ := json.Marshal(map[string]interface{}{
			"manufacturer": fmt.Sprintf("厂商%d", rand.Intn(10)+1),
			"model":        fmt.Sprintf("MODEL-%c%c%c", 'A'+rand.Intn(26), 'A'+rand.Intn(26), 'A'+rand.Intn(26)),
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO products (product_key, name, description, category_id, status, metadata, node_type, connectivity_method, access_protocol, organization_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("pk_product_%03d", i),
			fmt.Sprintf("产品%d", i),
			fmt.Sprintf("这是产品%d的描述信息", i),
			categoryID,
			statuses[rand.Intn(len(statuses))],
			string(metadata),
			nodeTypes[rand.Intn(len(nodeTypes))],
			connMethods[rand.Intn(len(connMethods))],
			protocols[rand.Intn(len(protocols))],
			orgID,
			time.Now().Add(-time.Duration(rand.Intn(365))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product insert error: %v", err)
		}
	}

	// --- devices (50) ---
	fmt.Println("Seeding devices...")
	productDevices := make(map[int64][]int64)
	for i := 1; i <= 50; i++ {
		prodID := int64((i-1)%10 + 1) // 50 台设备分布在前 10 个产品中，每个产品 5 台
		productDevices[prodID] = append(productDevices[prodID], int64(i))
		metadata, _ := json.Marshal(map[string]interface{}{
			"firmware_version": fmt.Sprintf("v%d.%d.%d", rand.Intn(3)+1, rand.Intn(10), rand.Intn(20)),
			"hardware_version": fmt.Sprintf("hw%d.%d", rand.Intn(5)+1, rand.Intn(3)),
		})
		orgID := productOrgMap[int(prodID)]
		if orgID == 0 {
			orgID = int64(1)
		}
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO devices (id, device_key, name, product_id, organization_id, enabled, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			int64(i),
			fmt.Sprintf("dk_device_%03d", i),
			fmt.Sprintf("设备%03d", i),
			prodID,
			orgID,
			true,
			string(metadata),
			time.Now().Add(-time.Duration(rand.Intn(365))*24*time.Hour),
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
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
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
			int64(rand.Intn(100)+1),
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
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
			int64(rand.Intn(50)+1),
			int64(i),
			sources[rand.Intn(len(sources))],
			string(desired),
			string(reported),
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
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
			int64(rand.Intn(50)+1),
			tagKeys[rand.Intn(len(tagKeys))],
			tagVals[rand.Intn(len(tagVals))],
			tagSources[rand.Intn(len(tagSources))],
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
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
			int64(rand.Intn(50)+1),
			eventTypes[rand.Intn(len(eventTypes))],
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
			string(data),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_event insert error: %v", err)
		}
	}

	// --- device_service_invocations (80) ---
	fmt.Println("Seeding device_service_invocations...")
	serviceDefinitions := []struct {
		identifier string
		name       string
	}{
		{"reboot", "重启设备"},
		{"set", "设置开关"},
		{"set_temperature", "设置目标温度"},
		{"firmware_upgrade", "固件升级"},
	}
	for i := 0; i < 80; i++ {
		service := serviceDefinitions[rand.Intn(len(serviceDefinitions))]
		input, _ := json.Marshal(map[string]interface{}{
			"requestId": fmt.Sprintf("seed-call-%03d", i+1),
			"value":     rand.Intn(100),
		})
		var output *string
		if i%5 != 0 {
			payload, _ := json.Marshal(map[string]interface{}{"success": true, "code": 0})
			value := string(payload)
			output = &value
		}
		invokedAt := time.Now().Add(-time.Duration(rand.Intn(7*24*60)) * time.Minute)
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_service_invocations (uuid, device_id, service_identifier, service_name, input_params, output_params, invoked_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			uuid.NewString(),
			int64(rand.Intn(50)+1),
			service.identifier,
			service.name,
			string(input),
			output,
			invokedAt,
			invokedAt,
			invokedAt,
		)
		if err != nil {
			log.Printf("device_service_invocation insert error: %v", err)
		}
	}

	// --- product_protocols (50) ---
	fmt.Println("Seeding product_protocols...")
	for i := 1; i <= 50; i++ {
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO product_protocols (id, product_id, transport_protocol, application_protocol, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			i,
			int64(i),
			protocols[rand.Intn(len(protocols))],
			"json",
			time.Now().Add(-time.Duration(rand.Intn(365))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product_protocol insert error: %v", err)
		}
	}

	// --- product_message_parsers (10) ---
	fmt.Println("Seeding product_message_parsers...")
	sampleScript := `function parse(payload, topic) {
  var data = JSON.parse(payload);
  return {
    temperature: data.temp || 0,
    humidity: data.hum || 0
  };
}`
	for i := 1; i <= 10; i++ {
		orgID := productOrgMap[i]
		if orgID == 0 {
			orgID = 1
		}
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO product_message_parsers (id, organization_id, product_id, language, script, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			i,
			orgID,
			int64(i),
			"javascript",
			sampleScript,
			time.Now().Add(-time.Duration(rand.Intn(100))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product_message_parser insert error: %v", err)
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
				map[string]interface{}{
					"identifier": "humidity",
					"name":       "湿度",
					"type":       "float",
				},
			},
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO product_ts_ls (product_key, tsl, created_at, updated_at) VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("pk_product_%03d", i),
			string(tsl),
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product_tsl insert error: %v", err)
		}
	}

	// --- device_endpoints (50) ---
	fmt.Println("Seeding device_endpoints...")
	for i := 1; i <= 50; i++ {
		prodID := int64((i-1)%10 + 1)
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_endpoints (id, device_id, product_protocol_id, endpoint, credentials, protocol_config, enabled, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			i,
			int64(i),
			prodID,
			fmt.Sprintf("mqtt://broker.0things.com:1883/devices/dk_device_%03d", i),
			fmt.Sprintf(`{"token":"tok_%03d"}`, i),
			`{"qos":1}`,
			true,
			"active",
			time.Now().Add(-time.Duration(rand.Intn(100))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_endpoint insert error: %v", err)
		}
	}

	// --- device_push_records (50) ---
	fmt.Println("Seeding device_push_records...")
	opTypes := []string{"command", "property_set", "config_update", "reboot"}
	for i := 1; i <= 50; i++ {
		dID := int64(rand.Intn(50) + 1)
		op := opTypes[rand.Intn(len(opTypes))]
		status := []string{"success", "failed", "pending"}[rand.Intn(3)]
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_push_records (id, device_id, operation_type, operation_name, payload, status, error_message, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			i,
			dID,
			op,
			fmt.Sprintf("下发指令-%s", op),
			fmt.Sprintf(`{"action":"%s","val":%d}`, op, rand.Intn(100)),
			status,
			"",
			"system",
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_push_record insert error: %v", err)
		}
	}

	// --- device_groups & members ---
	fmt.Println("Seeding device_groups & members...")
	groups := []struct {
		name      string
		groupType string
		desc      string
		rule      string
		orgID     int64
	}{
		{name: "1号车间温湿度传感器组", groupType: "manual", desc: "生产车间关键温湿度监测", orgID: 1},
		{name: "智能电表集中控制组", groupType: "manual", desc: "厂区各配电箱智能电表", orgID: 1},
		{name: "在线设备动态分组", groupType: "dynamic", desc: "自动匹配在线状态设备", rule: `{"state":"online"}`, orgID: 1},
		{name: "研发测试设备组", groupType: "manual", desc: "研发实验室设备", orgID: 2},
		{name: "测试环境仿真组", groupType: "manual", desc: "QA自动化测试仿真设备", orgID: 3},
	}
	for gIdx, g := range groups {
		gID := int64(gIdx + 1)
		gUUID := uuid.NewString()
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO device_groups (id, group_uuid, organization_id, name, type, description, rule, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			gID,
			gUUID,
			g.orgID,
			g.name,
			g.groupType,
			g.desc,
			g.rule,
			time.Now().Add(-20*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("device_group insert error: %v", err)
		}
		if g.groupType == "manual" {
			for m := 1; m <= 4; m++ {
				devID := int64((gIdx*4+m)%50 + 1)
				_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO device_group_members (group_id, device_id, created_at, updated_at) VALUES (?, ?, ?, ?)`,
					gID,
					devID,
					time.Now().Add(-15*24*time.Hour),
					time.Now(),
				)
			}
		}
	}

	// --- scene_linkage & detail ---
	fmt.Println("Seeding scene_linkage & details...")
	scenes := []struct {
		name    string
		desc    string
		orgID   int64
		enable  int
		trigger string
		action  string
	}{
		{
			name:    "高温自动联动排风",
			desc:    "当车间温度超过35℃时，自动开启排风机并发送预警通知",
			orgID:   1,
			enable:  1,
			trigger: `{"type":"property","property":"temperature","operator":">=","value":35}`,
			action:  `{"type":"device_command","command":"open_fan","notify":"webhook"}`,
		},
		{
			name:    "夜间安防布防告警",
			desc:    "夜间时段红外或门锁异常开启时，立即联动声光报警",
			orgID:   1,
			enable:  1,
			trigger: `{"type":"event","event":"door_opened","time_range":"22:00-06:00"}`,
			action:  `{"type":"alarm","level":"critical","push":true}`,
		},
		{
			name:    "电量超载保护联动",
			desc:    "检测到瞬时功率大于设定阈值时执行分闸保护",
			orgID:   2,
			enable:  0,
			trigger: `{"type":"property","property":"power","operator":">","value":5000}`,
			action:  `{"type":"device_command","command":"power_off"}`,
		},
	}
	for sIdx, sc := range scenes {
		sID := int64(sIdx + 1)
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO scene_linkage (id, organization_id, name, description, enable, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			sID,
			sc.orgID,
			sc.name,
			sc.desc,
			sc.enable,
			time.Now().Add(-10*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("scene_linkage insert error: %v", err)
		}
		_, err = deviceDB.Exec(`INSERT OR IGNORE INTO scene_linkage_detail (id, scene_id, trigger_config, action_config, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			sID,
			sID,
			sc.trigger,
			sc.action,
			time.Now().Add(-10*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("scene_linkage_detail insert error: %v", err)
		}
	}

	// --- OTA 数据体系 (ota_packages, ota_upgrade_batches, ota_device_upgrade_status) ---
	fmt.Println("Seeding OTA packages, upgrade batches, and device deployment records...")

	type otaSeedDef struct {
		name         string
		version      string
		prodID       int64
		pkgType      string
		status       string
		uploadType   string
		desc         string
		releaseNotes string
	}

	otaSeeds := []otaSeedDef{
		{
			name:         "Smart-Sensor-Firmware-V2.1.0",
			version:      "2.1.0",
			prodID:       1,
			pkgType:      "firmware",
			status:       "deploying",
			uploadType:   "oss",
			desc:         "智能温湿度传感器基础固件，优化低功耗唤醒机制与上报频次",
			releaseNotes: "1. 优化 BLE 连接稳定性；\n2. 降低待机功耗 30%；\n3. 修复传感器偶尔读取超时的问题。",
		},
		{
			name:         "Gateway-OS-Image-V3.0.4",
			version:      "3.0.4",
			prodID:       2,
			pkgType:      "full",
			status:       "success",
			uploadType:   "oss",
			desc:         "边缘计算智能网关系统完整镜像包，集成最新 MQTT Broker 桥接客户端",
			releaseNotes: "1. 升级底层 Linux 内核安全补丁；\n2. 支持 MQTT 5.0 协议属性；\n3. 增强离线数据本地缓存队列。",
		},
		{
			name:         "PowerMeter-Driver-Patch-V1.2.0",
			version:      "1.2.0",
			prodID:       3,
			pkgType:      "config",
			status:       "partial",
			uploadType:   "binary",
			desc:         "智能三相电表采集驱动补丁配置，更新 Modbus 寄存器映射表",
			releaseNotes: "1. 适配新版互感器变比配置；\n2. 增加谐波失真度计算上报。",
		},
		{
			name:         "GasSensor-Core-V1.0.1",
			version:      "1.0.1",
			prodID:       4,
			pkgType:      "firmware",
			status:       "deploying",
			uploadType:   "oss",
			desc:         "工业气体检测仪防爆固件，提升传感器预热与零点漂移校准精度",
			releaseNotes: "1. 新增零点自动校准算法；\n2. 优化高浓度气体声光报警响应时间。",
		},
		{
			name:         "SmartLock-Security-V2.0.0",
			version:      "2.0.0",
			prodID:       5,
			pkgType:      "firmware",
			status:       "draft",
			uploadType:   "binary",
			desc:         "智能门锁主控板固件，升级国密 SM4 加密通信协议",
			releaseNotes: "1. 引入双向认证安全芯片握手；\n2. 防暴力破解重试锁定逻辑优化。",
		},
		{
			name:         "HVAC-Controller-V1.5.2",
			version:      "1.5.2",
			prodID:       6,
			pkgType:      "firmware",
			status:       "success",
			uploadType:   "oss",
			desc:         "中央空调温控器 PID 控制算法升级包",
			releaseNotes: "1. 动态 PID 参数自整定支持；\n2. 优化节能运行模式逻辑。",
		},
		{
			name:         "Socket-Config-Template-V1.1",
			version:      "1.1.0",
			prodID:       7,
			pkgType:      "config",
			status:       "draft",
			uploadType:   "binary",
			desc:         "计量智能插座过载保护门限与定时任务默认配置文件",
			releaseNotes: "1. 更新默认最大电流保护值为 16A；\n2. 增加电量统计日结自动推送。",
		},
		{
			name:         "Camera-AI-Motion-V3.2.0",
			version:      "3.2.0",
			prodID:       8,
			pkgType:      "full",
			status:       "archived",
			uploadType:   "oss",
			desc:         "AI 摄像头人形检测算法模型历史归档版本",
			releaseNotes: "1. 基础人形检测模型发布；\n2. 历史稳定版本归档。",
		},
		{
			name:         "WaterMeter-NB-IoT-V1.0.8",
			version:      "1.0.8",
			prodID:       9,
			pkgType:      "firmware",
			status:       "deploying",
			uploadType:   "oss",
			desc:         "远传水表 NB-IoT 窄带通信射频参数优化固件",
			releaseNotes: "1. 增强弱信号区域 PSM 休眠与重连策略；\n2. 电池剩余寿命预测模型优化。",
		},
		{
			name:         "Switch-Touch-Firmware-V2.0.1",
			version:      "2.0.1",
			prodID:       10,
			pkgType:      "firmware",
			status:       "success",
			uploadType:   "oss",
			desc:         "智能触摸开关电容触摸灵敏度校准固件",
			releaseNotes: "1. 优化湿手触摸误触发滤波算法；\n2. 背光呼吸灯调光平滑度提升。",
		},
	}

	batchSeq := 1
	statusSeq := 1

	for pkgIdx, seed := range otaSeeds {
		pkgID := int64(pkgIdx + 1)
		pkgUUID := uuid.NewString()
		orgID := productOrgMap[int(seed.prodID)]
		if orgID == 0 {
			orgID = int64(1)
		}

		var releasedAt *time.Time
		createdAt := time.Now().Add(-time.Duration(rand.Intn(60)+5) * 24 * time.Hour)
		if seed.status != "draft" {
			t := createdAt.Add(time.Duration(rand.Intn(48)+12) * time.Hour)
			releasedAt = &t
		}

		fileSize := int64(rand.Intn(15000000) + 500000)
		checksum := fmt.Sprintf("%016x%016x", rand.Int63(), rand.Int63())
		fileURL := fmt.Sprintf("https://static.0things.com/firmware/%s/%s-%s.bin", seed.pkgType, seed.name, seed.version)

		fullDesc := seed.desc
		if seed.releaseNotes != "" {
			fullDesc = fmt.Sprintf("%s\n\n更新说明：\n%s", seed.desc, seed.releaseNotes)
		}

		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO ota_packages (
			id, uuid, package_name, version, product_id, organization_id, package_type, status,
			upload_type, file_url, file_size, checksum, description, release_notes,
			released_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			pkgID,
			pkgUUID,
			seed.name,
			seed.version,
			seed.prodID,
			orgID,
			seed.pkgType,
			seed.status,
			seed.uploadType,
			fileURL,
			fileSize,
			checksum,
			fullDesc,
			seed.releaseNotes,
			releasedAt,
			createdAt,
			time.Now(),
		)
		if err != nil {
			log.Printf("ota_package insert error: %v", err)
			continue
		}

		// 为部署中/已成功/部分成功的包生成升级批次和设备升级记录
		devList := productDevices[seed.prodID]
		if len(devList) == 0 || seed.status == "draft" || seed.status == "archived" {
			continue
		}

		pkgIDStr := strconv.FormatInt(pkgID, 10)

		switch seed.status {
		case "deploying":
			// 创建两个批次：批次 1（已完成），批次 2（进行中）
			b1UUID := uuid.NewString()
			b1Devices := devList[:len(devList)/2]
			if len(b1Devices) == 0 {
				b1Devices = devList[:1]
			}
			b1Time := createdAt.Add(2 * time.Hour)
			_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_upgrade_batches (
				id, batch_id, ota_package_id, batch_name, upgrade_strategy, status, target_device_count, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				batchSeq, b1UUID, pkgIDStr, "灰度一期-内部测试批次", "static", "completed", int32(len(b1Devices)), b1Time, b1Time.Add(1*time.Hour),
			)
			batchSeq++

			for _, dID := range b1Devices {
				ts := b1Time.Add(45 * time.Minute).Unix()
				_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_device_upgrade_status (
					id, device_id, ota_package_id, upgrade_batch_id, status, current_version, last_status_change_ts, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					statusSeq, dID, pkgIDStr, b1UUID, "success", seed.version, ts, b1Time, b1Time.Add(45*time.Minute),
				)
				statusSeq++
			}

			// 批次 2：进行中
			b2UUID := uuid.NewString()
			b2Devices := devList[len(devList)/2:]
			if len(b2Devices) == 0 {
				b2Devices = devList
			}
			b2Time := time.Now().Add(-6 * time.Hour)
			_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_upgrade_batches (
				id, batch_id, ota_package_id, batch_name, upgrade_strategy, status, target_device_count, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				batchSeq, b2UUID, pkgIDStr, "灰度二期-试点推广批次", "static", "in_progress", int32(len(b2Devices)), b2Time, time.Now(),
			)
			batchSeq++

			for idx, dID := range b2Devices {
				st := "in_progress"
				ver := "v1.0.0"
				if idx%2 == 1 {
					st = "pending"
				}
				ts := time.Now().Add(-time.Duration(idx*10+5) * time.Minute).Unix()
				_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_device_upgrade_status (
					id, device_id, ota_package_id, upgrade_batch_id, status, current_version, last_status_change_ts, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					statusSeq, dID, pkgIDStr, b2UUID, st, ver, ts, b2Time, time.Now(),
				)
				statusSeq++
			}

		case "success":
			// 全量批次全部成功
			bUUID := uuid.NewString()
			bTime := createdAt.Add(4 * time.Hour)
			_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_upgrade_batches (
				id, batch_id, ota_package_id, batch_name, upgrade_strategy, status, target_device_count, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				batchSeq, bUUID, pkgIDStr, "全量设备发布升级批次", "static", "completed", int32(len(devList)), bTime, bTime.Add(2*time.Hour),
			)
			batchSeq++

			for _, dID := range devList {
				ts := bTime.Add(30 * time.Minute).Unix()
				_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_device_upgrade_status (
					id, device_id, ota_package_id, upgrade_batch_id, status, current_version, last_status_change_ts, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					statusSeq, dID, pkgIDStr, bUUID, "success", seed.version, ts, bTime, bTime.Add(30*time.Minute),
				)
				statusSeq++
			}

		case "partial":
			// 部分成功批次
			bUUID := uuid.NewString()
			bTime := createdAt.Add(12 * time.Hour)
			_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_upgrade_batches (
				id, batch_id, ota_package_id, batch_name, upgrade_strategy, status, target_device_count, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				batchSeq, bUUID, pkgIDStr, "试点区域试运行批次", "static", "completed", int32(len(devList)), bTime, bTime.Add(3*time.Hour),
			)
			batchSeq++

			for idx, dID := range devList {
				st := "success"
				ver := seed.version
				if idx%2 == 0 {
					st = "failed"
					ver = "v1.0.0"
				}
				ts := bTime.Add(time.Duration(idx*20+10) * time.Minute).Unix()
				_, _ = deviceDB.Exec(`INSERT OR IGNORE INTO ota_device_upgrade_status (
					id, device_id, ota_package_id, upgrade_batch_id, status, current_version, last_status_change_ts, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					statusSeq, dID, pkgIDStr, bUUID, st, ver, ts, bTime, bTime.Add(time.Duration(idx*20+10)*time.Minute),
				)
				statusSeq++
			}
		}
	}

	fmt.Println("Seed complete!")
}
