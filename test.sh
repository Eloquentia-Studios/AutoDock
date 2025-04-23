#!/bin/bash

# Stop and remove the container
echo "Stopping and removing container..."
docker kill autodock-test > /dev/null
docker rm autodock-test > /dev/null

# Pull an old Nginx image
echo "Pulling old Nginx image..."
docker pull nginx:1.26 > /dev/null

# Re-tag the image to act as a new version
echo "Re-tagging image..."
docker tag nginx:1.26 nginx:latest > /dev/null

# Start a container with autodock enabled
echo "Starting container..."
docker run -d \
  --name autodock-test \
  --label autodock.enable=true \
  --label autodock.action=upgrade \
  nginx:latest \
  > /dev/null

# Print delimiter
echo

# Run AutoDock
go run cmd/main/service.go