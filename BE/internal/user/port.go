package user

import (
	"context"

	"github.com/google/uuid"
)

// ProfileRepository manages persistence of user profiles.
type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	Upsert(ctx context.Context, profile *UserProfile) error
}

// AllergyRepository manages persistence of user allergies.
type AllergyRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]UserAllergy, error)
	Set(ctx context.Context, userID uuid.UUID, allergies []UserAllergy) error
}

// AuthProvider abstracts the external authentication service (e.g. Supabase Auth).
type AuthProvider interface {
	Register(ctx context.Context, email, password, name string) (*AuthUser, *AuthTokens, error)
	Login(ctx context.Context, email, password string) (*AuthUser, *AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
	GetUser(ctx context.Context, token string) (*AuthUser, error)
}

// TDEECalculator abstracts the nutrition TDEE calculation so the user context
// does not import the nutrition context directly.
type TDEECalculator interface {
	CalculateBMR(gender string, weightKg, heightCm float64, age int) float64
	CalculateTDEE(bmr float64, activityLevel string) float64
	AdjustForGoal(tdee float64, goal string) float64
}
