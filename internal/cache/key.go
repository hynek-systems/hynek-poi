package cache

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/mmcloughlin/geohash"
)

const precision = 6

func BuildKey(query domain.SearchQuery) string {

	hash := geohash.EncodeWithPrecision(
		query.Latitude,
		query.Longitude,
		precision,
	)

	categoryPart := normalizeCategories(query.Categories)

	return fmt.Sprintf(
		"poi:%s:%d:%s",
		hash,
		query.Radius,
		categoryPart,
	)
}

func normalizeCategories(categories []string) string {

	if len(categories) == 0 {
		return "all"
	}

	normalized := make([]string, 0, len(categories))

	for _, c := range categories {

		c = strings.ToLower(strings.TrimSpace(c))

		if c != "" {
			normalized = append(normalized, c)
		}
	}

	sort.Strings(normalized)

	return strings.Join(normalized, ",")
}
