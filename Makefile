.PHONY: build run test clean docker-build docker-run help

# Binary name
BINARY_NAME=gateway

# Build the application
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Run the application
run: build
	@echo "Starting gateway..."
	@./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t llm-gateway:latest .
	@echo "Docker image built: llm-gateway:latest"

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -d \
		--name llm-gateway \
		-p 8080:8080 \
		-v $(PWD)/keys.json:/app/keys.json \
		llm-gateway:latest
	@echo "Container started: llm-gateway"
	@echo "Gateway available at http://localhost:8080"

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop llm-gateway || true
	@docker rm llm-gateway || true
	@echo "Container stopped"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run || go vet ./...
	@echo "Lint complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Setup: create keys.json from example
setup:
	@if [ ! -f keys.json ]; then \
		echo "Creating keys.json from example..."; \
		cp keys.example.json keys.json; \
		echo "Please edit keys.json with your actual API keys"; \
	else \
		echo "keys.json already exists"; \
	fi

# Help
help:
	@echo "LLM Gateway - Available Make Targets:"
	@echo ""
	@echo "  make build              - Build the application"
	@echo "  make run                - Build and run the application"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make clean              - Clean build artifacts"
	@echo ""
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Build and run Docker container"
	@echo "  make docker-stop        - Stop Docker container"
	@echo ""
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Run linter"
	@echo "  make deps               - Download and update dependencies"
	@echo "  make setup              - Create keys.json from example"
	@echo "  make help               - Show this help message"
