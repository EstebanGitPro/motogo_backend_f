package handlers

import (
	"strings"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// CreateDiagnosticRequest represents the request body for diagnostic creation (HU11)
type CreateDiagnosticRequest struct {
	BranchID           string   `json:"branch_id" binding:"required"`
	ProblemDescription *string  `json:"problem_description,omitempty"`
	EvidenceURLs       []string `json:"evidence_urls,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *CreateDiagnosticRequest) Sanitize() {
	r.BranchID = TrimString(r.BranchID)
	r.ProblemDescription = TrimStringPtr(r.ProblemDescription)
	for i := range r.EvidenceURLs {
		r.EvidenceURLs[i] = strings.TrimSpace(r.EvidenceURLs[i])
	}
}

// UpdateDiagnosticRequest represents the request body for diagnostic update (HU12)
// Note: possible_solution is NOT included here — only admin/workshop can set it via PATCH /admin/diagnostics/:id/solution
type UpdateDiagnosticRequest struct {
	ProblemDescription *string `json:"problem_description,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *UpdateDiagnosticRequest) Sanitize() {
	r.ProblemDescription = TrimStringPtr(r.ProblemDescription)
}

// ToDomain converts UpdateDiagnosticRequest to domain.Diagnostic
func (r *UpdateDiagnosticRequest) ToDomain() *domain.Diagnostic {
	return &domain.Diagnostic{
		ProblemDescription: r.ProblemDescription,
	}
}

// SetDiagnosticSolutionRequest represents the request body for admin setting diagnostic solution
type SetDiagnosticSolutionRequest struct {
	PossibleSolution *string `json:"possible_solution" binding:"required"`
}

// Sanitize trims whitespace from all string fields
func (r *SetDiagnosticSolutionRequest) Sanitize() {
	r.PossibleSolution = TrimStringPtr(r.PossibleSolution)
}

// DiagnosticResponse represents the API response for a diagnostic (HU11-14)
type DiagnosticResponse struct {
	ID                 string                       `json:"id"`
	MotorcycleID       string                       `json:"motorcycle_id"`
	BranchID           string                       `json:"branch_id"`
	Date               string                       `json:"date"`
	ProblemDescription *string                      `json:"problem_description,omitempty"`
	PossibleSolution   *string                      `json:"possible_solution,omitempty"`
	SentViaWhatsApp    bool                         `json:"sent_via_whatsapp"`
	Evidence           []DiagnosticEvidenceResponse `json:"evidence,omitempty"`
	Links              []Link                       `json:"_links,omitempty"`
}

// DiagnosticEvidenceResponse represents a diagnostic evidence photo in API responses
type DiagnosticEvidenceResponse struct {
	ID          string  `json:"id"`
	ImageURL    string  `json:"image_url"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// ToDiagnosticResponse converts domain.Diagnostic to DiagnosticResponse
func ToDiagnosticResponse(d *domain.Diagnostic) DiagnosticResponse {
	response := DiagnosticResponse{
		ID:                 d.ID,
		MotorcycleID:       d.MotorcycleID,
		BranchID:           d.BranchID,
		Date:               d.Date.Format("2006-01-02"),
		ProblemDescription: d.ProblemDescription,
		PossibleSolution:   d.PossibleSolution,
		SentViaWhatsApp:    d.SentViaWhatsApp,
	}

	if len(d.Evidence) > 0 {
		response.Evidence = make([]DiagnosticEvidenceResponse, len(d.Evidence))
		for i, e := range d.Evidence {
			response.Evidence[i] = DiagnosticEvidenceResponse{
				ID:          e.ID,
				ImageURL:    e.ImageURL,
				Description: e.Description,
				CreatedAt:   e.CreatedAt.Format(time.RFC3339),
			}
		}
	}

	return response
}

// ToDiagnosticResponseList converts a slice of domain.Diagnostic to []DiagnosticResponse
func ToDiagnosticResponseList(diagnostics []domain.Diagnostic) []DiagnosticResponse {
	responses := make([]DiagnosticResponse, len(diagnostics))
	for i := range diagnostics {
		responses[i] = ToDiagnosticResponse(&diagnostics[i])
	}
	return responses
}
