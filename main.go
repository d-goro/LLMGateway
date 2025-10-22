package main

import (
	"fmt"
	"llmgateway/config"
	"llmgateway/internal/handler"
	"llmgateway/internal/logger"
	"llmgateway/internal/middleware"
	"llmgateway/internal/tracker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration from KEYS_FILE_PATH env var or default "keys.json"
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: Cannot start gateway - failed to load configuration: %v\nPlease ensure keys.json exists and is valid.", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.LogToFile, cfg.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	appLogger.LogInfo("Starting LLM Gateway", map[string]any{
		"port":          cfg.ServerPort,
		"quota_enabled": cfg.QuotaEnabled,
		"quota_limit":   cfg.QuotaLimit,
	})

	// Initialize usage tracker
	usageTracker := tracker.NewTracker(cfg.QuotaEnabled, cfg.QuotaLimit)

	// Initialize handler
	h := handler.NewHandler(cfg, appLogger, usageTracker)

	// Create HTTP server with routes
	mux := http.NewServeMux()

	// Main endpoint - requires authentication
	authMiddleware := middleware.AuthMiddleware(cfg)
	mux.Handle("/chat/completions", authMiddleware(http.HandlerFunc(h.ChatCompletions)))

	// Health and metrics endpoints - no authentication required
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/metrics", h.Metrics)

	// Add a simple root handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"LLM Gateway","version":"1.0.0","endpoints":["/chat/completions","/health","/metrics"]}`)
	})

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		appLogger.LogInfo("Shutting down server...", nil)
		if err := server.Close(); err != nil {
			appLogger.LogError("Error during server shutdown", err)
		}
	}()

	// Start server
	appLogger.LogInfo("Server started", map[string]any{
		"address": "http://localhost:" + cfg.ServerPort,
	})
	fmt.Printf("LLM Gateway listening on port %s\n", cfg.ServerPort)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  POST /chat/completions\n")
	fmt.Printf("  GET  /health\n")
	fmt.Printf("  GET  /metrics\n")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		appLogger.LogError("Server failed to start", err)
		log.Fatalf("Server error: %v", err)
	}
}
