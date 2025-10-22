#!/bin/bash

# Example curl commands for testing the LLM Gateway

GATEWAY_URL="http://localhost:8080"

echo "LLM Gateway - cURL Examples"
echo "============================"
echo

# 1. Test chat completion with OpenAI
echo "1. Testing chat completion with OpenAI virtual key..."
curl -X POST "${GATEWAY_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Hello! Tell me a fun fact about programming."}
    ],
    "max_tokens": 100
  }'
echo
echo

# 2. Test chat completion with Anthropic
echo "2. Testing chat completion with Anthropic virtual key..."
curl -X POST "${GATEWAY_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user2_anthropic" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "messages": [
      {"role": "user", "content": "Hello! Tell me a fun fact about AI."}
    ],
    "max_tokens": 100
  }'
echo
echo

# 3. Test with invalid virtual key (should return 401)
echo "3. Testing with invalid virtual key (should fail)..."
curl -X POST "${GATEWAY_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid_key" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "This should fail"}
    ]
  }'
echo
echo

# 4. Test health endpoint
echo "4. Testing health endpoint..."
curl "${GATEWAY_URL}/health" | jq '.'
echo
echo

# 5. Test metrics endpoint
echo "5. Testing metrics endpoint..."
curl "${GATEWAY_URL}/metrics" | jq '.'
echo
echo

# 6. Test with missing authorization (should return 401)
echo "6. Testing with missing authorization header (should fail)..."
curl -X POST "${GATEWAY_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "This should fail"}
    ]
  }'
echo
echo

# 7. Test with invalid request format (should return 400)
echo "7. Testing with invalid request format (should fail)..."
curl -X POST "${GATEWAY_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "messages": [
      {"role": "user", "content": "Missing model field"}
    ]
  }'
echo
echo

echo "All tests completed!"
