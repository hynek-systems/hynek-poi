package provider

var osmCategoryMap = map[string]string{

	"restaurant": "restaurant",
	"cafe":       "cafe",
	"bar":        "bar",
	"pub":        "pub",
	"fast_food":  "fast_food",

	"atm":      "atm",
	"bank":     "bank",
	"hospital": "hospital",
	"pharmacy": "pharmacy",

	"hotel": "hotel",

	"fuel": "fuel",

	"parking": "parking",
}

func mapCategory(category string) (string, bool) {

	value, ok := osmCategoryMap[category]

	return value, ok
}
