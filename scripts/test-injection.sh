#!/usr/bin/env bash
# This script submits a phylax set update using the VAA injection admin command.
# First argument is node to submit to. Second argument is current set index.
set -e

node=$1
idx=$2
localPath=./scripts/new-phylaxset.prototxt
containerPath=/tmp/new-phylaxset.prototxt
sock=/tmp/admin.sock

# Create a phylax set update VAA, pipe stdout to a local file.
kubectl exec -n wormhole phylax-${node} -c phylaxd -- /phylaxd template phylax-set-update --num=1 --idx=${idx} > ${localPath}

# Copy the local VAA prototext to a pod's file system.
kubectl cp ${localPath} wormhole/phylax-${node}:${containerPath} -c phylaxd

# Verify the contents of the VAA prototext file and print the result. The digest incorporates the current time and is NOT deterministic.
kubectl exec -n wormhole phylax-${node} -c phylaxd -- /phylaxd admin governance-vaa-verify $containerPath

# Submit to node
kubectl exec -n wormhole phylax-${node} -c phylaxd -- /phylaxd admin governance-vaa-inject --socket $sock $containerPath
