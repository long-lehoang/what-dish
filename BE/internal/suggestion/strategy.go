package suggestion

import "context"

// Strategy defines the interface for suggestion algorithms.
type Strategy interface {
	Suggest(ctx context.Context, params map[string]any) (*SuggestionResult, error)
	Type() string
}
