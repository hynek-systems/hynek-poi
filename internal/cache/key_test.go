package cache

import (
	"strings"
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestBuildKey_Deterministic(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant", "cafe"},
	}

	// Build key twice
	key1 := BuildKey(query)
	key2 := BuildKey(query)

	if key1 != key2 {
		t.Errorf("Expected deterministic keys, got %s and %s", key1, key2)
	}
}

func TestBuildKey_Format(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
	}

	key := BuildKey(query)

	// Check format: poi:<geohash>:<radius>:<categories>
	if !strings.HasPrefix(key, "poi:") {
		t.Errorf("Expected key to start with 'poi:', got %s", key)
	}

	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		t.Errorf("Expected 4 parts in key, got %d: %s", len(parts), key)
	}

	if parts[2] != "1000" {
		t.Errorf("Expected radius 1000 in key, got %s", parts[2])
	}

	if parts[3] != "restaurant" {
		t.Errorf("Expected category restaurant in key, got %s", parts[3])
	}
}

func TestBuildKey_DifferentQueries(t *testing.T) {
	query1 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
	}

	query2 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     2000, // Different radius
		Categories: []string{"restaurant"},
	}

	query3 := domain.SearchQuery{
		Latitude:   60.0000, // Different location
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
	}

	key1 := BuildKey(query1)
	key2 := BuildKey(query2)
	key3 := BuildKey(query3)

	if key1 == key2 {
		t.Error("Expected different keys for different radius")
	}

	if key1 == key3 {
		t.Error("Expected different keys for different location")
	}
}

func TestBuildKey_MultipleCategories(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant", "cafe", "bar"},
	}

	key := BuildKey(query)

	// Categories are now sorted alphabetically
	if !strings.Contains(key, "bar,cafe,restaurant") {
		t.Errorf("Expected sorted categories in key, got %s", key)
	}
}

func TestBuildKey_NoCategories(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{},
	}

	key := BuildKey(query)

	// Should still build a valid key with "all" for empty categories
	if !strings.HasPrefix(key, "poi:") {
		t.Errorf("Expected valid key even with no categories, got %s", key)
	}
	if !strings.Contains(key, "all") {
		t.Errorf("Expected 'all' for empty categories, got %s", key)
	}
}

func TestBuildKey_WithBBox(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
		BBox: &domain.BBox{
			MinLat: 59.3000,
			MinLng: 18.0000,
			MaxLat: 59.4000,
			MaxLng: 18.1000,
		},
	}

	key := BuildKey(query)

	// Should use BBox instead of geohash
	if !strings.HasPrefix(key, "poi:") {
		t.Errorf("Expected valid key, got %s", key)
	}

	// Should contain BBox coordinates
	if !strings.Contains(key, "59.3") || !strings.Contains(key, "18.0") {
		t.Errorf("Expected BBox coordinates in key, got %s", key)
	}
}

func TestBuildKey_BBoxVsGeohash(t *testing.T) {
	query1 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
	}

	query2 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant"},
		BBox: &domain.BBox{
			MinLat: 59.3000,
			MinLng: 18.0000,
			MaxLat: 59.4000,
			MaxLng: 18.1000,
		},
	}

	key1 := BuildKey(query1)
	key2 := BuildKey(query2)

	// Keys should be different - one with geohash, one with BBox
	if key1 == key2 {
		t.Error("Expected different keys for geohash vs BBox")
	}
}

func TestBuildKey_CategoryNormalization(t *testing.T) {
	query1 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"Restaurant", "CAFE"},
	}

	query2 := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Categories: []string{"restaurant", "cafe"},
	}

	key1 := BuildKey(query1)
	key2 := BuildKey(query2)

	// Keys should be identical - categories normalized to lowercase
	if key1 != key2 {
		t.Errorf("Expected identical keys for normalized categories, got %s and %s", key1, key2)
	}
}
