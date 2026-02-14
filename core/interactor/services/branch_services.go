package services

import (
	"context"
	"errors"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/geocoding"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// branchService implements input.BranchService
type branchService struct {
	repository      output.BranchRepository
	locationRepo    output.LocationRepository
	geocodingClient geocoding.Client
}

// NewBranchService creates a new BranchService instance
func NewBranchService(repo output.BranchRepository, locationRepo output.LocationRepository, geocodingClient geocoding.Client) input.BranchService {
	return &branchService{
		repository:      repo,
		locationRepo:    locationRepo,
		geocodingClient: geocodingClient,
	}
}

// BeginTx starts a new database transaction
func (s *branchService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

// GeocodeLocation attempts to geocode a location if coordinates are not provided
// Returns: (coordsWereGenerated bool, error)
// - true if coordinates were successfully generated
// - false if coordinates already existed, geocoding failed, or no location
// This method modifies the location in-place
func (s *branchService) GeocodeLocation(ctx context.Context, location *domain.Location) (bool, error) {
	// No location to geocode
	if location == nil {
		return false, nil
	}

	// User already provided coordinates - skip geocoding
	if location.Latitude != nil && location.Longitude != nil {
		log.Debug(logger.LogGeocodingSkipped,
			"address", location.Address,
			"lat", *location.Latitude,
			"lng", *location.Longitude)
		return false, nil
	}

	// Validate we have the required data for geocoding
	if location.CityName == "" || location.DepartmentName == "" {
		log.Warn(logger.LogGeocodingCityError,
			"city_name", location.CityName,
			"department_name", location.DepartmentName,
			"reason", "missing city or department name for geocoding")
		return false, nil
	}

	// Attempt geocoding
	coords, err := s.geocodingClient.Geocode(ctx, location.Address, location.CityName, location.DepartmentName)
	if err != nil {
		log.Warn(logger.LogGeocodingError,
			"address", location.Address,
			"city", location.CityName,
			"department", location.DepartmentName,
			"error", err)
		return false, nil // Don't fail the registration, just log
	}

	// No results from geocoding
	if coords == nil {
		log.Warn(logger.LogGeocodingNoResults,
			"address", location.Address,
			"city", location.CityName,
			"department", location.DepartmentName)
		return false, nil
	}

	// Success! Update location with coordinates
	location.Latitude = &coords.Latitude
	location.Longitude = &coords.Longitude

	log.Info(logger.LogGeocodingSuccess,
		"address", location.Address,
		"city", location.CityName,
		"lat", coords.Latitude,
		"lng", coords.Longitude,
		"confidence", coords.Confidence)

	return true, nil
}

// RegisterBranch registers a new branch with validation
// Returns the branch with generated ID and default values
func (s *branchService) RegisterBranch(ctx context.Context, tx output.Tx, branch domain.Branch) (*domain.Branch, error) {
	// 1. Validate establishment type
	if !branch.IsValidEstablishmentType() {
		log.Warn(logger.LogBranchServiceInvalidType, "type", branch.EstablishmentType)
		return nil, domain.ErrInvalidBranchType
	}

	// 2. Check for duplicate name within franchise (only if franchise is set)
	if branch.FranchiseID != nil && *branch.FranchiseID != "" {
		existingBranch, err := s.repository.GetBranchByFranchiseAndName(ctx, *branch.FranchiseID, branch.Name)
		if err != nil && !errors.Is(err, domain.ErrBranchNotFound) {
			log.Error(logger.LogBranchServiceDupNameCheck, "error", err)
			return nil, err
		}
		if existingBranch != nil {
			log.Warn(logger.LogBranchServiceDupName, "name", branch.Name, "franchise_id", *branch.FranchiseID)
			return nil, domain.ErrDuplicateBranchName
		}
	}

	// 3. Generate UUID if not set
	if branch.ID == "" {
		branch.SetID()
	}

	// 4. Set default status if not set
	if branch.Status == "" {
		branch.Status = domain.BranchStatusActive
	}

	// 5. Save branch
	if err := s.repository.SaveBranch(ctx, tx, branch); err != nil {
		log.Error(logger.LogBranchServiceSaveError, "error", err, "branch_id", branch.ID)
		return nil, err
	}

	// 6. Save location if provided
	if branch.Location != nil {
		// 6.1 Check for duplicate address
		addressExists, err := s.locationRepo.CheckAddressExists(ctx, branch.Location.Address)
		if err != nil {
			log.Error(logger.LogBranchServiceLocSaveError, "error", err, "address", branch.Location.Address)
			return nil, err
		}
		if addressExists {
			log.Warn(logger.LogBranchServiceLocSaveError, "duplicate_address", branch.Location.Address)
			return nil, domain.ErrDuplicateAddress
		}

		branch.Location.BranchID = branch.ID
		if err := s.locationRepo.SaveLocation(ctx, tx, *branch.Location); err != nil {
			log.Error(logger.LogBranchServiceLocSaveError, "error", err, "branch_id", branch.ID)
			return nil, err
		}
	}

	// 7. Save brands if provided
	if len(branch.Brands) > 0 {
		if err := s.repository.SaveBranchBrands(ctx, tx, branch.ID, branch.Brands); err != nil {
			log.Error(logger.LogBranchServiceBrandSaveErr, "error", err, "branch_id", branch.ID)
			return nil, err
		}
	}

	// 8. Save displacement ranges if provided
	if len(branch.DisplacementRanges) > 0 {
		if err := s.repository.SaveBranchDisplacementRanges(ctx, tx, branch.ID, branch.DisplacementRanges); err != nil {
			log.Error(logger.LogBranchRepoDisplRangeSaveError, "error", err, "branch_id", branch.ID)
			return nil, err
		}
	}

	log.Info(logger.LogBranchServiceRegComplete, "branch_id", branch.ID, "name", branch.Name)
	return &branch, nil
}

// GetBranchByID retrieves a branch by its ID
func (s *branchService) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	branch, err := s.repository.GetBranchByID(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchServiceGetError, "error", err, "branch_id", branchID)
		return nil, err
	}
	return branch, nil
}

// ValidateBrands validates that all brands exist in motorcycle_references table
func (s *branchService) ValidateBrands(ctx context.Context, brands []string) error {
	if len(brands) == 0 {
		return nil
	}
	return s.repository.ValidateBrands(ctx, brands)
}

// SaveLocation saves a location for a branch
func (s *branchService) SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	return s.locationRepo.SaveLocation(ctx, tx, location)
}

// SaveBranchBrands saves brands for a branch
func (s *branchService) SaveBranchBrands(ctx context.Context, tx output.Tx, branchID string, brands []string) error {
	return s.repository.SaveBranchBrands(ctx, tx, branchID, brands)
}

// ValidateDisplacementRanges validates that all displacement ranges are valid ENUM values
func (s *branchService) ValidateDisplacementRanges(ranges []string) error {
	return domain.ValidateDisplacementRanges(ranges)
}

// SaveBranchDisplacementRanges saves displacement ranges for a branch
func (s *branchService) SaveBranchDisplacementRanges(ctx context.Context, tx output.Tx, branchID string, ranges []string) error {
	return s.repository.SaveBranchDisplacementRanges(ctx, tx, branchID, ranges)
}

// GetBranchesByRepresentative retrieves all branches for a representative (HU62)
func (s *branchService) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	return s.repository.GetBranchesByRepresentative(ctx, representativeID)
}

// UpdateBranch updates an existing branch (HU60)
func (s *branchService) UpdateBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error {
	// 1. Validate establishment type
	if !branch.IsValidEstablishmentType() {
		log.Warn(logger.LogBranchServiceInvalidType, "type", branch.EstablishmentType)
		return domain.ErrInvalidBranchType
	}

	// 2. Update branch core fields
	if err := s.repository.UpdateBranch(ctx, tx, branch); err != nil {
		log.Error(logger.LogBranchServiceSaveError, "error", err, "branch_id", branch.ID)
		return err
	}

	// 3. Update location if provided
	if branch.Location != nil {
		branch.Location.BranchID = branch.ID
		if err := s.locationRepo.UpdateLocation(ctx, tx, *branch.Location); err != nil {
			log.Error(logger.LogBranchServiceLocSaveError, "error", err, "branch_id", branch.ID)
			return err
		}
	}

	// 4. Update brands: delete existing and save new
	if err := s.repository.DeleteBranchBrands(ctx, tx, branch.ID); err != nil {
		log.Error(logger.LogBranchServiceBrandSaveErr, "error", err, "branch_id", branch.ID)
		return err
	}

	if len(branch.Brands) > 0 {
		if err := s.repository.SaveBranchBrands(ctx, tx, branch.ID, branch.Brands); err != nil {
			log.Error(logger.LogBranchServiceBrandSaveErr, "error", err, "branch_id", branch.ID)
			return err
		}
	}

	// 5. Update displacement ranges: delete existing and save new
	if err := s.repository.DeleteBranchDisplacementRanges(ctx, tx, branch.ID); err != nil {
		log.Error(logger.LogBranchRepoDisplRangeDelError, "error", err, "branch_id", branch.ID)
		return err
	}

	if len(branch.DisplacementRanges) > 0 {
		if err := s.repository.SaveBranchDisplacementRanges(ctx, tx, branch.ID, branch.DisplacementRanges); err != nil {
			log.Error(logger.LogBranchRepoDisplRangeSaveError, "error", err, "branch_id", branch.ID)
			return err
		}
	}

	log.Info(logger.LogBranchServiceRegComplete, "branch_id", branch.ID, "name", branch.Name)
	return nil
}

// DeleteBranch deletes a branch by ID (HU61)
// Related data (brands, location, schedules, services) is handled by CASCADE
// diagnostics and completed_services have RESTRICT - will error if exist
func (s *branchService) DeleteBranch(ctx context.Context, tx output.Tx, branchID string) error {
	if err := s.repository.DeleteBranch(ctx, tx, branchID); err != nil {
		log.Error(logger.LogBranchServiceDelError, "error", err, "branch_id", branchID)
		// Interpret FK constraint violations as domain errors
		if isForeignKeyError(err) {
			return domain.ErrBranchCannotDelete
		}
		return err
	}

	log.Info(logger.LogBranchServiceDelComplete, "branch_id", branchID)
	return nil
}

// isForeignKeyError checks if the error is a MySQL foreign key constraint violation
func isForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL error code 1451: Cannot delete or update a parent row: a foreign key constraint fails
	return strings.Contains(errStr, "1451") || strings.Contains(errStr, "foreign key constraint")
}

// GetBranchesNearby retrieves branches within radius of given coordinates (HU89)
// Uses bounding box pre-filter for optimization before Haversine calculation
func (s *branchService) GetBranchesNearby(ctx context.Context, lat, lng, radiusKm float64, establishmentType, brandID, displacementRange string) ([]domain.NearbyBranch, error) {
	const earthRadiusKm = 6371.0
	const degToRad = 3.141592653589793 / 180.0

	// Calculate bounding box for optimization
	// latDelta = radiusKm / earthRadiusKm * (180 / π)
	latDelta := radiusKm / earthRadiusKm * (180.0 / 3.141592653589793)

	// lngDelta = radiusKm / (earthRadiusKm * cos(lat in radians)) * (180 / π)
	cosLat := cosine(lat * degToRad)
	lngDelta := radiusKm / (earthRadiusKm * cosLat) * (180.0 / 3.141592653589793)

	latMin := lat - latDelta
	latMax := lat + latDelta
	lngMin := lng - lngDelta
	lngMax := lng + lngDelta

	return s.repository.GetBranchesNearby(ctx, lat, lng, radiusKm, establishmentType, latMin, latMax, lngMin, lngMax, brandID, displacementRange)
}

// cosine calculates cosine using Taylor series (avoids math import)
func cosine(x float64) float64 {
	// Normalize x to [-π, π]
	for x > 3.141592653589793 {
		x -= 2 * 3.141592653589793
	}
	for x < -3.141592653589793 {
		x += 2 * 3.141592653589793
	}

	// Taylor series: cos(x) = 1 - x²/2! + x⁴/4! - x⁶/6! + ...
	x2 := x * x
	return 1 - x2/2 + x2*x2/24 - x2*x2*x2/720 + x2*x2*x2*x2/40320
}
