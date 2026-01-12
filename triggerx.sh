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

# Setup persistent storage directories with proper permissions
setup_persistent_storage() {
  echo "Setting up persistent storage directories..."

  # Define directories for persistent storage
  dirs=(
    "$HOME/.triggerx/cache"
    "$HOME/.triggerx/logs/keeper"
    "$HOME/.triggerx/logs/performer"
    "$HOME/.triggerx/logs/othentic"
  )

  # Check if each directory exists before creating
  for dir in "${dirs[@]}"; do
    if [ ! -d "$dir" ]; then
      mkdir -p "$dir"
      echo "Created directory: $dir"

      # Set permissions to be writable by container user (UID 1000)
      chmod -R 777 "$dir"
    fi
  done

  echo "Persistent storage setup complete in user home (~/.triggerx/)"
}

# Clear peerstore directory before starting
clear_peerstore() {
  echo "Clearing peerstore directory..."
  peerstore_dir="$HOME/.triggerx/peerstore/othentic"
  
  if [ -d "$peerstore_dir" ]; then
    rm -rf "$peerstore_dir"/*
    echo "Cleared peerstore directory: $peerstore_dir"
  else
    echo "Peerstore directory does not exist, skipping clear: $peerstore_dir"
  fi
}

# Pull Docker images required for executor
pull_executor_images() {
  echo "Pulling Docker images for executor..."
  
  # Docker images from docker-executor.yaml
  images=(
    "golang:1.24-alpine"
    "node:22-alpine"
  )
  
  for image in "${images[@]}"; do
    # Check if image already exists locally
    if docker image inspect "$image" &>/dev/null; then
      echo "Image $image already exists locally, skipping pull."
    else
      echo "Pulling $image..."
      docker pull "$image" || {
        echo "Warning: Failed to pull $image. It will be pulled when needed."
      }
    fi
  done
  
  echo "Docker images check/pull complete"
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
        # Clear peerstore directory before starting
        clear_peerstore
        
        # echo "Setting up persistent storage..."
        setup_persistent_storage
        
        # echo "Detecting Docker group ID..."
        get_docker_group_id
        
        # Pull executor Docker images
        pull_executor_images
        
        echo "Starting TriggerX core services..."
        $DOCKER_COMPOSE_CMD --profile core up -d
        ;;

    start-mon)
        # echo "Setting up persistent storage..."
        setup_persistent_storage
        
        # echo "Detecting Docker group ID..."
        get_docker_group_id
        
        # Pull executor Docker images
        pull_executor_images
        
        source .env
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
        echo "Stopping TriggerX services gracefully..."
        $DOCKER_COMPOSE_CMD --profile core stop stop -t 40
        sleep 5
        $DOCKER_COMPOSE_CMD --profile core down --remove-orphans
        $DOCKER_COMPOSE_CMD --profile monitoring down --remove-orphans
        ;;

    stop-mon)
        echo "Stopping TriggerX monitoring services..."
        $DOCKER_COMPOSE_CMD --profile monitoring down --remove-orphans
        ;;

    logs)
        $DOCKER_COMPOSE_CMD --profile core logs -f
        ;;

    logs-keeper)
        $DOCKER_COMPOSE_CMD logs -f keeper
        ;;

    logs-othentic)
        $DOCKER_COMPOSE_CMD logs -f othentic
        ;;

    status)
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
