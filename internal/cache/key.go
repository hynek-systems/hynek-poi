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

	categoryPart := normalizeCategories(query.Categories)

	var bboxPart string

	if query.BBox != nil {

		bboxPart = fmt.Sprintf(
			"%f:%f:%f:%f",
			query.BBox.MinLat,
			query.BBox.MinLng,
			query.BBox.MaxLat,
			query.BBox.MaxLng,
		)

	} else {

		hash := geohash.EncodeWithPrecision(
			query.Latitude,
			query.Longitude,
			precision,
		)

		bboxPart = hash
	}

	return fmt.Sprintf(
		"poi:%s:%d:%s",
		bboxPart,
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
