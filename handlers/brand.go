package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// BrandItemResponse represents a single brand in the catalog
type BrandItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BrandListResponse represents the response for GET /brands
type BrandListResponse struct {
	Brands []BrandItemResponse `json:"brands"`
	Links  []Link              `json:"_links"`
}

// NewBrandListResponse creates a BrandListResponse from domain brands
func NewBrandListResponse(brands []domain.Brand, links []Link) BrandListResponse {
	items := make([]BrandItemResponse, len(brands))
	for i, brand := range brands {
		items[i] = BrandItemResponse{
			ID:   brand.ID,
			Name: brand.Name,
		}
	}
	return BrandListResponse{
		Brands: items,
		Links:  links,
	}
}
