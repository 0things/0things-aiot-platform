#!/bin/bash

# Continuous test script for telemetry and events publishing
# This script continuously sends telemetry events and device events to the HTTP API

API_URL="${API_URL:-http://localhost:8013}"
TELEMETRY_ENDPOINT="$API_URL/api/v1/telemetry"
EVENTS_ENDPOINT="$API_URL/api/v1/events"
INTERVAL="${INTERVAL:-2}"  # seconds between requests

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     Telemetry Service - Continuous Test Script            ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}API URL:${NC}              $API_URL"
echo -e "${BLUE}Telemetry Endpoint:${NC}   $TELEMETRY_ENDPOINT"
echo -e "${BLUE}Events Endpoint:${NC}      $EVENTS_ENDPOINT"
echo -e "${BLUE}Request Interval:${NC}     ${INTERVAL}s"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Counter for statistics
telemetry_success=0
telemetry_failed=0
events_success=0
events_failed=0
total_requests=0

# Trap Ctrl+C to show statistics
trap 'show_stats' INT

show_stats() {
    echo ""
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                    Test Statistics                         ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo -e "${BLUE}Total Requests:${NC}          $total_requests"
    echo -e "${BLUE}Telemetry Success:${NC}       ${GREEN}$telemetry_success${NC}"
    echo -e "${BLUE}Telemetry Failed:${NC}        ${RED}$telemetry_failed${NC}"
    echo -e "${BLUE}Events Success:${NC}          ${GREEN}$events_success${NC}"
    echo -e "${BLUE}Events Failed:${NC}           ${RED}$events_failed${NC}"
    echo ""
    exit 0
}

# Array of telemetry event types
TELEMETRY_EVENTS=(
    "metric.cpu"
    "metric.memory"
    "metric.network"
    "metric.temperature"
    "metric.humidity"
    "sensor.read"
    "data.upload"
)

# Array of device event types
DEVICE_EVENTS=(
    "device.online"
    "device.offline"
    "device.alert"
    "device.status"
    "device.heartbeat"
    "device.upgrade"
)

# Function to generate random product key
random_product_key() {
    echo "prod-$(uuidgen | cut -d'-' -f1 | tr '[:upper:]' '[:lower:]')"
}

# Function to generate random device key
random_device_key() {
    echo "dev-$(uuidgen | cut -d'-' -f1 | tr '[:upper:]' '[:lower:]')"
}

# Function to generate random session ID
random_session_id() {
    echo "session-$(uuidgen | cut -d'-' -f1 | tr '[:upper:]' '[:lower:]')"
}

# Function to send telemetry event
send_telemetry() {
    local event_type="${TELEMETRY_EVENTS[$RANDOM % ${#TELEMETRY_EVENTS[@]}]}"
    local product_key=$(random_product_key)
    local device_key=$(random_device_key)
    local session_id=$(random_session_id)

    local payload=$(cat <<EOF
{
  "type": "$event_type",
  "product_key": "$product_key",
  "device_key": "$device_key",
  "session_id": "$session_id",
  "data": {
    "value": $((RANDOM % 100 + 1)),
    "unit": "celsius"
  }
}
EOF
)

    local response=$(curl -s -w "\n%{http_code}" -X POST "$TELEMETRY_ENDPOINT" \
        -H "Content-Type: application/json" \
        -d "$payload")

    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n-1)

    ((total_requests++))

    if [ "$http_code" = "200" ]; then
        ((telemetry_success++))
        echo -e "[${GREEN}✓${NC}] Telemetry: ${BLUE}$event_type${NC} (product: ${product_key:0:12}..., device: ${device_key:0:12}...)"
    else
        ((telemetry_failed++))
        echo -e "[${RED}✗${NC}] Telemetry failed: HTTP $http_code - $body"
    fi
}

# Function to send device event
send_event() {
    local event_type="${DEVICE_EVENTS[$RANDOM % ${#DEVICE_EVENTS[@]}]}"
    local product_key=$(random_product_key)
    local device_key=$(random_device_key)

    local payload=$(cat <<EOF
{
  "type": "$event_type",
  "product_key": "$product_key",
  "device_key": "$device_key",
  "data": {
    "status": "online",
    "battery_level": $((RANDOM % 100 + 1)),
    "firmware_version": "1.0.$((RANDOM % 10))"
  }
}
EOF
)

    local response=$(curl -s -w "\n%{http_code}" -X POST "$EVENTS_ENDPOINT" \
        -H "Content-Type: application/json" \
        -d "$payload")

    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n-1)

    ((total_requests++))

    if [ "$http_code" = "200" ]; then
        ((events_success++))
        echo -e "[${GREEN}✓${NC}] Event:     ${BLUE}$event_type${NC} (product: ${product_key:0:12}..., device: ${device_key:0:12}...)"
    else
        ((events_failed++))
        echo -e "[${RED}✗${NC}] Event failed: HTTP $http_code - $body"
    fi
}

# Main loop
counter=0
while true; do
    # Alternate between telemetry and events
    if [ $((counter % 2)) -eq 0 ]; then
        send_telemetry
    else
        send_event
    fi

    ((counter++))

    # Show progress every 10 requests
    if [ $((counter % 10)) -eq 0 ]; then
        echo ""
        echo -e "${YELLOW}Stats:${NC} Total: $total_requests | Telemetry: ${GREEN}$telemetry_success${NC}/${RED}$telemetry_failed${NC} | Events: ${GREEN}$events_success${NC}/${RED}$events_failed${NC}"
        echo ""
    fi

    sleep "$INTERVAL"
done
