.PHONY: all test build mock-device mock-alarm

# 运行所有微服务单元测试
test:
	@echo "🧪 Running unit tests across all microservices and shared packages..."
	@(cd pkg/protocol && go test ./...)
	@(cd pkg/tsdb && go test ./...)
	@(cd pkg/event && go test ./...)
	@(cd transport-mqtt && go test ./...)
	@(cd data-engine && go test ./...)
	@(cd backend && go test ./...)
	@echo "✅ All tests passed successfully!"

# 编译验证所有微服务二进制
build:
	@echo "🔨 Building all microservice binaries..."
	@(cd transport-mqtt && go build -buildvcs=false -o /dev/null ./cmd/server)
	@(cd data-engine && go build -buildvcs=false -o /dev/null ./cmd/server)
	@(cd backend && go build -buildvcs=false -o /dev/null ./cmd/server)
	@echo "✅ All binaries built cleanly!"

# 启动虚拟设备模拟器 (正常遥测 + 监听 OTA 升级)
mock-device:
	@go run ./scripts/mock_device/main.go -device sensor_test_01 -interval 3

# 启动虚拟设备模拟器 (触发高温告警 88.8°C)
mock-alarm:
	@go run ./scripts/mock_device/main.go -device sensor_test_01 -interval 3 -alarm=true
