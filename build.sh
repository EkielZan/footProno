#!/bin/env bash
if [[ $1 == 1 ]];then
force="-a"
fi
cd src
echo "Build Main Binaries"
CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo $force -v -ldflags="-X main.Version=1.0" -o ../bin/footProno .
cd ../healthcheck
echo "Build Healthcheck Binaries"
CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo $force -v -o ../bin/healthcheck
cd ..
