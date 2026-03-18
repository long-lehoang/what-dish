package suggestion_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lehoanglong/whatdish/internal/suggestion"
)

func TestIntegration_SessionRepo_CreateAndList(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewSessionRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := uuid.New()

	session := &suggestion.SuggestionSession{
		ID:              uuid.New(),
		UserID:          &userID,
		SessionType:     "RANDOM",
		InputParams:     map[string]any{"max_cook_time": 30},
		ResultRecipeIDs: []uuid.UUID{recipeID},
		CreatedAt:       time.Now().UTC(),
	}

	err := repo.Create(ctx, session)
	require.NoError(t, err)

	sessions, total, err := repo.ListByUser(ctx, userID, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "RANDOM", sessions[0].SessionType)
	assert.Len(t, sessions[0].ResultRecipeIDs, 1)
}

func TestIntegration_SessionRepo_ListByType(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewSessionRepo(testDB.Pool)

	userID := uuid.New()
	now := time.Now().UTC()

	require.NoError(t, repo.Create(ctx, &suggestion.SuggestionSession{
		ID: uuid.New(), UserID: &userID, SessionType: "RANDOM",
		InputParams: map[string]any{}, ResultRecipeIDs: []uuid.UUID{}, CreatedAt: now,
	}))
	require.NoError(t, repo.Create(ctx, &suggestion.SuggestionSession{
		ID: uuid.New(), UserID: &userID, SessionType: "BY_CALORIES",
		InputParams: map[string]any{}, ResultRecipeIDs: []uuid.UUID{}, CreatedAt: now.Add(time.Second),
	}))
	require.NoError(t, repo.Create(ctx, &suggestion.SuggestionSession{
		ID: uuid.New(), UserID: &userID, SessionType: "RANDOM",
		InputParams: map[string]any{}, ResultRecipeIDs: []uuid.UUID{}, CreatedAt: now.Add(2 * time.Second),
	}))

	sessions, total, err := repo.ListByUser(ctx, userID, "RANDOM", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, sessions, 2)

	sessions, total, err = repo.ListByUser(ctx, userID, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, sessions, 3)
}

func TestIntegration_SessionRepo_ListPagination(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewSessionRepo(testDB.Pool)

	userID := uuid.New()
	now := time.Now().UTC()

	for i := 0; i < 5; i++ {
		require.NoError(t, repo.Create(ctx, &suggestion.SuggestionSession{
			ID: uuid.New(), UserID: &userID, SessionType: "RANDOM",
			InputParams: map[string]any{}, ResultRecipeIDs: []uuid.UUID{},
			CreatedAt: now.Add(time.Duration(i) * time.Second),
		}))
	}

	sessions, total, err := repo.ListByUser(ctx, userID, "", 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, sessions, 2)
}

func TestIntegration_ConfigRepo_ListByType(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewConfigRepo(testDB.Pool)

	_, err := testDB.Pool.Exec(ctx,
		`INSERT INTO suggestion_configs (config_type, name, description, params, is_active, sort_order)
		 VALUES ('CALORIE', 'Low Cal', 'Low calorie preset', '{"min": 200, "max": 400}', true, 1),
		        ('CALORIE', 'High Cal', 'High calorie preset', '{"min": 600, "max": 1000}', true, 2),
		        ('GROUP', 'Family', 'Family preset', '{"group_size": 4}', true, 1),
		        ('CALORIE', 'Inactive', 'Hidden', '{}', false, 3)`)
	require.NoError(t, err)

	configs, err := repo.ListByType(ctx, "CALORIE")
	require.NoError(t, err)
	assert.Len(t, configs, 2)

	configs, err = repo.ListByType(ctx, "")
	require.NoError(t, err)
	assert.Len(t, configs, 3)
}

func TestIntegration_ExclusionRepo_AddAndGet(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewExclusionRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := uuid.New()

	rule := &suggestion.ExclusionRule{
		ID:            uuid.New(),
		UserID:        userID,
		RecipeID:      recipeID,
		ExcludedUntil: time.Now().UTC().Add(24 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}

	require.NoError(t, repo.Add(ctx, rule))

	ids, err := repo.GetExcludedRecipeIDs(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.Equal(t, recipeID, ids[0])
}

func TestIntegration_ExclusionRepo_ExpiredNotReturned(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewExclusionRepo(testDB.Pool)

	userID := uuid.New()

	require.NoError(t, repo.Add(ctx, &suggestion.ExclusionRule{
		ID: uuid.New(), UserID: userID, RecipeID: uuid.New(),
		ExcludedUntil: time.Now().UTC().Add(24 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}))
	require.NoError(t, repo.Add(ctx, &suggestion.ExclusionRule{
		ID: uuid.New(), UserID: userID, RecipeID: uuid.New(),
		ExcludedUntil: time.Now().UTC().Add(-1 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}))

	ids, err := repo.GetExcludedRecipeIDs(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
}

func TestIntegration_ExclusionRepo_CleanExpired(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewExclusionRepo(testDB.Pool)

	userID := uuid.New()

	require.NoError(t, repo.Add(ctx, &suggestion.ExclusionRule{
		ID: uuid.New(), UserID: userID, RecipeID: uuid.New(),
		ExcludedUntil: time.Now().UTC().Add(-1 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}))

	require.NoError(t, repo.CleanExpired(ctx))

	var count int
	err := testDB.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM suggestion_exclusions WHERE user_id = $1", userID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestIntegration_ExclusionRepo_UpsertOnConflict(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := suggestion.NewExclusionRepo(testDB.Pool)

	userID := uuid.New()
	recipeID := uuid.New()

	require.NoError(t, repo.Add(ctx, &suggestion.ExclusionRule{
		ID: uuid.New(), UserID: userID, RecipeID: recipeID,
		ExcludedUntil: time.Now().UTC().Add(1 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}))

	require.NoError(t, repo.Add(ctx, &suggestion.ExclusionRule{
		ID: uuid.New(), UserID: userID, RecipeID: recipeID,
		ExcludedUntil: time.Now().UTC().Add(48 * time.Hour),
		CreatedAt:     time.Now().UTC(),
	}))

	var count int
	err := testDB.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM suggestion_exclusions WHERE user_id = $1 AND recipe_id = $2",
		userID, recipeID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
