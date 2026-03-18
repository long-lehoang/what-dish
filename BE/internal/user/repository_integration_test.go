package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lehoanglong/whatdish/internal/user"
)

func TestIntegration_ProfileRepo_UpsertAndGet(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := user.NewProfileRepo(testDB.Pool)

	userID := uuid.New()
	gender := "MALE"
	age := 30
	height := 175.0
	weight := 70.0
	activity := "MODERATE"
	goal := "MAINTAIN"
	bmr := 1650.0
	tdee := 2557.5
	dailyTarget := 2557.5

	profile := &user.UserProfile{
		ID:            uuid.New(),
		UserID:        userID,
		Gender:        &gender,
		Age:           &age,
		HeightCm:      &height,
		WeightKg:      &weight,
		ActivityLevel: &activity,
		Goal:          &goal,
		BMR:           &bmr,
		TDEE:          &tdee,
		DailyTarget:   &dailyTarget,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	err := repo.Upsert(ctx, profile)
	require.NoError(t, err)

	got, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, *profile.Gender, *got.Gender)
	assert.Equal(t, *profile.Age, *got.Age)
	assert.InDelta(t, *profile.BMR, *got.BMR, 0.01)

	newAge := 31
	profile.Age = &newAge
	profile.UpdatedAt = time.Now().UTC()
	err = repo.Upsert(ctx, profile)
	require.NoError(t, err)

	got, err = repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 31, *got.Age)
}

func TestIntegration_ProfileRepo_GetByUserID_NotFound(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := user.NewProfileRepo(testDB.Pool)

	_, err := repo.GetByUserID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestIntegration_AllergyRepo_SetAndList(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := user.NewAllergyRepo(testDB.Pool)

	userID := uuid.New()
	allergies := []user.UserAllergy{
		{ID: uuid.New(), UserID: userID, IngredientName: "Peanut", AllergyType: "ALLERGY"},
		{ID: uuid.New(), UserID: userID, IngredientName: "Shrimp", AllergyType: "ALLERGY"},
	}

	err := repo.Set(ctx, userID, allergies)
	require.NoError(t, err)

	got, err := repo.ListByUser(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, got, 2)

	newAllergies := []user.UserAllergy{
		{ID: uuid.New(), UserID: userID, IngredientName: "Milk", AllergyType: "DISLIKE"},
	}
	err = repo.Set(ctx, userID, newAllergies)
	require.NoError(t, err)

	got, err = repo.ListByUser(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Milk", got[0].IngredientName)
}

func TestIntegration_AllergyRepo_SetEmpty(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := user.NewAllergyRepo(testDB.Pool)

	userID := uuid.New()
	allergies := []user.UserAllergy{
		{ID: uuid.New(), UserID: userID, IngredientName: "Peanut", AllergyType: "ALLERGY"},
	}

	require.NoError(t, repo.Set(ctx, userID, allergies))
	require.NoError(t, repo.Set(ctx, userID, nil))

	got, err := repo.ListByUser(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, got)
}
