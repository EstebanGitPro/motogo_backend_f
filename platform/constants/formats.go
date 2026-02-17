package constants

import "regexp"

// Date/Time format constants — Go reference time: Mon Jan 2 15:04:05 MST 2006

// Date formats
const (
	DateFormat     = "2006-01-02"                // ISO 8601 date
	DateTimeFormat = "2006-01-02 15:04:05"       // Date + time for display
	DateTimeISO    = "2006-01-02T15:04:05Z07:00" // RFC3339 variant
)

// Time formats
const (
	TimeFormatLong  = "15:04:05" // HH:mm:ss (from DB)
	TimeFormatShort = "15:04"    // HH:mm (from API)
)

// TimeRegex validates HH:MM or HH:MM:SS format (24-hour, optional seconds)
var TimeRegex = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$`)
