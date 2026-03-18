package suggestion

import "github.com/google/uuid"

// RandomSuggestionRequest is the request body for POST /suggestions/random.
type RandomSuggestionRequest struct {
	Filters *RecipeFilter `json:"filters,omitempty"`
}

// CalorieSuggestionRequest is the request body for POST /suggestions/by-calories.
type CalorieSuggestionRequest struct {
	TargetCalories int           `json:"targetCalories" binding:"required,gt=0"`
	MealType       string        `json:"mealType"`
	TolerancePct   *float64      `json:"tolerancePct,omitempty"`
	Filters        *RecipeFilter `json:"filters,omitempty"`
}

// GroupSuggestionRequest is the request body for POST /suggestions/by-group.
type GroupSuggestionRequest struct {
	GroupSize int           `json:"groupSize" binding:"required,gt=0"`
	GroupType string        `json:"groupType"`
	MealType  string        `json:"mealType"`
	Filters   *RecipeFilter `json:"filters,omitempty"`
}

// SuggestionResponse is the unified response for all suggestion endpoints.
type SuggestionResponse struct {
	SessionID     uuid.UUID       `json:"sessionId"`
	Type          string          `json:"type"`
	Recipes       []RecipeSummary `json:"recipes"`
	TotalCalories *int            `json:"totalCalories,omitempty"`
}

// HistoryRequest holds query parameters for GET /suggestions/history.
type HistoryRequest struct {
	SessionType string `form:"sessionType"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}
