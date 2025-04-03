#!/bin/bash

set -e

CONTAINER_NAME="penilaian-360"

# Build the Docker image
echo "Building Docker image..."
docker build --network host -t $CONTAINER_NAME .

# Stop and remove existing container if running
EXISTING_CONTAINER=$(docker ps -aq --filter "name=$CONTAINER_NAME")
if [ -n "$EXISTING_CONTAINER" ]; then
    echo "Stopping existing container..."
    docker stop $EXISTING_CONTAINER
    echo "Removing existing container..."
    docker rm $EXISTING_CONTAINER
fi

# Run the new container
echo "Starting new container..."
docker run -d --network=host \
    -v "$PWD/params/.env:/app/params/.env:Z" \
    --name $CONTAINER_NAME \
    $CONTAINER_NAME

echo "Deployment completed successfully!"
