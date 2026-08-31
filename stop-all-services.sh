#!/bin/bash

# ==============================================================================
# 0things IoT Platform - 一键优雅停止所有服务
# ==============================================================================

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_DIR="$PROJECT_DIR/.service-pids"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${BOLD}${YELLOW}🛑 正在优雅停止 0things 所有微服务与前端...${NC}"
echo ""

stop_service() {
    local name=$1
    local pid_file="$PID_DIR/${name}.pid"

    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p "$pid" > /dev/null 2>&1; then
            echo -ne "${BLUE}→ 正在停止 ${name} (PID: ${pid})...${NC} "
            kill "$pid" 2>/dev/null || true
            
            # 等待最多 5 秒优雅关闭
            local count=0
            while ps -p "$pid" > /dev/null 2>&1 && [ $count -lt 5 ]; do
                sleep 1
                count=$((count + 1))
            done

            # 强制清理
            if ps -p "$pid" > /dev/null 2>&1; then
                kill -9 "$pid" 2>/dev/null || true
            fi
            echo -e "${GREEN}✓ 已停止${NC}"
        else
            echo -e "${YELLOW}ℹ [${name}] 进程未在运行${NC}"
        fi
        rm -f "$pid_file"
    else
        echo -e "${YELLOW}ℹ [${name}] 未找到 PID 记录${NC}"
    fi
}

# 停止所有服务
stop_service "frontend"
stop_service "coap-transport"
stop_service "http-transport"
stop_service "mqtt-transport"
stop_service "data-engine"
stop_service "backend"

# 清理空 PID 目录
rm -rf "$PID_DIR"

echo ""
echo -e "${BOLD}${GREEN}✓ 0things 所有服务已成功停止。${NC}"
