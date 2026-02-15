package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// MotorcycleCategory Tests
// ============================================

func TestMotorcycleCategory_Structure(t *testing.T) {
	cat := MotorcycleCategory{
		Name:      "Sport",
		LineCount: 15,
	}

	assert.Equal(t, "Sport", cat.Name)
	assert.Equal(t, 15, cat.LineCount)
}

func TestMotorcycleCategory_ZeroLineCount(t *testing.T) {
	cat := MotorcycleCategory{Name: "Unknown"}
	assert.Equal(t, 0, cat.LineCount)
}

// ============================================
// CategoryLine Tests
// ============================================

func TestCategoryLine_Structure(t *testing.T) {
	line := CategoryLine{
		Model:              "MT-07",
		BrandName:          "Yamaha",
		EngineDisplacement: 689,
	}

	assert.Equal(t, "MT-07", line.Model)
	assert.Equal(t, "Yamaha", line.BrandName)
	assert.Equal(t, 689, line.EngineDisplacement)
}

func TestCategoryLine_ZeroDisplacement(t *testing.T) {
	line := CategoryLine{Model: "Unknown", BrandName: "Test"}
	assert.Equal(t, 0, line.EngineDisplacement)
}

// ============================================
// EngineDisplacementRange Tests
// ============================================

func TestEngineDisplacementRange_Structure(t *testing.T) {
	r := EngineDisplacementRange{Range: DisplacementRangeLow}
	assert.Equal(t, DisplacementRangeLow, r.Range)
}

func TestEngineDisplacementRange_AllValues(t *testing.T) {
	ranges := []EngineDisplacementRange{
		{Range: DisplacementRangeLow},
		{Range: DisplacementRangeMedium},
		{Range: DisplacementRangeHigh},
	}

	assert.Equal(t, DisplacementRangeLow, ranges[0].Range)
	assert.Equal(t, DisplacementRangeMedium, ranges[1].Range)
	assert.Equal(t, DisplacementRangeHigh, ranges[2].Range)
}

// ============================================
// RatingRange Tests
// ============================================

func TestRatingRange_Structure(t *testing.T) {
	r := RatingRange{Value: 5, Label: "Excellent"}
	assert.Equal(t, 5, r.Value)
	assert.Equal(t, "Excellent", r.Label)
}

func TestRatingRange_AllValues(t *testing.T) {
	ratings := []RatingRange{
		{Value: 1, Label: "Very bad"},
		{Value: 2, Label: "Bad"},
		{Value: 3, Label: "Average"},
		{Value: 4, Label: "Good"},
		{Value: 5, Label: "Excellent"},
	}

	assert.Len(t, ratings, 5)
	assert.Equal(t, 1, ratings[0].Value)
	assert.Equal(t, 5, ratings[4].Value)
	assert.Equal(t, "Excellent", ratings[4].Label)
}

// ============================================
// DisplacementRange Constants Tests
// ============================================

func TestDisplacementRangeConstants(t *testing.T) {
	assert.Equal(t, DisplacementRange("BAJO"), DisplacementRangeLow)
	assert.Equal(t, DisplacementRange("MEDIO"), DisplacementRangeMedium)
	assert.Equal(t, DisplacementRange("ALTO"), DisplacementRangeHigh)
}

func TestIsValidDisplacementRange_Valid(t *testing.T) {
	assert.True(t, IsValidDisplacementRange(DisplacementRange("BAJO")))
	assert.True(t, IsValidDisplacementRange(DisplacementRange("MEDIO")))
	assert.True(t, IsValidDisplacementRange(DisplacementRange("ALTO")))
}

func TestIsValidDisplacementRange_Invalid(t *testing.T) {
	assert.False(t, IsValidDisplacementRange(DisplacementRange("INVALID")))
	assert.False(t, IsValidDisplacementRange(DisplacementRange("")))
	assert.False(t, IsValidDisplacementRange(DisplacementRange("bajo")))
}

func TestValidateDisplacementRanges_AllValid(t *testing.T) {
	err := ValidateDisplacementRanges([]DisplacementRange{DisplacementRangeLow, DisplacementRangeMedium, DisplacementRangeHigh})
	assert.NoError(t, err)
}

func TestValidateDisplacementRanges_OneInvalid(t *testing.T) {
	err := ValidateDisplacementRanges([]DisplacementRange{DisplacementRangeLow, DisplacementRange("INVALID")})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidDisplacementRange, err)
}

func TestValidateDisplacementRanges_Empty(t *testing.T) {
	err := ValidateDisplacementRanges([]DisplacementRange{})
	assert.NoError(t, err)
}
