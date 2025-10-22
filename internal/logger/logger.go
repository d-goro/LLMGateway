package logger

import (
	"encoding/json"
	"llmgateway/internal/models"
	"log"
	"os"
	"sync"
)

// Logger handles structured JSON logging of LLM interactions
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	toFile   bool
	toStdout bool
}

// NewLogger creates a new logger instance
func NewLogger(logToFile bool, logFilePath string) (*Logger, error) {
	logger := &Logger{
		toFile:   logToFile,
		toStdout: true, // Always log to stdout
	}

	// If logging to file is enabled, open the log file
	if logToFile {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		logger.file = file
	}

	return logger, nil
}

// LogInteraction logs an LLM interaction as structured JSON
func (l *Logger) LogInteraction(entry models.LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Marshal the log entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	// Add newline for readability
	jsonData = append(jsonData, '\n')

	// Write to stdout if enabled
	if l.toStdout {
		os.Stdout.Write(jsonData)
	}

	// Write to file if enabled
	if l.toFile && l.file != nil {
		if _, err := l.file.Write(jsonData); err != nil {
			log.Printf("Failed to write to log file: %v", err)
		}
	}
}

// LogError logs an error message
func (l *Logger) LogError(message string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	errorLog := map[string]any{
		"level":   "error",
		"message": message,
		"error":   err.Error(),
	}

	jsonData, _ := json.Marshal(errorLog)
	jsonData = append(jsonData, '\n')

	if l.toStdout {
		os.Stdout.Write(jsonData)
	}

	if l.toFile && l.file != nil {
		l.file.Write(jsonData)
	}
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(message string, data map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	infoLog := map[string]any{
		"level":   "info",
		"message": message,
	}

	// Add additional data fields
	for k, v := range data {
		infoLog[k] = v
	}

	jsonData, _ := json.Marshal(infoLog)
	jsonData = append(jsonData, '\n')

	if l.toStdout {
		os.Stdout.Write(jsonData)
	}

	if l.toFile && l.file != nil {
		l.file.Write(jsonData)
	}
}

// Close closes the log file if it's open
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
