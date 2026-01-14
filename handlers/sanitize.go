package handlers

import "strings"

// TrimString removes leading and trailing whitespace from a string
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// TrimStringPtr trims a string pointer, returns nil if pointer is nil
func TrimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

// Sanitizable interface for DTOs that can be sanitized
type Sanitizable interface {
	Sanitize()
}
