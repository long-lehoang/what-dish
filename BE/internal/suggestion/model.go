package suggestion

import (
	"time"

	"github.com/google/uuid"
)

// SuggestionSession represents a single suggestion request and its result.
type SuggestionSession struct {
	ID              uuid.UUID      `json:"id"`
	UserID          *uuid.UUID     `json:"userId,omitempty"`
	SessionType     string         `json:"sessionType"`
	InputParams     map[string]any `json:"inputParams"`
	ResultRecipeIDs []uuid.UUID    `json:"resultRecipeIds"`
	TotalCalories   *int           `json:"totalCalories,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
}

// SuggestionConfig holds preset configuration for suggestion modes.
type SuggestionConfig struct {
	ID          uuid.UUID      `json:"id"`
	ConfigType  string         `json:"configType"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Params      map[string]any `json:"params"`
	IsActive    bool           `json:"isActive"`
	SortOrder   int            `json:"sortOrder"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// ExclusionRule tracks recently suggested dishes to avoid repeats.
type ExclusionRule struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	RecipeID      uuid.UUID `json:"recipeId"`
	ExcludedUntil time.Time `json:"excludedUntil"`
	CreatedAt     time.Time `json:"createdAt"`
}

// SuggestionResult bundles a session with its resolved recipes.
type SuggestionResult struct {
	Session SuggestionSession `json:"session"`
	Recipes []RecipeSummary   `json:"recipes"`
}

// RecipeSummary is a cross-context read model for recipe data needed by suggestions.
type RecipeSummary struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	ImageURL   string    `json:"imageUrl"`
	Difficulty string    `json:"difficulty"`
	CookTime   int       `json:"cookTime"`
	Calories   *float64  `json:"calories,omitempty"`
}
