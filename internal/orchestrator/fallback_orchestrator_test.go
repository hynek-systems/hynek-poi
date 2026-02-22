package orchestrator

import (
	"errors"
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/provider"
)

type mockProvider struct {
	name       string
	searchFunc func(domain.SearchQuery) ([]domain.POI, error)
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {
	return m.searchFunc(query)
}

func TestFallbackOrchestrator_FirstProviderSucceeds(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "From Provider 1"}}, nil
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "2", Name: "From Provider 2"}}, nil
		},
	}

	orchestrator := NewFallback([]provider.Provider{provider1, provider2})

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "1" {
		t.Errorf("Expected result from provider1, got %s", results[0].ID)
	}
}

func TestFallbackOrchestrator_FirstProviderFails(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("provider1 failed")
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "2", Name: "From Provider 2"}}, nil
		},
	}

	orchestrator := NewFallback([]provider.Provider{provider1, provider2})

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "2" {
		t.Errorf("Expected result from provider2, got %s", results[0].ID)
	}
}

func TestFallbackOrchestrator_FirstProviderReturnsEmpty(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{}, nil // Empty results
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "2", Name: "From Provider 2"}}, nil
		},
	}

	orchestrator := NewFallback([]provider.Provider{provider1, provider2})

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "2" {
		t.Errorf("Expected result from provider2, got %s", results[0].ID)
	}
}

func TestFallbackOrchestrator_AllProvidersFail(t *testing.T) {
	provider1 := &mockProvider{
		name: "provider1",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("provider1 failed")
		},
	}

	provider2 := &mockProvider{
		name: "provider2",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errors.New("provider2 failed")
		},
	}

	orchestrator := NewFallback([]provider.Provider{provider1, provider2})

	query := domain.SearchQuery{Latitude: 59.0, Longitude: 18.0}
	results, err := orchestrator.Search(query)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}

	if err.Error() != "no provider returned results" {
		t.Errorf("Expected 'no provider returned results', got %s", err.Error())
	}
}
