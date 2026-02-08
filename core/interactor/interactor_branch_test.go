package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// RegisterBranch Tests (HU59)
// ============================================

func TestRegisterBranch_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	lat := 4.710989
	lng := -74.072092

	branch := domain.Branch{
		Name:              "Taller Norte",
		RepresentativeID:  "rep-123",
		EstablishmentType: domain.EstablishmentTypeWorkshop,
		Brands:            []string{"brand-1", "brand-2"},
		Location: &domain.Location{
			Address:      "Calle 100 #15-20",
			CityID:       "city-bogota",
			DepartmentID: "dept-cundinamarca",
		},
	}

	savedBranch := &domain.Branch{
		ID:                "branch-123",
		Name:              branch.Name,
		RepresentativeID:  branch.RepresentativeID,
		EstablishmentType: branch.EstablishmentType,
	}

	// Mock expectations
	mockBranchService.On("ValidateBrands", ctx, branch.Brands).Return(nil)
	// Mock GeocodeLocation to set coordinates on the location object
	mockBranchService.On("GeocodeLocation", ctx, mock.AnythingOfType("*domain.Location")).Run(func(args mock.Arguments) {
		loc := args.Get(1).(*domain.Location)
		loc.Latitude = &lat
		loc.Longitude = &lng
	}).Return(true, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("RegisterBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(savedBranch, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, geocodingSucceeded, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "branch-123", result.ID)
	assert.True(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestRegisterBranch_InvalidBrands(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branch := domain.Branch{
		Name:   "Taller Norte",
		Brands: []string{"invalid-brand"},
	}

	// Mock expectations
	mockBranchService.On("ValidateBrands", ctx, branch.Brands).Return(domain.ErrBrandNotFound)

	// Act
	result, geocodingSucceeded, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)
	assert.Equal(t, domain.ErrBrandNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestRegisterBranch_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branch := domain.Branch{
		Name: "Taller Norte",
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockBranchService.On("BeginTx", ctx).Return(nil, txError)

	// Act
	result, geocodingSucceeded, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
}

func TestRegisterBranch_SaveError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branch := domain.Branch{
		Name: "Taller Norte",
	}

	saveError := errors.New("database error")

	// Mock expectations
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("RegisterBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil, saveError)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, _, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, saveError, err)

	mockTx.AssertCalled(t, "Rollback")
	mockBranchService.AssertExpectations(t)
}

func TestRegisterBranch_CommitError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branch := domain.Branch{
		Name: "Taller Norte",
	}

	savedBranch := &domain.Branch{
		ID:   "branch-123",
		Name: branch.Name,
	}

	// Mock expectations
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("RegisterBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(savedBranch, nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	// Act
	result, geocodingSucceeded, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)

	mockTx.AssertCalled(t, "Commit")
	mockBranchService.AssertExpectations(t)
}

func TestRegisterBranch_RollbackError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branch := domain.Branch{
		Name: "Taller Norte",
	}

	saveError := errors.New("database error")

	// Mock expectations
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("RegisterBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil, saveError)
	mockTx.On("Rollback").Return(errors.New("rollback failed")) // Rollback also fails

	// Act
	result, _, err := branchInteractor.RegisterBranch(ctx, branch)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, saveError, err)

	mockTx.AssertCalled(t, "Rollback")
	mockBranchService.AssertExpectations(t)
}

// ============================================
// GetBranchByID Tests
// ============================================

func TestGetBranchByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	expectedBranch := &domain.Branch{
		ID:                branchID,
		Name:              "Taller Centro",
		RepresentativeID:  "rep-123",
		EstablishmentType: domain.EstablishmentTypeWorkshop,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(expectedBranch, nil)

	// Act
	result, err := branchInteractor.GetBranchByID(ctx, branchID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, branchID, result.ID)
	assert.Equal(t, "Taller Centro", result.Name)

	mockBranchService.AssertExpectations(t)
}

func TestGetBranchByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "non-existent"

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	result, err := branchInteractor.GetBranchByID(ctx, branchID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// GetBranchesByRepresentative Tests (HU62)
// ============================================

func TestGetBranchesByRepresentative_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	representativeID := "rep-123"
	expectedBranches := []domain.Branch{
		{ID: "branch-1", Name: "Taller Norte", RepresentativeID: representativeID},
		{ID: "branch-2", Name: "Taller Sur", RepresentativeID: representativeID},
	}

	// Mock expectations
	mockBranchService.On("GetBranchesByRepresentative", ctx, representativeID).Return(expectedBranches, nil)

	// Act
	result, err := branchInteractor.GetBranchesByRepresentative(ctx, representativeID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockBranchService.AssertExpectations(t)
}

func TestGetBranchesByRepresentative_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	representativeID := "rep-new"

	// Mock expectations
	mockBranchService.On("GetBranchesByRepresentative", ctx, representativeID).Return([]domain.Branch{}, nil)

	// Act
	result, err := branchInteractor.GetBranchesByRepresentative(ctx, representativeID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockBranchService.AssertExpectations(t)
}

func TestGetBranchesByRepresentative_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	representativeID := "rep-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockBranchService.On("GetBranchesByRepresentative", ctx, representativeID).Return(nil, dbError)

	// Act
	result, err := branchInteractor.GetBranchesByRepresentative(ctx, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// UpdateBranch Tests (HU60)
// ============================================

func TestUpdateBranch_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:                branchID,
		Name:              "Taller Original",
		RepresentativeID:  personID,
		Status:            "ACTIVE",
		EstablishmentType: domain.EstablishmentTypeWorkshop,
	}

	updatedData := domain.Branch{
		Name:              "Taller Actualizado",
		EstablishmentType: domain.EstablishmentTypeStore,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil).Once()
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("UpdateBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	// Second GetBranchByID for refetch
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:                branchID,
		Name:              "Taller Actualizado",
		RepresentativeID:  personID,
		Status:            "ACTIVE",
		EstablishmentType: domain.EstablishmentTypeStore,
	}, nil).Once()

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, updatedData, personID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Taller Actualizado", result.Name)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateBranch_NotOwned(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-other" // Different from owner

	existingBranch := &domain.Branch{
		ID:               branchID,
		Name:             "Taller Original",
		RepresentativeID: "rep-owner", // Original owner
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, domain.Branch{}, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateBranch_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "non-existent"
	personID := "rep-123"

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, domain.Branch{}, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateBranch_InvalidBrands(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	updateData := domain.Branch{
		Brands: []string{"invalid-brand-id"},
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("ValidateBrands", ctx, updateData.Brands).Return(domain.ErrBrandNotFound)

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, updateData, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)
	assert.Equal(t, domain.ErrBrandNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateBranch_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, domain.Branch{}, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateBranch_ServiceError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
		Status:           "ACTIVE",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil).Once()
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("UpdateBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(errors.New("update failed"))
	mockTx.On("Rollback").Return(nil)

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, domain.Branch{}, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Rollback")
}

func TestUpdateBranch_CommitError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
		Status:           "ACTIVE",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil).Once()
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("UpdateBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, domain.Branch{}, personID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateBranch_RefetchError_ReturnsOriginal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
		Status:           "ACTIVE",
	}

	updateData := domain.Branch{
		Name: "Updated Name",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil).Once()
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("UpdateBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	// Second GetBranchByID for refetch fails
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, errors.New("refetch failed")).Once()

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, updateData, personID)

	// Assert
	assert.NoError(t, err) // Update was successful, refetch failure is handled gracefully
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.Name)
	assert.False(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
}

func TestUpdateBranch_WithLocation_Geocoding(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
		Status:           "ACTIVE",
	}

	updateData := domain.Branch{
		Name: "Updated Name",
		Location: &domain.Location{
			Address: "Calle 100 #15-20, Bogotá",
		},
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil).Once()
	mockBranchService.On("GeocodeLocation", ctx, updateData.Location).Return(true, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("UpdateBranch", ctx, mockTx, mock.AnythingOfType("domain.Branch")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	// Second GetBranchByID for refetch
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(&domain.Branch{
		ID:               branchID,
		Name:             "Updated Name",
		RepresentativeID: personID,
		Status:           "ACTIVE",
	}, nil).Once()

	// Act
	result, geocodingSucceeded, err := branchInteractor.UpdateBranch(ctx, branchID, updateData, personID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, geocodingSucceeded)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeleteBranch Tests (HU61)
// ============================================

func TestDeleteBranch_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		Name:             "Taller a Eliminar",
		RepresentativeID: personID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", ctx, mockTx, branchID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.NoError(t, err)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteBranch_NotOwned(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-other" // Different from owner

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: "rep-owner",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)

	mockBranchService.AssertExpectations(t)
}

func TestDeleteBranch_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "non-existent"
	personID := "rep-123"

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(nil, domain.ErrBranchNotFound)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestDeleteBranch_HasAssociations(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-with-diagnostics"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	// The service now interprets FK constraint errors internally and returns domain error
	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", ctx, mockTx, branchID).Return(domain.ErrBranchCannotDelete)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchCannotDelete, err)

	mockBranchService.AssertExpectations(t)
}

func TestDeleteBranch_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)

	mockBranchService.AssertExpectations(t)
}

func TestDeleteBranch_CommitError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", ctx, mockTx, branchID).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)

	mockBranchService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteBranch_ServiceError_NonFK(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	mockTx := new(mocks.MockTx)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	branchID := "branch-123"
	personID := "rep-123"

	existingBranch := &domain.Branch{
		ID:               branchID,
		RepresentativeID: personID,
	}

	// Simulate a non-FK error (e.g., database connection error)
	serviceError := errors.New("database connection lost")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, branchID).Return(existingBranch, nil)
	mockBranchService.On("BeginTx", ctx).Return(mockTx, nil)
	mockBranchService.On("DeleteBranch", ctx, mockTx, branchID).Return(serviceError)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := branchInteractor.DeleteBranch(ctx, branchID, personID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, serviceError, err) // Should return the original error, not ErrBranchCannotDelete

	mockBranchService.AssertExpectations(t)
}

// ============================================
// GeocodeLocation Tests
// ============================================

func TestGeocodeLocation_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	lat := 4.710989
	lng := -74.072092
	location := &domain.Location{
		Address: "Calle 100 #15-20",
	}

	// Mock expectations - location is modified in-place with coordinates
	mockBranchService.On("GeocodeLocation", ctx, location).Run(func(args mock.Arguments) {
		loc := args.Get(1).(*domain.Location)
		loc.Latitude = &lat
		loc.Longitude = &lng
	}).Return(true, nil)

	// Act
	success, err := branchInteractor.GeocodeLocation(ctx, location)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
	assert.NotNil(t, location.Latitude)
	assert.NotNil(t, location.Longitude)

	mockBranchService.AssertExpectations(t)
}

func TestGeocodeLocation_FailsGracefully(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	location := &domain.Location{
		Address: "Unknown Address",
	}

	geocodeError := errors.New("geocoding service unavailable")

	// Mock expectations
	mockBranchService.On("GeocodeLocation", ctx, location).Return(false, geocodeError)

	// Act
	success, err := branchInteractor.GeocodeLocation(ctx, location)

	// Assert
	assert.Error(t, err)
	assert.False(t, success)

	mockBranchService.AssertExpectations(t)
}

// ============================================
// GetBranchesNearby Tests (HU89)
// ============================================

func TestGetBranchesNearby_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	lat := 4.710989
	lng := -74.072092
	radiusKm := 5.0
	establishmentType := "WORKSHOP"

	expectedBranches := []domain.NearbyBranch{
		{
			ID:         "branch-1",
			Name:       "Taller Cercano",
			DistanceKm: 1.5,
		},
		{
			ID:         "branch-2",
			Name:       "Taller Zona Sur",
			DistanceKm: 3.2,
		},
	}

	// Mock expectations
	mockBranchService.On("GetBranchesNearby", ctx, lat, lng, radiusKm, establishmentType).Return(expectedBranches, nil)

	// Act
	result, err := branchInteractor.GetBranchesNearby(ctx, lat, lng, radiusKm, establishmentType)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "branch-1", result[0].ID)
	assert.Equal(t, 1.5, result[0].DistanceKm)

	mockBranchService.AssertExpectations(t)
}

func TestGetBranchesNearby_EmptyResults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	lat := 0.0 // Middle of nowhere
	lng := 0.0
	radiusKm := 1.0
	establishmentType := "WORKSHOP"

	// Mock expectations - no branches found
	mockBranchService.On("GetBranchesNearby", ctx, lat, lng, radiusKm, establishmentType).Return([]domain.NearbyBranch{}, nil)

	// Act
	result, err := branchInteractor.GetBranchesNearby(ctx, lat, lng, radiusKm, establishmentType)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockBranchService.AssertExpectations(t)
}

func TestGetBranchesNearby_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBranchService := new(mocks.MockBranchService)

	branchInteractor := interactor.NewBranchInteractor(mockBranchService)

	lat := 4.710989
	lng := -74.072092
	radiusKm := 5.0
	establishmentType := "WORKSHOP"

	dbError := errors.New("database error")

	// Mock expectations
	mockBranchService.On("GetBranchesNearby", ctx, lat, lng, radiusKm, establishmentType).Return(nil, dbError)

	// Act
	result, err := branchInteractor.GetBranchesNearby(ctx, lat, lng, radiusKm, establishmentType)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockBranchService.AssertExpectations(t)
}
