# TriggerX Keeper Setup
This guide will walk you through the process of monitoring your keepr (operator) on TriggerX AVS.

## Prerequisites
- Git
- Docker

## Installation
1. Clone the repository
```bash
git clone https://github.com/trigg3rX/triggerx-keeper-setup.git
```
2. Create .env file
```bash
cp .env.example .env
```
3. Edit .env file
```bash
nano .env
```
4. Start Monitoring using Grafana
```bash
./start_monitoring.sh
```

## Grafana Dashboard

View the dashboard at `http://<host>:3000`, where `<host>` is the IP address of the machine running Keeper and Grafana. If firewall is enabled, you need to open port 3000.

- Username: `admin`
- Password: `triggerx`

## Navigate to TriggerX Keeper Dashboard

On the left sidebar, navigate to Dashboards > TriggerX Keeper Dashboard.

