package handlers_test

import (
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewScheduleResponse Tests
// ============================================

func TestNewScheduleResponse_WithEndDate(t *testing.T) {
	// Arrange
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	schedule := &domain.BranchSchedule{
		ID:        "schedule-123",
		BranchID:  "branch-123",
		Active:    true,
		StartDate: startDate,
		EndDate:   &endDate,
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/schedules/abc123", Method: "GET"},
	}

	// Act
	response := handlers.NewScheduleResponse(schedule, "encoded-sched-id", "encoded-branch-id", links)

	// Assert
	assert.Equal(t, "encoded-sched-id", response.ID)
	assert.Equal(t, "encoded-branch-id", response.BranchID)
	assert.True(t, response.Active)
	assert.Equal(t, "2026-01-01", response.StartDate)
	assert.NotNil(t, response.EndDate)
	assert.Equal(t, "2026-12-31", *response.EndDate)
	assert.Len(t, response.Links, 1)
}

func TestNewScheduleResponse_WithoutEndDate(t *testing.T) {
	// Arrange - indefinite schedule
	startDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	schedule := &domain.BranchSchedule{
		ID:        "schedule-456",
		BranchID:  "branch-456",
		Active:    false,
		StartDate: startDate,
		EndDate:   nil, // Indefinite
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewScheduleResponse(schedule, "enc-sched", "enc-branch", links)

	// Assert
	assert.Equal(t, "enc-sched", response.ID)
	assert.Equal(t, "enc-branch", response.BranchID)
	assert.False(t, response.Active)
	assert.Equal(t, "2026-06-15", response.StartDate)
	assert.Nil(t, response.EndDate)
	assert.Empty(t, response.Links)
}

// ============================================
// NewDaysCatalogResponse Tests (HU10)
// ============================================

func TestNewDaysCatalogResponse_Success(t *testing.T) {
	// Arrange
	days := domain.GetDayCatalog()
	baseURL := "https://api.motogo.com"

	// Act
	response := handlers.NewDaysCatalogResponse(days, baseURL)

	// Assert
	assert.Len(t, response.Days, 7) // 7 days of the week
	assert.Equal(t, 1, response.Days[0].Value)
	assert.Equal(t, "Lunes", response.Days[0].Name)
	assert.Equal(t, 7, response.Days[6].Value)
	assert.Equal(t, "Domingo", response.Days[6].Name)
	assert.Len(t, response.Links, 1)
	assert.Equal(t, "self", response.Links[0].Rel)
	assert.Equal(t, "https://api.motogo.com/schedules/days", response.Links[0].Href)
}

func TestNewDaysCatalogResponse_EmptyBaseURL(t *testing.T) {
	// Arrange
	days := domain.GetDayCatalog()

	// Act
	response := handlers.NewDaysCatalogResponse(days, "")

	// Assert
	assert.Len(t, response.Days, 7)
	assert.Equal(t, "/schedules/days", response.Links[0].Href)
}

// ============================================
// NewScheduleDeleteResponse Tests
// ============================================

func TestNewScheduleDeleteResponse_Success(t *testing.T) {
	// Arrange
	baseURL := "https://api.motogo.com"
	encodedBranchID := "xyz789"

	// Act
	response := handlers.NewScheduleDeleteResponse(baseURL, encodedBranchID)

	// Assert
	assert.Len(t, response.Links, 2)
	assert.Equal(t, "create", response.Links[0].Rel)
	assert.Contains(t, response.Links[0].Href, "/branches/xyz789/schedules")
	assert.Equal(t, "POST", response.Links[0].Method)
	assert.Equal(t, "branch", response.Links[1].Rel)
	assert.Contains(t, response.Links[1].Href, "/branches/xyz789")
	assert.Equal(t, "GET", response.Links[1].Method)
}

// ============================================
// UpdateScheduleRequest.Sanitize Tests
// ============================================

func TestUpdateScheduleRequest_Sanitize(t *testing.T) {
	// Arrange
	startDate := "  2026-01-01  "
	endDate := "  2026-12-31\t"
	req := &handlers.UpdateScheduleRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "2026-01-01", *req.StartDate)
	assert.Equal(t, "2026-12-31", *req.EndDate)
}

func TestUpdateScheduleRequest_Sanitize_NilFields(t *testing.T) {
	// Arrange
	req := &handlers.UpdateScheduleRequest{}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Nil(t, req.StartDate)
	assert.Nil(t, req.EndDate)
}

// ============================================
// BuildScheduleLinks Tests
// ============================================

func TestBuildScheduleLinks(t *testing.T) {
	// Arrange
	baseURL := "http://api.test.com"
	encodedBranchID := "branch-enc"
	encodedScheduleID := "schedule-enc"

	// Act
	links := handlers.BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID)

	// Assert
	assert.GreaterOrEqual(t, len(links), 4) // At least: self, delete, activate, deactivate

	// Find self link
	var selfLink *handlers.Link
	for i, l := range links {
		if l.Rel == "self" {
			selfLink = &links[i]
			break
		}
	}
	assert.NotNil(t, selfLink)
	assert.Equal(t, "GET", selfLink.Method)
	assert.Contains(t, selfLink.Href, "/branches/branch-enc/schedules")
}

func TestBuildScheduleLinks_ContainsExpectedRels(t *testing.T) {
	// Arrange
	links := handlers.BuildScheduleLinks("http://api", "branch-1", "sched-1")

	// Act - check all expected rels exist
	expectedRels := []string{"self", "delete", "activate", "deactivate", "details-list", "details-create", "branch", "days-catalog"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	// Assert
	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}

// ============================================
// BuildScheduleDetailLinks Tests
// ============================================

func TestBuildScheduleDetailLinks(t *testing.T) {
	// Act
	links := handlers.BuildScheduleDetailLinks("http://api", "branch-1", "detail-1")

	// Assert — self GET was removed (no GET /schedule-details/:id route exists)
	assert.GreaterOrEqual(t, len(links), 3)

	expectedRels := []string{"update", "delete", "schedule", "details-list"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}

// ============================================
// BuildScheduleDetailListLinks Tests
// ============================================

func TestBuildScheduleDetailListLinks(t *testing.T) {
	// Act
	links := handlers.BuildScheduleDetailListLinks("http://api", "branch-1")

	// Assert
	assert.GreaterOrEqual(t, len(links), 4)

	expectedRels := []string{"self", "create", "schedule", "branch", "days-catalog"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}
