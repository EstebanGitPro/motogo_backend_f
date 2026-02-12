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

// helper to create DiagnosticInteractor with fresh mocks
func setupDiagnosticInteractor() (*interactor.DiagnosticInteractor, *mocks.MockDiagnosticRepository, *mocks.MockMotorcycleRepository, *mocks.MockBranchRepository) {
	diagRepo := new(mocks.MockDiagnosticRepository)
	motoRepo := new(mocks.MockMotorcycleRepository)
	branchRepo := new(mocks.MockBranchRepository)
	di := interactor.NewDiagnosticInteractor(diagRepo, motoRepo, branchRepo)
	return di, diagRepo, motoRepo, branchRepo
}

// ============================================
// RegisterDiagnostic — CREATE path
// ============================================

func TestRegisterDiagnostic_CreateNew_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, branchRepo := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	desc := "Frenos hacen ruido"

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)
	diagRepo.On("GetByMotorcycleAndBranch", ctx, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Save", ctx, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	diagRepo.On("SaveEvidence", ctx, mockTx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", &desc, []string{"http://img1.jpg"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	assert.Equal(t, "branch-1", result.BranchID)
	assert.Len(t, result.Evidence, 1)
	motoRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
	diagRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterDiagnostic_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	di, _, motoRepo, _ := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-bad").Return(nil, errors.New("not found"))

	result, err := di.RegisterDiagnostic(ctx, "moto-bad", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestRegisterDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, _, motoRepo, _ := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "real-owner"}, nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "impostor", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestRegisterDiagnostic_BranchNotFound(t *testing.T) {
	ctx := context.Background()
	di, _, motoRepo, branchRepo := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-bad").Return(nil, errors.New("not found"))

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-bad", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestRegisterDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, branchRepo := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)
	diagRepo.On("GetByMotorcycleAndBranch", ctx, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestRegisterDiagnostic_SaveError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, branchRepo := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)
	diagRepo.On("GetByMotorcycleAndBranch", ctx, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Save", ctx, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(errors.New("save failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRegisterDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, branchRepo := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)
	diagRepo.On("GetByMotorcycleAndBranch", ctx, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Save", ctx, mockTx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// RegisterDiagnostic — UPSERT path
// ============================================

func TestRegisterDiagnostic_Upsert_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, branchRepo := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	desc := "Problema actualizado"
	existing := &domain.Diagnostic{ID: "diag-existing", MotorcycleID: "moto-1", BranchID: "branch-1"}

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	branchRepo.On("GetBranchByID", ctx, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)
	diagRepo.On("GetByMotorcycleAndBranch", ctx, "moto-1", "branch-1").Return(existing, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Update", ctx, mockTx, existing).Return(nil)
	diagRepo.On("DeleteEvidenceByDiagnosticID", ctx, mockTx, "diag-existing").Return(nil)
	diagRepo.On("SaveEvidence", ctx, mockTx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", &desc, []string{"http://new-img.jpg"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-existing", result.ID)
	assert.Len(t, result.Evidence, 1)
	diagRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================
// GetDiagnosticByID
// ============================================

func TestGetDiagnosticByID_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	diag := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	evidence := []domain.DiagnosticEvidence{{ID: "ev-1", DiagnosticID: "diag-1"}}

	diagRepo.On("GetByID", ctx, "diag-1").Return(diag, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", ctx, "diag-1").Return(evidence, nil)

	result, err := di.GetDiagnosticByID(ctx, "diag-1", "owner-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
	assert.Len(t, result.Evidence, 1)
	diagRepo.AssertExpectations(t)
	motoRepo.AssertExpectations(t)
}

func TestGetDiagnosticByID_NotFound(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagRepo.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	result, err := di.GetDiagnosticByID(ctx, "diag-bad", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetDiagnosticByID_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	diag := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(diag, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "real-owner"}, nil)

	result, err := di.GetDiagnosticByID(ctx, "diag-1", "impostor")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

// ============================================
// ListDiagnosticsByMotorcycle
// ============================================

func TestListDiagnosticsByMotorcycle_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
		{ID: "diag-2", MotorcycleID: "moto-1"},
	}

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("GetByMotorcycleID", ctx, "moto-1").Return(diagnostics, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", ctx, "diag-1").Return([]domain.DiagnosticEvidence{}, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", ctx, "diag-2").Return([]domain.DiagnosticEvidence{}, nil)

	result, err := di.ListDiagnosticsByMotorcycle(ctx, "moto-1", "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	diagRepo.AssertExpectations(t)
	motoRepo.AssertExpectations(t)
}

func TestListDiagnosticsByMotorcycle_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	di, _, motoRepo, _ := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-bad").Return(nil, errors.New("not found"))

	result, err := di.ListDiagnosticsByMotorcycle(ctx, "moto-bad", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListDiagnosticsByMotorcycle_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, _, motoRepo, _ := setupDiagnosticInteractor()

	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "real-owner"}, nil)

	result, err := di.ListDiagnosticsByMotorcycle(ctx, "moto-1", "impostor")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// UpdateDiagnostic
// ============================================

func TestUpdateDiagnostic_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	newDescription := "Problema actualizado"
	updates := &domain.Diagnostic{ProblemDescription: &newDescription}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Update", ctx, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &newDescription, result.ProblemDescription)
	diagRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUpdateDiagnostic_NotFound(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagRepo.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	result, err := di.UpdateDiagnostic(ctx, "diag-bad", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "real-owner"}, nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "impostor", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

func TestUpdateDiagnostic_UpdateError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Update", ctx, mockTx, existing).Return(errors.New("update failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestUpdateDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Update", ctx, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// DeleteDiagnostic
// ============================================

func TestDeleteDiagnostic_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Delete", ctx, mockTx, "diag-1").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.NoError(t, err)
	diagRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteDiagnostic_NotFound(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagRepo.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	err := di.DeleteDiagnostic(ctx, "diag-bad", "owner-1")

	assert.Error(t, err)
}

func TestDeleteDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "real-owner"}, nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "impostor")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDeleteDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}

func TestDeleteDiagnostic_DeleteError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Delete", ctx, mockTx, "diag-1").Return(errors.New("delete failed"))
	mockTx.On("Rollback").Return(nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestDeleteDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, motoRepo, _ := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	diagRepo.On("GetByID", ctx, "diag-1").Return(existing, nil)
	motoRepo.On("GetByID", ctx, "moto-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	diagRepo.On("BeginTx", ctx).Return(mockTx, nil)
	diagRepo.On("Delete", ctx, mockTx, "diag-1").Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// ListDiagnosticsByMotorcycleID (no ownership check — workshop use)
// ============================================

func TestListDiagnosticsByMotorcycleID_Success(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
	}

	diagRepo.On("GetByMotorcycleID", ctx, "moto-1").Return(diagnostics, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", ctx, "diag-1").Return([]domain.DiagnosticEvidence{}, nil)

	result, err := di.ListDiagnosticsByMotorcycleID(ctx, "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	diagRepo.AssertExpectations(t)
}

func TestListDiagnosticsByMotorcycleID_RepoError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagRepo.On("GetByMotorcycleID", ctx, "moto-1").Return(nil, errors.New("db error"))

	result, err := di.ListDiagnosticsByMotorcycleID(ctx, "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListDiagnosticsByMotorcycleID_EvidenceError(t *testing.T) {
	ctx := context.Background()
	di, diagRepo, _, _ := setupDiagnosticInteractor()

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
	}

	diagRepo.On("GetByMotorcycleID", ctx, "moto-1").Return(diagnostics, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", ctx, "diag-1").Return(nil, errors.New("evidence error"))

	result, err := di.ListDiagnosticsByMotorcycleID(ctx, "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}
