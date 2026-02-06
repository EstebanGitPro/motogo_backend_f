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

// Test fixtures
var (
	testMotorcycleID = "moto-123"
	testOwnerID      = "owner-456"
	testEvidenceID   = "evidence-789"
	testImageURL     = "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	testAngleFront   = domain.EvidenceAngleFrontal
)

func TestNewEvidenceInteractor(t *testing.T) {
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	assert.NotNil(t, interactorInstance)
}

// ============================================
// CreateEvidence Tests (HU16)
// ============================================

func TestCreateEvidence_Success(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
		Angle:    &testAngleFront,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(2, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.MotorcycleEvidence")).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, testMotorcycleID, result.MotorcycleID)
	assert.Equal(t, testImageURL, result.ImageURL)
	mockMotorcycleRepo.AssertExpectations(t)
	mockEvidenceRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestCreateEvidence_MotorcycleNotFound(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestCreateEvidence_NotOwner(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: "different-owner",
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err) // Security by obscurity
}

// NOTE: TestCreateEvidence_InvalidURL and TestCreateEvidence_InvalidAngle removed
// These validations are now handled by JSON Schema middleware (create_evidence_schema.json)

func TestCreateEvidence_LimitExceeded(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(5, nil) // Already at limit

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceLimitExceeded, err)
}

func TestCreateEvidence_CountError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(0, errors.New("db error"))

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

func TestCreateEvidence_BeginTxError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(2, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

func TestCreateEvidence_SaveError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(2, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.MotorcycleEvidence")).Return(errors.New("save error"))
	mockTx.On("Rollback").Return(nil)

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

func TestCreateEvidence_CommitError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidence := &domain.MotorcycleEvidence{
		ImageURL: testImageURL,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("CountByMotorcycleID", mock.Anything, testMotorcycleID).Return(2, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Save", mock.Anything, mockTx, mock.AnythingOfType("*domain.MotorcycleEvidence")).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))

	// Act
	result, err := interactorInstance.CreateEvidence(context.Background(), testMotorcycleID, testOwnerID, evidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceCannotSave, err)
}

// ============================================
// GetEvidenceByID Tests (HU18)
// ============================================

func TestGetEvidenceByID_Success(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
		CreatedAt:    time.Now(),
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)

	// Act
	result, err := interactorInstance.GetEvidenceByID(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testEvidenceID, result.ID)
}

func TestGetEvidenceByID_NotFound(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(nil, domain.ErrEvidenceNotFound)

	// Act
	result, err := interactorInstance.GetEvidenceByID(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestGetEvidenceByID_NotOwner(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: "different-owner",
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)

	// Act
	result, err := interactorInstance.GetEvidenceByID(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err) // Security by obscurity
}

// ============================================
// ListEvidenceByMotorcycle Tests (HU18)
// ============================================

func TestListEvidenceByMotorcycle_Success(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	evidences := []domain.MotorcycleEvidence{
		{ID: testEvidenceID, MotorcycleID: testMotorcycleID, ImageURL: testImageURL},
		{ID: "evidence-2", MotorcycleID: testMotorcycleID, ImageURL: testImageURL},
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("GetByMotorcycleID", mock.Anything, testMotorcycleID).Return(evidences, nil)

	// Act
	result, err := interactorInstance.ListEvidenceByMotorcycle(context.Background(), testMotorcycleID, testOwnerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListEvidenceByMotorcycle_MotorcycleNotFound(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	result, err := interactorInstance.ListEvidenceByMotorcycle(context.Background(), testMotorcycleID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListEvidenceByMotorcycle_NotOwner(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: "different-owner",
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)

	// Act
	result, err := interactorInstance.ListEvidenceByMotorcycle(context.Background(), testMotorcycleID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestListEvidenceByMotorcycle_EmptyList(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("GetByMotorcycleID", mock.Anything, testMotorcycleID).Return([]domain.MotorcycleEvidence{}, nil)

	// Act
	result, err := interactorInstance.ListEvidenceByMotorcycle(context.Background(), testMotorcycleID, testOwnerID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================
// DeleteEvidence Tests (HU19)
// ============================================

func TestDeleteEvidence_Success(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
		ImageURL:     testImageURL,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Delete", mock.Anything, mockTx, testEvidenceID).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.NoError(t, err)
	mockEvidenceRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestDeleteEvidence_NotFound(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(nil, domain.ErrEvidenceNotFound)

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestDeleteEvidence_NotOwner(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: "different-owner",
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
}

func TestDeleteEvidence_BeginTxError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("tx error"))

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}

func TestDeleteEvidence_DeleteError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Delete", mock.Anything, mockTx, testEvidenceID).Return(errors.New("delete error"))
	mockTx.On("Rollback").Return(nil)

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}

func TestDeleteEvidence_CommitError(t *testing.T) {
	// Arrange
	mockEvidenceRepo := new(mocks.MockEvidenceRepository)
	mockMotorcycleRepo := new(mocks.MockMotorcycleRepository)

	mockTx := new(mocks.MockTx)

	interactorInstance := interactor.NewEvidenceInteractor(mockEvidenceRepo, mockMotorcycleRepo)

	evidence := &domain.MotorcycleEvidence{
		ID:           testEvidenceID,
		MotorcycleID: testMotorcycleID,
	}

	motorcycle := &domain.Motorcycle{
		ID:      testMotorcycleID,
		OwnerID: testOwnerID,
	}

	mockEvidenceRepo.On("GetByID", mock.Anything, testEvidenceID).Return(evidence, nil)
	mockMotorcycleRepo.On("GetByID", mock.Anything, testMotorcycleID).Return(motorcycle, nil)
	mockEvidenceRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEvidenceRepo.On("Delete", mock.Anything, mockTx, testEvidenceID).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))

	// Act
	err := interactorInstance.DeleteEvidence(context.Background(), testEvidenceID, testOwnerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceCannotDelete, err)
}
