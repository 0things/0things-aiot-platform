#!/bin/bash

# ==============================================================================
# 0things IoT Platform - 一键启动所有微服务与前端
# ==============================================================================

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_DIR="$PROJECT_DIR/.service-pids"
LOG_DIR="$PROJECT_DIR/storage/logs"

mkdir -p "$PID_DIR"
mkdir -p "$LOG_DIR"

# 颜色输出定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${BOLD}${CYAN}==============================================================================${NC}"
echo -e "${BOLD}${CYAN}           🚀 正在启动 0things AIoT 平台全矩阵微服务与前端...               ${NC}"
echo -e "${BOLD}${CYAN}==============================================================================${NC}"
echo ""

# 检查单个进程是否已在运行
is_running() {
    local service_name=$1
    local pid_file="$PID_DIR/${service_name}.pid"
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        fi
    fi
    return 1
}

# 启动单个 Go 微服务
start_go_service() {
    local name=$1
    local dir="$PROJECT_DIR/$1"
    local conf="config/local.yml"
    local pid_file="$PID_DIR/${name}.pid"
    local log_file="$LOG_DIR/${name}.log"

    if is_running "$name"; then
        local pid=$(cat "$pid_file")
        echo -e "${YELLOW}ℹ [${name}] 已在运行中 (PID: ${pid})${NC}"
        return 0
    fi

    echo -ne "${BLUE}→ 正在启动 ${name}...${NC} "
    cd "$dir"
    go run ./cmd/server -conf "./$conf" > "$log_file" 2>&1 &
    local pid=$!
    echo "$pid" > "$pid_file"
    sleep 1

    if ps -p "$pid" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 启动成功 (PID: ${pid})${NC}"
    else
        echo -e "${RED}✗ 启动失败，请检查日志: ${log_file}${NC}"
    fi
    cd "$PROJECT_DIR"
}

# 启动前端 Vite 服务
start_frontend() {
    local name="frontend"
    local dir="$PROJECT_DIR/frontend"
    local pid_file="$PID_DIR/${name}.pid"
    local log_file="$LOG_DIR/${name}.log"

    if is_running "$name"; then
        local pid=$(cat "$pid_file")
        echo -e "${YELLOW}ℹ [${name}] 前端已在运行中 (PID: ${pid})${NC}"
        return 0
    fi

    echo -ne "${BLUE}→ 正在启动 ${name} (Vite Dev Server)...${NC} "
    cd "$dir"
    pnpm dev > "$log_file" 2>&1 &
    local pid=$!
    echo "$pid" > "$pid_file"
    sleep 2

    if ps -p "$pid" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 启动成功 (PID: ${pid})${NC}"
    else
        echo -e "${RED}✗ 启动失败，请检查日志: ${log_file}${NC}"
    fi
    cd "$PROJECT_DIR"
}

# 1. 依次启动核心微服务
start_go_service "backend"
start_go_service "data-engine"
start_go_service "mqtt-transport"
start_go_service "http-transport"
start_go_service "coap-transport"

# 2. 启动前端控制台
if [ -d "$PROJECT_DIR/frontend" ]; then
    start_frontend
fi

echo ""
echo -e "${BOLD}${GREEN}==============================================================================${NC}"
echo -e "${BOLD}${GREEN}  ✨ 0things 全部服务已就绪！控制台与接口访问导航：                          ${NC}"
echo -e "${BOLD}${GREEN}==============================================================================${NC}"
echo -e "  🌐 ${BOLD}前端管理控制台 (Web UI)${NC}:    ${CYAN}http://localhost:5173${NC}"
echo -e "  📡 ${BOLD}后端管理 API (REST Server)${NC}:    ${CYAN}http://localhost:8000${NC}"
echo -e "  📖 ${BOLD}Swagger API 交互文档${NC}:       ${CYAN}http://localhost:8000/swagger/index.html${NC}"
echo -e "  ⚡ ${BOLD}HTTP 设备协议网关${NC}:          ${CYAN}http://localhost:8081${NC}"
echo -e "  🔌 ${BOLD}MQTT 设备协议网关${NC}:          ${CYAN}tcp://localhost:1883${NC}"
echo -e "  📶 ${BOLD}CoAP 低功耗网关 (UDP)${NC}:       ${CYAN}coap://localhost:5683${NC}"
echo -e "  🧠 ${BOLD}数据计算与任务中心${NC}:          ${CYAN}data-engine (后台运行)${NC}"
echo -e "------------------------------------------------------------------------------"
echo -e "  📂 实时运行日志目录:  ${YELLOW}${LOG_DIR}/${NC}"
echo -e "  🛑 停止全部服务命令:  ${YELLOW}./stop-all-services.sh${NC}"
echo -e "${BOLD}${GREEN}==============================================================================${NC}"
