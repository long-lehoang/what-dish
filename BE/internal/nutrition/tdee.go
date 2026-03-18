package nutrition

import "math"

// Activity level multipliers (Mifflin-St Jeor).
var activityMultipliers = map[string]float64{
	"SEDENTARY":   1.2,
	"LIGHT":       1.375,
	"MODERATE":    1.55,
	"ACTIVE":      1.725,
	"VERY_ACTIVE": 1.9,
}

// Goal adjustment factors.
var goalAdjustments = map[string]float64{
	"LOSE_WEIGHT": 0.80,
	"MAINTAIN":    1.00,
	"GAIN_WEIGHT": 1.15,
}

// CalculateBMR calculates Basal Metabolic Rate using the Mifflin-St Jeor equation.
//
//	Male:   BMR = 10*weight(kg) + 6.25*height(cm) - 5*age + 5
//	Female: BMR = 10*weight(kg) + 6.25*height(cm) - 5*age - 161
func CalculateBMR(gender string, weightKg, heightCm float64, age int) float64 {
	bmr := 10*weightKg + 6.25*heightCm - 5*float64(age)
	if gender == "MALE" {
		bmr += 5
	} else {
		bmr -= 161
	}
	return math.Round(bmr*100) / 100
}

// CalculateTDEE calculates Total Daily Energy Expenditure from BMR and activity level.
func CalculateTDEE(bmr float64, activityLevel string) float64 {
	multiplier, ok := activityMultipliers[activityLevel]
	if !ok {
		multiplier = 1.2 // default to sedentary
	}
	return math.Round(bmr*multiplier*100) / 100
}

// AdjustForGoal adjusts TDEE based on the user's weight goal.
func AdjustForGoal(tdee float64, goal string) float64 {
	adjustment, ok := goalAdjustments[goal]
	if !ok {
		adjustment = 1.0 // default to maintain
	}
	return math.Round(tdee*adjustment*100) / 100
}

// CalculateMealBreakdown splits daily calorie target across meals.
// Breakfast: 30%, Lunch: 40%, Dinner: 30%.
func CalculateMealBreakdown(dailyTarget float64) map[string]float64 {
	return map[string]float64{
		"breakfast": math.Round(dailyTarget*0.30*100) / 100,
		"lunch":     math.Round(dailyTarget*0.40*100) / 100,
		"dinner":    math.Round(dailyTarget*0.30*100) / 100,
	}
}

// TDEECalculatorImpl satisfies the user.TDEECalculator interface so the
// user context can calculate TDEE without importing the nutrition package.
type TDEECalculatorImpl struct{}

// NewTDEECalculator creates a TDEECalculatorImpl.
func NewTDEECalculator() *TDEECalculatorImpl {
	return &TDEECalculatorImpl{}
}

func (t *TDEECalculatorImpl) CalculateBMR(gender string, weightKg, heightCm float64, age int) float64 {
	return CalculateBMR(gender, weightKg, heightCm, age)
}

func (t *TDEECalculatorImpl) CalculateTDEE(bmr float64, activityLevel string) float64 {
	return CalculateTDEE(bmr, activityLevel)
}

func (t *TDEECalculatorImpl) AdjustForGoal(tdee float64, goal string) float64 {
	return AdjustForGoal(tdee, goal)
}
