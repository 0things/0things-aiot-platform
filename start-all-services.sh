#!/bin/bash

# 0things - Start All Services Script (Local Development)
# Starts frontend, backend, and telemetry-service without Docker

set -e

echo "🚀 Starting all 0things services..."
echo ""

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"
TELEMETRY_SERVICE_DIR="$PROJECT_DIR/telemetry-service"
FRONTEND_DIR="$PROJECT_DIR/frontend"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

print_service() {
    echo -e "${BLUE}→${NC} $1"
}

# Check if services exist
if [ ! -d "$BACKEND_DIR" ]; then
    print_error "backend directory not found"
    exit 1
fi

if [ ! -d "$TELEMETRY_SERVICE_DIR" ]; then
    print_error "telemetry-service directory not found"
    exit 1
fi

if [ ! -d "$FRONTEND_DIR" ]; then
    print_error "frontend directory not found"
    exit 1
fi

# Create logs directory
mkdir -p "$PROJECT_DIR/.service-logs"

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_info "Shutting down services..."

    pkill -f "go run.*aiot-backend/cmd" || true
    pkill -f "go run.*telemetry-service/cmd" || true
    pkill -f "vite" || true

    sleep 1
    print_status "All services stopped"
}

trap cleanup EXIT

# 1. Start backend
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print_service "Starting backend..."
cd "$BACKEND_DIR"

# Kill any existing process
pkill -f "go run.*aiot-backend/cmd" || true
pkill -f "aiot-backend" || true
sleep 1

# Start in background with output to log file and terminal
LOG_FILE="$PROJECT_DIR/.service-logs/backend.log"
{
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting backend..."
    go run ./cmd/server -conf ./config
} > "$LOG_FILE" 2>&1 &

BACKEND_PID=$!
sleep 2

# Check if process started
if ps -p $BACKEND_PID > /dev/null 2>&1; then
    print_status "backend started (PID: $BACKEND_PID)"
    echo "  Log: $LOG_FILE"
else
    print_error "backend failed to start"
    cat "$LOG_FILE" | tail -20
    exit 1
fi

# 2. Start telemetry-service
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print_service "Starting telemetry-service..."
cd "$TELEMETRY_SERVICE_DIR"

# Kill any existing process
pkill -f "go run.*telemetry-service/cmd" || true
pkill -f "telemetry-service" || true
sleep 1

# Start in background with output to log file and terminal
LOG_FILE="$PROJECT_DIR/.service-logs/telemetry-service.log"
{
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting telemetry-service..."
    go run ./cmd/telemetry-service -conf ./configs
} > "$LOG_FILE" 2>&1 &

TELEMETRY_PID=$!
sleep 2

# Check if process started
if ps -p $TELEMETRY_PID > /dev/null 2>&1; then
    print_status "telemetry-service started (PID: $TELEMETRY_PID)"
    echo "  Log: $LOG_FILE"
else
    print_error "telemetry-service failed to start"
    cat "$LOG_FILE" | tail -20
    exit 1
fi

# 3. Start frontend
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print_service "Starting frontend dev server..."
cd "$FRONTEND_DIR"

# Kill any existing Vite process
pkill -f "vite" || true
sleep 1

# Start in background with output to log file
LOG_FILE="$PROJECT_DIR/.service-logs/frontend.log"
{
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting frontend dev server..."
    pnpm run dev
} > "$LOG_FILE" 2>&1 &

FRONTEND_PID=$!
sleep 3

# Check if process started
if ps -p $FRONTEND_PID > /dev/null 2>&1; then
    print_status "Frontend dev server started (PID: $FRONTEND_PID)"
    echo "  Log: $LOG_FILE"
else
    print_error "Frontend failed to start"
    cat "$LOG_FILE" | tail -20
    exit 1
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ All services started successfully!${NC}"
echo ""
echo "Service Status:"
echo "  • backend:            (running in background)"
echo "  • telemetry-service:  (running in background)"
echo "  • Frontend:           http://localhost:5173"
echo ""
echo "Process IDs:"
echo "  • backend:            $BACKEND_PID"
echo "  • telemetry-service:  $TELEMETRY_PID"
echo "  • Frontend:           $FRONTEND_PID"
echo ""
echo "Logs:"
echo "  • $PROJECT_DIR/.service-logs/backend.log"
echo "  • $PROJECT_DIR/.service-logs/telemetry-service.log"
echo "  • $PROJECT_DIR/.service-logs/frontend.log"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Keep the script running
wait
