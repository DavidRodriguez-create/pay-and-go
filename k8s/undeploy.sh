#!/bin/bash

echo "Removing deployments..."
kubectl delete -f all-services.yaml --ignore-not-found

echo "✓ All services removed"
