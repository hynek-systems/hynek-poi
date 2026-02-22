package dedupe

import (
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestDeduplicate_NoDuplicates(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Restaurant B", Latitude: 59.3400, Longitude: 18.0700},
		{ID: "3", Name: "Cafe C", Latitude: 59.3500, Longitude: 18.0800},
	}

	result := Deduplicate(pois)

	if len(result) != 3 {
		t.Errorf("Expected 3 POIs, got %d", len(result))
	}
}

func TestDeduplicate_ExactDuplicates(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686}, // Exact duplicate
		{ID: "3", Name: "Cafe C", Latitude: 59.3500, Longitude: 18.0800},
	}

	result := Deduplicate(pois)

	if len(result) != 2 {
		t.Errorf("Expected 2 POIs after deduplication, got %d", len(result))
	}

	// First occurrence should be kept
	if result[0].ID != "1" {
		t.Errorf("Expected first POI to be kept, got ID %s", result[0].ID)
	}
}

func TestDeduplicate_SameNameDifferentLocation(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Starbucks", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Starbucks", Latitude: 59.4000, Longitude: 18.1000}, // Same name, far location
	}

	result := Deduplicate(pois)

	// Should keep both - different locations
	if len(result) != 2 {
		t.Errorf("Expected 2 POIs (different locations), got %d", len(result))
	}
}

func TestDeduplicate_CloseProximity(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Restaurant A", Latitude: 59.32935, Longitude: 18.06865}, // Very close, ~5m away
	}

	result := Deduplicate(pois)

	if len(result) != 1 {
		t.Errorf("Expected 1 POI (duplicates within threshold), got %d", len(result))
	}
}

func TestDeduplicate_CaseInsensitive(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "RESTAURANT A", Latitude: 59.3293, Longitude: 18.0686}, // Same name, different case
	}

	result := Deduplicate(pois)

	if len(result) != 1 {
		t.Errorf("Expected 1 POI (case insensitive), got %d", len(result))
	}
}

func TestDeduplicate_WhitespaceNormalization(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "  Restaurant A  ", Latitude: 59.3293, Longitude: 18.0686}, // Extra whitespace
	}

	result := Deduplicate(pois)

	if len(result) != 1 {
		t.Errorf("Expected 1 POI (whitespace normalized), got %d", len(result))
	}
}

func TestDeduplicate_EmptyInput(t *testing.T) {
	pois := []domain.POI{}

	result := Deduplicate(pois)

	if len(result) != 0 {
		t.Errorf("Expected 0 POIs, got %d", len(result))
	}
}

func TestDeduplicate_MultipleOccurrences(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "3", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "4", Name: "Cafe B", Latitude: 59.3400, Longitude: 18.0700},
		{ID: "5", Name: "Cafe B", Latitude: 59.3400, Longitude: 18.0700},
	}

	result := Deduplicate(pois)

	if len(result) != 2 {
		t.Errorf("Expected 2 unique POIs, got %d", len(result))
	}
}

func TestDeduplicate_JustOutsideThreshold(t *testing.T) {
	pois := []domain.POI{
		{ID: "1", Name: "Restaurant A", Latitude: 59.3293, Longitude: 18.0686},
		{ID: "2", Name: "Restaurant A", Latitude: 59.3298, Longitude: 18.0686}, // ~55m away (just outside 50m threshold)
	}

	result := Deduplicate(pois)

	// Should keep both - just outside threshold
	if len(result) != 2 {
		t.Errorf("Expected 2 POIs (just outside threshold), got %d", len(result))
	}
}

func TestDistanceMeters(t *testing.T) {
	// Stockholm coordinates
	lat1, lon1 := 59.3293, 18.0686
	lat2, lon2 := 59.3293, 18.0686

	// Same location
	dist := distanceMeters(lat1, lon1, lat2, lon2)
	if dist != 0 {
		t.Errorf("Expected 0 meters for same location, got %f", dist)
	}

	// Moving east at latitude ~59 degrees
	dist = distanceMeters(59.3293, 18.0686, 59.3293, 18.0820)
	// At latitude 59°, 1 degree longitude ≈ 56km, so 0.0134° ≈ 750m
	if dist < 700 || dist > 800 {
		t.Errorf("Expected ~750 meters, got %f", dist)
	}

	// Moving north (1 degree latitude is always ~111km)
	dist = distanceMeters(59.3293, 18.0686, 59.3393, 18.0686)
	// 0.01 degrees latitude ≈ 1111m
	if dist < 1100 || dist > 1120 {
		t.Errorf("Expected ~1111 meters, got %f", dist)
	}
}
