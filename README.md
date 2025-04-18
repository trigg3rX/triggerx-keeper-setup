# TriggerX Keeper Setup
This guide will walk you through the process of setting up and monitoring your keeper (operator) on TriggerX AVS.

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
npm i -g @othentic/othentic-cli@1.10.0
```
5. Initialize Othentic-cli
```bash
othentic-cli init
```
6. Register your Keeper
```bash
othentic-cli operator register
```
7. Start Keeper
```bash
./triggerx.sh start
```
8. Start Monitoring using Grafana
```bash
./triggerx.sh start-mon
```
9. Stop Keeper (and monitoring)
```bash
./triggerx.sh stop
```
10. View logs
```bash
./triggerx.sh logs
```
11. View logs of Services
```bash
./triggerx.sh logs-keeper # Keeper logs
./triggerx.sh logs-othentic # Othentic logs
```
12. Status of Services
```bash
./triggerx.sh status
```
13. Help
```bash
./triggerx.sh help
```

### Grafana Dashboard

View the dashboard at `http://<host>:3000`, where `<host>` is the IP address of the machine running Keeper and Grafana. If firewall is enabled, you need to open port 3000.

- Username: `admin`
- Password: `triggerx`

#### Navigate to TriggerX Keeper Dashboard

On the left sidebar, navigate to Dashboards > TriggerX Keeper Dashboard.

