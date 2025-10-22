#!/bin/bash
echo "Testing /chat/completions endpoint with OpenAI virtual key..."
echo ""

curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Say hello!"}
    ],
    "max_tokens": 50
  }' | jq '.'

echo ""
echo "Note: This will fail with 502 if you don't have a real OpenAI API key in keys.json"
