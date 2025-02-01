# TriggerX Keeper Setup
This guide will walk you through the process of registering as an keepr (operator) to TriggerX AVS and running the TriggerX software.

## Prerequisites
1. **Registered Eigenlayer Operator Account:** Ensure you have a fully registered Eigenlayer operator account. If you don't have one, follow the steps in the [Eigenlayer User Guide](https://docs.eigenlayer.xyz/restaking-guides/restaking-user-guide) to create and fund your account.

## Software/Hardware Requirement 
* Will be released soon.
* Open Ports:
  * 9003 P2P network
  * 3000 Grafana dashboards
  * 9090 Prometheus 

## Keeper Setup
### ​Config File
Clone this [repo](https://github.com/trigg3rX/triggerx-keeper-setup) and run following commands
```bash
git clone https://github.com/trigg3rX/triggerx-keeper-setup.git
cd triggerx-keeper-setup
```
Edit `triggerx_keeper.yaml` and update the values for your setup.

### Generate a BLS pair (recommended)
The register process requires two sets of private keys: an ecdsa private key and a bls private key.  
We recommend creating a new BLS pair for security reasons.  
If you want to create a new BLS pair, follow the steps [here](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#create-and-list-keys).

### Register with eoracle AVS
Execute the following command to install the TriggerX CLI
```bash
curl -sSfL https://raw.githubusercontent.com/trigg3rX/triggerx-backend/main/scripts/binary/install.sh | sh -s
```
Now, run the following command to register with the TriggerX AVS
```bash
triggerx register --config triggerx_keeper.yaml
```

The output should look like
```
{demo output}
```

### Checking the status of TriggerX operator AVS

The following command will print the status of the operator
```bash
triggerx status --config triggerx_keeper.yaml
```