package orchestrator

import (
	"errors"
	"testing"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/cache"
	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type mockOrchestrator struct {
	searchFunc func(domain.SearchQuery) ([]domain.POI, error)
	callCount  int
}

func (m *mockOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {
	m.callCount++
	return m.searchFunc(query)
}

func TestCachedOrchestrator_CacheHit(t *testing.T) {
	memCache := cache.NewMemoryCache()
	mockInner := &mockOrchestrator{
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Test"}}, nil
		},
	}

	orchestrator := NewCached(mockInner, memCache, 1*time.Minute)

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
		Radius:    1000,
	}

	// First call - should hit inner orchestrator
	results, err := orchestrator.Search(query)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if mockInner.callCount != 1 {
		t.Errorf("Expected 1 call to inner orchestrator, got %d", mockInner.callCount)
	}

	// Second call - should hit cache
	results, err = orchestrator.Search(query)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if mockInner.callCount != 1 {
		t.Errorf("Expected still 1 call to inner orchestrator (cached), got %d", mockInner.callCount)
	}
}

func TestCachedOrchestrator_CacheMiss(t *testing.T) {
	memCache := cache.NewMemoryCache()
	mockInner := &mockOrchestrator{
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Test"}}, nil
		},
	}

	orchestrator := NewCached(mockInner, memCache, 1*time.Minute)

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
		Radius:    1000,
	}

	// Should miss cache and call inner orchestrator
	results, err := orchestrator.Search(query)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if mockInner.callCount != 1 {
		t.Errorf("Expected 1 call to inner orchestrator, got %d", mockInner.callCount)
	}
}

func TestCachedOrchestrator_ErrorNotCached(t *testing.T) {
	memCache := cache.NewMemoryCache()
	expectedErr := errors.New("provider error")
	mockInner := &mockOrchestrator{
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, expectedErr
		},
	}

	orchestrator := NewCached(mockInner, memCache, 1*time.Minute)

	query := domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
		Radius:    1000,
	}

	// First call - should return error
	_, err := orchestrator.Search(query)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	// Second call - error should not be cached, should call inner again
	_, err = orchestrator.Search(query)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if mockInner.callCount != 2 {
		t.Errorf("Expected 2 calls to inner orchestrator (errors not cached), got %d", mockInner.callCount)
	}
}
