package suggestion

import (
	"context"

	"github.com/google/uuid"
)

// SessionRepository persists suggestion sessions.
type SessionRepository interface {
	Create(ctx context.Context, session *SuggestionSession) error
	ListByUser(ctx context.Context, userID uuid.UUID, sessionType string, page, pageSize int) ([]SuggestionSession, int64, error)
}

// ConfigRepository reads suggestion configuration presets.
type ConfigRepository interface {
	ListByType(ctx context.Context, configType string) ([]SuggestionConfig, error)
}

// ExclusionRepository manages per-user dish exclusion rules.
type ExclusionRepository interface {
	GetExcludedRecipeIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	Add(ctx context.Context, rule *ExclusionRule) error
	CleanExpired(ctx context.Context) error
}

// DishReader provides read access to recipe data from the recipe bounded context.
type DishReader interface {
	GetAllPublishedIDs(ctx context.Context) ([]uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*RecipeSummary, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]RecipeSummary, error)
	GetIDsByFilter(ctx context.Context, filter RecipeFilter) ([]uuid.UUID, error)
}

// CalorieProvider provides calorie data from the nutrition bounded context.
type CalorieProvider interface {
	GetRecipeIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error)
	GetCalories(ctx context.Context, recipeID uuid.UUID) (*float64, error)
}

// RecipeFilter holds optional filter criteria for dish queries.
type RecipeFilter struct {
	DishTypeID       *uuid.UUID `json:"dishTypeId,omitempty"`
	RegionID         *uuid.UUID `json:"regionId,omitempty"`
	MainIngredientID *uuid.UUID `json:"mainIngredientId,omitempty"`
	MealTypeID       *uuid.UUID `json:"mealTypeId,omitempty"`
	Difficulty       *string    `json:"difficulty,omitempty"`
	MaxCookTime      *int       `json:"maxCookTime,omitempty"`
}
