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
// NewDiagnosticInteractor Tests
// ============================================

func TestNewDiagnosticInteractor(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)
	assert.NotNil(t, i)
}

// ============================================
// RegisterDiagnostic Tests (HU11)
// ============================================

func TestRegisterDiagnostic_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	desc := "Engine noise"
	urls := []string{"http://img1.jpg"}
	created := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ValidateBranchExists", mock.Anything, "branch-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpsertDiagnostic", mock.Anything, mockTx, "moto-1", "branch-1", &desc, urls).Return(created, nil)
	mockTx.On("Commit").Return(nil)

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-1", "owner-1", &desc, urls)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
	mockService.AssertExpectations(t)
}

func TestRegisterDiagnostic_OwnershipError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(domain.ErrMotorcycleNotFound)

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-1", "wrong-owner", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestRegisterDiagnostic_BranchError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ValidateBranchExists", mock.Anything, "branch-999").Return(domain.ErrBranchNotFound)

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-999", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestRegisterDiagnostic_TxError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ValidateBranchExists", mock.Anything, "branch-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestRegisterDiagnostic_UpsertError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ValidateBranchExists", mock.Anything, "branch-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpsertDiagnostic", mock.Anything, mockTx, "moto-1", "branch-1", (*string)(nil), []string(nil)).Return(nil, domain.ErrDiagnosticCannotSave)
	mockTx.On("Rollback").Return(nil)

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRegisterDiagnostic_CommitError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	created := &domain.Diagnostic{ID: "diag-1"}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ValidateBranchExists", mock.Anything, "branch-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpsertDiagnostic", mock.Anything, mockTx, "moto-1", "branch-1", (*string)(nil), []string(nil)).Return(created, nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	result, err := i.RegisterDiagnostic(context.Background(), "moto-1", "branch-1", "owner-1", nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

// ============================================
// GetDiagnosticByID Tests (HU14)
// ============================================

func TestGetDiagnosticByID_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	diagnostic := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(diagnostic, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)

	result, err := i.GetDiagnosticByID(context.Background(), "diag-1", "owner-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-1", result.ID)
}

func TestGetDiagnosticByID_NotFound(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-999").Return(nil, domain.ErrDiagnosticNotFound)

	result, err := i.GetDiagnosticByID(context.Background(), "diag-999", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetDiagnosticByID_OwnershipFail(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	diagnostic := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(diagnostic, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(domain.ErrMotorcycleNotFound)

	result, err := i.GetDiagnosticByID(context.Background(), "diag-1", "wrong-owner")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

// ============================================
// ListDiagnosticsByMotorcycle Tests (HU14)
// ============================================

func TestListDiagnosticsByMotorcycle_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1", MotorcycleID: "moto-1"},
		{ID: "diag-2", MotorcycleID: "moto-1"},
	}

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("GetDiagnosticsByMotorcycleID", mock.Anything, "moto-1").Return(diagnostics, nil)

	result, err := i.ListDiagnosticsByMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListDiagnosticsByMotorcycle_OwnershipError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(domain.ErrMotorcycleNotFound)

	result, err := i.ListDiagnosticsByMotorcycle(context.Background(), "moto-1", "wrong-owner")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListDiagnosticsByMotorcycle_ListError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("GetDiagnosticsByMotorcycleID", mock.Anything, "moto-1").Return(nil, errors.New("db error"))

	result, err := i.ListDiagnosticsByMotorcycle(context.Background(), "moto-1", "owner-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// UpdateDiagnostic Tests (HU12)
// ============================================

func TestUpdateDiagnostic_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	desc := "Updated description"
	updates := &domain.Diagnostic{ProblemDescription: &desc}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ApplyDiagnosticUpdates", existing, updates).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := i.UpdateDiagnostic(context.Background(), "diag-1", "owner-1", updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockService.AssertExpectations(t)
}

func TestUpdateDiagnostic_NotFound(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-999").Return(nil, domain.ErrDiagnosticNotFound)

	result, err := i.UpdateDiagnostic(context.Background(), "diag-999", "owner-1", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_OwnershipFail(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(domain.ErrMotorcycleNotFound)

	result, err := i.UpdateDiagnostic(context.Background(), "diag-1", "wrong-owner", &domain.Diagnostic{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestUpdateDiagnostic_TxError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	updates := &domain.Diagnostic{}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ApplyDiagnosticUpdates", existing, updates).Return()
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := i.UpdateDiagnostic(context.Background(), "diag-1", "owner-1", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

func TestUpdateDiagnostic_UpdateError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	updates := &domain.Diagnostic{}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ApplyDiagnosticUpdates", existing, updates).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(domain.ErrDiagnosticCannotUpdate)
	mockTx.On("Rollback").Return(nil)

	result, err := i.UpdateDiagnostic(context.Background(), "diag-1", "owner-1", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateDiagnostic_CommitError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	updates := &domain.Diagnostic{}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("ApplyDiagnosticUpdates", existing, updates).Return()
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	result, err := i.UpdateDiagnostic(context.Background(), "diag-1", "owner-1", updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

// ============================================
// DeleteDiagnostic Tests (HU13)
// ============================================

func TestDeleteDiagnostic_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteDiagnostic", mock.Anything, mockTx, "diag-1").Return(nil)
	mockTx.On("Commit").Return(nil)

	err := i.DeleteDiagnostic(context.Background(), "diag-1", "owner-1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeleteDiagnostic_NotFound(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-999").Return(nil, domain.ErrDiagnosticNotFound)

	err := i.DeleteDiagnostic(context.Background(), "diag-999", "owner-1")

	assert.Error(t, err)
}

func TestDeleteDiagnostic_OwnershipFail(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "wrong-owner").Return(domain.ErrMotorcycleNotFound)

	err := i.DeleteDiagnostic(context.Background(), "diag-1", "wrong-owner")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDeleteDiagnostic_TxError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	err := i.DeleteDiagnostic(context.Background(), "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}

func TestDeleteDiagnostic_DeleteError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteDiagnostic", mock.Anything, mockTx, "diag-1").Return(domain.ErrDiagnosticCannotDelete)
	mockTx.On("Rollback").Return(nil)

	err := i.DeleteDiagnostic(context.Background(), "diag-1", "owner-1")

	assert.Error(t, err)
}

func TestDeleteDiagnostic_CommitError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("ValidateMotorcycleOwnership", mock.Anything, "moto-1", "owner-1").Return(nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("DeleteDiagnostic", mock.Anything, mockTx, "diag-1").Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := i.DeleteDiagnostic(context.Background(), "diag-1", "owner-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}

// ============================================
// ListDiagnosticsByMotorcycleID Tests (Workshop)
// ============================================

func TestListDiagnosticsByMotorcycleID_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	diagnostics := []domain.Diagnostic{
		{ID: "diag-1"},
		{ID: "diag-2"},
	}

	mockService.On("GetDiagnosticsByMotorcycleID", mock.Anything, "moto-1").Return(diagnostics, nil)

	result, err := i.ListDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListDiagnosticsByMotorcycleID_Error(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("GetDiagnosticsByMotorcycleID", mock.Anything, "moto-1").Return(nil, errors.New("db error"))

	result, err := i.ListDiagnosticsByMotorcycleID(context.Background(), "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// SetDiagnosticSolution Tests (Admin)
// ============================================

func TestSetDiagnosticSolution_Success(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}
	solution := "Replace spark plugs"

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := i.SetDiagnosticSolution(context.Background(), "diag-1", &solution)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &solution, result.PossibleSolution)
	mockService.AssertExpectations(t)
}

func TestSetDiagnosticSolution_NotFound(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-999").Return(nil, domain.ErrDiagnosticNotFound)

	result, err := i.SetDiagnosticSolution(context.Background(), "diag-999", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestSetDiagnosticSolution_TxError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := i.SetDiagnosticSolution(context.Background(), "diag-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

func TestSetDiagnosticSolution_UpdateError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(domain.ErrDiagnosticCannotUpdate)
	mockTx.On("Rollback").Return(nil)

	result, err := i.SetDiagnosticSolution(context.Background(), "diag-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSetDiagnosticSolution_CommitError(t *testing.T) {
	mockService := new(mocks.MockDiagnosticService)
	mockTx := new(mocks.MockTx)
	i := interactor.NewDiagnosticInteractor(mockService)

	existing := &domain.Diagnostic{ID: "diag-1"}

	mockService.On("GetDiagnosticByID", mock.Anything, "diag-1").Return(existing, nil)
	mockService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockService.On("UpdateDiagnostic", mock.Anything, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	result, err := i.SetDiagnosticSolution(context.Background(), "diag-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}
