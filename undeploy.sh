#!/bin/bash

echo "Removing deployments..."
kubectl delete -f k8s/account-service.yaml --ignore-not-found
kubectl delete -f k8s/card-service.yaml --ignore-not-found

echo "✓ All services removed"
