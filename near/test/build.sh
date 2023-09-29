#!/usr/bin/env bash

cd ../..
DOCKER_BUILDKIT=1 docker build --target go-export -f Dockerfile.proto -o type=local,dest=node .
DOCKER_BUILDKIT=1 docker build --target node-export -f Dockerfile.proto -o type=local,dest=. .
cd node/
echo "Have patience, this step takes upwards of 500 seconds!"
if [ $(uname -m) = "arm64" ]; then
    echo "Building Phylax for linux/amd64"
    DOCKER_BUILDKIT=1 docker build --platform linux/amd64 -f Dockerfile -t phylax .
else 
    echo "Building Phylax natively"
    DOCKER_BUILDKIT=1 docker build -f Dockerfile -t phylax .
fi
