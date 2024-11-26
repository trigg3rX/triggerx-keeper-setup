#!/bin/bash

echo "Fetching contract ABIs and bytecode..."
go run ./triggerxinterface/fetchabi.go

echo "Generating bindings..."

script_path=$(
    cd "$(dirname "${BASH_SOURCE[0]}")"
    pwd -P
)

if [[ "$(docker images -q abigen-with-interfaces 2> /dev/null)" == "" ]]; then
    docker build -t abigen-with-interfaces -f ./triggerxinterface/abigen-with-interfaces.Dockerfile $script_path
fi

function create_binding {
    contract_dir=$1
    contract=$2
    binding_dir=$3
    echo "creating bindings for $contract..."
    mkdir -p $binding_dir/${contract}
    
    abi_file="$script_path/abis/${contract}.abi"
    bin_file="$script_path/abis/${contract}.bin"

    rm -f $binding_dir/${contract}/binding.go
    docker run --rm \
        --user $(id -u):$(id -g) \
        -v $(realpath $binding_dir):/home/binding_dir \
        -v $(realpath $script_path):/home/triggerxinterface \
        abigen-with-interfaces \
        --bin=/home/triggerxinterface/abis/${contract}.bin \
        --abi=/home/triggerxinterface/abis/${contract}.abi \
        --pkg=contract${contract} \
        --out=/home/binding_dir/${contract}/binding.go
}

cd $script_path

for abi_file in $script_path/abis/*.abi; do
    if [ -f "$abi_file" ]; then
        contract=$(basename "$abi_file" .abi)
        if [ -f "$script_path/abis/$contract.bin" ]; then
            create_binding . "$contract" ./bindings
        fi
    fi
done
