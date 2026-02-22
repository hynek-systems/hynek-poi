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

		pois = append(pois, domain.POI{
			ID:        fmt.Sprintf("%d", element.ID),
			Name:      name,
			Latitude:  element.Lat,
			Longitude: element.Lon,
			Category:  category,
			Source:    p.Name(),
		})
	}

	return pois, nil
}
