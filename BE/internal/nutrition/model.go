package nutrition

import (
	"time"

	"github.com/google/uuid"
)

// RecipeNutrition holds the per-serving nutritional data for a recipe.
type RecipeNutrition struct {
	ID          uuid.UUID `json:"id"`
	RecipeID    uuid.UUID `json:"recipeId"`
	Calories    *float64  `json:"calories,omitempty"`
	Protein     *float64  `json:"protein,omitempty"`
	Carbs       *float64  `json:"carbs,omitempty"`
	Fat         *float64  `json:"fat,omitempty"`
	Fiber       *float64  `json:"fiber,omitempty"`
	Sugar       *float64  `json:"sugar,omitempty"`
	Sodium      *float64  `json:"sodium,omitempty"`
	ServingSize *string   `json:"servingSize,omitempty"`
	DataSource  *string   `json:"dataSource,omitempty"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NutritionGoal represents a preset nutrition goal (e.g. weight loss, muscle gain).
type NutritionGoal struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	MealCaloriesMin  *int      `json:"mealCaloriesMin,omitempty"`
	MealCaloriesMax  *int      `json:"mealCaloriesMax,omitempty"`
	DailyCaloriesMin *int      `json:"dailyCaloriesMin,omitempty"`
	DailyCaloriesMax *int      `json:"dailyCaloriesMax,omitempty"`
	ProteinPct       *int      `json:"proteinPct,omitempty"`
	CarbsPct         *int      `json:"carbsPct,omitempty"`
	FatPct           *int      `json:"fatPct,omitempty"`
	IsActive         bool      `json:"isActive"`
	SortOrder        int       `json:"sortOrder"`
}

// TDEEResult holds the calculated TDEE values and meal breakdown.
type TDEEResult struct {
	BMR           float64            `json:"bmr"`
	TDEE          float64            `json:"tdee"`
	DailyTarget   float64            `json:"dailyTarget"`
	MealBreakdown map[string]float64 `json:"mealBreakdown"`
}
