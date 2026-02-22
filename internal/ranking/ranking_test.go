package ranking

import (
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestRank_ByDistance(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{
		{ID: "1", Name: "Far", Latitude: 59.4000, Longitude: 18.1000, Source: "provider1"},
		{ID: "2", Name: "Near", Latitude: 59.3294, Longitude: 18.0687, Source: "provider1"},
		{ID: "3", Name: "Medium", Latitude: 59.3500, Longitude: 18.0800, Source: "provider1"},
	}

	result := Rank(pois, query)

	// Should be sorted by distance (closest first)
	if result[0].ID != "2" {
		t.Errorf("Expected nearest POI first, got %s", result[0].ID)
	}
	if result[2].ID != "1" {
		t.Errorf("Expected farthest POI last, got %s", result[2].ID)
	}
}

func TestRank_ByProviderPriority(t *testing.T) {
	// Set provider priorities
	SetProviderPriorities(map[string]int{
		"high-priority": 1,
		"low-priority":  10,
	})

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{
		{ID: "1", Name: "Low Priority", Latitude: 59.3293, Longitude: 18.0686, Source: "low-priority"},
		{ID: "2", Name: "High Priority", Latitude: 59.3293, Longitude: 18.0686, Source: "high-priority"},
	}

	result := Rank(pois, query)

	// High priority provider should come first
	if result[0].ID != "2" {
		t.Errorf("Expected high priority POI first, got %s", result[0].ID)
	}
	if result[1].ID != "1" {
		t.Errorf("Expected low priority POI second, got %s", result[1].ID)
	}

	// Reset priorities
	SetProviderPriorities(map[string]int{})
}

func TestRank_PriorityThenDistance(t *testing.T) {
	SetProviderPriorities(map[string]int{
		"provider-a": 1,
		"provider-b": 2,
	})

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{
		{ID: "1", Name: "B Far", Latitude: 59.4000, Longitude: 18.1000, Source: "provider-b"},
		{ID: "2", Name: "A Far", Latitude: 59.4000, Longitude: 18.1000, Source: "provider-a"},
		{ID: "3", Name: "B Near", Latitude: 59.3294, Longitude: 18.0687, Source: "provider-b"},
		{ID: "4", Name: "A Near", Latitude: 59.3294, Longitude: 18.0687, Source: "provider-a"},
	}

	result := Rank(pois, query)

	// Priority first, then distance
	// Expected order: A Near, A Far, B Near, B Far
	if result[0].ID != "4" {
		t.Errorf("Expected A Near first, got %s", result[0].ID)
	}
	if result[1].ID != "2" {
		t.Errorf("Expected A Far second, got %s", result[1].ID)
	}
	if result[2].ID != "3" {
		t.Errorf("Expected B Near third, got %s", result[2].ID)
	}
	if result[3].ID != "1" {
		t.Errorf("Expected B Far fourth, got %s", result[3].ID)
	}

	// Reset priorities
	SetProviderPriorities(map[string]int{})
}

func TestRank_UnknownProvider(t *testing.T) {
	SetProviderPriorities(map[string]int{
		"known": 1,
	})

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{
		{ID: "1", Name: "Unknown", Latitude: 59.3293, Longitude: 18.0686, Source: "unknown"},
		{ID: "2", Name: "Known", Latitude: 59.3293, Longitude: 18.0686, Source: "known"},
	}

	result := Rank(pois, query)

	// Known provider should come first (priority 1 vs default 100)
	if result[0].ID != "2" {
		t.Errorf("Expected known provider first, got %s", result[0].ID)
	}

	// Reset priorities
	SetProviderPriorities(map[string]int{})
}

func TestRank_EmptyList(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{}

	result := Rank(pois, query)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestRank_SingleItem(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	pois := []domain.POI{
		{ID: "1", Name: "Only", Latitude: 59.3293, Longitude: 18.0686, Source: "provider1"},
	}

	result := Rank(pois, query)

	if len(result) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("Expected same item, got %s", result[0].ID)
	}
}

func TestRank_StableSort(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	}

	// All POIs at same location and same priority
	pois := []domain.POI{
		{ID: "1", Name: "First", Latitude: 59.3293, Longitude: 18.0686, Source: "provider1"},
		{ID: "2", Name: "Second", Latitude: 59.3293, Longitude: 18.0686, Source: "provider1"},
		{ID: "3", Name: "Third", Latitude: 59.3293, Longitude: 18.0686, Source: "provider1"},
	}

	result := Rank(pois, query)

	// Order should be preserved (stable sort)
	if result[0].ID != "1" || result[1].ID != "2" || result[2].ID != "3" {
		t.Error("Expected stable sort to preserve order for equal elements")
	}
}

func TestDistance_Calculation(t *testing.T) {
	query := domain.SearchQuery{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	poi1 := domain.POI{Latitude: 1.0, Longitude: 0.0}
	poi2 := domain.POI{Latitude: 2.0, Longitude: 0.0}

	d1 := distance(query, poi1)
	d2 := distance(query, poi2)

	// poi2 should be farther than poi1
	if d2 <= d1 {
		t.Errorf("Expected poi2 to be farther, got d1=%f, d2=%f", d1, d2)
	}
}
