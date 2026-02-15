package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// EstablishmentType represents the type of business establishment
type EstablishmentType string

const (
	EstablishmentTypeWorkshop      EstablishmentType = "WORKSHOP"
	EstablishmentTypeStore         EstablishmentType = "STORE"
	EstablishmentTypeWorkshopStore EstablishmentType = "WORKSHOP_STORE"
)

type EstablishmentTypeInfo struct {
	Code  EstablishmentType `json:"code"`
	Label string            `json:"label"`
}

func GetAllEstablishmentTypes() []EstablishmentTypeInfo {
	return []EstablishmentTypeInfo{
		{Code: EstablishmentTypeWorkshop, Label: "Taller"},
		{Code: EstablishmentTypeStore, Label: "Tienda"},
		{Code: EstablishmentTypeWorkshopStore, Label: "Taller y Tienda"},
	}
}

func GetEstablishmentTypeLabel(code EstablishmentType) string {
	labels := map[EstablishmentType]string{
		EstablishmentTypeWorkshop:      "Taller",
		EstablishmentTypeStore:         "Tienda",
		EstablishmentTypeWorkshopStore: "Taller y Tienda",
	}
	if label, ok := labels[code]; ok {
		return label
	}
	return string(code)
}

// BranchStatus represents the active state of a branch
type BranchStatus string

const (
	BranchStatusActive   BranchStatus = "ACTIVE"
	BranchStatusInactive BranchStatus = "INACTIVE"
)

// DisplacementRange represents motorcycle engine displacement ranges (ENUM values matching DB)
type DisplacementRange string

const (
	DisplacementRangeLow    DisplacementRange = "BAJO"  // 50-200cc
	DisplacementRangeMedium DisplacementRange = "MEDIO" // 201-400cc
	DisplacementRangeHigh   DisplacementRange = "ALTO"  // 401-3000cc
)

// IsValidDisplacementRange checks if a string is a valid displacement range ENUM value
func IsValidDisplacementRange(r DisplacementRange) bool {
	return r == DisplacementRangeLow ||
		r == DisplacementRangeMedium ||
		r == DisplacementRangeHigh
}

// ValidateDisplacementRanges checks that all displacement ranges are valid ENUM values
func ValidateDisplacementRanges(ranges []DisplacementRange) error {
	for _, r := range ranges {
		if !IsValidDisplacementRange(r) {
			return ErrInvalidDisplacementRange
		}
	}
	return nil
}

type Branch struct {
	ID                  string              `json:"id"`
	RepresentativeID    string              `json:"representative_id"`
	RepresentativePhone *string             `json:"representative_phone,omitempty"` // From JOIN with persons (read-only)
	FranchiseID         *string             `json:"franchise_id,omitempty"`
	Name                string              `json:"name"`
	EstablishmentType   EstablishmentType   `json:"establishment_type"`
	ProfileImageURL     *string             `json:"profile_image_url,omitempty"`
	Status              BranchStatus        `json:"status"`
	Location            *Location           `json:"location,omitempty"`
	Brands              []string            `json:"brands,omitempty"`
	DisplacementRanges  []DisplacementRange `json:"displacement_ranges,omitempty"`
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
func IsValidEstablishmentType(t EstablishmentType) bool {
	return t == EstablishmentTypeWorkshop ||
		t == EstablishmentTypeStore ||
		t == EstablishmentTypeWorkshopStore
}

func (b *Branch) ToLogger() []string {
	return []string{
		"id:" + b.ID,
		"name:" + b.Name,
		"type:" + string(b.EstablishmentType),
		"representative_id:" + b.RepresentativeID,
	}
}

// NearbyBranch represents a branch with distance information for proximity search (HU89)
type NearbyBranch struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	EstablishmentType  EstablishmentType   `json:"establishment_type"`
	ProfileImageURL    *string             `json:"profile_image_url,omitempty"`
	ContactPhone       *string             `json:"contact_phone,omitempty"` // Representative's phone
	Location           *NearbyLocation     `json:"location,omitempty"`
	DistanceKm         float64             `json:"distance_km"`
	Brands             []string            `json:"brands,omitempty"`
	DisplacementRanges []DisplacementRange `json:"displacement_ranges,omitempty"`
}

// NearbyLocation contains location info for nearby search results (HU89)
type NearbyLocation struct {
	Address        string  `json:"address"`
	CityName       string  `json:"city_name"`
	DepartmentName string  `json:"department_name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}
