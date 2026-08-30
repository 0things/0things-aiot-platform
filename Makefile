.PHONY: all start stop restart status test build

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
	@echo "🧪 Running unit tests across all microservices..."
	@(cd mqtt-transport && go test ./internal/...)
	@(cd http-transport && go test ./internal/...)
	@(cd data-engine && go test ./internal/...)
	@(cd backend && go test ./internal/service/... ./internal/handler/...)
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

