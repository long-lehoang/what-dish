package engagement

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// EngagementService provides business logic for favorites, views, and ratings.
type EngagementService struct {
	favorites FavoriteRepository
	views     ViewRepository
	ratings   RatingRepository
}

// NewEngagementService creates a new EngagementService.
func NewEngagementService(favorites FavoriteRepository, views ViewRepository, ratings RatingRepository) *EngagementService {
	return &EngagementService{
		favorites: favorites,
		views:     views,
		ratings:   ratings,
	}
}

// AddFavorite adds a recipe to a user's favorites.
func (s *EngagementService) AddFavorite(ctx context.Context, userID uuid.UUID, recipeID uuid.UUID) (*FavoriteResponse, error) {
	if err := s.favorites.Add(ctx, userID, recipeID); err != nil {
		return nil, fmt.Errorf("EngagementService.AddFavorite: %w", err)
	}

	slog.Info("favorite added", "user_id", userID, "recipe_id", recipeID)

	return &FavoriteResponse{
		RecipeID:  recipeID.String(),
		CreatedAt: time.Now().UTC(),
	}, nil
}

// RemoveFavorite removes a recipe from a user's favorites.
func (s *EngagementService) RemoveFavorite(ctx context.Context, userID uuid.UUID, recipeID uuid.UUID) error {
	if err := s.favorites.Remove(ctx, userID, recipeID); err != nil {
		return fmt.Errorf("EngagementService.RemoveFavorite: %w", err)
	}

	slog.Info("favorite removed", "user_id", userID, "recipe_id", recipeID)

	return nil
}

// ListFavorites returns a paginated list of a user's favorites.
func (s *EngagementService) ListFavorites(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]FavoriteResponse, int64, error) {
	items, total, err := s.favorites.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("EngagementService.ListFavorites: %w", err)
	}

	results := make([]FavoriteResponse, len(items))
	for i, f := range items {
		results[i] = FavoriteResponse{
			RecipeID:  f.RecipeID.String(),
			CreatedAt: f.CreatedAt,
		}
	}

	return results, total, nil
}

// CheckFavorites checks which of the given recipe IDs are favorited by the user.
func (s *EngagementService) CheckFavorites(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	result, err := s.favorites.Check(ctx, userID, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("EngagementService.CheckFavorites: %w", err)
	}

	return result, nil
}

// RecordView records a recipe view event.
func (s *EngagementService) RecordView(ctx context.Context, userID *uuid.UUID, sessionID string, recipeID uuid.UUID, source string) error {
	view := &ViewHistory{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
		RecipeID:  recipeID,
		Source:    source,
		ViewedAt:  time.Now().UTC(),
	}

	if err := s.views.Record(ctx, view); err != nil {
		return fmt.Errorf("EngagementService.RecordView: %w", err)
	}

	slog.Debug("view recorded", "recipe_id", recipeID, "source", source)

	return nil
}
