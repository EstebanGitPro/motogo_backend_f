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
// Mock repositories for EvidenceService
// ============================================

type mockEvidenceRepo struct{ mock.Mock }

func (m *mockEvidenceRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *mockEvidenceRepo) Save(ctx context.Context, tx output.Tx, e *domain.MotorcycleEvidence) error {
	return m.Called(ctx, tx, e).Error(0)
}
func (m *mockEvidenceRepo) Update(ctx context.Context, tx output.Tx, e *domain.MotorcycleEvidence) error {
	return m.Called(ctx, tx, e).Error(0)
}
func (m *mockEvidenceRepo) Delete(ctx context.Context, tx output.Tx, id string) error {
	return m.Called(ctx, tx, id).Error(0)
}
func (m *mockEvidenceRepo) GetByID(ctx context.Context, id string) (*domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MotorcycleEvidence), args.Error(1)
}
func (m *mockEvidenceRepo) GetByMotorcycleID(ctx context.Context, motoID string) ([]domain.MotorcycleEvidence, error) {
	args := m.Called(ctx, motoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MotorcycleEvidence), args.Error(1)
}
func (m *mockEvidenceRepo) CountByMotorcycleID(ctx context.Context, motoID string) (int, error) {
	args := m.Called(ctx, motoID)
	return args.Int(0), args.Error(1)
}

type mockStorageClient struct{ mock.Mock }

func (m *mockStorageClient) DeleteStorageFile(ctx context.Context, fileURL string) error {
	return m.Called(ctx, fileURL).Error(0)
}

type mockEvidTx struct{ mock.Mock }

func (m *mockEvidTx) Commit() error   { return m.Called().Error(0) }
func (m *mockEvidTx) Rollback() error { return m.Called().Error(0) }

// ============================================
// Helper
// ============================================

func setupEvidenceService() (*mockEvidenceRepo, *mockMotorcycleRepo, *services.EvidenceServiceImpl) {
	evidRepo := new(mockEvidenceRepo)
	motoRepo := new(mockMotorcycleRepo)
	svc := services.NewEvidenceService(evidRepo, motoRepo)
	return evidRepo, motoRepo, svc
}

// ============================================
// NewEvidenceService Tests
// ============================================

func TestNewEvidenceService_ReturnsInstance(t *testing.T) {
	_, _, svc := setupEvidenceService()
	assert.NotNil(t, svc)
}

// ============================================
// WithStorageClient Tests
// ============================================

func TestEvidSvc_WithStorageClient(t *testing.T) {
	_, _, svc := setupEvidenceService()
	client := new(mockStorageClient)
	result := svc.WithStorageClient(client)
	assert.NotNil(t, result)
	assert.Equal(t, svc, result)
}

// ============================================
// BeginTx Tests
// ============================================

func TestEvidSvc_BeginTx_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)
	evidRepo.On("BeginTx", mock.Anything).Return(tx, nil)

	result, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestEvidSvc_BeginTx_Error(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	evidRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.BeginTx(context.Background())
	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// ValidateMotorcycleOwnership Tests
// ============================================

func TestEvidSvc_ValidateOwnership_Success(t *testing.T) {
	_, motoRepo, svc := setupEvidenceService()
	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-1"}
	motoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")
	assert.NoError(t, err)
	assert.Equal(t, "moto-1", result.ID)
}

func TestEvidSvc_ValidateOwnership_NotFound(t *testing.T) {
	_, motoRepo, svc := setupEvidenceService()
	motoRepo.On("GetByID", mock.Anything, "moto-bad").Return(nil, errors.New("not found"))

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-bad", "owner-1")
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestEvidSvc_ValidateOwnership_WrongOwner(t *testing.T) {
	_, motoRepo, svc := setupEvidenceService()
	moto := &domain.Motorcycle{ID: "moto-1", OwnerID: "owner-other"}
	motoRepo.On("GetByID", mock.Anything, "moto-1").Return(moto, nil)

	result, err := svc.ValidateMotorcycleOwnership(context.Background(), "moto-1", "owner-1")
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// CheckEvidenceLimit Tests
// ============================================

func TestEvidSvc_CheckLimit_OK(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	evidRepo.On("CountByMotorcycleID", mock.Anything, "moto-1").Return(2, nil)

	err := svc.CheckEvidenceLimit(context.Background(), "moto-1")
	assert.NoError(t, err)
}

func TestEvidSvc_CheckLimit_Exceeded(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	evidRepo.On("CountByMotorcycleID", mock.Anything, "moto-1").Return(services.MaxEvidencePerMotorcycle, nil)

	err := svc.CheckEvidenceLimit(context.Background(), "moto-1")
	assert.Equal(t, domain.ErrEvidenceLimitExceeded, err)
}

func TestEvidSvc_CheckLimit_CountError(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	evidRepo.On("CountByMotorcycleID", mock.Anything, "moto-1").Return(0, errors.New("db error"))

	err := svc.CheckEvidenceLimit(context.Background(), "moto-1")
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

// ============================================
// CreateEvidence Tests
// ============================================

func TestEvidSvc_CreateEvidence_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)

	evidence := &domain.MotorcycleEvidence{
		ImageURL: "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
	}

	evidRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.MotorcycleEvidence")).Return(nil)

	result, err := svc.CreateEvidence(context.Background(), tx, "moto-1", evidence)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "moto-1", result.MotorcycleID)
}

func TestEvidSvc_CreateEvidence_SaveError(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)

	evidence := &domain.MotorcycleEvidence{
		ImageURL: "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
	}

	evidRepo.On("Save", mock.Anything, tx, mock.AnythingOfType("*domain.MotorcycleEvidence")).Return(errors.New("save error"))

	result, err := svc.CreateEvidence(context.Background(), tx, "moto-1", evidence)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

// ============================================
// GetByID Tests
// ============================================

func TestEvidSvc_GetByID_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	expected := &domain.MotorcycleEvidence{ID: "ev-1"}
	evidRepo.On("GetByID", mock.Anything, "ev-1").Return(expected, nil)

	result, err := svc.GetByID(context.Background(), "ev-1")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestEvidSvc_GetByID_Error(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	evidRepo.On("GetByID", mock.Anything, "bad").Return(nil, errors.New("not found"))

	result, err := svc.GetByID(context.Background(), "bad")
	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================
// GetByMotorcycleID Tests
// ============================================

func TestEvidSvc_GetByMotorcycleID_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	expected := []domain.MotorcycleEvidence{{ID: "ev-1"}, {ID: "ev-2"}}
	evidRepo.On("GetByMotorcycleID", mock.Anything, "moto-1").Return(expected, nil)

	result, err := svc.GetByMotorcycleID(context.Background(), "moto-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============================================
// ApplyUpdatesAndCleanup Tests
// ============================================

func TestEvidSvc_ApplyUpdates_ChangesImageURL(t *testing.T) {
	_, _, svc := setupEvidenceService()

	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "old-url"}
	updates := &domain.MotorcycleEvidence{ImageURL: "new-url"}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	assert.Equal(t, "new-url", existing.ImageURL)
}

func TestEvidSvc_ApplyUpdates_ChangesAngle(t *testing.T) {
	_, _, svc := setupEvidenceService()

	angle := domain.EvidenceAngleFrontal
	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "url"}
	updates := &domain.MotorcycleEvidence{Angle: &angle}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	assert.NotNil(t, existing.Angle)
	assert.Equal(t, domain.EvidenceAngleFrontal, *existing.Angle)
}

func TestEvidSvc_ApplyUpdates_ChangesDescription(t *testing.T) {
	_, _, svc := setupEvidenceService()

	desc := "nueva descripción"
	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "url"}
	updates := &domain.MotorcycleEvidence{Description: &desc}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	assert.NotNil(t, existing.Description)
	assert.Equal(t, "nueva descripción", *existing.Description)
}

func TestEvidSvc_ApplyUpdates_WithStorageCleanup(t *testing.T) {
	_, _, svc := setupEvidenceService()
	storageClient := new(mockStorageClient)
	svc.WithStorageClient(storageClient)

	storageClient.On("DeleteStorageFile", mock.Anything, "old-url").Return(nil)

	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "old-url"}
	updates := &domain.MotorcycleEvidence{ImageURL: "new-url"}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	assert.Equal(t, "new-url", existing.ImageURL)
	storageClient.AssertExpectations(t)
}

func TestEvidSvc_ApplyUpdates_StorageCleanupError(t *testing.T) {
	_, _, svc := setupEvidenceService()
	storageClient := new(mockStorageClient)
	svc.WithStorageClient(storageClient)

	storageClient.On("DeleteStorageFile", mock.Anything, "old-url").Return(errors.New("storage error"))

	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "old-url"}
	updates := &domain.MotorcycleEvidence{ImageURL: "new-url"}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	// Should still update the URL even if storage cleanup fails
	assert.Equal(t, "new-url", existing.ImageURL)
}

func TestEvidSvc_ApplyUpdates_SameURL_NoCleanup(t *testing.T) {
	_, _, svc := setupEvidenceService()

	existing := &domain.MotorcycleEvidence{ID: "ev-1", ImageURL: "same-url"}
	updates := &domain.MotorcycleEvidence{ImageURL: "same-url"}

	svc.ApplyUpdatesAndCleanup(context.Background(), existing, updates)
	assert.Equal(t, "same-url", existing.ImageURL)
}

// ============================================
// UpdateEvidence Tests
// ============================================

func TestEvidSvc_UpdateEvidence_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)
	evidence := &domain.MotorcycleEvidence{ID: "ev-1"}

	evidRepo.On("Update", mock.Anything, tx, evidence).Return(nil)

	err := svc.UpdateEvidence(context.Background(), tx, evidence)
	assert.NoError(t, err)
}

func TestEvidSvc_UpdateEvidence_Error(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)
	evidence := &domain.MotorcycleEvidence{ID: "ev-bad"}

	evidRepo.On("Update", mock.Anything, tx, evidence).Return(errors.New("update error"))

	err := svc.UpdateEvidence(context.Background(), tx, evidence)
	assert.Equal(t, domain.ErrEvidenceCannotUpdate, err)
}

// ============================================
// DeleteEvidence Tests
// ============================================

func TestEvidSvc_DeleteEvidence_Success(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)

	evidRepo.On("Delete", mock.Anything, tx, "ev-1").Return(nil)

	err := svc.DeleteEvidence(context.Background(), tx, "ev-1")
	assert.NoError(t, err)
}

func TestEvidSvc_DeleteEvidence_Error(t *testing.T) {
	evidRepo, _, svc := setupEvidenceService()
	tx := new(mockEvidTx)

	evidRepo.On("Delete", mock.Anything, tx, "ev-bad").Return(errors.New("delete error"))

	err := svc.DeleteEvidence(context.Background(), tx, "ev-bad")
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}

// ============================================
// DeleteStorageFile Tests
// ============================================

func TestEvidSvc_DeleteStorageFile_Success(t *testing.T) {
	_, _, svc := setupEvidenceService()
	storageClient := new(mockStorageClient)
	svc.WithStorageClient(storageClient)

	storageClient.On("DeleteStorageFile", mock.Anything, "https://storage/file.jpg").Return(nil)

	svc.DeleteStorageFile(context.Background(), "https://storage/file.jpg")
	storageClient.AssertExpectations(t)
}

func TestEvidSvc_DeleteStorageFile_Error(t *testing.T) {
	_, _, svc := setupEvidenceService()
	storageClient := new(mockStorageClient)
	svc.WithStorageClient(storageClient)

	storageClient.On("DeleteStorageFile", mock.Anything, "bad-url").Return(errors.New("storage error"))

	// Should not panic - best effort
	svc.DeleteStorageFile(context.Background(), "bad-url")
	storageClient.AssertExpectations(t)
}

func TestEvidSvc_DeleteStorageFile_EmptyURL(t *testing.T) {
	_, _, svc := setupEvidenceService()
	// Should be a no-op
	svc.DeleteStorageFile(context.Background(), "")
}

func TestEvidSvc_DeleteStorageFile_NoClient(t *testing.T) {
	_, _, svc := setupEvidenceService()
	// No storage client set - should be a no-op
	svc.DeleteStorageFile(context.Background(), "https://storage/file.jpg")
}
