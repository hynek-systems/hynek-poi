package dedupe

import (
	"math"
	"strings"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

const distanceThresholdMeters = 50

func Deduplicate(pois []domain.POI) []domain.POI {

	var result []domain.POI

	for _, poi := range pois {

		if !exists(result, poi) {
			result = append(result, poi)
		}
	}

	return result
}

func exists(pois []domain.POI, candidate domain.POI) bool {

	for _, existing := range pois {

		if sameName(existing.Name, candidate.Name) &&
			sameLocation(existing, candidate) {

			return true
		}
	}

	return false
}

func sameName(a, b string) bool {

	return normalize(a) == normalize(b)
}

func normalize(s string) string {

	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	return s
}

func sameLocation(a, b domain.POI) bool {

	return distanceMeters(
		a.Latitude,
		a.Longitude,
		b.Latitude,
		b.Longitude,
	) < distanceThresholdMeters
}

// Haversine formula
func distanceMeters(lat1, lon1, lat2, lon2 float64) float64 {

	const R = 6371000

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*
			math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func toRad(d float64) float64 {

	return d * math.Pi / 180
}
