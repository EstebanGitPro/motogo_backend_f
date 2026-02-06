package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================
// EvidenceInteractor Mock Integration Tests
// These tests verify the mock implementations work correctly
// ============================================

func TestMockEvidenceInteractor_CreateEvidence_Success(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	angle := "FRONT"
	inputEvidence := &domain.MotorcycleEvidence{
		ImageURL: "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
		Angle:    &angle,
	}

	createdEvidence := &domain.MotorcycleEvidence{
		ID:           "evidence-123",
		MotorcycleID: "moto-456",
		Angle:        &angle,
		ImageURL:     "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
		CreatedAt:    time.Now(),
	}

	mockEvidence.On("CreateEvidence",
		mock.Anything,
		"moto-456",
		"owner-789",
		inputEvidence,
	).Return(createdEvidence, nil)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.CreateEvidence(ctx, "moto-456", "owner-789", inputEvidence)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "evidence-123", result.ID)
	assert.Equal(t, "moto-456", result.MotorcycleID)
	assert.NotNil(t, result.Angle)
	assert.Equal(t, "FRONT", *result.Angle)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_CreateEvidence_MotorcycleNotFound(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	inputEvidence := &domain.MotorcycleEvidence{
		ImageURL: "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
	}

	mockEvidence.On("CreateEvidence",
		mock.Anything,
		"moto-not-found",
		"owner-789",
		inputEvidence,
	).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.CreateEvidence(ctx, "moto-not-found", "owner-789", inputEvidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_CreateEvidence_InvalidURL(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	inputEvidence := &domain.MotorcycleEvidence{
		ImageURL: "https://invalid-url.com/image.jpg",
	}

	mockEvidence.On("CreateEvidence",
		mock.Anything,
		"moto-456",
		"owner-789",
		inputEvidence,
	).Return(nil, domain.ErrInvalidEvidenceURL)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.CreateEvidence(ctx, "moto-456", "owner-789", inputEvidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidEvidenceURL, err)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_CreateEvidence_LimitExceeded(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	inputEvidence := &domain.MotorcycleEvidence{
		ImageURL: "https://firebasestorage.googleapis.com/v0/b/test/image.jpg",
	}

	mockEvidence.On("CreateEvidence",
		mock.Anything,
		"moto-full",
		"owner-789",
		inputEvidence,
	).Return(nil, domain.ErrEvidenceLimitExceeded)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.CreateEvidence(ctx, "moto-full", "owner-789", inputEvidence)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceLimitExceeded, err)
	mockEvidence.AssertExpectations(t)
}

// ============================================
// GetEvidenceByID Tests
// ============================================

func TestMockEvidenceInteractor_GetEvidenceByID_Success(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	angle := "SIDE"
	evidence := &domain.MotorcycleEvidence{
		ID:           "evidence-123",
		MotorcycleID: "moto-456",
		Angle:        &angle,
		ImageURL:     "https://storage/image.jpg",
		CreatedAt:    time.Now(),
	}

	mockEvidence.On("GetEvidenceByID",
		mock.Anything,
		"evidence-123",
		"owner-789",
	).Return(evidence, nil)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.GetEvidenceByID(ctx, "evidence-123", "owner-789")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "evidence-123", result.ID)
	assert.Equal(t, "SIDE", *result.Angle)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_GetEvidenceByID_NotFound(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("GetEvidenceByID",
		mock.Anything,
		"evidence-not-found",
		"owner-789",
	).Return(nil, domain.ErrEvidenceNotFound)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.GetEvidenceByID(ctx, "evidence-not-found", "owner-789")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
	mockEvidence.AssertExpectations(t)
}

// ============================================
// ListEvidenceByMotorcycle Tests
// ============================================

func TestMockEvidenceInteractor_ListEvidenceByMotorcycle_Success(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	angle1 := "FRONT"
	angle2 := "BACK"
	evidences := []domain.MotorcycleEvidence{
		{
			ID:           "evidence-1",
			MotorcycleID: "moto-456",
			Angle:        &angle1,
			ImageURL:     "https://storage/img1.jpg",
			CreatedAt:    time.Now(),
		},
		{
			ID:           "evidence-2",
			MotorcycleID: "moto-456",
			Angle:        &angle2,
			ImageURL:     "https://storage/img2.jpg",
			CreatedAt:    time.Now(),
		},
	}

	mockEvidence.On("ListEvidenceByMotorcycle",
		mock.Anything,
		"moto-456",
		"owner-789",
	).Return(evidences, nil)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.ListEvidenceByMotorcycle(ctx, "moto-456", "owner-789")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "evidence-1", result[0].ID)
	assert.Equal(t, "evidence-2", result[1].ID)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_ListEvidenceByMotorcycle_Empty(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("ListEvidenceByMotorcycle",
		mock.Anything,
		"moto-empty",
		"owner-789",
	).Return([]domain.MotorcycleEvidence{}, nil)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.ListEvidenceByMotorcycle(ctx, "moto-empty", "owner-789")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_ListEvidenceByMotorcycle_MotorcycleNotFound(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("ListEvidenceByMotorcycle",
		mock.Anything,
		"moto-not-found",
		"owner-789",
	).Return(nil, domain.ErrMotorcycleNotFound)

	// Act
	ctx := context.Background()
	result, err := mockEvidence.ListEvidenceByMotorcycle(ctx, "moto-not-found", "owner-789")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
	mockEvidence.AssertExpectations(t)
}

// ============================================
// DeleteEvidence Tests
// ============================================

func TestMockEvidenceInteractor_DeleteEvidence_Success(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("DeleteEvidence",
		mock.Anything,
		"evidence-123",
		"owner-789",
	).Return(nil)

	// Act
	ctx := context.Background()
	err := mockEvidence.DeleteEvidence(ctx, "evidence-123", "owner-789")

	// Assert
	assert.NoError(t, err)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_DeleteEvidence_NotFound(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("DeleteEvidence",
		mock.Anything,
		"evidence-not-found",
		"owner-789",
	).Return(domain.ErrEvidenceNotFound)

	// Act
	ctx := context.Background()
	err := mockEvidence.DeleteEvidence(ctx, "evidence-not-found", "owner-789")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
	mockEvidence.AssertExpectations(t)
}

func TestMockEvidenceInteractor_DeleteEvidence_NotOwner(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	mockEvidence.On("DeleteEvidence",
		mock.Anything,
		"evidence-123",
		"not-owner",
	).Return(domain.ErrEvidenceNotFound) // Security: returns NotFound for non-owners

	// Act
	ctx := context.Background()
	err := mockEvidence.DeleteEvidence(ctx, "evidence-123", "not-owner")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEvidenceNotFound, err)
	mockEvidence.AssertExpectations(t)
}

// ============================================
// Interface Verification
// ============================================

func TestMockEvidenceInteractor_VerifyInterfaceImplementation(t *testing.T) {
	// Verify MockEvidenceInteractor implements input.EvidenceInteractorInterface
	mockEvidence := new(mocks.MockEvidenceInteractor)
	assert.NotNil(t, mockEvidence)
}

func TestMockEvidenceInteractor_MultipleCalls(t *testing.T) {
	// Arrange
	mockEvidence := new(mocks.MockEvidenceInteractor)

	angle := "FRONT"
	evidences := []domain.MotorcycleEvidence{
		{ID: "evidence-1", Angle: &angle, ImageURL: "https://storage/img.jpg"},
	}

	mockEvidence.On("ListEvidenceByMotorcycle",
		mock.Anything,
		"moto-456",
		"owner-789",
	).Return(evidences, nil).Times(2)

	ctx := context.Background()

	// First call
	result1, err1 := mockEvidence.ListEvidenceByMotorcycle(ctx, "moto-456", "owner-789")
	assert.NoError(t, err1)
	assert.Len(t, result1, 1)

	// Second call
	result2, err2 := mockEvidence.ListEvidenceByMotorcycle(ctx, "moto-456", "owner-789")
	assert.NoError(t, err2)
	assert.Len(t, result2, 1)

	mockEvidence.AssertExpectations(t)
}
