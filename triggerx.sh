#!/bin/bash

# TriggerX Management Script

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
        echo "Starting TriggerX core services..."
        docker compose --profile core up -d
        ;;

    start-mon)
        echo "Starting TriggerX with monitoring services..."
        docker compose --profile monitoring up -d
        ;;

    stop)
        echo "Stopping TriggerX services..."
        docker compose --profile core down
        docker compose --profile monitoring down
        docker volume rm othentic_peerstore_temp || true
        ;;

    stop-mon)
        echo "Stopping TriggerX monitoring services..."
        docker compose --profile monitoring down
        ;;

    logs)
        docker compose --profile core logs -f
        ;;

    logs-keeper)
        docker compose logs -f keeper
        ;;

    logs-othentic)
        docker compose logs -f othentic
        ;;

    status)
        echo "TriggerX Service Status:"
        docker compose ps
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
