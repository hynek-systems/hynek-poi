package domain

type POI struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  string  `json:"category"`
	Source    string  `json:"source"`

	Rating       float64  `json:"rating,omitempty"`
	RatingCount  int      `json:"rating_count,omitempty"`
	Website      string   `json:"website,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	OpeningHours []string `json:"opening_hours,omitempty"`
	Cuisine      string   `json:"cuisine,omitempty"`
	PriceLevel   int      `json:"price_level,omitempty"`
	MenuURL      string   `json:"menu_url,omitempty"`
}
