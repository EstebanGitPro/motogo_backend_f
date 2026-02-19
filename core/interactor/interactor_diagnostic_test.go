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

// helper to create DiagnosticInteractor with fresh mock service
func setupDiagnosticInteractor() (*interactor.DiagnosticInteractor, *mocks.MockDiagnosticService) {
	svc := new(mocks.MockDiagnosticService)
	di := interactor.NewDiagnosticInteractor(svc)
	return di, svc
}

// ============================================
// RegisterDiagnostic — CREATE path
// ============================================

func TestRegisterDiagnostic_CreateNew_Success(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	desc := "Frenos hacen ruido"
	createdDiag := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ValidateBranchExists", ctx, "branch-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RegisterOrUpdateDiagnostic", ctx, mockTx, "moto-1", "branch-1", &desc).Return(createdDiag, nil)
	mockTx.On("Commit").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", &desc)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRegisterDiagnostic_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-bad", "owner-1").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := di.RegisterDiagnostic(ctx, "moto-bad", "branch-1", "owner-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestRegisterDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "impostor").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "impostor", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestRegisterDiagnostic_BranchNotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ValidateBranchExists", ctx, "branch-bad").Return(domain.ErrBranchNotFound)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-bad", "owner-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestRegisterDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ValidateBranchExists", ctx, "branch-1").Return(nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestRegisterDiagnostic_ServiceError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ValidateBranchExists", ctx, "branch-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RegisterOrUpdateDiagnostic", ctx, mockTx, "moto-1", "branch-1", (*string)(nil)).Return(nil, domain.ErrDiagnosticCannotSave)
	mockTx.On("Rollback").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRegisterDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ValidateBranchExists", ctx, "branch-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RegisterOrUpdateDiagnostic", ctx, mockTx, "moto-1", "branch-1", (*string)(nil)).Return(&domain.Diagnostic{ID: "diag-1"}, nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := di.RegisterDiagnostic(ctx, "moto-1", "branch-1", "owner-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GetDiagnosticByID
// ============================================

func TestGetDiagnosticByID_Success(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	diag := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(diag, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)

	result, err := di.GetDiagnosticByID(ctx, "diag-1", "owner-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
	svc.AssertExpectations(t)
}

func TestGetDiagnosticByID_NotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	result, err := di.GetDiagnosticByID(ctx, "diag-bad", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetDiagnosticByID_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	diag := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(diag, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "impostor").Return(nil, domain.ErrMotorcycleNotFound)

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
	di, svc := setupDiagnosticInteractor()

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
		{ID: "diag-2", MotorcycleID: "moto-1"},
	}

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("GetByMotorcycleID", ctx, "moto-1").Return(diagnostics, nil)

	result, err := di.ListDiagnosticsByMotorcycle(ctx, "moto-1", "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestListDiagnosticsByMotorcycle_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-bad", "owner-1").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := di.ListDiagnosticsByMotorcycle(ctx, "moto-bad", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListDiagnosticsByMotorcycle_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "impostor").Return(nil, domain.ErrMotorcycleNotFound)

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
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	newDescription := "Problema actualizado"
	updates := &domain.Diagnostic{ProblemDescription: &newDescription}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ApplyDiagnosticUpdates", existing, updates).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateDiagnostic", ctx, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUpdateDiagnostic_NotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	result, err := di.UpdateDiagnostic(ctx, "diag-bad", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "impostor").Return(nil, domain.ErrMotorcycleNotFound)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "impostor", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ApplyDiagnosticUpdates", existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

func TestUpdateDiagnostic_UpdateError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ApplyDiagnosticUpdates", existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateDiagnostic", ctx, mockTx, existing).Return(domain.ErrDiagnosticCannotUpdate)
	mockTx.On("Rollback").Return(nil)

	result, err := di.UpdateDiagnostic(ctx, "diag-1", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestUpdateDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("ApplyDiagnosticUpdates", existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateDiagnostic", ctx, mockTx, existing).Return(nil)
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
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteDiagnostic", ctx, mockTx, "diag-1").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteDiagnostic_NotFound(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("GetByID", ctx, "diag-bad").Return(nil, errors.New("not found"))

	err := di.DeleteDiagnostic(ctx, "diag-bad", "owner-1")

	assert.Error(t, err)
}

func TestDeleteDiagnostic_OwnershipError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "impostor").Return(nil, domain.ErrMotorcycleNotFound)

	err := di.DeleteDiagnostic(ctx, "diag-1", "impostor")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDeleteDiagnostic_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}

func TestDeleteDiagnostic_DeleteError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteDiagnostic", ctx, mockTx, "diag-1").Return(domain.ErrDiagnosticCannotDelete)
	mockTx.On("Rollback").Return(nil)

	err := di.DeleteDiagnostic(ctx, "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestDeleteDiagnostic_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	svc.On("GetByID", ctx, "diag-1").Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, "moto-1", "owner-1").Return(&domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteDiagnostic", ctx, mockTx, "diag-1").Return(nil)
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
	di, svc := setupDiagnosticInteractor()

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
	}

	svc.On("GetByMotorcycleID", ctx, "moto-1").Return(diagnostics, nil)

	result, err := di.ListDiagnosticsByMotorcycleID(ctx, "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestListDiagnosticsByMotorcycleID_RepoError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("GetByMotorcycleID", ctx, "moto-1").Return(nil, errors.New("db error"))

	result, err := di.ListDiagnosticsByMotorcycleID(ctx, "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// SetSolution
// ============================================

func TestSetSolution_Success(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SetSolution", ctx, mockTx, "diag-1", "Cambiar pastillas de freno").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := di.SetSolution(ctx, "diag-1", "Cambiar pastillas de freno")

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestSetSolution_BeginTxError(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()

	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := di.SetSolution(ctx, "diag-1", "Solución")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

func TestSetSolution_ServiceError_RollsBack(t *testing.T) {
	ctx := context.Background()
	di, svc := setupDiagnosticInteractor()
	mockTx := new(mocks.MockTx)

	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SetSolution", ctx, mockTx, "diag-1", "Solución").Return(domain.ErrDiagnosticNotFound)
	mockTx.On("Rollback").Return(nil)

	err := di.SetSolution(ctx, "diag-1", "Solución")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
	mockTx.AssertCalled(t, "Rollback")
}
