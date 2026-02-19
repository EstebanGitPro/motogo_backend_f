package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Mock repositories (local to this test file)
// ============================================

type mockCSRepo struct{ mock.Mock }

func (m *mockCSRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *mockCSRepo) Save(ctx context.Context, tx output.Tx, service *domain.CompletedService) error {
	return m.Called(ctx, tx, service).Error(0)
}
func (m *mockCSRepo) SaveItems(ctx context.Context, tx output.Tx, items []domain.CompletedServiceItem) error {
	return m.Called(ctx, tx, items).Error(0)
}
func (m *mockCSRepo) SaveStatusHistory(ctx context.Context, tx output.Tx, history *domain.ServiceStatusHistory) error {
	return m.Called(ctx, tx, history).Error(0)
}
func (m *mockCSRepo) GetByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedService), args.Error(1)
}
func (m *mockCSRepo) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error) {
	args := m.Called(ctx, motorcycleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CompletedService), args.Error(1)
}
func (m *mockCSRepo) GetByBranchID(ctx context.Context, branchID string) ([]domain.CompletedService, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CompletedService), args.Error(1)
}
func (m *mockCSRepo) GetItemsByCompletedServiceID(ctx context.Context, csID string) ([]domain.CompletedServiceItem, error) {
	args := m.Called(ctx, csID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CompletedServiceItem), args.Error(1)
}
func (m *mockCSRepo) ValidateBranchServices(ctx context.Context, branchID string, serviceIDs []string) error {
	return m.Called(ctx, branchID, serviceIDs).Error(0)
}
func (m *mockCSRepo) HasActiveService(ctx context.Context, branchID, motorcycleID string) (bool, error) {
	args := m.Called(ctx, branchID, motorcycleID)
	return args.Bool(0), args.Error(1)
}
func (m *mockCSRepo) Delete(ctx context.Context, tx output.Tx, serviceID string) error {
	return m.Called(ctx, tx, serviceID).Error(0)
}
func (m *mockCSRepo) SoftDelete(ctx context.Context, tx output.Tx, serviceID string) error {
	return m.Called(ctx, tx, serviceID).Error(0)
}
func (m *mockCSRepo) UpdateStatus(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time) error {
	return m.Called(ctx, tx, serviceID, status, completionDate).Error(0)
}
func (m *mockCSRepo) GetStatusHistory(ctx context.Context, serviceID string) ([]domain.ServiceStatusHistory, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ServiceStatusHistory), args.Error(1)
}
func (m *mockCSRepo) UpdateStatusWithPrice(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time, finalPrice *float64) error {
	return m.Called(ctx, tx, serviceID, status, completionDate, finalPrice).Error(0)
}
func (m *mockCSRepo) UpdateDetails(ctx context.Context, tx output.Tx, serviceID string, quotedPrice, finalPrice *float64, notes *string) error {
	return m.Called(ctx, tx, serviceID, quotedPrice, finalPrice, notes).Error(0)
}
func (m *mockCSRepo) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string, isOffensive bool) error {
	return m.Called(ctx, tx, itemID, rating, comment, isOffensive).Error(0)
}
func (m *mockCSRepo) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedServiceItem), args.Error(1)
}

type mockDiagRepo struct{ mock.Mock }

func (m *mockDiagRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *mockDiagRepo) Save(ctx context.Context, tx output.Tx, d *domain.Diagnostic) error {
	return m.Called(ctx, tx, d).Error(0)
}
func (m *mockDiagRepo) Update(ctx context.Context, tx output.Tx, d *domain.Diagnostic) error {
	return m.Called(ctx, tx, d).Error(0)
}
func (m *mockDiagRepo) Delete(ctx context.Context, tx output.Tx, id string) error {
	return m.Called(ctx, tx, id).Error(0)
}
func (m *mockDiagRepo) SaveEvidence(ctx context.Context, tx output.Tx, e *domain.DiagnosticEvidence) error {
	return m.Called(ctx, tx, e).Error(0)
}
func (m *mockDiagRepo) GetByID(ctx context.Context, id string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}
func (m *mockDiagRepo) GetByMotorcycleID(ctx context.Context, motoID string) ([]domain.Diagnostic, error) {
	args := m.Called(ctx, motoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Diagnostic), args.Error(1)
}
func (m *mockDiagRepo) GetByMotorcycleAndBranch(ctx context.Context, motoID, branchID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, motoID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}
func (m *mockDiagRepo) GetEvidenceByDiagnosticID(ctx context.Context, diagID string) ([]domain.DiagnosticEvidence, error) {
	args := m.Called(ctx, diagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticEvidence), args.Error(1)
}
func (m *mockDiagRepo) DeleteEvidenceByDiagnosticID(ctx context.Context, tx output.Tx, diagID string) error {
	return m.Called(ctx, tx, diagID).Error(0)
}

// ============================================
// Helper
// ============================================

func setupCSService() (*mockCSRepo, *mockDiagRepo, input.CompletedServiceService) {
	csRepo := new(mockCSRepo)
	diagRepo := new(mockDiagRepo)
	svc := services.NewCompletedServiceService(csRepo, diagRepo)
	return csRepo, diagRepo, svc
}

// ============================================
// ValidateBranchServices
// ============================================

func TestCSValidateBranchServices_EmptyIDs(t *testing.T) {
	_, _, svc := setupCSService()

	err := svc.ValidateBranchServices(context.Background(), "branch-1", []string{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one service_id is required")
}

func TestCSValidateBranchServices_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("ValidateBranchServices", mock.Anything, "branch-1", []string{"svc-1"}).Return(nil)

	err := svc.ValidateBranchServices(context.Background(), "branch-1", []string{"svc-1"})

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

func TestCSValidateBranchServices_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("ValidateBranchServices", mock.Anything, "branch-1", []string{"svc-bad"}).Return(errors.New("invalid services"))

	err := svc.ValidateBranchServices(context.Background(), "branch-1", []string{"svc-bad"})

	assert.Error(t, err)
}

// ============================================
// ValidateDiagnosticForMotorcycle
// ============================================

func TestCSValidateDiagnostic_NotFound(t *testing.T) {
	_, diagRepo, svc := setupCSService()

	diagRepo.On("GetByID", mock.Anything, "diag-bad").Return(nil, errors.New("not found"))

	err := svc.ValidateDiagnosticForMotorcycle(context.Background(), "diag-bad", "moto-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "diagnostic not found")
}

func TestCSValidateDiagnostic_WrongMotorcycle(t *testing.T) {
	_, diagRepo, svc := setupCSService()

	diagRepo.On("GetByID", mock.Anything, "diag-1").Return(
		&domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-OTHER"}, nil,
	)

	err := svc.ValidateDiagnosticForMotorcycle(context.Background(), "diag-1", "moto-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "diagnostic does not belong to this motorcycle")
}

func TestCSValidateDiagnostic_Success(t *testing.T) {
	_, diagRepo, svc := setupCSService()

	diagRepo.On("GetByID", mock.Anything, "diag-1").Return(
		&domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1"}, nil,
	)

	err := svc.ValidateDiagnosticForMotorcycle(context.Background(), "diag-1", "moto-1")

	assert.NoError(t, err)
}

// ============================================
// ValidateNoActiveService
// ============================================

func TestCSValidateNoActive_HasActive(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("HasActiveService", mock.Anything, "branch-1", "moto-1").Return(true, nil)

	err := svc.ValidateNoActiveService(context.Background(), "branch-1", "moto-1")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrActiveServiceExists, err)
}

func TestCSValidateNoActive_NoActive(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("HasActiveService", mock.Anything, "branch-1", "moto-1").Return(false, nil)

	err := svc.ValidateNoActiveService(context.Background(), "branch-1", "moto-1")

	assert.NoError(t, err)
}

func TestCSValidateNoActive_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("HasActiveService", mock.Anything, "branch-1", "moto-1").Return(false, errors.New("db error"))

	err := svc.ValidateNoActiveService(context.Background(), "branch-1", "moto-1")

	assert.Error(t, err)
}

// ============================================
// GetByID (with item hydration)
// ============================================

func TestCSGetByID_Success(t *testing.T) {
	csRepo, _, svc := setupCSService()

	cs := &domain.CompletedService{ID: "cs-1", BranchID: "branch-1"}
	items := []domain.CompletedServiceItem{{ID: "item-1", CompletedServiceID: "cs-1"}}

	csRepo.On("GetByID", mock.Anything, "cs-1").Return(cs, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(items, nil)

	result, err := svc.GetByID(context.Background(), "cs-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 1)
}

func TestCSGetByID_NotFound(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetByID", mock.Anything, "cs-bad").Return(nil, errors.New("not found"))

	result, err := svc.GetByID(context.Background(), "cs-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCSGetByID_ItemsError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	cs := &domain.CompletedService{ID: "cs-1"}
	csRepo.On("GetByID", mock.Anything, "cs-1").Return(cs, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(nil, errors.New("items error"))

	result, err := svc.GetByID(context.Background(), "cs-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// GetByMotorcycleID — hydration
// ============================================

func TestCSGetByMotorcycleID_HydratesItems(t *testing.T) {
	csRepo, _, svc := setupCSService()

	services := []domain.CompletedService{
		{ID: "cs-1"},
		{ID: "cs-2"},
	}
	csRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(services, nil)

	items1 := []domain.CompletedServiceItem{
		{ID: "item-1", ServiceID: "svc-1"},
		{ID: "item-2", ServiceID: "svc-2"},
	}
	items2 := []domain.CompletedServiceItem{
		{ID: "item-3", ServiceID: "svc-3"},
	}
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(items1, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-2").Return(items2, nil)

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].Services, 2)
	assert.Len(t, result[1].Services, 1)
	assert.Equal(t, "item-1", result[0].Services[0].ID)
	assert.Equal(t, "item-3", result[1].Services[0].ID)
	csRepo.AssertExpectations(t)
}

func TestCSGetByMotorcycleID_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetByMotorcycleID", mock.Anything, "moto-bad").Return(nil, errors.New("db error"))

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCSGetByMotorcycleID_HydrationError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	services := []domain.CompletedService{{ID: "cs-1"}}
	csRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(services, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(nil, errors.New("item load error"))

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "item load error")
}

func TestCSGetByMotorcycleID_EmptyList(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetByMotorcycleID", mock.Anything, "moto-empty").Return([]domain.CompletedService{}, nil)

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-empty")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// GetByBranchID — hydration
// ============================================

func TestCSGetByBranchID_HydratesItems(t *testing.T) {
	csRepo, _, svc := setupCSService()

	services := []domain.CompletedService{
		{ID: "cs-1"},
		{ID: "cs-2"},
	}
	csRepo.On("GetByBranchID", mock.Anything, "branch-1").Return(services, nil)

	items1 := []domain.CompletedServiceItem{
		{ID: "item-a", ServiceID: "svc-a"},
	}
	items2 := []domain.CompletedServiceItem{
		{ID: "item-b", ServiceID: "svc-b"},
		{ID: "item-c", ServiceID: "svc-c"},
	}
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(items1, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-2").Return(items2, nil)

	result, err := svc.GetByBranchID(context.Background(), "branch-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].Services, 1)
	assert.Len(t, result[1].Services, 2)
	csRepo.AssertExpectations(t)
}

func TestCSGetByBranchID_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetByBranchID", mock.Anything, "branch-bad").Return(nil, errors.New("db error"))

	result, err := svc.GetByBranchID(context.Background(), "branch-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCSGetByBranchID_HydrationError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	services := []domain.CompletedService{{ID: "cs-1"}, {ID: "cs-2"}}
	csRepo.On("GetByBranchID", mock.Anything, "branch-1").Return(services, nil)
	csRepo.On("GetItemsByCompletedServiceID", mock.Anything, "cs-1").Return(nil, errors.New("hydration fail"))

	result, err := svc.GetByBranchID(context.Background(), "branch-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCSGetByBranchID_EmptyList(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetByBranchID", mock.Anything, "branch-empty").Return([]domain.CompletedService{}, nil)

	result, err := svc.GetByBranchID(context.Background(), "branch-empty")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCSBeginTx_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	tx, err := svc.BeginTx(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

// mockCSTx for BeginTx test
type mockCSTx struct{ mock.Mock }

func (m *mockCSTx) Commit() error   { return m.Called().Error(0) }
func (m *mockCSTx) Rollback() error { return m.Called().Error(0) }

// ============================================
// SaveCompletedService
// ============================================

func TestCSSaveCompletedService_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	cs := &domain.CompletedService{ID: "cs-1"}
	csRepo.On("Save", mock.Anything, mockTx, cs).Return(nil)

	err := svc.SaveCompletedService(context.Background(), mockTx, cs)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

// ============================================
// SaveItems
// ============================================

func TestCSSaveItems_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	items := []domain.CompletedServiceItem{{ID: "item-1"}}
	csRepo.On("SaveItems", mock.Anything, mockTx, items).Return(nil)

	err := svc.SaveItems(context.Background(), mockTx, items)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

// ============================================
// SaveStatusHistory
// ============================================

func TestCSSaveStatusHistory_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	history := &domain.ServiceStatusHistory{ID: "hist-1"}
	csRepo.On("SaveStatusHistory", mock.Anything, mockTx, history).Return(nil)

	err := svc.SaveStatusHistory(context.Background(), mockTx, history)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

// ============================================
// DeleteCompletedService (HU65)
// ============================================

func TestCSDeleteCompletedService_HardDelete_Pendiente(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("Delete", mock.Anything, mockTx, "cs-1").Return(nil)

	err := svc.DeleteCompletedService(context.Background(), mockTx, "cs-1", domain.ServiceStatusPending)

	assert.NoError(t, err)
	csRepo.AssertCalled(t, "Delete", mock.Anything, mockTx, "cs-1")
	csRepo.AssertNotCalled(t, "SoftDelete")
}

func TestCSDeleteCompletedService_HardDelete_EnProceso(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("Delete", mock.Anything, mockTx, "cs-1").Return(nil)

	err := svc.DeleteCompletedService(context.Background(), mockTx, "cs-1", domain.ServiceStatusInProgress)

	assert.NoError(t, err)
	csRepo.AssertCalled(t, "Delete", mock.Anything, mockTx, "cs-1")
	csRepo.AssertNotCalled(t, "SoftDelete")
}

func TestCSDeleteCompletedService_SoftDelete_Finalizado(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("SoftDelete", mock.Anything, mockTx, "cs-1").Return(nil)

	err := svc.DeleteCompletedService(context.Background(), mockTx, "cs-1", domain.ServiceStatusCompleted)

	assert.NoError(t, err)
	csRepo.AssertCalled(t, "SoftDelete", mock.Anything, mockTx, "cs-1")
	csRepo.AssertNotCalled(t, "Delete")
}

func TestCSDeleteCompletedService_SoftDelete_Cancelado(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("SoftDelete", mock.Anything, mockTx, "cs-1").Return(nil)

	err := svc.DeleteCompletedService(context.Background(), mockTx, "cs-1", domain.ServiceStatusCancelled)

	assert.NoError(t, err)
	csRepo.AssertCalled(t, "SoftDelete", mock.Anything, mockTx, "cs-1")
	csRepo.AssertNotCalled(t, "Delete")
}

func TestCSDeleteCompletedService_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("Delete", mock.Anything, mockTx, "cs-bad").Return(errors.New("delete error"))

	err := svc.DeleteCompletedService(context.Background(), mockTx, "cs-bad", domain.ServiceStatusPending)

	assert.Error(t, err)
}

// ============================================
// UpdateStatus (HU74)
// ============================================

func TestCSUpdateStatus_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("UpdateStatus", mock.Anything, mockTx, "cs-1", "EN_PROCESO", (*time.Time)(nil)).Return(nil)

	err := svc.UpdateStatus(context.Background(), mockTx, "cs-1", "EN_PROCESO", nil)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

func TestCSUpdateStatus_WithCompletionDate(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	now := time.Now()
	csRepo.On("UpdateStatus", mock.Anything, mockTx, "cs-1", "FINALIZADO", &now).Return(nil)

	err := svc.UpdateStatus(context.Background(), mockTx, "cs-1", "FINALIZADO", &now)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

func TestCSUpdateStatus_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("UpdateStatus", mock.Anything, mockTx, "cs-1", "INVALID", (*time.Time)(nil)).Return(errors.New("update error"))

	err := svc.UpdateStatus(context.Background(), mockTx, "cs-1", "INVALID", nil)

	assert.Error(t, err)
}

// ============================================
// GetStatusHistory (HU73)
// ============================================

func TestCSGetStatusHistory_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	expected := []domain.ServiceStatusHistory{
		{ID: "hist-1", NewStatus: domain.ServiceStatusPending},
		{ID: "hist-2", NewStatus: domain.ServiceStatusInProgress},
	}
	csRepo.On("GetStatusHistory", mock.Anything, "cs-1").Return(expected, nil)

	result, err := svc.GetStatusHistory(context.Background(), "cs-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	csRepo.AssertExpectations(t)
}

func TestCSGetStatusHistory_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetStatusHistory", mock.Anything, "cs-bad").Return(nil, errors.New("history error"))

	result, err := svc.GetStatusHistory(context.Background(), "cs-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCSGetStatusHistory_Empty(t *testing.T) {
	csRepo, _, svc := setupCSService()

	csRepo.On("GetStatusHistory", mock.Anything, "cs-new").Return([]domain.ServiceStatusHistory{}, nil)

	result, err := svc.GetStatusHistory(context.Background(), "cs-new")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// UpdateStatusWithPrice
// ============================================

func TestCSUpdateStatusWithPrice_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	completionDate := time.Now()
	finalPrice := 150000.0

	csRepo.On("UpdateStatusWithPrice", mock.Anything, mockTx, "cs-1", "FINALIZADO", &completionDate, &finalPrice).Return(nil)

	err := svc.UpdateStatusWithPrice(context.Background(), mockTx, "cs-1", "FINALIZADO", &completionDate, &finalPrice)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

func TestCSUpdateStatusWithPrice_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	csRepo.On("UpdateStatusWithPrice", mock.Anything, mockTx, "cs-1", "FINALIZADO", (*time.Time)(nil), (*float64)(nil)).Return(errors.New("update error"))

	err := svc.UpdateStatusWithPrice(context.Background(), mockTx, "cs-1", "FINALIZADO", nil, nil)

	assert.Error(t, err)
}

// ============================================
// UpdateDetails
// ============================================

func TestCSUpdateDetails_DelegatesToRepo(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	quotedPrice := 200000.0
	finalPrice := 180000.0
	notes := "Revisión completada"

	csRepo.On("UpdateDetails", mock.Anything, mockTx, "cs-1", &quotedPrice, &finalPrice, &notes).Return(nil)

	err := svc.UpdateDetails(context.Background(), mockTx, "cs-1", &quotedPrice, &finalPrice, &notes)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}

func TestCSUpdateDetails_RepoError(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	quotedPrice := 100000.0
	csRepo.On("UpdateDetails", mock.Anything, mockTx, "cs-bad", &quotedPrice, (*float64)(nil), (*string)(nil)).Return(errors.New("update error"))

	err := svc.UpdateDetails(context.Background(), mockTx, "cs-bad", &quotedPrice, nil, nil)

	assert.Error(t, err)
}

func TestCSUpdateDetails_NilOptionals(t *testing.T) {
	csRepo, _, svc := setupCSService()

	mockTx := new(mockCSTx)
	notes := "Solo notas"
	csRepo.On("UpdateDetails", mock.Anything, mockTx, "cs-1", (*float64)(nil), (*float64)(nil), &notes).Return(nil)

	err := svc.UpdateDetails(context.Background(), mockTx, "cs-1", nil, nil, &notes)

	assert.NoError(t, err)
	csRepo.AssertExpectations(t)
}
