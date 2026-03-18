package user

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile holds a user's physical stats and nutrition targets.
type UserProfile struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	Gender        *string   `json:"gender,omitempty"`
	Age           *int      `json:"age,omitempty"`
	HeightCm      *float64  `json:"heightCm,omitempty"`
	WeightKg      *float64  `json:"weightKg,omitempty"`
	ActivityLevel *string   `json:"activityLevel,omitempty"`
	Goal          *string   `json:"goal,omitempty"`
	BMR           *float64  `json:"bmr,omitempty"`
	TDEE          *float64  `json:"tdee,omitempty"`
	DailyTarget   *float64  `json:"dailyTarget,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UserAllergy represents a user's allergy or dislike for a specific ingredient.
type UserAllergy struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"userId"`
	IngredientID   *uuid.UUID `json:"ingredientId,omitempty"`
	IngredientName string     `json:"ingredientName"`
	AllergyType    string     `json:"allergyType"` // ALLERGY or DISLIKE
}

// AuthTokens holds the access and refresh tokens returned by the auth provider.
type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// AuthUser represents the authenticated user identity from the auth provider.
type AuthUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
