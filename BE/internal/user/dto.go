package user

// RegisterRequest is the payload for POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
}

// LoginRequest is the payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest is the payload for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// AuthResponse is the response for auth endpoints.
type AuthResponse struct {
	User   AuthUser   `json:"user"`
	Tokens AuthTokens `json:"tokens"`
}

// UpdateProfileRequest is the payload for PUT /users/me/profile.
type UpdateProfileRequest struct {
	Gender        *string      `json:"gender" validate:"omitempty,oneof=MALE FEMALE"`
	Age           *int         `json:"age" validate:"omitempty,min=1,max=150"`
	HeightCm      *float64     `json:"heightCm" validate:"omitempty,min=50,max=300"`
	WeightKg      *float64     `json:"weightKg" validate:"omitempty,min=10,max=500"`
	ActivityLevel *string      `json:"activityLevel" validate:"omitempty,oneof=SEDENTARY LIGHT MODERATE ACTIVE VERY_ACTIVE"`
	Goal          *string      `json:"goal" validate:"omitempty,oneof=LOSE_WEIGHT MAINTAIN GAIN_WEIGHT"`
	Allergies     []AllergyDTO `json:"allergies" validate:"omitempty,dive"`
}

// AllergyDTO represents a single allergy or dislike entry in a request.
type AllergyDTO struct {
	IngredientName string `json:"ingredientName" validate:"required,min=1,max=200"`
	AllergyType    string `json:"allergyType" validate:"required,oneof=ALLERGY DISLIKE"`
}

// ProfileResponse is the response for profile endpoints.
type ProfileResponse struct {
	ID            string       `json:"id"`
	UserID        string       `json:"userId"`
	Gender        *string      `json:"gender,omitempty"`
	Age           *int         `json:"age,omitempty"`
	HeightCm      *float64     `json:"heightCm,omitempty"`
	WeightKg      *float64     `json:"weightKg,omitempty"`
	ActivityLevel *string      `json:"activityLevel,omitempty"`
	Goal          *string      `json:"goal,omitempty"`
	BMR           *float64     `json:"bmr,omitempty"`
	TDEE          *float64     `json:"tdee,omitempty"`
	DailyTarget   *float64     `json:"dailyTarget,omitempty"`
	Allergies     []AllergyDTO `json:"allergies"`
}
