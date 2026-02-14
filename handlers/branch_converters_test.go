package handlers

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// toDisplacementRanges Tests
// ============================================

func TestToDisplacementRanges_WithValues(t *testing.T) {
	input := []string{"BAJO", "MEDIO", "ALTO"}

	result := toDisplacementRanges(input)

	assert.Len(t, result, 3)
	assert.Equal(t, domain.DisplacementRangeLow, result[0])
	assert.Equal(t, domain.DisplacementRangeMedium, result[1])
	assert.Equal(t, domain.DisplacementRangeHigh, result[2])
}

func TestToDisplacementRanges_Nil(t *testing.T) {
	result := toDisplacementRanges(nil)
	assert.Nil(t, result)
}

func TestToDisplacementRanges_Empty(t *testing.T) {
	result := toDisplacementRanges([]string{})
	assert.Len(t, result, 0)
}

// ============================================
// fromDisplacementRanges Tests
// ============================================

func TestFromDisplacementRanges_WithValues(t *testing.T) {
	input := []domain.DisplacementRange{
		domain.DisplacementRangeLow,
		domain.DisplacementRangeMedium,
		domain.DisplacementRangeHigh,
	}

	result := fromDisplacementRanges(input)

	assert.Len(t, result, 3)
	assert.Equal(t, "BAJO", result[0])
	assert.Equal(t, "MEDIO", result[1])
	assert.Equal(t, "ALTO", result[2])
}

func TestFromDisplacementRanges_Nil(t *testing.T) {
	result := fromDisplacementRanges(nil)
	assert.Nil(t, result)
}

func TestFromDisplacementRanges_Empty(t *testing.T) {
	result := fromDisplacementRanges([]domain.DisplacementRange{})
	assert.Len(t, result, 0)
}
