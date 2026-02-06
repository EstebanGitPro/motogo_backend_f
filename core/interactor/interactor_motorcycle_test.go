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
// RegisterMotorcycle Tests (HU43)
// ============================================

func TestRegisterMotorcycle_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	year := 2023
	mileage := 5000
	motorcycle := &domain.Motorcycle{
		LicensePlate:   "ABC123",
		OwnerID:        "owner-123",
		ReferenceID:    "ref-honda-cbr",
		Year:           &year,
		CurrentMileage: &mileage,
	}

	// Mock expectations
	mockRepo.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(true, nil)
	mockRepo.On("CheckLicensePlateExists", ctx, "ABC123").Return(false, nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("Save", ctx, mockTx, mock.AnythingOfType("*domain.Motorcycle")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID) // UUID should be generated

	mockRepo.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestRegisterMotorcycle_ReferenceRequired(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "", // Missing reference
	}

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceRequired, err)
}

func TestRegisterMotorcycle_ReferenceNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "invalid-ref",
	}

	// Mock expectations
	mockRepo.On("ValidateReferenceExists", ctx, "invalid-ref").Return(false, nil)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterMotorcycle_DuplicateLicensePlate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	// Mock expectations
	mockRepo.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(true, nil)
	mockRepo.On("CheckLicensePlateExists", ctx, "ABC123").Return(true, nil) // Already exists

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDuplicateLicensePlate, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterMotorcycle_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockRepo.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(true, nil)
	mockRepo.On("CheckLicensePlateExists", ctx, "ABC123").Return(false, nil)
	mockRepo.On("BeginTx", ctx).Return(nil, txError)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterMotorcycle_SaveError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	saveError := errors.New("database error")

	// Mock expectations
	mockRepo.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(true, nil)
	mockRepo.On("CheckLicensePlateExists", ctx, "ABC123").Return(false, nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("Save", ctx, mockTx, mock.AnythingOfType("*domain.Motorcycle")).Return(saveError)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)

	mockTx.AssertCalled(t, "Rollback")
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMotorcycleByID Tests (HU46)
// ============================================

func TestGetMotorcycleByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	expectedMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(expectedMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByID(ctx, motorcycleID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, motorcycleID, result.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetMotorcycleByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "non-existent"

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByID(ctx, motorcycleID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMotorcyclesByOwner Tests (HU47)
// ============================================

func TestGetMotorcyclesByOwner_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	ownerID := "owner-123"
	expectedMotorcycles := []domain.Motorcycle{
		{ID: "moto-1", LicensePlate: "ABC123", OwnerID: ownerID},
		{ID: "moto-2", LicensePlate: "DEF456", OwnerID: ownerID},
	}

	// Mock expectations
	mockRepo.On("GetByOwnerID", ctx, ownerID).Return(expectedMotorcycles, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockRepo.AssertExpectations(t)
}

func TestGetMotorcyclesByOwner_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	ownerID := "new-owner"

	// Mock expectations
	mockRepo.On("GetByOwnerID", ctx, ownerID).Return([]domain.Motorcycle{}, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestGetMotorcyclesByOwner_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	ownerID := "owner-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetByOwnerID", ctx, ownerID).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMotorcycleByLicensePlate Tests (HU47)
// ============================================

func TestGetMotorcycleByLicensePlate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	licensePlate := "ABC123"
	expectedMotorcycle := &domain.Motorcycle{
		ID:           "moto-123",
		LicensePlate: licensePlate,
		OwnerID:      "owner-123",
	}

	// Mock expectations
	mockRepo.On("GetByLicensePlate", ctx, licensePlate).Return(expectedMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByLicensePlate(ctx, licensePlate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, licensePlate, result.LicensePlate)

	mockRepo.AssertExpectations(t)
}

func TestGetMotorcycleByLicensePlate_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	licensePlate := "INVALID"

	// Mock expectations
	mockRepo.On("GetByLicensePlate", ctx, licensePlate).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByLicensePlate(ctx, licensePlate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// UpdateMotorcycle Tests (HU44)
// ============================================

func TestUpdateMotorcycle_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	year := 2023
	mileage := 10000

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
		ReferenceID:  "ref-honda-cbr",
	}

	updates := &domain.Motorcycle{
		Year:           &year,
		CurrentMileage: &mileage,
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil).Once()
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("Update", ctx, mockTx, mock.AnythingOfType("*domain.Motorcycle")).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil).Once()

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateMotorcycle_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "non-existent"
	ownerID := "owner-123"
	updates := &domain.Motorcycle{}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMotorcycle_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "other-owner" // Different from actual owner

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      "owner-123", // Actual owner
	}

	updates := &domain.Motorcycle{}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockRepo.AssertExpectations(t)
}

func TestUpdateMotorcycle_ReferenceNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
		ReferenceID:  "ref-honda-cbr",
	}

	updates := &domain.Motorcycle{
		ReferenceID: "invalid-ref", // New invalid reference
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockRepo.On("ValidateReferenceExists", ctx, "invalid-ref").Return(false, nil)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// DeleteMotorcycle Tests (HU45)
// ============================================

func TestDeleteMotorcycle_Success_SoftDelete(t *testing.T) {
	// Arrange - motorcycle WITH service history -> soft delete
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
	}

	// Mock expectations - has history -> soft delete
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockRepo.On("HasServiceHistory", ctx, motorcycleID).Return(true, nil) // Has history
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("Delete", ctx, mockTx, motorcycleID).Return(nil) // Soft delete
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteMotorcycle_Success_HardDelete(t *testing.T) {
	// Arrange - motorcycle WITHOUT service history -> hard delete
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-456"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "DEF456",
		OwnerID:      ownerID,
	}

	// Mock expectations - no history -> hard delete
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockRepo.On("HasServiceHistory", ctx, motorcycleID).Return(false, nil) // No history
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("HardDelete", ctx, mockTx, motorcycleID).Return(nil) // Hard delete
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteMotorcycle_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "non-existent"
	ownerID := "owner-123"

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteMotorcycle_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "other-owner" // Different from actual owner

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      "owner-123", // Actual owner
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockRepo.AssertExpectations(t)
}

func TestDeleteMotorcycle_TxError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	deleteError := errors.New("delete failed")

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockRepo.On("HasServiceHistory", ctx, motorcycleID).Return(true, nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("Delete", ctx, mockTx, motorcycleID).Return(deleteError)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotDelete, err)

	mockTx.AssertCalled(t, "Rollback")
	mockRepo.AssertExpectations(t)
}

// ============================================
// GetMotorcycleReferences Tests (HU50)
// ============================================

func TestGetMotorcycleReferences_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	expectedRefs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "CBR 600"},
		{ID: "ref-2", Model: "Ninja 650"},
	}

	// Mock expectations
	mockRepo.On("GetAllReferences", ctx).Return(expectedRefs, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleReferences(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockRepo.AssertExpectations(t)
}

func TestGetMotorcycleReferences_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetAllReferences", ctx).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleReferences(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// GetReferencesByBrandID Tests (HU40)
// ============================================

func TestGetReferencesByBrandID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	brandID := "brand-honda"
	expectedRefs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "CBR 600", BrandID: brandID},
		{ID: "ref-2", Model: "CBR 1000", BrandID: brandID},
	}

	// Mock expectations
	mockRepo.On("GetReferencesByBrandID", ctx, brandID).Return(expectedRefs, nil)

	// Act
	result, err := motorcycleInteractor.GetReferencesByBrandID(ctx, brandID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockRepo.AssertExpectations(t)
}

func TestGetReferencesByBrandID_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	brandID := "brand-honda"
	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetReferencesByBrandID", ctx, brandID).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetReferencesByBrandID(ctx, brandID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// DeleteProfileImage Tests (HU39)
// ============================================

func TestDeleteProfileImage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)
	mockTx := new(mocks.MockTx)
	mockStorage := new(mocks.MockStorageClient)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo).
		WithStorageClient(mockStorage)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		LicensePlate:    "ABC123",
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockStorage.On("DeleteStorageFile", ctx, imageURL).Return(nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteProfileImage_Success_NoImage(t *testing.T) {
	// Arrange - motorcycle without profile image
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		LicensePlate:    "ABC123",
		OwnerID:         ownerID,
		ProfileImageURL: nil, // No image
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert - should succeed without doing anything
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_Success_EmptyImageURL(t *testing.T) {
	// Arrange - motorcycle with empty profile image URL
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	emptyURL := ""

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		LicensePlate:    "ABC123",
		OwnerID:         ownerID,
		ProfileImageURL: &emptyURL, // Empty string
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert - should succeed without doing anything
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "non-existent"
	ownerID := "owner-123"

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "other-owner" // Different from actual owner

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      "owner-123", // Actual owner
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_StorageError_ContinuesSuccessfully(t *testing.T) {
	// Arrange - storage deletion fails but operation continues
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)
	mockTx := new(mocks.MockTx)
	mockStorage := new(mocks.MockStorageClient)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo).
		WithStorageClient(mockStorage)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	storageError := errors.New("storage unavailable")

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockStorage.On("DeleteStorageFile", ctx, imageURL).Return(storageError)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert - should succeed even with storage error (best effort)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestDeleteProfileImage_ClearURLError_Rollback(t *testing.T) {
	// Arrange - database clear fails
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)
	mockTx := new(mocks.MockTx)
	mockStorage := new(mocks.MockStorageClient)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo).
		WithStorageClient(mockStorage)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	dbError := errors.New("database error")

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockStorage.On("DeleteStorageFile", ctx, imageURL).Return(nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(dbError)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)

	mockTx.AssertCalled(t, "Rollback")
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_CommitError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)
	mockTx := new(mocks.MockTx)
	mockStorage := new(mocks.MockStorageClient)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo).
		WithStorageClient(mockStorage)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	commitError := errors.New("commit failed")

	// Mock expectations
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockStorage.On("DeleteStorageFile", ctx, imageURL).Return(nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(commitError)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProfileImage_Success_WithoutStorageClient(t *testing.T) {
	// Arrange - no storage client configured (skips storage deletion)
	ctx := context.Background()
	mockRepo := new(mocks.MockMotorcycleRepository)
	mockTx := new(mocks.MockTx)

	// No storage client configured
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockRepo)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	// Mock expectations - no storage deletion expected
	mockRepo.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)
	mockRepo.On("BeginTx", ctx).Return(mockTx, nil)
	mockRepo.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}
