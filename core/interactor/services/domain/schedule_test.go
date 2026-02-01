package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// DayOfWeek Tests
// ============================================

func TestDayOfWeek_DayName(t *testing.T) {
	tests := []struct {
		day      DayOfWeek
		expected string
	}{
		{Monday, "Lunes"},
		{Tuesday, "Martes"},
		{Wednesday, "Miércoles"},
		{Thursday, "Jueves"},
		{Friday, "Viernes"},
		{Saturday, "Sábado"},
		{Sunday, "Domingo"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.day.DayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDayOfWeek_IsValid(t *testing.T) {
	tests := []struct {
		day      DayOfWeek
		expected bool
	}{
		{Monday, true},
		{Tuesday, true},
		{Wednesday, true},
		{Thursday, true},
		{Friday, true},
		{Saturday, true},
		{Sunday, true},
		{DayOfWeek(0), false},
		{DayOfWeek(8), false},
		{DayOfWeek(-1), false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.day.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllDaysOfWeek(t *testing.T) {
	days := AllDaysOfWeek()

	assert.Len(t, days, 7)
	assert.Equal(t, Monday, days[0])
	assert.Equal(t, Sunday, days[6])
}

func TestGetDayCatalog(t *testing.T) {
	catalog := GetDayCatalog()

	assert.Len(t, catalog, 7)

	// Check first entry
	assert.Equal(t, 1, catalog[0].Value)
	assert.Equal(t, "Lunes", catalog[0].Name)

	// Check last entry
	assert.Equal(t, 7, catalog[6].Value)
	assert.Equal(t, "Domingo", catalog[6].Name)
}

// ============================================
// EntryType Tests
// ============================================

func TestEntryType_Constants(t *testing.T) {
	assert.Equal(t, EntryType("REGULAR"), EntryTypeRegular)
	assert.Equal(t, EntryType("EXCEPTION"), EntryTypeException)
}

// ============================================
// BranchSchedule Tests
// ============================================

func TestBranchSchedule_Structure(t *testing.T) {
	schedule := BranchSchedule{
		ID:       "schedule-001",
		BranchID: "branch-001",
		Active:   true,
	}

	assert.Equal(t, "schedule-001", schedule.ID)
	assert.Equal(t, "branch-001", schedule.BranchID)
	assert.True(t, schedule.Active)
	assert.Nil(t, schedule.EndDate)
}

// ============================================
// ScheduleDetail Tests
// ============================================

func TestScheduleDetail_Structure(t *testing.T) {
	dayOfWeek := 1
	openingTime := "08:00"
	closingTime := "17:00"

	detail := ScheduleDetail{
		ID:          "detail-001",
		ScheduleID:  "schedule-001",
		EntryType:   EntryTypeRegular,
		DayOfWeek:   &dayOfWeek,
		OpeningTime: &openingTime,
		ClosingTime: &closingTime,
		IsClosed:    false,
		Active:      true,
	}

	assert.Equal(t, "detail-001", detail.ID)
	assert.Equal(t, EntryTypeRegular, detail.EntryType)
	assert.NotNil(t, detail.DayOfWeek)
	assert.Equal(t, 1, *detail.DayOfWeek)
	assert.Equal(t, "08:00", *detail.OpeningTime)
	assert.False(t, detail.IsClosed)
	assert.True(t, detail.Active)
}

func TestScheduleDetail_ClosedDay(t *testing.T) {
	dayOfWeek := 7 // Sunday

	detail := ScheduleDetail{
		ID:         "detail-002",
		ScheduleID: "schedule-001",
		EntryType:  EntryTypeRegular,
		DayOfWeek:  &dayOfWeek,
		IsClosed:   true,
		Active:     true,
	}

	assert.True(t, detail.IsClosed)
	assert.Nil(t, detail.OpeningTime)
	assert.Nil(t, detail.ClosingTime)
}
