package handlers

import (
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
)

// ==========================================
// Request DTOs
// ==========================================

// CreateCompletedServiceRequest represents the POST /completed-services request body.
type CreateCompletedServiceRequest struct {
	BranchID            string   `json:"branch_id" binding:"required"`
	MotorcycleID        string   `json:"motorcycle_id" binding:"required"`
	DiagnosticID        *string  `json:"diagnostic_id,omitempty"`
	ServiceIDs          []string `json:"service_ids" binding:"required,min=1"`
	QuotedPrice         *float64 `json:"quoted_price,omitempty"`
	FinalPrice          *float64 `json:"final_price,omitempty"`
	RepresentativeNotes *string  `json:"representative_notes,omitempty"`
}

// Sanitize trims whitespace from string fields.
func (r *CreateCompletedServiceRequest) Sanitize() {
	r.BranchID = strings.TrimSpace(r.BranchID)
	r.MotorcycleID = strings.TrimSpace(r.MotorcycleID)
	if r.DiagnosticID != nil {
		trimmed := strings.TrimSpace(*r.DiagnosticID)
		r.DiagnosticID = &trimmed
	}

	if r.RepresentativeNotes != nil {
		trimmed := strings.TrimSpace(*r.RepresentativeNotes)
		r.RepresentativeNotes = &trimmed
	}
	for i := range r.ServiceIDs {
		r.ServiceIDs[i] = strings.TrimSpace(r.ServiceIDs[i])
	}
}

// UpdateStatusRequest represents the request body for status updates
type UpdateStatusRequest struct {
	Status     string   `json:"status" binding:"required"`
	FinalPrice *float64 `json:"final_price,omitempty"`
}

// UpdateDetailsRequest represents the request body for updating service details (prices/notes)
type UpdateDetailsRequest struct {
	QuotedPrice         *float64 `json:"quoted_price,omitempty"`
	FinalPrice          *float64 `json:"final_price,omitempty"`
	RepresentativeNotes *string  `json:"representative_notes,omitempty"`
}

// Sanitize trims whitespace from string fields.
func (r *UpdateDetailsRequest) Sanitize() {
	if r.RepresentativeNotes != nil {
		trimmed := strings.TrimSpace(*r.RepresentativeNotes)
		r.RepresentativeNotes = &trimmed
	}
}

// ==========================================
// Response DTOs
// ==========================================

// CompletedServiceResponse represents the API response for a completed service.
type CompletedServiceResponse struct {
	ID                  string                         `json:"id"`
	BranchID            string                         `json:"branch_id"`
	BranchName          *string                        `json:"branch_name,omitempty"`
	MotorcycleID        string                         `json:"motorcycle_id"`
	DiagnosticID        *string                        `json:"diagnostic_id,omitempty"`
	Status              string                         `json:"status"`
	RequestDate         string                         `json:"request_date"`
	QuotedPrice         *float64                       `json:"quoted_price,omitempty"`
	FinalPrice          *float64                       `json:"final_price,omitempty"`
	RepresentativeNotes *string                        `json:"representative_notes,omitempty"`
	Services            []CompletedServiceItemResponse `json:"services,omitempty"`
}

// CompletedServiceItemResponse represents a single service item in the response.
type CompletedServiceItemResponse struct {
	ID          string  `json:"id"`
	ServiceID   string  `json:"service_id"`
	ServiceName *string `json:"service_name,omitempty"`
	Rating      *int    `json:"rating,omitempty"`
	Comment     *string `json:"comment,omitempty"`
	RatedAt     *string `json:"rated_at,omitempty"`
}

// StatusTransitionResponse represents a single status transition in the API response
type StatusTransitionResponse struct {
	ID             string  `json:"id"`
	PreviousStatus *string `json:"previous_status"`
	NewStatus      string  `json:"new_status"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
}

// ServiceStatusItem represents a single status option for the frontend
type ServiceStatusItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ==========================================
// Converters
// ==========================================

// ToCompletedServiceResponse converts a domain CompletedService to an API response.
func ToCompletedServiceResponse(cs *domain.CompletedService) CompletedServiceResponse {
	resp := CompletedServiceResponse{
		ID:                  cs.ID,
		BranchID:            cs.BranchID,
		BranchName:          cs.BranchName,
		MotorcycleID:        cs.MotorcycleID,
		DiagnosticID:        cs.DiagnosticID,
		Status:              string(cs.Status),
		RequestDate:         cs.RequestDate.Format(constants.DateFormat),
		QuotedPrice:         cs.QuotedPrice,
		FinalPrice:          cs.FinalPrice,
		RepresentativeNotes: cs.RepresentativeNotes,
	}

	if cs.Services != nil {
		resp.Services = make([]CompletedServiceItemResponse, len(cs.Services))
		for i, item := range cs.Services {
			itemResp := CompletedServiceItemResponse{
				ID:          item.ID,
				ServiceID:   item.ServiceID,
				ServiceName: item.ServiceName,
				Rating:      item.Rating,
				Comment:     item.Comment,
			}
			if item.RatedAt != nil {
				formatted := item.RatedAt.Format(constants.DateFormat)
				itemResp.RatedAt = &formatted
			}
			resp.Services[i] = itemResp
		}
	}
	return resp
}
