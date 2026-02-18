package handlers

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
)

const (
	pathScheduleExceptions = "/schedule-exceptions/"
	pathBranches           = "/branches/"
	pathServices           = "/services"
)

// ==========================================
// Request DTOs
// ==========================================

// CreateScheduleExceptionRequest represents the request body for creating an exception
type CreateScheduleExceptionRequest struct {
	ExceptionStartDate string  `json:"exception_start_date" binding:"required"` // YYYY-MM-DD
	ExceptionEndDate   string  `json:"exception_end_date"`                      // YYYY-MM-DD (optional, defaults to start date)
	OpeningTime        *string `json:"opening_time"`                            // HH:mm
	ClosingTime        *string `json:"closing_time"`                            // HH:mm
	IsClosed           bool    `json:"is_closed"`
}

// Sanitize trims whitespace from all string fields
func (r *CreateScheduleExceptionRequest) Sanitize() {
	r.ExceptionStartDate = TrimString(r.ExceptionStartDate)
	r.ExceptionEndDate = TrimString(r.ExceptionEndDate)
	r.OpeningTime = TrimStringPtr(r.OpeningTime)
	r.ClosingTime = TrimStringPtr(r.ClosingTime)
}

// UpdateScheduleExceptionRequest represents the request body for updating an exception
type UpdateScheduleExceptionRequest struct {
	OpeningTime *string `json:"opening_time"` // HH:mm
	ClosingTime *string `json:"closing_time"` // HH:mm
	IsClosed    bool    `json:"is_closed"`
}

// Sanitize trims whitespace from all string fields
func (r *UpdateScheduleExceptionRequest) Sanitize() {
	r.OpeningTime = TrimStringPtr(r.OpeningTime)
	r.ClosingTime = TrimStringPtr(r.ClosingTime)
}

// ==========================================
// Response DTOs
// ==========================================

// ScheduleExceptionResponse represents the response for a schedule exception
type ScheduleExceptionResponse struct {
	ID                          string `json:"id"`
	ScheduleID                  string `json:"schedule_id"`
	ExceptionStartDate          string `json:"exception_start_date"`
	ExceptionEndDate            string `json:"exception_end_date"`
	ExceptionStartDateFormatted string `json:"exception_start_date_formatted"`
	ExceptionEndDateFormatted   string `json:"exception_end_date_formatted,omitempty"`
	DayName                     string `json:"day_name"`
	OpeningTime                 string `json:"opening_time,omitempty"`
	ClosingTime                 string `json:"closing_time,omitempty"`
	IsClosed                    bool   `json:"is_closed"`
	Active                      bool   `json:"active"`
	Links                       []Link `json:"_links"`
}

// ScheduleExceptionListResponse represents the list response
type ScheduleExceptionListResponse struct {
	Exceptions []ScheduleExceptionResponse `json:"exceptions"`
	Links      []Link                      `json:"_links"`
}

// ==========================================
// Converters & Helpers
// ==========================================

// NewScheduleExceptionResponse creates a response from domain entity
func NewScheduleExceptionResponse(
	exception *domain.ScheduleDetail,
	encodedExceptionID, encodedScheduleID string,
	links []Link,
) ScheduleExceptionResponse {
	response := ScheduleExceptionResponse{
		ID:         encodedExceptionID,
		ScheduleID: encodedScheduleID,
		IsClosed:   exception.IsClosed,
		Active:     exception.Active,
		Links:      links,
	}

	if exception.ExceptionStartDate != nil {
		response.ExceptionStartDate = exception.ExceptionStartDate.Format(constants.DateFormat)
		response.ExceptionStartDateFormatted = formatDateSpanish(*exception.ExceptionStartDate)
		response.DayName = getDayNameSpanish(*exception.ExceptionStartDate)
	}

	if exception.ExceptionEndDate != nil {
		response.ExceptionEndDate = exception.ExceptionEndDate.Format(constants.DateFormat)
		if !exception.ExceptionEndDate.Equal(*exception.ExceptionStartDate) {
			response.ExceptionEndDateFormatted = formatDateSpanish(*exception.ExceptionEndDate)
		}
	}

	if exception.OpeningTime != nil {
		response.OpeningTime = *exception.OpeningTime
	}
	if exception.ClosingTime != nil {
		response.ClosingTime = *exception.ClosingTime
	}

	return response
}

// formatDateSpanish formats a date in Spanish (e.g., "24 de Diciembre, 2026")
func formatDateSpanish(date time.Time) string {
	months := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}
	return date.Format("2") + " de " + months[date.Month()-1] + ", " + date.Format("2006")
}

// getDayNameSpanish returns the Spanish day name for a date
func getDayNameSpanish(date time.Time) string {
	days := []string{
		"Domingo", "Lunes", "Martes", "Miércoles",
		"Jueves", "Viernes", "Sábado",
	}
	return days[date.Weekday()]
}

// BuildScheduleExceptionLinks builds HATEOAS links for a schedule exception
func BuildScheduleExceptionLinks(baseURL, encodedBranchID, encodedExceptionID string) []Link {
	exceptionURL := baseURL + pathScheduleExceptions + encodedExceptionID
	return []Link{
		{Rel: "self", Href: exceptionURL, Method: "GET"},
		{Rel: "update", Href: exceptionURL, Method: "PUT"},
		{Rel: "delete", Href: exceptionURL, Method: "DELETE"},
		{Rel: "activate", Href: exceptionURL + "/activate", Method: "PUT"},
		{Rel: "deactivate", Href: exceptionURL + "/deactivate", Method: "PUT"},
		{Rel: "branch", Href: baseURL + pathBranches + encodedBranchID, Method: "GET"},
	}
}

// BuildScheduleExceptionListLinks builds HATEOAS links for exception list
func BuildScheduleExceptionListLinks(baseURL, encodedBranchID string) []Link {
	branchURL := baseURL + pathBranches + encodedBranchID
	return []Link{
		{Rel: "self", Href: branchURL + "/schedules/exceptions", Method: "GET"},
		{Rel: "create", Href: branchURL + "/schedules/exceptions", Method: "POST"},
		{Rel: "branch", Href: branchURL, Method: "GET"},
	}
}
