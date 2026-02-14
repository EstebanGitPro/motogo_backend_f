package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// ============================================
// Motorcycle Category DTOs (HU41)
// ============================================

// CategoryItemResponse represents a single motorcycle category with HATEOAS drill-down links
type CategoryItemResponse struct {
	Name      string `json:"name"`
	LineCount int    `json:"line_count"`
	Links     []Link `json:"_links"`
}

// CategoryListResponse represents the list of motorcycle categories
type CategoryListResponse struct {
	Categories []CategoryItemResponse `json:"categories"`
	Links      []Link                 `json:"_links"`
}

// CategoryLineItemResponse represents a motorcycle line within a category
type CategoryLineItemResponse struct {
	Model              string `json:"model"`
	Brand              string `json:"brand"`
	EngineDisplacement int    `json:"engine_displacement"`
}

// CategoryLinesResponse represents the lines within a specific category
type CategoryLinesResponse struct {
	Category string                     `json:"category"`
	Lines    []CategoryLineItemResponse `json:"lines"`
	Links    []Link                     `json:"_links"`
}

// NewCategoryListResponse builds a CategoryListResponse from domain categories
func NewCategoryListResponse(categories []domain.MotorcycleCategory, baseURL string) CategoryListResponse {
	items := make([]CategoryItemResponse, len(categories))
	for i, cat := range categories {
		items[i] = CategoryItemResponse{
			Name:      cat.Name,
			LineCount: cat.LineCount,
			Links:     BuildMotorcycleCategoryItemLinks(baseURL, cat.Name),
		}
	}
	return CategoryListResponse{
		Categories: items,
		Links:      BuildMotorcycleCategoryListLinks(baseURL),
	}
}

// NewCategoryLinesResponse builds a CategoryLinesResponse from domain category lines
func NewCategoryLinesResponse(categoryName string, lines []domain.CategoryLine, baseURL string) CategoryLinesResponse {
	items := make([]CategoryLineItemResponse, len(lines))
	for i, line := range lines {
		items[i] = CategoryLineItemResponse{
			Model:              line.Model,
			Brand:              line.BrandName,
			EngineDisplacement: line.EngineDisplacement,
		}
	}
	return CategoryLinesResponse{
		Category: categoryName,
		Lines:    items,
		Links:    BuildCategoryLinesLinks(baseURL, categoryName),
	}
}

// ============================================
// Engine Displacement DTOs (HU49)
// ============================================

// DisplacementItemResponse represents a displacement range category
type DisplacementItemResponse struct {
	Range string `json:"range"`
}

// DisplacementListResponse represents the list of displacement range categories
type DisplacementListResponse struct {
	Displacements []DisplacementItemResponse `json:"displacements"`
	Links         []Link                     `json:"_links"`
}

// NewDisplacementListResponse builds a DisplacementListResponse from domain data
func NewDisplacementListResponse(displacements []domain.EngineDisplacementRange, baseURL string) DisplacementListResponse {
	items := make([]DisplacementItemResponse, len(displacements))
	for i, d := range displacements {
		items[i] = DisplacementItemResponse{
			Range: d.Range,
		}
	}
	return DisplacementListResponse{
		Displacements: items,
		Links:         BuildEngineDisplacementLinks(baseURL),
	}
}

// ============================================
// Rating Range DTOs (HU48)
// ============================================

// RatingRangeItemResponse represents a single rating value with its label
type RatingRangeItemResponse struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// RatingRangeListResponse represents the list of valid rating ranges
type RatingRangeListResponse struct {
	Ratings []RatingRangeItemResponse `json:"ratings"`
	Links   []Link                    `json:"_links"`
}

// NewRatingRangeListResponse builds a RatingRangeListResponse from domain data
func NewRatingRangeListResponse(ranges []domain.RatingRange, baseURL string) RatingRangeListResponse {
	items := make([]RatingRangeItemResponse, len(ranges))
	for i, r := range ranges {
		items[i] = RatingRangeItemResponse{
			Value: r.Value,
			Label: r.Label,
		}
	}
	return RatingRangeListResponse{
		Ratings: items,
		Links:   BuildRatingRangeLinks(baseURL),
	}
}
