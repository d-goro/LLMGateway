#!/bin/bash
echo "Testing without Authorization header (should return 401)..."
echo ""

curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello"}
    ]
  }' | jq '.'

echo ""
echo "Expected: 401 Unauthorized - missing Authorization header"
