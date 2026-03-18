package suggestion

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SuggestionService orchestrates suggestion operations by delegating to
// the appropriate strategy based on session type.
type SuggestionService struct {
	strategies map[string]Strategy
	sessions   SessionRepository
	configs    ConfigRepository
}

// NewSuggestionService creates a new SuggestionService.
func NewSuggestionService(strategies map[string]Strategy, sessionRepo SessionRepository, configRepo ConfigRepository) *SuggestionService {
	return &SuggestionService{
		strategies: strategies,
		sessions:   sessionRepo,
		configs:    configRepo,
	}
}

// Suggest delegates to the appropriate strategy for the given session type.
func (s *SuggestionService) Suggest(ctx context.Context, sessionType string, params map[string]any) (*SuggestionResult, error) {
	strategy, ok := s.strategies[sessionType]
	if !ok {
		return nil, fmt.Errorf("suggestion.SuggestionService.Suggest: unknown session type %q", sessionType)
	}
	result, err := strategy.Suggest(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("suggestion.SuggestionService.Suggest: %w", err)
	}
	return result, nil
}

// GetHistory returns paginated suggestion history for a user.
func (s *SuggestionService) GetHistory(ctx context.Context, userID uuid.UUID, sessionType string, page, pageSize int) ([]SuggestionSession, int64, error) {
	sessions, total, err := s.sessions.ListByUser(ctx, userID, sessionType, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("suggestion.SuggestionService.GetHistory: %w", err)
	}
	return sessions, total, nil
}

// ListConfigs returns suggestion configurations filtered by type.
func (s *SuggestionService) ListConfigs(ctx context.Context, configType string) ([]SuggestionConfig, error) {
	configs, err := s.configs.ListByType(ctx, configType)
	if err != nil {
		return nil, fmt.Errorf("suggestion.SuggestionService.ListConfigs: %w", err)
	}
	return configs, nil
}
