package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderEndpoint(t *testing.T) {
	tests := []struct {
		provider Provider
		expected string
	}{
		{ProviderOpenAI, "https://api.openai.com/v1/chat/completions"},
		{ProviderAnthropic, "https://api.anthropic.com/v1/messages"},
		{Provider("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			endpoint := tt.provider.Endpoint()
			assert.Equal(t, tt.expected, endpoint)
		})
	}
}

func TestProviderConstants(t *testing.T) {
	assert.Equal(t, "openai", string(ProviderOpenAI))
	assert.Equal(t, "anthropic", string(ProviderAnthropic))
}
