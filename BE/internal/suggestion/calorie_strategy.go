package suggestion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// CalorieStrategy implements the BY_CALORIES suggestion type.
// It finds dishes within a calorie range computed from a target +/- tolerance.
type CalorieStrategy struct {
	dishes   DishReader
	calories CalorieProvider
	sessions SessionRepository
}

// NewCalorieStrategy creates a new CalorieStrategy.
func NewCalorieStrategy(dishReader DishReader, calorieProvider CalorieProvider, sessionRepo SessionRepository) *CalorieStrategy {
	return &CalorieStrategy{
		dishes:   dishReader,
		calories: calorieProvider,
		sessions: sessionRepo,
	}
}

// Type returns the strategy identifier.
func (s *CalorieStrategy) Type() string {
	return "BY_CALORIES"
}

// Suggest finds a random dish within the target calorie range (default +/-15%).
// Optional filters can further narrow results.
func (s *CalorieStrategy) Suggest(ctx context.Context, params map[string]any) (*SuggestionResult, error) {
	targetCalories, tolerancePct, filter, userID := s.parseParams(params)
	if targetCalories <= 0 {
		return nil, fmt.Errorf("suggestion.CalorieStrategy.Suggest: target_calories is required and must be positive")
	}

	// Compute calorie range.
	minCal := float64(targetCalories) * (1.0 - tolerancePct)
	maxCal := float64(targetCalories) * (1.0 + tolerancePct)

	// Get recipe IDs within calorie range.
	calorieIDs, err := s.calories.GetRecipeIDsByCalorieRange(ctx, minCal, maxCal)
	if err != nil {
		return nil, fmt.Errorf("suggestion.CalorieStrategy.Suggest: %w", err)
	}
	if len(calorieIDs) == 0 {
		return nil, fmt.Errorf("suggestion.CalorieStrategy.Suggest: %w: no dishes in calorie range %.0f-%.0f", ErrNoCandidates, minCal, maxCal)
	}

	// Apply optional recipe filters.
	candidateIDs := calorieIDs
	if filter != nil {
		filteredIDs, err := s.dishes.GetIDsByFilter(ctx, *filter)
		if err != nil {
			slog.Warn("suggestion.CalorieStrategy.Suggest: filter failed, using calorie results only", "error", err)
		} else {
			candidateIDs = intersectIDs(calorieIDs, filteredIDs)
			if len(candidateIDs) == 0 {
				candidateIDs = calorieIDs
				slog.Info("suggestion.CalorieStrategy.Suggest: filters too restrictive, using calorie results only",
					"calorieMatches", len(calorieIDs))
			}
		}
	}

	// Pick a random candidate.
	pickedID, err := pickRandom(candidateIDs)
	if err != nil {
		return nil, fmt.Errorf("suggestion.CalorieStrategy.Suggest: %w", err)
	}

	// Fetch recipe detail.
	recipe, err := s.dishes.GetByID(ctx, pickedID)
	if err != nil {
		return nil, fmt.Errorf("suggestion.CalorieStrategy.Suggest: %w", err)
	}

	// Fetch actual calories.
	cal, err := s.calories.GetCalories(ctx, pickedID)
	if err != nil {
		slog.Warn("suggestion.CalorieStrategy.Suggest: failed to get calories", "error", err, "recipeID", pickedID)
	}
	recipe.Calories = cal

	// Calculate total calories for session.
	var totalCalories *int
	if cal != nil {
		tc := int(*cal)
		totalCalories = &tc
	}

	// Build and save session.
	session := SuggestionSession{
		ID:              uuid.New(),
		UserID:          userID,
		SessionType:     s.Type(),
		InputParams:     params,
		ResultRecipeIDs: []uuid.UUID{pickedID},
		TotalCalories:   totalCalories,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.sessions.Create(ctx, &session); err != nil {
		slog.Error("suggestion.CalorieStrategy.Suggest: failed to save session", "error", err)
	}

	return &SuggestionResult{
		Session: session,
		Recipes: []RecipeSummary{*recipe},
	}, nil
}

// parseParams extracts calorie strategy parameters from the params map.
func (s *CalorieStrategy) parseParams(params map[string]any) (int, float64, *RecipeFilter, *uuid.UUID) {
	targetCalories := 0
	tolerancePct := 0.15
	var filter *RecipeFilter
	var userID *uuid.UUID

	if tc, ok := params["target_calories"]; ok {
		switch v := tc.(type) {
		case float64:
			targetCalories = int(v)
		case int:
			targetCalories = v
		case json.Number:
			if n, err := v.Int64(); err == nil {
				targetCalories = int(n)
			}
		}
	}

	if tp, ok := params["tolerance_pct"]; ok {
		switch v := tp.(type) {
		case float64:
			if v > 0 && v <= 1 {
				tolerancePct = v
			}
		}
	}

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

	return targetCalories, tolerancePct, filter, userID
}

// intersectIDs returns the intersection of two UUID slices.
func intersectIDs(a, b []uuid.UUID) []uuid.UUID {
	set := make(map[uuid.UUID]struct{}, len(b))
	for _, id := range b {
		set[id] = struct{}{}
	}
	result := make([]uuid.UUID, 0)
	for _, id := range a {
		if _, ok := set[id]; ok {
			result = append(result, id)
		}
	}
	return result
}
