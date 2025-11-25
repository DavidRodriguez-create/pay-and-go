#!/bin/bash

echo "Building account-service image..."
podman build -f Dockerfile.account -t account-service:latest .

# echo "Building card-service image..."
# podman build -f Dockerfile.card -t card-service:latest .

echo "✓ All service images built successfully!"
