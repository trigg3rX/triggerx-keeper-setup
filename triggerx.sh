#!/bin/bash

# TriggerX Management Script

# Check which docker compose command is available
check_docker_compose() {
  if command -v docker &>/dev/null; then
    if docker compose version &>/dev/null; then
      echo "docker compose"
      return 0
    elif command -v docker-compose &>/dev/null; then
      echo "docker-compose"
      return 0
    else
      echo "Neither 'docker compose' nor 'docker-compose' found. Please install Docker Compose."
      return 1
    fi
  else
    echo "Docker is not installed. Please install Docker and Docker Compose."
    return 1
  fi
}

# Get the appropriate docker compose command
DOCKER_COMPOSE_CMD=$(check_docker_compose)
if [ $? -ne 0 ]; then
  echo "$DOCKER_COMPOSE_CMD"
  exit 1
fi

# Get Docker group ID for container access
get_docker_group_id() {
  if command -v getent &> /dev/null; then
    DOCKER_GROUP_ID=$(getent group docker | cut -d: -f3)
    if [ -n "$DOCKER_GROUP_ID" ]; then
      echo "Detected Docker group ID: $DOCKER_GROUP_ID"
      export DOCKER_GROUP_ID
    else
      echo "Warning: Could not detect Docker group ID, using default 999"
      export DOCKER_GROUP_ID=999
    fi
  else
    echo "Warning: getent not available, using default Docker group ID 999"
    export DOCKER_GROUP_ID=999
  fi
}

# Load environment variables from .env (if present)
load_env_file() {
  if [ -f ".env" ]; then
    # shellcheck disable=SC2046,SC1090
    set -a
    source ".env"
    set +a
  else
    echo "Warning: .env file not found in the current directory."
  fi
}

# Ensure storage-related environment variables are exported
ensure_storage_variables() {
  load_env_file

  if [ -n "${TRIGGERX_CACHE_DIR:-}" ]; then
    return
  fi

  local address="${OPERATOR_ADDRESS:-default}"
  if [ "$address" = "default" ]; then
    echo "Warning: OPERATOR_ADDRESS not set. Using shared storage directory."
  fi

  local sanitized_address
  sanitized_address=$(echo "$address" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]/-/g')
  if [ -z "$sanitized_address" ]; then
    sanitized_address="default"
  fi

  export TRIGGERX_STORAGE_ROOT="$HOME/.triggerx/$sanitized_address"
  export TRIGGERX_CACHE_DIR="$TRIGGERX_STORAGE_ROOT/cache"
  export TRIGGERX_KEEPER_LOG_DIR="$TRIGGERX_STORAGE_ROOT/logs/keeper"
  export TRIGGERX_OTHENTIC_LOG_DIR="$TRIGGERX_STORAGE_ROOT/logs/othentic"
  export TRIGGERX_PEERSTORE_DIR="$TRIGGERX_STORAGE_ROOT/peerstore/othentic"

  echo "Using keeper storage root: $TRIGGERX_STORAGE_ROOT"
}

# Setup persistent storage directories with proper permissions
setup_persistent_storage() {
  ensure_storage_variables

  echo "Setting up persistent storage directories..."

  dirs=(
    "$TRIGGERX_CACHE_DIR"
    "$TRIGGERX_KEEPER_LOG_DIR"
    "$TRIGGERX_OTHENTIC_LOG_DIR"
    "$TRIGGERX_PEERSTORE_DIR"
  )

  for dir in "${dirs[@]}"; do
    if [ ! -d "$dir" ]; then
      mkdir -p "$dir"
      echo "Created directory: $dir"
      chmod -R 777 "$dir"
    fi
  done

  echo "Persistent storage setup complete."
}

# Display help information
show_help() {
    echo "TriggerX Management Script"
    echo "Usage: ./triggerx.sh COMMAND [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start         Start the core services (keeper and othentic)"
    echo "  start-mon     Start the core services with monitoring (includes Prometheus and Grafana)"
    echo "  stop          Stop all services"
    echo "  stop-mon      Stop monitoring services (Prometheus and Grafana)"
    echo "  logs          View logs from all services"
    echo "  logs-keeper   View logs from the keeper service only"
    echo "  logs-othentic View logs from the othentic service only"
    echo "  status        Show the status of all services"
    echo "  help          Show this help message"
    echo ""
}

# Handle command line arguments
case "$1" in
    start)
        ensure_storage_variables
        setup_persistent_storage
        get_docker_group_id
        
        echo "Starting TriggerX core services..."
        $DOCKER_COMPOSE_CMD --profile core up -d
        ;;

    start-mon)
        ensure_storage_variables
        setup_persistent_storage
        get_docker_group_id

cat > prometheus.yaml << EOF
global:
  scrape_interval: 15s

scrape_configs:      
  - job_name: 'triggerx-keeper-attester'
    static_configs:
      - targets: ['othentic:${OPERATOR_METRICS_PORT}']

  - job_name: 'triggerx-keeper-server'
    static_configs:
      - targets: ['keeper:${OPERATOR_RPC_PORT}']
    metrics_path: /metrics
EOF
        echo "Starting TriggerX with monitoring services..."
        $DOCKER_COMPOSE_CMD --profile core --profile monitoring up -d
        ;;

    stop)
        ensure_storage_variables
        echo "Stopping TriggerX services gracefully..."
        $DOCKER_COMPOSE_CMD --profile core stop stop -t 40
        sleep 5
        $DOCKER_COMPOSE_CMD --profile core down --remove-orphans
        $DOCKER_COMPOSE_CMD --profile monitoring down --remove-orphans
        ;;

    stop-mon)
        ensure_storage_variables
        echo "Stopping TriggerX monitoring services..."
        $DOCKER_COMPOSE_CMD --profile monitoring down --remove-orphans
        ;;

    logs)
        ensure_storage_variables
        $DOCKER_COMPOSE_CMD --profile core logs -f
        ;;

    logs-keeper)
        ensure_storage_variables
        $DOCKER_COMPOSE_CMD logs -f keeper
        ;;

    logs-othentic)
        ensure_storage_variables
        $DOCKER_COMPOSE_CMD logs -f othentic
        ;;

    status)
        ensure_storage_variables
        echo "TriggerX Service Status:"
        $DOCKER_COMPOSE_CMD ps
        ;;

    help|--help|-h)
        show_help
        ;;

    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac

exit 0 