package suggestion

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// RandomStrategy implements the RANDOM suggestion type.
// It picks a random published dish, filtering out recently suggested ones.
type RandomStrategy struct {
	dishes    DishReader
	exclusion ExclusionRepository
	sessions  SessionRepository
}

// NewRandomStrategy creates a new RandomStrategy.
func NewRandomStrategy(dishReader DishReader, exclusionRepo ExclusionRepository, sessionRepo SessionRepository) *RandomStrategy {
	return &RandomStrategy{
		dishes:    dishReader,
		exclusion: exclusionRepo,
		sessions:  sessionRepo,
	}
}

// Type returns the strategy identifier.
func (s *RandomStrategy) Type() string {
	return "RANDOM"
}

// Suggest picks a random dish from all published recipes, excluding recently
// suggested ones. If no candidates remain after exclusion, exclusions are
// ignored so a result is always returned.
func (s *RandomStrategy) Suggest(ctx context.Context, params map[string]any) (*SuggestionResult, error) {
	// Parse optional filters from params.
	filter, userID := s.parseParams(params)

	// Get candidate IDs.
	candidateIDs, err := s.getCandidateIDs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("suggestion.RandomStrategy.Suggest: %w", err)
	}
	if len(candidateIDs) == 0 {
		return nil, fmt.Errorf("suggestion.RandomStrategy.Suggest: %w: no dishes available", ErrNoCandidates)
	}

	// Remove excluded IDs for authenticated users.
	if userID != nil {
		excludedIDs, err := s.exclusion.GetExcludedRecipeIDs(ctx, *userID)
		if err != nil {
			slog.Warn("suggestion.RandomStrategy.Suggest: failed to get exclusions, proceeding without",
				"error", err, "userID", userID)
		} else {
			filtered := removeIDs(candidateIDs, excludedIDs)
			if len(filtered) > 0 {
				candidateIDs = filtered
			} else {
				slog.Info("suggestion.RandomStrategy.Suggest: all candidates excluded, ignoring exclusions",
					"userID", userID, "total", len(candidateIDs))
			}
		}
	}

	// Pick a random candidate.
	pickedID, err := pickRandom(candidateIDs)
	if err != nil {
		return nil, fmt.Errorf("suggestion.RandomStrategy.Suggest: %w", err)
	}

	// Fetch full recipe detail.
	recipe, err := s.dishes.GetByID(ctx, pickedID)
	if err != nil {
		return nil, fmt.Errorf("suggestion.RandomStrategy.Suggest: %w", err)
	}

	// Build and save session.
	session := SuggestionSession{
		ID:              uuid.New(),
		UserID:          userID,
		SessionType:     s.Type(),
		InputParams:     params,
		ResultRecipeIDs: []uuid.UUID{pickedID},
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.sessions.Create(ctx, &session); err != nil {
		slog.Error("suggestion.RandomStrategy.Suggest: failed to save session", "error", err)
	}

	// Add exclusion rule for authenticated users.
	if userID != nil {
		rule := ExclusionRule{
			ID:            uuid.New(),
			UserID:        *userID,
			RecipeID:      pickedID,
			ExcludedUntil: time.Now().UTC().Add(7 * 24 * time.Hour),
			CreatedAt:     time.Now().UTC(),
		}
		if err := s.exclusion.Add(ctx, &rule); err != nil {
			slog.Warn("suggestion.RandomStrategy.Suggest: failed to add exclusion", "error", err)
		}
	}

	return &SuggestionResult{
		Session: session,
		Recipes: []RecipeSummary{*recipe},
	}, nil
}

// getCandidateIDs returns either filtered IDs or all published IDs.
func (s *RandomStrategy) getCandidateIDs(ctx context.Context, filter *RecipeFilter) ([]uuid.UUID, error) {
	if filter != nil {
		ids, err := s.dishes.GetIDsByFilter(ctx, *filter)
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			return ids, nil
		}
	}
	return s.dishes.GetAllPublishedIDs(ctx)
}

// parseParams extracts filter and userID from the params map.
func (s *RandomStrategy) parseParams(params map[string]any) (*RecipeFilter, *uuid.UUID) {
	var filter *RecipeFilter
	var userID *uuid.UUID

	if f, ok := params["filters"]; ok {
		filter = parseRecipeFilter(f)
	}
	if uid, ok := params["userId"]; ok {
		if uidStr, ok := uid.(string); ok {
			if parsed, err := uuid.Parse(uidStr); err == nil {
				userID = &parsed
			}
		}
	}
	return filter, userID
}

// pickRandom selects a random element from a slice using crypto/rand.
func pickRandom(ids []uuid.UUID) (uuid.UUID, error) {
	if len(ids) == 0 {
		return uuid.UUID{}, fmt.Errorf("pickRandom: empty slice")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(ids))))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pickRandom: %w", err)
	}
	return ids[n.Int64()], nil
}

// removeIDs returns a new slice with excluded IDs removed.
func removeIDs(candidates, excluded []uuid.UUID) []uuid.UUID {
	excludeSet := make(map[uuid.UUID]struct{}, len(excluded))
	for _, id := range excluded {
		excludeSet[id] = struct{}{}
	}
	result := make([]uuid.UUID, 0, len(candidates))
	for _, id := range candidates {
		if _, ok := excludeSet[id]; !ok {
			result = append(result, id)
		}
	}
	return result
}

// parseRecipeFilter attempts to convert an interface{} to a RecipeFilter.
func parseRecipeFilter(v any) *RecipeFilter {
	switch f := v.(type) {
	case RecipeFilter:
		return &f
	case *RecipeFilter:
		return f
	case map[string]any:
		data, err := json.Marshal(f)
		if err != nil {
			return nil
		}
		var filter RecipeFilter
		if err := json.Unmarshal(data, &filter); err != nil {
			return nil
		}
		return &filter
	default:
		return nil
	}
}

// ErrNoCandidates is returned when no recipes match the criteria.
var ErrNoCandidates = fmt.Errorf("no candidates available")
