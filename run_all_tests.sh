#!/bin/bash

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Usage information
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Runs all end-to-end tests for the LLM Gateway."
    echo ""
    echo "Options:"
    echo "  --cleanup    Automatically stop the gateway after tests (if started by this script)"
    echo "  -h, --help   Show this help message"
    echo ""
    echo "The script will automatically start the gateway with Docker if it's not running."
    echo ""
    exit 0
fi

# Parse command line arguments
AUTO_CLEANUP=false
if [[ "$1" == "--cleanup" ]]; then
    AUTO_CLEANUP=true
fi

GATEWAY_STARTED=false

# Cleanup function
cleanup() {
    if [ "$GATEWAY_STARTED" = true ] && [ "$AUTO_CLEANUP" = true ]; then
        echo ""
        echo -e "${YELLOW}Cleaning up: stopping gateway...${NC}"
        make docker-stop > /dev/null 2>&1 || true
        echo -e "${GREEN}✓ Gateway stopped${NC}"
    fi
}

# Trap cleanup on script exit (including errors)
trap cleanup EXIT

echo "=========================================="
echo "LLM Gateway - Full Test Suite"
echo "=========================================="
echo ""

# Check if gateway is already running
echo "Checking if gateway is running..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Gateway is already running${NC}"
    GATEWAY_STARTED=false
else
    echo -e "${YELLOW}Gateway not running. Starting with Docker...${NC}"

    # Check if container exists but is stopped
    if docker ps -a --format '{{.Names}}' | grep -q "^llm-gateway$"; then
        echo "Removing existing container..."
        docker rm -f llm-gateway > /dev/null 2>&1 || true
    fi

    # Start the gateway with Docker
    echo "Running: make docker-run"
    make docker-run

    GATEWAY_STARTED=true

    # Wait for gateway to be ready
    echo -n "Waiting for gateway to be ready"
    MAX_WAIT=30
    WAITED=0
    while ! curl -s http://localhost:8080/health > /dev/null 2>&1; do
        if [ $WAITED -ge $MAX_WAIT ]; then
            echo -e "\n${RED}✗ Gateway failed to start within ${MAX_WAIT} seconds${NC}"
            echo "Check docker logs: docker logs llm-gateway"
            exit 1
        fi
        echo -n "."
        sleep 1
        WAITED=$((WAITED + 1))
    done
    echo -e "\n${GREEN}✓ Gateway is ready!${NC}"
fi

echo ""
echo "=========================================="
echo "Test 1: Health Check"
echo "=========================================="
./test_health.sh

echo ""
echo ""
echo "=========================================="
echo "Test 2: Metrics Check"
echo "=========================================="
./test_metrics.sh

echo ""
echo ""
echo "=========================================="
echo "Test 3: Invalid Virtual Key (401)"
echo "=========================================="
./test_invalid_key.sh

echo ""
echo ""
echo "=========================================="
echo "Test 4: Missing Authorization (401)"
echo "=========================================="
./test_no_auth.sh

echo ""
echo ""
echo "=========================================="
echo "Test 5: Chat Completion Request"
echo "=========================================="
echo "Note: This will fail with 502 Bad Gateway if you don't have real API keys"
./test_chat.sh

echo ""
echo ""
echo "=========================================="
echo "Test 6: Check Metrics Again"
echo "=========================================="
./test_metrics.sh

echo ""
echo ""
echo "=========================================="
echo "All Tests Complete!"
echo "=========================================="
echo ""

# Cleanup if we started the gateway (and not using auto-cleanup)
if [ "$GATEWAY_STARTED" = true ] && [ "$AUTO_CLEANUP" = false ]; then
    echo ""
    echo "The gateway was started by this script."
    echo -n "Do you want to stop the gateway container? (y/N): "
    read -r RESPONSE
    if [[ "$RESPONSE" =~ ^[Yy]$ ]]; then
        echo "Stopping gateway..."
        make docker-stop
        echo -e "${GREEN}✓ Gateway stopped${NC}"
    else
        echo -e "${YELLOW}Gateway is still running. Stop it with: make docker-stop${NC}"
    fi
elif [ "$GATEWAY_STARTED" = true ] && [ "$AUTO_CLEANUP" = true ]; then
    echo -e "${GREEN}Gateway will be stopped automatically (--cleanup flag)${NC}"
fi

echo ""
echo "To test with real API keys:"
echo "  1. Edit keys.json with your real OpenAI/Anthropic keys"
echo "  2. Restart the gateway: make docker-run"
echo "  3. Run ./test_chat.sh again"
echo ""
echo "For more testing info, see: TESTING_GUIDE.md"
echo ""
