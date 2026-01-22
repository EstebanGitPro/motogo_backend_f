package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// ============================================
// Schedule Detail DTOs (HU6-9)
// ============================================

// CreateScheduleDetailRequest is the DTO for creating a schedule detail (HU6)
type CreateScheduleDetailRequest struct {
	DayOfWeek   int     `json:"day_of_week"`            // 1=Monday to 7=Sunday
	OpeningTime *string `json:"opening_time,omitempty"` // HH:MM format, required if is_closed=false
	ClosingTime *string `json:"closing_time,omitempty"` // HH:MM format, required if is_closed=false
	IsClosed    bool    `json:"is_closed"`              // true if branch is closed on this day
}

// Sanitize trims whitespace from string fields
func (r *CreateScheduleDetailRequest) Sanitize() {
	r.OpeningTime = TrimStringPtr(r.OpeningTime)
	r.ClosingTime = TrimStringPtr(r.ClosingTime)
}

// UpdateScheduleDetailRequest is the DTO for updating a schedule detail (HU7)
type UpdateScheduleDetailRequest struct {
	DayOfWeek   *int    `json:"day_of_week,omitempty"`  // 1=Monday to 7=Sunday
	OpeningTime *string `json:"opening_time,omitempty"` // HH:MM format
	ClosingTime *string `json:"closing_time,omitempty"` // HH:MM format
	IsClosed    *bool   `json:"is_closed,omitempty"`    // true if branch is closed
	Active      *bool   `json:"active,omitempty"`       // true if detail is active
}

// Sanitize trims whitespace from string fields
func (r *UpdateScheduleDetailRequest) Sanitize() {
	r.OpeningTime = TrimStringPtr(r.OpeningTime)
	r.ClosingTime = TrimStringPtr(r.ClosingTime)
}

// ScheduleDetailResponse represents a schedule detail in API responses (HATEOAS)
type ScheduleDetailResponse struct {
	ID          string  `json:"id"`
	ScheduleID  string  `json:"schedule_id"`
	DayOfWeek   int     `json:"day_of_week"`
	DayName     string  `json:"day_name"`
	OpeningTime *string `json:"opening_time,omitempty"`
	ClosingTime *string `json:"closing_time,omitempty"`
	IsClosed    bool    `json:"is_closed"`
	Active      bool    `json:"active"`
	Links       []Link  `json:"_links"`
}

// NewScheduleDetailResponse creates a ScheduleDetailResponse from domain.ScheduleDetail
func NewScheduleDetailResponse(
	detail *domain.ScheduleDetail,
	encodedID, encodedScheduleID string,
	links []Link,
) ScheduleDetailResponse {
	dayOfWeek := 0
	dayName := ""
	if detail.DayOfWeek != nil {
		dayOfWeek = *detail.DayOfWeek
		dayName = domain.DayOfWeek(dayOfWeek).DayName()
	}

	return ScheduleDetailResponse{
		ID:          encodedID,
		ScheduleID:  encodedScheduleID,
		DayOfWeek:   dayOfWeek,
		DayName:     dayName,
		OpeningTime: detail.OpeningTime,
		ClosingTime: detail.ClosingTime,
		IsClosed:    detail.IsClosed,
		Active:      detail.Active,
		Links:       links,
	}
}

// ScheduleDetailListResponse represents a list of schedule details (HU9)
type ScheduleDetailListResponse struct {
	Details []ScheduleDetailResponse `json:"details"`
	Links   []Link                   `json:"_links"`
}

// NewScheduleDetailListResponse creates a list response
func NewScheduleDetailListResponse(details []ScheduleDetailResponse, links []Link) ScheduleDetailListResponse {
	return ScheduleDetailListResponse{
		Details: details,
		Links:   links,
	}
}

// ScheduleDetailDeleteResponse represents the response after deleting a detail
type ScheduleDetailDeleteResponse struct {
	Links []Link `json:"_links"`
}

// NewScheduleDetailDeleteResponse creates a response with next action links
func NewScheduleDetailDeleteResponse(baseURL, encodedBranchID string) ScheduleDetailDeleteResponse {
	return ScheduleDetailDeleteResponse{
		Links: []Link{
			{Rel: "schedule", Href: BuildResourceURL(baseURL, "branches", encodedBranchID) + "/schedules", Method: "GET"},
			{Rel: "details-list", Href: BuildResourceURL(baseURL, "branches", encodedBranchID) + "/schedules/details", Method: "GET"},
			{Rel: "create", Href: BuildResourceURL(baseURL, "branches", encodedBranchID) + "/schedules/details", Method: "POST"},
		},
	}
}
