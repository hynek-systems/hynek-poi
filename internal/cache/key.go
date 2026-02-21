package cache

import (
	"fmt"
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

	categoryPart := strings.Join(query.Categories, ",")

	return fmt.Sprintf(
		"poi:%s:%d:%s",
		hash,
		query.Radius,
		categoryPart,
	)
}
