# LLM Gateway

A lightweight, high-performance LLM gateway service in Go that routes requests to OpenAI and Anthropic APIs, manages virtual API keys, and logs all interactions.

## Features

### Core Features
- **Unified API Endpoint**: Single `/chat/completions` endpoint for all LLM providers
- **Virtual Key Management**: Map virtual API keys to actual provider keys with automatic routing
- **Multi-Provider Support**: OpenAI and Anthropic integration
- **Request Proxying**: Transparent forwarding with header management
- **Structured Logging**: JSON-formatted logs of all interactions
- **Configuration-Based**: Simple JSON configuration for key management

### Bonus Features
- **Usage Tracking**: Track request counts and token usage per virtual key
- **Rate Limiting**: Configurable quotas (default: 100 requests/hour per key)
- **Health Checks**: Monitor gateway and provider availability
- **Metrics Endpoint**: Real-time usage statistics
- **Request Validation**: Pre-flight validation of request format
- **Concurrent Handling**: Efficient goroutine-based request processing
- **Request Timeouts**: Configurable timeout handling

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get up and running in 5 minutes
- **[TESTING_GUIDE.md](TESTING_GUIDE.md)** - Complete testing guide with examples
- **[test_with_real_keys.md](test_with_real_keys.md)** - Testing with real API keys

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ Virtual Key: vk_user1_openai
       │
       ▼
┌─────────────────────────────────────┐
│        LLM Gateway (Port 8080)      │
│                                     │
│  ┌──────────────────────────────┐  │
│  │   Auth Middleware            │  │
│  │   - Validate Virtual Key     │  │
│  │   - Map to Provider + API Key│  │
│  └──────────────────────────────┘  │
│                                     │
│  ┌──────────────────────────────┐  │
│  │   Handler                    │  │
│  │   - Check Quota              │  │
│  │   - Validate Request         │  │
│  │   - Proxy to Provider        │  │
│  │   - Log Interaction          │  │
│  └──────────────────────────────┘  │
│                                     │
└─────────┬───────────────┬───────────┘
          │               │
          ▼               ▼
    ┌──────────┐    ┌──────────┐
    │  OpenAI  │    │Anthropic │
    │   API    │    │   API    │
    └──────────┘    └──────────┘
```

## Project Structure

```
.
├── main.go                      # Application entry point
├── config/
│   ├── config.go                # Configuration management
│   └── config_test.go           # Config tests
├── internal/
│   ├── handler/
│   │   └── handler.go           # HTTP request handlers
│   ├── logger/
│   │   └── logger.go            # Structured JSON logging
│   ├── middleware/
│   │   └── auth.go              # Authentication middleware
│   ├── models/
│   │   ├── models.go            # Data models
│   │   └── models_test.go       # Model tests
│   ├── proxy/
│   │   ├── proxy.go             # Provider proxy logic
│   │   └── proxy_test.go        # Proxy tests
│   └── tracker/
│       ├── tracker.go           # Usage tracking and quotas
│       └── tracker_test.go      # Tracker tests
├── examples/
│   ├── python_client.py         # Python example
│   ├── go_client.go             # Go example
│   └── curl_examples.sh         # cURL examples
├── test_health.sh               # E2E: Health check test
├── test_metrics.sh              # E2E: Metrics endpoint test
├── test_chat.sh                 # E2E: Chat completions test
├── test_invalid_key.sh          # E2E: Invalid auth test
├── test_no_auth.sh              # E2E: Missing auth test
├── run_all_tests.sh             # E2E: Run all tests
├── keys.example.json            # Example configuration
├── Dockerfile                   # Container image
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
├── README.md                    # This file
├── QUICKSTART.md                # Quick start guide
├── TESTING_GUIDE.md             # Complete testing guide
└── test_with_real_keys.md       # Real API keys testing
```

## Quick Start

### Prerequisites

**Choose one:**
- **Docker** (recommended - no Go installation needed), OR
- **Go 1.20+** (to build from source)

**API Keys (optional for testing):**
- OpenAI and/or Anthropic API keys (not required to run tests with example keys)


### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd LLMGateway
```

2. Create your `keys.json` configuration file:
```bash
make setup
# Edit keys.json with your actual API keys
```

3. Build and run the application:
```bash
make run
```

Or using Docker:
```bash
make docker-run
```

The server will start on `http://localhost:8080`.

### Configuration

Create a `keys.json` file in the project root:

```json
{
  "virtual_keys": {
    "vk_user1_openai": {
      "provider": "openai",
      "api_key": "sk-your-real-openai-key-here"
    },
    "vk_user2_anthropic": {
      "provider": "anthropic",
      "api_key": "sk-ant-your-real-anthropic-key-here"
    },
    "vk_admin_openai": {
      "provider": "openai",
      "api_key": "sk-another-openai-key-here"
    }
  }
}
```

### Environment Variables

The gateway supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `KEYS_FILE_PATH` | `keys.json` | Path to the keys configuration file |
| `SERVER_PORT` | `8080` | Server port |
| `LOG_TO_FILE` | `false` | Enable logging to file |
| `LOG_FILE_PATH` | `gateway.log` | Path to log file |
| `QUOTA_ENABLED` | `true` | Enable rate limiting |
| `QUOTA_LIMIT` | `100` | Max requests per hour per key |
| `REQUEST_TIMEOUT` | `30` | Request timeout in seconds |

Example:
```bash
export QUOTA_LIMIT=200
export REQUEST_TIMEOUT=60
./gateway
```

## Usage

### API Endpoints

#### POST /chat/completions

Main endpoint for chat completions. Routes requests to the appropriate provider based on the virtual key.

**Headers:**
- `Authorization: Bearer <virtual-key>` (required)
- `Content-Type: application/json` (required)

**Request Body:**
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "max_tokens": 100
}
```

**Response:**
Returns the provider's response unchanged.

**Error Responses:**
- `401`: Invalid or missing virtual key
- `400`: Invalid request format
- `429`: Quota exceeded
- `502`: Provider request failed

#### GET /health

Health check endpoint. Returns gateway status and provider availability.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "providers": {
    "openai": {
      "healthy": true
    },
    "anthropic": {
      "healthy": true
    }
  }
}
```

#### GET /metrics

Returns usage statistics.

**Response:**
```json
{
  "total_requests": 150,
  "requests_by_provider": {
    "openai": 100,
    "anthropic": 50
  },
  "average_response_ms": 1250.5,
  "last_updated": "2024-01-15T10:30:00Z"
}
```

### Example Clients

#### Python (using OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    api_key="vk_user1_openai",  # Your virtual key
    base_url="http://localhost:8080",
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)
```

See [examples/python_client.py](examples/python_client.py) for a complete example.

#### cURL

```bash
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

See [examples/curl_examples.sh](examples/curl_examples.sh) for more examples.

#### Go

See [examples/go_client.go](examples/go_client.go) for a complete Go client implementation.

## Logging

All LLM interactions are logged as structured JSON:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "virtual_key": "vk_user1_openai",
  "provider": "openai",
  "method": "POST",
  "status": 200,
  "duration_ms": 1250,
  "request": {
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  },
  "response": {
    "choices": [
      {
        "message": {
          "role": "assistant",
          "content": "Hello! How can I help you?"
        }
      }
    ]
  }
}
```

Logs are written to stdout by default. Enable file logging with `LOG_TO_FILE=true`.

## Rate Limiting

The gateway includes built-in rate limiting:

- Configurable quota per virtual key (default: 100 requests/hour)
- Sliding window implementation
- Returns `429 Too Many Requests` when quota exceeded
- Independent quotas for each virtual key

Disable rate limiting:
```bash
QUOTA_ENABLED=false ./gateway
```

Adjust quota limit:
```bash
QUOTA_LIMIT=200 ./gateway
```

## Testing

The project includes comprehensive testing at multiple levels:

### Unit Tests

Run the Go test suite:

```bash
# Using Make
make test

# Or directly with go
go test ./... -v

# With coverage
make test-coverage
```

Unit test coverage includes:
- Configuration loading and validation
- Virtual key authentication
- Request validation
- Usage tracking and quota enforcement
- Provider endpoint resolution
- Unsupported provider handling

### End-to-End Tests

Test the running gateway with real HTTP requests:

```bash
# Run all E2E tests (autonomous - starts gateway if needed)
./run_all_tests.sh

# Run with auto-cleanup
./run_all_tests.sh --cleanup

# Or individual tests (require gateway to be running)
./test_health.sh          # Health check endpoint
./test_metrics.sh         # Metrics endpoint
./test_chat.sh            # Chat completions
./test_invalid_key.sh     # Invalid authentication
./test_no_auth.sh         # Missing authentication
```

**Note**: `run_all_tests.sh` automatically starts the gateway with Docker if it's not running. Individual test scripts require the gateway to be started manually.

### Testing Guides

- **[TESTING_GUIDE.md](TESTING_GUIDE.md)** - Complete testing guide with examples and troubleshooting
- **[test_with_real_keys.md](test_with_real_keys.md)** - Guide for testing with real OpenAI/Anthropic API keys
- **[QUICKSTART.md](QUICKSTART.md)** - Quick start guide for getting up and running

## Docker

### Quick Start with Docker

```bash
# Build and run
make docker-run

# View logs
docker logs -f llm-gateway

# Stop the container
make docker-stop
```

### Advanced Docker Usage

<details>
<summary>Run with custom environment variables</summary>

```bash
# Build image first
make docker-build

# Run with custom settings
docker run -d \
  --name llm-gateway \
  -p 8080:8080 \
  -e QUOTA_LIMIT=200 \
  -e REQUEST_TIMEOUT=60 \
  -v $(pwd)/keys.json:/app/keys.json \
  llm-gateway:latest

# View logs
docker logs -f llm-gateway

# Stop the container
make docker-stop
```
</details>

## Performance

- **Concurrent Request Handling**: Uses Go goroutines for efficient concurrent processing
- **Low Memory Footprint**: Minimal memory usage, suitable for containerized environments
- **Request Timeout**: Prevents hanging requests (default: 30s)
- **Connection Pooling**: Reuses HTTP connections for improved performance

## Security Considerations

- Virtual keys should be treated as secrets
- Use HTTPS in production (reverse proxy recommended)
- Keep `keys.json` secure and out of version control
- Regularly rotate API keys
- Monitor logs for suspicious activity
- Consider implementing IP whitelisting for production

## Production Deployment

### Recommended Setup

1. **Reverse Proxy**: Use nginx or similar for:
   - HTTPS termination
   - Rate limiting (additional layer)
   - Load balancing
   - Access logging

2. **Monitoring**: Integrate with:
   - Prometheus (metrics)
   - ELK Stack (log aggregation)
   - Health check monitoring

3. **Scaling**: Deploy multiple instances behind a load balancer

4. **Configuration**:
   ```bash
   export LOG_TO_FILE=true
   export LOG_FILE_PATH=/var/log/gateway/gateway.log
   export QUOTA_ENABLED=true
   export QUOTA_LIMIT=1000
   export REQUEST_TIMEOUT=30
   ```

## Development

### Building from Source

```bash
# Clone repository
git clone <repository-url>
cd LLMGateway

# Install dependencies
make deps

# Run tests
make test

# Build and run
make run

# Or build Docker
make docker-run
```

### Code Structure

The project follows Go best practices:
- Clear package separation
- Interface-based design for extensibility
- Comprehensive error handling
- Thread-safe implementations
- Well-documented code

### Adding a New Provider

1. Add provider constant in [internal/models/models.go](internal/models/models.go)
2. Add endpoint in `Provider.Endpoint()` method
3. Add provider-specific headers in [internal/proxy/proxy.go](internal/proxy/proxy.go)
4. Update documentation

## Troubleshooting

### Common Issues

**Gateway won't start:**
- Check if port 8080 is available
- Verify `keys.json` exists and is valid JSON
- Ensure at least one virtual key is configured

**401 Unauthorized:**
- Verify virtual key exists in `keys.json`
- Check Authorization header format: `Bearer <key>`

**429 Too Many Requests:**
- Check quota limit with `/metrics` endpoint
- Adjust `QUOTA_LIMIT` if needed
- Wait for quota window to reset (1 hour)

**502 Bad Gateway:**
- Verify actual API keys are valid
- Check provider API status
- Review gateway logs for error details
- Check network connectivity to provider APIs

## License

MIT License

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Support

For issues and questions:
- Check the [Troubleshooting](#troubleshooting) section
- Review example code in [examples/](examples/)
- Open an issue on GitHub

## Acknowledgments

Built with Go and designed for simplicity, performance, and extensibility.
