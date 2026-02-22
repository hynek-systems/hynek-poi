package cache

import (
	"testing"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestMemoryCache_GetSet(t *testing.T) {
	cache := NewMemoryCache()
	key := "test-key"
	pois := []domain.POI{
		{
			ID:        "1",
			Name:      "Test POI",
			Latitude:  59.3293,
			Longitude: 18.0686,
			Category:  "restaurant",
			Source:    "test",
		},
	}

	// Test cache miss
	_, found := cache.Get(key)
	if found {
		t.Error("Expected cache miss, got hit")
	}

	// Test cache set and hit
	cache.Set(key, pois, 1*time.Hour)
	result, found := cache.Get(key)
	if !found {
		t.Fatal("Expected cache hit, got miss")
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 POI, got %d", len(result))
	}

	if result[0].ID != pois[0].ID {
		t.Errorf("Expected ID %s, got %s", pois[0].ID, result[0].ID)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache()
	key := "expiring-key"
	pois := []domain.POI{
		{ID: "1", Name: "Test"},
	}

	// Set with short TTL
	cache.Set(key, pois, 100*time.Millisecond)

	// Should be available immediately
	_, found := cache.Get(key)
	if !found {
		t.Error("Expected cache hit immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, found = cache.Get(key)
	if found {
		t.Error("Expected cache miss after expiration")
	}
}

func TestMemoryCache_Concurrent(t *testing.T) {
	cache := NewMemoryCache()
	pois := []domain.POI{{ID: "1", Name: "Test"}}

	// Write concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := "key"
			cache.Set(key, pois, 1*time.Minute)
			done <- true
		}(i)
	}

	// Wait for all writes
	for i := 0; i < 10; i++ {
		<-done
	}

	// Read concurrently
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = cache.Get("key")
			done <- true
		}()
	}

	// Wait for all reads
	for i := 0; i < 10; i++ {
		<-done
	}
}
