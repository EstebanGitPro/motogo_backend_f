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
// CreateFranchiseWithBranches Tests
// ============================================

func TestCreateFranchiseWithBranches_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{
		Name:        "Motozona Colombia",
		Description: stringPtr("Franquicia de motos"),
	}
	branchIDs := []string{"branch-1", "branch-2"}
	representativeID := "rep-123"

	branch1 := &domain.Branch{ID: "branch-1", RepresentativeID: representativeID}
	branch2 := &domain.Branch{ID: "branch-2", RepresentativeID: representativeID}

	createdFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: franchise.Name,
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, "branch-1").Return(branch1, nil)
	mockBranchService.On("GetBranchByID", ctx, "branch-2").Return(branch2, nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(createdFranchise, nil)
	mockFranchiseService.On("AssociateBranches", ctx, mockTx, "franchise-123", branchIDs).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "franchise-123", result.ID)

	mockFranchiseService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestCreateFranchiseWithBranches_NoBranches(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{} // Empty - should fail
	representativeID := "rep-123"

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseNoBranches, err)
}

func TestCreateFranchiseWithBranches_BranchNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"non-existent-branch"}
	representativeID := "rep-123"

	// Mock expectations - branch not found
	mockBranchService.On("GetBranchByID", ctx, "non-existent-branch").Return(nil, domain.ErrBranchNotFound)

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateFranchiseWithBranches_BranchNotOwned(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	// Branch belongs to different representative
	otherRepBranch := &domain.Branch{
		ID:               "branch-1",
		RepresentativeID: "other-rep",
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, "branch-1").Return(otherRepBranch, nil)

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateFranchiseWithBranches_BranchAlreadyInFranchise(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	existingFranchiseID := "another-franchise"
	branchWithFranchise := &domain.Branch{
		ID:               "branch-1",
		RepresentativeID: representativeID,
		FranchiseID:      &existingFranchiseID, // Already in franchise
	}

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, "branch-1").Return(branchWithFranchise, nil)

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)

	mockBranchService.AssertExpectations(t)
}

func TestCreateFranchiseWithBranches_CreateFails_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	branch1 := &domain.Branch{ID: "branch-1", RepresentativeID: representativeID}
	createError := errors.New("create failed")

	// Mock expectations
	mockBranchService.On("GetBranchByID", ctx, "branch-1").Return(branch1, nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(nil, createError)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, createError, err)

	mockTx.AssertCalled(t, "Rollback")
	mockFranchiseService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// ============================================
// DeleteFranchise Tests
// ============================================

func TestDeleteFranchise_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchiseID := "franchise-123"
	representativeID := "rep-123"

	// Mock expectations
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("DissociateBranches", ctx, mockTx, franchiseID).Return(nil)
	mockFranchiseService.On("DeleteFranchise", ctx, mockTx, franchiseID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := franchiseInteractor.DeleteFranchise(ctx, franchiseID, representativeID)

	// Assert
	assert.NoError(t, err)

	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteFranchise_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	franchiseID := "non-existent"
	representativeID := "rep-123"

	// Mock expectations - service DeleteFranchise will check existence
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("DissociateBranches", ctx, mockTx, franchiseID).Return(nil)
	mockFranchiseService.On("DeleteFranchise", ctx, mockTx, franchiseID).Return(domain.ErrFranchiseNotFound)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := franchiseInteractor.DeleteFranchise(ctx, franchiseID, representativeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// GetFranchiseByID Tests
// ============================================

func TestGetFranchiseByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	expectedFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona",
	}

	// Mock expectations
	mockFranchiseService.On("GetFranchiseByID", ctx, "franchise-123").Return(expectedFranchise, nil)

	// Act
	result, err := franchiseInteractor.GetFranchiseByID(ctx, "franchise-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "franchise-123", result.ID)

	mockFranchiseService.AssertExpectations(t)
}

func TestGetFranchiseByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	// Mock expectations
	mockFranchiseService.On("GetFranchiseByID", ctx, "not-found").Return(nil, domain.ErrFranchiseNotFound)

	// Act
	result, err := franchiseInteractor.GetFranchiseByID(ctx, "not-found")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// GetFranchisesByRepresentative Tests
// ============================================

func TestGetFranchisesByRepresentative_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockBranchService := new(mocks.MockBranchService)
	mockLogger := new(mocks.MockLogger)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService, mockBranchService, mockLogger)

	expectedFranchises := []domain.Franchise{
		{ID: "f1", Name: "Motozona Norte"},
		{ID: "f2", Name: "Motozona Sur"},
	}

	// Mock expectations
	mockFranchiseService.On("GetFranchisesByRepresentative", ctx, "rep-123").Return(expectedFranchises, nil)

	// Act
	result, err := franchiseInteractor.GetFranchisesByRepresentative(ctx, "rep-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// Helper Functions
// ============================================

func stringPtr(s string) *string {
	return &s
}
