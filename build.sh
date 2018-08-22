#!/bin/bash

echo "Building for Docker..."
echo
echo "GOOS set to Linux..."
export GOOS=linux
echo
echo "Building Infra Client..."
go build -o iclient client/client.go client/handlers.go
echo "Building Infra Server..."
go build -o iserver server/*.go
echo
echo "Building client container..."
docker build -f iclient.Dockerfile -t infra_client .
echo
echo "Building server container..."
docker build -f iserver.Dockerfile -t infra_server .
echo
echo "Finished"