package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

const (
	EstablishmentTypeWorkshop      = "WORKSHOP"
	EstablishmentTypeStore         = "STORE"
	EstablishmentTypeWorkshopStore = "WORKSHOP_STORE"
)

type EstablishmentTypeInfo struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

func GetAllEstablishmentTypes() []EstablishmentTypeInfo {
	return []EstablishmentTypeInfo{
		{Code: EstablishmentTypeWorkshop, Label: "Taller"},
		{Code: EstablishmentTypeStore, Label: "Tienda"},
		{Code: EstablishmentTypeWorkshopStore, Label: "Taller y Tienda"},
	}
}

func GetEstablishmentTypeLabel(code string) string {
	labels := map[string]string{
		EstablishmentTypeWorkshop:      "Taller",
		EstablishmentTypeStore:         "Tienda",
		EstablishmentTypeWorkshopStore: "Taller y Tienda",
	}
	if label, ok := labels[code]; ok {
		return label
	}
	return code
}

const (
	BranchStatusActive   = "ACTIVE"
	BranchStatusInactive = "INACTIVE"
)

// Displacement range constants (ENUM values matching DB)
const (
	DisplacementRangeLow    = "BAJO"  // 50-200cc
	DisplacementRangeMedium = "MEDIO" // 201-400cc
	DisplacementRangeHigh   = "ALTO"  // 401-3000cc
)

// IsValidDisplacementRange checks if a string is a valid displacement range ENUM value
func IsValidDisplacementRange(r string) bool {
	return r == DisplacementRangeLow ||
		r == DisplacementRangeMedium ||
		r == DisplacementRangeHigh
}

// ValidateDisplacementRanges checks that all displacement ranges are valid ENUM values
func ValidateDisplacementRanges(ranges []string) error {
	for _, r := range ranges {
		if !IsValidDisplacementRange(r) {
			return ErrInvalidDisplacementRange
		}
	}
	return nil
}

type Branch struct {
	ID                  string    `json:"id"`
	RepresentativeID    string    `json:"representative_id"`
	RepresentativePhone *string   `json:"representative_phone,omitempty"` // From JOIN with persons (read-only)
	FranchiseID         *string   `json:"franchise_id,omitempty"`
	Name                string    `json:"name"`
	EstablishmentType   string    `json:"establishment_type"`
	ProfileImageURL     *string   `json:"profile_image_url,omitempty"`
	Status              string    `json:"status"`
	Location            *Location `json:"location,omitempty"`
	Brands              []string  `json:"brands,omitempty"`
	DisplacementRanges  []string  `json:"displacement_ranges,omitempty"`
}

type Location struct {
	ID             string   `json:"id,omitempty"`
	BranchID       string   `json:"branch_id,omitempty"`
	DepartmentID   string   `json:"department_id"`
	CityID         string   `json:"city_id"`
	Address        string   `json:"address"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	CityName       string   `json:"-"`
	DepartmentName string   `json:"-"`
}

func (b *Branch) SetID() {
	b.ID = uuid.Generate()
}

// IsValidEstablishmentType checks if the establishment type is valid
func (b *Branch) IsValidEstablishmentType() bool {
	return b.EstablishmentType == EstablishmentTypeWorkshop ||
		b.EstablishmentType == EstablishmentTypeStore ||
		b.EstablishmentType == EstablishmentTypeWorkshopStore
}

// IsValidEstablishmentType checks if a string is a valid establishment type
func IsValidEstablishmentType(t string) bool {
	return t == EstablishmentTypeWorkshop ||
		t == EstablishmentTypeStore ||
		t == EstablishmentTypeWorkshopStore
}

func (b *Branch) ToLogger() []string {
	return []string{
		"id:" + b.ID,
		"name:" + b.Name,
		"type:" + b.EstablishmentType,
		"representative_id:" + b.RepresentativeID,
	}
}

// NearbyBranch represents a branch with distance information for proximity search (HU89)
type NearbyBranch struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	EstablishmentType  string          `json:"establishment_type"`
	ProfileImageURL    *string         `json:"profile_image_url,omitempty"`
	ContactPhone       *string         `json:"contact_phone,omitempty"` // Representative's phone
	Location           *NearbyLocation `json:"location,omitempty"`
	DistanceKm         float64         `json:"distance_km"`
	Brands             []string        `json:"brands,omitempty"`
	DisplacementRanges []string        `json:"displacement_ranges,omitempty"`
}

// NearbyLocation contains location info for nearby search results (HU89)
type NearbyLocation struct {
	Address        string  `json:"address"`
	CityName       string  `json:"city_name"`
	DepartmentName string  `json:"department_name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}
