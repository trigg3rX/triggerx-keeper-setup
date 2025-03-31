#!/bin/bash

source .env

docker-compose up -d

echo "Monitoring started for keeper: $OPERATOR_ADDRESS"
echo "Access Grafana at http://localhost:$GRAFANA_PORT (admin/triggerx)" 