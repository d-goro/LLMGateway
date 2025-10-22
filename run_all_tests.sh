#!/bin/bash

echo "=========================================="
echo "LLM Gateway - Full Test Suite"
echo "=========================================="
echo ""
echo "Make sure the gateway is running in docker or in another terminal:"
echo ""
echo "Press Enter to continue..."
read

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
echo "To test with real API keys:"
echo "  1. Edit keys.json with your real OpenAI/Anthropic keys"
echo "  2. Restart the gateway"
echo "  3. Run ./test_chat.sh again"
echo ""
