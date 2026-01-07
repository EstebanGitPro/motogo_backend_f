package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// DepartmentResponse represents a single department in the API response
type DepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DepartmentListResponse represents the response for GET /departments
type DepartmentListResponse struct {
	Departments []DepartmentResponse `json:"departments"`
	Links       []Link               `json:"_links"`
}

// CityResponse represents a single city in the API response
type CityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CityListResponse represents the response for GET /departments/:id/cities
type CityListResponse struct {
	Cities []CityResponse `json:"cities"`
	Links  []Link         `json:"_links"`
}

// NewDepartmentListResponse creates a DepartmentListResponse from domain entities
func NewDepartmentListResponse(departments []domain.Department, links []Link) DepartmentListResponse {
	deptResponses := make([]DepartmentResponse, len(departments))
	for i, dept := range departments {
		deptResponses[i] = DepartmentResponse{
			ID:   dept.ID,
			Name: dept.Name,
		}
	}
	return DepartmentListResponse{
		Departments: deptResponses,
		Links:       links,
	}
}

// NewCityListResponse creates a CityListResponse from domain entities
func NewCityListResponse(cities []domain.City, links []Link) CityListResponse {
	cityResponses := make([]CityResponse, len(cities))
	for i, city := range cities {
		cityResponses[i] = CityResponse{
			ID:   city.ID,
			Name: city.Name,
		}
	}
	return CityListResponse{
		Cities: cityResponses,
		Links:  links,
	}
}
