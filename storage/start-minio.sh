#!/bin/bash

# Script to start MinIO for the storage service

echo "Starting MinIO server..."

# Determine script directory and load .env if present
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$DIR/.env"
if [ -f "$ENV_FILE" ]; then
    # export variables defined in .env
    set -o allexport
    # shellcheck source=/dev/null
    source "$ENV_FILE"
    set +o allexport
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker Desktop first."
    exit 1
fi

# Check if MinIO container already exists
if docker ps -a --format '{{.Names}}' | grep -q '^fate-vault-minio$'; then
    echo "MinIO container exists. Starting it..."
    docker start fate-vault-minio
else
    echo "Creating new MinIO container..."
    docker run -d \
        --name fate-vault-minio \
        -p 9000:9000 \
        -p 9001:9001 \
        -e MINIO_ROOT_USER="${MINIO_USER}" \
        -e MINIO_ROOT_PASSWORD="${MINIO_PASSWORD}" \
        -v minio_data:/data \
        minio/minio server /data --console-address ":9001"
fi

echo ""
echo "MinIO is starting..."
echo "API endpoint: http://localhost:9000"
echo "Console: http://localhost:9001"
echo "Username: ${MINIO_USER}"
echo "Password: ${MINIO_PASSWORD}"
echo ""
echo "Wait a few seconds for MinIO to be ready, then run your storage service."


