package domain

type SearchQuery struct {
	Latitude   float64
	Longitude  float64
	Radius     int
	Limit      int
	Categories []string
}
