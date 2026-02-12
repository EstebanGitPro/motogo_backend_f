package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// GrantPermissionRequest represents the request body for granting diagnostic permission
type GrantPermissionRequest struct {
	BranchID string `json:"branch_id" binding:"required"`
	Active   *bool  `json:"active"`
}

// Sanitize trims whitespace from all string fields
func (r *GrantPermissionRequest) Sanitize() {
	r.BranchID = TrimString(r.BranchID)
}

// DiagnosticPermissionResponse represents the API response for a diagnostic permission
type DiagnosticPermissionResponse struct {
	ID           string `json:"id"`
	MotorcycleID string `json:"motorcycle_id"`
	BranchID     string `json:"branch_id"`
	Active       bool   `json:"active"`
	Links        []Link `json:"_links,omitempty"`
}

// ToDiagnosticPermissionResponse converts domain to DTO
func ToDiagnosticPermissionResponse(p *domain.DiagnosticPermission) DiagnosticPermissionResponse {
	return DiagnosticPermissionResponse{
		ID:           p.ID,
		MotorcycleID: p.MotorcycleID,
		BranchID:     p.BranchID,
		Active:       p.Active,
	}
}

// ToDiagnosticPermissionResponseList converts domain slice to DTO slice
func ToDiagnosticPermissionResponseList(permissions []domain.DiagnosticPermission) []DiagnosticPermissionResponse {
	responses := make([]DiagnosticPermissionResponse, len(permissions))
	for i, p := range permissions {
		responses[i] = ToDiagnosticPermissionResponse(&p)
	}
	return responses
}
