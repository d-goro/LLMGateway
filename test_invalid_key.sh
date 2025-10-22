#!/bin/bash
echo "Testing with invalid virtual key (should return 401)..."
echo ""

curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid_key_123" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello"}
    ]
  }' | jq '.'

echo ""
echo "Expected: 401 Unauthorized with error message"
