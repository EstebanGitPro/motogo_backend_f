package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// LocationInteractor handles location-related use cases
type LocationInteractor struct {
	locationService input.LocationService
}

// NewLocationInteractor creates a new LocationInteractor
func NewLocationInteractor(locationService input.LocationService) *LocationInteractor {
	return &LocationInteractor{
		locationService: locationService,
	}
}

// GetAllDepartments retrieves all departments
func (i *LocationInteractor) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	log.Info(logger.LogLocationInteractorGetDepartments)

	departments, err := i.locationService.GetAllDepartments(ctx)
	if err != nil {
		log.Error(logger.LogLocationInteractorGetDepartmentsError, "error", err)
		return nil, err
	}

	log.Success(logger.LogLocationInteractorGetDepartmentsOK, "count", len(departments))
	return departments, nil
}

// GetCitiesByDepartment retrieves all cities for a specific department
func (i *LocationInteractor) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	log.Info(logger.LogLocationInteractorGetCities, "department_id", departmentID)

	cities, err := i.locationService.GetCitiesByDepartment(ctx, departmentID)
	if err != nil {
		log.Error(logger.LogLocationInteractorGetCitiesError, "error", err, "department_id", departmentID)
		return nil, err
	}

	log.Success(logger.LogLocationInteractorGetCitiesOK, "count", len(cities), "department_id", departmentID)
	return cities, nil
}
