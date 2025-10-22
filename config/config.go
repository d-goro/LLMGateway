package config

import (
	"encoding/json"
	"fmt"
	"llmgateway/internal/models"
	"os"
)

// Config holds the application configuration
type Config struct {
	KeysConfig     models.KeysConfig
	ServerPort     string
	LogToFile      bool
	LogFilePath    string
	QuotaEnabled   bool
	QuotaLimit     int64 // Max requests per hour per virtual key
	RequestTimeout int   // Request timeout in seconds
}

// Load loads the configuration from environment variables
// All config values can be set via environment variables:
// - KEYS_FILE_PATH: path to keys.json (default: "keys.json")
// - SERVER_PORT: server port (default: "8080")
// - LOG_TO_FILE: enable file logging (default: false)
// - LOG_FILE_PATH: log file path (default: "gateway.log")
// - QUOTA_ENABLED: enable rate limiting (default: true)
// - QUOTA_LIMIT: max requests per hour per key (default: 100)
// - REQUEST_TIMEOUT: request timeout in seconds (default: 30)
func Load() (*Config, error) {
	// Get keys file path from environment
	keysFilePath := getEnvOrDefault("KEYS_FILE_PATH", "keys.json")

	// Read keys.json file
	data, err := os.ReadFile(keysFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys file: %w", err)
	}

	// Parse JSON
	var keysConfig models.KeysConfig
	if err := json.Unmarshal(data, &keysConfig); err != nil {
		return nil, fmt.Errorf("failed to parse keys file: %w", err)
	}

	// Validate that we have at least one virtual key
	if len(keysConfig.VirtualKeys) == 0 {
		return nil, fmt.Errorf("no virtual keys configured")
	}

	// Create config with all values from environment variables
	cfg := &Config{
		KeysConfig:     keysConfig,
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		LogToFile:      getEnvBoolOrDefault("LOG_TO_FILE", false),
		LogFilePath:    getEnvOrDefault("LOG_FILE_PATH", "gateway.log"),
		QuotaEnabled:   getEnvBoolOrDefault("QUOTA_ENABLED", true),
		QuotaLimit:     getEnvInt64OrDefault("QUOTA_LIMIT", 100),
		RequestTimeout: getEnvIntOrDefault("REQUEST_TIMEOUT", 30),
	}

	return cfg, nil
}

// ValidateVirtualKey checks if a virtual key exists and returns its configuration
func (c *Config) ValidateVirtualKey(virtualKey string) (models.VirtualKeyConfig, bool) {
	keyConfig, exists := c.KeysConfig.VirtualKeys[virtualKey]
	return keyConfig, exists
}

// Helper functions to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		var intValue int64
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
