package orchestrator

import (
	"errors"
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestWeightedOrchestrator_Success(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Result 1"}}, nil
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "2", Name: "Result 2"}}, nil
		},
	}

	configs := []ProviderConfig{
		{Provider: provider1, Weight: 10},
		{Provider: provider2, Weight: 5},
	}

	orchestrator := NewWeighted(configs)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected results, got none")
	}

	// Should get result from one of the providers
	if results[0].ID != "1" && results[0].ID != "2" {
		t.Errorf("Expected result from provider1 or provider2, got %s", results[0].ID)
	}
}

func TestWeightedOrchestrator_HigherWeightFirst(t *testing.T) {
	lowWeightProvider := &mockProvider{
		name: "low",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "low", Name: "Low Weight"}}, nil
		},
	}

	highWeightProvider := &mockProvider{
		name: "high",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "high", Name: "High Weight"}}, nil
		},
	}

	configs := []ProviderConfig{
		{Provider: lowWeightProvider, Weight: 1},
		{Provider: highWeightProvider, Weight: 100},
	}

	orchestrator := NewWeighted(configs)

	// Run multiple times to check weight bias
	highWeightCount := 0
	totalRuns := 20

	for i := 0; i < totalRuns; i++ {
		query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
		results, err := orchestrator.Search(query)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(results) > 0 && results[0].ID == "high" {
			highWeightCount++
		}
	}

	// High weight provider should be selected more often
	// With weight 100 vs 1, we expect high weight to win most of the time
	if highWeightCount < totalRuns/2 {
		t.Errorf("Expected high weight provider to win more often, got %d/%d", highWeightCount, totalRuns)
	}
}

func TestWeightedOrchestrator_FallbackOnError(t *testing.T) {
	failingProvider := &mockProvider{
		name: "failing",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("failed")
		},
	}

	workingProvider := &mockProvider{
		name: "working",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Working"}}, nil
		},
	}

	configs := []ProviderConfig{
		{Provider: failingProvider, Weight: 10},
		{Provider: workingProvider, Weight: 5},
	}

	orchestrator := NewWeighted(configs)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "1" {
		t.Errorf("Expected result from working provider, got %s", results[0].ID)
	}
}

func TestWeightedOrchestrator_AllProvidersFail(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("failed")
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("failed")
		},
	}

	configs := []ProviderConfig{
		{Provider: provider1, Weight: 10},
		{Provider: provider2, Weight: 5},
	}

	orchestrator := NewWeighted(configs)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

func TestWeightedOrchestrator_EmptyResults(t *testing.T) {
	emptyProvider := &mockProvider{
		name: "empty",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{}, nil
		},
	}

	workingProvider := &mockProvider{
		name: "working",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Working"}}, nil
		},
	}

	configs := []ProviderConfig{
		{Provider: emptyProvider, Weight: 10},
		{Provider: workingProvider, Weight: 5},
	}

	orchestrator := NewWeighted(configs)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should fallback to working provider
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
}
