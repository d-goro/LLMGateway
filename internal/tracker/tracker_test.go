package tracker

import (
	"llmgateway/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracker(t *testing.T) {
	tracker := NewTracker(true, 100)
	require.NotNil(t, tracker)
	assert.True(t, tracker.quotaEnabled)
	assert.Equal(t, int64(100), tracker.quotaLimit)
}

func TestCheckQuotaDisabled(t *testing.T) {
	tracker := NewTracker(false, 100)

	// When quota is disabled, all requests should be allowed
	for i := 0; i < 200; i++ {
		allowed, err := tracker.CheckQuota("test_key")
		assert.True(t, allowed, "Request %d should be allowed when quota disabled", i)
		require.NoError(t, err)
	}
}

func TestCheckQuotaEnabled(t *testing.T) {
	tracker := NewTracker(true, 10)

	// First 10 requests should be allowed
	for i := 0; i < 10; i++ {
		allowed, err := tracker.CheckQuota("test_key")
		assert.True(t, allowed, "Request %d should be allowed", i)
		require.NoError(t, err)
	}

	// 11th request should be denied
	allowed, err := tracker.CheckQuota("test_key")
	assert.False(t, allowed, "11th request should be denied")
	require.Error(t, err, "Expected error for quota exceeded")
}

func TestCheckQuotaPerKey(t *testing.T) {
	tracker := NewTracker(true, 5)

	// Each key should have its own quota
	for i := 0; i < 5; i++ {
		allowed1, err1 := tracker.CheckQuota("key1")
		allowed2, err2 := tracker.CheckQuota("key2")

		assert.True(t, allowed1, "key1 request %d should be allowed", i)
		assert.True(t, allowed2, "key2 request %d should be allowed", i)
		require.NoError(t, err1)
		require.NoError(t, err2)
	}

	// Both keys should now be at limit
	allowed1, _ := tracker.CheckQuota("key1")
	allowed2, _ := tracker.CheckQuota("key2")

	assert.False(t, allowed1, "key1 should exceed quota")
	assert.False(t, allowed2, "key2 should exceed quota")
}

func TestRecordRequest(t *testing.T) {
	tracker := NewTracker(false, 0)

	// Record some requests
	tracker.RecordRequest(models.ProviderOpenAI, 100)
	tracker.RecordRequest(models.ProviderOpenAI, 200)
	tracker.RecordRequest(models.ProviderAnthropic, 150)

	stats := tracker.GetStats()

	assert.Equal(t, int64(3), stats.TotalRequests)
	assert.Equal(t, int64(2), stats.RequestsByProvider[models.ProviderOpenAI])
	assert.Equal(t, int64(1), stats.RequestsByProvider[models.ProviderAnthropic])
	assert.Equal(t, 150.0, stats.AverageResponseMs)
}

func TestQuotaWindowReset(t *testing.T) {
	tracker := NewTracker(true, 5)

	// Use up quota
	for i := 0; i < 5; i++ {
		tracker.CheckQuota("test_key")
	}

	// Should be denied now
	allowed, _ := tracker.CheckQuota("test_key")
	assert.False(t, allowed, "Request should be denied")

	// Manually set window start to past hour to simulate time passage
	tracker.mu.Lock()
	tracker.quotas["test_key"].WindowStart = time.Now().Add(-2 * time.Hour)
	tracker.mu.Unlock()

	// Should be allowed again after window reset
	allowed, _ = tracker.CheckQuota("test_key")
	assert.True(t, allowed, "Request should be allowed after window reset")
}
