package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llmgateway/internal/models"
	"net/http"
	"time"
)

// ProxyRequest forwards a request to the appropriate LLM provider
func ProxyRequest(
	ctx context.Context,
	provider models.Provider,
	apiKey string,
	requestBody []byte,
	originalHeaders http.Header,
	timeout time.Duration,
) (responseBody []byte, statusCode int, err error) {
	// Get the provider endpoint
	endpoint := provider.Endpoint()
	if endpoint == "" {
		return nil, 0, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers from original request (except Authorization)
	for key, values := range originalHeaders {
		// Skip Authorization header as we'll set it with the real API key
		if key == "Authorization" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set the appropriate authorization header based on provider
	switch provider {
	case models.ProviderOpenAI:
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
	case models.ProviderAnthropic:
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		return nil, 0, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request to %s: %w", provider, err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

// CheckProviderHealth checks if a provider's API is reachable
func CheckProviderHealth(provider models.Provider, apiKey string) (bool, error) {
	endpoint := provider.Endpoint()
	if endpoint == "" {
		return false, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Create a minimal test request based on provider
	var testBody []byte
	var err error

	switch provider {
	case models.ProviderOpenAI:
		// Minimal OpenAI request
		testBody, err = json.Marshal(map[string]any{
			"model":      "gpt-3.5-turbo",
			"messages":   []map[string]string{{"role": "user", "content": "test"}},
			"max_tokens": 1,
		})
	case models.ProviderAnthropic:
		// Minimal Anthropic request
		testBody, err = json.Marshal(map[string]any{
			"model":      "claude-3-haiku-20240307",
			"messages":   []map[string]string{{"role": "user", "content": "test"}},
			"max_tokens": 1,
		})
	default:
		return false, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return false, fmt.Errorf("failed to create test request: %w", err)
	}

	// Try to send a request with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, statusCode, err := ProxyRequest(ctx, provider, apiKey, testBody, http.Header{}, 5*time.Second)
	if err != nil {
		return false, err
	}

	// Consider the provider healthy if we get a response (even if it's an error due to invalid request)
	// Status codes in the 200-499 range indicate the API is reachable
	return statusCode >= 200 && statusCode < 500, nil
}

// ValidateRequestFormat performs basic validation on the request body
func ValidateRequestFormat(requestBody []byte) error {
	var req map[string]any
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Check for required fields
	if _, hasModel := req["model"]; !hasModel {
		return fmt.Errorf("missing required field: model")
	}

	if _, hasMessages := req["messages"]; !hasMessages {
		return fmt.Errorf("missing required field: messages")
	}

	return nil
}
