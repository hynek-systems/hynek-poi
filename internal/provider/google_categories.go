package provider

import "strings"

var googleCategoryMap = map[string]string{

	"restaurant": "restaurant",

	"cafe": "cafe",

	"bar": "bar",

	"hotel": "lodging",

	"atm": "atm",

	"hospital": "hospital",
}

func mapGoogleCategory(category string) string {

	if v, ok := googleCategoryMap[strings.ToLower(category)]; ok {

		return v
	}

	return ""
}
