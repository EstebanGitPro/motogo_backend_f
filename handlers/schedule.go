package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
)

// ============================================
// Schedule DTOs (HU30-35, HU10)
// ============================================

// ScheduleResponse represents a schedule in API responses (HATEOAS)
type ScheduleResponse struct {
	ID        string  `json:"id"`
	BranchID  string  `json:"branch_id"`
	Active    bool    `json:"active"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"` // nullable - NULL means indefinite
	Links     []Link  `json:"_links"`
}

// NewScheduleResponse creates a ScheduleResponse from domain.BranchSchedule
func NewScheduleResponse(schedule *domain.BranchSchedule, encodedID, encodedBranchID string, links []Link) ScheduleResponse {
	startDate, endDate := formatScheduleDates(schedule)
	return ScheduleResponse{
		ID:        encodedID,
		BranchID:  encodedBranchID,
		Active:    schedule.Active,
		StartDate: startDate,
		EndDate:   endDate,
		Links:     links,
	}
}

// formatScheduleDates formats StartDate and EndDate for API response
func formatScheduleDates(schedule *domain.BranchSchedule) (string, *string) {
	startDate := schedule.StartDate.Format(constants.DateFormat)
	var endDate *string
	if schedule.EndDate != nil {
		formatted := schedule.EndDate.Format(constants.DateFormat)
		endDate = &formatted
	}
	return startDate, endDate
}

// DayOfWeekResponse represents a day of week for the catalog endpoint (HU10)
type DayOfWeekResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// DaysCatalogResponse represents the response for GET /schedules/days (HU10)
type DaysCatalogResponse struct {
	Days  []domain.DayCatalogEntry `json:"days"`
	Links []Link                   `json:"_links"`
}

// NewDaysCatalogResponse creates a DaysCatalogResponse
func NewDaysCatalogResponse(days []domain.DayCatalogEntry, baseURL string) DaysCatalogResponse {
	return DaysCatalogResponse{
		Days: days,
		Links: []Link{
			{Rel: "self", Href: baseURL + "/schedules/days", Method: "GET"},
		},
	}
}

// ScheduleDeleteResponse represents the response after deleting a schedule
type ScheduleDeleteResponse struct {
	Links []Link `json:"_links"`
}

// NewScheduleDeleteResponse creates a response with next action links
func NewScheduleDeleteResponse(baseURL, encodedBranchID string) ScheduleDeleteResponse {
	return ScheduleDeleteResponse{
		Links: []Link{
			{Rel: "create", Href: BuildResourceURL(baseURL, "branches", encodedBranchID) + "/schedules", Method: "POST"},
			{Rel: "branch", Href: BuildResourceURL(baseURL, "branches", encodedBranchID), Method: "GET"},
		},
	}
}

// ============================================
// Schedule Request DTOs (HU31)
// ============================================

// UpdateScheduleRequest is the DTO for updating schedule (HU31)
type UpdateScheduleRequest struct {
	Active    *bool   `json:"active,omitempty"`
	StartDate *string `json:"start_date,omitempty"` // YYYY-MM-DD format
	EndDate   *string `json:"end_date,omitempty"`   // YYYY-MM-DD format, null = indefinite
}

// Sanitize trims whitespace from string fields
func (r *UpdateScheduleRequest) Sanitize() {
	r.StartDate = TrimStringPtr(r.StartDate)
	r.EndDate = TrimStringPtr(r.EndDate)
}
