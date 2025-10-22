package proxy

import (
	"context"
	"llmgateway/internal/models"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRequestFormat(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
	}{
		{
			name: "valid request",
			requestBody: `{
				"model": "gpt-3.5-turbo",
				"messages": [{"role": "user", "content": "Hello"}]
			}`,
			expectError: false,
		},
		{
			name:        "missing model",
			requestBody: `{"messages": [{"role": "user", "content": "Hello"}]}`,
			expectError: true,
		},
		{
			name:        "missing messages",
			requestBody: `{"model": "gpt-3.5-turbo"}`,
			expectError: true,
		},
		{
			name:        "invalid json",
			requestBody: `{invalid json}`,
			expectError: true,
		},
		{
			name:        "empty object",
			requestBody: `{}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequestFormat([]byte(tt.requestBody))
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProxyRequestUnsupportedProvider(t *testing.T) {
	ctx := context.Background()
	unsupportedProvider := models.Provider("unsupported")
	requestBody := []byte(`{"model":"test","messages":[{"role":"user","content":"test"}]}`)

	_, statusCode, err := ProxyRequest(ctx, unsupportedProvider, "test-key", requestBody, http.Header{}, 5*time.Second)

	require.Error(t, err)
	assert.Equal(t, 0, statusCode)
	assert.Contains(t, err.Error(), "unsupported provider")
}

func TestCheckProviderHealthUnsupportedProvider(t *testing.T) {
	unsupportedProvider := models.Provider("unsupported")

	healthy, err := CheckProviderHealth(unsupportedProvider, "test-key")

	require.Error(t, err)
	assert.False(t, healthy)
	assert.Contains(t, err.Error(), "unsupported provider")
}
