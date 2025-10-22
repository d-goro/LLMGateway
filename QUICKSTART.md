# Quick Start Guide

Get the LLM Gateway up and running in 5 minutes!

## Prerequisites

Choose one:
- **Docker** (recommended - no Go installation needed), OR
- **Go 1.20+** (to build from source)

## Step 1: Setup Keys

Create your `keys.json` configuration file:

```bash
make setup
# Or manually:
# cp keys.example.json keys.json
```

Edit `keys.json` with your actual API keys:

```json
{
  "virtual_keys": {
    "vk_user1_openai": {
      "provider": "openai",
      "api_key": "sk-your-actual-openai-key"
    },
    "vk_user2_anthropic": {
      "provider": "anthropic",
      "api_key": "sk-ant-your-actual-anthropic-key"
    }
  }
}
```

> **Note:** Without real API keys, the gateway will start but requests will fail with 502 errors. For testing without API keys, see [TESTING_GUIDE.md](TESTING_GUIDE.md).

## Step 2: Run the Gateway

### Option A: Docker (Recommended)

No Go installation required!

```bash
# Build and run
make docker-run

# View logs
docker logs -f llm-gateway

# Stop when done
make docker-stop
```

<details>
<summary>Or use Docker commands directly</summary>

```bash
# Build image
docker build -t llm-gateway .

# Run container
docker run -d \
  --name llm-gateway \
  -p 8080:8080 \
  -v $(pwd)/keys.json:/app/keys.json \
  llm-gateway

# Stop when done
docker stop llm-gateway && docker rm llm-gateway
```
</details>

### Option B: Build from Source

Requires Go 1.20+:

```bash
make run
```

<details>
<summary>Or build manually</summary>

```bash
go build -o gateway .
./gateway
```
</details>

The gateway will start on http://localhost:8080

## Step 3: Test It

### Quick Health Check

```bash
curl http://localhost:8080/health | jq
```

### Test Chat Completion

```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer vk_user1_openai" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }' | jq
```

> **Note:** This will return 502 without real API keys in keys.json

### View Metrics

```bash
curl http://localhost:8080/metrics | jq
```

### Run Full Test Suite

```bash
./run_all_tests.sh
```

See [TESTING_GUIDE.md](TESTING_GUIDE.md) for complete testing instructions.

## Step 4: Use with OpenAI SDK

### Python

```bash
pip install openai
python examples/python_client.py
```

Or manually:

```python
from openai import OpenAI

client = OpenAI(
    api_key="vk_user1_openai",
    base_url="http://localhost:8080",
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello!"}]
)

print(response.choices[0].message.content)
```

### Go

```bash
go run examples/go_client.go
```

## Environment Variables

Customize the gateway:

```bash
# Different port
export SERVER_PORT=9090

# Disable rate limiting
export QUOTA_ENABLED=false

# Custom quota
export QUOTA_LIMIT=200

# Then run
./gateway
```

## Docker with Custom Settings

```bash
docker run -d \
  --name llm-gateway \
  -p 8080:8080 \
  -e QUOTA_LIMIT=200 \
  -e REQUEST_TIMEOUT=60 \
  -v $(pwd)/keys.json:/app/keys.json \
  llm-gateway
```

## Next Steps

- ✅ Add real API keys to [keys.json](keys.json)
- ✅ Read [TESTING_GUIDE.md](TESTING_GUIDE.md) for comprehensive testing
- ✅ Read [README.md](README.md) for detailed documentation
- ✅ Try the example clients in [examples/](examples/)
- ✅ Check [examples/curl_examples.sh](examples/curl_examples.sh) for more examples

## Troubleshooting

### Port 8080 already in use

```bash
# For Docker
docker run -p 9090:8080 ...

# For Go binary
export SERVER_PORT=9090
./gateway
```

### Gateway won't start - missing keys.json

```bash
cp keys.example.json keys.json
# Edit keys.json with your API keys
```

### Requests fail with 502 Bad Gateway

- You need real API keys in [keys.json](keys.json)
- Get keys from:
  - OpenAI: https://platform.openai.com/api-keys
  - Anthropic: https://console.anthropic.com/settings/keys

### 401 Unauthorized errors

- Verify virtual key exists in [keys.json](keys.json)
- Check Authorization header format: `Bearer vk_user1_openai`
- Virtual key name must match exactly

## Need More Help?

- 📖 Full documentation: [README.md](README.md)
- 🧪 Testing guide: [TESTING_GUIDE.md](TESTING_GUIDE.md)
- 💡 Example code: [examples/](examples/)
