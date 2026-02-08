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
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	created := &domain.Motorcycle{
		ID:           "new-uuid",
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(nil)
	mockService.On("ValidateLicensePlateUnique", mock.Anything, "ABC123").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("CreateMotorcycle", mock.Anything, mockTx, motorcycle).Return(created, nil)
	mockTx.On("Commit").Return(nil)

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-uuid", result.ID)
	mockService.AssertExpectations(t)
}

func TestRegisterMotorcycle_ReferenceRequired(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "").Return(domain.ErrReferenceRequired)

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceRequired, err)
}

func TestRegisterMotorcycle_ReferenceNotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "invalid-ref",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "invalid-ref").Return(domain.ErrReferenceNotFound)

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)
}

func TestRegisterMotorcycle_DuplicateLicensePlate(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "DUP123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(nil)
	mockService.On("ValidateLicensePlateUnique", mock.Anything, "DUP123").Return(domain.ErrDuplicateLicensePlate)

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDuplicateLicensePlate, err)
}

func TestRegisterMotorcycle_TxError(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(nil)
	mockService.On("ValidateLicensePlateUnique", mock.Anything, "ABC123").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

func TestRegisterMotorcycle_SaveError_Rollback(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	mockService.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(nil)
	mockService.On("ValidateLicensePlateUnique", mock.Anything, "ABC123").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("CreateMotorcycle", mock.Anything, mockTx, motorcycle).Return(nil, domain.ErrMotorcycleCannotSave)
	mockTx.On("Rollback").Return(nil)

	result, err := i.RegisterMotorcycle(context.Background(), motorcycle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

// ============================================
// GetMotorcycleByID Tests (HU46)
// ============================================

func TestGetMotorcycleByID_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:           "moto-1",
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
	}

	mockService.On("GetMotorcycleByID", mock.Anything, "moto-1").Return(motorcycle, nil)

	result, err := i.GetMotorcycleByID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.ID)
}

func TestGetMotorcycleByID_NotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetMotorcycleByID", mock.Anything, "moto-999").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.GetMotorcycleByID(context.Background(), "moto-999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// GetMotorcyclesByOwner Tests (HU47)
// ============================================

func TestGetMotorcyclesByOwner_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycles := []domain.Motorcycle{
		{ID: "moto-1", OwnerID: "owner-1"},
		{ID: "moto-2", OwnerID: "owner-1"},
	}

	mockService.On("GetMotorcyclesByOwner", mock.Anything, "owner-1").Return(motorcycles, nil)

	result, err := i.GetMotorcyclesByOwner(context.Background(), "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetMotorcyclesByOwner_Empty(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetMotorcyclesByOwner", mock.Anything, "owner-1").Return([]domain.Motorcycle{}, nil)

	result, err := i.GetMotorcyclesByOwner(context.Background(), "owner-1")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetMotorcyclesByOwner_Error(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetMotorcyclesByOwner", mock.Anything, "owner-1").Return(nil, errors.New("db error"))

	result, err := i.GetMotorcyclesByOwner(context.Background(), "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetMotorcycleByLicensePlate Tests (HU47)
// ============================================

func TestGetMotorcycleByLicensePlate_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:           "moto-1",
		LicensePlate: "ABC123",
	}

	mockService.On("GetMotorcycleByLicensePlate", mock.Anything, "ABC123").Return(motorcycle, nil)

	result, err := i.GetMotorcycleByLicensePlate(context.Background(), "ABC123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ABC123", result.LicensePlate)
}

func TestGetMotorcycleByLicensePlate_NotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetMotorcycleByLicensePlate", mock.Anything, "ZZZ999").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.GetMotorcycleByLicensePlate(context.Background(), "ZZZ999")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// UpdateMotorcycle Tests (HU44)
// ============================================

func TestUpdateMotorcycle_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:           "moto-1",
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
	}

	year := 2023
	updates := &domain.Motorcycle{Year: &year}

	updated := &domain.Motorcycle{
		ID:           "moto-1",
		LicensePlate: "ABC123",
		OwnerID:      "owner-1",
		ReferenceID:  "ref-1",
		Year:         &year,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("ApplyMotorcycleUpdates", mock.Anything, motorcycle, updates).Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateMotorcycle", mock.Anything, mockTx, motorcycle).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockService.On("GetMotorcycleByID", mock.Anything, "moto-1").Return(updated, nil)

	result, err := i.UpdateMotorcycle(context.Background(), "moto-1", "owner-1", updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2023, *result.Year)
	mockService.AssertExpectations(t)
}

func TestUpdateMotorcycle_NotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	updates := &domain.Motorcycle{}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-999", "owner-1").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.UpdateMotorcycle(context.Background(), "moto-999", "owner-1", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestUpdateMotorcycle_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	updates := &domain.Motorcycle{}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.UpdateMotorcycle(context.Background(), "moto-1", "wrong-owner", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestUpdateMotorcycle_ReferenceNotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:          "moto-1",
		OwnerID:     "owner-1",
		ReferenceID: "ref-1",
	}

	updates := &domain.Motorcycle{ReferenceID: "invalid-ref"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("ApplyMotorcycleUpdates", mock.Anything, motorcycle, updates).Return(domain.ErrReferenceNotFound)

	result, err := i.UpdateMotorcycle(context.Background(), "moto-1", "owner-1", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrReferenceNotFound, err)
}

// ============================================
// DeleteMotorcycle Tests (HU45)
// ============================================

func TestDeleteMotorcycle_Success_SoftDelete(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:      "moto-1",
		OwnerID: "owner-1",
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("CheckServiceHistory", mock.Anything, "moto-1").Return(true, nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteMotorcycle", mock.Anything, mockTx, "moto-1", true).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.DeleteMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeleteMotorcycle_Success_HardDelete(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:      "moto-1",
		OwnerID: "owner-1",
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("CheckServiceHistory", mock.Anything, "moto-1").Return(false, nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteMotorcycle", mock.Anything, mockTx, "moto-1", false).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.DeleteMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
}

func TestDeleteMotorcycle_NotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-999", "owner-1").Return(nil, domain.ErrMotorcycleNotFound)

	err := i.DeleteMotorcycle(context.Background(), "moto-999", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDeleteMotorcycle_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	err := i.DeleteMotorcycle(context.Background(), "moto-1", "wrong-owner")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDeleteMotorcycle_TxError_Rollback(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:      "moto-1",
		OwnerID: "owner-1",
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("CheckServiceHistory", mock.Anything, "moto-1").Return(false, nil)
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	err := i.DeleteMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotDelete, err)
}

func TestDeleteMotorcycle_WithProfileImage(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/o/image.jpg"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: &imageURL,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("CheckServiceHistory", mock.Anything, "moto-1").Return(false, nil)
	mockService.On("DeleteStorageFile", mock.Anything, imageURL).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteMotorcycle", mock.Anything, mockTx, "moto-1", false).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.DeleteMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// ============================================
// GetMotorcycleReferences Tests (HU50)
// ============================================

func TestGetMotorcycleReferences_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	refs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "Reference 1"},
		{ID: "ref-2", Model: "Reference 2"},
	}

	mockService.On("GetAllReferences", mock.Anything).Return(refs, nil)

	result, err := i.GetMotorcycleReferences(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetMotorcycleReferences_Error(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetAllReferences", mock.Anything).Return(nil, errors.New("db error"))

	result, err := i.GetMotorcycleReferences(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetReferencesByBrandID Tests (HU40)
// ============================================

func TestGetReferencesByBrandID_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	refs := []domain.MotorcycleReference{
		{ID: "ref-1", Model: "Reference 1"},
	}

	mockService.On("GetReferencesByBrandID", mock.Anything, "brand-1").Return(refs, nil)

	result, err := i.GetReferencesByBrandID(context.Background(), "brand-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetReferencesByBrandID_Error(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("GetReferencesByBrandID", mock.Anything, "brand-1").Return(nil, errors.New("db error"))

	result, err := i.GetReferencesByBrandID(context.Background(), "brand-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// DeleteProfileImage Tests (HU39)
// ============================================

func TestDeleteProfileImage_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/o/image.jpg"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: &imageURL,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("DeleteStorageFile", mock.Anything, imageURL).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteProfileImage", mock.Anything, mockTx, "moto-1").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeleteProfileImage_Success_NoImage(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: nil,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
}

func TestDeleteProfileImage_Success_EmptyImageURL(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	emptyURL := ""
	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: &emptyURL,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
}

func TestDeleteProfileImage_NotFound(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-999", "owner-1").Return(nil, domain.ErrMotorcycleNotFound)

	err := i.DeleteProfileImage(context.Background(), "moto-999", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDeleteProfileImage_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "wrong-owner")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDeleteProfileImage_ClearURLError_Rollback(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/o/image.jpg"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: &imageURL,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("DeleteStorageFile", mock.Anything, imageURL).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteProfileImage", mock.Anything, mockTx, "moto-1").Return(errors.New("clear error"))
	mockTx.On("Rollback").Return(nil)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)
}

func TestDeleteProfileImage_CommitError(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/o/image.jpg"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-1",
		OwnerID:         "owner-1",
		ProfileImageURL: &imageURL,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("DeleteStorageFile", mock.Anything, imageURL).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteProfileImage", mock.Anything, mockTx, "moto-1").Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := i.DeleteProfileImage(context.Background(), "moto-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)
}

// ============================================
// GrantDiagnosticPermission Tests
// ============================================

func TestGrantDiagnosticPermission_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	permission := &domain.DiagnosticPermission{
		ID:           "perm-1",
		MotorcycleID: "moto-1",
		BranchID:     "branch-1",
		Active:       true,
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("BeginPermissionTx", mock.Anything).Return(mockTx, nil)
	mockService.On("GrantDiagnosticPermission", mock.Anything, mockTx, "moto-1", "branch-1").Return(permission, nil)
	mockTx.On("Commit").Return(nil)

	result, err := i.GrantDiagnosticPermission(context.Background(), "moto-1", "branch-1", "owner-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	mockService.AssertExpectations(t)
}

func TestGrantDiagnosticPermission_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.GrantDiagnosticPermission(context.Background(), "moto-1", "branch-1", "wrong-owner")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGrantDiagnosticPermission_TxError(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("BeginPermissionTx", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := i.GrantDiagnosticPermission(context.Background(), "moto-1", "branch-1", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrPermissionCannotSave, err)
}

func TestGrantDiagnosticPermission_SaveError(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("BeginPermissionTx", mock.Anything).Return(mockTx, nil)
	mockService.On("GrantDiagnosticPermission", mock.Anything, mockTx, "moto-1", "branch-1").Return(nil, domain.ErrPermissionCannotSave)
	mockTx.On("Rollback").Return(nil)

	result, err := i.GrantDiagnosticPermission(context.Background(), "moto-1", "branch-1", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// RevokeDiagnosticPermission Tests
// ============================================

func TestRevokeDiagnosticPermission_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("BeginPermissionTx", mock.Anything).Return(mockTx, nil)
	mockService.On("RevokeDiagnosticPermission", mock.Anything, mockTx, "moto-1", "branch-1").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.RevokeDiagnosticPermission(context.Background(), "moto-1", "branch-1", "owner-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestRevokeDiagnosticPermission_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	err := i.RevokeDiagnosticPermission(context.Background(), "moto-1", "branch-1", "wrong-owner")

	assert.Error(t, err)
}

func TestRevokeDiagnosticPermission_DeleteError(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("BeginPermissionTx", mock.Anything).Return(mockTx, nil)
	mockService.On("RevokeDiagnosticPermission", mock.Anything, mockTx, "moto-1", "branch-1").Return(domain.ErrPermissionNotFound)
	mockTx.On("Rollback").Return(nil)

	err := i.RevokeDiagnosticPermission(context.Background(), "moto-1", "branch-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPermissionNotFound, err)
}

// ============================================
// ListDiagnosticPermissions Tests
// ============================================

func TestListDiagnosticPermissions_Success(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	permissions := []domain.DiagnosticPermission{
		{ID: "perm-1", MotorcycleID: "moto-1", BranchID: "branch-1"},
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("ListDiagnosticPermissions", mock.Anything, "moto-1").Return(permissions, nil)

	result, err := i.ListDiagnosticPermissions(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestListDiagnosticPermissions_NotOwner(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := i.ListDiagnosticPermissions(context.Background(), "moto-1", "wrong-owner")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListDiagnosticPermissions_Error(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	i := interactor.NewMotorcycleInteractor(mockService)

	motorcycle := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(motorcycle, nil)
	mockService.On("ListDiagnosticPermissions", mock.Anything, "moto-1").Return(nil, errors.New("db error"))

	result, err := i.ListDiagnosticPermissions(context.Background(), "moto-1", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// WithStorageClient Tests
// ============================================

func TestMotorcycleWithStorageClient(t *testing.T) {
	mockService := new(mocks.MockMotorcycleService)
	mockStorage := new(mocks.MockStorageClient)
	i := interactor.NewMotorcycleInteractor(mockService)

	mockService.On("WithStorageClient", mockStorage).Return()

	result := i.WithStorageClient(mockStorage)

	assert.NotNil(t, result)
	assert.Equal(t, i, result)
	mockService.AssertExpectations(t)
}
