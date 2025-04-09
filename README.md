# TriggerX Keeper Setup
This guide will walk you through the process of setting up and monitoring your keepr (operator) on TriggerX AVS.

## Prerequisites
- Git
- Docker
- Nodes.js (strictly v22.6.0) for Othentic-cli

## Setup Keeper

### Installation
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
4. Install Othentic-cli
```bash
npm i -g @othentic/othentic-cli@1.10.1
```
5. Initialize Othentic-cli
```bash
othentic-cli init
```
6. Register your Keeper
```bash
othentic-cli operator register
```
7. Pull the latest Keeper image
```bash
docker pull trigg3rx/triggerx-keeper:latest
```
8. Start Keeper
```bash
docker run -d -v ./.env:/root/.env --name triggerx-keeper --network host trigg3rx/triggerx-keeper:latest
```
9. Start Monitoring using Grafana
```bash
./start_monitoring.sh
```

### Grafana Dashboard

View the dashboard at `http://<host>:3000`, where `<host>` is the IP address of the machine running Keeper and Grafana. If firewall is enabled, you need to open port 3000.

- Username: `admin`
- Password: `triggerx`

#### Navigate to TriggerX Keeper Dashboard

On the left sidebar, navigate to Dashboards > TriggerX Keeper Dashboard.

