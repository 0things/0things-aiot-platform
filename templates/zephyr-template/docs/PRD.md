可以。你这个需求不应该只是写一个“Zephyr MQTT Demo”，而应该把它定义成一个**可用于真实 IoT 产品的 Zephyr OS 工程脚手架**：传感器插件化、采集与上报解耦、Wi-Fi/MQTT 自动恢复、离线缓存、运行时配置、线程安全、低功耗扩展、测试和多板卡支持都提前设计好。

下面这份提示词可以直接丢给 **Claude Code / Codex / Cursor / OpenCode**。我已经按照你之前的使用场景，把“不同传感器不同采集周期 + 不同上传周期”作为核心架构要求。

```text
你是一名资深嵌入式架构师、Zephyr OS 内核工程师和 IoT 平台架构师。

请为我从 0 设计并实现一个：

# Production-Grade Zephyr IoT Application Scaffold

目标不是 Demo，而是一套可以长期用于多个 IoT 项目的“专业 Zephyr OS 应用脚手架”。

技术栈：

- Zephyr RTOS
- C
- CMake
- west
- Kconfig
- Devicetree
- Zephyr Settings
- Wi-Fi
- MQTT
- TLS
- JSON
- Flash/NVS
- Logging
- Work Queue
- Message Queue / FIFO
- 多传感器
- OTA 预留
- Watchdog 预留
- 低功耗预留

第一阶段主要支持：

- ESP32 / ESP32-S3
- 后续可以扩展 Nordic nRF70 / STM32 + Wi-Fi 模组

不要把代码写死到 ESP32。

==================================================
一、核心设计目标
==================================================

整个框架必须实现：

Sensor Collection
        ↓
Collector
        ↓
Local Cache
        ↓
Telemetry Queue
        ↓
Uploader
        ↓
MQTT
        ↓
Cloud

必须严格做到：

采集逻辑
≠
数据缓存
≠
MQTT 上传
≠
网络连接

几个模块必须完全解耦。

即：

传感器采集失败

不能影响：

Wi-Fi
MQTT
其他传感器
其他线程

MQTT 断开

不能影响：

传感器继续采集。

网络恢复以后：

自动补传缓存数据。

==================================================
二、核心需求：多 Collector 架构
==================================================

整个系统中存在多个 Collector。

例如：

temperature
humidity
imu
gps
fuel_level
battery
engine
custom_sensor

每一个 Collector 都必须拥有自己独立的：

1. collection_interval
2. upload_interval
3. enabled
4. cache_size
5. retry
6. priority

例如：

temperature:

collection_interval = 300s
upload_interval = 1800s

imu:

collection_interval = 10s
upload_interval = 600s

battery:

collection_interval = 60s
upload_interval = 300s

GPS:

collection_interval = 5s
upload_interval = 60s

必须支持：

采集频率 != 上传频率。

例如：

温度：

每 5 分钟采集一次

但是：

每 30 分钟上传一次。

在 30 分钟期间的数据：

必须存储到 Collector 自己的 buffer/cache。

上传时：

一次可以发送：

{
  "device_id": "...",
  "collector": "temperature",
  "timestamp": ...,
  "records": [
    ...
  ]
}

也可以设计成为：

Telemetry Batch。

请根据工程实践选择合理方案。

==================================================
三、Collector 抽象
==================================================

必须设计统一 Collector API。

例如：

struct collector;

struct collector_ops {
    int (*init)(struct collector *collector);
    int (*start)(struct collector *collector);
    int (*stop)(struct collector *collector);
    int (*collect)(struct collector *collector,
                   struct telemetry_record *record);
};

Collector 不应该知道：

MQTT
Wi-Fi
Cloud
Topic

Collector 只负责：

采集数据。

例如：

temperature_collector.c
imu_collector.c
gps_collector.c
battery_collector.c

必须支持：

collector_register(...)

或者：

STRUCT_SECTION_ITERABLE

等 Zephyr 风格的自动注册机制。

优先选择：

符合 Zephyr 原生设计风格的实现。

最终希望未来新增一个传感器，只需要增加：

xxx_collector.c

而不需要修改：

main.c
scheduler.c
mqtt.c

==================================================
四、调度系统
==================================================

不要写：

while (1) {
    k_sleep(...)
}

这种简单 Demo。

设计专业 Scheduler。

优先考虑：

k_work_delayable

每一个 Collector：

拥有独立 k_work_delayable。

例如：

collector_runtime {
    struct collector *collector;

    struct k_work_delayable collect_work;
    struct k_work_delayable upload_work;

    bool running;
};

要求：

每一个 Sensor：

独立 Collect Timer

独立 Upload Timer。

例如：

temperature

collect:
300s

upload:
1800s

IMU

collect:
10s

upload:
600s

禁止一个全局 timer 控制所有 Sensor。

==================================================
五、Collector 配置
==================================================

配置必须支持：

默认配置
+
运行时配置。

编译时：

Kconfig

例如：

CONFIG_APP_COLLECTOR_TEMP_INTERVAL=300
CONFIG_APP_COLLECTOR_TEMP_UPLOAD_INTERVAL=1800

但运行以后：

云端可以修改配置：

{
  "collectors": {
    "temperature": {
      "enabled": true,
      "collection_interval": 120,
      "upload_interval": 600
    },

    "imu": {
      "enabled": true,
      "collection_interval": 10,
      "upload_interval": 300
    }
  }
}

配置必须保存到：

Zephyr Settings subsystem

或者：

NVS。

重启以后：

仍然生效。

优先使用：

settings subsystem

如果底层使用 NVS：

由 settings 管理。

不要让业务直接依赖 NVS。

==================================================
六、Wi-Fi Manager
==================================================

独立设计：

wifi_manager

不要把 Wi-Fi 连接逻辑放：

main.c。

要求支持：

状态机：

DISCONNECTED
CONNECTING
CONNECTED
RECONNECTING
FAILED

必须使用 Zephyr：

net_mgmt

事件：

NET_EVENT_WIFI_CONNECT_RESULT
NET_EVENT_WIFI_DISCONNECT_RESULT
NET_EVENT_IPV4_ADDR_ADD

Wi-Fi Manager 提供：

wifi_manager_init()

wifi_manager_connect()

wifi_manager_disconnect()

wifi_manager_is_connected()

wifi_manager_wait_connected()

Wi-Fi 参数：

SSID
Password

通过：

Kconfig
或者
Settings

配置。

必须支持：

自动重连。

重连：

Exponential Backoff

例如：

1s
2s
4s
8s
16s
30s
60s

设置最大值。

不要：

Wi-Fi 一断就疯狂重连。

==================================================
七、MQTT Manager
==================================================

设计：

mqtt_manager

MQTT 与 Wi-Fi 完全解耦。

MQTT Manager 不直接控制 Wi-Fi。

只监听：

Network Ready Event。

状态机：

DISCONNECTED
CONNECTING
CONNECTED
RECONNECTING

要求：

MQTT 连接断开自动重连。

支持：

Client ID
Username
Password
Broker
Port

必须支持：

mqtt://

和：

mqtts://

TLS 必须作为正式能力实现。

不要只是 TODO。

支持：

CA Certificate

证书可以：

编译到 firmware

或者：

从 settings/storage 获取。

MQTT Manager API：

mqtt_manager_init()

mqtt_manager_start()

mqtt_manager_stop()

mqtt_manager_publish()

mqtt_manager_subscribe()

mqtt_manager_is_connected()

要求 MQTT publish：

禁止业务模块直接调用 Zephyr mqtt_publish()。

统一通过：

mqtt_manager_publish()

进行封装。

==================================================
八、MQTT Topic 设计
==================================================

设计专业 Topic：

devices/{device_id}/telemetry

devices/{device_id}/event

devices/{device_id}/state

devices/{device_id}/config

devices/{device_id}/command

或者：

devices/{device_id}/telemetry/{collector}

请比较两种设计。

最终选择一个适合：

多传感器 IoT 设备

的设计。

必须支持：

Cloud → Device

配置：

config

命令：

command

Device → Cloud

telemetry
event
state

==================================================
九、Telemetry Model
==================================================

设计统一：

telemetry_record

不要让每一个 Sensor 自己拼 JSON。

例如：

struct telemetry_record {
    int64_t timestamp;
    const char *collector;
    ...
};

但是必须解决：

不同 Sensor 字段不同的问题。

例如：

temperature:

temperature_c

IMU:

accel_x_mg
accel_y_mg
accel_z_mg

gyro_x_mdps
gyro_y_mdps
gyro_z_mdps

GPS:

latitude
longitude
altitude

Battery:

voltage_mv
battery_percent

请设计一个合理的数据模型。

可以考虑：

typed payload
union
TLV
key-value
CBOR

但第一版 Cloud Payload 使用：

JSON。

要求：

Sensor Driver
↓
Collector Data
↓
Telemetry Serializer
↓
JSON

Serializer 独立模块。

Collector 不负责 JSON。

==================================================
十、时间系统
==================================================

所有 telemetry 必须包含：

timestamp。

设计：

time_manager。

时间来源优先级：

1. SNTP
2. RTC
3. uptime fallback

Wi-Fi 建立以后：

自动进行 SNTP 时间同步。

必须能判断：

当前 wall clock 是否有效。

例如：

time_manager_is_synced()

time_manager_now_ms()

如果还没有同步：

数据仍然允许采集。

不要直接丢弃数据。

可以保存：

uptime

等时间恢复后处理。

==================================================
十一、本地缓存
==================================================

非常重要。

当：

Wi-Fi 断开

或者：

MQTT 断开

Sensor 必须继续采集。

数据进入：

local telemetry cache。

至少设计：

RAM Cache

进一步预留：

Flash Persistent Queue。

设计：

telemetry_store

API：

telemetry_store_append()

telemetry_store_peek()

telemetry_store_pop()

telemetry_store_count()

telemetry_store_clear()

第一版可以：

RAM ring buffer

但架构必须支持以后增加：

Flash Queue。

例如：

telemetry_store_ram.c

telemetry_store_flash.c

不要让 Collector 感知 storage 类型。

==================================================
十二、Upload Scheduler
==================================================

Upload Scheduler 与 Collect Scheduler 分开。

例如：

temperature：

5min collect

30min upload。

upload_work 到时间后：

检查：

MQTT connected？

YES：

将 Temperature Cache 打包。

publish。

publish 成功：

清除成功数据。

NO：

保留数据。

等待下一次：

upload interval

或者：

network reconnect trigger。

网络恢复以后：

可以触发一次：

flush。

必须注意：

不要瞬间把所有 Collector 数据一起发送造成 MQTT Flood。

设计：

upload queue
priority
rate limit

可以采用：

k_msgq

或：

k_fifo。

==================================================
十三、系统事件 Event Bus
==================================================

不要各个模块互相直接调用。

设计统一事件机制。

例如：

APP_EVENT_WIFI_CONNECTED

APP_EVENT_WIFI_DISCONNECTED

APP_EVENT_MQTT_CONNECTED

APP_EVENT_MQTT_DISCONNECTED

APP_EVENT_TIME_SYNCED

APP_EVENT_CONFIG_UPDATED

APP_EVENT_UPLOAD_REQUEST

APP_EVENT_SENSOR_ERROR

可以使用：

Zephyr zbus

优先考虑：

zbus。

如果不适合：

解释原因。

网络状态变化：

通过 Event Bus 通知其他模块。

==================================================
十四、远程配置
==================================================

Cloud 可以发送：

devices/{device_id}/config

例如：

{
  "version": 12,

  "collectors": {
    "temperature": {
      "enabled": true,
      "collection_interval_sec": 300,
      "upload_interval_sec": 1800
    },

    "imu": {
      "enabled": true,
      "collection_interval_sec": 10,
      "upload_interval_sec": 600
    }
  }
}

设备收到以后：

JSON Parse

验证：

范围

例如：

collection_interval >= 1

upload_interval >= collection_interval

或者：

允许 upload < collection。

你需要分析：

是否应该强制：

upload_interval >= collection_interval。

配置更新：

1 Validate
2 Apply
3 Persist
4 Reschedule

不能重启 MCU 才生效。

必须：

实时重新调度对应：

k_work_delayable。

==================================================
十五、Configuration Version
==================================================

配置必须具有：

version。

例如：

{
  "version": 13
}

设备保存：

last_config_version。

旧版本：

拒绝。

避免：

MQTT retained message

覆盖新的配置。

==================================================
十六、Device Identity
==================================================

实现：

device_identity

提供：

device_id

来源优先级：

1 MCU unique ID
2 MAC
3 Settings
4 Kconfig

API：

device_identity_get()

Cloud Topic：

全部使用这个 device_id。

==================================================
十七、错误处理
==================================================

禁止：

assert everywhere。

对生产代码：

可恢复错误：

返回 errno。

例如：

-EIO
-EINVAL
-ENOMEM
-ENOTCONN
-ETIMEDOUT

不可恢复：

才 assert。

每一个 subsystem：

有自己的 LOG_MODULE_REGISTER。

例如：

LOG_MODULE_REGISTER(wifi_manager)
LOG_MODULE_REGISTER(mqtt_manager)
LOG_MODULE_REGISTER(collector_manager)
LOG_MODULE_REGISTER(telemetry)

日志等级通过：

Kconfig。

==================================================
十八、Health Manager
==================================================

设计：

health_manager。

统计：

uptime
wifi reconnect count
mqtt reconnect count
collector errors
upload errors
queue size
dropped telemetry
free heap
stack usage

周期发送：

Device Health。

例如：

devices/{device_id}/state

JSON：

{
  "uptime": 12345,
  "wifi": "connected",
  "mqtt": "connected",
  "queue_size": 5,
  "dropped": 0
}

==================================================
十九、Watchdog
==================================================

预留：

watchdog_manager。

不要直接：

main feed watchdog。

设计：

各关键 subsystem heartbeat。

例如：

network
collector
uploader

Health Manager：

判断系统是否健康。

后期可接：

Zephyr watchdog。

==================================================
二十、低功耗设计
==================================================

架构必须考虑：

未来设备可能是 Battery Powered。

因此：

禁止：

高频 busy polling。

优先：

event-driven
k_work_delayable
interrupt
sleep

Collector：

如果没有任务：

不应该占 CPU。

设计需要兼容：

CONFIG_PM
CONFIG_PM_DEVICE

Sensor Driver：

可以 suspend。

第一版：

不要求完整 Deep Sleep。

但是架构必须能够未来加入。

==================================================
二十一、Thread Model
==================================================

请认真设计线程模型。

不要：

一个 Sensor 一个 Thread

除非确实需要。

优先使用：

System Work Queue
Custom Work Queue

建议：

Network Thread
MQTT Thread
Upload Work Queue
Collector Work Queue

或者：

合理利用 Zephyr networking 内部线程。

请给出：

线程表：

Name
Priority
Stack Size
Responsibilities

并解释：

哪些模块：

使用 Thread

哪些：

使用 k_work

哪些：

使用 k_work_delayable

哪些：

使用 k_msgq。

==================================================
二十二、内存管理
==================================================

嵌入式系统禁止随意：

malloc/free。

优先：

静态分配。

例如：

K_MSGQ_DEFINE
K_MEM_SLAB_DEFINE
Ring Buffer

JSON：

限制最大 payload。

例如：

CONFIG_APP_MQTT_PAYLOAD_MAX_SIZE

必须防止：

memory fragmentation。

如果必须动态分配：

解释为什么。

==================================================
二十三、Sensor Driver 和 Collector 分层
==================================================

必须区分：

Sensor Driver

和：

Collector。

例如：

Zephyr Sensor Driver：

sensor_sample_fetch()
sensor_channel_get()

Collector：

决定：

什么时候采集
采集哪些字段
如何转换业务单位
如何形成 telemetry_record

例如：

imu_collector

调用：

device_get_binding / DEVICE_DT_GET

然后：

sensor_sample_fetch()

sensor_channel_get()

不要自己实现底层 I2C。

==================================================
二十四、Devicetree
==================================================

传感器 hardware config：

必须使用：

Devicetree。

例如：

&i2c0 {

    status = "okay";

    sensor@68 {
        compatible = "...";
        reg = <0x68>;
    };
};

不要：

在 C 里面写：

I2C_ADDR 0x68

所有硬件信息：

Devicetree。

==================================================
二十五、Kconfig
==================================================

提供：

app/Kconfig

例如：

menu "Application"

config APP_WIFI
config APP_MQTT
config APP_COLLECTOR_TEMP
config APP_COLLECTOR_IMU

config APP_TEMP_COLLECTION_INTERVAL
config APP_TEMP_UPLOAD_INTERVAL

config APP_TELEMETRY_CACHE_SIZE

config APP_MQTT_PAYLOAD_MAX_SIZE

config APP_HEALTH_INTERVAL

endmenu

prj.conf：

只负责选择配置。

不要把业务配置散落：

#define。

==================================================
二十六、工程目录
==================================================

设计专业目录，例如：

iot-zephyr/
│
├── CMakeLists.txt
├── Kconfig
├── prj.conf
├── west.yml
│
├── boards/
│
├── dts/
│
├── include/
│   └── app/
│
├── src/
│
│   ├── main.c
│
│   ├── core/
│   │   ├── app.c
│   │   ├── app_event.c
│   │   ├── health_manager.c
│   │   ├── device_identity.c
│   │   └── time_manager.c
│
│   ├── network/
│   │   ├── wifi_manager.c
│   │   └── mqtt_manager.c
│
│   ├── collector/
│   │   ├── collector.c
│   │   ├── collector_manager.c
│   │   ├── temperature_collector.c
│   │   ├── imu_collector.c
│   │   └── battery_collector.c
│
│   ├── telemetry/
│   │   ├── telemetry.c
│   │   ├── telemetry_store.c
│   │   ├── telemetry_serializer.c
│   │   └── telemetry_uploader.c
│
│   ├── config/
│   │   ├── config_manager.c
│   │   └── remote_config.c
│
│   ├── storage/
│   │   └── settings_storage.c
│
│   └── utils/
│
├── tests/
│
├── samples/
│
└── docs/

可以根据 Zephyr 最佳实践调整目录。

但必须保持：

模块职责清楚。

==================================================
二十七、main.c 要求
==================================================

main.c 必须非常干净。

类似：

int main(void)
{
    app_init();
    app_start();

    return 0;
}

或者：

int main(void)
{
    return app_run();
}

禁止 main.c 出现：

Wi-Fi 连接细节
MQTT 连接细节
Sensor 逻辑
JSON
while loop
复杂业务。

==================================================
二十八、状态机
==================================================

Network：

INIT
WAIT_NETWORK
CONNECTED

MQTT：

DISCONNECTED
CONNECTING
CONNECTED

Collector：

DISABLED
IDLE
COLLECTING
ERROR

Uploader：

IDLE
WAIT_NETWORK
UPLOADING
RETRY

需要避免：

多个 bool：

wifi_ready
mqtt_ready
mqtt_connecting
mqtt_failed

造成状态混乱。

优先：

enum state。

==================================================
二十九、MQTT QoS
==================================================

设计合理策略。

Telemetry：

QoS 0 / QoS 1

Config：

QoS 1

Command：

QoS 1

State：

QoS 1

分析：

Telemetry 使用 QoS 0 还是 QoS 1。

生产系统优先：

可配置。

==================================================
三十、Retained Message
==================================================

分析：

config

是否：

retain = true

state：

是否：

retain = true

telemetry：

retain = false

给出推荐。

==================================================
三十一、LWT
==================================================

MQTT 必须实现：

Last Will Testament。

例如：

devices/{device_id}/presence

连接：

{
  "online": true
}

LWT：

{
  "online": false
}

==================================================
三十二、安全
==================================================

不要：

LOG_INF("password=%s")

任何日志禁止输出：

Wi-Fi Password
MQTT Password
Token
TLS Key

MQTT TLS：

必须验证 Server Certificate。

禁止：

CERT_NONE。

如果使用：

PSK
Client Certificate

预留扩展。

==================================================
三十三、测试
==================================================

使用：

Zephyr ztest。

至少实现：

collector scheduler test

config validation test

telemetry serialization test

ring buffer test

config version test

mqtt topic test

mock collector test

重点测试：

temperature:

collect interval = 5 min

upload interval = 30 min

确认：

30min 内：

collect 6 次。

upload：

1 次。

==================================================
三十四、Native Simulator
==================================================

尽量让：

Collector
Scheduler
Telemetry
Config

可以通过：

native_sim

测试。

不要让：

所有模块

都依赖：

ESP32。

==================================================
三十五、示例 Collector
==================================================

必须实现三个完整示例：

1 Temperature Collector

字段：

temperature_c

2 IMU Collector

字段：

accel_x_mg
accel_y_mg
accel_z_mg

gyro_x_mdps
gyro_y_mdps
gyro_z_mdps

3 Battery Collector

字段：

voltage_mv
battery_percent

如果没有真实 sensor：

允许使用：

mock sensor

并通过：

CONFIG_APP_MOCK_SENSOR

启用。

==================================================
三十六、JSON 示例
==================================================

Telemetry：

{
  "device_id": "device-001",
  "timestamp": 1780000000000,
  "collector": "imu",
  "records": [
    {
      "timestamp": 1780000000000,
      "accel_x_mg": 10,
      "accel_y_mg": 20,
      "accel_z_mg": 1001,
      "gyro_x_mdps": 100,
      "gyro_y_mdps": 200,
      "gyro_z_mdps": 300
    }
  ]
}

你可以提出更专业的格式。

==================================================
三十七、代码质量
==================================================

代码必须：

C11

Zephyr Coding Style

统一：

snake_case

所有公共接口：

.h

内部函数：

static

模块封装良好。

每一个 public function：

写清楚：

parameters
return
thread safety

不要：

God Object。

不要：

main.c 1000 行。

不要：

app.c 2000 行。

==================================================
三十八、并发安全
==================================================

明确说明每个：

API

是否：

Thread Safe。

Collector Cache：

可能同时被：

Collector Work

Uploader

访问。

必须：

mutex
spinlock
message queue

之一进行保护。

不能存在：

Race Condition。

==================================================
三十九、断网场景
==================================================

必须完整考虑：

Boot
↓
WiFi unavailable
↓
Collectors Start
↓
collect data
↓
cache
↓
WiFi reconnect
↓
MQTT reconnect
↓
flush cache

绝对不能要求：

有网络以后

才启动：

Sensor。

==================================================
四十、MQTT Broker 不可用
==================================================

Wi-Fi：

Connected

但是：

Broker：

Down。

Collector：

仍然工作。

Telemetry：

继续缓存。

MQTT：

独立 retry。

==================================================
四十一、Config 无效
==================================================

Cloud：

发送：

{
  "collection_interval_sec": -1
}

不能 crash。

返回：

config rejected。

发送 event：

devices/{device_id}/event

例如：

{
  "type": "config_rejected",
  "reason": "invalid collection interval"
}

==================================================
四十二、数据溢出策略
==================================================

如果：

Offline 3 days

Buffer 满。

必须定义：

Drop Oldest

或者：

Drop Newest。

默认推荐：

Drop Oldest。

同时：

dropped_records++

Health 中：

上报。

==================================================
四十三、未来能力预留
==================================================

架构需要预留：

BLE
LTE
Ethernet

Transport：

不能完全写死 Wi-Fi。

最好设计：

network_manager

Wi-Fi：

只是：

network backend。

但第一版不要过度工程化。

需要判断：

哪些抽象现在值得做。

==================================================
四十四、OTA
==================================================

预留：

ota_manager。

支持未来：

MCUboot

Zephyr DFU

MQTT Command：

firmware_update

第一版：

可以只有 interface。

不要实现完整 OTA。

==================================================
四十五、CLI / Shell
==================================================

支持：

Zephyr Shell。

命令：

app status

app collectors

app collector temperature

app config

wifi status

mqtt status

telemetry status

例如：

uart console：

uart:~$ app status

输出：

WiFi: connected
MQTT: connected
Time: synced

Collectors:

temperature
collect: 300s
upload: 1800s
records: 4

imu
collect: 10s
upload: 600s
records: 32

这个功能非常重要。

==================================================
四十六、README
==================================================

最终 README 必须包含：

Architecture

Build

Flash

Configuration

MQTT

Sensor

Remote Config

Testing

Adding Collector

Adding Board

FAQ

例如：

west build -b esp32s3_devkitc/esp32s3/procpu app

或者针对当前 Zephyr board qualifier：

给出正确命令。

==================================================
四十七、Architecture 文档
==================================================

创建：

docs/architecture.md

必须包括：

模块图

数据流图

线程模型

状态机

Collector 生命周期

MQTT reconnect

Remote Config

Offline Cache。

使用：

Mermaid。

例如：

flowchart TD

Sensor --> Collector
Collector --> Store
Store --> Uploader
Uploader --> MQTT
MQTT --> Cloud

==================================================
四十八、设计原则
==================================================

优先：

Simple
Explicit
Testable
Event Driven
Static Allocation
Failure Isolation

不要为了：

设计模式

而设计模式。

避免：

过度抽象。

例如：

不要写：

IoTAbstractGenericTransportFactoryManagerInterface。

接口命名保持：

简单。

==================================================
四十九、实施步骤
==================================================

不要一次性乱写所有代码。

严格按照：

Phase 1
Architecture

Phase 2
Project Skeleton

Phase 3
Core

Phase 4
Collector Framework

Phase 5
Telemetry

Phase 6
Wi-Fi

Phase 7
MQTT

Phase 8
Remote Config

Phase 9
Storage

Phase 10
Tests

Phase 11
Documentation

逐步实施。

每一个 Phase：

完成以后：

确保：

可以编译。

禁止：

最后才编译。

==================================================
五十、最终交付
==================================================

最终必须给我：

1 完整目录结构

2 Architecture 说明

3 所有核心 .h

4 所有核心 .c

5 Kconfig

6 CMakeLists.txt

7 prj.conf

8 overlay 示例

9 Wi-Fi Manager

10 MQTT Manager

11 Collector Framework

12 Scheduler

13 Telemetry Model

14 Telemetry Store

15 Telemetry Uploader

16 Remote Config

17 Settings Persistence

18 Time Manager

19 Device Identity

20 Event Bus

21 Health Manager

22 Zephyr Shell

23 ztest

24 README.md

25 docs/architecture.md

26 sample collector

27 mock collector

28 sample MQTT payload

29 sample cloud config

30 ESP32-S3 build example

==================================================
关键要求
==================================================

这是一个：

Production Scaffold。

不是：

Tutorial。

不要只给伪代码。

需要提供：

可以真正 build 的代码。

如果某个 Zephyr API：

存在版本差异，

优先采用：

当前 Zephyr 4.x API。

如果不确定：

不要编造 API。

先检查：

Zephyr API

然后再实现。

代码优先：

Zephyr Native API。

避免：

自己重新实现 Zephyr 已经提供的能力。

第一阶段重点保证：

Architecture Correctness

而不是：

堆很多 Feature。

==================================================
最重要的一条设计要求
==================================================

系统必须完美支持：

Sensor A

collect = 10 sec
upload = 10 min

Sensor B

collect = 5 min
upload = 30 min

Sensor C

collect = 1 sec
upload = 5 min

三个 Sensor：

独立运行。

互不影响。

MQTT Down：

继续 Collect。

Wi-Fi Down：

继续 Collect。

Network Recovery：

自动 Upload。

Remote Config：

实时修改：

collect interval

upload interval

无需：

Reboot。

这应该成为整个框架设计的核心。
```

我建议你的核心架构最终收敛成下面这个模型，而不是“每个传感器自己定时 + 自己 MQTT”：

```text
                         ┌──────────────────┐
                         │  Remote Config   │
                         └────────┬─────────┘
                                  │
                                  ▼
                         ┌──────────────────┐
                         │ Config Manager   │
                         │ settings / NVS   │
                         └────────┬─────────┘
                                  │
                  reschedule      │
          ┌───────────────────────┼──────────────────────┐
          ▼                       ▼                      ▼
  ┌──────────────┐        ┌──────────────┐       ┌──────────────┐
  │ Temperature  │        │     IMU      │       │   Battery    │
  │ Collector    │        │ Collector    │       │ Collector    │
  │              │        │              │       │              │
  │ collect 5m   │        │ collect 10s  │       │ collect 1m   │
  │ upload 30m   │        │ upload 10m   │       │ upload 5m    │
  └──────┬───────┘        └──────┬───────┘       └──────┬───────┘
         │                       │                      │
         └───────────────────────┼──────────────────────┘
                                 ▼
                     ┌─────────────────────┐
                     │  Telemetry Store    │
                     │                     │
                     │ RAM Ring Buffer     │
                     │ Flash Queue later   │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │ Telemetry Uploader  │
                     │ batch / retry       │
                     │ rate limit          │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │    MQTT Manager     │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │   Network Manager   │
                     │       Wi-Fi         │
                     └─────────────────────┘
```

其中我特别建议你坚持 **`Collector ≠ MQTT`**。以后你加温度、IMU、油位、GPS、发动机 CAN、Modbus 等采集器时，只需要实现一个新的 `xxx_collector.c`，完全不用碰 MQTT 和 Wi-Fi。这才是这套脚手架后面能长期复用的关键。

另外，你之前的实际业务类似 **温度 5 分钟采集 / 30 分钟上传、IMU 10 秒采集 / 10 分钟上传**，这套设计正好可以原生覆盖，不需要为不同设备再重新写调度逻辑。
