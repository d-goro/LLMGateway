package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ChatRequest represents the request structure for chat completions
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a single message in the chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the response from the chat completion API
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single choice in the response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GatewayClient is a simple HTTP client for the LLM Gateway
type GatewayClient struct {
	BaseURL    string
	VirtualKey string
	HTTPClient *http.Client
}

// NewGatewayClient creates a new gateway client
func NewGatewayClient(baseURL, virtualKey string) *GatewayClient {
	return &GatewayClient{
		BaseURL:    baseURL,
		VirtualKey: virtualKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatCompletion sends a chat completion request through the gateway
func (c *GatewayClient) ChatCompletion(req ChatRequest) (*ChatResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.VirtualKey)

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResp, nil
}

// GetHealth checks the health of the gateway
func (c *GatewayClient) GetHealth() (map[string]any, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/health")
	if err != nil {
		return nil, fmt.Errorf("failed to get health: %w", err)
	}
	defer resp.Body.Close()

	var health map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode health response: %w", err)
	}

	return health, nil
}

// GetMetrics retrieves usage metrics from the gateway
func (c *GatewayClient) GetMetrics() (map[string]any, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	var metrics map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode metrics response: %w", err)
	}

	return metrics, nil
}

func main() {
	fmt.Println("LLM Gateway Go Client Example")
	fmt.Println("===============================")

	// Create a client with OpenAI virtual key
	client := NewGatewayClient("http://localhost:8080", "vk_user1_openai")

	// Test chat completion
	fmt.Println("\nSending chat completion request...")
	req := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello! Can you tell me about Go programming language?"},
		},
	}

	resp, err := client.ChatCompletion(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\nResponse received:")
	if len(resp.Choices) > 0 {
		fmt.Println(resp.Choices[0].Message.Content)
		fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)
	}

	// Test health endpoint
	fmt.Println("\n" + "===============================")
	fmt.Println("Checking gateway health...")
	health, err := client.GetHealth()
	if err != nil {
		fmt.Printf("Error getting health: %v\n", err)
	} else {
		healthJSON, _ := json.MarshalIndent(health, "", "  ")
		fmt.Printf("Health status:\n%s\n", healthJSON)
	}

	// Test metrics endpoint
	fmt.Println("\n" + "===============================")
	fmt.Println("Fetching metrics...")
	metrics, err := client.GetMetrics()
	if err != nil {
		fmt.Printf("Error getting metrics: %v\n", err)
	} else {
		metricsJSON, _ := json.MarshalIndent(metrics, "", "  ")
		fmt.Printf("Metrics:\n%s\n", metricsJSON)
	}

	fmt.Println("\nAll tests completed!")
}
