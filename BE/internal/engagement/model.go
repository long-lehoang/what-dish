package engagement

import (
	"time"

	"github.com/google/uuid"
)

// Favorite represents a user's bookmarked recipe.
type Favorite struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	RecipeID  uuid.UUID `json:"recipeId"`
	CreatedAt time.Time `json:"createdAt"`
}

// ViewHistory records a single recipe view event.
type ViewHistory struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"userId,omitempty"`
	SessionID string     `json:"sessionId"`
	RecipeID  uuid.UUID  `json:"recipeId"`
	Source    string     `json:"source"`
	ViewedAt  time.Time  `json:"viewedAt"`
}

// Rating represents a user's star rating and optional comment on a recipe.
type Rating struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	RecipeID  uuid.UUID `json:"recipeId"`
	Score     int       `json:"score"` // 1-5
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
