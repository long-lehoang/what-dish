package engagement

import (
	"context"

	"github.com/google/uuid"
)

// FavoriteRepository manages persistence of user favorites.
type FavoriteRepository interface {
	Add(ctx context.Context, userID, recipeID uuid.UUID) error
	Remove(ctx context.Context, userID, recipeID uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Favorite, int64, error)
	Check(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

// ViewRepository manages persistence of view history.
type ViewRepository interface {
	Record(ctx context.Context, view *ViewHistory) error
}

// RatingRepository manages persistence of recipe ratings.
type RatingRepository interface {
	Upsert(ctx context.Context, rating *Rating) error
	GetByRecipe(ctx context.Context, recipeID uuid.UUID) ([]Rating, error)
}
