package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// CreateFranchise Tests
// ============================================

func TestCreateFranchise_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{
		Name:        "Motozona Colombia",
		Description: stringPtr("Franquicia de motos"),
	}

	// Mock expectations - no duplicate name
	mockRepo.On("GetFranchiseByName", ctx, franchise.Name).Return(nil, nil)
	mockRepo.On("SaveFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(nil)

	// Act
	result, err := service.CreateFranchise(ctx, mockTx, franchise)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, franchise.Name, result.Name)
	assert.NotEmpty(t, result.ID)

	mockRepo.AssertExpectations(t)
}

func TestCreateFranchise_DuplicateName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{
		Name: "Motozona Colombia",
	}

	existingFranchise := &domain.Franchise{
		ID:   "existing-id",
		Name: franchise.Name,
	}

	// Mock expectations - name already exists
	mockRepo.On("GetFranchiseByName", ctx, franchise.Name).Return(existingFranchise, nil)

	// Act
	result, err := service.CreateFranchise(ctx, mockTx, franchise)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseDuplicateName, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateFranchise_SaveError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{
		Name: "Motozona Colombia",
	}

	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetFranchiseByName", ctx, franchise.Name).Return(nil, nil)
	mockRepo.On("SaveFranchise", ctx, mockTx, mock.AnythingOfType("domain.Franchise")).Return(dbError)

	// Act
	result, err := service.CreateFranchise(ctx, mockTx, franchise)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseCannotSave, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// GetFranchiseByID Tests
// ============================================

func TestGetFranchiseByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	service := services.NewFranchiseService(mockRepo, nil)

	expectedFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Colombia",
	}

	// Mock expectations
	mockRepo.On("GetFranchiseByID", ctx, "franchise-123").Return(expectedFranchise, nil)

	// Act
	result, err := service.GetFranchiseByID(ctx, "franchise-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedFranchise.ID, result.ID)
	assert.Equal(t, expectedFranchise.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetFranchiseByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	service := services.NewFranchiseService(mockRepo, nil)

	// Mock expectations - not found
	mockRepo.On("GetFranchiseByID", ctx, "not-found-id").Return(nil, domain.ErrFranchiseNotFound)

	// Act
	result, err := service.GetFranchiseByID(ctx, "not-found-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// GetFranchisesByRepresentative Tests
// ============================================

func TestGetFranchisesByRepresentative_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	service := services.NewFranchiseService(mockRepo, nil)

	expectedFranchises := []domain.Franchise{
		{ID: "f1", Name: "Motozona Norte"},
		{ID: "f2", Name: "Motozona Sur"},
	}

	// Mock expectations
	mockRepo.On("GetFranchisesByRepresentative", ctx, "rep-123").Return(expectedFranchises, nil)

	// Act
	result, err := service.GetFranchisesByRepresentative(ctx, "rep-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Motozona Norte", result[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetFranchisesByRepresentative_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	service := services.NewFranchiseService(mockRepo, nil)

	// Mock expectations - no franchises
	mockRepo.On("GetFranchisesByRepresentative", ctx, "rep-no-franchises").Return([]domain.Franchise{}, nil)

	// Act
	result, err := service.GetFranchisesByRepresentative(ctx, "rep-no-franchises")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// UpdateFranchise Tests
// ============================================

func TestUpdateFranchise_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	existingFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Original",
	}

	updatedFranchise := domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Updated",
	}

	// Mock expectations
	mockRepo.On("GetFranchiseByID", ctx, updatedFranchise.ID).Return(existingFranchise, nil)
	mockRepo.On("GetFranchiseByName", ctx, updatedFranchise.Name).Return(nil, nil)
	mockRepo.On("UpdateFranchise", ctx, mockTx, updatedFranchise).Return(nil)

	// Act
	err := service.UpdateFranchise(ctx, mockTx, updatedFranchise)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateFranchise_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	franchise := domain.Franchise{
		ID:   "not-found-id",
		Name: "Motozona",
	}

	// Mock expectations - not found
	mockRepo.On("GetFranchiseByID", ctx, franchise.ID).Return(nil, domain.ErrFranchiseNotFound)

	// Act
	err := service.UpdateFranchise(ctx, mockTx, franchise)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateFranchise_DuplicateName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	existingFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Original",
	}

	updatedFranchise := domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Duplicada",
	}

	anotherFranchise := &domain.Franchise{
		ID:   "franchise-456",
		Name: "Motozona Duplicada",
	}

	// Mock expectations - trying to use a name that already exists
	mockRepo.On("GetFranchiseByID", ctx, updatedFranchise.ID).Return(existingFranchise, nil)
	mockRepo.On("GetFranchiseByName", ctx, updatedFranchise.Name).Return(anotherFranchise, nil)

	// Act
	err := service.UpdateFranchise(ctx, mockTx, updatedFranchise)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseDuplicateName, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateFranchise_SameName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	existingFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona Colombia",
	}

	updatedFranchise := domain.Franchise{
		ID:          "franchise-123",
		Name:        "Motozona Colombia", // Same name, just updating description
		Description: stringPtr("Nueva descripción"),
	}

	// Mock expectations - same name, no duplicate check needed
	mockRepo.On("GetFranchiseByID", ctx, updatedFranchise.ID).Return(existingFranchise, nil)
	mockRepo.On("UpdateFranchise", ctx, mockTx, updatedFranchise).Return(nil)

	// Act
	err := service.UpdateFranchise(ctx, mockTx, updatedFranchise)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// DeleteFranchise Tests
// ============================================

func TestDeleteFranchise_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	existingFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona",
	}

	// Mock expectations
	mockRepo.On("GetFranchiseByID", ctx, "franchise-123").Return(existingFranchise, nil)
	mockRepo.On("DeleteFranchise", ctx, mockTx, "franchise-123").Return(nil)

	// Act
	err := service.DeleteFranchise(ctx, mockTx, "franchise-123")

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteFranchise_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	// Mock expectations - not found
	mockRepo.On("GetFranchiseByID", ctx, "not-found-id").Return(nil, domain.ErrFranchiseNotFound)

	// Act
	err := service.DeleteFranchise(ctx, mockTx, "not-found-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteFranchise_DeleteError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	existingFranchise := &domain.Franchise{
		ID:   "franchise-123",
		Name: "Motozona",
	}

	dbError := errors.New("delete error")

	// Mock expectations
	mockRepo.On("GetFranchiseByID", ctx, "franchise-123").Return(existingFranchise, nil)
	mockRepo.On("DeleteFranchise", ctx, mockTx, "franchise-123").Return(dbError)

	// Act
	err := service.DeleteFranchise(ctx, mockTx, "franchise-123")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseCannotDelete, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// AssociateBranches Tests
// ============================================

func TestAssociateBranches_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	branchIDs := []string{"branch-1", "branch-2"}

	// Mock expectations
	mockRepo.On("AssociateBranchesToFranchise", ctx, mockTx, "franchise-123", branchIDs).Return(nil)

	// Act
	err := service.AssociateBranches(ctx, mockTx, "franchise-123", branchIDs)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAssociateBranches_NoBranches(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	branchIDs := []string{} // Empty

	// Act
	err := service.AssociateBranches(ctx, mockTx, "franchise-123", branchIDs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrFranchiseNoBranches, err)
}

// ============================================
// DissociateBranches Tests
// ============================================

func TestDissociateBranches_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	mockTx := new(mocks.MockTx)

	service := services.NewFranchiseService(mockRepo, nil)

	// Mock expectations
	mockRepo.On("DissociateBranchesFromFranchise", ctx, mockTx, "franchise-123").Return(nil)

	// Act
	err := service.DissociateBranches(ctx, mockTx, "franchise-123")

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// CountBranches Tests
// ============================================

func TestCountBranches_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockFranchiseRepository)

	service := services.NewFranchiseService(mockRepo, nil)

	// Mock expectations
	mockRepo.On("CountBranchesByFranchise", ctx, "franchise-123").Return(5, nil)

	// Act
	count, err := service.CountBranches(ctx, "franchise-123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	mockRepo.AssertExpectations(t)
}

// ============================================
// Helper Functions
// ============================================

func stringPtr(s string) *string {
	return &s
}
