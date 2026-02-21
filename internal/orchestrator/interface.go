package orchestrator

import "github.com/hynek-systems/hynek-poi/internal/domain"

type Orchestrator interface {
	Search(query domain.SearchQuery) ([]domain.POI, error)
}
