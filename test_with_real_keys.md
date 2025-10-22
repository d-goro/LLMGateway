# Testing with Real API Keys

If you have real OpenAI or Anthropic API keys, follow these steps:

## 1. Update keys.json

```json
{
  "virtual_keys": {
    "vk_user1_openai": {
      "provider": "openai",
      "api_key": "sk-proj-YOUR-REAL-OPENAI-KEY"
    },
    "vk_user2_anthropic": {
      "provider": "anthropic",
      "api_key": "sk-ant-YOUR-REAL-ANTHROPIC-KEY"
    }
  }
}
```

## 2. Start/Restart the Gateway

**Option A: Docker (Recommended)**
```bash
# First time build and run:
make docker-run

# If already running, restart to reload keys:
docker restart llm-gateway

# Or fully rebuild:
make docker-stop
make docker-run
```

**Option B: Build from Source**
```bash
# Build and run:
make run

# If already running, stop (Ctrl+C) and restart:
make run
```

## 3. Test OpenAI

```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Say hello in Spanish!"}
    ],
    "max_tokens": 50
  }'
```

## 4. Test Anthropic

```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user2_anthropic" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "messages": [
      {"role": "user", "content": "Say hello in French!"}
    ],
    "max_tokens": 50
  }'
```

## 5. Check Logs

The gateway logs all interactions to stdout:

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

## 6. Check Metrics

After making requests:

```bash
curl http://localhost:8080/metrics | jq '.'
```

You should see:
```json
{
  "total_requests": 2,
  "requests_by_provider": {
    "openai": 1,
    "anthropic": 1
  },
  "average_response_ms": 1150.5,
  "last_updated": "2025-01-15T10:35:00Z"
}
```
