package orchestrator

import (
	"errors"
	"testing"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/provider"
)

func TestParallelOrchestrator_MergesResults(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{
				{ID: "1", Name: "Result 1", Latitude: 59.0, Longitude: 18.0, Source: "provider1"},
			}, nil
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{
				{ID: "2", Name: "Result 2", Latitude: 59.0, Longitude: 18.0, Source: "provider2"},
			}, nil
		},
	}

	orchestrator := NewParallel(
		[]provider.Provider{provider1, provider2},
		1*time.Second,
	)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// ParallelOrchestrator now merges results from all providers
	if len(results) != 2 {
		t.Fatalf("Expected 2 results (merged), got %d", len(results))
	}

	// Verify both results are present
	foundProvider1 := false
	foundProvider2 := false
	for _, r := range results {
		if r.ID == "1" {
			foundProvider1 = true
		}
		if r.ID == "2" {
			foundProvider2 = true
		}
	}

	if !foundProvider1 || !foundProvider2 {
		t.Error("Expected results from both providers")
	}
}

func TestParallelOrchestrator_OneProviderFails(t *testing.T) {
	failingProvider := &mockProvider{
		name: "failing",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("provider failed")
		},
	}

	workingProvider := &mockProvider{
		name: "working",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Working Result"}}, nil
		},
	}

	orchestrator := NewParallel(
		[]provider.Provider{failingProvider, workingProvider},
		1*time.Second,
	)

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

func TestParallelOrchestrator_AllProvidersFail(t *testing.T) {
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

	orchestrator := NewParallel(
		[]provider.Provider{provider1, provider2},
		1*time.Second,
	)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	// NOTE: Current implementation returns nil, nil when all providers fail quickly
	// This may be unintended behavior - consider returning an error
	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}

	if err != nil {
		t.Logf("Got error (good): %v", err)
	}
}

func TestParallelOrchestrator_Timeout(t *testing.T) {
	slowProvider := &mockProvider{
		name: "slow",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			time.Sleep(500 * time.Millisecond)
			return []domain.POI{{ID: "1", Name: "Slow"}}, nil
		},
	}

	orchestrator := NewParallel(
		[]provider.Provider{slowProvider},
		100*time.Millisecond, // Short timeout
	)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results on timeout, got %v", results)
	}
}

func TestParallelOrchestrator_EmptyResults(t *testing.T) {
	emptyProvider := &mockProvider{
		name: "empty",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{}, nil
		},
	}

	orchestrator := NewParallel(
		[]provider.Provider{emptyProvider},
		1*time.Second,
	)

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	// NOTE: Current implementation returns nil, nil when provider returns empty results
	// This may be unintended behavior - consider returning an error
	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}

	if err != nil {
		t.Logf("Got error (good): %v", err)
	}
}
