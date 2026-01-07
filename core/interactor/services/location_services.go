package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

type locationService struct {
	repository output.LocationRepository
}

// NewLocationService creates a new LocationService implementation
func NewLocationService(repository output.LocationRepository) *locationService {
	return &locationService{
		repository: repository,
	}
}

// GetAllDepartments retrieves all departments
func (s *locationService) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	return s.repository.GetAllDepartments(ctx)
}

// GetCitiesByDepartment retrieves all cities for a specific department
func (s *locationService) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	return s.repository.GetCitiesByDepartment(ctx, departmentID)
}

// ValidateCityInDepartment checks if the city belongs to the specified department
func (s *locationService) ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error {
	return s.repository.ValidateCityInDepartment(ctx, cityID, departmentID)
}
