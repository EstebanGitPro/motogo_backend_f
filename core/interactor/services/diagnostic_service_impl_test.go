package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// Mock repositories for DiagnosticService
// ============================================

type mockDiagnosticRepo struct{ mock.Mock }

func (m *mockDiagnosticRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *mockDiagnosticRepo) Save(ctx context.Context, tx output.Tx, d *domain.Diagnostic) error {
	return m.Called(ctx, tx, d).Error(0)
}
func (m *mockDiagnosticRepo) Update(ctx context.Context, tx output.Tx, d *domain.Diagnostic) error {
	return m.Called(ctx, tx, d).Error(0)
}
func (m *mockDiagnosticRepo) Delete(ctx context.Context, tx output.Tx, id string) error {
	return m.Called(ctx, tx, id).Error(0)
}
func (m *mockDiagnosticRepo) SaveEvidence(ctx context.Context, tx output.Tx, e *domain.DiagnosticEvidence) error {
	return m.Called(ctx, tx, e).Error(0)
}
func (m *mockDiagnosticRepo) GetByID(ctx context.Context, id string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}
func (m *mockDiagnosticRepo) GetByMotorcycleID(ctx context.Context, motoID string) ([]domain.Diagnostic, error) {
	args := m.Called(ctx, motoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Diagnostic), args.Error(1)
}
func (m *mockDiagnosticRepo) GetByMotorcycleAndBranch(ctx context.Context, motoID, branchID string) (*domain.Diagnostic, error) {
	args := m.Called(ctx, motoID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Diagnostic), args.Error(1)
}
func (m *mockDiagnosticRepo) GetEvidenceByDiagnosticID(ctx context.Context, diagID string) ([]domain.DiagnosticEvidence, error) {
	args := m.Called(ctx, diagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DiagnosticEvidence), args.Error(1)
}
func (m *mockDiagnosticRepo) DeleteEvidenceByDiagnosticID(ctx context.Context, tx output.Tx, diagID string) error {
	return m.Called(ctx, tx, diagID).Error(0)
}

type mockMotorcycleRepo struct{ mock.Mock }

func (m *mockMotorcycleRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) Save(_ context.Context, _ output.Tx, _ *domain.Motorcycle) error {
	return nil
}
func (m *mockMotorcycleRepo) Update(_ context.Context, _ output.Tx, _ *domain.Motorcycle) error {
	return nil
}
func (m *mockMotorcycleRepo) Delete(_ context.Context, _ output.Tx, _ string) error { return nil }
func (m *mockMotorcycleRepo) HardDelete(_ context.Context, _ output.Tx, _ string) error {
	return nil
}
func (m *mockMotorcycleRepo) ClearProfileImageURL(_ context.Context, _ output.Tx, _ string) error {
	return nil
}
func (m *mockMotorcycleRepo) GetByID(ctx context.Context, id string) (*domain.Motorcycle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Motorcycle), args.Error(1)
}
func (m *mockMotorcycleRepo) GetByOwnerID(_ context.Context, _ string) ([]domain.Motorcycle, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetByLicensePlate(_ context.Context, _ string) (*domain.Motorcycle, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetAllReferences(_ context.Context) ([]domain.MotorcycleReference, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetReferencesByBrandID(_ context.Context, _ string) ([]domain.MotorcycleReference, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetDistinctCategories(_ context.Context) ([]domain.MotorcycleCategory, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetLinesByCategory(_ context.Context, _ string) ([]domain.CategoryLine, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) GetDistinctDisplacements(_ context.Context) ([]domain.EngineDisplacementRange, error) {
	return nil, nil
}
func (m *mockMotorcycleRepo) ValidateReferenceExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockMotorcycleRepo) CheckLicensePlateExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockMotorcycleRepo) HasServiceHistory(_ context.Context, _ string) (bool, error) {
	return false, nil
}

type mockBranchRepo struct{ mock.Mock }

func (m *mockBranchRepo) BeginTx(_ context.Context) (output.Tx, error) { return nil, nil }
func (m *mockBranchRepo) SaveBranch(_ context.Context, _ output.Tx, _ domain.Branch) error {
	return nil
}
func (m *mockBranchRepo) UpdateBranch(_ context.Context, _ output.Tx, _ domain.Branch) error {
	return nil
}
func (m *mockBranchRepo) DeleteBranch(_ context.Context, _ output.Tx, _ string) error { return nil }
func (m *mockBranchRepo) SaveBranchBrands(_ context.Context, _ output.Tx, _ string, _ []string) error {
	return nil
}
func (m *mockBranchRepo) DeleteBranchBrands(_ context.Context, _ output.Tx, _ string) error {
	return nil
}
func (m *mockBranchRepo) SaveBranchDisplacementRanges(_ context.Context, _ output.Tx, _ string, _ []string) error {
	return nil
}
func (m *mockBranchRepo) DeleteBranchDisplacementRanges(_ context.Context, _ output.Tx, _ string) error {
	return nil
}
func (m *mockBranchRepo) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}
func (m *mockBranchRepo) GetBranchByFranchiseAndName(_ context.Context, _, _ string) (*domain.Branch, error) {
	return nil, nil
}
func (m *mockBranchRepo) GetBranchesByRepresentative(_ context.Context, _ string) ([]domain.Branch, error) {
	return nil, nil
}
func (m *mockBranchRepo) HasBranchesByRepresentative(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockBranchRepo) ValidateBrands(_ context.Context, _ []string) error { return nil }
func (m *mockBranchRepo) GetBranchesNearby(_ context.Context, _, _, _ float64, _ string, _, _, _, _ float64, _, _ string) ([]domain.NearbyBranch, error) {
	return nil, nil
}

type mockDiagTx struct{ mock.Mock }

func (m *mockDiagTx) Commit() error   { return m.Called().Error(0) }
func (m *mockDiagTx) Rollback() error { return m.Called().Error(0) }

// ============================================
// Helper
// ============================================

func setupDiagnosticService() (*mockDiagnosticRepo, *mockMotorcycleRepo, *mockBranchRepo, *services.DiagnosticServiceImpl) {
	diagRepo := new(mockDiagnosticRepo)
	motoRepo := new(mockMotorcycleRepo)
	branchRepo := new(mockBranchRepo)
	svc := services.NewDiagnosticService(diagRepo, motoRepo, branchRepo)
	return diagRepo, motoRepo, branchRepo, svc
}

// ============================================
// NewDiagnosticService Tests
// ============================================

func TestNewDiagnosticService_ReturnsInstance(t *testing.T) {
	_, _, _, svc := setupDiagnosticService()
	assert.NotNil(t, svc)
}

// ============================================
// BeginTx Tests
// ============================================

func TestDiagSvc_BeginTx_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	diagRepo.On("BeginTx", mock.Anything).Return(tx, nil)

	result, err := svc.BeginTx(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	diagRepo.AssertExpectations(t)
}

func TestDiagSvc_BeginTx_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.BeginTx(context.Background())

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// ValidateMotorcycleOwnership Tests
// ============================================

func TestDiagSvc_ValidateMotorcycleOwnership_Success(t *testing.T) {
	_, motoRepo, _, svc := setupDiagnosticService()

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	motoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.ID)
}

func TestDiagSvc_ValidateMotorcycleOwnership_NotFound(t *testing.T) {
	_, motoRepo, _, svc := setupDiagnosticService()

	motoRepo.On("GetByID", mock.Anything, "moto-bad").Return(nil, errors.New("not found"))

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-bad", "owner-1")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDiagSvc_ValidateMotorcycleOwnership_WrongOwner(t *testing.T) {
	_, motoRepo, _, svc := setupDiagnosticService()

	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-other"}
	motoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// ValidateBranchExists Tests
// ============================================

func TestDiagSvc_ValidateBranchExists_Success(t *testing.T) {
	_, _, branchRepo, svc := setupDiagnosticService()

	branchRepo.On("GetBranchByID", mock.Anything, "branch-1").Return(&domain.Branch{ID: "branch-1"}, nil)

	err := svc.ValidateBranchExists(context.Background(), "branch-1")

	assert.NoError(t, err)
}

func TestDiagSvc_ValidateBranchExists_NotFound(t *testing.T) {
	_, _, branchRepo, svc := setupDiagnosticService()

	branchRepo.On("GetBranchByID", mock.Anything, "branch-bad").Return(nil, errors.New("not found"))

	err := svc.ValidateBranchExists(context.Background(), "branch-bad")

	assert.Equal(t, domain.ErrBranchNotFound, err)
}

// ============================================
// RegisterOrUpdateDiagnostic Tests
// ============================================

func TestDiagSvc_RegisterOrUpdate_CreateNew(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	desc := "Motor con ruido"

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	diagRepo.On("SaveEvidence", mock.Anything, tx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", &desc, []string{"http://img1.jpg"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	assert.Equal(t, "branch-1", result.BranchID)
	assert.Len(t, result.Evidence, 1)
}

func TestDiagSvc_RegisterOrUpdate_CreateNew_NoEvidence(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, []string{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Evidence)
}

func TestDiagSvc_RegisterOrUpdate_UpdateExisting(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	existing := &domain.Diagnostic{ID: "diag-existing", MotorcycleID: "moto-1", BranchID: "branch-1"}
	desc := "Descripción actualizada"

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	diagRepo.On("Update", mock.Anything, tx, existing).Return(nil)
	diagRepo.On("DeleteEvidenceByDiagnosticID", mock.Anything, tx, "diag-existing").Return(nil)
	diagRepo.On("SaveEvidence", mock.Anything, tx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(nil)

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", &desc, []string{"http://new-img.jpg"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "diag-existing", result.ID)
	assert.Len(t, result.Evidence, 1)
}

func TestDiagSvc_RegisterOrUpdate_LookupError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, errors.New("db error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagSvc_RegisterOrUpdate_SaveError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.Diagnostic")).Return(errors.New("save error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagSvc_RegisterOrUpdate_SaveEvidenceError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(nil, nil)
	diagRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.Diagnostic")).Return(nil)
	diagRepo.On("SaveEvidence", mock.Anything, tx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(errors.New("evidence error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, []string{"http://img.jpg"})

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagSvc_RegisterOrUpdate_UpdateError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	diagRepo.On("Update", mock.Anything, tx, existing).Return(errors.New("update error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagSvc_RegisterOrUpdate_DeleteEvidenceError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	diagRepo.On("Update", mock.Anything, tx, existing).Return(nil)
	diagRepo.On("DeleteEvidenceByDiagnosticID", mock.Anything, tx, "diag-1").Return(errors.New("delete evid error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

func TestDiagSvc_RegisterOrUpdate_UpsertSaveEvidenceError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	existing := &domain.Diagnostic{ID: "diag-1", MotorcycleID: "moto-1", BranchID: "branch-1"}

	diagRepo.On("GetByMotorcycleAndBranch", mock.Anything, "moto-1", "branch-1").Return(existing, nil)
	diagRepo.On("Update", mock.Anything, tx, existing).Return(nil)
	diagRepo.On("DeleteEvidenceByDiagnosticID", mock.Anything, tx, "diag-1").Return(nil)
	diagRepo.On("SaveEvidence", mock.Anything, tx, mock.AnythingOfType("*domain.DiagnosticEvidence")).Return(errors.New("save evid error"))

	result, err := svc.RegisterOrUpdateDiagnostic(context.Background(), tx, "moto-1", "branch-1", nil, []string{"http://img.jpg"})

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDiagnosticCannotSave, err)
}

// ============================================
// GetByID Tests
// ============================================

func TestDiagSvc_GetByID_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	expected := &domain.Diagnostic{ID: "diag-1"}
	diagRepo.On("GetByID", mock.Anything, "diag-1").Return(expected, nil)

	result, err := svc.GetByID(context.Background(), "diag-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDiagSvc_GetByID_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagRepo.On("GetByID", mock.Anything, "bad").Return(nil, errors.New("not found"))

	result, err := svc.GetByID(context.Background(), "bad")

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// GetByMotorcycleID Tests
// ============================================

func TestDiagSvc_GetByMotorcycleID_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	expected := []domain.Diagnostic{{ID: "d1"}, {ID: "d2"}}
	diagRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(expected, nil)

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDiagSvc_GetByMotorcycleID_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagRepo.On("GetByMotorcycleID", mock.Anything, "bad").Return(nil, errors.New("error"))

	result, err := svc.GetByMotorcycleID(context.Background(), "bad")

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// ApplyDiagnosticUpdates Tests
// ============================================

func TestDiagSvc_ApplyUpdates_UpdatesBoth(t *testing.T) {
	_, _, _, svc := setupDiagnosticService()

	existing := &domain.Diagnostic{ID: "d1"}
	newDesc := "new description"
	newSol := "new solution"
	updates := &domain.Diagnostic{ProblemDescription: &newDesc, PossibleSolution: &newSol}

	svc.ApplyDiagnosticUpdates(existing, updates)

	assert.NotNil(t, existing.ProblemDescription)
	assert.Equal(t, "new description", *existing.ProblemDescription)
	assert.NotNil(t, existing.PossibleSolution)
	assert.Equal(t, "new solution", *existing.PossibleSolution)
}

func TestDiagSvc_ApplyUpdates_PartialUpdate(t *testing.T) {
	_, _, _, svc := setupDiagnosticService()

	oldDesc := "old desc"
	existing := &domain.Diagnostic{ID: "d1", ProblemDescription: &oldDesc}
	newSol := "solution only"
	updates := &domain.Diagnostic{PossibleSolution: &newSol}

	svc.ApplyDiagnosticUpdates(existing, updates)

	// ProblemDescription should remain unchanged
	assert.Equal(t, "old desc", *existing.ProblemDescription)
	assert.NotNil(t, existing.PossibleSolution)
	assert.Equal(t, "solution only", *existing.PossibleSolution)
}

func TestDiagSvc_ApplyUpdates_NoUpdates(t *testing.T) {
	_, _, _, svc := setupDiagnosticService()

	existing := &domain.Diagnostic{ID: "d1"}
	updates := &domain.Diagnostic{}

	svc.ApplyDiagnosticUpdates(existing, updates)

	assert.Nil(t, existing.ProblemDescription)
	assert.Nil(t, existing.PossibleSolution)
}

// ============================================
// UpdateDiagnostic Tests
// ============================================

func TestDiagSvc_UpdateDiagnostic_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	diag := &domain.Diagnostic{ID: "diag-1"}

	diagRepo.On("Update", mock.Anything, tx, diag).Return(nil)

	err := svc.UpdateDiagnostic(context.Background(), tx, diag)

	assert.NoError(t, err)
}

func TestDiagSvc_UpdateDiagnostic_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	diag := &domain.Diagnostic{ID: "diag-bad"}

	diagRepo.On("Update", mock.Anything, tx, diag).Return(errors.New("update error"))

	err := svc.UpdateDiagnostic(context.Background(), tx, diag)

	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

// ============================================
// DeleteDiagnostic Tests
// ============================================

func TestDiagSvc_DeleteDiagnostic_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("Delete", mock.Anything, tx, "diag-1").Return(nil)

	err := svc.DeleteDiagnostic(context.Background(), tx, "diag-1")

	assert.NoError(t, err)
}

func TestDiagSvc_DeleteDiagnostic_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("Delete", mock.Anything, tx, "diag-bad").Return(errors.New("delete error"))

	err := svc.DeleteDiagnostic(context.Background(), tx, "diag-bad")

	assert.Equal(t, domain.ErrDiagnosticCannotDelete, err)
}

// ============================================
// SetSolution Tests
// ============================================

func TestDiagSvc_SetSolution_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	diag := &domain.Diagnostic{ID: "diag-1"}

	diagRepo.On("GetByID", mock.Anything, "diag-1").Return(diag, nil)
	diagRepo.On("Update", mock.Anything, tx, diag).Return(nil)

	err := svc.SetSolution(context.Background(), tx, "diag-1", "Cambiar aceite")

	assert.NoError(t, err)
	assert.NotNil(t, diag.PossibleSolution)
	assert.Equal(t, "Cambiar aceite", *diag.PossibleSolution)
}

func TestDiagSvc_SetSolution_DiagNotFound(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)

	diagRepo.On("GetByID", mock.Anything, "diag-bad").Return(nil, errors.New("not found"))

	err := svc.SetSolution(context.Background(), tx, "diag-bad", "solution")

	assert.Equal(t, domain.ErrDiagnosticNotFound, err)
}

func TestDiagSvc_SetSolution_UpdateError(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	tx := new(mockDiagTx)
	diag := &domain.Diagnostic{ID: "diag-1"}

	diagRepo.On("GetByID", mock.Anything, "diag-1").Return(diag, nil)
	diagRepo.On("Update", mock.Anything, tx, diag).Return(errors.New("update error"))

	err := svc.SetSolution(context.Background(), tx, "diag-1", "solution")

	assert.Equal(t, domain.ErrDiagnosticCannotUpdate, err)
}

// ============================================
// LoadEvidence Tests
// ============================================

func TestDiagSvc_LoadEvidence_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	expected := []domain.DiagnosticEvidence{{ID: "ev-1"}, {ID: "ev-2"}}
	diagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "diag-1").Return(expected, nil)

	result, err := svc.LoadEvidence(context.Background(), "diag-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDiagSvc_LoadEvidence_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "bad").Return(nil, errors.New("error"))

	result, err := svc.LoadEvidence(context.Background(), "bad")

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// LoadEvidenceForDiagnostics Tests
// ============================================

func TestDiagSvc_LoadEvidenceForDiagnostics_Success(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagnostics := []domain.Diagnostic{
		{ID: "d1"},
		{ID: "d2"},
	}

	diagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "d1").Return([]domain.DiagnosticEvidence{{ID: "ev-1"}}, nil)
	diagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "d2").Return([]domain.DiagnosticEvidence{{ID: "ev-2"}}, nil)

	err := svc.LoadEvidenceForDiagnostics(context.Background(), diagnostics)

	assert.NoError(t, err)
	assert.Len(t, diagnostics[0].Evidence, 1)
	assert.Len(t, diagnostics[1].Evidence, 1)
}

func TestDiagSvc_LoadEvidenceForDiagnostics_Error(t *testing.T) {
	diagRepo, _, _, svc := setupDiagnosticService()

	diagnostics := []domain.Diagnostic{
		{ID: "d1"},
	}

	diagRepo.On("GetEvidenceByDiagnosticID", mock.Anything, "d1").Return(nil, errors.New("load error"))

	err := svc.LoadEvidenceForDiagnostics(context.Background(), diagnostics)

	assert.Error(t, err)
}

func TestDiagSvc_LoadEvidenceForDiagnostics_Empty(t *testing.T) {
	_, _, _, svc := setupDiagnosticService()

	diagnostics := []domain.Diagnostic{}

	err := svc.LoadEvidenceForDiagnostics(context.Background(), diagnostics)

	assert.NoError(t, err)
}
