package handlers_test

import (
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// CreateScheduleExceptionRequest.Sanitize Tests
// ============================================

func TestCreateScheduleExceptionRequest_Sanitize(t *testing.T) {
	// Arrange
	openingTime := "  09:00  "
	closingTime := "  18:00\t"
	req := &handlers.CreateScheduleExceptionRequest{
		ExceptionStartDate: "  2026-12-24  ",
		ExceptionEndDate:   "  2026-12-25\t",
		OpeningTime:        &openingTime,
		ClosingTime:        &closingTime,
		IsClosed:           false,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "2026-12-24", req.ExceptionStartDate)
	assert.Equal(t, "2026-12-25", req.ExceptionEndDate)
	assert.Equal(t, "09:00", *req.OpeningTime)
	assert.Equal(t, "18:00", *req.ClosingTime)
}

func TestCreateScheduleExceptionRequest_Sanitize_NilFields(t *testing.T) {
	// Arrange
	req := &handlers.CreateScheduleExceptionRequest{
		ExceptionStartDate: "2026-12-24",
		IsClosed:           true,
	}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Equal(t, "2026-12-24", req.ExceptionStartDate)
	assert.Empty(t, req.ExceptionEndDate)
	assert.Nil(t, req.OpeningTime)
	assert.Nil(t, req.ClosingTime)
}

// ============================================
// UpdateScheduleExceptionRequest.Sanitize Tests
// ============================================

func TestUpdateScheduleExceptionRequest_Sanitize(t *testing.T) {
	// Arrange
	openingTime := "  10:00  "
	closingTime := "  20:00\n"
	req := &handlers.UpdateScheduleExceptionRequest{
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
		IsClosed:    false,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "10:00", *req.OpeningTime)
	assert.Equal(t, "20:00", *req.ClosingTime)
}

func TestUpdateScheduleExceptionRequest_Sanitize_NilFields(t *testing.T) {
	// Arrange
	req := &handlers.UpdateScheduleExceptionRequest{
		IsClosed: true,
	}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Nil(t, req.OpeningTime)
	assert.Nil(t, req.ClosingTime)
	assert.True(t, req.IsClosed)
}

// ============================================
// NewScheduleExceptionResponse Tests
// ============================================

func TestNewScheduleExceptionResponse_FullDetails(t *testing.T) {
	// Arrange
	startDate := time.Date(2026, 12, 24, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local)
	openingTime := "09:00"
	closingTime := "14:00"

	exception := &domain.ScheduleDetail{
		ID:                 "exception-123",
		ScheduleID:         "schedule-456",
		ExceptionStartDate: &startDate,
		ExceptionEndDate:   &endDate,
		OpeningTime:        &openingTime,
		ClosingTime:        &closingTime,
		IsClosed:           false,
		Active:             true,
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/schedule-exceptions/abc", Method: "GET"},
	}

	// Act
	response := handlers.NewScheduleExceptionResponse(exception, "enc-exception", "enc-schedule", links)

	// Assert
	assert.Equal(t, "enc-exception", response.ID)
	assert.Equal(t, "enc-schedule", response.ScheduleID)
	assert.Equal(t, "2026-12-24", response.ExceptionStartDate)
	assert.Equal(t, "2026-12-25", response.ExceptionEndDate)
	assert.Contains(t, response.ExceptionStartDateFormatted, "24 de Diciembre")
	assert.Contains(t, response.ExceptionEndDateFormatted, "25 de Diciembre")
	assert.Equal(t, "Jueves", response.DayName) // December 24, 2026 is Thursday
	assert.Equal(t, "09:00", response.OpeningTime)
	assert.Equal(t, "14:00", response.ClosingTime)
	assert.False(t, response.IsClosed)
	assert.True(t, response.Active)
	assert.Len(t, response.Links, 1)
}

func TestNewScheduleExceptionResponse_ClosedDay(t *testing.T) {
	// Arrange
	date := time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local)

	exception := &domain.ScheduleDetail{
		ID:                 "exception-789",
		ScheduleID:         "schedule-456",
		ExceptionStartDate: &date,
		ExceptionEndDate:   &date, // Same date
		IsClosed:           true,
		Active:             true,
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewScheduleExceptionResponse(exception, "enc-ex", "enc-sch", links)

	// Assert
	assert.Equal(t, "enc-ex", response.ID)
	assert.Equal(t, "2026-12-25", response.ExceptionStartDate)
	assert.Equal(t, "2026-12-25", response.ExceptionEndDate)
	// End date formatted should be empty when same as start date
	assert.Empty(t, response.ExceptionEndDateFormatted)
	assert.True(t, response.IsClosed)
	assert.Empty(t, response.OpeningTime)
	assert.Empty(t, response.ClosingTime)
}

func TestNewScheduleExceptionResponse_NilDates(t *testing.T) {
	// Arrange
	exception := &domain.ScheduleDetail{
		ID:         "exception-nil",
		ScheduleID: "schedule-nil",
		IsClosed:   true,
		Active:     false,
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewScheduleExceptionResponse(exception, "enc-nil", "enc-sch", links)

	// Assert
	assert.Equal(t, "enc-nil", response.ID)
	assert.Empty(t, response.ExceptionStartDate)
	assert.Empty(t, response.ExceptionEndDate)
	assert.Empty(t, response.ExceptionStartDateFormatted)
	assert.Empty(t, response.DayName)
}

// ============================================
// BuildScheduleExceptionLinks Tests
// ============================================

func TestBuildScheduleExceptionLinks(t *testing.T) {
	// Act
	links := handlers.BuildScheduleExceptionLinks("http://api.test.com", "branch-enc", "exception-enc")

	// Assert — self GET was removed (no GET /schedule-exceptions/:id route exists)
	assert.Len(t, links, 5)

	// Verify all expected rels exist
	expectedRels := []string{"update", "delete", "activate", "deactivate", "branch"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}

func TestBuildScheduleExceptionLinks_URLs(t *testing.T) {
	// Act
	links := handlers.BuildScheduleExceptionLinks("http://api", "branch-1", "exc-1")

	// Assert - verify URLs are correctly constructed (no self-GET link)
	var updateLink, deleteLink, branchLink *handlers.Link
	for i, l := range links {
		switch l.Rel {
		case "update":
			updateLink = &links[i]
		case "delete":
			deleteLink = &links[i]
		case "branch":
			branchLink = &links[i]
		}
	}

	assert.NotNil(t, updateLink)
	assert.Equal(t, "http://api/schedule-exceptions/exc-1", updateLink.Href)
	assert.Equal(t, "PUT", updateLink.Method)

	assert.NotNil(t, deleteLink)
	assert.Equal(t, "http://api/schedule-exceptions/exc-1", deleteLink.Href)
	assert.Equal(t, "DELETE", deleteLink.Method)

	assert.NotNil(t, branchLink)
	assert.Equal(t, "http://api/branches/branch-1", branchLink.Href)
	assert.Equal(t, "GET", branchLink.Method)
}

func TestBuildScheduleExceptionLinks_ActivateDeactivate(t *testing.T) {
	// Act
	links := handlers.BuildScheduleExceptionLinks("http://api", "br", "ex")

	// Assert
	var activateLink, deactivateLink *handlers.Link
	for i, l := range links {
		if l.Rel == "activate" {
			activateLink = &links[i]
		}
		if l.Rel == "deactivate" {
			deactivateLink = &links[i]
		}
	}

	assert.NotNil(t, activateLink)
	assert.Equal(t, "http://api/schedule-exceptions/ex/activate", activateLink.Href)
	assert.Equal(t, "PUT", activateLink.Method)

	assert.NotNil(t, deactivateLink)
	assert.Equal(t, "http://api/schedule-exceptions/ex/deactivate", deactivateLink.Href)
	assert.Equal(t, "PUT", deactivateLink.Method)
}

// ============================================
// BuildScheduleExceptionListLinks Tests
// ============================================

func TestBuildScheduleExceptionListLinks(t *testing.T) {
	// Act
	links := handlers.BuildScheduleExceptionListLinks("http://api.test.com", "branch-enc")

	// Assert
	assert.Len(t, links, 3)

	expectedRels := []string{"self", "create", "branch"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}

func TestBuildScheduleExceptionListLinks_URLs(t *testing.T) {
	// Act
	links := handlers.BuildScheduleExceptionListLinks("http://api", "branch-123")

	// Assert
	var selfLink, createLink, branchLink *handlers.Link
	for i, l := range links {
		switch l.Rel {
		case "self":
			selfLink = &links[i]
		case "create":
			createLink = &links[i]
		case "branch":
			branchLink = &links[i]
		}
	}

	assert.NotNil(t, selfLink)
	assert.Equal(t, "http://api/branches/branch-123/schedules/exceptions", selfLink.Href)
	assert.Equal(t, "GET", selfLink.Method)

	assert.NotNil(t, createLink)
	assert.Equal(t, "http://api/branches/branch-123/schedules/exceptions", createLink.Href)
	assert.Equal(t, "POST", createLink.Method)

	assert.NotNil(t, branchLink)
	assert.Equal(t, "http://api/branches/branch-123", branchLink.Href)
	assert.Equal(t, "GET", branchLink.Method)
}
