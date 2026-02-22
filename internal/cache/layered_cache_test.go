package cache

import (
	"testing"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestLayeredCache_L1Hit(t *testing.T) {
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	cache := NewLayeredCache(l1, l2)

	key := "test-key"
	pois := []domain.POI{{ID: "1", Name: "Test"}}

	// Set in L1 only
	l1.Set(key, pois, 1*time.Hour)

	// Should get from L1
	result, found := cache.Get(key)
	if !found {
		t.Fatal("Expected cache hit from L1")
	}

	if len(result) != 1 || result[0].ID != "1" {
		t.Error("Got incorrect data from L1")
	}
}

func TestLayeredCache_L2HitWithL1Promotion(t *testing.T) {
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	cache := NewLayeredCache(l1, l2)

	key := "test-key"
	pois := []domain.POI{{ID: "2", Name: "Test L2"}}

	// Set in L2 only
	l2.Set(key, pois, 1*time.Hour)

	// First get should hit L2 and promote to L1
	result, found := cache.Get(key)
	if !found {
		t.Fatal("Expected cache hit from L2")
	}

	if len(result) != 1 || result[0].ID != "2" {
		t.Error("Got incorrect data from L2")
	}

	// Verify it was promoted to L1
	resultL1, foundL1 := l1.Get(key)
	if !foundL1 {
		t.Error("Expected data to be promoted to L1")
	}

	if len(resultL1) != 1 || resultL1[0].ID != "2" {
		t.Error("Promoted data in L1 is incorrect")
	}
}

func TestLayeredCache_Miss(t *testing.T) {
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	cache := NewLayeredCache(l1, l2)

	// Try to get non-existent key
	_, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected cache miss for nonexistent key")
	}
}

func TestLayeredCache_SetBothLayers(t *testing.T) {
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	cache := NewLayeredCache(l1, l2)

	key := "test-key"
	pois := []domain.POI{{ID: "3", Name: "Test Both"}}

	// Set through layered cache
	cache.Set(key, pois, 1*time.Hour)

	// Verify both layers have the data
	resultL1, foundL1 := l1.Get(key)
	if !foundL1 {
		t.Error("Expected data in L1")
	}

	resultL2, foundL2 := l2.Get(key)
	if !foundL2 {
		t.Error("Expected data in L2")
	}

	if len(resultL1) != 1 || resultL1[0].ID != "3" {
		t.Error("Incorrect data in L1")
	}

	if len(resultL2) != 1 || resultL2[0].ID != "3" {
		t.Error("Incorrect data in L2")
	}
}
