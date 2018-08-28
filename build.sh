#!/bin/bash

echo "Building client container..."
docker build -f iclient.Dockerfile -t infra_client .
echo
echo "Building server container..."
docker build -f iserver.Dockerfile -t infra_server .
echo
echo "Finished"