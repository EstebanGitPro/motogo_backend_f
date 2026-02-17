package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
)

// ============================================
// BeginTx Tests (franchise)
// ============================================

func TestFranchiseSvc_BeginTx_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	tx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("BeginTx", ctx).Return(tx, nil)

	result, err := service.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestFranchiseSvc_BeginTx_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("BeginTx", ctx).Return(nil, errors.New("db error"))

	result, err := service.BeginTx(ctx)
	assert.Nil(t, result)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============================================
// DissociateSingleBranch Tests
// ============================================

func TestFranchiseSvc_DissociateSingleBranch_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("DissociateSingleBranch", ctx, mockTx, "branch-1").Return(nil)

	err := service.DissociateSingleBranch(ctx, mockTx, "branch-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFranchiseSvc_DissociateSingleBranch_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("DissociateSingleBranch", ctx, mockTx, "bad").Return(errors.New("db error"))

	err := service.DissociateSingleBranch(ctx, mockTx, "bad")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============================================
// CanRemoveBranch Tests
// ============================================

func TestFranchiseSvc_CanRemoveBranch_OK(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("CountBranchesByFranchise", ctx, "franchise-1").Return(3, nil)

	err := service.CanRemoveBranch(ctx, "franchise-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFranchiseSvc_CanRemoveBranch_MinBranches(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("CountBranchesByFranchise", ctx, "franchise-1").Return(1, nil)

	err := service.CanRemoveBranch(ctx, "franchise-1")
	assert.Equal(t, domain.ErrFranchiseMinBranches, err)
	mockRepo.AssertExpectations(t)
}

func TestFranchiseSvc_CanRemoveBranch_CountError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("CountBranchesByFranchise", ctx, "franchise-1").Return(0, errors.New("db error"))

	err := service.CanRemoveBranch(ctx, "franchise-1")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============================================
// ValidateBranchOwnership Tests
// ============================================

func TestFranchiseSvc_ValidateBranchOwnership_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	branch := &domain.Branch{ID: "branch-1", RepresentativeID: "rep-1", FranchiseID: nil}
	mockBranchRepo.On("GetBranchByID", ctx, "branch-1").Return(branch, nil)

	err := service.ValidateBranchOwnership(ctx, "branch-1", "rep-1")
	assert.NoError(t, err)
	mockBranchRepo.AssertExpectations(t)
}

func TestFranchiseSvc_ValidateBranchOwnership_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	mockBranchRepo.On("GetBranchByID", ctx, "bad-branch").Return(nil, errors.New("not found"))

	err := service.ValidateBranchOwnership(ctx, "bad-branch", "rep-1")
	assert.Equal(t, domain.ErrBranchNotFound, err)
	mockBranchRepo.AssertExpectations(t)
}

func TestFranchiseSvc_ValidateBranchOwnership_WrongOwner(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	branch := &domain.Branch{ID: "branch-1", RepresentativeID: "rep-other", FranchiseID: nil}
	mockBranchRepo.On("GetBranchByID", ctx, "branch-1").Return(branch, nil)

	err := service.ValidateBranchOwnership(ctx, "branch-1", "rep-1")
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)
	mockBranchRepo.AssertExpectations(t)
}

func TestFranchiseSvc_ValidateBranchOwnership_AlreadyInFranchise(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	franchiseID := "existing-franchise"
	branch := &domain.Branch{ID: "branch-1", RepresentativeID: "rep-1", FranchiseID: &franchiseID}
	mockBranchRepo.On("GetBranchByID", ctx, "branch-1").Return(branch, nil)

	err := service.ValidateBranchOwnership(ctx, "branch-1", "rep-1")
	assert.Equal(t, domain.ErrFranchiseBranchNotOwned, err)
	mockBranchRepo.AssertExpectations(t)
}

// ============================================
// ValidateBranchesForFranchise Tests
// ============================================

func TestFranchiseSvc_ValidateBranchesForFranchise_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	branch1 := &domain.Branch{ID: "b1", RepresentativeID: "rep-1", FranchiseID: nil}
	branch2 := &domain.Branch{ID: "b2", RepresentativeID: "rep-1", FranchiseID: nil}
	mockBranchRepo.On("GetBranchByID", ctx, "b1").Return(branch1, nil)
	mockBranchRepo.On("GetBranchByID", ctx, "b2").Return(branch2, nil)

	err := service.ValidateBranchesForFranchise(ctx, []string{"b1", "b2"}, "rep-1")
	assert.NoError(t, err)
	mockBranchRepo.AssertExpectations(t)
}

func TestFranchiseSvc_ValidateBranchesForFranchise_EmptyBranches(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	err := service.ValidateBranchesForFranchise(ctx, []string{}, "rep-1")
	assert.Equal(t, domain.ErrFranchiseNoBranches, err)
}

func TestFranchiseSvc_ValidateBranchesForFranchise_OneInvalid(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	service := services.NewFranchiseService(mockRepo, mockBranchRepo)

	branch1 := &domain.Branch{ID: "b1", RepresentativeID: "rep-1", FranchiseID: nil}
	mockBranchRepo.On("GetBranchByID", ctx, "b1").Return(branch1, nil)
	mockBranchRepo.On("GetBranchByID", ctx, "b2").Return(nil, errors.New("not found"))

	err := service.ValidateBranchesForFranchise(ctx, []string{"b1", "b2"}, "rep-1")
	assert.Equal(t, domain.ErrBranchNotFound, err)
	mockBranchRepo.AssertExpectations(t)
}

// ============================================
// CreateFranchise - NameLookupError
// ============================================

func TestFranchiseSvc_CreateFranchise_NameLookupError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{Name: "Test"}
	mockRepo.On("GetFranchiseByName", ctx, "Test").Return(nil, errors.New("db error"))

	result, err := service.CreateFranchise(ctx, mockTx, franchise)
	assert.Nil(t, result)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============================================
// UpdateFranchise - ID exists but gets nil
// ============================================

func TestFranchiseSvc_UpdateFranchise_ExistsButNil(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{ID: "f-1", Name: "Test"}
	mockRepo.On("GetFranchiseByID", ctx, "f-1").Return((*domain.Franchise)(nil), nil)

	err := service.UpdateFranchise(ctx, mockTx, franchise)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)
}

func TestFranchiseSvc_UpdateFranchise_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	existing := &domain.Franchise{ID: "f-1", Name: "Old"}
	franchise := domain.Franchise{ID: "f-1", Name: "New"}
	mockRepo.On("GetFranchiseByID", ctx, "f-1").Return(existing, nil)
	mockRepo.On("GetFranchiseByName", ctx, "New").Return(nil, nil)
	mockRepo.On("UpdateFranchise", ctx, mockTx, franchise).Return(errors.New("db error"))

	err := service.UpdateFranchise(ctx, mockTx, franchise)
	assert.Equal(t, domain.ErrFranchiseCannotUpdate, err)
}

func TestFranchiseSvc_UpdateFranchise_NameCheckError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	existing := &domain.Franchise{ID: "f-1", Name: "Old"}
	franchise := domain.Franchise{ID: "f-1", Name: "New"}
	mockRepo.On("GetFranchiseByID", ctx, "f-1").Return(existing, nil)
	mockRepo.On("GetFranchiseByName", ctx, "New").Return(nil, errors.New("db error"))

	err := service.UpdateFranchise(ctx, mockTx, franchise)
	assert.Error(t, err)
}

// ============================================
// DeleteFranchise - nil existing
// ============================================

func TestFranchiseSvc_DeleteFranchise_ExistsButNil(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	mockTx := new(mocks.MockTx)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("GetFranchiseByID", ctx, "f-1").Return((*domain.Franchise)(nil), nil)

	err := service.DeleteFranchise(ctx, mockTx, "f-1")
	assert.Equal(t, domain.ErrFranchiseNotFound, err)
}

// ============================================
// GetFranchiseByID - DB error (not NotFound)
// ============================================

func TestFranchiseSvc_GetFranchiseByID_DBError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)
	service := services.NewFranchiseService(mockRepo, nil)

	mockRepo.On("GetFranchiseByID", ctx, "f-1").Return(nil, errors.New("connection error"))

	result, err := service.GetFranchiseByID(ctx, "f-1")
	assert.Nil(t, result)
	assert.Error(t, err)
}

// mockFranchiseTx ensures BeginTx returns proper types
var _ output.Tx = (*mocks.MockTx)(nil)
