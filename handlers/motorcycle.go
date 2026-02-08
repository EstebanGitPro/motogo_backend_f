package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// MotorcycleResponse represents the API response for a motorcycle (HU46)
type MotorcycleResponse struct {
	ID              string                  `json:"id"`
	LicensePlate    string                  `json:"license_plate"`
	Year            *int                    `json:"year,omitempty"`
	CurrentMileage  *int                    `json:"current_mileage,omitempty"`
	OwnerNotes      *string                 `json:"owner_notes,omitempty"`
	ProfileImageURL *string                 `json:"profile_image_url,omitempty"`
	Reference       *MotorcycleReferenceDTO `json:"reference,omitempty"`
	Links           []Link                  `json:"_links,omitempty"`
}

// MotorcycleReferenceDTO represents the motorcycle reference in API responses (owner view)
type MotorcycleReferenceDTO struct {
	ID                 string `json:"id"`
	BrandID            string `json:"brand_id"`
	BrandName          string `json:"brand_name"`
	Model              string `json:"model"`
	Category           string `json:"category"`
	EngineDisplacement int    `json:"engine_displacement_cc"`
}

// MotorcycleLookupResponse represents the API response for motorcycle lookup by workshops (HU47)
// This DTO excludes private owner data (notes) and unnecessary IDs (brand_id)
type MotorcycleLookupResponse struct {
	ID              string                        `json:"id"`
	LicensePlate    string                        `json:"license_plate"`
	Year            *int                          `json:"year,omitempty"`
	CurrentMileage  *int                          `json:"current_mileage,omitempty"`
	ProfileImageURL *string                       `json:"profile_image_url,omitempty"`
	Reference       *MotorcycleLookupReferenceDTO `json:"reference,omitempty"`
	Diagnostics     []DiagnosticResponse          `json:"diagnostics,omitempty"`
	Links           []Link                        `json:"_links,omitempty"`
}

// MotorcycleLookupReferenceDTO represents motorcycle reference for workshop lookup (no brand_id)
type MotorcycleLookupReferenceDTO struct {
	BrandName          string `json:"brand_name"`
	Model              string `json:"model"`
	Category           string `json:"category"`
	EngineDisplacement int    `json:"engine_displacement_cc"`
}

// ToMotorcycleResponse converts domain.Motorcycle to MotorcycleResponse
func ToMotorcycleResponse(m *domain.Motorcycle) MotorcycleResponse {
	response := MotorcycleResponse{
		ID:              m.ID,
		LicensePlate:    m.LicensePlate,
		Year:            m.Year,
		CurrentMileage:  m.CurrentMileage,
		OwnerNotes:      m.OwnerNotes,
		ProfileImageURL: m.ProfileImageURL,
	}

	if m.Reference != nil {
		response.Reference = &MotorcycleReferenceDTO{
			ID:                 m.Reference.ID,
			BrandID:            m.Reference.BrandID,
			BrandName:          m.Reference.BrandName,
			Model:              m.Reference.Model,
			Category:           m.Reference.Category,
			EngineDisplacement: m.Reference.EngineDisplacement,
		}
	}

	return response
}

// ToMotorcycleLookupResponse converts domain.Motorcycle to MotorcycleLookupResponse (workshop view)
// Excludes: owner_notes (private), brand_id (not needed for workshops)
func ToMotorcycleLookupResponse(m *domain.Motorcycle) MotorcycleLookupResponse {
	response := MotorcycleLookupResponse{
		ID:              m.ID,
		LicensePlate:    m.LicensePlate,
		Year:            m.Year,
		CurrentMileage:  m.CurrentMileage,
		ProfileImageURL: m.ProfileImageURL,
		// owner_notes intentionally excluded - private to owner
	}

	if m.Reference != nil {
		response.Reference = &MotorcycleLookupReferenceDTO{
			BrandName:          m.Reference.BrandName,
			Model:              m.Reference.Model,
			Category:           m.Reference.Category,
			EngineDisplacement: m.Reference.EngineDisplacement,
			// brand_id intentionally excluded - not needed for workshops
		}
	}

	return response
}

// RegisterMotorcycleRequest represents the request body for motorcycle registration (HU43)
// NOTE: reference_id is optional for Release 9 testing (catalog comes in Release 11)
type RegisterMotorcycleRequest struct {
	LicensePlate    string  `json:"license_plate" binding:"required"`
	ReferenceID     *string `json:"reference_id,omitempty"` // Optional until Release 11
	Year            *int    `json:"year,omitempty"`
	CurrentMileage  *int    `json:"current_mileage,omitempty"`
	OwnerNotes      *string `json:"owner_notes,omitempty"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"` // HU36: Main profile photo
}

// Sanitize trims whitespace from all string fields
func (r *RegisterMotorcycleRequest) Sanitize() {
	r.LicensePlate = TrimString(r.LicensePlate)
	r.ReferenceID = TrimStringPtr(r.ReferenceID)
	r.OwnerNotes = TrimStringPtr(r.OwnerNotes)
	r.ProfileImageURL = TrimStringPtr(r.ProfileImageURL)
}

// ToDomain converts RegisterMotorcycleRequest to domain.Motorcycle
func (r *RegisterMotorcycleRequest) ToDomain(ownerID string) *domain.Motorcycle {
	moto := &domain.Motorcycle{
		LicensePlate:    r.LicensePlate,
		OwnerID:         ownerID,
		Year:            r.Year,
		CurrentMileage:  r.CurrentMileage,
		OwnerNotes:      r.OwnerNotes,
		ProfileImageURL: r.ProfileImageURL,
	}
	if r.ReferenceID != nil {
		moto.ReferenceID = *r.ReferenceID
	}
	return moto
}

// ToMotorcycleResponseList converts a slice of domain.Motorcycle to []MotorcycleResponse
func ToMotorcycleResponseList(motorcycles []domain.Motorcycle) []MotorcycleResponse {
	responses := make([]MotorcycleResponse, len(motorcycles))
	for i, m := range motorcycles {
		responses[i] = ToMotorcycleResponse(&m)
	}
	return responses
}

// UpdateMotorcycleRequest represents the request body for motorcycle update (HU44)
// NOTE: license_plate is NOT updateable - it's a business identifier
type UpdateMotorcycleRequest struct {
	ReferenceID     *string `json:"reference_id,omitempty"`
	Year            *int    `json:"year,omitempty"`
	CurrentMileage  *int    `json:"current_mileage,omitempty"`
	OwnerNotes      *string `json:"owner_notes,omitempty"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"` // HU36: Main profile photo
}

// Sanitize trims whitespace from all string fields
func (r *UpdateMotorcycleRequest) Sanitize() {
	r.ReferenceID = TrimStringPtr(r.ReferenceID)
	r.OwnerNotes = TrimStringPtr(r.OwnerNotes)
	r.ProfileImageURL = TrimStringPtr(r.ProfileImageURL)
}

// ToDomain converts UpdateMotorcycleRequest to domain.Motorcycle
func (r *UpdateMotorcycleRequest) ToDomain() *domain.Motorcycle {
	moto := &domain.Motorcycle{
		Year:            r.Year,
		CurrentMileage:  r.CurrentMileage,
		OwnerNotes:      r.OwnerNotes,
		ProfileImageURL: r.ProfileImageURL,
	}
	if r.ReferenceID != nil {
		moto.ReferenceID = *r.ReferenceID
	}
	return moto
}

// MotorcycleReferenceCatalogItem represents a single reference in the catalog list (HU50)
type MotorcycleReferenceCatalogItem struct {
	ID                 string `json:"id"`
	BrandID            string `json:"brand_id"`
	BrandName          string `json:"brand_name"`
	Model              string `json:"model"`
	Category           string `json:"category"`
	EngineDisplacement int    `json:"engine_displacement_cc"`
}

// ToMotorcycleReferenceCatalogList converts domain references to catalog items
func ToMotorcycleReferenceCatalogList(refs []domain.MotorcycleReference) []MotorcycleReferenceCatalogItem {
	result := make([]MotorcycleReferenceCatalogItem, len(refs))
	for i, ref := range refs {
		result[i] = MotorcycleReferenceCatalogItem{
			ID:                 ref.ID,
			BrandID:            ref.BrandID,
			BrandName:          ref.BrandName,
			Model:              ref.Model,
			Category:           ref.Category,
			EngineDisplacement: ref.EngineDisplacement,
		}
	}
	return result
}

// BrandLineItem represents a simplified line item for brand lines endpoint (HU40)
// Only includes brand_name and model as requested
type BrandLineItem struct {
	BrandName string `json:"brand_name"`
	Model     string `json:"model"`
}
