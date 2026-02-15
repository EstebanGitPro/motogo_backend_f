package handlers_test

import (
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// CreateEvidenceRequest.Sanitize Tests
// ============================================

func TestCreateEvidenceRequest_Sanitize(t *testing.T) {
	// Arrange
	angle := "  FRONTAL  "
	description := "  Foto del tanque\t"
	req := &handlers.CreateEvidenceRequest{
		ImageURL:    "  https://firebasestorage.googleapis.com/image.jpg  ",
		Angle:       &angle,
		Description: &description,
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "https://firebasestorage.googleapis.com/image.jpg", req.ImageURL)
	assert.Equal(t, "FRONTAL", *req.Angle)
	assert.Equal(t, "Foto del tanque", *req.Description)
}

func TestCreateEvidenceRequest_Sanitize_NilFields(t *testing.T) {
	// Arrange
	req := &handlers.CreateEvidenceRequest{
		ImageURL: "https://storage.com/img.jpg",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "https://storage.com/img.jpg", req.ImageURL)
	assert.Nil(t, req.Angle)
	assert.Nil(t, req.Description)
}

// ============================================
// CreateEvidenceRequest.ToDomain Tests
// ============================================

func TestCreateEvidenceRequest_ToDomain(t *testing.T) {
	// Arrange
	angle := "LATERAL"
	description := "Vista lateral"
	req := &handlers.CreateEvidenceRequest{
		ImageURL:    "https://storage.com/lateral.jpg",
		Angle:       &angle,
		Description: &description,
	}

	// Act
	result := req.ToDomain()

	// Assert
	assert.Equal(t, "https://storage.com/lateral.jpg", result.ImageURL)
	assert.NotNil(t, result.Angle)
	assert.Equal(t, domain.EvidenceAngle("LATERAL"), *result.Angle)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Vista lateral", *result.Description)
}

func TestCreateEvidenceRequest_ToDomain_NilOptionals(t *testing.T) {
	// Arrange
	req := &handlers.CreateEvidenceRequest{
		ImageURL: "https://storage.com/image.jpg",
	}

	// Act
	result := req.ToDomain()

	// Assert
	assert.Equal(t, "https://storage.com/image.jpg", result.ImageURL)
	assert.Nil(t, result.Angle)
	assert.Nil(t, result.Description)
}

// ============================================
// ToEvidenceResponse Tests
// ============================================

func TestToEvidenceResponse_FullData(t *testing.T) {
	// Arrange
	angle := domain.EvidenceAngle("FRONTAL")
	description := "Foto frontal"
	evidence := &domain.MotorcycleEvidence{
		ID:           "evidence-123",
		MotorcycleID: "moto-456",
		Angle:        &angle,
		ImageURL:     "https://storage.com/image.jpg",
		Description:  &description,
		CreatedAt:    time.Date(2026, 2, 6, 10, 30, 0, 0, time.UTC),
	}

	// Act
	response := handlers.ToEvidenceResponse(evidence)

	// Assert
	assert.Equal(t, "evidence-123", response.ID)
	assert.Equal(t, "moto-456", response.MotorcycleID)
	assert.NotNil(t, response.Angle)
	assert.Equal(t, "FRONTAL", *response.Angle)
	assert.Equal(t, "https://storage.com/image.jpg", response.ImageURL)
	assert.NotNil(t, response.Description)
	assert.Equal(t, "Foto frontal", *response.Description)
	assert.Contains(t, response.CreatedAt, "2026-02-06")
}

func TestToEvidenceResponse_MinimalData(t *testing.T) {
	// Arrange
	evidence := &domain.MotorcycleEvidence{
		ID:           "evidence-789",
		MotorcycleID: "moto-abc",
		ImageURL:     "https://storage.com/minimal.jpg",
		CreatedAt:    time.Now(),
	}

	// Act
	response := handlers.ToEvidenceResponse(evidence)

	// Assert
	assert.Equal(t, "evidence-789", response.ID)
	assert.Equal(t, "moto-abc", response.MotorcycleID)
	assert.Nil(t, response.Angle)
	assert.Nil(t, response.Description)
	assert.NotEmpty(t, response.CreatedAt)
}

// ============================================
// ToEvidenceResponseList Tests
// ============================================

func TestToEvidenceResponseList_MultipleItems(t *testing.T) {
	// Arrange
	evidences := []domain.MotorcycleEvidence{
		{ID: "ev-1", MotorcycleID: "moto-1", ImageURL: "url1", CreatedAt: time.Now()},
		{ID: "ev-2", MotorcycleID: "moto-1", ImageURL: "url2", CreatedAt: time.Now()},
		{ID: "ev-3", MotorcycleID: "moto-1", ImageURL: "url3", CreatedAt: time.Now()},
	}

	// Act
	responses := handlers.ToEvidenceResponseList(evidences)

	// Assert
	assert.Len(t, responses, 3)
	assert.Equal(t, "ev-1", responses[0].ID)
	assert.Equal(t, "ev-2", responses[1].ID)
	assert.Equal(t, "ev-3", responses[2].ID)
}

func TestToEvidenceResponseList_EmptyList(t *testing.T) {
	// Arrange
	evidences := []domain.MotorcycleEvidence{}

	// Act
	responses := handlers.ToEvidenceResponseList(evidences)

	// Assert
	assert.Empty(t, responses)
	assert.NotNil(t, responses)
}
