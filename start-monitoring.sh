#!/bin/bash

# Load environment variables
source .env

# Create directory if it doesn't exist
mkdir -p grafana/provisioning/datasources

# Generate datasource configuration with proper variable substitution
cat > grafana/provisioning/datasources/prometheus.yaml << EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://localhost:$PROMETHEUS_PORT
    isDefault: true
    editable: false
EOF

echo "Generated Grafana datasource configuration with proper environment variable substitution." 

# Create prometheus.yaml with proper variable substitution
cat > prometheus.yaml << EOF
global:
  scrape_interval: 15s

scrape_configs:      
  - job_name: 'triggerx-keeper-othentic'
    static_configs:
      - targets: ['localhost:6060']

  - job_name: 'triggerx-keeper-exec'
    params:
      address: ['$OPERATOR_ADDRESS']
    static_configs:
      - targets: ['157.173.218.229:8081']
    metrics_path: /metrics/keeper
EOF

echo "Generated prometheus.yaml with proper environment variable substitution." 

# Start the containers
docker compose up -d

echo "Monitoring started for keeper: $OPERATOR_ADDRESS"
echo "Access Grafana at http://localhost:$GRAFANA_PORT (admin/triggerx)" 