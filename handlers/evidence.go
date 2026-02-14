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

// Sanitize trims whitespace from all string fields
func (r *CreateEvidenceRequest) Sanitize() {
	r.ImageURL = TrimString(r.ImageURL)
	r.Angle = TrimStringPtr(r.Angle)
	r.Description = TrimStringPtr(r.Description)
}

// ToDomain converts CreateEvidenceRequest to domain.MotorcycleEvidence
func (r *CreateEvidenceRequest) ToDomain() *domain.MotorcycleEvidence {
	e := &domain.MotorcycleEvidence{
		ImageURL:    r.ImageURL,
		Description: r.Description,
	}
	if r.Angle != nil {
		angle := domain.EvidenceAngle(*r.Angle)
		e.Angle = &angle
	}
	return e
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
	resp := EvidenceResponse{
		ID:           e.ID,
		MotorcycleID: e.MotorcycleID,
		ImageURL:     e.ImageURL,
		Description:  e.Description,
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
	}
	if e.Angle != nil {
		angle := string(*e.Angle)
		resp.Angle = &angle
	}
	return resp
}

// ToEvidenceResponseList converts a slice of domain.MotorcycleEvidence to []EvidenceResponse
func ToEvidenceResponseList(evidences []domain.MotorcycleEvidence) []EvidenceResponse {
	responses := make([]EvidenceResponse, len(evidences))
	for i, e := range evidences {
		responses[i] = ToEvidenceResponse(&e)
	}
	return responses
}
