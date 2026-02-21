package orchestrator

import (
	"errors"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/provider"
)

type FallbackOrchestrator struct {
	providers []provider.Provider
}

var _ Orchestrator = (*FallbackOrchestrator)(nil)

func NewFallback(providers []provider.Provider) *FallbackOrchestrator {
	return &FallbackOrchestrator{
		providers: providers,
	}
}

func (o *FallbackOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {

	for _, provider := range o.providers {

		results, err := provider.Search(query)

		if err != nil {
			continue
		}

		if len(results) > 0 {
			return results, nil
		}
	}

	return nil, errors.New("no provider returned results")
}
