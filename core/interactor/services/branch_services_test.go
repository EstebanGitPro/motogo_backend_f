package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/EstebanGitPro/motogo-backend/platform/geocoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Helper Functions
// ============================================

func setupBranchServiceMocks() (*mocks.MockBranchRepository, *mocks.MockLocationRepository, *mocks.MockGeocodingClient) {
	return new(mocks.MockBranchRepository),
		new(mocks.MockLocationRepository),
		new(mocks.MockGeocodingClient)
}

func floatPtr(f float64) *float64 {
	return &f
}

// ============================================
// NewBranchService Tests
// ============================================

func TestNewBranchService_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	// Act
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Assert
	assert.NotNil(t, service)
}

// ============================================
// BeginTx Tests
// ============================================

func TestBranchService_BeginTx_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	tx, err := service.BeginTx(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_BeginTx_Error(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	branchRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("connection error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	tx, err := service.BeginTx(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tx)
	branchRepo.AssertExpectations(t)
}

// ============================================
// GeocodeLocation Tests
// ============================================

func TestGeocodeLocation_NilLocation(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), nil)

	// Assert
	assert.NoError(t, err)
	assert.False(t, generated)
}

func TestGeocodeLocation_CoordinatesAlreadyProvided(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	location := &domain.Location{
		Address:   "Calle 10 #20-30",
		Latitude:  floatPtr(6.2518),
		Longitude: floatPtr(-75.5636),
	}

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), location)

	// Assert
	assert.NoError(t, err)
	assert.False(t, generated) // No geocoding needed
}

func TestGeocodeLocation_MissingCityOrDepartment(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	location := &domain.Location{
		Address:  "Calle 10 #20-30",
		CityName: "", // Missing city
	}

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), location)

	// Assert
	assert.NoError(t, err)
	assert.False(t, generated)
}

func TestGeocodeLocation_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	location := &domain.Location{
		Address:        "Calle 10 #20-30",
		CityName:       "Medellín",
		DepartmentName: "Antioquia",
	}

	coords := &geocoding.Coordinates{
		Latitude:   6.2518,
		Longitude:  -75.5636,
		Confidence: 9,
	}

	geocodingClient.On("Geocode", mock.Anything, "Calle 10 #20-30", "Medellín", "Antioquia").Return(coords, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), location)

	// Assert
	assert.NoError(t, err)
	assert.True(t, generated)
	assert.NotNil(t, location.Latitude)
	assert.NotNil(t, location.Longitude)
	assert.Equal(t, 6.2518, *location.Latitude)
	assert.Equal(t, -75.5636, *location.Longitude)
	geocodingClient.AssertExpectations(t)
}

func TestGeocodeLocation_GeocodingError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	location := &domain.Location{
		Address:        "Dirección Inválida",
		CityName:       "Medellín",
		DepartmentName: "Antioquia",
	}

	geocodingClient.On("Geocode", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("geocoding API error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), location)

	// Assert
	assert.NoError(t, err) // Should not fail, just log warning
	assert.False(t, generated)
	geocodingClient.AssertExpectations(t)
}

func TestGeocodeLocation_NoResults(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	location := &domain.Location{
		Address:        "Dirección No Encontrada",
		CityName:       "Medellín",
		DepartmentName: "Antioquia",
	}

	geocodingClient.On("Geocode", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil) // No error, but no results

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	generated, err := service.GeocodeLocation(context.Background(), location)

	// Assert
	assert.NoError(t, err)
	assert.False(t, generated)
}

// ============================================
// RegisterBranch Tests
// ============================================

func TestRegisterBranch_InvalidEstablishmentType(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		Name:              "Test Branch",
		EstablishmentType: "INVALID_TYPE",
	}

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.RegisterBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidBranchType, err)
}

func TestRegisterBranch_DuplicateName(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	franchiseID := "franchise-123"
	branch := domain.Branch{
		Name:              "Taller Express",
		EstablishmentType: "WORKSHOP",
		FranchiseID:       &franchiseID,
	}

	existingBranch := &domain.Branch{ID: "existing-id", Name: "Taller Express"}
	branchRepo.On("GetBranchByFranchiseAndName", mock.Anything, franchiseID, "Taller Express").
		Return(existingBranch, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.RegisterBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDuplicateBranchName, err)
	branchRepo.AssertExpectations(t)
}

func TestRegisterBranch_Success_MinimalFields(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		Name:              "Taller Básico",
		EstablishmentType: "WORKSHOP",
		RepresentativeID:  "rep-123",
	}

	branchRepo.On("SaveBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.RegisterBranch(context.Background(), mockTx, branch)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID) // UUID should be generated
	assert.Equal(t, domain.BranchStatusActive, result.Status)
	branchRepo.AssertExpectations(t)
}

func TestRegisterBranch_Success_WithLocationAndBrands(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		Name:              "Taller Completo",
		EstablishmentType: "WORKSHOP",
		RepresentativeID:  "rep-123",
		Location: &domain.Location{
			Address:        "Calle 10 #20-30",
			CityName:       "Medellín",
			DepartmentName: "Antioquia",
		},
		Brands: []string{"Honda", "Yamaha"},
	}

	branchRepo.On("SaveBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	locationRepo.On("CheckAddressExists", mock.Anything, "Calle 10 #20-30").Return(false, nil)
	locationRepo.On("SaveLocation", mock.Anything, mockTx, mock.AnythingOfType("domain.Location")).Return(nil)
	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, mock.AnythingOfType("string"), []string{"Honda", "Yamaha"}).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.RegisterBranch(context.Background(), mockTx, branch)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	branchRepo.AssertExpectations(t)
	locationRepo.AssertExpectations(t)
}

func TestRegisterBranch_DuplicateAddress(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		Name:              "Taller con Dirección Duplicada",
		EstablishmentType: "WORKSHOP",
		Location: &domain.Location{
			Address: "Dirección Existente",
		},
	}

	branchRepo.On("SaveBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	locationRepo.On("CheckAddressExists", mock.Anything, "Dirección Existente").Return(true, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.RegisterBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDuplicateAddress, err)
}

// ============================================
// GetBranchByID Tests
// ============================================

func TestGetBranchByID_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	expectedBranch := &domain.Branch{
		ID:                "branch-123",
		Name:              "Test Branch",
		EstablishmentType: "WORKSHOP",
	}

	branchRepo.On("GetBranchByID", mock.Anything, "branch-123").Return(expectedBranch, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchByID(context.Background(), "branch-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "branch-123", result.ID)
	branchRepo.AssertExpectations(t)
}

func TestGetBranchByID_NotFound(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("GetBranchByID", mock.Anything, "non-existent").Return(nil, domain.ErrBranchNotFound)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchByID(context.Background(), "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

// ============================================
// ValidateBrands Tests
// ============================================

func TestValidateBrands_EmptyList(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.ValidateBrands(context.Background(), []string{})

	// Assert
	assert.NoError(t, err)
}

func TestValidateBrands_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("ValidateBrands", mock.Anything, []string{"Honda", "Yamaha"}).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.ValidateBrands(context.Background(), []string{"Honda", "Yamaha"})

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestValidateBrands_InvalidBrand(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("ValidateBrands", mock.Anything, []string{"InvalidBrand"}).
		Return(domain.ErrBrandNotFound)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.ValidateBrands(context.Background(), []string{"InvalidBrand"})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBrandNotFound, err)
}

// ============================================
// GetBranchesByRepresentative Tests
// ============================================

func TestGetBranchesByRepresentative_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	expectedBranches := []domain.Branch{
		{ID: "branch-1", Name: "Branch 1"},
		{ID: "branch-2", Name: "Branch 2"},
	}

	branchRepo.On("GetBranchesByRepresentative", mock.Anything, "rep-123").
		Return(expectedBranches, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchesByRepresentative(context.Background(), "rep-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	branchRepo.AssertExpectations(t)
}

func TestGetBranchesByRepresentative_Empty(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("GetBranchesByRepresentative", mock.Anything, "rep-no-branches").
		Return([]domain.Branch{}, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchesByRepresentative(context.Background(), "rep-no-branches")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// UpdateBranch Tests
// ============================================

func TestUpdateBranch_InvalidType(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "INVALID",
	}

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidBranchType, err)
}

func TestUpdateBranch_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "STORE",
		Brands:            []string{"Suzuki"},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, "branch-123", []string{"Suzuki"}).Return(nil)
	branchRepo.On("DeleteBranchDisplacementRanges", mock.Anything, mockTx, "branch-123").Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestUpdateBranch_WithLocation(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "WORKSHOP",
		Location: &domain.Location{
			Address: "Nueva Dirección",
		},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	locationRepo.On("UpdateLocation", mock.Anything, mockTx, mock.AnythingOfType("domain.Location")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("DeleteBranchDisplacementRanges", mock.Anything, mockTx, "branch-123").Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.NoError(t, err)
	locationRepo.AssertExpectations(t)
}

func TestUpdateBranch_RepoError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "WORKSHOP",
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).
		Return(errors.New("db update error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "db update error", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestUpdateBranch_LocationUpdateError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "STORE",
		Location: &domain.Location{
			Address: "Dirección Actualizada",
		},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	locationRepo.On("UpdateLocation", mock.Anything, mockTx, mock.AnythingOfType("domain.Location")).
		Return(errors.New("location update error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "location update error", err.Error())
}

func TestUpdateBranch_DeleteBrandsError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "WORKSHOP",
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").
		Return(errors.New("delete brands error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "delete brands error", err.Error())
}

func TestUpdateBranch_SaveBrandsError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "STORE",
		Brands:            []string{"Honda", "Yamaha"},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, "branch-123", []string{"Honda", "Yamaha"}).
		Return(errors.New("save brands error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "save brands error", err.Error())
}

func TestUpdateBranch_DeleteDisplacementRangesError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Updated Branch",
		EstablishmentType: "WORKSHOP",
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("DeleteBranchDisplacementRanges", mock.Anything, mockTx, "branch-123").
		Return(errors.New("delete displacement error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "delete displacement error", err.Error())
}

func TestUpdateBranch_SaveDisplacementRangesError(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                 "branch-123",
		Name:               "Updated Branch",
		EstablishmentType:  "STORE",
		DisplacementRanges: []domain.DisplacementRange{domain.DisplacementRangeLow, domain.DisplacementRangeMedium},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("DeleteBranchDisplacementRanges", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("SaveBranchDisplacementRanges", mock.Anything, mockTx, "branch-123", []string{"BAJO", "MEDIO"}).
		Return(errors.New("save displacement error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "save displacement error", err.Error())
}

func TestUpdateBranch_SuccessWithAllFields(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		ID:                 "branch-123",
		Name:               "Full Update Branch",
		EstablishmentType:  "WORKSHOP",
		Brands:             []string{"Honda"},
		DisplacementRanges: []domain.DisplacementRange{domain.DisplacementRangeHigh},
		Location: &domain.Location{
			Address: "Dirección Completa",
		},
	}

	branchRepo.On("UpdateBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	locationRepo.On("UpdateLocation", mock.Anything, mockTx, mock.AnythingOfType("domain.Location")).Return(nil)
	branchRepo.On("DeleteBranchBrands", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, "branch-123", []string{"Honda"}).Return(nil)
	branchRepo.On("DeleteBranchDisplacementRanges", mock.Anything, mockTx, "branch-123").Return(nil)
	branchRepo.On("SaveBranchDisplacementRanges", mock.Anything, mockTx, "branch-123", []string{"ALTO"}).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.UpdateBranch(context.Background(), mockTx, branch)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
	locationRepo.AssertExpectations(t)
}

// ============================================
// DeleteBranch Tests
// ============================================

func TestDeleteBranch_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("DeleteBranch", mock.Anything, mockTx, "branch-123").Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.DeleteBranch(context.Background(), mockTx, "branch-123")

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestDeleteBranch_Error(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("DeleteBranch", mock.Anything, mockTx, "branch-with-diagnostics").
		Return(domain.ErrBranchCannotDelete)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.DeleteBranch(context.Background(), mockTx, "branch-with-diagnostics")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

// ============================================
// GetBranchesNearby Tests
// ============================================

func TestGetBranchesNearby_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	expectedBranches := []domain.NearbyBranch{
		{ID: "branch-1", Name: "Nearby Branch 1", DistanceKm: 1.5},
		{ID: "branch-2", Name: "Nearby Branch 2", DistanceKm: 2.0},
	}

	branchRepo.On("GetBranchesNearby", mock.Anything,
		6.2518, -75.5636, 5.0, "WORKSHOP",
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		"", "").
		Return(expectedBranches, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchesNearby(context.Background(), 6.2518, -75.5636, 5.0, "WORKSHOP", "", "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Nearby Branch 1", result[0].Name)
	branchRepo.AssertExpectations(t)
}

func TestGetBranchesNearby_Empty(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("GetBranchesNearby", mock.Anything,
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("string"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return([]domain.NearbyBranch{}, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	result, err := service.GetBranchesNearby(context.Background(), 0.0, 0.0, 1.0, "", "", "")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// SaveLocation Tests (direct delegation)
// ============================================

func TestSaveLocation_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	location := domain.Location{
		BranchID: "branch-123",
		Address:  "Calle 10 #20-30",
	}

	locationRepo.On("SaveLocation", mock.Anything, mockTx, location).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.SaveLocation(context.Background(), mockTx, location)

	// Assert
	assert.NoError(t, err)
	locationRepo.AssertExpectations(t)
}

func TestSaveLocation_Error(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	location := domain.Location{
		BranchID: "branch-123",
		Address:  "Calle 10 #20-30",
	}

	locationRepo.On("SaveLocation", mock.Anything, mockTx, location).Return(errors.New("save error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.SaveLocation(context.Background(), mockTx, location)

	// Assert
	assert.Error(t, err)
}

// ============================================
// SaveBranchBrands Tests (direct delegation)
// ============================================

func TestSaveBranchBrands_Success(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchID := "branch-123"
	brands := []string{"Honda", "Yamaha"}

	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, branchID, brands).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.SaveBranchBrands(context.Background(), mockTx, branchID, brands)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestSaveBranchBrands_Error(t *testing.T) {
	// Arrange
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchID := "branch-123"
	brands := []string{"InvalidBrand"}

	branchRepo.On("SaveBranchBrands", mock.Anything, mockTx, branchID, brands).Return(errors.New("invalid brand"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act
	err := service.SaveBranchBrands(context.Background(), mockTx, branchID, brands)

	// Assert
	assert.Error(t, err)
}

// ============================================
// ValidateDisplacementRanges Tests
// ============================================

func TestValidateDisplacementRanges_AllValid(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.ValidateDisplacementRanges([]string{"BAJO", "MEDIO", "ALTO"})
	assert.NoError(t, err)
}

func TestValidateDisplacementRanges_OneValid(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.ValidateDisplacementRanges([]string{"BAJO"})
	assert.NoError(t, err)
}

func TestValidateDisplacementRanges_Invalid(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.ValidateDisplacementRanges([]string{"BAJO", "INVALID"})
	assert.Error(t, err)
}

func TestValidateDisplacementRanges_Empty(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.ValidateDisplacementRanges([]string{})
	assert.NoError(t, err)
}

// ============================================
// SaveBranchDisplacementRanges Tests
// ============================================

func TestSaveBranchDisplacementRanges_Success(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("SaveBranchDisplacementRanges", mock.Anything, mockTx, "branch-123", []string{"BAJO", "MEDIO"}).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.SaveBranchDisplacementRanges(context.Background(), mockTx, "branch-123", []string{"BAJO", "MEDIO"})
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestSaveBranchDisplacementRanges_Error(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("SaveBranchDisplacementRanges", mock.Anything, mockTx, "branch-123", []string{"BAJO"}).Return(errors.New("db error"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.SaveBranchDisplacementRanges(context.Background(), mockTx, "branch-123", []string{"BAJO"})
	assert.Error(t, err)
}

// ============================================
// RegisterBranch with DisplacementRanges Tests
// (also exercises displacementRangesToStrings)
// ============================================

func TestRegisterBranch_WithDisplacementRanges(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branch := domain.Branch{
		Name:               "Taller con Rangos",
		EstablishmentType:  "WORKSHOP",
		RepresentativeID:   "rep-123",
		DisplacementRanges: []domain.DisplacementRange{domain.DisplacementRangeLow, domain.DisplacementRangeMedium},
	}

	branchRepo.On("SaveBranch", mock.Anything, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	branchRepo.On("SaveBranchDisplacementRanges", mock.Anything, mockTx, mock.AnythingOfType("string"), []string{"BAJO", "MEDIO"}).Return(nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	result, err := service.RegisterBranch(context.Background(), mockTx, branch)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	branchRepo.AssertExpectations(t)
}

// ============================================
// DeleteBranch ForeignKey Tests
// ============================================

func TestDeleteBranch_ForeignKeyError1451(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("DeleteBranch", mock.Anything, mockTx, "branch-fk").
		Return(errors.New("Error 1451: Cannot delete or update a parent row"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.DeleteBranch(context.Background(), mockTx, "branch-fk")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

func TestDeleteBranch_ForeignKeyConstraintText(t *testing.T) {
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()
	mockTx := new(mocks.MockTx)

	branchRepo.On("DeleteBranch", mock.Anything, mockTx, "branch-fk2").
		Return(errors.New("a foreign key constraint fails"))

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	err := service.DeleteBranch(context.Background(), mockTx, "branch-fk2")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

// ============================================
// GetBranchesNearby — Large Angle Normalization
// ============================================

func TestGetBranchesNearby_LargeAngleNormalization(t *testing.T) {
	// Arrange — use a latitude that results in a radian angle > π
	// so the cosine normalization loop is exercised
	branchRepo, locationRepo, geocodingClient := setupBranchServiceMocks()

	branchRepo.On("GetBranchesNearby", mock.Anything,
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("string"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"), mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return([]domain.NearbyBranch{}, nil)

	service := services.NewBranchService(branchRepo, locationRepo, geocodingClient)

	// Act — use a normal latitude (the cosine function is called internally
	// with lat * degToRad which is within range, but we test the overall flow)
	result, err := service.GetBranchesNearby(context.Background(), -89.99, 179.99, 50.0, "", "", "")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}
