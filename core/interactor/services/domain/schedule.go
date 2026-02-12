package domain

import (
	"time"

	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// DayOfWeek represents a day of the week (1=Monday to 7=Sunday)
type DayOfWeek int

const (
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
	Sunday    DayOfWeek = 7
)

// DayName returns the Spanish name for the day (used in API responses)
func (d DayOfWeek) DayName() string {
	names := map[DayOfWeek]string{
		Monday:    "Lunes",
		Tuesday:   "Martes",
		Wednesday: "Miércoles",
		Thursday:  "Jueves",
		Friday:    "Viernes",
		Saturday:  "Sábado",
		Sunday:    "Domingo",
	}
	return names[d]
}

// IsValid checks if the DayOfWeek value is valid (1-7)
func (d DayOfWeek) IsValid() bool {
	return d >= Monday && d <= Sunday
}

// AllDaysOfWeek returns all days of the week (for HU10 - Day Catalog)
func AllDaysOfWeek() []DayOfWeek {
	return []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}
}

// DayCatalogEntry represents a day in the catalog (HU10)
type DayCatalogEntry struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

// GetDayCatalog returns the full day catalog for API responses
func GetDayCatalog() []DayCatalogEntry {
	days := AllDaysOfWeek()
	catalog := make([]DayCatalogEntry, len(days))
	for i, d := range days {
		catalog[i] = DayCatalogEntry{
			Value: int(d),
			Name:  d.DayName(),
		}
	}
	return catalog
}

// EntryType defines the type of schedule entry
type EntryType string

const (
	EntryTypeRegular   EntryType = "REGULAR"
	EntryTypeException EntryType = "EXCEPTION"
)

// BranchSchedule represents the schedule configuration for a branch (HU30-35)
type BranchSchedule struct {
	ID        string     `json:"id"`
	BranchID  string     `json:"branch_id"`
	Active    bool       `json:"active"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"` // nullable - NULL means indefinite
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *BranchSchedule) SetID() {
	s.ID = uuid.Generate()
}

// ScheduleDetail represents a specific time slot (HU6-9, HU20-25)
type ScheduleDetail struct {
	ID                 string     `json:"id"`
	ScheduleID         string     `json:"schedule_id"`
	EntryType          EntryType  `json:"entry_type"`
	DayOfWeek          *int       `json:"day_of_week,omitempty"`
	ExceptionStartDate *time.Time `json:"exception_start_date,omitempty"`
	ExceptionEndDate   *time.Time `json:"exception_end_date,omitempty"`
	OpeningTime        *string    `json:"opening_time,omitempty"`
	ClosingTime        *string    `json:"closing_time,omitempty"`
	IsClosed           bool       `json:"is_closed"`
	Active             bool       `json:"active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (d *ScheduleDetail) SetID() {
	d.ID = uuid.Generate()
}
