package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// RegisterBranchRequest is the DTO for branch registration (HU59)
type RegisterBranchRequest struct {
	Name              string      `json:"name" binding:"required"`
	EstablishmentType string      `json:"establishment_type" binding:"required"` // WORKSHOP or STORE
	FranchiseID       *string     `json:"franchise_id,omitempty"`
	ProfileImageURL   *string     `json:"profile_image_url,omitempty"`
	Location          LocationDTO `json:"location" binding:"required"`
	Brands            []string    `json:"brands,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (r *RegisterBranchRequest) Sanitize() {
	r.Name = TrimString(r.Name)
	r.EstablishmentType = TrimString(r.EstablishmentType)
	r.FranchiseID = TrimStringPtr(r.FranchiseID)
	// ProfileImageURL is intentionally NOT trimmed - it's a URL
	r.Location.Sanitize()
	// Brands are trimmed individually
	for i := range r.Brands {
		r.Brands[i] = TrimString(r.Brands[i])
	}
}

// ToDomain maps RegisterBranchRequest to domain.Branch
func (r *RegisterBranchRequest) ToDomain(representativeID string) domain.Branch {
	branch := domain.Branch{
		RepresentativeID:  representativeID,
		Name:              r.Name,
		EstablishmentType: r.EstablishmentType,
		FranchiseID:       r.FranchiseID,
		ProfileImageURL:   r.ProfileImageURL,
		Brands:            r.Brands,
	}

	// Location is now required, always map it
	branch.Location = &domain.Location{
		DepartmentID:   r.Location.DepartmentID,
		CityID:         r.Location.CityID,
		CityName:       r.Location.CityName,       // For geocoding
		DepartmentName: r.Location.DepartmentName, // For geocoding
		Address:        r.Location.Address,
		Latitude:       r.Location.Latitude,
		Longitude:      r.Location.Longitude,
	}

	return branch
}

// LocationDTO represents location data in request/response (HU59)
type LocationDTO struct {
	DepartmentID   string   `json:"department_id" binding:"required"`
	CityID         string   `json:"city_id" binding:"required"`
	CityName       string   `json:"city_name,omitempty"`       // For geocoding (not persisted)
	DepartmentName string   `json:"department_name,omitempty"` // For geocoding (not persisted)
	Address        string   `json:"address" binding:"required"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
}

// Sanitize trims whitespace from all string fields
func (l *LocationDTO) Sanitize() {
	l.DepartmentID = TrimString(l.DepartmentID)
	l.CityID = TrimString(l.CityID)
	l.CityName = TrimString(l.CityName)
	l.DepartmentName = TrimString(l.DepartmentName)
	l.Address = TrimString(l.Address)
	// Latitude/Longitude are floats, no trimming needed
}

// GeocodingStatus indicates the result of automatic geocoding
type GeocodingStatus string

const (
	GeocodingStatusSuccess GeocodingStatus = "SUCCESS" // Coordinates were generated automatically
	GeocodingStatusFailed  GeocodingStatus = "FAILED"  // Geocoding failed, retry via update
	GeocodingStatusSkipped GeocodingStatus = "SKIPPED" // User provided coordinates, no geocoding needed
)

// BranchResponse is the DTO for branch responses (HU59)
type BranchResponse struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	EstablishmentType string          `json:"establishment_type"`
	Status            string          `json:"status"`
	FranchiseID       *string         `json:"franchise_id,omitempty"`
	ProfileImageURL   *string         `json:"profile_image_url,omitempty"`
	Location          *LocationDTO    `json:"location,omitempty"`
	Brands            []string        `json:"brands,omitempty"`
	GeocodingStatus   GeocodingStatus `json:"geocoding_status,omitempty"` // Indicates geocoding result
	Links             []Link          `json:"_links"`
}

// NewBranchResponse creates a BranchResponse from domain.Branch
func NewBranchResponse(branch *domain.Branch, encodedID string, geocodingStatus GeocodingStatus, links []Link) BranchResponse {
	response := BranchResponse{
		ID:                encodedID,
		Name:              branch.Name,
		EstablishmentType: branch.EstablishmentType,
		Status:            branch.Status,
		FranchiseID:       branch.FranchiseID,
		ProfileImageURL:   branch.ProfileImageURL,
		Brands:            branch.Brands,
		GeocodingStatus:   geocodingStatus,
		Links:             links,
	}

	if branch.Location != nil {
		response.Location = &LocationDTO{
			CityID:    branch.Location.CityID,
			Address:   branch.Location.Address,
			Latitude:  branch.Location.Latitude,
			Longitude: branch.Location.Longitude,
		}
	}

	return response
}
