package handlers

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// CreateDiagnosticRequest represents the request body for diagnostic creation (HU11)
type CreateDiagnosticRequest struct {
	BranchID           string  `json:"branch_id" binding:"required"`
	ProblemDescription *string `json:"problem_description,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *CreateDiagnosticRequest) Sanitize() {
	r.BranchID = TrimString(r.BranchID)
	r.ProblemDescription = TrimStringPtr(r.ProblemDescription)
}

// UpdateDiagnosticRequest represents the request body for diagnostic update (HU12)
type UpdateDiagnosticRequest struct {
	ProblemDescription *string `json:"problem_description,omitempty"`
	PossibleSolution   *string `json:"possible_solution,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *UpdateDiagnosticRequest) Sanitize() {
	r.ProblemDescription = TrimStringPtr(r.ProblemDescription)
	r.PossibleSolution = TrimStringPtr(r.PossibleSolution)
}

// ToDomain converts UpdateDiagnosticRequest to domain.Diagnostic
func (r *UpdateDiagnosticRequest) ToDomain() *domain.Diagnostic {
	return &domain.Diagnostic{
		ProblemDescription: r.ProblemDescription,
		PossibleSolution:   r.PossibleSolution,
	}
}

// SetSolutionRequest represents the request body for setting a diagnostic solution (representative)
type SetSolutionRequest struct {
	PossibleSolution string `json:"possible_solution" binding:"required"`
}

// Sanitize trims whitespace from all string fields
func (r *SetSolutionRequest) Sanitize() {
	r.PossibleSolution = TrimString(r.PossibleSolution)
}

// DiagnosticResponse represents the API response for a diagnostic (HU11-14)
type DiagnosticResponse struct {
	ID                 string  `json:"id"`
	MotorcycleID       string  `json:"motorcycle_id"`
	BranchID           string  `json:"branch_id"`
	BranchName         string  `json:"branch_name"`
	Date               string  `json:"date"`
	ProblemDescription *string `json:"problem_description,omitempty"`
	PossibleSolution   *string `json:"possible_solution,omitempty"`
	Links              []Link  `json:"_links,omitempty"`
}

// ToDiagnosticResponse converts domain.Diagnostic to DiagnosticResponse
func ToDiagnosticResponse(d *domain.Diagnostic) DiagnosticResponse {
	return DiagnosticResponse{
		ID:                 d.ID,
		MotorcycleID:       d.MotorcycleID,
		BranchID:           d.BranchID,
		BranchName:         d.BranchName,
		Date:               d.Date.Format(time.RFC3339),
		ProblemDescription: d.ProblemDescription,
		PossibleSolution:   d.PossibleSolution,
	}
}

// ToDiagnosticResponseList converts a slice of domain.Diagnostic to []DiagnosticResponse
func ToDiagnosticResponseList(diagnostics []domain.Diagnostic) []DiagnosticResponse {
	responses := make([]DiagnosticResponse, len(diagnostics))
	for i := range diagnostics {
		responses[i] = ToDiagnosticResponse(&diagnostics[i])
	}
	return responses
}
