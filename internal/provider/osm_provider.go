package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type OSMProvider struct {
	endpoint string
	client   *http.Client
}

func NewOSMProvider() *OSMProvider {
	return &OSMProvider{
		endpoint: "https://overpass-api.de/api/interpreter",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *OSMProvider) Name() string {
	return "osm"
}

type overpassResponse struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	ID   int64             `json:"id"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags"`
}

func (p *OSMProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	amenityFilter := ""

	if len(query.Categories) > 0 {

		var mapped []string

		for _, cat := range query.Categories {

			if amenity, ok := mapCategory(cat); ok {
				mapped = append(mapped, amenity)
			}
		}

		if len(mapped) > 0 {

			regex := strings.Join(mapped, "|")

			amenityFilter = fmt.Sprintf(`["amenity"~"%s"]`, regex)

		} else {

			amenityFilter = `["amenity"]`
		}

	} else {

		amenityFilter = `["amenity"]`
	}

	var overpassQuery string

	if query.BBox != nil {

		overpassQuery = fmt.Sprintf(
			`[out:json][timeout:5];node%s(%f,%f,%f,%f);out body %d;`,
			amenityFilter,
			query.BBox.MinLat,
			query.BBox.MinLng,
			query.BBox.MaxLat,
			query.BBox.MaxLng,
			query.Limit,
		)

	} else {

		overpassQuery = fmt.Sprintf(
			`[out:json][timeout:5];node%s(around:%d,%f,%f);out body %d;`,
			amenityFilter,
			query.Radius,
			query.Latitude,
			query.Longitude,
			query.Limit,
		)
	}

	form := url.Values{}
	form.Add("data", overpassQuery)

	resp, err := p.client.PostForm(p.endpoint, form)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var overpassResp overpassResponse

	err = json.NewDecoder(resp.Body).Decode(&overpassResp)

	if err != nil {
		return nil, err
	}

	var pois []domain.POI

	for _, element := range overpassResp.Elements {

		name := element.Tags["name"]
		if name == "" {
			continue
		}

		category := element.Tags["amenity"]

		poi := domain.POI{
			ID:        fmt.Sprintf("%d", element.ID),
			Name:      name,
			Latitude:  element.Lat,
			Longitude: element.Lon,
			Category:  category,
			Source:    p.Name(),
			Website:   element.Tags["website"],
			Phone:     element.Tags["phone"],
			Cuisine:   element.Tags["cuisine"],
			Email:     element.Tags["email"],
			Address:   buildOSMAddress(element.Tags),
		}

		if hours := element.Tags["opening_hours"]; hours != "" {
			poi.OpeningHours = []string{hours}
		}

		if v, ok := element.Tags["wheelchair"]; ok {
			b := v == "yes"
			poi.WheelchairAccessible = &b
		}

		if v, ok := element.Tags["outdoor_seating"]; ok {
			b := v == "yes"
			poi.OutdoorSeating = &b
		}

		if v, ok := element.Tags["takeaway"]; ok {
			b := v == "yes"
			poi.Takeaway = &b
		}

		if v, ok := element.Tags["delivery"]; ok {
			b := v == "yes"
			poi.Delivery = &b
		}

		pois = append(pois, poi)
	}

	return pois, nil
}

func buildOSMAddress(tags map[string]string) string {

	var parts []string

	if v := tags["addr:street"]; v != "" {
		street := v
		if num := tags["addr:housenumber"]; num != "" {
			street += " " + num
		}
		parts = append(parts, street)
	}

	if v := tags["addr:postcode"]; v != "" {
		parts = append(parts, v)
	}

	if v := tags["addr:city"]; v != "" {
		parts = append(parts, v)
	}

	return strings.Join(parts, ", ")
}
