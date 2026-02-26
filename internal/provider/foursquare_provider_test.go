package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

func TestFoursquareProvider_Name(t *testing.T) {

	p := NewFoursquareProvider("test-key")

	if p.Name() != "foursquare" {
		t.Errorf("Expected name 'foursquare', got '%s'", p.Name())
	}
}

func TestFoursquareProvider_Search(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "test-key" {
			t.Errorf("Expected Authorization header 'test-key', got '%s'", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept header 'application/json', got '%s'", r.Header.Get("Accept"))
		}

		resp := foursquareResponse{
			Results: []foursquarePlace{
				{
					FsqID: "abc123",
					Name:  "Test Restaurant",
					Categories: []foursquareCategory{
						{ID: 13065, Name: "Restaurant"},
					},
					Geocodes: foursquareGeocodes{
						Main: foursquareLatLng{
							Latitude:  59.3293,
							Longitude: 18.0686,
						},
					},
					Rating:  8.5,
					Price:   2,
					Tel:     "+46812345678",
					Website: "https://testrestaurant.se",
					Menu:    "https://testrestaurant.se/menu",
					Hours: &foursquareHours{
						Display: "Mon-Fri 11:00-22:00",
						OpenNow: true,
					},
					Tastes: []string{"Swedish", "Seafood"},
					Location: &foursquareLocation{
						FormattedAddress: "Storgatan 1, 111 23 Stockholm",
					},
					Description: "A cozy Swedish restaurant",
					Email:       "info@testrestaurant.se",
					Verified:    boolPtr(true),
					Popularity:  0.85,
				},
				{
					FsqID: "def456",
					Name:  "Test Cafe",
					Categories: []foursquareCategory{
						{ID: 13032, Name: "Cafe"},
					},
					Geocodes: foursquareGeocodes{
						Main: foursquareLatLng{
							Latitude:  59.3300,
							Longitude: 18.0700,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))

	defer server.Close()

	p := NewFoursquareProvider("test-key")
	p.endpoint = server.URL

	query := domain.SearchQuery{
		Latitude:   59.3293,
		Longitude:  18.0686,
		Radius:     1000,
		Limit:      50,
		Categories: []string{"restaurant"},
	}

	results, err := p.Search(query)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].ID != "abc123" {
		t.Errorf("Expected ID 'abc123', got '%s'", results[0].ID)
	}

	if results[0].Name != "Test Restaurant" {
		t.Errorf("Expected name 'Test Restaurant', got '%s'", results[0].Name)
	}

	if results[0].Category != "Restaurant" {
		t.Errorf("Expected category 'Restaurant', got '%s'", results[0].Category)
	}

	if results[0].Source != "foursquare" {
		t.Errorf("Expected source 'foursquare', got '%s'", results[0].Source)
	}

	if results[0].Latitude != 59.3293 {
		t.Errorf("Expected latitude 59.3293, got %f", results[0].Latitude)
	}

	if results[0].Rating != 8.5 {
		t.Errorf("Expected rating 8.5, got %f", results[0].Rating)
	}

	if results[0].PriceLevel != 2 {
		t.Errorf("Expected price level 2, got %d", results[0].PriceLevel)
	}

	if results[0].Phone != "+46812345678" {
		t.Errorf("Expected phone '+46812345678', got '%s'", results[0].Phone)
	}

	if results[0].Website != "https://testrestaurant.se" {
		t.Errorf("Expected website 'https://testrestaurant.se', got '%s'", results[0].Website)
	}

	if results[0].MenuURL != "https://testrestaurant.se/menu" {
		t.Errorf("Expected menu URL 'https://testrestaurant.se/menu', got '%s'", results[0].MenuURL)
	}

	if len(results[0].OpeningHours) != 1 || results[0].OpeningHours[0] != "Mon-Fri 11:00-22:00" {
		t.Errorf("Expected opening hours ['Mon-Fri 11:00-22:00'], got %v", results[0].OpeningHours)
	}

	if results[0].Cuisine != "Swedish, Seafood" {
		t.Errorf("Expected cuisine 'Swedish, Seafood', got '%s'", results[0].Cuisine)
	}

	if results[0].Address != "Storgatan 1, 111 23 Stockholm" {
		t.Errorf("Expected address 'Storgatan 1, 111 23 Stockholm', got '%s'", results[0].Address)
	}

	if results[0].Description != "A cozy Swedish restaurant" {
		t.Errorf("Expected description 'A cozy Swedish restaurant', got '%s'", results[0].Description)
	}

	if results[0].Email != "info@testrestaurant.se" {
		t.Errorf("Expected email 'info@testrestaurant.se', got '%s'", results[0].Email)
	}

	if results[0].OpenNow == nil || !*results[0].OpenNow {
		t.Errorf("Expected open_now true, got %v", results[0].OpenNow)
	}

	if results[0].Verified == nil || !*results[0].Verified {
		t.Errorf("Expected verified true, got %v", results[0].Verified)
	}

	if results[0].Popularity != 0.85 {
		t.Errorf("Expected popularity 0.85, got %f", results[0].Popularity)
	}

	if results[1].ID != "def456" {
		t.Errorf("Expected ID 'def456', got '%s'", results[1].ID)
	}

	// Second result has no enriched fields
	if results[1].Rating != 0 {
		t.Errorf("Expected rating 0, got %f", results[1].Rating)
	}

	if results[1].Website != "" {
		t.Errorf("Expected empty website, got '%s'", results[1].Website)
	}
}

func TestFoursquareProvider_SearchEmptyResults(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := foursquareResponse{
			Results: []foursquarePlace{},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))

	defer server.Close()

	p := NewFoursquareProvider("test-key")
	p.endpoint = server.URL

	results, err := p.Search(domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
		Radius:    1000,
		Limit:     50,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestFoursquareProvider_SearchNonOKStatus(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))

	defer server.Close()

	p := NewFoursquareProvider("bad-key")
	p.endpoint = server.URL

	_, err := p.Search(domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := "foursquare status 401"

	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestFoursquareProvider_SearchNoCategories(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := foursquareResponse{
			Results: []foursquarePlace{
				{
					FsqID:      "xyz789",
					Name:       "Unknown Place",
					Categories: []foursquareCategory{},
					Geocodes: foursquareGeocodes{
						Main: foursquareLatLng{
							Latitude:  59.3293,
							Longitude: 18.0686,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))

	defer server.Close()

	p := NewFoursquareProvider("test-key")
	p.endpoint = server.URL

	results, err := p.Search(domain.SearchQuery{
		Latitude:  59.3293,
		Longitude: 18.0686,
		Radius:    1000,
		Limit:     50,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Category != "" {
		t.Errorf("Expected empty category, got '%s'", results[0].Category)
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestMapFoursquareCategory(t *testing.T) {

	tests := []struct {
		input    string
		expected int
		found    bool
	}{
		{"restaurant", 13065, true},
		{"cafe", 13032, true},
		{"bar", 13003, true},
		{"hotel", 19014, true},
		{"Restaurant", 13065, true},
		{"unknown", 0, false},
	}

	for _, tt := range tests {

		id, ok := mapFoursquareCategory(tt.input)

		if ok != tt.found {
			t.Errorf("mapFoursquareCategory(%q) found = %v, want %v", tt.input, ok, tt.found)
		}

		if id != tt.expected {
			t.Errorf("mapFoursquareCategory(%q) = %d, want %d", tt.input, id, tt.expected)
		}
	}
}
