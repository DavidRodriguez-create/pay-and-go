#!/bin/bash

echo "Building service images..."
echo ""

echo "Building account-service image..."
podman build -f podman/Dockerfile.account -t account-service:latest .

echo ""
echo "Building card-service image..."
podman build -f podman/Dockerfile.card -t card-service:latest .

echo ""
echo "✓ All service images built successfully!"
