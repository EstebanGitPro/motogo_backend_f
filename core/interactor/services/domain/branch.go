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

type Branch struct {
	ID                string    `json:"id"`
	RepresentativeID  string    `json:"representative_id"`
	FranchiseID       *string   `json:"franchise_id,omitempty"`
	Name              string    `json:"name"`
	EstablishmentType string    `json:"establishment_type"`
	ProfileImageURL   *string   `json:"profile_image_url,omitempty"`
	Status            string    `json:"status"`
	Location          *Location `json:"location,omitempty"`
	Brands            []string  `json:"brands,omitempty"`
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

func (b *Branch) ToLogger() []string {
	return []string{
		"id:" + b.ID,
		"name:" + b.Name,
		"type:" + b.EstablishmentType,
		"representative_id:" + b.RepresentativeID,
	}
}
