package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// MotorcycleResponse represents the API response for a motorcycle (HU46)
type MotorcycleResponse struct {
	ID             string                  `json:"id"`
	LicensePlate   string                  `json:"license_plate"`
	Year           *int                    `json:"year,omitempty"`
	CurrentMileage *int                    `json:"current_mileage,omitempty"`
	OwnerNotes     *string                 `json:"owner_notes,omitempty"`
	Reference      *MotorcycleReferenceDTO `json:"reference,omitempty"`
	Links          []Link                  `json:"_links,omitempty"`
}

// MotorcycleReferenceDTO represents the motorcycle reference in API responses
type MotorcycleReferenceDTO struct {
	ID                 string `json:"id"`
	BrandID            string `json:"brand_id"`
	BrandName          string `json:"brand_name"`
	Model              string `json:"model"`
	Category           string `json:"category"`
	EngineDisplacement int    `json:"engine_displacement_cc"`
}

// ToMotorcycleResponse converts domain.Motorcycle to MotorcycleResponse
func ToMotorcycleResponse(m *domain.Motorcycle) MotorcycleResponse {
	response := MotorcycleResponse{
		ID:             m.ID,
		LicensePlate:   m.LicensePlate,
		Year:           m.Year,
		CurrentMileage: m.CurrentMileage,
		OwnerNotes:     m.OwnerNotes,
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

// RegisterMotorcycleRequest represents the request body for motorcycle registration (HU43)
// NOTE: reference_id is optional for Release 9 testing (catalog comes in Release 11)
type RegisterMotorcycleRequest struct {
	LicensePlate   string  `json:"license_plate" binding:"required"`
	ReferenceID    *string `json:"reference_id,omitempty"` // Optional until Release 11
	Year           *int    `json:"year,omitempty"`
	CurrentMileage *int    `json:"current_mileage,omitempty"`
	OwnerNotes     *string `json:"owner_notes,omitempty"`
}

// ToDomain converts RegisterMotorcycleRequest to domain.Motorcycle
func (r *RegisterMotorcycleRequest) ToDomain(ownerID string) *domain.Motorcycle {
	moto := &domain.Motorcycle{
		LicensePlate:   r.LicensePlate,
		OwnerID:        ownerID,
		Year:           r.Year,
		CurrentMileage: r.CurrentMileage,
		OwnerNotes:     r.OwnerNotes,
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
