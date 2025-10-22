# Testing Guide

Complete guide to testing the LLM Gateway.

## Quick Start

### Run All Tests (Autonomous)

The easiest way to test is to run the autonomous test script:

```bash
./run_all_tests.sh
```

This script will:
- ✓ Check if the gateway is running
- ✓ Automatically start it with Docker if needed
- ✓ Wait for the gateway to be ready
- ✓ Run all test scenarios
- ✓ Optionally stop the gateway after tests

**Options:**
```bash
./run_all_tests.sh            # Interactive: prompts to stop if started
./run_all_tests.sh --cleanup  # Automatic: stops gateway after tests
./run_all_tests.sh --help     # Show usage information
```

### Manual Testing

If you prefer to manage the gateway manually:

**Option A: Docker (Recommended)**
```bash
# Terminal 1
make docker-run
```

**Option B: Build from Source**
```bash
# Terminal 1
make run
```

Then run tests in another terminal:
```bash
# Terminal 2
./run_all_tests.sh
```

## Test Scripts

We've created several test scripts for you:

| Script | Description | Expected Result |
|--------|-------------|-----------------|
| `test_health.sh` | Health check endpoint | Shows service status and provider availability |
| `test_metrics.sh` | Usage metrics | Shows request counts and timing |
| `test_chat.sh` | Chat completion request | 502 error without real API keys |
| `test_invalid_key.sh` | Invalid virtual key | 401 Unauthorized |
| `test_no_auth.sh` | Missing authorization | 401 Unauthorized |
| `run_all_tests.sh` | All tests together | Complete test suite |

## Manual Testing

### Using cURL

```bash
# Health check (no auth required)
curl http://localhost:8080/health | jq '.'

# Metrics (no auth required)
curl http://localhost:8080/metrics | jq '.'

# Chat completion (requires auth)
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Using Python

```bash
pip install openai
python examples/python_client.py
```

### Using Go

```bash
go run examples/go_client.go
```

### Using the Full cURL Examples

```bash
./examples/curl_examples.sh
```

## Testing with Real API Keys

### 1. Get API Keys

- **OpenAI**: https://platform.openai.com/api-keys
- **Anthropic**: https://console.anthropic.com/settings/keys

### 2. Update Configuration

Edit `keys.json`:

```json
{
  "virtual_keys": {
    "vk_user1_openai": {
      "provider": "openai",
      "api_key": "sk-proj-YOUR_REAL_KEY_HERE"
    },
    "vk_user2_anthropic": {
      "provider": "anthropic",
      "api_key": "sk-ant-YOUR_REAL_KEY_HERE"
    }
  }
}
```

### 3. Restart Gateway

**Option A: Docker**
```bash
# Restart to reload keys
docker restart llm-gateway

# Or fully rebuild
make docker-stop
make docker-run
```

**Option B: From Source**
```bash
# Stop the running gateway (Ctrl+C) and restart:
make run
```

### 4. Test Real Requests

```bash
# OpenAI
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Say hello in Spanish!"}
    ],
    "max_tokens": 50
  }' | jq '.'

# Anthropic
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user2_anthropic" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "messages": [
      {"role": "user", "content": "Say hello in French!"}
    ],
    "max_tokens": 50
  }' | jq '.'
```

## What to Expect

### Without Real API Keys

- ✅ `/health` - Returns 200, shows providers as unhealthy
- ✅ `/metrics` - Returns 200, shows 0 requests
- ✅ Invalid auth - Returns 401 with error
- ❌ `/chat/completions` - Returns 502 (cannot reach provider)

### With Real API Keys

- ✅ `/health` - Returns 200, shows providers as healthy
- ✅ `/chat/completions` - Returns 200 with LLM response
- ✅ `/metrics` - Shows increasing request counts
- ✅ Logs appear in terminal with full request/response

## Viewing Logs

All requests are logged to stdout:

```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "virtual_key": "vk_user1_openai",
  "provider": "openai",
  "method": "POST",
  "status": 200,
  "duration_ms": 1250,
  "request": {
    "model": "gpt-3.5-turbo",
    "messages": [...]
  },
  "response": {
    "choices": [...]
  }
}
```

To save logs to a file:

**Docker:**
```bash
# View logs
docker logs llm-gateway

# Follow logs in real-time
docker logs -f llm-gateway

# Save to file
docker logs llm-gateway > gateway.log
```

**From Source:**
```bash
# Redirect output to file
./gateway > gateway.log 2>&1

# Or configure file logging:
export LOG_TO_FILE=true
export LOG_FILE_PATH=./logs/gateway.log
./gateway
```

## Testing Rate Limiting

The gateway has a default quota of 100 requests/hour per virtual key.

```bash
# Test quota (need a loop)
for i in {1..101}; do
  echo "Request $i"
  curl -X POST http://localhost:8080/chat/completions \
    -H "Authorization: Bearer vk_user1_openai" \
    -H "Content-Type: application/json" \
    -d '{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"test"}]}' \
    -s -o /dev/null -w "%{http_code}\n"
done
```

After 100 requests, you should get `429 Too Many Requests`.

## Environment Variables

Test different configurations:

**From Source:**
```bash
# Different port
export SERVER_PORT=9090
make run

# Disable quota
export QUOTA_ENABLED=false
make run

# Custom quota limit
export QUOTA_LIMIT=10
make run

# Different keys file
export KEYS_FILE_PATH=/path/to/other/keys.json
make run
```

**Docker (with custom env vars):**
```bash
# Different port
docker run -d --name llm-gateway -p 9090:9090 -e SERVER_PORT=9090 \
  -v $(pwd)/keys.json:/app/keys.json llm-gateway:latest

# Disable quota
docker run -d --name llm-gateway -p 8080:8080 -e QUOTA_ENABLED=false \
  -v $(pwd)/keys.json:/app/keys.json llm-gateway:latest

# Custom quota limit
docker run -d --name llm-gateway -p 8080:8080 -e QUOTA_LIMIT=10 \
  -v $(pwd)/keys.json:/app/keys.json llm-gateway:latest

# Different keys file
docker run -d --name llm-gateway -p 8080:8080 \
  -v /path/to/other/keys.json:/app/keys.json llm-gateway:latest
```

## Troubleshooting

### Gateway won't start

```bash
# Check if port 8080 is in use
lsof -i :8080

# Use different port (from source)
export SERVER_PORT=8081
make run

# Or Docker
docker run -d --name llm-gateway -p 8081:8081 -e SERVER_PORT=8081 \
  -v $(pwd)/keys.json:/app/keys.json llm-gateway:latest
```

### Can't connect to gateway

```bash
# Check if running
curl http://localhost:8080/health

# Check logs (Docker)
docker logs llm-gateway

# Or check logs (from source)
# Look at terminal where ./gateway is running
```

### 401 Unauthorized

- Check virtual key exists in `keys.json`
- Check Authorization header: `Bearer vk_user1_openai`

### 502 Bad Gateway

- This is expected without real API keys
- Provider APIs are unreachable or keys are invalid
- Check gateway logs for details

## Next Steps

Once basic testing works:

1. ✅ Test with real API keys
2. ✅ Monitor logs and metrics
3. ✅ Test rate limiting
4. ✅ Test with Python/Go clients
5. ✅ Try the Docker deployment

For production deployment, see [README.md](README.md).
