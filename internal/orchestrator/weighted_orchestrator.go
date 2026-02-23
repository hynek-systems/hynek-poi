package orchestrator

import (
	"errors"
	"math/rand"
	"sort"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type WeightedOrchestrator struct {
	providers []ProviderConfig
}

func NewWeighted(configs []ProviderConfig) *WeightedOrchestrator {

	// sort highest weight first
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Weight > configs[j].Weight
	})

	return &WeightedOrchestrator{
		providers: configs,
	}
}

func (o *WeightedOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {

	// shuffle providers with weight bias
	providers := o.weightedShuffle()

	for _, config := range providers {

		results, err := config.Provider.Search(query)

		if err != nil {
			continue
		}

		if len(results) > 0 {
			return results, nil
		}
	}

	return nil, errors.New("no provider returned results")
}

func (o *WeightedOrchestrator) weightedShuffle() []ProviderConfig {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var shuffled []ProviderConfig

	for _, config := range o.providers {

		for i := 0; i < config.Weight; i++ {
			shuffled = append(shuffled, config)
		}
	}

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// remove duplicates, keep order
	seen := make(map[string]bool)

	var result []ProviderConfig

	for _, config := range shuffled {

		name := config.Provider.Name()

		if !seen[name] {
			result = append(result, config)
			seen[name] = true
		}
	}

	return result
}
