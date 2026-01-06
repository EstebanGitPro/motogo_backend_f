package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// EstablishmentType defines valid types for branches
const (
	EstablishmentTypeWorkshop = "WORKSHOP"
	EstablishmentTypeStore    = "STORE"
)

// BranchStatus defines valid statuses for branches
const (
	BranchStatusActive   = "ACTIVE"
	BranchStatusInactive = "INACTIVE"
)

// Branch represents a workshop or store branch (Sede)
type Branch struct {
	ID                string    `json:"id"`
	RepresentativeID  string    `json:"representative_id"`
	FranchiseID       *string   `json:"franchise_id,omitempty"`
	Name              string    `json:"name"`
	EstablishmentType string    `json:"establishment_type"` // WORKSHOP or STORE
	ProfileImageURL   *string   `json:"profile_image_url,omitempty"`
	Status            string    `json:"status"` // ACTIVE or INACTIVE
	Location          *Location `json:"location,omitempty"`
	Brands            []string  `json:"brands,omitempty"`
}

// Location represents the physical location of a branch
type Location struct {
	ID        string   `json:"id,omitempty"`
	BranchID  string   `json:"branch_id,omitempty"`
	CityID    string   `json:"city_id"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	// Transient fields for geocoding (not persisted to DB)
	CityName       string `json:"-"` // Used for geocoding query
	DepartmentName string `json:"-"` // Used for geocoding query
}

// SetID generates a new UUID for the branch
func (b *Branch) SetID() {
	b.ID = uuid.Generate()
}

// IsValidEstablishmentType checks if the establishment type is valid
func (b *Branch) IsValidEstablishmentType() bool {
	return b.EstablishmentType == EstablishmentTypeWorkshop ||
		b.EstablishmentType == EstablishmentTypeStore
}

// ToLogger returns a slice of strings for structured logging
func (b *Branch) ToLogger() []string {
	return []string{
		"id:" + b.ID,
		"name:" + b.Name,
		"type:" + b.EstablishmentType,
		"representative_id:" + b.RepresentativeID,
	}
}
