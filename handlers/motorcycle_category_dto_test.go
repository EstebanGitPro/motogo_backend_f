package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewCategoryListResponse Tests
// ============================================

func TestNewCategoryListResponse_WithCategories(t *testing.T) {
	categories := []domain.MotorcycleCategory{
		{Name: "Sport", LineCount: 5},
		{Name: "Scooter", LineCount: 3},
	}

	resp := handlers.NewCategoryListResponse(categories, "http://localhost:8080")

	assert.Len(t, resp.Categories, 2)
	assert.Equal(t, "Sport", resp.Categories[0].Name)
	assert.Equal(t, 5, resp.Categories[0].LineCount)
	assert.Equal(t, "Scooter", resp.Categories[1].Name)
	assert.Equal(t, 3, resp.Categories[1].LineCount)
	assert.NotEmpty(t, resp.Categories[0].Links)
	assert.NotEmpty(t, resp.Links)
}

func TestNewCategoryListResponse_Empty(t *testing.T) {
	resp := handlers.NewCategoryListResponse([]domain.MotorcycleCategory{}, "http://localhost:8080")

	assert.Empty(t, resp.Categories)
	assert.NotEmpty(t, resp.Links) // Collection-level links should still be present
}

// ============================================
// NewCategoryLinesResponse Tests
// ============================================

func TestNewCategoryLinesResponse_WithLines(t *testing.T) {
	lines := []domain.CategoryLine{
		{Model: "CB 190R", BrandName: "Honda", EngineDisplacement: 190},
		{Model: "Gixxer 250", BrandName: "Suzuki", EngineDisplacement: 249},
	}

	resp := handlers.NewCategoryLinesResponse("Sport", lines, "http://localhost:8080")

	assert.Equal(t, "Sport", resp.Category)
	assert.Len(t, resp.Lines, 2)
	assert.Equal(t, "CB 190R", resp.Lines[0].Model)
	assert.Equal(t, "Honda", resp.Lines[0].Brand)
	assert.Equal(t, 190, resp.Lines[0].EngineDisplacement)
	assert.Equal(t, "Gixxer 250", resp.Lines[1].Model)
	assert.Equal(t, "Suzuki", resp.Lines[1].Brand)
	assert.NotEmpty(t, resp.Links)
}

func TestNewCategoryLinesResponse_Empty(t *testing.T) {
	resp := handlers.NewCategoryLinesResponse("Unknown", []domain.CategoryLine{}, "http://localhost:8080")

	assert.Equal(t, "Unknown", resp.Category)
	assert.Empty(t, resp.Lines)
	assert.NotEmpty(t, resp.Links)
}

// ============================================
// NewDisplacementListResponse Tests
// ============================================

func TestNewDisplacementListResponse_WithDisplacements(t *testing.T) {
	displacements := []domain.EngineDisplacementRange{
		{Range: domain.DisplacementRangeLow},
		{Range: domain.DisplacementRangeMedium},
		{Range: domain.DisplacementRangeHigh},
	}

	resp := handlers.NewDisplacementListResponse(displacements, "http://localhost:8080")

	assert.Len(t, resp.Displacements, 3)
	assert.Equal(t, "BAJO", resp.Displacements[0].Range)
	assert.Equal(t, "MEDIO", resp.Displacements[1].Range)
	assert.Equal(t, "ALTO", resp.Displacements[2].Range)
	assert.NotEmpty(t, resp.Links)
}

func TestNewDisplacementListResponse_Empty(t *testing.T) {
	resp := handlers.NewDisplacementListResponse([]domain.EngineDisplacementRange{}, "http://localhost:8080")

	assert.Empty(t, resp.Displacements)
	assert.NotEmpty(t, resp.Links)
}

// ============================================
// NewRatingRangeListResponse Tests
// ============================================

func TestNewRatingRangeListResponse_WithRatings(t *testing.T) {
	ranges := []domain.RatingRange{
		{Value: 1, Label: "Muy malo"},
		{Value: 2, Label: "Malo"},
		{Value: 3, Label: "Regular"},
		{Value: 4, Label: "Bueno"},
		{Value: 5, Label: "Excelente"},
	}

	resp := handlers.NewRatingRangeListResponse(ranges, "http://localhost:8080")

	assert.Len(t, resp.Ratings, 5)
	assert.Equal(t, 1, resp.Ratings[0].Value)
	assert.Equal(t, "Muy malo", resp.Ratings[0].Label)
	assert.Equal(t, 5, resp.Ratings[4].Value)
	assert.Equal(t, "Excelente", resp.Ratings[4].Label)
	assert.NotEmpty(t, resp.Links)
}

func TestNewRatingRangeListResponse_Empty(t *testing.T) {
	resp := handlers.NewRatingRangeListResponse([]domain.RatingRange{}, "http://localhost:8080")

	assert.Empty(t, resp.Ratings)
	assert.NotEmpty(t, resp.Links)
}
