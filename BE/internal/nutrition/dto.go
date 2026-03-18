package nutrition

// CalculateTDEERequest is the payload for POST /nutrition/calculate-tdee.
type CalculateTDEERequest struct {
	Gender        string  `json:"gender" validate:"required,oneof=MALE FEMALE"`
	Age           int     `json:"age" validate:"required,min=1,max=150"`
	HeightCm      float64 `json:"heightCm" validate:"required,min=50,max=300"`
	WeightKg      float64 `json:"weightKg" validate:"required,min=10,max=500"`
	ActivityLevel string  `json:"activityLevel" validate:"required,oneof=SEDENTARY LIGHT MODERATE ACTIVE VERY_ACTIVE"`
	Goal          string  `json:"goal" validate:"required,oneof=LOSE_WEIGHT MAINTAIN GAIN_WEIGHT"`
}

// TDEEResponse is the response for the TDEE calculation endpoint.
type TDEEResponse struct {
	BMR           float64            `json:"bmr"`
	TDEE          float64            `json:"tdee"`
	DailyTarget   float64            `json:"dailyTarget"`
	MealBreakdown map[string]float64 `json:"mealBreakdown"`
}

// RecipeNutritionResponse is the response for recipe nutrition queries.
type RecipeNutritionResponse struct {
	ID          string   `json:"id"`
	RecipeID    string   `json:"recipeId"`
	Calories    *float64 `json:"calories,omitempty"`
	Protein     *float64 `json:"protein,omitempty"`
	Carbs       *float64 `json:"carbs,omitempty"`
	Fat         *float64 `json:"fat,omitempty"`
	Fiber       *float64 `json:"fiber,omitempty"`
	Sugar       *float64 `json:"sugar,omitempty"`
	Sodium      *float64 `json:"sodium,omitempty"`
	ServingSize *string  `json:"servingSize,omitempty"`
	DataSource  *string  `json:"dataSource,omitempty"`
	IsVerified  bool     `json:"isVerified"`
}

// UpsertNutritionRequest is the payload for POST /nutrition/recipes/:id (admin).
type UpsertNutritionRequest struct {
	Calories    *float64 `json:"calories" validate:"omitempty,min=0"`
	Protein     *float64 `json:"protein" validate:"omitempty,min=0"`
	Carbs       *float64 `json:"carbs" validate:"omitempty,min=0"`
	Fat         *float64 `json:"fat" validate:"omitempty,min=0"`
	Fiber       *float64 `json:"fiber" validate:"omitempty,min=0"`
	Sugar       *float64 `json:"sugar" validate:"omitempty,min=0"`
	Sodium      *float64 `json:"sodium" validate:"omitempty,min=0"`
	ServingSize *string  `json:"servingSize"`
	DataSource  *string  `json:"dataSource"`
	IsVerified  bool     `json:"isVerified"`
}
