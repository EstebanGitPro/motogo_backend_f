package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
)

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
// IDs are encoded using the provided encoder for security
func NewBrandListResponse(brands []domain.Brand, links []Link, encoder *idencoder.HashidsEncoder) BrandListResponse {
	items := make([]BrandItemResponse, len(brands))
	for i, brand := range brands {
		encodedID, err := encoder.Encode(brand.ID)
		if err != nil {
			// Use original ID if encoding fails (shouldn't happen with valid UUIDs)
			encodedID = brand.ID
		}
		items[i] = BrandItemResponse{
			ID:   encodedID,
			Name: brand.Name,
		}
	}
	return BrandListResponse{
		Brands: items,
		Links:  links,
	}
}
