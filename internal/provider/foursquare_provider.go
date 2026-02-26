package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type FoursquareProvider struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

func NewFoursquareProvider(apiKey string) *FoursquareProvider {

	return &FoursquareProvider{
		apiKey:   apiKey,
		endpoint: "https://api.foursquare.com/v3/places/search",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *FoursquareProvider) Name() string {

	return "foursquare"
}

type foursquareResponse struct {
	Results []foursquarePlace `json:"results"`
}

type foursquarePlace struct {
	FsqID       string               `json:"fsq_id"`
	Name        string               `json:"name"`
	Categories  []foursquareCategory `json:"categories"`
	Geocodes    foursquareGeocodes   `json:"geocodes"`
	Rating      float64              `json:"rating"`
	Price       int                  `json:"price"`
	Tel         string               `json:"tel"`
	Website     string               `json:"website"`
	Menu        string               `json:"menu"`
	Hours       *foursquareHours     `json:"hours"`
	Tastes      []string             `json:"tastes"`
	Location    *foursquareLocation  `json:"location"`
	Description string               `json:"description"`
	Email       string               `json:"email"`
	Verified    *bool                `json:"verified"`
	Popularity  float64              `json:"popularity"`
}

type foursquareLocation struct {
	Address          string `json:"address"`
	Locality         string `json:"locality"`
	Region           string `json:"region"`
	Postcode         string `json:"postcode"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formatted_address"`
}

type foursquareHours struct {
	Display string                  `json:"display"`
	OpenNow bool                    `json:"open_now"`
	Regular []foursquareHoursEntry  `json:"regular"`
}

type foursquareHoursEntry struct {
	Day   int    `json:"day"`
	Open  string `json:"open"`
	Close string `json:"close"`
}

type foursquareCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type foursquareGeocodes struct {
	Main foursquareLatLng `json:"main"`
}

type foursquareLatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (p *FoursquareProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	req, err := http.NewRequest("GET", p.endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", p.apiKey)
	req.Header.Set("Accept", "application/json")

	params := req.URL.Query()

	params.Set("ll", fmt.Sprintf("%f,%f", query.Latitude, query.Longitude))

	if query.Radius > 0 {
		params.Set("radius", fmt.Sprintf("%d", query.Radius))
	}

	if query.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}

	if len(query.Categories) > 0 {

		var ids []string

		for _, cat := range query.Categories {

			if id, ok := mapFoursquareCategory(cat); ok {
				ids = append(ids, fmt.Sprintf("%d", id))
			}
		}

		if len(ids) > 0 {
			params.Set("categories", strings.Join(ids, ","))
		}
	}

	params.Set("fields", "fsq_id,name,categories,geocodes,rating,price,tel,website,hours,menu,tastes,location,description,email,verified,popularity")

	req.URL.RawQuery = params.Encode()

	resp, err := p.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		return nil, fmt.Errorf("foursquare status %d", resp.StatusCode)
	}

	var fsqResp foursquareResponse

	err = json.NewDecoder(resp.Body).Decode(&fsqResp)

	if err != nil {
		return nil, err
	}

	var pois []domain.POI

	for _, place := range fsqResp.Results {

		category := ""

		if len(place.Categories) > 0 {
			category = place.Categories[0].Name
		}

		poi := domain.POI{
			ID:          place.FsqID,
			Name:        place.Name,
			Latitude:    place.Geocodes.Main.Latitude,
			Longitude:   place.Geocodes.Main.Longitude,
			Category:    category,
			Source:      p.Name(),
			Rating:      place.Rating,
			PriceLevel:  place.Price,
			Phone:       place.Tel,
			Website:     place.Website,
			MenuURL:     place.Menu,
			Description: place.Description,
			Email:       place.Email,
			Verified:    place.Verified,
			Popularity:  place.Popularity,
		}

		if place.Location != nil && place.Location.FormattedAddress != "" {
			poi.Address = place.Location.FormattedAddress
		}

		if place.Hours != nil {
			if place.Hours.Display != "" {
				poi.OpeningHours = []string{place.Hours.Display}
			}
			openNow := place.Hours.OpenNow
			poi.OpenNow = &openNow
		}

		if len(place.Tastes) > 0 {
			poi.Cuisine = strings.Join(place.Tastes, ", ")
		}

		pois = append(pois, poi)
	}

	return pois, nil
}
