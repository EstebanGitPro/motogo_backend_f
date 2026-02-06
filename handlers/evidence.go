package handlers

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// CreateEvidenceRequest represents the request body for evidence creation (HU16)
type CreateEvidenceRequest struct {
	ImageURL    string  `json:"image_url" binding:"required"`
	Angle       *string `json:"angle,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ToDomain converts CreateEvidenceRequest to domain.MotorcycleEvidence
func (r *CreateEvidenceRequest) ToDomain() *domain.MotorcycleEvidence {
	return &domain.MotorcycleEvidence{
		ImageURL:    r.ImageURL,
		Angle:       r.Angle,
		Description: r.Description,
	}
}

// EvidenceResponse represents the API response for motorcycle evidence (HU16-19)
type EvidenceResponse struct {
	ID           string  `json:"id"`
	MotorcycleID string  `json:"motorcycle_id"`
	Angle        *string `json:"angle,omitempty"`
	ImageURL     string  `json:"image_url"`
	Description  *string `json:"description,omitempty"`
	CreatedAt    string  `json:"created_at"`
	Links        []Link  `json:"_links,omitempty"`
}

// ToEvidenceResponse converts domain.MotorcycleEvidence to EvidenceResponse
func ToEvidenceResponse(e *domain.MotorcycleEvidence) EvidenceResponse {
	return EvidenceResponse{
		ID:           e.ID,
		MotorcycleID: e.MotorcycleID,
		Angle:        e.Angle,
		ImageURL:     e.ImageURL,
		Description:  e.Description,
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
	}
}

// ToEvidenceResponseList converts a slice of domain.MotorcycleEvidence to []EvidenceResponse
func ToEvidenceResponseList(evidences []domain.MotorcycleEvidence) []EvidenceResponse {
	responses := make([]EvidenceResponse, len(evidences))
	for i, e := range evidences {
		responses[i] = ToEvidenceResponse(&e)
	}
	return responses
}
