package provider

import (
	"strings"
	"testing"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
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

func TestTimeoutProvider_ReturnsResultBeforeTimeout(t *testing.T) {

	base := &mockProvider{
		name: "fast",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return []domain.POI{{ID: "1", Name: "Result"}}, nil
		},
	}

	tp := NewTimeoutProvider(base, 1*time.Second)

	results, err := tp.Search(domain.SearchQuery{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 || results[0].ID != "1" {
		t.Fatalf("Expected 1 result with ID '1', got %v", results)
	}
}

func TestTimeoutProvider_TimesOut(t *testing.T) {

	base := &mockProvider{
		name: "slow",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			time.Sleep(500 * time.Millisecond)
			return []domain.POI{{ID: "1", Name: "Late"}}, nil
		},
	}

	tp := NewTimeoutProvider(base, 50*time.Millisecond)

	results, err := tp.Search(domain.SearchQuery{})

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "provider timeout") {
		t.Errorf("Expected 'provider timeout' error, got: %v", err)
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

func TestTimeoutProvider_PropagatesError(t *testing.T) {

	base := &mockProvider{
		name: "failing",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			return nil, errMock("provider error")
		},
	}

	tp := NewTimeoutProvider(base, 1*time.Second)

	results, err := tp.Search(domain.SearchQuery{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "provider error" {
		t.Errorf("Expected 'provider error', got: %v", err)
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

func TestTimeoutProvider_Name(t *testing.T) {

	base := &mockProvider{name: "test-provider"}

	tp := NewTimeoutProvider(base, 1*time.Second)

	if tp.Name() != "test-provider" {
		t.Errorf("Expected name 'test-provider', got '%s'", tp.Name())
	}
}

func TestTimeoutProvider_ZeroTimeoutTimesOutImmediately(t *testing.T) {

	base := &mockProvider{
		name: "any",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			time.Sleep(10 * time.Millisecond)
			return []domain.POI{{ID: "1", Name: "Result"}}, nil
		},
	}

	tp := NewTimeoutProvider(base, 0)

	_, err := tp.Search(domain.SearchQuery{})

	if err == nil {
		t.Fatal("Expected timeout error with zero duration, got nil")
	}

	if !strings.Contains(err.Error(), "provider timeout") {
		t.Errorf("Expected 'provider timeout' error, got: %v", err)
	}
}

type errMock string

func (e errMock) Error() string {
	return string(e)
}
