#!/bin/bash

# 0things - Stop All Services Script

echo "🛑 Stopping all 0things services..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Stop backend
if pkill -f "go run.*0things-backend/cmd" 2>/dev/null; then
    print_status "Stopped backend"
else
    echo "  backend not running"
fi

# Stop telemetry-service
if pkill -f "go run.*telemetry-service/cmd" 2>/dev/null; then
    print_status "Stopped telemetry-service"
else
    echo "  telemetry-service not running"
fi

# Stop frontend
if pkill -f "vite" 2>/dev/null; then
    print_status "Stopped frontend dev server"
else
    echo "  Frontend not running"
fi

sleep 1

echo ""
echo -e "${GREEN}✅ All services stopped${NC}"
