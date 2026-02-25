package provider

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestRetryProvider_SucceedsFirstAttempt(t *testing.T) {

	var calls int32

	base := &mockProvider{
		name: "ok",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			atomic.AddInt32(&calls, 1)
			return []domain.POI{{ID: "1", Name: "Result"}}, nil
		},
	}

	rp := NewRetryProvider(base, 3)

	results, err := rp.Search(domain.SearchQuery{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 || results[0].ID != "1" {
		t.Fatalf("Expected 1 result with ID '1', got %v", results)
	}

	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}

func TestRetryProvider_SucceedsAfterRetries(t *testing.T) {

	var calls int32

	base := &mockProvider{
		name: "flaky",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			n := atomic.AddInt32(&calls, 1)

			if n < 3 {
				return nil, errors.New("temporary failure")
			}

			return []domain.POI{{ID: "1", Name: "Result"}}, nil
		},
	}

	rp := NewRetryProvider(base, 3)

	results, err := rp.Search(domain.SearchQuery{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}

func TestRetryProvider_ExhaustsRetries(t *testing.T) {

	var calls int32

	base := &mockProvider{
		name: "broken",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			atomic.AddInt32(&calls, 1)
			return nil, errors.New("permanent failure")
		},
	}

	rp := NewRetryProvider(base, 2)

	results, err := rp.Search(domain.SearchQuery{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "permanent failure" {
		t.Errorf("Expected 'permanent failure', got: %v", err)
	}

	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}

	// 1 initial + 2 retries = 3 total calls
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("Expected 3 calls (1 + 2 retries), got %d", calls)
	}
}

func TestRetryProvider_ZeroRetries(t *testing.T) {

	var calls int32

	base := &mockProvider{
		name: "once",
		searchFunc: func(q domain.SearchQuery) ([]domain.POI, error) {
			atomic.AddInt32(&calls, 1)
			return nil, errors.New("failed")
		},
	}

	rp := NewRetryProvider(base, 0)

	_, err := rp.Search(domain.SearchQuery{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// 0 retries = only 1 attempt
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("Expected 1 call with 0 retries, got %d", calls)
	}
}

func TestRetryProvider_Name(t *testing.T) {

	base := &mockProvider{name: "test-provider"}

	rp := NewRetryProvider(base, 1)

	if rp.Name() != "test-provider" {
		t.Errorf("Expected name 'test-provider', got '%s'", rp.Name())
	}
}
