package nutrition

import (
	"context"

	"github.com/google/uuid"
)

// NutritionRepository manages persistence of recipe nutrition data.
type NutritionRepository interface {
	GetByRecipeID(ctx context.Context, recipeID uuid.UUID) (*RecipeNutrition, error)
	GetByRecipeIDs(ctx context.Context, recipeIDs []uuid.UUID) ([]RecipeNutrition, error)
	Upsert(ctx context.Context, nutrition *RecipeNutrition) error
	GetIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error)
}

// GoalRepository manages persistence of nutrition goals.
type GoalRepository interface {
	List(ctx context.Context) ([]NutritionGoal, error)
}
