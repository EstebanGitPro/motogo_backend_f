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
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{
		Name:        "Motozona Colombia",
		Description: stringPtr("Franquicia de motos"),
	}
	branchIDs := []string{"branch-1", "branch-2"}
	representativeID := "rep-123"

	createdFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: franchise.Name,
	}

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(createdFranchise, nil)
	mockFranchiseService.On("AssociateBranches", ctx, mockTx, "franchise-123", branchIDs).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "franchise-123", result.ID)

	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestCreateFranchiseWithBranches_NoBranches(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{}
	representativeID := "rep-123"

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(domain.ErrFranchiseNoBranches)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseNoBranches, err)
}

func TestCreateFranchiseWithBranches_BranchNotFound(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"non-existent-branch"}
	representativeID := "rep-123"

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(domain.ErrBranchNotFound)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestCreateFranchiseWithBranches_BranchNotOwned(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(domain.ErrFranchiseBranchNotOwned)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)
}

func TestCreateFranchiseWithBranches_BranchAlreadyInFranchise(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(domain.ErrFranchiseBranchNotOwned)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)
}

func TestCreateFranchiseWithBranches_CreateFails_Rollback(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{Name: "Motozona"}
	branchIDs := []string{"branch-1"}
	representativeID := "rep-123"

	createError := errors.New("create failed")

	mockFranchiseService.On("ValidateBranchesForFranchise", ctx, branchIDs, representativeID).Return(nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("CreateFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(nil, createError)
	mockTx.On("Rollback").Return(nil)

	result, err := franchiseInteractor.CreateFranchiseWithBranches(ctx, franchise, branchIDs, representativeID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, createError, err)

	mockTx.AssertCalled(t, "Rollback")
	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// DeleteFranchise Tests
// ============================================

func TestDeleteFranchise_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchiseID := "franchise-123"
	representativeID := "rep-123"

	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("DissociateBranches", ctx, mockTx, franchiseID).Return(nil)
	mockFranchiseService.On("DeleteFranchise", ctx, mockTx, franchiseID).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := franchiseInteractor.DeleteFranchise(ctx, franchiseID, representativeID)

	assert.NoError(t, err)

	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteFranchise_NotFound(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchiseID := "non-existent"
	representativeID := "rep-123"

	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("DissociateBranches", ctx, mockTx, franchiseID).Return(nil)
	mockFranchiseService.On("DeleteFranchise", ctx, mockTx, franchiseID).Return(domain.ErrFranchiseNotFound)
	mockTx.On("Rollback").Return(nil)

	err := franchiseInteractor.DeleteFranchise(ctx, franchiseID, representativeID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// GetFranchiseByID Tests
// ============================================

func TestGetFranchiseByID_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	expectedFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona",
	}

	mockFranchiseService.On("GetFranchiseByID", ctx, "franchise-123").Return(expectedFranchise, nil)

	result, err := franchiseInteractor.GetFranchiseByID(ctx, "franchise-123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "franchise-123", result.ID)

	mockFranchiseService.AssertExpectations(t)
}

func TestGetFranchiseByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	mockFranchiseService.On("GetFranchiseByID", ctx, "not-found").Return(nil, domain.ErrFranchiseNotFound)

	result, err := franchiseInteractor.GetFranchiseByID(ctx, "not-found")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockFranchiseService.AssertExpectations(t)
}

// ============================================
// GetFranchisesByRepresentative Tests
// ============================================

func TestGetFranchisesByRepresentative_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	expectedFranchises := []domain.Franchise{
		{ID: "f1", Name: "Motozona Norte"},
		{ID: "f2", Name: "Motozona Sur"},
	}

	mockFranchiseService.On("GetFranchisesByRepresentative", ctx, "rep-123").Return(expectedFranchises, nil)

	result, err := franchiseInteractor.GetFranchisesByRepresentative(ctx, "rep-123")

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

// ============================================
// UpdateFranchise Tests
// ============================================

func TestUpdateFranchise_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{ID: "franchise-123", Name: "Updated Name"}
	representativeID := "rep-123"

	mockFranchiseService.On("CountBranches", ctx, "franchise-123").Return(2, nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("UpdateFranchise", ctx, mockTx, franchise).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	err := franchiseInteractor.UpdateFranchise(ctx, franchise, representativeID)

	assert.NoError(t, err)
	mockFranchiseService.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateFranchise_NotFound(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{ID: "non-existent", Name: "Name"}
	representativeID := "rep-123"

	mockFranchiseService.On("CountBranches", ctx, "non-existent").Return(0, nil)

	err := franchiseInteractor.UpdateFranchise(ctx, franchise, representativeID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)
	mockFranchiseService.AssertExpectations(t)
}

func TestUpdateFranchise_CountBranchesError(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchise := domain.Franchise{ID: "franchise-123", Name: "Name"}
	representativeID := "rep-123"

	dbError := errors.New("database error")
	mockFranchiseService.On("CountBranches", ctx, "franchise-123").Return(0, dbError)

	err := franchiseInteractor.UpdateFranchise(ctx, franchise, representativeID)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
}

// ============================================
// AddBranchToFranchise Tests
// ============================================

func TestAddBranchToFranchise_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchiseID := "franchise-123"
	branchID := "branch-456"
	representativeID := "rep-123"

	mockFranchiseService.On("ValidateBranchOwnership", ctx, branchID, representativeID).Return(nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("AssociateBranches", ctx, mockTx, franchiseID, []string{branchID}).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	err := franchiseInteractor.AddBranchToFranchise(ctx, franchiseID, branchID, representativeID)

	assert.NoError(t, err)
	mockFranchiseService.AssertExpectations(t)
}

func TestAddBranchToFranchise_BranchNotFound(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	mockFranchiseService.On("ValidateBranchOwnership", ctx, "non-existent", "rep-123").Return(domain.ErrBranchNotFound)

	err := franchiseInteractor.AddBranchToFranchise(ctx, "franchise-123", "non-existent", "rep-123")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestAddBranchToFranchise_BranchNotOwned(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	mockFranchiseService.On("ValidateBranchOwnership", ctx, "branch-1", "rep-123").Return(domain.ErrFranchiseBranchNotOwned)

	err := franchiseInteractor.AddBranchToFranchise(ctx, "franchise-123", "branch-1", "rep-123")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)
}

// ============================================
// RemoveBranchFromFranchise Tests
// ============================================

func TestRemoveBranchFromFranchise_Success(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)
	mockTx := new(mocks.MockTx)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchiseID := "franchise-123"
	branchID := "branch-456"
	representativeID := "rep-123"

	mockFranchiseService.On("CanRemoveBranch", ctx, franchiseID).Return(nil)
	mockFranchiseService.On("BeginTx", ctx).Return(mockTx, nil)
	mockFranchiseService.On("DissociateSingleBranch", ctx, mockTx, branchID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	err := franchiseInteractor.RemoveBranchFromFranchise(ctx, franchiseID, branchID, representativeID)

	assert.NoError(t, err)
}

func TestRemoveBranchFromFranchise_LastBranch(t *testing.T) {
	ctx := context.Background()
	mockFranchiseService := new(mocks.MockFranchiseService)

	franchiseInteractor := interactor.NewFranchiseInteractor(mockFranchiseService)

	franchiseID := "franchise-123"
	branchID := "branch-456"
	representativeID := "rep-123"

	mockFranchiseService.On("CanRemoveBranch", ctx, franchiseID).Return(domain.ErrFranchiseMinBranches)

	err := franchiseInteractor.RemoveBranchFromFranchise(ctx, franchiseID, branchID, representativeID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseMinBranches, err)
}
