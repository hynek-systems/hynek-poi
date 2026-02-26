package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type GoogleProvider struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

func NewGoogleProvider(apiKey string) *GoogleProvider {

	return &GoogleProvider{
		apiKey:   apiKey,
		endpoint: "https://maps.googleapis.com/maps/api/place/nearbysearch/json",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *GoogleProvider) Name() string {

	return "google"
}

type googleResponse struct {
	Results []googleResult `json:"results"`
}

type googleResult struct {
	PlaceID          string   `json:"place_id"`
	Name             string   `json:"name"`
	Types            []string `json:"types"`
	Rating           float64  `json:"rating"`
	UserRatingsTotal int      `json:"user_ratings_total"`
	PriceLevel       int      `json:"price_level"`

	Geometry struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`

	OpeningHours *googleOpeningHours `json:"opening_hours"`
}

type googleOpeningHours struct {
	WeekdayText []string `json:"weekday_text"`
}

func (p *GoogleProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	params := url.Values{}

	params.Set("key", p.apiKey)

	params.Set(
		"location",
		fmt.Sprintf("%f,%f", query.Latitude, query.Longitude),
	)

	params.Set(
		"radius",
		fmt.Sprintf("%d", query.Radius),
	)

	if len(query.Categories) > 0 {

		params.Set(
			"type",
			mapGoogleCategory(query.Categories[0]),
		)
	}

	reqURL := p.endpoint + "?" + params.Encode()

	resp, err := p.client.Get(reqURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		return nil, fmt.Errorf("google status %d", resp.StatusCode)
	}

	var gr googleResponse

	err = json.NewDecoder(resp.Body).Decode(&gr)

	if err != nil {
		return nil, err
	}

	var pois []domain.POI

	for _, r := range gr.Results {

		category := ""

		if len(r.Types) > 0 {
			category = r.Types[0]
		}

		poi := domain.POI{
			ID:          r.PlaceID,
			Name:        r.Name,
			Latitude:    r.Geometry.Location.Lat,
			Longitude:   r.Geometry.Location.Lng,
			Category:    category,
			Source:      p.Name(),
			Rating:      r.Rating,
			RatingCount: r.UserRatingsTotal,
			PriceLevel:  r.PriceLevel,
		}

		if r.OpeningHours != nil && len(r.OpeningHours.WeekdayText) > 0 {
			poi.OpeningHours = r.OpeningHours.WeekdayText
		}

		pois = append(pois, poi)
	}

	return pois, nil
}
