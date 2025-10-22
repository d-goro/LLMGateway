package tracker

import (
	"fmt"
	"llmgateway/internal/models"
	"maps"
	"sync"
	"time"
)

// Tracker manages usage tracking and quota enforcement
type Tracker struct {
	mu              sync.RWMutex
	quotas          map[string]*models.QuotaInfo // Virtual key -> quota info
	quotaLimit      int64
	quotaEnabled    bool
	stats           models.UsageStats
	totalDurationMs int64 // For calculating average
}

// NewTracker creates a new usage tracker
func NewTracker(quotaEnabled bool, quotaLimit int64) *Tracker {
	return &Tracker{
		quotas:       make(map[string]*models.QuotaInfo),
		quotaLimit:   quotaLimit,
		quotaEnabled: quotaEnabled,
		stats: models.UsageStats{
			RequestsByProvider: make(map[models.Provider]int64),
			LastUpdated:        time.Now(),
		},
	}
}

// CheckQuota checks if a virtual key has exceeded its quota
// Returns true if request is allowed, false if quota exceeded
func (t *Tracker) CheckQuota(virtualKey string) (bool, error) {
	if !t.quotaEnabled {
		return true, nil // Quota disabled, allow all requests
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	quota, exists := t.quotas[virtualKey]

	if !exists {
		// First request for this key
		t.quotas[virtualKey] = &models.QuotaInfo{
			RequestCount: 1,
			WindowStart:  now,
			MaxRequests:  t.quotaLimit,
		}
		return true, nil
	}

	// Check if we're still in the same hour window
	windowEnd := quota.WindowStart.Add(1 * time.Hour)
	if now.After(windowEnd) {
		// New window, reset counter
		quota.RequestCount = 1
		quota.WindowStart = now
		return true, nil
	}

	// Check if quota exceeded
	if quota.RequestCount >= quota.MaxRequests {
		return false, fmt.Errorf("quota exceeded: %d requests per hour limit reached", quota.MaxRequests)
	}

	// Increment counter
	quota.RequestCount++
	return true, nil
}

// RecordRequest records a completed request for statistics
func (t *Tracker) RecordRequest(provider models.Provider, durationMs int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stats.TotalRequests++
	t.stats.RequestsByProvider[provider]++
	t.totalDurationMs += durationMs

	// Calculate average response time
	if t.stats.TotalRequests > 0 {
		t.stats.AverageResponseMs = float64(t.totalDurationMs) / float64(t.stats.TotalRequests)
	}

	t.stats.LastUpdated = time.Now()
}

// GetStats returns current usage statistics
func (t *Tracker) GetStats() models.UsageStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Create a copy to avoid race conditions
	statsCopy := models.UsageStats{
		TotalRequests:      t.stats.TotalRequests,
		RequestsByProvider: make(map[models.Provider]int64),
		AverageResponseMs:  t.stats.AverageResponseMs,
		LastUpdated:        t.stats.LastUpdated,
	}

	maps.Copy(statsCopy.RequestsByProvider, t.stats.RequestsByProvider)

	return statsCopy
}
