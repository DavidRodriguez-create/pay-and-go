#!/bin/bash

set -e

echo "========================================="
echo "  Pay-and-Go Services Deployment"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Step 1: Clean up existing containers
echo "🧹 Cleaning up existing containers..."
podman rm -f account-service card-service kafka zookeeper 2>/dev/null || true
echo -e "${GREEN}✓${NC} Cleanup complete"
echo ""

# Step 2: Build images
echo "🔨 Building service images..."
./build-images.sh
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Images built successfully"
else
    echo -e "${RED}✗${NC} Image build failed"
    exit 1
fi
echo ""

# Step 3: Pull Kafka images if not present
echo "📥 Pulling Kafka images (if needed)..."
podman pull confluentinc/cp-zookeeper:latest > /dev/null 2>&1 || true
podman pull confluentinc/cp-kafka:7.5.0 > /dev/null 2>&1 || true
echo -e "${GREEN}✓${NC} Kafka images ready"
echo ""

# Step 4: Create network
echo "🌐 Creating network..."
podman network create pay-and-go-network 2>/dev/null || echo "Network already exists"
echo -e "${GREEN}✓${NC} Network ready"
echo ""

# Step 5: Start Zookeeper
echo "🚀 Starting Zookeeper..."
podman run -d --name zookeeper \
    --network pay-and-go-network \
    -p 2181:2181 \
    -e ZOOKEEPER_CLIENT_PORT=2181 \
    -e ZOOKEEPER_TICK_TIME=2000 \
    confluentinc/cp-zookeeper:latest > /dev/null

echo -e "${YELLOW}⏳${NC} Waiting for Zookeeper to be ready..."
sleep 5
echo -e "${GREEN}✓${NC} Zookeeper started"
echo ""

# Step 6: Start Kafka
echo "🚀 Starting Kafka..."
podman run -d --name kafka \
    --network pay-and-go-network \
    -p 9092:9092 \
    -p 9093:9093 \
    -e KAFKA_BROKER_ID=1 \
    -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
    -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9093,PLAINTEXT_HOST://localhost:9092 \
    -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT \
    -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
    -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9093,PLAINTEXT_HOST://0.0.0.0:9092 \
    -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
    confluentinc/cp-kafka:7.5.0 > /dev/null

echo -e "${YELLOW}⏳${NC} Waiting for Kafka to be ready..."
sleep 10
echo -e "${GREEN}✓${NC} Kafka started"
echo ""

# Step 7: Start Account Service
echo "🚀 Starting Account Service..."
podman run -d --name account-service \
    --network pay-and-go-network \
    -p 8081:8081 \
    -e PORT=8081 \
    -e KAFKA_BROKERS=kafka:9093 \
    -e KAFKA_TOPIC=account-events \
    localhost/account-service:latest > /dev/null

sleep 3
echo -e "${GREEN}✓${NC} Account Service started"
echo ""

# Step 8: Start Card Service
echo "🚀 Starting Card Service..."
podman run -d --name card-service \
    --network pay-and-go-network \
    -p 8082:8082 \
    -e PORT=8082 \
    -e KAFKA_BROKERS=kafka:9093 \
    -e KAFKA_TOPIC=account-events \
    -e KAFKA_GROUP_ID=card-service \
    localhost/card-service:latest > /dev/null

sleep 3
echo -e "${GREEN}✓${NC} Card Service started"
echo ""

# Step 9: Verify deployment
echo "========================================="
echo "  Deployment Verification"
echo "========================================="
echo ""

echo "📋 Running services:"
podman ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "NAMES|zookeeper|kafka|account|card"
echo ""

# Final status
echo "========================================="
echo "  Deployment Complete!"
echo "========================================="
echo ""
echo "Services available at:"
echo "  • Account Service: http://localhost:8081"
echo "  • Card Service:    http://localhost:8082"
echo "  • Kafka Broker:    localhost:9092"
echo "  • Zookeeper:       localhost:2181"
echo ""
echo "========================================="
echo "  Sample Test Cases"
echo "========================================="
echo ""
echo "1️⃣  Test Health Endpoints:"
echo "   curl http://localhost:8081/health"
echo "   curl http://localhost:8082/health"
echo ""
echo "2️⃣  Create an Account (triggers Kafka event):"
echo "   curl -X POST http://localhost:8081/account \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -d '{\"beholder_name\":\"John Doe\",\"country_code\":\"US\"}'"
echo ""
echo "3️⃣  List All Accounts:"
echo "   curl http://localhost:8081/accounts"
echo ""
echo "4️⃣  Get Account by ID:"
echo "   curl 'http://localhost:8081/accounts?id=<ACCOUNT_ID>'"
echo ""
echo "5️⃣  Update Account Status (triggers Kafka event on status change):"
echo "   curl -X PUT 'http://localhost:8081/accounts?id=<ACCOUNT_ID>' \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -d '{\"status\":\"SUSPENDED\"}'"
echo ""
echo "6️⃣  Create a Card (requires account synced via Kafka):"
echo "   curl -X POST http://localhost:8082/card \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -d '{\"account_id\":\"<ACCOUNT_ID>\",\"card_type\":\"DEBIT\",\"country\":\"US\"}'"
echo ""
echo "7️⃣  List All Cards:"
echo "   curl http://localhost:8082/cards"
echo ""
echo "8️⃣  Get Card by ID:"
echo "   curl 'http://localhost:8082/cards?id=<CARD_ID>'"
echo ""
echo "========================================="
echo "  Verification Commands"
echo "========================================="
echo ""
echo "📋 Check running containers:"
echo "   podman ps"
echo ""
echo "📊 View service logs:"
echo "   podman logs -f account-service"
echo "   podman logs -f card-service"
echo "   podman logs -f kafka"
echo ""
echo "🔍 Verify Kafka events:"
echo "   podman logs card-service | grep 'Received account event'"
echo ""
echo "⚙️  Check account service Kafka connection:"
echo "   podman logs account-service | grep -i kafka"
echo ""
echo "🛑 Stop all services:"
echo "   podman rm -f account-service card-service kafka zookeeper"
echo ""
echo "🔄 Restart a service:"
echo "   podman restart <service-name>"
echo ""
