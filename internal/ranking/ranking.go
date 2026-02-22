package ranking

import (
	"sort"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

var providerPriority = map[string]int{}

func SetProviderPriorities(priorities map[string]int) {

	providerPriority = priorities
}

func Rank(pois []domain.POI, query domain.SearchQuery) []domain.POI {

	sort.SliceStable(pois, func(i, j int) bool {

		a := pois[i]
		b := pois[j]

		// 1. provider priority
		pa := priority(a.Source)
		pb := priority(b.Source)

		if pa != pb {
			return pa < pb
		}

		// 2. distance
		da := distance(query, a)
		db := distance(query, b)

		return da < db
	})

	return pois
}

func priority(provider string) int {

	if p, ok := providerPriority[provider]; ok {

		return p
	}

	return 100
}

func distance(query domain.SearchQuery, poi domain.POI) float64 {

	dx := query.Latitude - poi.Latitude
	dy := query.Longitude - poi.Longitude

	return dx*dx + dy*dy
}
