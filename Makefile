.PHONY: all start stop restart status test build mock-device mock-alarm

# 一键启动所有微服务与前端
start:
	@bash ./start-all-services.sh

# 一键优雅停止所有服务
stop:
	@bash ./stop-all-services.sh

# 一键重启所有服务
restart: stop start

# 运行所有微服务单元测试
test:
	@echo "🧪 Running unit tests across all microservices and shared packages..."
	@(cd pkg/protocol && go test ./...)
	@(cd pkg/tsdb && go test ./...)
	@(cd mqtt-transport && go test ./internal/...)
	@(cd http-transport && go test ./internal/...)
	@(cd data-engine && go test ./internal/...)
	@(cd backend && go test ./...)
	@echo "✅ All tests passed successfully!"

# 编译验证所有微服务二进制
build:
	@echo "🔨 Building all microservice binaries..."
	@(cd mqtt-transport && go build -o /dev/null ./cmd/server)
	@(cd http-transport && go build -o /dev/null ./cmd/server)
	@(cd coap-transport && go build -o /dev/null ./cmd/server)
	@(cd data-engine && go build -o /dev/null ./cmd/server)
	@(cd backend && go build -o /dev/null ./cmd/server)
	@echo "✅ All binaries built cleanly!"

# 启动虚拟设备模拟器 (正常遥测 + 监听 OTA 升级)
mock-device:
	@go run ./scripts/mock_device/main.go -device sensor_test_01 -interval 3

# 启动虚拟设备模拟器 (触发高温告警 88.8°C)
mock-alarm:
	@go run ./scripts/mock_device/main.go -device sensor_test_01 -interval 3 -alarm=true
