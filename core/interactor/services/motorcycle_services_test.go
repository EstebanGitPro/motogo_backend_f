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
// Constructor Tests
// ============================================

func TestNewMotorcycleService(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)

	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)
	assert.NotNil(t, svc)
}

func TestMotorcycleService_WithStorageClient(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockStorage := new(mocks.MockStorageClient)

	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)
	svc.WithStorageClient(mockStorage)
	// No panic = storage client set correctly
}

// ============================================
// BeginTx / BeginPermissionTx Tests
// ============================================

func TestMotorcycleService_BeginTx_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	tx, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestMotorcycleService_BeginPermissionTx_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockPermRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	tx, err := svc.BeginPermissionTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

// ============================================
// ValidateMotorcycleOwnership Tests
// ============================================

func TestMotorcycleService_ValidateMotorcycleOwnership_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	mockMotoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.ID)
}

func TestMotorcycleService_ValidateMotorcycleOwnership_NotFound(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("GetByID", mock.Anything, "moto-999").Return(nil, errors.New("not found"))

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-999", "owner-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMotorcycleService_ValidateMotorcycleOwnership_WrongOwner(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	mockMotoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "wrong-owner")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// ValidateReferenceExists Tests
// ============================================

func TestMotorcycleService_ValidateReferenceExists_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(true, nil)

	err := svc.ValidateReferenceExists(context.Background(), "ref-1")
	assert.NoError(t, err)
}

func TestMotorcycleService_ValidateReferenceExists_Empty(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	err := svc.ValidateReferenceExists(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrReferenceRequired, err)
}

func TestMotorcycleService_ValidateReferenceExists_RepoError(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-1").Return(false, errors.New("db error"))

	err := svc.ValidateReferenceExists(context.Background(), "ref-1")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

func TestMotorcycleService_ValidateReferenceExists_NotFound(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-999").Return(false, nil)

	err := svc.ValidateReferenceExists(context.Background(), "ref-999")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrReferenceNotFound, err)
}

// ============================================
// ValidateLicensePlateUnique Tests
// ============================================

func TestMotorcycleService_ValidateLicensePlateUnique_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(false, nil)

	err := svc.ValidateLicensePlateUnique(context.Background(), "ABC123")
	assert.NoError(t, err)
}

func TestMotorcycleService_ValidateLicensePlateUnique_Duplicate(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(true, nil)

	err := svc.ValidateLicensePlateUnique(context.Background(), "ABC123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateLicensePlate, err)
}

func TestMotorcycleService_ValidateLicensePlateUnique_RepoError(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("CheckLicensePlateExists", mock.Anything, "ABC123").Return(false, errors.New("db error"))

	err := svc.ValidateLicensePlateUnique(context.Background(), "ABC123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

// ============================================
// CreateMotorcycle Tests
// ============================================

func TestMotorcycleService_CreateMotorcycle_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{LicensePlate: "ABC123", ReferenceID: "ref-1", OwnerID: "owner-1"}
	mockMotoRepo.On("Save", mock.Anything, mockTx, moto).Return(nil)

	result, err := svc.CreateMotorcycle(context.Background(), mockTx, moto)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID) // UUID generated
}

func TestMotorcycleService_CreateMotorcycle_SaveError(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{LicensePlate: "ABC123"}
	mockMotoRepo.On("Save", mock.Anything, mockTx, moto).Return(errors.New("save error"))

	result, err := svc.CreateMotorcycle(context.Background(), mockTx, moto)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

// ============================================
// Read Delegation Tests
// ============================================

func TestMotorcycleService_GetMotorcycleByID_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1"}
	mockMotoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.GetMotorcycleByID(context.Background(), "moto-1")
	assert.NoError(t, err)
	assert.Equal(t, "moto-1", result.ID)
}

func TestMotorcycleService_GetMotorcyclesByOwner_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	motos := []domain.Motorcycle{{ID: "moto-1"}, {ID: "moto-2"}}
	mockMotoRepo.On("GetByOwnerID", mock.Anything, "owner-1").Return(motos, nil)

	result, err := svc.GetMotorcyclesByOwner(context.Background(), "owner-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMotorcycleService_GetMotorcycleByLicensePlate_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", LicensePlate: "ABC123"}
	mockMotoRepo.On("GetByLicensePlate", mock.Anything, "ABC123").Return(moto, nil)

	result, err := svc.GetMotorcycleByLicensePlate(context.Background(), "ABC123")
	assert.NoError(t, err)
	assert.Equal(t, "ABC123", result.LicensePlate)
}

// ============================================
// ApplyMotorcycleUpdates Tests
// ============================================

func TestMotorcycleService_ApplyMotorcycleUpdates_Year(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	year := 2024
	updates := &domain.Motorcycle{Year: &year}

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.NoError(t, err)
	assert.Equal(t, &year, moto.Year)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_Mileage(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	mileage := 15000
	updates := &domain.Motorcycle{CurrentMileage: &mileage}

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.NoError(t, err)
	assert.Equal(t, &mileage, moto.CurrentMileage)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_Notes(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	notes := "Some owner notes"
	updates := &domain.Motorcycle{OwnerNotes: &notes}

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.NoError(t, err)
	assert.Equal(t, &notes, moto.OwnerNotes)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_ProfileImage(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	imgURL := "http://img.jpg"
	updates := &domain.Motorcycle{ProfileImageURL: &imgURL}

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.NoError(t, err)
	assert.Equal(t, &imgURL, moto.ProfileImageURL)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_ReferenceChange_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-2"}

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-2").Return(true, nil)

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.NoError(t, err)
	assert.Equal(t, "ref-2", moto.ReferenceID)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_ReferenceChange_NotFound(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-999"}

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-999").Return(false, nil)

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrReferenceNotFound, err)
}

func TestMotorcycleService_ApplyMotorcycleUpdates_ReferenceChange_RepoError(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1", ReferenceID: "ref-1"}
	updates := &domain.Motorcycle{ReferenceID: "ref-2"}

	mockMotoRepo.On("ValidateReferenceExists", mock.Anything, "ref-2").Return(false, errors.New("db error"))

	err := svc.ApplyMotorcycleUpdates(context.Background(), moto, updates)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotUpdate, err)
}

// ============================================
// UpdateMotorcycle Tests
// ============================================

func TestMotorcycleService_UpdateMotorcycle_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	moto := &domain.Motorcycle{ID: "moto-1"}
	mockMotoRepo.On("Update", mock.Anything, mockTx, moto).Return(nil)

	err := svc.UpdateMotorcycle(context.Background(), mockTx, moto)
	assert.NoError(t, err)
}

// ============================================
// CheckServiceHistory Tests
// ============================================

func TestMotorcycleService_CheckServiceHistory_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("HasServiceHistory", mock.Anything, "moto-1").Return(true, nil)

	has, err := svc.CheckServiceHistory(context.Background(), "moto-1")
	assert.NoError(t, err)
	assert.True(t, has)
}

// ============================================
// DeleteMotorcycle Tests
// ============================================

func TestMotorcycleService_DeleteMotorcycle_SoftDelete(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("Delete", mock.Anything, mockTx, "moto-1").Return(nil)

	err := svc.DeleteMotorcycle(context.Background(), mockTx, "moto-1", true)
	assert.NoError(t, err)
	mockMotoRepo.AssertCalled(t, "Delete", mock.Anything, mockTx, "moto-1")
}

func TestMotorcycleService_DeleteMotorcycle_HardDelete(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("HardDelete", mock.Anything, mockTx, "moto-1").Return(nil)

	err := svc.DeleteMotorcycle(context.Background(), mockTx, "moto-1", false)
	assert.NoError(t, err)
	mockMotoRepo.AssertCalled(t, "HardDelete", mock.Anything, mockTx, "moto-1")
}

// ============================================
// DeleteProfileImage & DeleteStorageFile Tests
// ============================================

func TestMotorcycleService_DeleteProfileImage_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockMotoRepo.On("ClearProfileImageURL", mock.Anything, mockTx, "moto-1").Return(nil)

	err := svc.DeleteProfileImage(context.Background(), mockTx, "moto-1")
	assert.NoError(t, err)
}

func TestMotorcycleService_DeleteStorageFile_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockStorage := new(mocks.MockStorageClient)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)
	svc.WithStorageClient(mockStorage)

	mockStorage.On("DeleteStorageFile", mock.Anything, "http://img.jpg").Return(nil)

	// Should not panic
	svc.DeleteStorageFile(context.Background(), "http://img.jpg")
	mockStorage.AssertCalled(t, "DeleteStorageFile", mock.Anything, "http://img.jpg")
}

func TestMotorcycleService_DeleteStorageFile_EmptyURL(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockStorage := new(mocks.MockStorageClient)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)
	svc.WithStorageClient(mockStorage)

	// Should not call storage for empty URL
	svc.DeleteStorageFile(context.Background(), "")
	mockStorage.AssertNotCalled(t, "DeleteStorageFile")
}

func TestMotorcycleService_DeleteStorageFile_NilClient(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	// Should not panic even without storage client
	svc.DeleteStorageFile(context.Background(), "http://img.jpg")
}

// ============================================
// Reference Catalog Tests
// ============================================

func TestMotorcycleService_GetAllReferences_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	refs := []domain.MotorcycleReference{{ID: "ref-1"}}
	mockMotoRepo.On("GetAllReferences", mock.Anything).Return(refs, nil)

	result, err := svc.GetAllReferences(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestMotorcycleService_GetReferencesByBrandID_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	refs := []domain.MotorcycleReference{{ID: "ref-1"}}
	mockMotoRepo.On("GetReferencesByBrandID", mock.Anything, "brand-1").Return(refs, nil)

	result, err := svc.GetReferencesByBrandID(context.Background(), "brand-1")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// ============================================
// Diagnostic Permission Tests
// ============================================

func TestMotorcycleService_GrantDiagnosticPermission_Success(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockPermRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.DiagnosticPermission")).Return(nil)

	result, err := svc.GrantDiagnosticPermission(context.Background(), mockTx, "moto-1", "branch-1", true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	assert.Equal(t, "branch-1", result.BranchID)
	assert.True(t, result.Active)
	assert.NotEmpty(t, result.ID)
}

func TestMotorcycleService_GrantDiagnosticPermission_SaveError(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockPermRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.DiagnosticPermission")).Return(errors.New("save error"))

	result, err := svc.GrantDiagnosticPermission(context.Background(), mockTx, "moto-1", "branch-1", true)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrPermissionCannotSave, err)
}

func TestMotorcycleService_RevokeDiagnosticPermission_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockPermRepo.On("Deactivate", mock.Anything, mockTx, "moto-1", "branch-1").Return(nil)

	err := svc.RevokeDiagnosticPermission(context.Background(), mockTx, "moto-1", "branch-1")
	assert.NoError(t, err)
}

func TestMotorcycleService_ListDiagnosticPermissions_Delegates(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	perms := []domain.DiagnosticPermission{{ID: "perm-1"}}
	mockPermRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(perms, nil)

	result, err := svc.ListDiagnosticPermissions(context.Background(), "moto-1")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// ============================================
// ValidateBranchPermission Tests
// ============================================

func TestMotorcycleService_ValidateBranchPermission_Authorized(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	perm := &domain.DiagnosticPermission{ID: "perm-1", Active: true}
	mockPermRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(perm, nil)

	err := svc.ValidateBranchPermission(context.Background(), "moto-1", []string{"branch-1"})
	assert.NoError(t, err)
}

func TestMotorcycleService_ValidateBranchPermission_NotAuthorized(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	mockPermRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, errors.New("not found"))

	err := svc.ValidateBranchPermission(context.Background(), "moto-1", []string{"branch-1"})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotAuthorized, err)
}

func TestMotorcycleService_ValidateBranchPermission_MultipleBranches(t *testing.T) {
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockPermRepo := new(mocks.MockDiagnosticPermissionRepository)
	svc := services.NewMotorcycleService(mockMotoRepo, mockPermRepo)

	// First branch is not authorized, second is
	mockPermRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, errors.New("not found"))
	perm := &domain.DiagnosticPermission{ID: "perm-2", Active: true}
	mockPermRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-2").Return(perm, nil)

	err := svc.ValidateBranchPermission(context.Background(), "moto-1", []string{"branch-1", "branch-2"})
	assert.NoError(t, err) // At least one branch authorized
}
