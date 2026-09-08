.PHONY: init wire
init:
	go install github.com/google/wire/cmd/wire@latest

wire:
	wire ./cmd/server/wire

.PHONY: test
test:
	go test -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out -o coverage.html

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: docker
docker:
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/task -t 1.1.1.1:5000/demo-task:v1 .
	docker run --rm -i 1.1.1.1:5000/demo-task:v1

.PHONY: swag
swag:
	swag init  -g cmd/server/main.go -o ./docs
