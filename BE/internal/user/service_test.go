package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockAuthProvider struct {
	registerFn     func(ctx context.Context, email, password, name string) (*AuthUser, *AuthTokens, error)
	loginFn        func(ctx context.Context, email, password string) (*AuthUser, *AuthTokens, error)
	refreshTokenFn func(ctx context.Context, refreshToken string) (*AuthTokens, error)
	verifyTokenFn  func(ctx context.Context, token string) (uuid.UUID, error)
	getUserFn      func(ctx context.Context, token string) (*AuthUser, error)
}

func (m *mockAuthProvider) Register(ctx context.Context, email, password, name string) (*AuthUser, *AuthTokens, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, email, password, name)
	}
	return nil, nil, nil
}

func (m *mockAuthProvider) Login(ctx context.Context, email, password string) (*AuthUser, *AuthTokens, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, email, password)
	}
	return nil, nil, nil
}

func (m *mockAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	if m.refreshTokenFn != nil {
		return m.refreshTokenFn(ctx, refreshToken)
	}
	return nil, nil
}

func (m *mockAuthProvider) VerifyToken(ctx context.Context, token string) (uuid.UUID, error) {
	if m.verifyTokenFn != nil {
		return m.verifyTokenFn(ctx, token)
	}
	return uuid.Nil, nil
}

func (m *mockAuthProvider) GetUser(ctx context.Context, token string) (*AuthUser, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, token)
	}
	return nil, nil
}

type mockProfileRepo struct {
	getByUserIDFn func(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	upsertFn      func(ctx context.Context, profile *UserProfile) error
}

func (m *mockProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockProfileRepo) Upsert(ctx context.Context, profile *UserProfile) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, profile)
	}
	return nil
}

type mockAllergyRepo struct {
	listByUserFn func(ctx context.Context, userID uuid.UUID) ([]UserAllergy, error)
	setFn        func(ctx context.Context, userID uuid.UUID, allergies []UserAllergy) error
}

func (m *mockAllergyRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]UserAllergy, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockAllergyRepo) Set(ctx context.Context, userID uuid.UUID, allergies []UserAllergy) error {
	if m.setFn != nil {
		return m.setFn(ctx, userID, allergies)
	}
	return nil
}

type mockTDEECalculator struct {
	calculateBMRFn  func(gender string, weightKg, heightCm float64, age int) float64
	calculateTDEEFn func(bmr float64, activityLevel string) float64
	adjustForGoalFn func(tdee float64, goal string) float64
}

func (m *mockTDEECalculator) CalculateBMR(gender string, weightKg, heightCm float64, age int) float64 {
	if m.calculateBMRFn != nil {
		return m.calculateBMRFn(gender, weightKg, heightCm, age)
	}
	return 0
}

func (m *mockTDEECalculator) CalculateTDEE(bmr float64, activityLevel string) float64 {
	if m.calculateTDEEFn != nil {
		return m.calculateTDEEFn(bmr, activityLevel)
	}
	return 0
}

func (m *mockTDEECalculator) AdjustForGoal(tdee float64, goal string) float64 {
	if m.adjustForGoalFn != nil {
		return m.adjustForGoalFn(tdee, goal)
	}
	return 0
}

// ---------------------------------------------------------------------------
// Tests: AuthService.Register
// ---------------------------------------------------------------------------

func TestAuthService_Register_Success(t *testing.T) {
	userID := uuid.New()
	expectedUser := &AuthUser{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	expectedTokens := &AuthTokens{
		AccessToken:  "access-token-xyz",
		RefreshToken: "refresh-token-xyz",
		ExpiresIn:    3600,
	}

	svc := NewAuthService(&mockAuthProvider{
		registerFn: func(_ context.Context, email, password, name string) (*AuthUser, *AuthTokens, error) {
			assert.Equal(t, "test@example.com", email)
			assert.Equal(t, "password123", password)
			assert.Equal(t, "Test User", name)
			return expectedUser, expectedTokens, nil
		},
	})

	resp, err := svc.Register(context.Background(), RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID, resp.User.ID)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "access-token-xyz", resp.Tokens.AccessToken)
	assert.Equal(t, "refresh-token-xyz", resp.Tokens.RefreshToken)
	assert.Equal(t, 3600, resp.Tokens.ExpiresIn)
}

func TestAuthService_Register_ProviderError(t *testing.T) {
	providerErr := errors.New("email already registered")

	svc := NewAuthService(&mockAuthProvider{
		registerFn: func(_ context.Context, _, _, _ string) (*AuthUser, *AuthTokens, error) {
			return nil, nil, providerErr
		},
	})

	resp, err := svc.Register(context.Background(), RegisterRequest{
		Email:    "dup@example.com",
		Password: "password123",
		Name:     "Dup User",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providerErr)
	assert.Contains(t, err.Error(), "AuthService.Register")
	assert.Nil(t, resp)
}

// ---------------------------------------------------------------------------
// Tests: AuthService.Login
// ---------------------------------------------------------------------------

func TestAuthService_Login_Success(t *testing.T) {
	userID := uuid.New()
	expectedUser := &AuthUser{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	expectedTokens := &AuthTokens{
		AccessToken:  "access-abc",
		RefreshToken: "refresh-abc",
		ExpiresIn:    3600,
	}

	svc := NewAuthService(&mockAuthProvider{
		loginFn: func(_ context.Context, email, password string) (*AuthUser, *AuthTokens, error) {
			assert.Equal(t, "test@example.com", email)
			assert.Equal(t, "password123", password)
			return expectedUser, expectedTokens, nil
		},
	})

	resp, err := svc.Login(context.Background(), LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID, resp.User.ID)
	assert.Equal(t, "access-abc", resp.Tokens.AccessToken)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	authErr := errors.New("invalid credentials")

	svc := NewAuthService(&mockAuthProvider{
		loginFn: func(_ context.Context, _, _ string) (*AuthUser, *AuthTokens, error) {
			return nil, nil, authErr
		},
	})

	resp, err := svc.Login(context.Background(), LoginRequest{
		Email:    "test@example.com",
		Password: "wrong",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, authErr)
	assert.Contains(t, err.Error(), "AuthService.Login")
	assert.Nil(t, resp)
}

// ---------------------------------------------------------------------------
// Tests: ProfileService.GetProfile
// ---------------------------------------------------------------------------

func TestProfileService_GetProfile_Success(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	gender := "MALE"
	age := 25
	height := 175.0
	weight := 70.0
	activity := "MODERATE"
	goal := "MAINTAIN"
	bmr := 1673.75
	tdee := 2594.31
	daily := 2594.31

	profile := &UserProfile{
		ID:            profileID,
		UserID:        userID,
		Gender:        &gender,
		Age:           &age,
		HeightCm:      &height,
		WeightKg:      &weight,
		ActivityLevel: &activity,
		Goal:          &goal,
		BMR:           &bmr,
		TDEE:          &tdee,
		DailyTarget:   &daily,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	allergies := []UserAllergy{
		{ID: uuid.New(), UserID: userID, IngredientName: "Peanuts", AllergyType: "ALLERGY"},
		{ID: uuid.New(), UserID: userID, IngredientName: "Cilantro", AllergyType: "DISLIKE"},
	}

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, id uuid.UUID) (*UserProfile, error) {
				assert.Equal(t, userID, id)
				return profile, nil
			},
		},
		&mockAllergyRepo{
			listByUserFn: func(_ context.Context, id uuid.UUID) ([]UserAllergy, error) {
				assert.Equal(t, userID, id)
				return allergies, nil
			},
		},
		&mockTDEECalculator{},
	)

	got, err := svc.GetProfile(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, profileID.String(), got.ID)
	assert.Equal(t, userID.String(), got.UserID)
	assert.Equal(t, &gender, got.Gender)
	assert.Equal(t, &age, got.Age)
	assert.Equal(t, &height, got.HeightCm)
	assert.Equal(t, &weight, got.WeightKg)
	assert.Equal(t, &activity, got.ActivityLevel)
	assert.Equal(t, &goal, got.Goal)
	assert.Equal(t, &bmr, got.BMR)
	assert.Equal(t, &tdee, got.TDEE)
	assert.Equal(t, &daily, got.DailyTarget)
	assert.Len(t, got.Allergies, 2)
	assert.Equal(t, "Peanuts", got.Allergies[0].IngredientName)
	assert.Equal(t, "ALLERGY", got.Allergies[0].AllergyType)
	assert.Equal(t, "Cilantro", got.Allergies[1].IngredientName)
	assert.Equal(t, "DISLIKE", got.Allergies[1].AllergyType)
}

func TestProfileService_GetProfile_NoAllergies(t *testing.T) {
	userID := uuid.New()
	gender := "FEMALE"

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				return &UserProfile{
					ID:     uuid.New(),
					UserID: userID,
					Gender: &gender,
				}, nil
			},
		},
		&mockAllergyRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID) ([]UserAllergy, error) {
				return []UserAllergy{}, nil
			},
		},
		&mockTDEECalculator{},
	)

	got, err := svc.GetProfile(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Empty(t, got.Allergies)
}

func TestProfileService_GetProfile_ProfileNotFound(t *testing.T) {
	notFoundErr := errors.New("not found")

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				return nil, notFoundErr
			},
		},
		&mockAllergyRepo{},
		&mockTDEECalculator{},
	)

	got, err := svc.GetProfile(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.ErrorIs(t, err, notFoundErr)
	assert.Contains(t, err.Error(), "ProfileService.GetProfile")
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: ProfileService.UpdateProfile
// ---------------------------------------------------------------------------

func TestProfileService_UpdateProfile_CreatesNewProfile(t *testing.T) {
	userID := uuid.New()
	gender := "MALE"
	age := 25
	height := 175.0
	weight := 70.0
	activity := "MODERATE"
	goal := "MAINTAIN"

	profileNotFoundErr := errors.New("not found")

	// Track calls for verification.
	var upsertedProfile *UserProfile
	getCallCount := 0

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, id uuid.UUID) (*UserProfile, error) {
				getCallCount++
				if getCallCount == 1 {
					// First call: profile does not exist yet.
					return nil, profileNotFoundErr
				}
				// Second call (from GetProfile at the end of UpdateProfile):
				// Return the upserted profile.
				return upsertedProfile, nil
			},
			upsertFn: func(_ context.Context, p *UserProfile) error {
				upsertedProfile = p
				return nil
			},
		},
		&mockAllergyRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID) ([]UserAllergy, error) {
				return []UserAllergy{}, nil
			},
		},
		&mockTDEECalculator{
			calculateBMRFn: func(g string, w, h float64, a int) float64 {
				assert.Equal(t, "MALE", g)
				return 1673.75
			},
			calculateTDEEFn: func(bmr float64, al string) float64 {
				assert.Equal(t, 1673.75, bmr)
				assert.Equal(t, "MODERATE", al)
				return 2594.31
			},
			adjustForGoalFn: func(tdee float64, g string) float64 {
				assert.Equal(t, 2594.31, tdee)
				assert.Equal(t, "MAINTAIN", g)
				return 2594.31
			},
		},
	)

	resp, err := svc.UpdateProfile(context.Background(), userID, UpdateProfileRequest{
		Gender:        &gender,
		Age:           &age,
		HeightCm:      &height,
		WeightKg:      &weight,
		ActivityLevel: &activity,
		Goal:          &goal,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify the profile was created with correct fields.
	assert.NotNil(t, upsertedProfile)
	assert.Equal(t, userID, upsertedProfile.UserID)
	assert.Equal(t, &gender, upsertedProfile.Gender)
	assert.Equal(t, &age, upsertedProfile.Age)
	assert.Equal(t, &height, upsertedProfile.HeightCm)
	assert.Equal(t, &weight, upsertedProfile.WeightKg)
	assert.Equal(t, &activity, upsertedProfile.ActivityLevel)
	assert.Equal(t, &goal, upsertedProfile.Goal)

	// TDEE should have been calculated.
	expectedBMR := 1673.75
	expectedTDEE := 2594.31
	expectedDaily := 2594.31
	assert.Equal(t, &expectedBMR, upsertedProfile.BMR)
	assert.Equal(t, &expectedTDEE, upsertedProfile.TDEE)
	assert.Equal(t, &expectedDaily, upsertedProfile.DailyTarget)
}

func TestProfileService_UpdateProfile_UpdatesExistingProfile(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	oldGender := "MALE"
	oldAge := 25
	oldHeight := 175.0
	oldWeight := 70.0
	oldActivity := "SEDENTARY"
	oldGoal := "MAINTAIN"

	existingProfile := &UserProfile{
		ID:            profileID,
		UserID:        userID,
		Gender:        &oldGender,
		Age:           &oldAge,
		HeightCm:      &oldHeight,
		WeightKg:      &oldWeight,
		ActivityLevel: &oldActivity,
		Goal:          &oldGoal,
		CreatedAt:     time.Now().UTC().Add(-24 * time.Hour),
		UpdatedAt:     time.Now().UTC().Add(-24 * time.Hour),
	}

	var upsertedProfile *UserProfile
	getCallCount := 0

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				getCallCount++
				if getCallCount == 1 {
					return existingProfile, nil
				}
				return upsertedProfile, nil
			},
			upsertFn: func(_ context.Context, p *UserProfile) error {
				upsertedProfile = p
				return nil
			},
		},
		&mockAllergyRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID) ([]UserAllergy, error) {
				return []UserAllergy{}, nil
			},
		},
		&mockTDEECalculator{
			calculateBMRFn: func(_ string, _ float64, _ float64, _ int) float64 {
				return 1673.75
			},
			calculateTDEEFn: func(_ float64, _ string) float64 {
				return 2302.66
			},
			adjustForGoalFn: func(_ float64, _ string) float64 {
				return 2302.66
			},
		},
	)

	newActivity := "ACTIVE"
	resp, err := svc.UpdateProfile(context.Background(), userID, UpdateProfileRequest{
		ActivityLevel: &newActivity,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Should keep the existing profile ID.
	assert.Equal(t, profileID, upsertedProfile.ID)
	// Activity level should be updated.
	assert.Equal(t, &newActivity, upsertedProfile.ActivityLevel)
	// Other fields should remain unchanged.
	assert.Equal(t, &oldGender, upsertedProfile.Gender)
	assert.Equal(t, &oldAge, upsertedProfile.Age)
}

func TestProfileService_UpdateProfile_WithAllergies(t *testing.T) {
	userID := uuid.New()
	gender := "FEMALE"
	age := 30
	height := 165.0
	weight := 60.0
	activity := "LIGHT"
	goal := "LOSE_WEIGHT"

	var capturedAllergies []UserAllergy
	getCallCount := 0

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				getCallCount++
				if getCallCount == 1 {
					return nil, errors.New("not found")
				}
				return &UserProfile{
					ID:            uuid.New(),
					UserID:        userID,
					Gender:        &gender,
					Age:           &age,
					HeightCm:      &height,
					WeightKg:      &weight,
					ActivityLevel: &activity,
					Goal:          &goal,
				}, nil
			},
			upsertFn: func(_ context.Context, _ *UserProfile) error {
				return nil
			},
		},
		&mockAllergyRepo{
			setFn: func(_ context.Context, uid uuid.UUID, allergies []UserAllergy) error {
				assert.Equal(t, userID, uid)
				capturedAllergies = allergies
				return nil
			},
			listByUserFn: func(_ context.Context, _ uuid.UUID) ([]UserAllergy, error) {
				return capturedAllergies, nil
			},
		},
		&mockTDEECalculator{
			calculateBMRFn:  func(_ string, _, _ float64, _ int) float64 { return 1320.25 },
			calculateTDEEFn: func(_ float64, _ string) float64 { return 1815.34 },
			adjustForGoalFn: func(_ float64, _ string) float64 { return 1452.27 },
		},
	)

	resp, err := svc.UpdateProfile(context.Background(), userID, UpdateProfileRequest{
		Gender:        &gender,
		Age:           &age,
		HeightCm:      &height,
		WeightKg:      &weight,
		ActivityLevel: &activity,
		Goal:          &goal,
		Allergies: []AllergyDTO{
			{IngredientName: "Shrimp", AllergyType: "ALLERGY"},
			{IngredientName: "Onion", AllergyType: "DISLIKE"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, capturedAllergies, 2)
	assert.Equal(t, userID, capturedAllergies[0].UserID)
	assert.Equal(t, "Shrimp", capturedAllergies[0].IngredientName)
	assert.Equal(t, "ALLERGY", capturedAllergies[0].AllergyType)
	assert.Equal(t, "Onion", capturedAllergies[1].IngredientName)
	assert.Equal(t, "DISLIKE", capturedAllergies[1].AllergyType)
}

func TestProfileService_UpdateProfile_UpsertError(t *testing.T) {
	upsertErr := errors.New("db write failed")

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				return nil, errors.New("not found")
			},
			upsertFn: func(_ context.Context, _ *UserProfile) error {
				return upsertErr
			},
		},
		&mockAllergyRepo{},
		&mockTDEECalculator{},
	)

	resp, err := svc.UpdateProfile(context.Background(), uuid.New(), UpdateProfileRequest{})

	assert.Error(t, err)
	assert.ErrorIs(t, err, upsertErr)
	assert.Contains(t, err.Error(), "ProfileService.UpdateProfile")
	assert.Nil(t, resp)
}

func TestProfileService_UpdateProfile_PartialFields_NoTDEE(t *testing.T) {
	// When not all TDEE-required fields are present, TDEE should NOT be calculated.
	userID := uuid.New()
	gender := "MALE"

	var upsertedProfile *UserProfile
	getCallCount := 0

	svc := NewProfileService(
		&mockProfileRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*UserProfile, error) {
				getCallCount++
				if getCallCount == 1 {
					return nil, errors.New("not found")
				}
				return upsertedProfile, nil
			},
			upsertFn: func(_ context.Context, p *UserProfile) error {
				upsertedProfile = p
				return nil
			},
		},
		&mockAllergyRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID) ([]UserAllergy, error) {
				return []UserAllergy{}, nil
			},
		},
		&mockTDEECalculator{
			// These should NOT be called since not all fields are present.
			calculateBMRFn: func(_ string, _, _ float64, _ int) float64 {
				t.Error("CalculateBMR should not be called with partial fields")
				return 0
			},
		},
	)

	resp, err := svc.UpdateProfile(context.Background(), userID, UpdateProfileRequest{
		Gender: &gender,
		// Age, HeightCm, WeightKg, ActivityLevel, Goal are all missing.
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// BMR/TDEE/DailyTarget should remain nil.
	assert.Nil(t, upsertedProfile.BMR)
	assert.Nil(t, upsertedProfile.TDEE)
	assert.Nil(t, upsertedProfile.DailyTarget)
}
