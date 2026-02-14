package services_test

import (
	"context"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewMotorcycleService Tests
// ============================================

func TestNewMotorcycleService_ReturnsNonNil(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	assert.NotNil(t, svc)
}

// ============================================
// GetDistinctDisplacements Tests
// ============================================

func TestGetDistinctDisplacements_ReturnsThreeRanges(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, err := svc.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Len(t, displacements, 3)
}

func TestGetDistinctDisplacements_CorrectValues(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, err := svc.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, domain.DisplacementRangeLow, displacements[0].Range)
	assert.Equal(t, domain.DisplacementRangeMedium, displacements[1].Range)
	assert.Equal(t, domain.DisplacementRangeHigh, displacements[2].Range)
}

func TestGetDistinctDisplacements_CorrectOrder(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, _ := svc.GetDistinctDisplacements(context.Background())

	// Order should be BAJO, MEDIO, ALTO (ascending)
	assert.Equal(t, "BAJO", displacements[0].Range)
	assert.Equal(t, "MEDIO", displacements[1].Range)
	assert.Equal(t, "ALTO", displacements[2].Range)
}

// ============================================
// GetRatingRanges Tests
// ============================================

func TestGetRatingRanges_ReturnsFiveRanges(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, err := svc.GetRatingRanges(context.Background())

	assert.NoError(t, err)
	assert.Len(t, ratings, 5)
}

func TestGetRatingRanges_CorrectValues(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, err := svc.GetRatingRanges(context.Background())

	assert.NoError(t, err)

	expected := []struct {
		value int
		label string
	}{
		{1, "Very bad"},
		{2, "Bad"},
		{3, "Average"},
		{4, "Good"},
		{5, "Excellent"},
	}

	for i, exp := range expected {
		assert.Equal(t, exp.value, ratings[i].Value)
		assert.Equal(t, exp.label, ratings[i].Label)
	}
}

func TestGetRatingRanges_FirstIsOne(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, _ := svc.GetRatingRanges(context.Background())
	assert.Equal(t, 1, ratings[0].Value)
}

func TestGetRatingRanges_LastIsFive(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, _ := svc.GetRatingRanges(context.Background())
	assert.Equal(t, 5, ratings[len(ratings)-1].Value)
}
