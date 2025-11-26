#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display usage
usage() {
    echo "Usage: $0 {start|stop|restart|status} [--no-browser]"
    echo ""
    echo "Commands:"
    echo "  start    - Build and deploy all services (Zookeeper, Kafka, Account, Card)"
    echo "  stop     - Stop and remove all services"
    echo "  restart  - Stop and then start all services"
    echo "  status   - Show status of all running services"
    echo ""
    echo "Options:"
    echo "  --no-browser  - Don't open UI in browser after starting services"
    exit 1
}

# Function to print colored messages
print_header() {
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}=========================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Function to start all services
start_services() {
    print_header "Starting Pay-and-Go Services"

    echo "🧹 Cleaning up existing containers..."
    podman rm -f account-service card-service kafka zookeeper 2>/dev/null || true
    print_success "Cleanup complete"
    echo ""

    echo "🔨 Building service images..."
    echo "Building account-service image..."
    podman build -f podman/Dockerfile.account -t account-service:latest . > /dev/null 2>&1
    echo "Building card-service image..."
    podman build -f podman/Dockerfile.card -t card-service:latest . > /dev/null 2>&1
    print_success "Images built successfully"
    echo ""

    echo "📥 Pulling Kafka images..."
    podman pull confluentinc/cp-zookeeper:latest > /dev/null 2>&1 || true
    podman pull confluentinc/cp-kafka:7.5.0 > /dev/null 2>&1 || true
    print_success "Kafka images ready"
    echo ""

    echo "🌐 Creating network..."
    podman network exists pay-and-go-network || podman network create pay-and-go-network > /dev/null 2>&1
    print_success "Network ready"
    echo ""

    echo "🚀 Starting Zookeeper..."
    podman run -d \
        --name zookeeper \
        --network pay-and-go-network \
        -p 2181:2181 \
        -e ZOOKEEPER_CLIENT_PORT=2181 \
        -e ZOOKEEPER_TICK_TIME=2000 \
        confluentinc/cp-zookeeper:latest > /dev/null 2>&1
    print_success "Zookeeper started"
    echo ""

    echo "🚀 Starting Kafka..."
    podman run -d \
        --name kafka \
        --network pay-and-go-network \
        -p 9092:9092 \
        -p 9093:9093 \
        -e KAFKA_BROKER_ID=1 \
        -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
        -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092,PLAINTEXT_INTERNAL://kafka:9093 \
        -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_INTERNAL:PLAINTEXT \
        -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT_INTERNAL \
        -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
        confluentinc/cp-kafka:7.5.0 > /dev/null 2>&1
    
    echo "⏳ Waiting for Kafka to be ready..."
    sleep 8
    print_success "Kafka started"
    echo ""

    echo "🚀 Starting Account Service..."
    podman run -d \
        --name account-service \
        --network pay-and-go-network \
        -p 8081:8081 \
        -e PORT=8081 \
        -e KAFKA_BROKERS=kafka:9093 \
        -e KAFKA_TOPIC=account-events \
        localhost/account-service:latest > /dev/null 2>&1
    print_success "Account Service started"
    echo ""

    echo "🚀 Starting Card Service..."
    podman run -d \
        --name card-service \
        --network pay-and-go-network \
        -p 8082:8082 \
        -e PORT=8082 \
        -e KAFKA_BROKERS=kafka:9093 \
        -e KAFKA_TOPIC=account-events \
        -e KAFKA_GROUP_ID=card-service \
        localhost/card-service:latest > /dev/null 2>&1
    print_success "Card Service started"
    echo ""

    print_header "Deployment Complete!"
    echo ""
    echo "Services available at:"
    echo "  💼 Account Service: http://localhost:8081"
    echo "  💳 Card Service:    http://localhost:8082"
    echo "  📨 Kafka Broker:    localhost:9092"
    echo "  🔧 Zookeeper:       localhost:2181"
    echo ""
    
    # Open UI in default browser (unless --no-browser flag is set)
    if [[ "$OPEN_BROWSER" == "true" ]]; then
        UI_PATH="$(pwd)/ui.html"
        if [[ -f "$UI_PATH" ]]; then
            echo "🌐 Opening UI in your default browser..."
            if [[ "$OSTYPE" == "linux-gnu"* ]]; then
                xdg-open "$UI_PATH" 2>/dev/null || sensible-browser "$UI_PATH" 2>/dev/null || true
            elif [[ "$OSTYPE" == "darwin"* ]]; then
                open "$UI_PATH"
            elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
                start "$UI_PATH" 2>/dev/null || cmd.exe /c start "" "$UI_PATH" 2>/dev/null || true
            fi
            print_success "UI opened in browser"
            echo ""
        else
            print_error "ui.html not found in current directory"
            echo ""
        fi
    else
        echo "  🌐 UI:              Open ui.html in your browser"
        echo ""
    fi
    
    print_info "Run '$0 status' to check service health"
    print_info "Run '$0 stop' to stop all services"
}

# Function to stop all services
stop_services() {
    print_header "Stopping Pay-and-Go Services"

    echo "🛑 Stopping and removing containers..."
    podman rm -f account-service card-service kafka zookeeper 2>/dev/null || true
    print_success "All services stopped and removed"
    echo ""

    echo "🌐 Removing network..."
    podman network rm pay-and-go-network 2>/dev/null || true
    print_success "Network removed"
    echo ""

    print_header "Services Stopped"
    print_info "Run '$0 start' to start services again"
}

# Function to show service status
show_status() {
    print_header "Service Status"

    if ! podman ps --filter "name=zookeeper" --format "{{.Names}}" | grep -q "zookeeper"; then
        print_error "Services are not running"
        echo ""
        print_info "Run '$0 start' to start all services"
        exit 0
    fi

    echo "📋 Running containers:"
    podman ps --filter "name=zookeeper|kafka|account-service|card-service" \
        --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""

    echo "🏥 Health checks:"
    
    # Check Account Service
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        print_success "Account Service is healthy (http://localhost:8081)"
    else
        print_error "Account Service is not responding"
    fi

    # Check Card Service
    if curl -s http://localhost:8082/health > /dev/null 2>&1; then
        print_success "Card Service is healthy (http://localhost:8082)"
    else
        print_error "Card Service is not responding"
    fi

    echo ""
    print_info "View logs: podman logs -f <service-name>"
    print_info "Open UI: Open ui.html in your browser"
}

# Function to restart services
restart_services() {
    print_header "Restarting Pay-and-Go Services"
    stop_services
    echo ""
    start_services
}

# Parse command-line arguments
COMMAND="${1:-}"
OPEN_BROWSER="true"

# Check for --no-browser flag
for arg in "$@"; do
    if [[ "$arg" == "--no-browser" ]]; then
        OPEN_BROWSER="false"
    fi
done

# Main script logic
case "$COMMAND" in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    status)
        show_status
        ;;
    *)
        usage
        ;;
esac
