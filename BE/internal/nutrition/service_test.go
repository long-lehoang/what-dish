package nutrition

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Mock implementations (shared across all _test.go files in this package)
// ---------------------------------------------------------------------------

type mockNutritionRepo struct {
	getByRecipeIDFn        func(ctx context.Context, recipeID uuid.UUID) (*RecipeNutrition, error)
	getByRecipeIDsFn       func(ctx context.Context, recipeIDs []uuid.UUID) ([]RecipeNutrition, error)
	upsertFn               func(ctx context.Context, nutrition *RecipeNutrition) error
	getIDsByCalorieRangeFn func(ctx context.Context, min, max float64) ([]uuid.UUID, error)
}

func (m *mockNutritionRepo) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) (*RecipeNutrition, error) {
	if m.getByRecipeIDFn != nil {
		return m.getByRecipeIDFn(ctx, recipeID)
	}
	return nil, nil
}

func (m *mockNutritionRepo) GetByRecipeIDs(ctx context.Context, recipeIDs []uuid.UUID) ([]RecipeNutrition, error) {
	if m.getByRecipeIDsFn != nil {
		return m.getByRecipeIDsFn(ctx, recipeIDs)
	}
	return nil, nil
}

func (m *mockNutritionRepo) Upsert(ctx context.Context, nutrition *RecipeNutrition) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, nutrition)
	}
	return nil
}

func (m *mockNutritionRepo) GetIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error) {
	if m.getIDsByCalorieRangeFn != nil {
		return m.getIDsByCalorieRangeFn(ctx, min, max)
	}
	return nil, nil
}

type mockGoalRepo struct {
	listFn func(ctx context.Context) ([]NutritionGoal, error)
}

func (m *mockGoalRepo) List(ctx context.Context) ([]NutritionGoal, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests: NutritionService.CalculateTDEE
// ---------------------------------------------------------------------------

func TestNutritionService_CalculateTDEE_MaleModerateMaintain(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})

	req := CalculateTDEERequest{
		Gender:        "MALE",
		Age:           25,
		HeightCm:      175,
		WeightKg:      70,
		ActivityLevel: "MODERATE",
		Goal:          "MAINTAIN",
	}

	resp, err := svc.CalculateTDEE(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// BMR = 10*70 + 6.25*175 - 5*25 + 5 = 1673.75
	expectedBMR := CalculateBMR("MALE", 70, 175, 25)
	assert.Equal(t, expectedBMR, resp.BMR)

	// TDEE = BMR * 1.55
	expectedTDEE := CalculateTDEE(expectedBMR, "MODERATE")
	assert.Equal(t, expectedTDEE, resp.TDEE)

	// DailyTarget = TDEE * 1.0 (MAINTAIN)
	expectedDaily := AdjustForGoal(expectedTDEE, "MAINTAIN")
	assert.Equal(t, expectedDaily, resp.DailyTarget)

	// Meal breakdown should have 3 keys.
	assert.Len(t, resp.MealBreakdown, 3)
	assert.Contains(t, resp.MealBreakdown, "breakfast")
	assert.Contains(t, resp.MealBreakdown, "lunch")
	assert.Contains(t, resp.MealBreakdown, "dinner")
}

func TestNutritionService_CalculateTDEE_FemaleLoseWeight(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})

	req := CalculateTDEERequest{
		Gender:        "FEMALE",
		Age:           30,
		HeightCm:      165,
		WeightKg:      60,
		ActivityLevel: "LIGHT",
		Goal:          "LOSE_WEIGHT",
	}

	resp, err := svc.CalculateTDEE(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	expectedBMR := CalculateBMR("FEMALE", 60, 165, 30)
	expectedTDEE := CalculateTDEE(expectedBMR, "LIGHT")
	expectedDaily := AdjustForGoal(expectedTDEE, "LOSE_WEIGHT")

	assert.Equal(t, expectedBMR, resp.BMR)
	assert.Equal(t, expectedTDEE, resp.TDEE)
	assert.Equal(t, expectedDaily, resp.DailyTarget)

	// LOSE_WEIGHT daily target should be less than TDEE.
	assert.Less(t, resp.DailyTarget, resp.TDEE)
}

func TestNutritionService_CalculateTDEE_ResponseStructure(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})

	req := CalculateTDEERequest{
		Gender:        "MALE",
		Age:           40,
		HeightCm:      180,
		WeightKg:      85,
		ActivityLevel: "ACTIVE",
		Goal:          "GAIN_WEIGHT",
	}

	resp, err := svc.CalculateTDEE(context.Background(), req)

	assert.NoError(t, err)
	assert.Greater(t, resp.BMR, 0.0)
	assert.Greater(t, resp.TDEE, resp.BMR)
	assert.Greater(t, resp.DailyTarget, resp.TDEE) // GAIN_WEIGHT means daily > TDEE

	// Meal breakdown should sum approximately to daily target.
	sum := resp.MealBreakdown["breakfast"] + resp.MealBreakdown["lunch"] + resp.MealBreakdown["dinner"]
	assert.InDelta(t, resp.DailyTarget, sum, 0.02)
}

// ---------------------------------------------------------------------------
// Tests: NutritionService.GetByRecipeID
// ---------------------------------------------------------------------------

func TestNutritionService_GetByRecipeID_Success(t *testing.T) {
	recipeID := uuid.New()
	cal := 350.0
	protein := 25.0
	nutrition := &RecipeNutrition{
		ID:       uuid.New(),
		RecipeID: recipeID,
		Calories: &cal,
		Protein:  &protein,
	}

	svc := NewNutritionService(
		&mockNutritionRepo{
			getByRecipeIDFn: func(_ context.Context, id uuid.UUID) (*RecipeNutrition, error) {
				assert.Equal(t, recipeID, id)
				return nutrition, nil
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetByRecipeID(context.Background(), recipeID)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, recipeID.String(), got.RecipeID)
	assert.Equal(t, &cal, got.Calories)
	assert.Equal(t, &protein, got.Protein)
}

func TestNutritionService_GetByRecipeID_NotFound(t *testing.T) {
	notFoundErr := errors.New("not found")
	recipeID := uuid.New()

	svc := NewNutritionService(
		&mockNutritionRepo{
			getByRecipeIDFn: func(_ context.Context, _ uuid.UUID) (*RecipeNutrition, error) {
				return nil, notFoundErr
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetByRecipeID(context.Background(), recipeID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, notFoundErr)
	assert.Contains(t, err.Error(), "NutritionService.GetByRecipeID")
	assert.Nil(t, got)
}

func TestNutritionService_GetByRecipeID_AllFieldsMapped(t *testing.T) {
	recipeID := uuid.New()
	cal := 400.0
	protein := 30.0
	carbs := 50.0
	fat := 15.0
	fiber := 5.0
	sugar := 8.0
	sodium := 600.0
	servingSize := "1 serving (250g)"
	dataSource := "manual"

	nutrition := &RecipeNutrition{
		ID:          uuid.New(),
		RecipeID:    recipeID,
		Calories:    &cal,
		Protein:     &protein,
		Carbs:       &carbs,
		Fat:         &fat,
		Fiber:       &fiber,
		Sugar:       &sugar,
		Sodium:      &sodium,
		ServingSize: &servingSize,
		DataSource:  &dataSource,
		IsVerified:  true,
	}

	svc := NewNutritionService(
		&mockNutritionRepo{
			getByRecipeIDFn: func(_ context.Context, _ uuid.UUID) (*RecipeNutrition, error) {
				return nutrition, nil
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetByRecipeID(context.Background(), recipeID)

	assert.NoError(t, err)
	assert.Equal(t, &cal, got.Calories)
	assert.Equal(t, &protein, got.Protein)
	assert.Equal(t, &carbs, got.Carbs)
	assert.Equal(t, &fat, got.Fat)
	assert.Equal(t, &fiber, got.Fiber)
	assert.Equal(t, &sugar, got.Sugar)
	assert.Equal(t, &sodium, got.Sodium)
	assert.Equal(t, &servingSize, got.ServingSize)
	assert.Equal(t, &dataSource, got.DataSource)
	assert.True(t, got.IsVerified)
}

// ---------------------------------------------------------------------------
// Tests: NutritionService.ListGoals
// ---------------------------------------------------------------------------

func TestNutritionService_ListGoals_Success(t *testing.T) {
	goals := []NutritionGoal{
		{ID: uuid.New(), Name: "Weight Loss", Description: "Lose weight safely", IsActive: true, SortOrder: 1},
		{ID: uuid.New(), Name: "Maintenance", Description: "Maintain current weight", IsActive: true, SortOrder: 2},
		{ID: uuid.New(), Name: "Muscle Gain", Description: "Gain muscle mass", IsActive: true, SortOrder: 3},
	}

	svc := NewNutritionService(
		&mockNutritionRepo{},
		&mockGoalRepo{
			listFn: func(_ context.Context) ([]NutritionGoal, error) {
				return goals, nil
			},
		},
	)

	got, err := svc.ListGoals(context.Background())

	assert.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "Weight Loss", got[0].Name)
	assert.Equal(t, "Maintenance", got[1].Name)
	assert.Equal(t, "Muscle Gain", got[2].Name)
}

func TestNutritionService_ListGoals_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	svc := NewNutritionService(
		&mockNutritionRepo{},
		&mockGoalRepo{
			listFn: func(_ context.Context) ([]NutritionGoal, error) {
				return nil, repoErr
			},
		},
	)

	got, err := svc.ListGoals(context.Background())

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "NutritionService.ListGoals")
	assert.Nil(t, got)
}

func TestNutritionService_ListGoals_EmptyResult(t *testing.T) {
	svc := NewNutritionService(
		&mockNutritionRepo{},
		&mockGoalRepo{
			listFn: func(_ context.Context) ([]NutritionGoal, error) {
				return []NutritionGoal{}, nil
			},
		},
	)

	got, err := svc.ListGoals(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, got)
}

// ---------------------------------------------------------------------------
// Tests: NutritionService.GetIDsByCalorieRange
// ---------------------------------------------------------------------------

func TestNutritionService_GetIDsByCalorieRange_Success(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	svc := NewNutritionService(
		&mockNutritionRepo{
			getIDsByCalorieRangeFn: func(_ context.Context, min, max float64) ([]uuid.UUID, error) {
				assert.Equal(t, 300.0, min)
				assert.Equal(t, 500.0, max)
				return []uuid.UUID{id1, id2, id3}, nil
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetIDsByCalorieRange(context.Background(), 300, 500)

	assert.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Contains(t, got, id1)
	assert.Contains(t, got, id2)
	assert.Contains(t, got, id3)
}

func TestNutritionService_GetIDsByCalorieRange_Empty(t *testing.T) {
	svc := NewNutritionService(
		&mockNutritionRepo{
			getIDsByCalorieRangeFn: func(_ context.Context, _, _ float64) ([]uuid.UUID, error) {
				return []uuid.UUID{}, nil
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetIDsByCalorieRange(context.Background(), 10, 50)

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestNutritionService_GetIDsByCalorieRange_RepoError(t *testing.T) {
	repoErr := errors.New("query failed")

	svc := NewNutritionService(
		&mockNutritionRepo{
			getIDsByCalorieRangeFn: func(_ context.Context, _, _ float64) ([]uuid.UUID, error) {
				return nil, repoErr
			},
		},
		&mockGoalRepo{},
	)

	got, err := svc.GetIDsByCalorieRange(context.Background(), 300, 500)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "NutritionService.GetIDsByCalorieRange")
	assert.Nil(t, got)
}
