#!/usr/bin/env bash
set -euo pipefail

if [ -z "${NUM_PHYLAXS}" ]; then
    echo "Error: NUM_PHYLAXS is unset, cannot create deltachain config."
    exit 1
fi

# Get the hostname
hostname=$(hostname)

# # for local development/debugging, set the hostname as it  would be in a devnet container
# if [ ! -z "${INST}" ]; then
#   hostname="deltachain-$INST"
#   echo "set hostname with INST value: $hostname"
# fi

# Check if the hostname starts with "deltachain-"
if [[ ! $hostname =~ ^deltachain- ]]; then
  # If the hostname does not start with "deltachain-", print an error message and exit
  echo "Error: hostname does not start with 'deltachain-'"
  exit 1
fi

# Split the hostname on "-"
instance=$(echo $hostname | cut -d'-' -f2)

# get the context of this script call, so it can be prepended to file paths,
# so this script will work in the tilt docker container, and when run locally.
pwd=$(pwd)

# config dir path for the instance, which is passed to deltachaind via --home
home_path="$pwd/devnet/$hostname"

# copy config from devnet/base
cp -r $pwd/devnet/base/* ${home_path}/

# update the moniker
sed -i "s/moniker = \"deltachain\"/moniker = \"$hostname\"/g" ${home_path}/config/config.toml

# set the external address so deltachain peers can resolve each other
sed -i "s/external_address = \"\"/external_address = \"${hostname}:26656\"/g" ${home_path}/config/config.toml

if [ $instance -eq 0 ] && [ $NUM_PHYLAXS -ge 2 ]; then
  echo "$hostname: enabling seed mode in config.toml."
  sed -i "s/pex = false/pex = true/g" ${home_path}/config/config.toml
  sed -i "s/seed_mode = false/seed_mode = true/g" ${home_path}/config/config.toml
elif [ $instance -ge 1 ]; then
  echo "$hostname: adding seed address to config.toml."
  sed -i "s/seeds = \"\"/seeds = \"90ea40bee73abfda5226a0e8ddb18b0e324d2a29@deltachain-0:26656\"/g" ${home_path}/config/config.toml
  sed -i "s/persistent_peers = \"\"/persistent_peers = \"90ea40bee73abfda5226a0e8ddb18b0e324d2a29@deltachain-0:26656\"/g" ${home_path}/config/config.toml
fi

# copy the config to tendermint's default location, ~/.{chain-id}
mkdir -p /root/.deltachain && cp -r ${home_path}/* /root/.deltachain/

echo "$hostname: done with create-config.sh, exiting success."
