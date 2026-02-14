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
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(nil)
	mockSvc.On("CheckLicensePlateUnique", ctx, "ABC123").Return(nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("RegisterMotorcycle", ctx, mockTx, motorcycle).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestRegisterMotorcycle_ReferenceRequired(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "", // Missing reference
	}

	// Mock expectations
	mockSvc.On("ValidateReferenceExists", ctx, "").Return(domain.ErrReferenceRequired)

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
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "invalid-ref",
	}

	// Mock expectations
	mockSvc.On("ValidateReferenceExists", ctx, "invalid-ref").Return(domain.ErrReferenceNotFound)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)

	mockSvc.AssertExpectations(t)
}

func TestRegisterMotorcycle_DuplicateLicensePlate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	// Mock expectations
	mockSvc.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(nil)
	mockSvc.On("CheckLicensePlateUnique", ctx, "ABC123").Return(domain.ErrDuplicateLicensePlate)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDuplicateLicensePlate, err)

	mockSvc.AssertExpectations(t)
}

func TestRegisterMotorcycle_TxError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	txError := errors.New("transaction error")

	// Mock expectations
	mockSvc.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(nil)
	mockSvc.On("CheckLicensePlateUnique", ctx, "ABC123").Return(nil)
	mockSvc.On("BeginTx", ctx).Return(nil, txError)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)

	mockSvc.AssertExpectations(t)
}

func TestRegisterMotorcycle_SaveError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		ReferenceID:  "ref-honda-cbr",
	}

	// Mock expectations
	mockSvc.On("ValidateReferenceExists", ctx, "ref-honda-cbr").Return(nil)
	mockSvc.On("CheckLicensePlateUnique", ctx, "ABC123").Return(nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("RegisterMotorcycle", ctx, mockTx, motorcycle).Return(domain.ErrMotorcycleCannotSave)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := motorcycleInteractor.RegisterMotorcycle(ctx, motorcycle)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockTx.AssertCalled(t, "Rollback")
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetMotorcycleByID Tests (HU46)
// ============================================

func TestGetMotorcycleByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	expectedMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
	}

	// Mock expectations
	mockSvc.On("GetByID", ctx, motorcycleID).Return(expectedMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByID(ctx, motorcycleID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, motorcycleID, result.ID)

	mockSvc.AssertExpectations(t)
}

func TestGetMotorcycleByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "non-existent"

	// Mock expectations
	mockSvc.On("GetByID", ctx, motorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByID(ctx, motorcycleID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockSvc.AssertExpectations(t)
}

// ============================================
// GetMotorcyclesByOwner Tests (HU47)
// ============================================

func TestGetMotorcyclesByOwner_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	ownerID := "owner-123"
	expectedMotorcycles := []domain.Motorcycle{
		{ID: "moto-1", LicensePlate: "ABC123", OwnerID: ownerID},
		{ID: "moto-2", LicensePlate: "DEF456", OwnerID: ownerID},
	}

	// Mock expectations
	mockSvc.On("GetByOwnerID", ctx, ownerID).Return(expectedMotorcycles, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockSvc.AssertExpectations(t)
}

func TestGetMotorcyclesByOwner_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	ownerID := "new-owner"

	// Mock expectations
	mockSvc.On("GetByOwnerID", ctx, ownerID).Return([]domain.Motorcycle{}, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockSvc.AssertExpectations(t)
}

func TestGetMotorcyclesByOwner_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	ownerID := "owner-123"
	dbError := errors.New("database error")

	// Mock expectations
	mockSvc.On("GetByOwnerID", ctx, ownerID).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetMotorcyclesByOwner(ctx, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockSvc.AssertExpectations(t)
}

// ============================================
// GetMotorcycleByLicensePlate Tests (HU47)
// ============================================

func TestGetMotorcycleByLicensePlate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	licensePlate := "ABC123"
	expectedMotorcycle := &domain.Motorcycle{
		ID:           "moto-123",
		LicensePlate: licensePlate,
		OwnerID:      "owner-123",
	}

	// Mock expectations
	mockSvc.On("GetByLicensePlate", ctx, licensePlate).Return(expectedMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByLicensePlate(ctx, licensePlate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, licensePlate, result.LicensePlate)

	mockSvc.AssertExpectations(t)
}

func TestGetMotorcycleByLicensePlate_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	licensePlate := "INVALID"

	// Mock expectations
	mockSvc.On("GetByLicensePlate", ctx, licensePlate).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleByLicensePlate(ctx, licensePlate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockSvc.AssertExpectations(t)
}

// ============================================
// UpdateMotorcycle Tests (HU44)
// ============================================

func TestUpdateMotorcycle_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("ApplyUpdates", existingMotorcycle, updates).Return(nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("UpdateMotorcycle", ctx, mockTx, existingMotorcycle).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockSvc.On("GetByID", ctx, motorcycleID).Return(existingMotorcycle, nil)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateMotorcycle_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "non-existent"
	ownerID := "owner-123"
	updates := &domain.Motorcycle{}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockSvc.AssertExpectations(t)
}

func TestUpdateMotorcycle_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "other-owner" // Different from actual owner
	updates := &domain.Motorcycle{}

	// Mock expectations - service returns not found for non-owner (security by obscurity)
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockSvc.AssertExpectations(t)
}

func TestUpdateMotorcycle_ReferenceNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("ApplyUpdates", existingMotorcycle, updates).Return(domain.ErrReferenceNotFound)

	// Act
	result, err := motorcycleInteractor.UpdateMotorcycle(ctx, motorcycleID, ownerID, updates)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)

	mockSvc.AssertExpectations(t)
}

// ============================================
// DeleteMotorcycle Tests (HU45)
// ============================================

func TestDeleteMotorcycle_Success_WithHistory(t *testing.T) {
	// Arrange - motorcycle WITH service history -> service handles soft/hard delete internally
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:           motorcycleID,
		LicensePlate: "ABC123",
		OwnerID:      ownerID,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("DeleteMotorcycle", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteMotorcycle_Success_WithProfileImage(t *testing.T) {
	// Arrange - motorcycle with profile image -> storage cleanup
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-456"
	ownerID := "owner-123"
	imageURL := "https://storage.example.com/image.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		LicensePlate:    "DEF456",
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("DeleteStorageFile", ctx, imageURL).Return()
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("DeleteMotorcycle", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteMotorcycle_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "non-existent"
	ownerID := "owner-123"

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockSvc.AssertExpectations(t)
}

func TestDeleteMotorcycle_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "other-owner"

	// Mock expectations - service returns not found for non-owner (security by obscurity)
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockSvc.AssertExpectations(t)
}

func TestDeleteMotorcycle_DeleteError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("DeleteMotorcycle", ctx, mockTx, motorcycleID).Return(domain.ErrMotorcycleCannotDelete)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteMotorcycle(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)

	mockTx.AssertCalled(t, "Rollback")
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetMotorcycleReferences Tests (HU50)
// ============================================

func TestGetMotorcycleReferences_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expectedRefs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "CBR 600"},
		{ID: "ref-2", Model: "Ninja 650"},
	}

	// Mock expectations
	mockSvc.On("GetAllReferences", ctx).Return(expectedRefs, nil)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleReferences(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockSvc.AssertExpectations(t)
}

func TestGetMotorcycleReferences_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	dbError := errors.New("database error")

	// Mock expectations
	mockSvc.On("GetAllReferences", ctx).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetMotorcycleReferences(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockSvc.AssertExpectations(t)
}

// ============================================
// GetReferencesByBrandID Tests (HU40)
// ============================================

func TestGetReferencesByBrandID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	brandID := "brand-honda"
	expectedRefs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "CBR 600", BrandID: brandID},
		{ID: "ref-2", Model: "CBR 1000", BrandID: brandID},
	}

	// Mock expectations
	mockSvc.On("GetReferencesByBrandID", ctx, brandID).Return(expectedRefs, nil)

	// Act
	result, err := motorcycleInteractor.GetReferencesByBrandID(ctx, brandID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockSvc.AssertExpectations(t)
}

func TestGetReferencesByBrandID_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	brandID := "brand-honda"
	dbError := errors.New("database error")

	// Mock expectations
	mockSvc.On("GetReferencesByBrandID", ctx, brandID).Return(nil, dbError)

	// Act
	result, err := motorcycleInteractor.GetReferencesByBrandID(ctx, brandID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockSvc.AssertExpectations(t)
}

// ============================================
// DeleteProfileImage Tests (HU39)
// ============================================

func TestDeleteProfileImage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("DeleteStorageFile", ctx, imageURL).Return()
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestDeleteProfileImage_Success_NoImage(t *testing.T) {
	// Arrange - motorcycle without profile image
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		LicensePlate:    "ABC123",
		OwnerID:         ownerID,
		ProfileImageURL: nil, // No image
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert - should succeed without doing anything
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileImage_Success_EmptyImageURL(t *testing.T) {
	// Arrange - motorcycle with empty profile image URL
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert - should succeed without doing anything
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileImage_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "non-existent"
	ownerID := "owner-123"

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)

	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileImage_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "other-owner"

	// Mock expectations - service returns not found for non-owner (security by obscurity)
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Returns 404 for security

	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileImage_ClearURLError_Rollback(t *testing.T) {
	// Arrange - database clear fails
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "owner-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/motogo.appspot.com/o/motorcycles%2Fprofile.jpg"

	existingMotorcycle := &domain.Motorcycle{
		ID:              motorcycleID,
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("DeleteStorageFile", ctx, imageURL).Return()
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(domain.ErrMotorcycleCannotUpdate)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)

	mockTx.AssertCalled(t, "Rollback")
	mockSvc.AssertExpectations(t)
}

func TestDeleteProfileImage_CommitError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

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
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("DeleteStorageFile", ctx, imageURL).Return()
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("ClearProfileImageURL", ctx, mockTx, motorcycleID).Return(nil)
	mockTx.On("Commit").Return(commitError)
	mockTx.On("Rollback").Return(nil) // defer rollback fires on commit failure

	// Act
	err := motorcycleInteractor.DeleteProfileImage(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GrantDiagnosticPermission Tests
// ============================================

func TestGrantDiagnosticPermission_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	expectedPermission := &domain.DiagnosticPermission{
		MotorcycleID: motorcycleID,
		BranchID:     branchID,
		Active:       true,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("GrantPermission", ctx, mockTx, motorcycleID, branchID, true).Return(expectedPermission, nil)
	mockTx.On("Commit").Return(nil)

	// Act
	result, err := motorcycleInteractor.GrantDiagnosticPermission(ctx, motorcycleID, branchID, ownerID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, motorcycleID, result.MotorcycleID)
	assert.Equal(t, branchID, result.BranchID)
	assert.True(t, result.Active)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestGrantDiagnosticPermission_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "other-owner"

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.GrantDiagnosticPermission(ctx, motorcycleID, branchID, ownerID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockSvc.AssertExpectations(t)
}

func TestGrantDiagnosticPermission_SaveError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("GrantPermission", ctx, mockTx, motorcycleID, branchID, true).Return(nil, domain.ErrPermissionCannotSave)
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := motorcycleInteractor.GrantDiagnosticPermission(ctx, motorcycleID, branchID, ownerID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockTx.AssertCalled(t, "Rollback")
	mockSvc.AssertExpectations(t)
}

// ============================================
// RevokeDiagnosticPermission Tests
// ============================================

func TestRevokeDiagnosticPermission_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("RevokePermission", ctx, mockTx, motorcycleID, branchID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := motorcycleInteractor.RevokeDiagnosticPermission(ctx, motorcycleID, branchID, ownerID)

	// Assert
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
	mockTx.AssertCalled(t, "Commit")
}

func TestRevokeDiagnosticPermission_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "other-owner"

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	err := motorcycleInteractor.RevokeDiagnosticPermission(ctx, motorcycleID, branchID, ownerID)

	// Assert
	assert.Error(t, err)

	mockSvc.AssertExpectations(t)
}

func TestRevokeDiagnosticPermission_DeleteError_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	branchID := "branch-456"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("BeginTx", ctx).Return(mockTx, nil)
	mockSvc.On("RevokePermission", ctx, mockTx, motorcycleID, branchID).Return(domain.ErrPermissionCannotDelete)
	mockTx.On("Rollback").Return(nil)

	// Act
	err := motorcycleInteractor.RevokeDiagnosticPermission(ctx, motorcycleID, branchID, ownerID)

	// Assert
	assert.Error(t, err)

	mockTx.AssertCalled(t, "Rollback")
	mockSvc.AssertExpectations(t)
}

// ============================================
// ListDiagnosticPermissions Tests
// ============================================

func TestListDiagnosticPermissions_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "owner-123"

	existingMotorcycle := &domain.Motorcycle{
		ID:      motorcycleID,
		OwnerID: ownerID,
	}

	expectedPermissions := []domain.DiagnosticPermission{
		{MotorcycleID: motorcycleID, BranchID: "branch-1", Active: true},
		{MotorcycleID: motorcycleID, BranchID: "branch-2", Active: true},
	}

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(existingMotorcycle, nil)
	mockSvc.On("ListPermissions", ctx, motorcycleID).Return(expectedPermissions, nil)

	// Act
	result, err := motorcycleInteractor.ListDiagnosticPermissions(ctx, motorcycleID, ownerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockSvc.AssertExpectations(t)
}

func TestListDiagnosticPermissions_NotOwner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	motorcycleID := "moto-123"
	ownerID := "other-owner"

	// Mock expectations
	mockSvc.On("ValidateOwnership", ctx, motorcycleID, ownerID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := motorcycleInteractor.ListDiagnosticPermissions(ctx, motorcycleID, ownerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockSvc.AssertExpectations(t)
}

// Ignore unused import for mock
var _ = mock.Anything

// ============================================
// GetDistinctCategories Interactor Tests (HU41)
// ============================================

func TestGetDistinctCategories_Interactor_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expected := []domain.MotorcycleCategory{{Name: "Sport", LineCount: 5}}
	mockSvc.On("GetDistinctCategories", ctx).Return(expected, nil)

	result, err := motorcycleInteractor.GetDistinctCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Sport", result[0].Name)
	mockSvc.AssertExpectations(t)
}

func TestGetDistinctCategories_Interactor_Error(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	mockSvc.On("GetDistinctCategories", ctx).Return(nil, errors.New("database error"))

	result, err := motorcycleInteractor.GetDistinctCategories(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetLinesByCategory Interactor Tests (HU41)
// ============================================

func TestGetLinesByCategory_Interactor_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expected := []domain.CategoryLine{{Model: "Ninja"}}
	mockSvc.On("GetLinesByCategory", ctx, "Sport").Return(expected, nil)

	result, err := motorcycleInteractor.GetLinesByCategory(ctx, "Sport")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Ninja", result[0].Model)
	mockSvc.AssertExpectations(t)
}

func TestGetLinesByCategory_Interactor_Error(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	mockSvc.On("GetLinesByCategory", ctx, "Sport").Return(nil, errors.New("database error"))

	result, err := motorcycleInteractor.GetLinesByCategory(ctx, "Sport")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetDistinctDisplacements Interactor Tests (HU49)
// ============================================

func TestGetDistinctDisplacements_Interactor_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expected := []domain.EngineDisplacementRange{
		{Range: domain.DisplacementRangeLow},
		{Range: domain.DisplacementRangeMedium},
		{Range: domain.DisplacementRangeHigh},
	}
	mockSvc.On("GetDistinctDisplacements", ctx).Return(expected, nil)

	result, err := motorcycleInteractor.GetDistinctDisplacements(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	mockSvc.AssertExpectations(t)
}

func TestGetDistinctDisplacements_Interactor_Error(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	mockSvc.On("GetDistinctDisplacements", ctx).Return(nil, errors.New("database error"))

	result, err := motorcycleInteractor.GetDistinctDisplacements(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockSvc.AssertExpectations(t)
}

// ============================================
// GetRatingRanges Interactor Tests (HU48)
// ============================================

func TestGetRatingRanges_Interactor_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expected := []domain.RatingRange{
		{Value: 1, Label: "Very bad"},
		{Value: 5, Label: "Excellent"},
	}
	mockSvc.On("GetRatingRanges", ctx).Return(expected, nil)

	result, err := motorcycleInteractor.GetRatingRanges(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockSvc.AssertExpectations(t)
}

func TestGetRatingRanges_Interactor_Error(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	mockSvc.On("GetRatingRanges", ctx).Return(nil, errors.New("database error"))

	result, err := motorcycleInteractor.GetRatingRanges(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockSvc.AssertExpectations(t)
}

// ============================================
// LookupPermissions Interactor Tests
// ============================================

func TestLookupPermissions_Interactor_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	expected := []domain.DiagnosticPermission{
		{MotorcycleID: "m-1", BranchID: "b-1", Active: true},
		{MotorcycleID: "m-1", BranchID: "b-2", Active: true},
	}
	mockSvc.On("ListPermissions", ctx, "m-1").Return(expected, nil)

	result, err := motorcycleInteractor.LookupPermissions(ctx, "m-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockSvc.AssertExpectations(t)
}

func TestLookupPermissions_Interactor_Error(t *testing.T) {
	ctx := context.Background()
	mockSvc := new(mocks.MockMotorcycleService)
	motorcycleInteractor := interactor.NewMotorcycleInteractor(mockSvc)

	mockSvc.On("ListPermissions", ctx, "m-1").Return(nil, errors.New("database error"))

	result, err := motorcycleInteractor.LookupPermissions(ctx, "m-1")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockSvc.AssertExpectations(t)
}
