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
// NewDiagnosticService Tests
// ============================================

func TestNewDiagnosticService(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)

	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)
	assert.NotNil(t, svc)
}

// ============================================
// BeginTx Tests
// ============================================

func TestDiagnosticService_BeginTx_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	tx, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestDiagnosticService_BeginTx_Error(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("db error"))

	tx, err := svc.BeginTx(context.Background())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

// ============================================
// ValidateMotorcycleOwnership Tests
// ============================================

func TestDiagnosticService_ValidateMotorcycleOwnership_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	mockMotoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")
	assert.NoError(t, err)
}

func TestDiagnosticService_ValidateMotorcycleOwnership_NotFound(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockMotoRepo.On("GetByID", mock.Anything, "moto-999").Return(nil, errors.New("not found"))

	err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-999", "owner-1")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDiagnosticService_ValidateMotorcycleOwnership_WrongOwner(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	mockMotoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "wrong-owner")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// ValidateBranchExists Tests
// ============================================

func TestDiagnosticService_ValidateBranchExists_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	branch := &domain.Branch{ID: "branch-1"}
	mockBranchRepo.On("GetBranchByID", mock.Anything, "branch-1").Return(branch, nil)

	err := svc.ValidateBranchExists(context.Background(), "branch-1")
	assert.NoError(t, err)
}

func TestDiagnosticService_ValidateBranchExists_NotFound(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockBranchRepo.On("GetBranchByID", mock.Anything, "branch-999").Return(nil, errors.New("not found"))

	err := svc.ValidateBranchExists(context.Background(), "branch-999")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

// ============================================
// UpsertDiagnostic Tests
// ============================================

func TestDiagnosticService_UpsertDiagnostic_CreateNew(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	desc := "Engine noise"
	urls := []string{"http://img1.jpg"}

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	mockDiagRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	mockDiagRepo.On("SaveEvidence", mock.Anything, mockTx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", &desc, urls)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	assert.Equal(t, "branch-1", result.BranchID)
	assert.Len(t, result.Evidence, 1)
	mockDiagRepo.AssertExpectations(t)
}

func TestDiagnosticService_UpsertDiagnostic_UpdateExisting(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}
	desc := "Updated description"
	urls := []string{"http://new-img.jpg"}

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	mockDiagRepo.On("Update", mock.Anything, mockTx, existing).Return(nil)
	mockDiagRepo.On("DeleteEvidenceByDiagnosticID", mock.Anything, mockTx, "diag-1").Return(nil)
	mockDiagRepo.On("SaveEvidence", mock.Anything, mockTx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", &desc, urls)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
	assert.Equal(t, &desc, result.ProblemDescription)
	mockDiagRepo.AssertExpectations(t)
}

func TestDiagnosticService_UpsertDiagnostic_LookupError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, errors.New("db error"))

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagnosticService_UpsertDiagnostic_UpdateError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	mockDiagRepo.On("Update", mock.Anything, mockTx, existing).Return(errors.New("update error"))

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagnosticService_UpsertDiagnostic_DeleteEvidenceError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	mockDiagRepo.On("Update", mock.Anything, mockTx, existing).Return(nil)
	mockDiagRepo.On("DeleteEvidenceByDiagnosticID", mock.Anything, mockTx, "diag-1").Return(errors.New("delete error"))

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagnosticService_UpsertDiagnostic_SaveNewError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	mockDiagRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(errors.New("save error"))

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagnosticService_UpsertDiagnostic_SaveEvidenceError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	urls := []string{"http://img1.jpg"}

	mockDiagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	mockDiagRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	mockDiagRepo.On("SaveEvidence", mock.Anything, mockTx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(errors.New("save evidence error"))

	result, err := svc.UpsertDiagnostic(context.Background(), mockTx, "moto-1", "branch-1", nil, urls)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

// ============================================
// GetDiagnosticByID Tests
// ============================================

func TestDiagnosticService_GetDiagnosticByID_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diagnostic := &domain.Diagnostic{ID: "diag-1"}
	evidence := []domain.DiagnosticEvidence{{ID: "ev-1"}}

	mockDiagRepo.On("GetByID", mock.Anything, "diag-1").Return(diagnostic, nil)
	mockDiagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-1").Return(evidence, nil)

	result, err := svc.GetDiagnosticByID(context.Background(), "diag-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Evidence, 1)
}

func TestDiagnosticService_GetDiagnosticByID_NotFound(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("GetByID", mock.Anything, "diag-999").Return(nil, domain.ErrDiagnosticNotFound)

	result, err := svc.GetDiagnosticByID(context.Background(), "diag-999")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDiagnosticService_GetDiagnosticByID_EvidenceError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diagnostic := &domain.Diagnostic{ID: "diag-1"}
	mockDiagRepo.On("GetByID", mock.Anything, "diag-1").Return(diagnostic, nil)
	mockDiagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-1").Return(nil, errors.New("evidence error"))

	result, err := svc.GetDiagnosticByID(context.Background(), "diag-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetDiagnosticsByMotorcycleID Tests
// ============================================

func TestDiagnosticService_GetDiagnosticsByMotorcycleID_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diags := []domain.Diagnostic{
		{ID: "diag-1"},
		{ID: "diag-2"},
	}
	ev1 := []domain.DiagnosticEvidence{{ID: "ev-1"}}
	ev2 := []domain.DiagnosticEvidence{{ID: "ev-2"}}

	mockDiagRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(diags, nil)
	mockDiagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-1").Return(ev1, nil)
	mockDiagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-2").Return(ev2, nil)

	result, err := svc.GetDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].Evidence, 1)
	assert.Len(t, result[1].Evidence, 1)
}

func TestDiagnosticService_GetDiagnosticsByMotorcycleID_Empty(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return([]domain.Diagnostic{}, nil)

	result, err := svc.GetDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestDiagnosticService_GetDiagnosticsByMotorcycleID_ListError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(nil, errors.New("db error"))

	result, err := svc.GetDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDiagnosticService_GetDiagnosticsByMotorcycleID_EvidenceError(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diags := []domain.Diagnostic{{ID: "diag-1"}}
	mockDiagRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(diags, nil)
	mockDiagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-1").Return(nil, errors.New("evidence error"))

	result, err := svc.GetDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// ApplyDiagnosticUpdates Tests
// ============================================

func TestDiagnosticService_ApplyDiagnosticUpdates_AppliesDescription(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	existing := &domain.Diagnostic{ID: "diag-1"}
	newDesc := "Updated problem"
	updates := &domain.Diagnostic{ProblemDescription: &newDesc}

	svc.ApplyDiagnosticUpdates(existing, updates)

	assert.Equal(t, &newDesc, existing.ProblemDescription)
}

func TestDiagnosticService_ApplyDiagnosticUpdates_NilIgnored(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	originalDesc := "Original"
	existing := &domain.Diagnostic{ID: "diag-1", ProblemDescription: &originalDesc}
	updates := &domain.Diagnostic{}

	svc.ApplyDiagnosticUpdates(existing, updates)

	assert.Equal(t, &originalDesc, existing.ProblemDescription)
}

// ============================================
// UpdateDiagnostic Tests
// ============================================

func TestDiagnosticService_UpdateDiagnostic_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diag := &domain.Diagnostic{ID: "diag-1"}
	mockDiagRepo.On("Update", mock.Anything, mockTx, diag).Return(nil)

	err := svc.UpdateDiagnostic(context.Background(), mockTx, diag)
	assert.NoError(t, err)
}

func TestDiagnosticService_UpdateDiagnostic_Error(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	diag := &domain.Diagnostic{ID: "diag-1"}
	mockDiagRepo.On("Update", mock.Anything, mockTx, diag).Return(errors.New("update error"))

	err := svc.UpdateDiagnostic(context.Background(), mockTx, diag)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

// ============================================
// DeleteDiagnostic Tests
// ============================================

func TestDiagnosticService_DeleteDiagnostic_Success(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("Delete", mock.Anything, mockTx, "diag-1").Return(nil)

	err := svc.DeleteDiagnostic(context.Background(), mockTx, "diag-1")
	assert.NoError(t, err)
}

func TestDiagnosticService_DeleteDiagnostic_Error(t *testing.T) {
	mockDiagRepo := new(mocks.MockDiagnosticRepository)
	mockMotoRepo := new(mocks.MockMotorcycleRepository)
	mockBranchRepo := new(mocks.MockBranchRepository)
	mockTx := new(mocks.MockTx)
	svc := services.NewDiagnosticService(mockDiagRepo, mockMotoRepo, mockBranchRepo)

	mockDiagRepo.On("Delete", mock.Anything, mockTx, "diag-1").Return(errors.New("delete error"))

	err := svc.DeleteDiagnostic(context.Background(), mockTx, "diag-1")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}
