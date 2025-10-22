package middleware

import (
	"context"
	"encoding/json"
	"llmgateway/config"
	"llmgateway/internal/models"
	"net/http"
	"strings"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// VirtualKeyContextKey is the key for storing the virtual key in context
	VirtualKeyContextKey ContextKey = "virtualKey"
	// KeyConfigContextKey is the key for storing the key config in context
	KeyConfigContextKey ContextKey = "keyConfig"
)

// AuthMiddleware validates the virtual API key from the Authorization header
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "missing Authorization header")
				return
			}

			// Parse the Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeJSONError(w, http.StatusUnauthorized, "invalid Authorization header format")
				return
			}

			virtualKey := parts[1]

			// Validate the virtual key
			keyConfig, valid := cfg.ValidateVirtualKey(virtualKey)
			if !valid {
				writeJSONError(w, http.StatusUnauthorized, "invalid virtual key")
				return
			}

			// Store the virtual key and config in the request context
			ctx := context.WithValue(r.Context(), VirtualKeyContextKey, virtualKey)
			ctx = context.WithValue(ctx, KeyConfigContextKey, keyConfig)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetVirtualKey retrieves the virtual key from the request context
func GetVirtualKey(ctx context.Context) (string, bool) {
	virtualKey, ok := ctx.Value(VirtualKeyContextKey).(string)
	return virtualKey, ok
}

// GetKeyConfig retrieves the key configuration from the request context
func GetKeyConfig(ctx context.Context) (models.VirtualKeyConfig, bool) {
	keyConfig, ok := ctx.Value(KeyConfigContextKey).(models.VirtualKeyConfig)
	return keyConfig, ok
}

// writeJSONError writes a JSON error response
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]any{
		"error": map[string]any{
			"message": message,
			"type":    "invalid_request_error",
			"code":    statusCode,
		},
	}
	json.NewEncoder(w).Encode(response)
}
