package interactor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper to create CompletedServiceInteractor with fresh mock service
func setupCompletedServiceInteractor() (*interactor.CompletedServiceInteractor, *mocks.MockCompletedServiceService) {
	svc := new(mocks.MockCompletedServiceService)
	ci := interactor.NewCompletedServiceInteractor(svc)
	return ci, svc
}

// ============================================
// RegisterCompletedService — success
// ============================================

func TestRegisterCompletedService_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	diagID := "diag-1"
	cs := &domain.CompletedService{
		BranchID:     "branch-1",
		MotorcycleID: "moto-1",
		DiagnosticID: &diagID,
	}
	serviceIDs := []string{"svc-1", "svc-2"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateDiagnosticForMotorcycle", ctx, "diag-1", "moto-1").Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(nil)
	svc.On("SaveItems", ctx, mockTx, mock.AnythingOfType("[]domain.CompletedServiceItem")).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, domain.ServiceStatusPending, result.Status)
	assert.Len(t, result.Services, 2)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// ============================================
// RegisterCompletedService — without diagnostic
// ============================================

func TestRegisterCompletedService_NoDiagnostic_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{
		BranchID:     "branch-1",
		MotorcycleID: "moto-1",
		DiagnosticID: nil,
	}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	// ValidateDiagnosticForMotorcycle should NOT be called (DiagnosticID is nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(nil)
	svc.On("SaveItems", ctx, mockTx, mock.AnythingOfType("[]domain.CompletedServiceItem")).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	svc.AssertNotCalled(t, "ValidateDiagnosticForMotorcycle", mock.Anything, mock.Anything, mock.Anything)
	svc.AssertExpectations(t)
}

// ============================================
// RegisterCompletedService — validation errors
// ============================================

func TestRegisterCompletedService_ValidateBranchServicesError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-bad"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(domain.ErrInvalidBranchServices)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidBranchServices, err)
}

func TestRegisterCompletedService_ValidateDiagnosticError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	diagID := "diag-bad"
	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1", DiagnosticID: &diagID}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateDiagnosticForMotorcycle", ctx, "diag-bad", "moto-1").Return(domain.ErrDiagnosticNotForMotorcycle)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticNotForMotorcycle, err)
}

func TestRegisterCompletedService_ActiveServiceExists(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(domain.ErrActiveServiceExists)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrActiveServiceExists, err)
}

// ============================================
// RegisterCompletedService — transaction errors
// ============================================

func TestRegisterCompletedService_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRegisterCompletedService_SaveError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(errors.New("save error"))
	mockTx.On("Rollback").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRegisterCompletedService_SaveItemsError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(nil)
	svc.On("SaveItems", ctx, mockTx, mock.AnythingOfType("[]domain.CompletedServiceItem")).Return(errors.New("items error"))
	mockTx.On("Rollback").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRegisterCompletedService_SaveHistoryError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(nil)
	svc.On("SaveItems", ctx, mockTx, mock.AnythingOfType("[]domain.CompletedServiceItem")).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(errors.New("history error"))
	mockTx.On("Rollback").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRegisterCompletedService_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{BranchID: "branch-1", MotorcycleID: "moto-1"}
	serviceIDs := []string{"svc-1"}

	svc.On("ValidateBranchServices", ctx, "branch-1", serviceIDs).Return(nil)
	svc.On("ValidateNoActiveService", ctx, "branch-1", "moto-1").Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("SaveCompletedService", ctx, mockTx, cs).Return(nil)
	svc.On("SaveItems", ctx, mockTx, mock.AnythingOfType("[]domain.CompletedServiceItem")).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit failed"))
	mockTx.On("Rollback").Return(nil)

	result, err := ci.RegisterCompletedService(ctx, cs, serviceIDs, "person-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GetCompletedServiceByID
// ============================================

func TestGetCompletedServiceByID_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", BranchID: "branch-1"}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)

	result, err := ci.GetCompletedServiceByID(ctx, "cs-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cs-1", result.ID)
	svc.AssertExpectations(t)
}

func TestGetCompletedServiceByID_Error(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByID", ctx, "cs-bad").Return(nil, errors.New("not found"))

	result, err := ci.GetCompletedServiceByID(ctx, "cs-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetCompletedServicesByBranch
// ============================================

func TestGetCompletedServicesByBranch_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	services := []domain.CompletedService{
		{ID: "cs-1", BranchID: "branch-1"},
		{ID: "cs-2", BranchID: "branch-1"},
	}
	svc.On("GetByBranchID", ctx, "branch-1").Return(services, nil)

	result, err := ci.GetCompletedServicesByBranch(ctx, "branch-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestGetCompletedServicesByBranch_Error(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByBranchID", ctx, "branch-bad").Return(nil, errors.New("db error"))

	result, err := ci.GetCompletedServicesByBranch(ctx, "branch-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetCompletedServicesByMotorcycle
// ============================================

func TestGetCompletedServicesByMotorcycle_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	services := []domain.CompletedService{
		{ID: "cs-1", MotorcycleID: "moto-1"},
	}
	svc.On("GetByMotorcycleID", ctx, "moto-1").Return(services, nil)

	result, err := ci.GetCompletedServicesByMotorcycle(ctx, "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestGetCompletedServicesByMotorcycle_Error(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByMotorcycleID", ctx, "moto-bad").Return(nil, errors.New("db error"))

	result, err := ci.GetCompletedServicesByMotorcycle(ctx, "moto-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// DeleteCompletedService (HU65)
// ============================================

func TestDeleteCompletedService_Success_Pendiente(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteCompletedService", ctx, mockTx, "cs-1", domain.ServiceStatusPending).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := ci.DeleteCompletedService(ctx, "cs-1")

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteCompletedService_Success_Finalizado(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteCompletedService", ctx, mockTx, "cs-1", domain.ServiceStatusCompleted).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := ci.DeleteCompletedService(ctx, "cs-1")

	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestDeleteCompletedService_NotFound(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByID", ctx, "cs-bad").Return(nil, errors.New("not found"))

	err := ci.DeleteCompletedService(ctx, "cs-bad")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceNotFound, err)
}

func TestDeleteCompletedService_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := ci.DeleteCompletedService(ctx, "cs-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotDelete, err)
}

func TestDeleteCompletedService_DeleteError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteCompletedService", ctx, mockTx, "cs-1", domain.ServiceStatusPending).Return(errors.New("delete error"))
	mockTx.On("Rollback").Return(nil)

	err := ci.DeleteCompletedService(ctx, "cs-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestDeleteCompletedService_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteCompletedService", ctx, mockTx, "cs-1", domain.ServiceStatusPending).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := ci.DeleteCompletedService(ctx, "cs-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// TransitionStatus (HU74)
// ============================================

func TestTransitionStatus_Success_PendienteToEnProceso(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateStatus", ctx, mockTx, "cs-1", "EN_PROCESO", mock.Anything).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := ci.TransitionStatus(ctx, "cs-1", "EN_PROCESO", "person-1", nil)

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTransitionStatus_Success_EnProcesoToFinalizado(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusInProgress}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	// When transitioning to FINALIZADO, completion_date should NOT be nil
	svc.On("UpdateStatus", ctx, mockTx, "cs-1", "FINALIZADO", mock.MatchedBy(func(t *time.Time) bool {
		return t != nil
	})).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := ci.TransitionStatus(ctx, "cs-1", "FINALIZADO", "person-1", nil)

	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestTransitionStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByID", ctx, "cs-bad").Return(nil, errors.New("not found"))

	err := ci.TransitionStatus(ctx, "cs-bad", "EN_PROCESO", "person-1", nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceNotFound, err)
}

func TestTransitionStatus_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	// PENDIENTE → FINALIZADO is not allowed
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)

	err := ci.TransitionStatus(ctx, "cs-1", "FINALIZADO", "person-1", nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidStatusTransition, err)
}

func TestTransitionStatus_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := ci.TransitionStatus(ctx, "cs-1", "EN_PROCESO", "person-1", nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidStatusTransition, err)
}

func TestTransitionStatus_UpdateError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateStatus", ctx, mockTx, "cs-1", "EN_PROCESO", mock.Anything).Return(errors.New("update error"))
	mockTx.On("Rollback").Return(nil)

	err := ci.TransitionStatus(ctx, "cs-1", "EN_PROCESO", "person-1", nil)

	assert.Error(t, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestTransitionStatus_SaveHistoryError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateStatus", ctx, mockTx, "cs-1", "EN_PROCESO", mock.Anything).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(errors.New("history error"))
	mockTx.On("Rollback").Return(nil)

	err := ci.TransitionStatus(ctx, "cs-1", "EN_PROCESO", "person-1", nil)

	assert.Error(t, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestTransitionStatus_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()
	mockTx := new(mocks.MockTx)

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateStatus", ctx, mockTx, "cs-1", "EN_PROCESO", mock.Anything).Return(nil)
	svc.On("SaveStatusHistory", ctx, mockTx, mock.AnythingOfType("*domain.ServiceStatusHistory")).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := ci.TransitionStatus(ctx, "cs-1", "EN_PROCESO", "person-1", nil)

	assert.Error(t, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GetStatusHistory (HU73)
// ============================================

func TestGetStatusHistory_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1"}
	history := []domain.ServiceStatusHistory{
		{ID: "h-1", NewStatus: domain.ServiceStatusPending},
		{ID: "h-2", NewStatus: domain.ServiceStatusInProgress},
	}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("GetStatusHistory", ctx, "cs-1").Return(history, nil)

	result, err := ci.GetStatusHistory(ctx, "cs-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestGetStatusHistory_NotFound(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByID", ctx, "cs-bad").Return(nil, errors.New("not found"))

	result, err := ci.GetStatusHistory(ctx, "cs-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrCompletedServiceNotFound, err)
}

func TestGetStatusHistory_HistoryError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1"}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("GetStatusHistory", ctx, "cs-1").Return(nil, errors.New("history error"))

	result, err := ci.GetStatusHistory(ctx, "cs-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// UpdateCompletedServiceDetails (HU78)
// ============================================

func TestUpdateCompletedServiceDetails_Success(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)

	mockTx := new(mocks.MockTx)
	svc.On("BeginTx", ctx).Return(mockTx, nil)

	quoted := 150000.0
	final := 140000.0
	notes := "Todo revisado"
	svc.On("UpdateDetails", ctx, mockTx, "cs-1", &quoted, &final, &notes).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	err := ci.UpdateCompletedServiceDetails(ctx, "cs-1", &quoted, &final, &notes)

	assert.NoError(t, err)
	svc.AssertCalled(t, "UpdateDetails", ctx, mockTx, "cs-1", &quoted, &final, &notes)
	mockTx.AssertCalled(t, "Commit")
}

func TestUpdateCompletedServiceDetails_NotFound(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	svc.On("GetByID", ctx, "cs-missing").Return(nil, errors.New("not found"))

	quoted := 100000.0
	err := ci.UpdateCompletedServiceDetails(ctx, "cs-missing", &quoted, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceNotFound, err)
}

func TestUpdateCompletedServiceDetails_CannotUpdate_CompletedStatus(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)

	quoted := 100000.0
	err := ci.UpdateCompletedServiceDetails(ctx, "cs-1", &quoted, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotUpdate, err)
}

func TestUpdateCompletedServiceDetails_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusInProgress}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	quoted := 100000.0
	err := ci.UpdateCompletedServiceDetails(ctx, "cs-1", &quoted, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotUpdate, err)
}

func TestUpdateCompletedServiceDetails_UpdateDetailsError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ci, svc := setupCompletedServiceInteractor()

	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusPending}
	svc.On("GetByID", ctx, "cs-1").Return(cs, nil)

	mockTx := new(mocks.MockTx)
	svc.On("BeginTx", ctx).Return(mockTx, nil)

	quoted := 100000.0
	svc.On("UpdateDetails", ctx, mockTx, "cs-1", &quoted, mock.Anything, mock.Anything).Return(errors.New("update err"))
	mockTx.On("Rollback").Return(nil)

	err := ci.UpdateCompletedServiceDetails(ctx, "cs-1", &quoted, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCompletedServiceCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}
