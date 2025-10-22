package handler

import (
	"context"
	"encoding/json"
	"io"
	"llmgateway/config"
	"llmgateway/internal/logger"
	"llmgateway/internal/middleware"
	"llmgateway/internal/models"
	"llmgateway/internal/proxy"
	"llmgateway/internal/tracker"
	"net/http"
	"time"
)

// Handler manages HTTP request handling
type Handler struct {
	config  *config.Config
	logger  *logger.Logger
	tracker *tracker.Tracker
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, log *logger.Logger, track *tracker.Tracker) *Handler {
	return &Handler{
		config:  cfg,
		logger:  log,
		tracker: track,
	}
}

// ChatCompletions handles the /chat/completions endpoint
func (h *Handler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get virtual key and config from context (set by auth middleware)
	virtualKey, ok := middleware.GetVirtualKey(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "authentication failed")
		return
	}

	keyConfig, ok := middleware.GetKeyConfig(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "authentication failed")
		return
	}

	// Check quota if enabled
	if h.config.QuotaEnabled {
		allowed, err := h.tracker.CheckQuota(virtualKey)
		if !allowed {
			h.writeError(w, http.StatusTooManyRequests, err.Error())
			return
		}
	}

	// Read the request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	// Validate request format
	if err := proxy.ValidateRequestFormat(requestBody); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request format: "+err.Error())
		return
	}

	// Parse request body for logging
	var requestData map[string]any
	json.Unmarshal(requestBody, &requestData)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.config.RequestTimeout)*time.Second)
	defer cancel()

	// Proxy the request to the appropriate provider
	responseBody, statusCode, err := proxy.ProxyRequest(
		ctx,
		keyConfig.Provider,
		keyConfig.APIKey,
		requestBody,
		r.Header,
		time.Duration(h.config.RequestTimeout)*time.Second,
	)

	duration := time.Since(startTime)
	durationMs := duration.Milliseconds()

	// Parse response body for logging
	var responseData map[string]any
	if err == nil && len(responseBody) > 0 {
		json.Unmarshal(responseBody, &responseData)
	}

	// Create log entry
	logEntry := models.LogEntry{
		Timestamp:  startTime.Format(time.RFC3339),
		VirtualKey: virtualKey,
		Provider:   keyConfig.Provider,
		Method:     r.Method,
		Status:     statusCode,
		DurationMs: durationMs,
		Request:    requestData,
		Response:   responseData,
	}

	if err != nil {
		logEntry.Error = err.Error()
		logEntry.Status = http.StatusBadGateway
		h.logger.LogInteraction(logEntry)
		h.writeError(w, http.StatusBadGateway, "failed to proxy request: "+err.Error())
		return
	}

	// Record the request in tracker for statistics
	h.tracker.RecordRequest(keyConfig.Provider, durationMs)

	// Log the interaction
	h.logger.LogInteraction(logEntry)

	// Write the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}

// Health handles the /health endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	health := map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"providers": make(map[string]any),
	}

	// Check each provider's health using the first available key for that provider
	providerKeys := make(map[models.Provider]string)
	for _, keyConfig := range h.config.KeysConfig.VirtualKeys {
		if _, exists := providerKeys[keyConfig.Provider]; !exists {
			providerKeys[keyConfig.Provider] = keyConfig.APIKey
		}
	}

	allHealthy := true
	for provider, apiKey := range providerKeys {
		healthy, err := proxy.CheckProviderHealth(provider, apiKey)
		providerStatus := map[string]any{
			"healthy": healthy,
		}
		if err != nil {
			providerStatus["error"] = err.Error()
			allHealthy = false
		}
		health["providers"].(map[string]any)[string(provider)] = providerStatus
	}

	if !allHealthy {
		health["status"] = "degraded"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Metrics handles the /metrics endpoint
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	stats := h.tracker.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// writeError writes a JSON error response
func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := map[string]any{
		"error": map[string]any{
			"message": message,
			"type":    "api_error",
			"code":    statusCode,
		},
	}
	json.NewEncoder(w).Encode(errorResponse)
}
