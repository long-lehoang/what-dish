package nutrition

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// NutritionService provides business logic for recipe nutrition and TDEE calculation.
type NutritionService struct {
	nutritionRepo NutritionRepository
	goalRepo      GoalRepository
}

// NewNutritionService creates a new NutritionService.
func NewNutritionService(nutritionRepo NutritionRepository, goalRepo GoalRepository) *NutritionService {
	return &NutritionService{
		nutritionRepo: nutritionRepo,
		goalRepo:      goalRepo,
	}
}

// GetByRecipeID returns nutrition data for a single recipe.
func (s *NutritionService) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) (*RecipeNutritionResponse, error) {
	n, err := s.nutritionRepo.GetByRecipeID(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("NutritionService.GetByRecipeID: %w", err)
	}

	return toNutritionResponse(n), nil
}

// GetByRecipeIDs returns nutrition data for multiple recipes.
func (s *NutritionService) GetByRecipeIDs(ctx context.Context, recipeIDs []uuid.UUID) ([]RecipeNutritionResponse, error) {
	items, err := s.nutritionRepo.GetByRecipeIDs(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("NutritionService.GetByRecipeIDs: %w", err)
	}

	results := make([]RecipeNutritionResponse, len(items))
	for i, n := range items {
		results[i] = *toNutritionResponse(&n)
	}

	return results, nil
}

// CalculateTDEE performs a stateless TDEE calculation.
func (s *NutritionService) CalculateTDEE(_ context.Context, req CalculateTDEERequest) (*TDEEResponse, error) {
	bmr := CalculateBMR(req.Gender, req.WeightKg, req.HeightCm, req.Age)
	tdee := CalculateTDEE(bmr, req.ActivityLevel)
	dailyTarget := AdjustForGoal(tdee, req.Goal)
	mealBreakdown := CalculateMealBreakdown(dailyTarget)

	slog.Debug("tdee calculated",
		"gender", req.Gender,
		"bmr", bmr,
		"tdee", tdee,
		"daily_target", dailyTarget,
	)

	return &TDEEResponse{
		BMR:           bmr,
		TDEE:          tdee,
		DailyTarget:   dailyTarget,
		MealBreakdown: mealBreakdown,
	}, nil
}

// Upsert creates or updates nutrition data for a recipe (admin endpoint).
func (s *NutritionService) Upsert(ctx context.Context, recipeID uuid.UUID, req UpsertNutritionRequest) (*RecipeNutritionResponse, error) {
	now := time.Now().UTC()

	n := &RecipeNutrition{
		ID:          uuid.New(),
		RecipeID:    recipeID,
		Calories:    req.Calories,
		Protein:     req.Protein,
		Carbs:       req.Carbs,
		Fat:         req.Fat,
		Fiber:       req.Fiber,
		Sugar:       req.Sugar,
		Sodium:      req.Sodium,
		ServingSize: req.ServingSize,
		DataSource:  req.DataSource,
		IsVerified:  req.IsVerified,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.nutritionRepo.Upsert(ctx, n); err != nil {
		return nil, fmt.Errorf("NutritionService.Upsert: %w", err)
	}

	slog.Info("recipe nutrition upserted", "recipe_id", recipeID)

	return toNutritionResponse(n), nil
}

// ListGoals returns all active nutrition goals.
func (s *NutritionService) ListGoals(ctx context.Context) ([]NutritionGoal, error) {
	goals, err := s.goalRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("NutritionService.ListGoals: %w", err)
	}

	return goals, nil
}

// GetIDsByCalorieRange returns recipe IDs within a calorie range.
func (s *NutritionService) GetIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error) {
	ids, err := s.nutritionRepo.GetIDsByCalorieRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("NutritionService.GetIDsByCalorieRange: %w", err)
	}

	return ids, nil
}

func toNutritionResponse(n *RecipeNutrition) *RecipeNutritionResponse {
	return &RecipeNutritionResponse{
		ID:          n.ID.String(),
		RecipeID:    n.RecipeID.String(),
		Calories:    n.Calories,
		Protein:     n.Protein,
		Carbs:       n.Carbs,
		Fat:         n.Fat,
		Fiber:       n.Fiber,
		Sugar:       n.Sugar,
		Sodium:      n.Sodium,
		ServingSize: n.ServingSize,
		DataSource:  n.DataSource,
		IsVerified:  n.IsVerified,
	}
}
