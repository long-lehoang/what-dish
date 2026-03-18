package nutrition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateBMR_Male(t *testing.T) {
	// Male, 70kg, 175cm, age 25
	// BMR = 10*70 + 6.25*175 - 5*25 + 5 = 700 + 1093.75 - 125 + 5 = 1673.75
	got := CalculateBMR("MALE", 70, 175, 25)
	assert.Equal(t, 1673.75, got)
}

func TestCalculateBMR_Female(t *testing.T) {
	// Female, 60kg, 165cm, age 30
	// BMR = 10*60 + 6.25*165 - 5*30 - 161 = 600 + 1031.25 - 150 - 161 = 1320.25
	got := CalculateBMR("FEMALE", 60, 165, 30)
	assert.Equal(t, 1320.25, got)
}

func TestCalculateBMR_NonMaleDefaultsToFemaleFormula(t *testing.T) {
	// Any gender that is not "MALE" should use the female formula (-161).
	got := CalculateBMR("OTHER", 60, 165, 30)
	expected := CalculateBMR("FEMALE", 60, 165, 30)
	assert.Equal(t, expected, got)
}

func TestCalculateTDEE_AllActivityLevels(t *testing.T) {
	bmr := 1500.0
	tests := []struct {
		name           string
		activityLevel  string
		wantMultiplier float64
	}{
		{"Sedentary", "SEDENTARY", 1.2},
		{"Light", "LIGHT", 1.375},
		{"Moderate", "MODERATE", 1.55},
		{"Active", "ACTIVE", 1.725},
		{"VeryActive", "VERY_ACTIVE", 1.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTDEE(bmr, tt.activityLevel)
			want := roundTwo(bmr * tt.wantMultiplier)
			assert.Equal(t, want, got)
		})
	}
}

func TestCalculateTDEE_UnknownActivityLevelDefaultsToSedentary(t *testing.T) {
	bmr := 1500.0
	got := CalculateTDEE(bmr, "UNKNOWN")
	want := CalculateTDEE(bmr, "SEDENTARY")
	assert.Equal(t, want, got)
}

func TestAdjustForGoal_AllGoals(t *testing.T) {
	tdee := 2000.0
	tests := []struct {
		name           string
		goal           string
		wantMultiplier float64
	}{
		{"LoseWeight", "LOSE_WEIGHT", 0.80},
		{"Maintain", "MAINTAIN", 1.00},
		{"GainWeight", "GAIN_WEIGHT", 1.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdjustForGoal(tdee, tt.goal)
			want := roundTwo(tdee * tt.wantMultiplier)
			assert.Equal(t, want, got)
		})
	}
}

func TestAdjustForGoal_UnknownGoalDefaultsToMaintain(t *testing.T) {
	tdee := 2000.0
	got := AdjustForGoal(tdee, "BULK")
	want := AdjustForGoal(tdee, "MAINTAIN")
	assert.Equal(t, want, got)
}

func TestCalculateMealBreakdown_SumsToTarget(t *testing.T) {
	target := 2000.0
	breakdown := CalculateMealBreakdown(target)

	assert.Equal(t, roundTwo(target*0.30), breakdown["breakfast"])
	assert.Equal(t, roundTwo(target*0.40), breakdown["lunch"])
	assert.Equal(t, roundTwo(target*0.30), breakdown["dinner"])

	sum := breakdown["breakfast"] + breakdown["lunch"] + breakdown["dinner"]
	assert.InDelta(t, target, sum, 0.01)
}

func TestCalculateMealBreakdown_Keys(t *testing.T) {
	breakdown := CalculateMealBreakdown(1800)
	assert.Contains(t, breakdown, "breakfast")
	assert.Contains(t, breakdown, "lunch")
	assert.Contains(t, breakdown, "dinner")
	assert.Len(t, breakdown, 3)
}

func TestTDEECalculatorImpl_SatisfiesInterface(t *testing.T) {
	// TDEECalculatorImpl should satisfy the user.TDEECalculator interface.
	// We verify by calling the same methods and checking results match the
	// package-level functions.
	calc := NewTDEECalculator()

	bmr := calc.CalculateBMR("MALE", 70, 175, 25)
	assert.Equal(t, CalculateBMR("MALE", 70, 175, 25), bmr)

	tdee := calc.CalculateTDEE(bmr, "MODERATE")
	assert.Equal(t, CalculateTDEE(bmr, "MODERATE"), tdee)

	adjusted := calc.AdjustForGoal(tdee, "LOSE_WEIGHT")
	assert.Equal(t, AdjustForGoal(tdee, "LOSE_WEIGHT"), adjusted)
}

// roundTwo rounds a float64 to 2 decimal places, matching the production code.
func roundTwo(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
