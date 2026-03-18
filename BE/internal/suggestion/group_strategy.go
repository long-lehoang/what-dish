package suggestion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
)

// Slot types for group meal composition.
const (
	slotMain = "main"
	slotSoup = "soup"
	slotSide = "side"
)

// GroupStrategy implements the BY_GROUP suggestion type.
// It builds a balanced meal combo for N people with diverse dishes.
type GroupStrategy struct {
	dishes   DishReader
	calories CalorieProvider
	sessions SessionRepository
}

// NewGroupStrategy creates a new GroupStrategy.
func NewGroupStrategy(dishReader DishReader, calorieProvider CalorieProvider, sessionRepo SessionRepository) *GroupStrategy {
	return &GroupStrategy{
		dishes:   dishReader,
		calories: calorieProvider,
		sessions: sessionRepo,
	}
}

// Type returns the strategy identifier.
func (s *GroupStrategy) Type() string {
	return "BY_GROUP"
}

// Suggest builds a balanced group meal: 1 main dish, 1 soup/canh, and
// additional side dishes. Dish count is determined by group size.
func (s *GroupStrategy) Suggest(ctx context.Context, params map[string]any) (*SuggestionResult, error) {
	groupSize, filter, userID := s.parseParams(params)
	if groupSize <= 0 {
		return nil, fmt.Errorf("suggestion.GroupStrategy.Suggest: group_size is required and must be positive")
	}

	// Determine dish count: min(ceil(groupSize/2)+1, 5), min=3.
	dishCount := int(math.Ceil(float64(groupSize)/2.0)) + 1
	if dishCount < 3 {
		dishCount = 3
	}
	if dishCount > 5 {
		dishCount = 5
	}

	// Allocate slots: 1 main, 1 soup, rest sides.
	slots := make([]string, dishCount)
	slots[0] = slotMain
	slots[1] = slotSoup
	for i := 2; i < dishCount; i++ {
		slots[i] = slotSide
	}

	// Get candidate IDs (either filtered or all published).
	candidateIDs, err := s.getCandidateIDs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("suggestion.GroupStrategy.Suggest: %w", err)
	}
	if len(candidateIDs) < dishCount {
		return nil, fmt.Errorf("suggestion.GroupStrategy.Suggest: %w: need %d dishes but only %d available",
			ErrNoCandidates, dishCount, len(candidateIDs))
	}

	// Pick dishes ensuring no duplicates.
	pickedIDs, err := s.pickDishes(candidateIDs, slots)
	if err != nil {
		return nil, fmt.Errorf("suggestion.GroupStrategy.Suggest: %w", err)
	}

	// Fetch recipe details.
	recipes, err := s.dishes.GetByIDs(ctx, pickedIDs)
	if err != nil {
		return nil, fmt.Errorf("suggestion.GroupStrategy.Suggest: %w", err)
	}

	// Calculate total nutrition.
	totalCalories := 0
	hasCalories := false
	for i := range recipes {
		cal, err := s.calories.GetCalories(ctx, recipes[i].ID)
		if err != nil {
			slog.Warn("suggestion.GroupStrategy.Suggest: failed to get calories",
				"error", err, "recipeID", recipes[i].ID)
			continue
		}
		if cal != nil {
			recipes[i].Calories = cal
			totalCalories += int(*cal)
			hasCalories = true
		}
	}

	var totalCalPtr *int
	if hasCalories {
		totalCalPtr = &totalCalories
	}

	// Build and save session.
	session := SuggestionSession{
		ID:              uuid.New(),
		UserID:          userID,
		SessionType:     s.Type(),
		InputParams:     params,
		ResultRecipeIDs: pickedIDs,
		TotalCalories:   totalCalPtr,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.sessions.Create(ctx, &session); err != nil {
		slog.Error("suggestion.GroupStrategy.Suggest: failed to save session", "error", err)
	}

	return &SuggestionResult{
		Session: session,
		Recipes: recipes,
	}, nil
}

// getCandidateIDs returns filtered or all published recipe IDs.
func (s *GroupStrategy) getCandidateIDs(ctx context.Context, filter *RecipeFilter) ([]uuid.UUID, error) {
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

// pickDishes selects unique dishes for each slot, ensuring no duplicate recipes.
func (s *GroupStrategy) pickDishes(candidateIDs []uuid.UUID, slots []string) ([]uuid.UUID, error) {
	used := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0, len(slots))
	remaining := make([]uuid.UUID, len(candidateIDs))
	copy(remaining, candidateIDs)

	for range slots {
		available := make([]uuid.UUID, 0, len(remaining))
		for _, id := range remaining {
			if _, ok := used[id]; !ok {
				available = append(available, id)
			}
		}
		if len(available) == 0 {
			return nil, fmt.Errorf("pickDishes: not enough unique dishes")
		}

		picked, err := pickRandom(available)
		if err != nil {
			return nil, err
		}
		result = append(result, picked)
		used[picked] = struct{}{}
	}

	return result, nil
}

// parseParams extracts group strategy parameters from the params map.
func (s *GroupStrategy) parseParams(params map[string]any) (int, *RecipeFilter, *uuid.UUID) {
	groupSize := 0
	var filter *RecipeFilter
	var userID *uuid.UUID

	if gs, ok := params["group_size"]; ok {
		switch v := gs.(type) {
		case float64:
			groupSize = int(v)
		case int:
			groupSize = v
		case json.Number:
			if n, err := v.Int64(); err == nil {
				groupSize = int(n)
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

	return groupSize, filter, userID
}
