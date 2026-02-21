package orchestrator

import "github.com/hynek-systems/hynek-poi/internal/provider"

type ProviderConfig struct {
	Provider provider.Provider
	Weight   int
}
