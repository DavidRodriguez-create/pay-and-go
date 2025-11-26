#!/bin/bash

echo "Deploying all services to Kubernetes..."
echo ""

# Deploy all services (Kafka, Zookeeper, Account Service, Card Service)
kubectl apply -f all-services.yaml

echo ""
echo "✓ Deployment complete!"
echo ""
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=zookeeper --timeout=60s
kubectl wait --for=condition=ready pod -l app=kafka --timeout=60s
kubectl wait --for=condition=ready pod -l app=account-service --timeout=60s
kubectl wait --for=condition=ready pod -l app=card-service --timeout=60s

echo ""
echo "Check status:"
echo "  kubectl get pods"
echo "  kubectl get services"
echo ""
echo "Access services:"
echo "  Account Service: http://localhost:8081"
echo "  Card Service: http://localhost:8082"
echo ""
echo "View logs:"
echo "  kubectl logs -f deployment/account-service"
echo "  kubectl logs -f deployment/card-service"
echo "  kubectl logs -f deployment/kafka"
