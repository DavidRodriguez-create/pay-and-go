#!/bin/bash

echo "Deploying services to Kubernetes..."

# Deploy account service
kubectl apply -f k8s/account-service.yaml

# Uncomment when ready to deploy card service
# kubectl apply -f k8s/card-service.yaml

echo ""
echo "✓ Deployment complete!"
echo ""
echo "Check status:"
echo "  kubectl get pods"
echo "  kubectl get services"
echo ""
echo "Access services:"
echo "  Account Service: http://localhost:8081"
echo "  Card Service: http://localhost:8082"
