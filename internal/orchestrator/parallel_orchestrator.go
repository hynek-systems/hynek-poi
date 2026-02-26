package orchestrator

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/dedupe"
	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/provider"
	"github.com/hynek-systems/hynek-poi/internal/ranking"
)

type ParallelOrchestrator struct {
	providers []provider.Provider
	timeout   time.Duration
}

func NewParallel(
	providers []provider.Provider,
	timeout time.Duration,
) *ParallelOrchestrator {

	return &ParallelOrchestrator{
		providers: providers,
		timeout:   timeout,
	}
}

func (o *ParallelOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	var wg sync.WaitGroup

	resultsChan := make(chan []domain.POI, len(o.providers))

	for _, p := range o.providers {

		wg.Add(1)

		go func(provider provider.Provider) {

			defer wg.Done()

			results, err := provider.Search(query)

			if err != nil {
				log.Printf("provider %s failed: %v", provider.Name(), err)
				return
			}

			if len(results) == 0 {
				log.Printf("provider %s returned 0 results", provider.Name())
				return
			}

			select {

			case resultsChan <- results:

			case <-ctx.Done():
				return
			}

		}(p)
	}

	// close channel when all providers finished
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var all []domain.POI

	for {

		select {

		case results, ok := <-resultsChan:

			if !ok {

				if len(all) == 0 {
					return nil, errors.New("all providers failed or timeout")
				}

				deduped := dedupe.Deduplicate(all)

				ranked := ranking.Rank(deduped, query)

				return ranked, nil
			}

			all = append(all, results...)

		case <-ctx.Done():

			if len(all) == 0 {
				return nil, errors.New("all providers failed or timeout")
			}

			deduped := dedupe.Deduplicate(all)

			ranked := ranking.Rank(deduped, query)

			return ranked, nil
		}
	}
}
