package engagement_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lehoanglong/whatdish/internal/engagement"
)

func seedRecipeForEngagement(t *testing.T, name string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var id uuid.UUID
	err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO recipes (name, slug, status) VALUES ($1, $2, 'PUBLISHED') RETURNING id`,
		name, "slug-"+uuid.New().String()[:8],
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestIntegration_FavoriteRepo_AddAndList(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := seedRecipeForEngagement(t, "Phở")

	err := repo.Add(ctx, userID, recipeID)
	require.NoError(t, err)

	favs, total, err := repo.ListByUser(ctx, userID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, recipeID, favs[0].RecipeID)
}

func TestIntegration_FavoriteRepo_AddDuplicate(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := seedRecipeForEngagement(t, "Phở")

	require.NoError(t, repo.Add(ctx, userID, recipeID))
	require.NoError(t, repo.Add(ctx, userID, recipeID))

	_, total, err := repo.ListByUser(ctx, userID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
}

func TestIntegration_FavoriteRepo_Remove(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := seedRecipeForEngagement(t, "Phở")

	require.NoError(t, repo.Add(ctx, userID, recipeID))
	require.NoError(t, repo.Remove(ctx, userID, recipeID))

	_, total, err := repo.ListByUser(ctx, userID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
}

func TestIntegration_FavoriteRepo_RemoveNotFound(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	err := repo.Remove(ctx, uuid.New(), uuid.New())
	assert.Error(t, err)
}

func TestIntegration_FavoriteRepo_Check(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	userID := uuid.New()
	r1 := seedRecipeForEngagement(t, "Phở")
	r2 := seedRecipeForEngagement(t, "Bún")

	require.NoError(t, repo.Add(ctx, userID, r1))

	result, err := repo.Check(ctx, userID, []uuid.UUID{r1, r2})
	require.NoError(t, err)
	assert.True(t, result[r1])
	assert.False(t, result[r2])
}

func TestIntegration_FavoriteRepo_ListPagination(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewFavoriteRepo(testDB.Pool)

	userID := uuid.New()
	for i := 0; i < 5; i++ {
		r := seedRecipeForEngagement(t, "Dish "+uuid.New().String()[:4])
		require.NoError(t, repo.Add(ctx, userID, r))
	}

	favs, total, err := repo.ListByUser(ctx, userID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, favs, 2)
}

func TestIntegration_ViewRepo_Record(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewViewRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := seedRecipeForEngagement(t, "Phở")

	view := &engagement.ViewHistory{
		ID:        uuid.New(),
		UserID:    &userID,
		SessionID: "session-123",
		RecipeID:  recipeID,
		Source:    "suggestion",
		ViewedAt:  time.Now().UTC(),
	}

	err := repo.Record(ctx, view)
	require.NoError(t, err)

	var count int
	err = testDB.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM engagement_views WHERE recipe_id = $1", recipeID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestIntegration_RatingRepo_Upsert(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := engagement.NewRatingRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := seedRecipeForEngagement(t, "Phở")

	rating := &engagement.Rating{
		ID:        uuid.New(),
		UserID:    userID,
		RecipeID:  recipeID,
		Score:     4,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	require.NoError(t, repo.Upsert(ctx, rating))

	ratings, err := repo.GetByRecipe(ctx, recipeID)
	require.NoError(t, err)
	require.Len(t, ratings, 1)
	assert.Equal(t, 4, ratings[0].Score)

	rating.Score = 5
	rating.UpdatedAt = time.Now().UTC()
	require.NoError(t, repo.Upsert(ctx, rating))

	ratings, err = repo.GetByRecipe(ctx, recipeID)
	require.NoError(t, err)
	require.Len(t, ratings, 1)
	assert.Equal(t, 5, ratings[0].Score)
}
