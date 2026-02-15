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
// Helper
// ============================================

func setupMotorcycleMocks() (*mocks.MockMotorcycleRepository, *mocks.MockDiagnosticPermissionRepository) {
	return &mocks.MockMotorcycleRepository{}, &mocks.MockDiagnosticPermissionRepository{}
}

// ============================================
// NewMotorcycleService Tests
// ============================================

func TestNewMotorcycleService_ReturnsNonNil(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	assert.NotNil(t, svc)
}

// ============================================
// WithStorageClient Tests
// ============================================

func TestWithStorageClient_SetsClient(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	storageMock := &mocks.MockStorageClient{}
	svc.WithStorageClient(storageMock)
	// No error means success - field is set
}

// ============================================
// BeginTx Tests
// ============================================

func TestMotorcycleService_BeginTx_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)

	tx, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	repoMock.AssertExpectations(t)
}

func TestMotorcycleService_BeginTx_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("BeginTx", mock.Anything).Return(nil, errors.New("db error"))

	tx, err := svc.BeginTx(context.Background())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

// ============================================
// ValidateReferenceExists Tests
// ============================================

func TestValidateReferenceExists_EmptyID(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	err := svc.ValidateReferenceExists(context.Background(), "")
	assert.ErrorIs(t, err, domain.ErrReferenceRequired)
}

func TestValidateReferenceExists_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-123").Return(true, nil)

	err := svc.ValidateReferenceExists(context.Background(), "ref-123")
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestValidateReferenceExists_NotFound(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-999").Return(false, nil)

	err := svc.ValidateReferenceExists(context.Background(), "ref-999")
	assert.ErrorIs(t, err, domain.ErrReferenceNotFound)
}

func TestValidateReferenceExists_RepoError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-123").Return(false, errors.New("db err"))

	err := svc.ValidateReferenceExists(context.Background(), "ref-123")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotSave)
}

// ============================================
// CheckLicensePlateUnique Tests
// ============================================

func TestCheckLicensePlateUnique_Unique(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(false, nil)

	err := svc.CheckLicensePlateUnique(context.Background(), "ABC123")
	assert.NoError(t, err)
}

func TestCheckLicensePlateUnique_Duplicate(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(true, nil)

	err := svc.CheckLicensePlateUnique(context.Background(), "ABC123")
	assert.ErrorIs(t, err, domain.ErrDuplicateLicensePlate)
}

func TestCheckLicensePlateUnique_RepoError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(false, errors.New("db err"))

	err := svc.CheckLicensePlateUnique(context.Background(), "ABC123")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotSave)
}

// ============================================
// ValidateOwnership Tests
// ============================================

func TestValidateOwnership_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	moto := &domain.Motorcycle{ID: "m-1", OwnerID: "owner-1"}

	repoMock.On("GetByID", mock.Anything, "m-1").Return(moto, nil)

	result, err := svc.ValidateOwnership(context.Background(), "m-1", "owner-1")
	assert.NoError(t, err)
	assert.Equal(t, "m-1", result.ID)
}

func TestValidateOwnership_NotOwner(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	moto := &domain.Motorcycle{ID: "m-1", OwnerID: "owner-1"}

	repoMock.On("GetByID", mock.Anything, "m-1").Return(moto, nil)

	result, err := svc.ValidateOwnership(context.Background(), "m-1", "other-owner")
	assert.ErrorIs(t, err, domain.ErrMotorcycleNotFound)
	assert.Nil(t, result)
}

func TestValidateOwnership_RepoError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetByID", mock.Anything, "m-1").Return(nil, errors.New("db err"))

	result, err := svc.ValidateOwnership(context.Background(), "m-1", "owner-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// RegisterMotorcycle Tests
// ============================================

func TestRegisterMotorcycle_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}
	moto := &domain.Motorcycle{LicensePlate: "ABC123"}

	repoMock.On("Save", mock.Anything, txMock, moto).Return(nil)

	err := svc.RegisterMotorcycle(context.Background(), txMock, moto)
	assert.NoError(t, err)
	assert.NotEmpty(t, moto.ID) // ID should have been set
	repoMock.AssertExpectations(t)
}

func TestRegisterMotorcycle_SaveError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}
	moto := &domain.Motorcycle{LicensePlate: "ABC123"}

	repoMock.On("Save", mock.Anything, txMock, moto).Return(errors.New("db err"))

	err := svc.RegisterMotorcycle(context.Background(), txMock, moto)
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotSave)
}

// ============================================
// GetByID Tests
// ============================================

func TestGetByID_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	moto := &domain.Motorcycle{ID: "m-1"}

	repoMock.On("GetByID", mock.Anything, "m-1").Return(moto, nil)

	result, err := svc.GetByID(context.Background(), "m-1")
	assert.NoError(t, err)
	assert.Equal(t, "m-1", result.ID)
}

func TestGetByID_NotFound(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetByID", mock.Anything, "m-999").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := svc.GetByID(context.Background(), "m-999")
	assert.ErrorIs(t, err, domain.ErrMotorcycleNotFound)
	assert.Nil(t, result)
}

// ============================================
// GetByOwnerID Tests
// ============================================

func TestGetByOwnerID_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	list := []domain.Motorcycle{{ID: "m-1"}, {ID: "m-2"}}

	repoMock.On("GetByOwnerID", mock.Anything, "owner-1").Return(list, nil)

	result, err := svc.GetByOwnerID(context.Background(), "owner-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetByOwnerID_Empty(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetByOwnerID", mock.Anything, "owner-1").Return([]domain.Motorcycle{}, nil)

	result, err := svc.GetByOwnerID(context.Background(), "owner-1")
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// GetByLicensePlate Tests
// ============================================

func TestGetByLicensePlate_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	moto := &domain.Motorcycle{ID: "m-1", LicensePlate: "ABC123"}

	repoMock.On("GetByLicensePlate", mock.Anything, "ABC123").Return(moto, nil)

	result, err := svc.GetByLicensePlate(context.Background(), "ABC123")
	assert.NoError(t, err)
	assert.Equal(t, "ABC123", result.LicensePlate)
}

func TestGetByLicensePlate_NotFound(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetByLicensePlate", mock.Anything, "ZZZ999").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := svc.GetByLicensePlate(context.Background(), "ZZZ999")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// ApplyUpdates Tests
// ============================================

func TestApplyUpdates_NoReferenceChange(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	year := 2022
	existing := &domain.Motorcycle{ID: "m-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{Year: &year}

	err := svc.ApplyUpdates(existing, updates)
	assert.NoError(t, err)
	assert.Equal(t, &year, existing.Year)
	assert.Equal(t, "ref-1", existing.ReferenceID) // unchanged
}

func TestApplyUpdates_ReferenceChange_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	existing := &domain.Motorcycle{ID: "m-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-2"}

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-2").Return(true, nil)

	err := svc.ApplyUpdates(existing, updates)
	assert.NoError(t, err)
	assert.Equal(t, "ref-2", existing.ReferenceID) // updated
}

func TestApplyUpdates_ReferenceChange_NotFound(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	existing := &domain.Motorcycle{ID: "m-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-999"}

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-999").Return(false, nil)

	err := svc.ApplyUpdates(existing, updates)
	assert.ErrorIs(t, err, domain.ErrReferenceNotFound)
}

func TestApplyUpdates_ReferenceChange_RepoError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	existing := &domain.Motorcycle{ID: "m-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-2"}

	repoMock.On("ValidateReferenceExists", mock.Anything, "ref-2").Return(false, errors.New("db err"))

	err := svc.ApplyUpdates(existing, updates)
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotUpdate)
}

func TestApplyUpdates_AllOptionalFields(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	year := 2023
	mileage := 15000
	notes := "Well maintained"
	imageURL := "https://example.com/img.jpg"

	existing := &domain.Motorcycle{ID: "m-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{
		Year:            &year,
		CurrentMileage:  &mileage,
		OwnerNotes:      &notes,
		ProfileImageURL: &imageURL,
	}

	err := svc.ApplyUpdates(existing, updates)
	assert.NoError(t, err)
	assert.Equal(t, &year, existing.Year)
	assert.Equal(t, &mileage, existing.CurrentMileage)
	assert.Equal(t, &notes, existing.OwnerNotes)
	assert.Equal(t, &imageURL, existing.ProfileImageURL)
}

// ============================================
// UpdateMotorcycle Tests
// ============================================

func TestUpdateMotorcycle_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}
	moto := &domain.Motorcycle{ID: "m-1"}

	repoMock.On("Update", mock.Anything, txMock, moto).Return(nil)

	err := svc.UpdateMotorcycle(context.Background(), txMock, moto)
	assert.NoError(t, err)
}

func TestUpdateMotorcycle_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}
	moto := &domain.Motorcycle{ID: "m-1"}

	repoMock.On("Update", mock.Anything, txMock, moto).Return(errors.New("db err"))

	err := svc.UpdateMotorcycle(context.Background(), txMock, moto)
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotUpdate)
}

// ============================================
// DeleteMotorcycle Tests (Hybrid Strategy)
// ============================================

func TestDeleteMotorcycle_SoftDelete_WithHistory(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(true, nil)
	repoMock.On("Delete", mock.Anything, txMock, "m-1").Return(nil)

	err := svc.DeleteMotorcycle(context.Background(), txMock, "m-1")
	assert.NoError(t, err)
	repoMock.AssertCalled(t, "Delete", mock.Anything, txMock, "m-1")
	repoMock.AssertNotCalled(t, "HardDelete")
}

func TestDeleteMotorcycle_HardDelete_NoHistory(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(false, nil)
	repoMock.On("HardDelete", mock.Anything, txMock, "m-1").Return(nil)

	err := svc.DeleteMotorcycle(context.Background(), txMock, "m-1")
	assert.NoError(t, err)
	repoMock.AssertCalled(t, "HardDelete", mock.Anything, txMock, "m-1")
	repoMock.AssertNotCalled(t, "Delete")
}

func TestDeleteMotorcycle_HistoryCheckError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(false, errors.New("db err"))

	err := svc.DeleteMotorcycle(context.Background(), txMock, "m-1")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotDelete)
}

func TestDeleteMotorcycle_SoftDeleteError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(true, nil)
	repoMock.On("Delete", mock.Anything, txMock, "m-1").Return(errors.New("db err"))

	err := svc.DeleteMotorcycle(context.Background(), txMock, "m-1")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotDelete)
}

func TestDeleteMotorcycle_HardDeleteError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(false, nil)
	repoMock.On("HardDelete", mock.Anything, txMock, "m-1").Return(errors.New("db err"))

	err := svc.DeleteMotorcycle(context.Background(), txMock, "m-1")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotDelete)
}

// ============================================
// ClearProfileImageURL Tests
// ============================================

func TestClearProfileImageURL_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("ClearProfileImageURL", mock.Anything, txMock, "m-1").Return(nil)

	err := svc.ClearProfileImageURL(context.Background(), txMock, "m-1")
	assert.NoError(t, err)
}

func TestClearProfileImageURL_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	repoMock.On("ClearProfileImageURL", mock.Anything, txMock, "m-1").Return(errors.New("db err"))

	err := svc.ClearProfileImageURL(context.Background(), txMock, "m-1")
	assert.ErrorIs(t, err, domain.ErrMotorcycleCannotUpdate)
}

// ============================================
// HasServiceHistory Tests
// ============================================

func TestHasServiceHistory_True(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(true, nil)

	has, err := svc.HasServiceHistory(context.Background(), "m-1")
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestHasServiceHistory_False(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("HasServiceHistory", mock.Anything, "m-1").Return(false, nil)

	has, err := svc.HasServiceHistory(context.Background(), "m-1")
	assert.NoError(t, err)
	assert.False(t, has)
}

// ============================================
// DeleteStorageFile Tests
// ============================================

func TestDeleteStorageFile_EmptyURL(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	storageMock := &mocks.MockStorageClient{}
	svc.WithStorageClient(storageMock)

	// Should not call storage client with empty URL
	svc.DeleteStorageFile(context.Background(), "")
	storageMock.AssertNotCalled(t, "DeleteStorageFile")
}

func TestDeleteStorageFile_NilClient(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	// No storage client set - should not panic
	svc.DeleteStorageFile(context.Background(), "https://storage.example.com/file.jpg")
}

func TestDeleteStorageFile_Success(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	storageMock := &mocks.MockStorageClient{}
	svc.WithStorageClient(storageMock)

	storageMock.On("DeleteStorageFile", mock.Anything, "https://storage.example.com/file.jpg").Return(nil)

	svc.DeleteStorageFile(context.Background(), "https://storage.example.com/file.jpg")
	storageMock.AssertExpectations(t)
}

func TestDeleteStorageFile_ErrorIsSwallowed(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	storageMock := &mocks.MockStorageClient{}
	svc.WithStorageClient(storageMock)

	storageMock.On("DeleteStorageFile", mock.Anything, "https://storage.example.com/file.jpg").Return(errors.New("storage err"))

	// Should not panic - errors are logged but not returned
	svc.DeleteStorageFile(context.Background(), "https://storage.example.com/file.jpg")
	storageMock.AssertExpectations(t)
}

// ============================================
// GetAllReferences Tests
// ============================================

func TestGetAllReferences_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	refs := []domain.MotorcycleReference{{ID: "r-1"}}

	repoMock.On("GetAllReferences", mock.Anything).Return(refs, nil)

	result, err := svc.GetAllReferences(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// ============================================
// GetReferencesByBrandID Tests
// ============================================

func TestGetReferencesByBrandID_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	refs := []domain.MotorcycleReference{{ID: "r-1"}}

	repoMock.On("GetReferencesByBrandID", mock.Anything, "brand-1").Return(refs, nil)

	result, err := svc.GetReferencesByBrandID(context.Background(), "brand-1")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// ============================================
// GetDistinctCategories Tests
// ============================================

func TestGetDistinctCategories_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	cats := []domain.MotorcycleCategory{{Name: "Sport"}}

	repoMock.On("GetDistinctCategories", mock.Anything).Return(cats, nil)

	result, err := svc.GetDistinctCategories(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Sport", result[0].Name)
}

func TestGetDistinctCategories_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetDistinctCategories", mock.Anything).Return(nil, errors.New("db err"))

	result, err := svc.GetDistinctCategories(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetLinesByCategory Tests
// ============================================

func TestGetLinesByCategory_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	lines := []domain.CategoryLine{{Model: "Ninja"}}

	repoMock.On("GetLinesByCategory", mock.Anything, "Sport").Return(lines, nil)

	result, err := svc.GetLinesByCategory(context.Background(), "Sport")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetLinesByCategory_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	repoMock.On("GetLinesByCategory", mock.Anything, "Sport").Return(nil, errors.New("db err"))

	result, err := svc.GetLinesByCategory(context.Background(), "Sport")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetDistinctDisplacements Tests
// ============================================

func TestGetDistinctDisplacements_ReturnsThreeRanges(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, err := svc.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Len(t, displacements, 3)
}

func TestGetDistinctDisplacements_CorrectValues(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, err := svc.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, domain.DisplacementRangeLow, displacements[0].Range)
	assert.Equal(t, domain.DisplacementRangeMedium, displacements[1].Range)
	assert.Equal(t, domain.DisplacementRangeHigh, displacements[2].Range)
}

func TestGetDistinctDisplacements_CorrectOrder(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	displacements, _ := svc.GetDistinctDisplacements(context.Background())

	// Order should be BAJO, MEDIO, ALTO (ascending)
	assert.Equal(t, domain.DisplacementRange("BAJO"), displacements[0].Range)
	assert.Equal(t, domain.DisplacementRange("MEDIO"), displacements[1].Range)
	assert.Equal(t, domain.DisplacementRange("ALTO"), displacements[2].Range)
}

// ============================================
// GetRatingRanges Tests
// ============================================

func TestGetRatingRanges_ReturnsFiveRanges(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, err := svc.GetRatingRanges(context.Background())

	assert.NoError(t, err)
	assert.Len(t, ratings, 5)
}

func TestGetRatingRanges_CorrectValues(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, err := svc.GetRatingRanges(context.Background())

	assert.NoError(t, err)

	expected := []struct {
		value int
		label string
	}{
		{1, "Very bad"},
		{2, "Bad"},
		{3, "Average"},
		{4, "Good"},
		{5, "Excellent"},
	}

	for i, exp := range expected {
		assert.Equal(t, exp.value, ratings[i].Value)
		assert.Equal(t, exp.label, ratings[i].Label)
	}
}

func TestGetRatingRanges_FirstIsOne(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, _ := svc.GetRatingRanges(context.Background())
	assert.Equal(t, 1, ratings[0].Value)
}

func TestGetRatingRanges_LastIsFive(t *testing.T) {
	svc := services.NewMotorcycleService(nil, nil)
	ratings, _ := svc.GetRatingRanges(context.Background())
	assert.Equal(t, 5, ratings[len(ratings)-1].Value)
}

// ============================================
// GrantPermission Tests
// ============================================

func TestGrantPermission_NewPermission(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	// No existing permission
	diagMock.On("GetByMotorcycleAndBranch", mock.Anything, "m-1", "b-1").Return(nil, nil)
	diagMock.On("Save", mock.Anything, txMock, mock.AnythingOfType("*domain.DiagnosticPermission")).Return(nil)

	result, err := svc.GrantPermission(context.Background(), txMock, "m-1", "b-1", true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "m-1", result.MotorcycleID)
	assert.Equal(t, "b-1", result.BranchID)
	assert.True(t, result.Active)
	assert.NotEmpty(t, result.ID) // new ID generated
}

func TestGrantPermission_ExistingPermission_ReusesID(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	existing := &domain.DiagnosticPermission{ID: "perm-existing", MotorcycleID: "m-1", BranchID: "b-1"}
	diagMock.On("GetByMotorcycleAndBranch", mock.Anything, "m-1", "b-1").Return(existing, nil)
	diagMock.On("Save", mock.Anything, txMock, mock.AnythingOfType("*domain.DiagnosticPermission")).Return(nil)

	result, err := svc.GrantPermission(context.Background(), txMock, "m-1", "b-1", true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "perm-existing", result.ID) // reuses existing ID
}

func TestGrantPermission_SaveError(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	diagMock.On("GetByMotorcycleAndBranch", mock.Anything, "m-1", "b-1").Return(nil, nil)
	diagMock.On("Save", mock.Anything, txMock, mock.AnythingOfType("*domain.DiagnosticPermission")).Return(errors.New("db err"))

	result, err := svc.GrantPermission(context.Background(), txMock, "m-1", "b-1", true)
	assert.ErrorIs(t, err, domain.ErrPermissionCannotSave)
	assert.Nil(t, result)
}

// ============================================
// RevokePermission Tests
// ============================================

func TestRevokePermission_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	diagMock.On("Delete", mock.Anything, txMock, "m-1", "b-1").Return(nil)

	err := svc.RevokePermission(context.Background(), txMock, "m-1", "b-1")
	assert.NoError(t, err)
}

func TestRevokePermission_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	txMock := &mocks.MockTx{}

	diagMock.On("Delete", mock.Anything, txMock, "m-1", "b-1").Return(errors.New("db err"))

	err := svc.RevokePermission(context.Background(), txMock, "m-1", "b-1")
	assert.Error(t, err)
}

// ============================================
// ListPermissions Tests
// ============================================

func TestListPermissions_Success(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)
	perms := []domain.DiagnosticPermission{{ID: "p-1"}, {ID: "p-2"}}

	diagMock.On("GetByMotorcycleID", mock.Anything, "m-1").Return(perms, nil)

	result, err := svc.ListPermissions(context.Background(), "m-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListPermissions_Empty(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	diagMock.On("GetByMotorcycleID", mock.Anything, "m-1").Return([]domain.DiagnosticPermission{}, nil)

	result, err := svc.ListPermissions(context.Background(), "m-1")
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestListPermissions_Error(t *testing.T) {
	repoMock, diagMock := setupMotorcycleMocks()
	svc := services.NewMotorcycleService(repoMock, diagMock)

	diagMock.On("GetByMotorcycleID", mock.Anything, "m-1").Return(nil, errors.New("db err"))

	result, err := svc.ListPermissions(context.Background(), "m-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}
