package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Helper functions
// ============================================

func intPtrDetail(i int) *int {
	return &i
}

func stringPtrDetail(s string) *string {
	return &s
}

// ============================================
// NewScheduleDetailResponse Tests
// ============================================

func TestNewScheduleDetailResponse_Success(t *testing.T) {
	// Arrange
	detail := &domain.ScheduleDetail{
		ID:          "detail-123",
		ScheduleID:  "schedule-123",
		DayOfWeek:   intPtrDetail(1), // Monday
		OpeningTime: stringPtrDetail("08:00"),
		ClosingTime: stringPtrDetail("18:00"),
		IsClosed:    false,
		Active:      true,
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/details/abc", Method: "GET"},
	}

	// Act
	response := handlers.NewScheduleDetailResponse(detail, "enc-detail", "enc-schedule", links)

	// Assert
	assert.Equal(t, "enc-detail", response.ID)
	assert.Equal(t, "enc-schedule", response.ScheduleID)
	assert.Equal(t, 1, response.DayOfWeek)
	assert.Equal(t, "Lunes", response.DayName)
	assert.Equal(t, "08:00", *response.OpeningTime)
	assert.Equal(t, "18:00", *response.ClosingTime)
	assert.False(t, response.IsClosed)
	assert.True(t, response.Active)
	assert.Len(t, response.Links, 1)
}

func TestNewScheduleDetailResponse_ClosedDay(t *testing.T) {
	// Arrange
	detail := &domain.ScheduleDetail{
		ID:         "detail-456",
		ScheduleID: "schedule-456",
		DayOfWeek:  intPtrDetail(7), // Sunday
		IsClosed:   true,
		Active:     true,
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewScheduleDetailResponse(detail, "enc-det", "enc-sch", links)

	// Assert
	assert.Equal(t, 7, response.DayOfWeek)
	assert.Equal(t, "Domingo", response.DayName)
	assert.Nil(t, response.OpeningTime)
	assert.Nil(t, response.ClosingTime)
	assert.True(t, response.IsClosed)
}

func TestNewScheduleDetailResponse_NilDayOfWeek(t *testing.T) {
	// Arrange - exception type detail (no day of week)
	detail := &domain.ScheduleDetail{
		ID:          "detail-789",
		ScheduleID:  "schedule-789",
		DayOfWeek:   nil, // No day - exception type
		OpeningTime: stringPtrDetail("09:00"),
		ClosingTime: stringPtrDetail("17:00"),
		IsClosed:    false,
		Active:      true,
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewScheduleDetailResponse(detail, "enc-d", "enc-s", links)

	// Assert
	assert.Equal(t, 0, response.DayOfWeek)
	assert.Equal(t, "", response.DayName) // No day name for exceptions
}

func TestNewScheduleDetailResponse_AllDays(t *testing.T) {
	// Test all day names are correct
	days := []struct {
		dayNum   int
		expected string
	}{
		{1, "Lunes"},
		{2, "Martes"},
		{3, "Miércoles"},
		{4, "Jueves"},
		{5, "Viernes"},
		{6, "Sábado"},
		{7, "Domingo"},
	}

	for _, tc := range days {
		detail := &domain.ScheduleDetail{
			ID:         "detail-test",
			ScheduleID: "schedule-test",
			DayOfWeek:  intPtrDetail(tc.dayNum),
		}

		response := handlers.NewScheduleDetailResponse(detail, "e1", "e2", nil)

		assert.Equal(t, tc.dayNum, response.DayOfWeek)
		assert.Equal(t, tc.expected, response.DayName, "Day %d should be %s", tc.dayNum, tc.expected)
	}
}

// ============================================
// NewScheduleDetailListResponse Tests
// ============================================

func TestNewScheduleDetailListResponse_Success(t *testing.T) {
	// Arrange
	details := []handlers.ScheduleDetailResponse{
		{ID: "enc-1", DayOfWeek: 1, DayName: "Lunes"},
		{ID: "enc-2", DayOfWeek: 2, DayName: "Martes"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/details", Method: "GET"},
	}

	// Act
	response := handlers.NewScheduleDetailListResponse(details, links)

	// Assert
	assert.Len(t, response.Details, 2)
	assert.Equal(t, "enc-1", response.Details[0].ID)
	assert.Equal(t, "enc-2", response.Details[1].ID)
	assert.Len(t, response.Links, 1)
}

func TestNewScheduleDetailListResponse_Empty(t *testing.T) {
	// Arrange
	details := []handlers.ScheduleDetailResponse{}
	links := []handlers.Link{
		{Rel: "self", Href: "/details", Method: "GET"},
	}

	// Act
	response := handlers.NewScheduleDetailListResponse(details, links)

	// Assert
	assert.Empty(t, response.Details)
	assert.Len(t, response.Links, 1)
}

// ============================================
// NewScheduleDetailDeleteResponse Tests
// ============================================

func TestNewScheduleDetailDeleteResponse_Success(t *testing.T) {
	// Arrange
	baseURL := "https://api.motogo.com"
	encodedBranchID := "xyz789"

	// Act
	response := handlers.NewScheduleDetailDeleteResponse(baseURL, encodedBranchID)

	// Assert
	assert.Len(t, response.Links, 3)
	assert.Equal(t, "schedule", response.Links[0].Rel)
	assert.Equal(t, "details-list", response.Links[1].Rel)
	assert.Equal(t, "create", response.Links[2].Rel)
}

// ============================================
// CreateScheduleDetailRequest.Sanitize Tests
// ============================================

func TestCreateScheduleDetailRequest_Sanitize(t *testing.T) {
	// Arrange
	openTime := "  08:00  "
	closeTime := "  18:00\t"
	req := &handlers.CreateScheduleDetailRequest{
		DayOfWeek:   1,
		OpeningTime: &openTime,
		ClosingTime: &closeTime,
		IsClosed:    false,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "08:00", *req.OpeningTime)
	assert.Equal(t, "18:00", *req.ClosingTime)
}

func TestCreateScheduleDetailRequest_Sanitize_NilTimes(t *testing.T) {
	// Arrange
	req := &handlers.CreateScheduleDetailRequest{
		DayOfWeek: 7,
		IsClosed:  true,
	}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Nil(t, req.OpeningTime)
	assert.Nil(t, req.ClosingTime)
}

// ============================================
// UpdateScheduleDetailRequest.Sanitize Tests
// ============================================

func TestUpdateScheduleDetailRequest_Sanitize(t *testing.T) {
	// Arrange
	openTime := "\n09:00\n"
	closeTime := " 17:00 "
	req := &handlers.UpdateScheduleDetailRequest{
		OpeningTime: &openTime,
		ClosingTime: &closeTime,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "09:00", *req.OpeningTime)
	assert.Equal(t, "17:00", *req.ClosingTime)
}

func TestUpdateScheduleDetailRequest_Sanitize_NilTimes(t *testing.T) {
	// Arrange
	dayOfWeek := 3
	req := &handlers.UpdateScheduleDetailRequest{
		DayOfWeek: &dayOfWeek,
	}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Nil(t, req.OpeningTime)
	assert.Nil(t, req.ClosingTime)
}
