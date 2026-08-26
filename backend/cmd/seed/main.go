package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	userDB.Exec("DELETE FROM organization_users")
	userDB.Exec("DELETE FROM organizations")
	userDB.Exec("DELETE FROM users")
	deviceDB.Exec("DELETE FROM device_events")
	deviceDB.Exec("DELETE FROM device_tags")
	deviceDB.Exec("DELETE FROM device_shadow_histories")
	deviceDB.Exec("DELETE FROM device_shadows")
	deviceDB.Exec("DELETE FROM device_states")
	deviceDB.Exec("DELETE FROM devices")
	deviceDB.Exec("DELETE FROM products")
	deviceDB.Exec("DELETE FROM product_ts_ls")
	deviceDB.Exec("DELETE FROM ota_device_upgrade_status")
	deviceDB.Exec("DELETE FROM ota_upgrade_batches")
	deviceDB.Exec("DELETE FROM ota_packages")
	userDB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('users','organizations','organization_users')")
	deviceDB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('device_events','device_tags','device_shadow_histories','device_shadows','device_states','devices','products','product_ts_ls','ota_packages','ota_upgrade_batches','ota_device_upgrade_status')")

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
	categories := []string{"传感器", "网关", "摄像头", "温控器", "电表", "水表", "气体检测", "门锁", "开关", "插座"}
	statuses := []string{"active", "inactive", "draft"}
	nodeTypes := []string{"device", "gateway", "sub_device"}
	connMethods := []string{"wifi", "ble", "zigbee", "lora", "4g", "ethernet"}
	protocols := []string{"mqtt", "http", "coap", "modbus", "tcp"}

	productOrgMap := make(map[int]int64)
	for i := 1; i <= 50; i++ {
		orgID := int64(rand.Intn(3) + 1)
		productOrgMap[i] = orgID
		metadata, _ := json.Marshal(map[string]interface{}{
			"manufacturer": fmt.Sprintf("厂商%d", rand.Intn(10)+1),
			"model":        fmt.Sprintf("MODEL-%c%c%c", 'A'+rand.Intn(26), 'A'+rand.Intn(26), 'A'+rand.Intn(26)),
		})
		_, err := deviceDB.Exec(`INSERT OR IGNORE INTO products (product_key, name, description, category, status, metadata, node_type, connectivity_method, access_protocol, organization_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("pk_product_%03d", i),
			fmt.Sprintf("产品%d", i),
			fmt.Sprintf("这是产品%d的描述信息", i),
			categories[rand.Intn(len(categories))],
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
			time.Now().Add(-time.Duration(rand.Intn(30))*24*time.Hour),
			time.Now(),
		)
		if err != nil {
			log.Printf("product_tsl insert error: %v", err)
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
			seed.desc,
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
