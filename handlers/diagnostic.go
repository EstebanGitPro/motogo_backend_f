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
type UpdateDiagnosticRequest struct {
	ProblemDescription *string  `json:"problem_description,omitempty"`
	PossibleSolution   *string  `json:"possible_solution,omitempty"`
	LaborQuote         *float64 `json:"labor_quote,omitempty"`
	PartsQuote         *float64 `json:"parts_quote,omitempty"`
	EstimatedTime      *string  `json:"estimated_time,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *UpdateDiagnosticRequest) Sanitize() {
	r.ProblemDescription = TrimStringPtr(r.ProblemDescription)
	r.PossibleSolution = TrimStringPtr(r.PossibleSolution)
	r.EstimatedTime = TrimStringPtr(r.EstimatedTime)
}

// ToDomain converts UpdateDiagnosticRequest to domain.Diagnostic
func (r *UpdateDiagnosticRequest) ToDomain() *domain.Diagnostic {
	return &domain.Diagnostic{
		ProblemDescription: r.ProblemDescription,
		PossibleSolution:   r.PossibleSolution,
		LaborQuote:         r.LaborQuote,
		PartsQuote:         r.PartsQuote,
		EstimatedTime:      r.EstimatedTime,
	}
}

// DiagnosticResponse represents the API response for a diagnostic (HU11-14)
type DiagnosticResponse struct {
	ID                 string                       `json:"id"`
	MotorcycleID       string                       `json:"motorcycle_id"`
	BranchID           string                       `json:"branch_id"`
	Date               string                       `json:"date"`
	ProblemDescription *string                      `json:"problem_description,omitempty"`
	PossibleSolution   *string                      `json:"possible_solution,omitempty"`
	LaborQuote         *float64                     `json:"labor_quote,omitempty"`
	PartsQuote         *float64                     `json:"parts_quote,omitempty"`
	EstimatedTime      *string                      `json:"estimated_time,omitempty"`
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
		Date:               d.Date.Format(time.RFC3339),
		ProblemDescription: d.ProblemDescription,
		PossibleSolution:   d.PossibleSolution,
		LaborQuote:         d.LaborQuote,
		PartsQuote:         d.PartsQuote,
		EstimatedTime:      d.EstimatedTime,
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
