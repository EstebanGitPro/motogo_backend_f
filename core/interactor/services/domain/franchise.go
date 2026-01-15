package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// Franchise represents a franchise that groups multiple branches (HU26-29)
type Franchise struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Branches    []Branch `json:"branches,omitempty"` // Populated on query via JOIN
}

// SetID generates a new UUID for the franchise
func (f *Franchise) SetID() {
	f.ID = uuid.Generate()
}

// ToLogger returns a slice of strings for structured logging
func (f *Franchise) ToLogger() []string {
	return []string{
		"id:" + f.ID,
		"name:" + f.Name,
	}
}
