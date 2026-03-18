package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// AuthService handles authentication operations by delegating to AuthProvider.
type AuthService struct {
	auth AuthProvider
}

// NewAuthService creates a new AuthService.
func NewAuthService(auth AuthProvider) *AuthService {
	return &AuthService{auth: auth}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	user, tokens, err := s.auth.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register: %w", err)
	}

	slog.Info("user registered", "user_id", user.ID, "email", user.Email)

	return &AuthResponse{User: *user, Tokens: *tokens}, nil
}

// Login authenticates an existing user.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, tokens, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login: %w", err)
	}

	slog.Info("user logged in", "user_id", user.ID)

	return &AuthResponse{User: *user, Tokens: *tokens}, nil
}

// RefreshToken exchanges a refresh token for new access tokens.
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshRequest) (*AuthTokens, error) {
	tokens, err := s.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("AuthService.RefreshToken: %w", err)
	}

	return tokens, nil
}

// ProfileService manages user profile and allergy data.
type ProfileService struct {
	profiles  ProfileRepository
	allergies AllergyRepository
	tdee      TDEECalculator
}

// NewProfileService creates a new ProfileService.
func NewProfileService(profiles ProfileRepository, allergies AllergyRepository, tdee TDEECalculator) *ProfileService {
	return &ProfileService{
		profiles:  profiles,
		allergies: allergies,
		tdee:      tdee,
	}
}

// GetProfile returns the full profile with allergies for a user.
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ProfileService.GetProfile: %w", err)
	}

	allergies, err := s.allergies.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ProfileService.GetProfile: %w", err)
	}

	allergyDTOs := make([]AllergyDTO, len(allergies))
	for i, a := range allergies {
		allergyDTOs[i] = AllergyDTO{
			IngredientName: a.IngredientName,
			AllergyType:    a.AllergyType,
		}
	}

	return &ProfileResponse{
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
		Gender:        profile.Gender,
		Age:           profile.Age,
		HeightCm:      profile.HeightCm,
		WeightKg:      profile.WeightKg,
		ActivityLevel: profile.ActivityLevel,
		Goal:          profile.Goal,
		BMR:           profile.BMR,
		TDEE:          profile.TDEE,
		DailyTarget:   profile.DailyTarget,
		Allergies:     allergyDTOs,
	}, nil
}

// UpdateProfile upserts the user profile, recalculates TDEE, and sets allergies.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*ProfileResponse, error) {
	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		// If profile doesn't exist, create a new one.
		profile = &UserProfile{
			ID:        uuid.New(),
			UserID:    userID,
			CreatedAt: time.Now().UTC(),
		}
	}

	// Update fields from request.
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.Age != nil {
		profile.Age = req.Age
	}
	if req.HeightCm != nil {
		profile.HeightCm = req.HeightCm
	}
	if req.WeightKg != nil {
		profile.WeightKg = req.WeightKg
	}
	if req.ActivityLevel != nil {
		profile.ActivityLevel = req.ActivityLevel
	}
	if req.Goal != nil {
		profile.Goal = req.Goal
	}
	profile.UpdatedAt = time.Now().UTC()

	// Recalculate TDEE if all required fields are present.
	if profile.Gender != nil && profile.WeightKg != nil && profile.HeightCm != nil && profile.Age != nil && profile.ActivityLevel != nil && profile.Goal != nil {
		bmr := s.tdee.CalculateBMR(*profile.Gender, *profile.WeightKg, *profile.HeightCm, *profile.Age)
		tdee := s.tdee.CalculateTDEE(bmr, *profile.ActivityLevel)
		dailyTarget := s.tdee.AdjustForGoal(tdee, *profile.Goal)

		profile.BMR = &bmr
		profile.TDEE = &tdee
		profile.DailyTarget = &dailyTarget
	}

	if err := s.profiles.Upsert(ctx, profile); err != nil {
		return nil, fmt.Errorf("ProfileService.UpdateProfile: %w", err)
	}

	// Set allergies.
	if req.Allergies != nil {
		allergies := make([]UserAllergy, len(req.Allergies))
		for i, a := range req.Allergies {
			allergies[i] = UserAllergy{
				ID:             uuid.New(),
				UserID:         userID,
				IngredientName: a.IngredientName,
				AllergyType:    a.AllergyType,
			}
		}
		if err := s.allergies.Set(ctx, userID, allergies); err != nil {
			return nil, fmt.Errorf("ProfileService.UpdateProfile: %w", err)
		}
	}

	slog.Info("user profile updated", "user_id", userID)

	return s.GetProfile(ctx, userID)
}
