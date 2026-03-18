package engagement

import (
	"time"
)

// AddFavoriteRequest is the payload for POST /favorites.
type AddFavoriteRequest struct {
	RecipeID string `json:"recipeId" validate:"required,uuid"`
}

// FavoriteResponse is the response for favorite operations.
type FavoriteResponse struct {
	RecipeID  string    `json:"recipeId"`
	CreatedAt time.Time `json:"createdAt"`
}

// CheckFavoritesRequest holds recipe IDs to check (parsed from query).
type CheckFavoritesRequest struct {
	RecipeIDs []string `json:"recipeIds" validate:"required,min=1,dive,uuid"`
}

// RecordViewRequest is the payload for POST /views.
type RecordViewRequest struct {
	RecipeID string `json:"recipeId" validate:"required,uuid"`
	Source   string `json:"source" validate:"required,min=1,max=50"`
}

// RatingRequest is the payload for rating a recipe.
type RatingRequest struct {
	RecipeID string  `json:"recipeId" validate:"required,uuid"`
	Score    int     `json:"score" validate:"required,min=1,max=5"`
	Comment  *string `json:"comment" validate:"omitempty,max=1000"`
}
