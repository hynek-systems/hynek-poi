package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	overpassQuery := fmt.Sprintf(`
[out:json][timeout:10];
node["amenity"](around:%d,%f,%f);
out;
`,
		query.Radius,
		query.Latitude,
		query.Longitude,
	)

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
