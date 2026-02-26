package domain

type SearchQuery struct {
	Latitude  float64
	Longitude float64

	// optional
	BBox *BBox

	Radius int
	Limit  int

	Categories []string
}

type BBox struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
}

type PaginatedResponse struct {
	Data       []POI `json:"data"`
	Total      int   `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}
