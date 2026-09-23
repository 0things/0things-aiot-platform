#!/usr/bin/env bash

# Starts infrastructure with Docker and application services on the host.

set -o pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$PROJECT_DIR/.service-pids"
LOG_DIR="$PROJECT_DIR/storage/logs"

mkdir -p "$PID_DIR" "$LOG_DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

is_running() {
    local name="$1"
    local pid_file="$PID_DIR/${name}.pid"

    [ -f "$pid_file" ] || return 1

    local pid
    pid=$(<"$pid_file")
    if kill -0 "$pid" 2>/dev/null; then
        return 0
    fi

    rm -f "$pid_file"
    return 1
}

is_port_in_use() {
    local port="$1"
    lsof -nP -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1
}

is_process_in_directory() {
    local dir="$1"
    shift

    local pid cwd
    while read -r pid; do
        [ -n "$pid" ] || continue
        cwd=$(lsof -a -p "$pid" -d cwd -Fn 2>/dev/null | sed -n 's/^n//p')
        if [ "$cwd" = "$dir" ]; then
            return 0
        fi
    done < <(pgrep -f "$*" 2>/dev/null || true)

    return 1
}

check_dependency() {
    local name="$1"
    local port="$2"

    if is_port_in_use "$port"; then
        echo -e "${GREEN}✓ ${name} 已就绪 (localhost:${port})${NC}"
    else
        echo -e "${YELLOW}⚠ ${name} 未监听 localhost:${port}；相关服务可能无法连接${NC}"
    fi
}

start_infrastructure() {
    local -a compose_command

    if docker compose version >/dev/null 2>&1; then
        compose_command=(docker compose)
    elif command -v docker-compose >/dev/null 2>&1; then
        compose_command=(docker-compose)
    else
        echo -e "${RED}✗ 未找到 Docker Compose，请先安装 Docker Compose。${NC}"
        exit 1
    fi

    echo -e "${BOLD}启动 Docker 基础设施：${NC}"
    "${compose_command[@]}" -f "$PROJECT_DIR/docker-compose.yml" up -d \
        postgres redis nats emqx tdengine logto-postgres logto
}

start_service() {
    local name="$1"
    local dir="$2"
    local port="$3"
    shift 3

    local pid_file="$PID_DIR/${name}.pid"
    local log_file="$LOG_DIR/${name}.log"

    if is_running "$name"; then
        echo -e "${YELLOW}ℹ [${name}] 已由本脚本启动 (PID: $(<"$pid_file"))${NC}"
        return
    fi

    if is_process_in_directory "$dir" "$@"; then
        echo -e "${YELLOW}ℹ [${name}] 已有本地进程运行，跳过启动${NC}"
        return
    fi

    if [ -n "$port" ] && is_port_in_use "$port"; then
        echo -e "${YELLOW}ℹ [${name}] 端口 ${port} 已被其他进程占用，跳过启动${NC}"
        return
    fi

    if [ ! -d "$dir" ]; then
        echo -e "${RED}✗ [${name}] 目录不存在: ${dir}${NC}"
        return
    fi

    echo -ne "${BLUE}→ 启动 ${name}（本地）...${NC} "
    (
        cd "$dir" || exit 1
        exec nohup "$@"
    ) >>"$log_file" 2>&1 &
    local pid=$!
    echo "$pid" >"$pid_file"

    # Allow go run/pnpm to compile before deciding whether startup failed.
    local count=0
    while kill -0 "$pid" 2>/dev/null && [ "$count" -lt 5 ]; do
        sleep 1
        count=$((count + 1))
    done
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "${GREEN}✓ 已启动 (PID: ${pid})${NC}"
    else
        rm -f "$pid_file"
        echo -e "${RED}✗ 启动失败，请查看 ${log_file}${NC}"
    fi
}

echo -e "${BOLD}${CYAN}启动 0things：Docker 基础设施 + 本地应用服务${NC}"
echo

start_infrastructure
echo

echo -e "${BOLD}检查本地基础设施：${NC}"
check_dependency "PostgreSQL" 5432
check_dependency "Redis" 6379
check_dependency "NATS" 4222
check_dependency "EMQX" 1883
check_dependency "TDengine" 6041
check_dependency "Logto" 3001
echo

start_service "backend" "$PROJECT_DIR/backend" 8000 go run ./cmd/server -conf config/local.yml
start_service "data-engine" "$PROJECT_DIR/data-engine" "" go run ./cmd/server -conf config/local.yml
start_service "transport-mqtt" "$PROJECT_DIR/transport-mqtt" "" go run ./cmd/server -conf config/local.yml
start_service "mcp-server" "$PROJECT_DIR/backend" 8009 go run ./cmd/mcp -conf config/local.yml

if [ -d "$PROJECT_DIR/ai-copilot" ]; then
    start_service "ai-copilot" "$PROJECT_DIR/ai-copilot" 8005 pnpm dev
fi

if [ -d "$PROJECT_DIR/frontend" ]; then
    start_service "frontend" "$PROJECT_DIR/frontend" 5173 pnpm dev
fi

echo
echo -e "${BOLD}${GREEN}本地应用服务启动命令已提交。${NC}"
echo -e "  前端:      ${CYAN}http://localhost:5173${NC}"
echo -e "  后端 API:  ${CYAN}http://localhost:8000${NC}"
echo -e "  服务日志:  ${YELLOW}${LOG_DIR}/${NC}"
echo -e "  停止服务:  ${YELLOW}./scripts/stop-all-services.sh${NC}"
