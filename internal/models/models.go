package models

import "time"

// Provider represents an LLM provider type
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
)

// VirtualKeyConfig represents the configuration for a single virtual key
type VirtualKeyConfig struct {
	Provider Provider `json:"provider"`
	APIKey   string   `json:"api_key"`
}

// KeysConfig represents the structure of keys.json file
type KeysConfig struct {
	VirtualKeys map[string]VirtualKeyConfig `json:"virtual_keys"`
}

// LogEntry represents a single LLM interaction log entry
type LogEntry struct {
	Timestamp  string         `json:"timestamp"`
	VirtualKey string         `json:"virtual_key"`
	Provider   Provider       `json:"provider"`
	Method     string         `json:"method"`
	Status     int            `json:"status"`
	DurationMs int64          `json:"duration_ms"`
	Request    map[string]any `json:"request,omitempty"`
	Response   map[string]any `json:"response,omitempty"`
	Error      string         `json:"error,omitempty"`
}

// ProviderEndpoint returns the API endpoint URL for a given provider
func (p Provider) Endpoint() string {
	switch p {
	case ProviderOpenAI:
		return "https://api.openai.com/v1/chat/completions"
	case ProviderAnthropic:
		return "https://api.anthropic.com/v1/messages"
	default:
		return ""
	}
}

// UsageStats tracks usage statistics for metrics
type UsageStats struct {
	TotalRequests      int64              `json:"total_requests"`
	RequestsByProvider map[Provider]int64 `json:"requests_by_provider"`
	AverageResponseMs  float64            `json:"average_response_ms"`
	LastUpdated        time.Time          `json:"last_updated"`
}

// QuotaInfo tracks rate limiting information per virtual key
type QuotaInfo struct {
	RequestCount int64
	WindowStart  time.Time
	MaxRequests  int64
}
