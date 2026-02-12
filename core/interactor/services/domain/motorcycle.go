package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// Motorcycle represents a registered motorcycle in the system
type Motorcycle struct {
	ID              string               // UUID primary key
	LicensePlate    string               // Unique motorcycle plate
	ReferenceID     string               // FK to motorcycle_references catalog
	OwnerID         string               // FK to person with USER role
	Year            *int                 // Motorcycle year (optional)
	CurrentMileage  *int                 // Current mileage (optional)
	OwnerNotes      *string              // Owner notes (optional)
	ProfileImageURL *string              // Main profile photo URL (Firebase Storage) - HU36
	Reference       *MotorcycleReference // Motorcycle reference with brand info (populated on read)
}

func (m *Motorcycle) SetID() {
	m.ID = uuid.Generate()
}

// MotorcycleReference represents the motorcycle catalog reference
type MotorcycleReference struct {
	ID                 string // UUID primary key
	BrandID            string // FK to brand catalog
	BrandName          string // Brand name (denormalized for display)
	Model              string // Model name (e.g., "CB 190R")
	Category           string // Category (Sport, Scooter, etc.)
	EngineDisplacement int    // Engine displacement in cc
}
