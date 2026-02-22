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
