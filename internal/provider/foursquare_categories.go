package provider

import "strings"

// Foursquare Places API v3 category IDs.
// Reference: https://docs.foursquare.com/data-products/docs/categories
var foursquareCategoryMap = map[string]int{

	"restaurant": 13065,

	"cafe": 13032,

	"bar": 13003,

	"pub": 13025,

	"fast_food": 13145,

	"hotel": 19014,

	"atm": 11044,

	"bank": 11045,

	"hospital": 15014,

	"pharmacy": 15026,

	"fuel": 19007,

	"parking": 19020,
}

func mapFoursquareCategory(category string) (int, bool) {

	id, ok := foursquareCategoryMap[strings.ToLower(category)]

	return id, ok
}
