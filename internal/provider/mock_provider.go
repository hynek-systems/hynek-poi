package provider

import "github.com/hynek-systems/hynek-poi/internal/domain"

type MockProvider struct{}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	pois := []domain.POI{
		{
			ID:           "1",
			Name:         "Test Restaurant",
			Latitude:     query.Latitude,
			Longitude:    query.Longitude,
			Category:     "restaurant",
			Source:       p.Name(),
			Rating:       4.5,
			RatingCount:  120,
			Website:      "https://example.com",
			Phone:        "+46701234567",
			OpeningHours: []string{"Mon-Sun 10:00-22:00"},
			Cuisine:      "Swedish",
			PriceLevel:   2,
		},
	}

	return pois, nil
}
