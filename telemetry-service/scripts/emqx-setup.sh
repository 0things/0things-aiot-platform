#!/bin/bash

# EMQX Configuration Script
# This script configures EMQX with:
# 1. PostgreSQL authentication
# 2. Kafka connector

set -e

# Configuration - modify these values as needed
EMQX_HOST="${EMQX_HOST:-http://localhost:18083}"
EMQX_API_KEY="${EMQX_API_KEY:-admin}"
EMQX_API_SECRET="${EMQX_API_SECRET:-public}"

# PostgreSQL Configuration
PG_HOST="${PG_HOST:-sbp-81n5qyn2id6j45cr.supabase.opentrust.net}"
PG_PORT="${PG_PORT:-5432}"
PG_DATABASE="${PG_DATABASE:-postgres}"
PG_USERNAME="${PG_USERNAME:-postgres}"
PG_PASSWORD="${PG_PASSWORD:-Denglu654...}"

# Kafka Configuration
KAFKA_HOST="${KAFKA_HOST:-129.211.224.19:30828}"
KAFKA_CONNECTOR_NAME="${KAFKA_CONNECTOR_NAME:-pve13_kafka_consumer}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to make API calls
emqx_api() {
    local method=$1
    local endpoint=$2
    local data=$3

    if [ -n "$data" ]; then
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -u "${EMQX_API_KEY}:${EMQX_API_SECRET}" \
            -d "$data" \
            "${EMQX_HOST}/api/v5${endpoint}"
    else
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -u "${EMQX_API_KEY}:${EMQX_API_SECRET}" \
            "${EMQX_HOST}/api/v5${endpoint}"
    fi
}

# 1. Configure PostgreSQL Authentication
configure_postgresql_auth() {
    log_info "Configuring PostgreSQL authentication..."

    local auth_config='{
        "mechanism": "password_based",
        "backend": "postgresql",
        "enable": true,
        "server": "'"${PG_HOST}:${PG_PORT}"'",
        "database": "'"${PG_DATABASE}"'",
        "username": "'"${PG_USERNAME}"'",
        "password": "'"${PG_PASSWORD}"'",
        "pool_size": 8,
        "auto_reconnect": true,
        "ssl": {
            "enable": false
        },
        "password_hash_algorithm": {
            "name": "plain",
            "salt_position": "disable"
        },
        "query": "SELECT password_hash, salt FROM mqtt_user WHERE username = ${username} LIMIT 1"
    }'

    # Check if authentication already exists
    local existing=$(emqx_api GET "/authentication")

    # Create or update PostgreSQL authentication
    local response=$(emqx_api POST "/authentication" "$auth_config")

    if echo "$response" | grep -q "error"; then
        # Try to update existing
        log_warn "Authentication may already exist, trying to update..."
        response=$(emqx_api PUT "/authentication/password_based:postgresql" "$auth_config")
    fi

    if echo "$response" | grep -q "error"; then
        log_error "Failed to configure PostgreSQL authentication"
        echo "$response"
        return 1
    else
        log_info "PostgreSQL authentication configured successfully"
    fi
}

# 2. Configure Kafka Connector
configure_kafka_connector() {
    log_info "Configuring Kafka connector..."

    local connector_config='{
        "name": "'"${KAFKA_CONNECTOR_NAME}"'",
        "type": "kafka_producer",
        "enable": true,
        "bootstrap_hosts": "'"${KAFKA_HOST}"'",
        "connect_timeout": "5s",
        "metadata_request_timeout": "5s",
        "min_metadata_refresh_interval": "3s",
        "socket_opts": {
            "nodelay": true,
            "recbuf": "1024KB",
            "sndbuf": "1024KB"
        },
        "authentication": "none",
        "ssl": {
            "enable": false
        },
        "resource_opts": {
            "health_check_interval": "15s",
            "start_after_created": true,
            "start_timeout": "5s"
        }
    }'

    # Create Kafka connector
    local response=$(emqx_api POST "/connectors" "$connector_config")

    if echo "$response" | grep -q "error"; then
        # Try to update existing
        log_warn "Connector may already exist, trying to update..."
        response=$(emqx_api PUT "/connectors/${KAFKA_CONNECTOR_NAME}" "$connector_config")
    fi

    if echo "$response" | grep -q "error"; then
        log_error "Failed to configure Kafka connector"
        echo "$response"
        return 1
    else
        log_info "Kafka connector configured successfully"
    fi
}

# Main
main() {
    log_info "Starting EMQX configuration..."
    log_info "EMQX Host: ${EMQX_HOST}"

    # Test connection
    log_info "Testing EMQX API connection..."
    local status=$(emqx_api GET "/status")
    if echo "$status" | grep -q "error\|failed\|refused"; then
        log_error "Cannot connect to EMQX API at ${EMQX_HOST}"
        exit 1
    fi
    log_info "EMQX API connection successful"

    # Configure PostgreSQL authentication
    configure_postgresql_auth

    # Configure Kafka connector
    configure_kafka_connector

    log_info "EMQX configuration completed!"
}

# Run main function
main "$@"
