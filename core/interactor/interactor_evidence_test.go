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

// Test fixtures
var (
	testMotorcycleID = "moto-123"
	testOwnerID      = "owner-456"
	testEvidenceID   = "evidence-789"
	testImageURL     = "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	testAngleFront   = domain.EvidenceAngleFrontal
)

// helper to create EvidenceInteractor with fresh mock service
func setupEvidenceInteractor() (*interactor.EvidenceInteractor, *mocks.MockEvidenceService) {
	svc := new(mocks.MockEvidenceService)
	ei := interactor.NewEvidenceInteractor(svc)
	return ei, svc
}

func TestNewEvidenceInteractor(t *testing.T) {
	svc := new(mocks.MockEvidenceService)
	ei := interactor.NewEvidenceInteractor(svc)
	assert.NotNil(t, ei)
}

// ============================================
// CreateEvidence Tests (HU16)
// ============================================

func TestCreateEvidence_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
		Angle:    &testAngleFront,
	}

	created := &domain.MotorcycleEvidence{
		ID:           "new-ev-id",
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
		Angle:        &testAngleFront,
	}

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("CreateEvidence", ctx, mockTx, testMotorcycleID, evidence).Return(created, nil)
	mockTx.On("Commit").Return(nil)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, evidence)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-ev-id", result.ID)
	assert.Equal(t, testMotorcycleID, result.MotorcycleID)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestCreateEvidence_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, &domain.MotorcycleEvidence{ImageURL: testImageURL})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestCreateEvidence_NotOwner(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, &domain.MotorcycleEvidence{ImageURL: testImageURL})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestCreateEvidence_LimitExceeded(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(domain.ErrEvidenceLimitExceeded)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, &domain.MotorcycleEvidence{ImageURL: testImageURL})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceLimitExceeded, err)
}

func TestCreateEvidence_CountError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(domain.ErrEvidenceCannotSave)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, &domain.MotorcycleEvidence{ImageURL: testImageURL})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

func TestCreateEvidence_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, &domain.MotorcycleEvidence{ImageURL: testImageURL})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

func TestCreateEvidence_SaveError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{ImageURL: testImageURL}

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("CreateEvidence", ctx, mockTx, testMotorcycleID, evidence).Return(nil, domain.ErrEvidenceCannotSave)
	mockTx.On("Rollback").Return(nil)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, evidence)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestCreateEvidence_CommitError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{ImageURL: testImageURL}
	created := &domain.MotorcycleEvidence{ID: "ev-1", MotorcycleID: testMotorcycleID}

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("CheckEvidenceLimit", ctx, testMotorcycleID).Return(nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("CreateEvidence", ctx, mockTx, testMotorcycleID, evidence).Return(created, nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	result, err := ei.CreateEvidence(ctx, testMotorcycleID, testOwnerID, evidence)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GetEvidenceByID Tests (HU18)
// ============================================

func TestGetEvidenceByID_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)

	result, err := ei.GetEvidenceByID(ctx, testEvidenceID, testOwnerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testEvidenceID, result.ID)
	svc.AssertExpectations(t)
}

func TestGetEvidenceByID_NotFound(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("GetByID", ctx, testEvidenceID).Return(nil, domain.ErrEvidenceNotFound)

	result, err := ei.GetEvidenceByID(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestGetEvidenceByID_NotOwner(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.GetEvidenceByID(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err) // Security by obscurity
}

// ============================================
// ListEvidenceByMotorcycle Tests (HU18)
// ============================================

func TestListEvidenceByMotorcycle_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidences := []domain.MotorcycleEvidence{
		{ID: testEvidenceID, MotorcycleID: testMotorcycleID, ImageURL: testImageURL},
		{ID: "evidence-2", MotorcycleID: testMotorcycleID, ImageURL: testImageURL},
	}

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("GetByMotorcycleID", ctx, testMotorcycleID).Return(evidences, nil)

	result, err := ei.ListEvidenceByMotorcycle(ctx, testMotorcycleID, testOwnerID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestListEvidenceByMotorcycle_MotorcycleNotFound(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.ListEvidenceByMotorcycle(ctx, testMotorcycleID, testOwnerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListEvidenceByMotorcycle_NotOwner(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.ListEvidenceByMotorcycle(ctx, testMotorcycleID, testOwnerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListEvidenceByMotorcycle_EmptyList(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("GetByMotorcycleID", ctx, testMotorcycleID).Return([]domain.MotorcycleEvidence{}, nil)

	result, err := ei.ListEvidenceByMotorcycle(ctx, testMotorcycleID, testOwnerID)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// DeleteEvidence Tests (HU19)
// ============================================

func TestDeleteEvidence_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("DeleteStorageFile", ctx, testImageURL).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteEvidence", ctx, mockTx, testEvidenceID).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteEvidence_NotFound(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("GetByID", ctx, testEvidenceID).Return(nil, domain.ErrEvidenceNotFound)

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestDeleteEvidence_NotOwner(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestDeleteEvidence_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("DeleteStorageFile", ctx, testImageURL).Return()
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}

func TestDeleteEvidence_DeleteError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("DeleteStorageFile", ctx, testImageURL).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteEvidence", ctx, mockTx, testEvidenceID).Return(domain.ErrEvidenceCannotDelete)
	mockTx.On("Rollback").Return(nil)

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestDeleteEvidence_CommitError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(evidence, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("DeleteStorageFile", ctx, testImageURL).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("DeleteEvidence", ctx, mockTx, testEvidenceID).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := ei.DeleteEvidence(ctx, testEvidenceID, testOwnerID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// UpdateEvidence Tests (HU17)
// ============================================

func TestUpdateEvidence_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
		Angle:        &testAngleFront,
	}

	newAngle := domain.EvidenceAngleLateral
	updates := &domain.MotorcycleEvidence{Angle: &newAngle}

	svc.On("GetByID", ctx, testEvidenceID).Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("ApplyUpdatesAndCleanup", ctx, existing, updates).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateEvidence", ctx, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(nil)

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testEvidenceID, result.ID)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUpdateEvidence_NotFound(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("GetByID", ctx, testEvidenceID).Return(nil, domain.ErrEvidenceNotFound)

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, &domain.MotorcycleEvidence{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestUpdateEvidence_NotOwner(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	existing := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(nil, domain.ErrMotorcycleNotFound)

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, &domain.MotorcycleEvidence{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err) // Security by obscurity
}

func TestUpdateEvidence_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	existing := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("ApplyUpdatesAndCleanup", ctx, existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, &domain.MotorcycleEvidence{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotUpdate, err)
}

func TestUpdateEvidence_UpdateError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("ApplyUpdatesAndCleanup", ctx, existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateEvidence", ctx, mockTx, existing).Return(domain.ErrEvidenceCannotUpdate)
	mockTx.On("Rollback").Return(nil)

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, &domain.MotorcycleEvidence{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestUpdateEvidence_CommitError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()
	mockTx := new(mocks.MockTx)

	existing := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	svc.On("GetByID", ctx, testEvidenceID).Return(existing, nil)
	svc.On("ValidateMotorcycleOwnership", ctx, testMotorcycleID, testOwnerID).Return(&domain.Motorcycle{ID: testMotorcycleID, OwnerID: testOwnerID}, nil)
	svc.On("ApplyUpdatesAndCleanup", ctx, existing, mock.Anything).Return()
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("UpdateEvidence", ctx, mockTx, existing).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	result, err := ei.UpdateEvidence(ctx, testEvidenceID, testOwnerID, &domain.MotorcycleEvidence{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotUpdate, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// LookupEvidence Tests (workshop use — no ownership)
// ============================================

func TestLookupEvidence_Success(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	evidences := []domain.MotorcycleEvidence{
		{ID: "ev-1", MotorcycleID: testMotorcycleID},
	}

	svc.On("GetByMotorcycleID", ctx, testMotorcycleID).Return(evidences, nil)

	result, err := ei.LookupEvidence(ctx, testMotorcycleID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestLookupEvidence_RepoError(t *testing.T) {
	ctx := context.Background()
	ei, svc := setupEvidenceInteractor()

	svc.On("GetByMotorcycleID", ctx, testMotorcycleID).Return(nil, errors.New("db error"))

	result, err := ei.LookupEvidence(ctx, testMotorcycleID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
