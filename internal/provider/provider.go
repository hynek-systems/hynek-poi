package provider

import "github.com/hynek-systems/hynek-poi/internal/domain"

type Provider interface {
	Name() string

	Search(query domain.SearchQuery) ([]domain.POI, error)
}
