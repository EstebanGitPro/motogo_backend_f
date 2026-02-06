package services_test

import (
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewEvidence Tests
// ============================================

func TestNewEvidence_GeneratesUUID(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, nil, nil)

	// Assert
	assert.NotEmpty(t, evidence.ID)
	assert.Len(t, evidence.ID, 36) // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
}

func TestNewEvidence_SetsMotorcycleID(t *testing.T) {
	// Arrange
	motorcycleID := "moto-456"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, nil, nil)

	// Assert
	assert.Equal(t, motorcycleID, evidence.MotorcycleID)
}

func TestNewEvidence_SetsImageURL(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/custom-image.jpg"

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, nil, nil)

	// Assert
	assert.Equal(t, imageURL, evidence.ImageURL)
}

func TestNewEvidence_SetsOptionalAngle(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	angle := domain.EvidenceAngleFrontal

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, &angle, nil)

	// Assert
	assert.NotNil(t, evidence.Angle)
	assert.Equal(t, domain.EvidenceAngleFrontal, *evidence.Angle)
}

func TestNewEvidence_SetsOptionalDescription(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	description := "Foto frontal del tanque"

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, nil, &description)

	// Assert
	assert.NotNil(t, evidence.Description)
	assert.Equal(t, "Foto frontal del tanque", *evidence.Description)
}

func TestNewEvidence_SetsCreatedAtTimestamp(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	beforeCall := time.Now()

	// Act
	evidence := services.NewEvidence(motorcycleID, imageURL, nil, nil)

	// Assert
	afterCall := time.Now()
	assert.True(t, evidence.CreatedAt.After(beforeCall) || evidence.CreatedAt.Equal(beforeCall))
	assert.True(t, evidence.CreatedAt.Before(afterCall) || evidence.CreatedAt.Equal(afterCall))
}

func TestNewEvidence_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"

	// Act
	evidence1 := services.NewEvidence(motorcycleID, imageURL, nil, nil)
	evidence2 := services.NewEvidence(motorcycleID, imageURL, nil, nil)

	// Assert
	assert.NotEqual(t, evidence1.ID, evidence2.ID)
}

// ============================================
// IsValidAngle Tests
// ============================================

func TestIsValidAngle_Frontal(t *testing.T) {
	assert.True(t, services.IsValidAngle("FRONTAL"))
}

func TestIsValidAngle_Lateral(t *testing.T) {
	assert.True(t, services.IsValidAngle("LATERAL"))
}

func TestIsValidAngle_Rear(t *testing.T) {
	assert.True(t, services.IsValidAngle("REAR"))
}

func TestIsValidAngle_InvalidAngle(t *testing.T) {
	assert.False(t, services.IsValidAngle("INVALID"))
}

func TestIsValidAngle_EmptyString(t *testing.T) {
	assert.False(t, services.IsValidAngle(""))
}

func TestIsValidAngle_LowercaseInvalid(t *testing.T) {
	// Angles must be uppercase
	assert.False(t, services.IsValidAngle("frontal"))
	assert.False(t, services.IsValidAngle("lateral"))
	assert.False(t, services.IsValidAngle("rear"))
}

func TestIsValidAngle_OldFormatsInvalid(t *testing.T) {
	// Old formats should not be valid anymore
	assert.False(t, services.IsValidAngle("FRONT"))
	assert.False(t, services.IsValidAngle("SIDE"))
	assert.False(t, services.IsValidAngle("BACK"))
}

// ============================================
// IsValidFirebaseURL Tests
// ============================================

func TestIsValidFirebaseURL_ValidURL(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/project/o/image.jpg"
	assert.True(t, services.IsValidFirebaseURL(url))
}

func TestIsValidFirebaseURL_ValidURLWithParams(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/project/o/image.jpg?alt=media&token=abc123"
	assert.True(t, services.IsValidFirebaseURL(url))
}

func TestIsValidFirebaseURL_InvalidExternalURL(t *testing.T) {
	url := "https://example.com/image.jpg"
	assert.False(t, services.IsValidFirebaseURL(url))
}

func TestIsValidFirebaseURL_InvalidS3URL(t *testing.T) {
	url := "https://s3.amazonaws.com/bucket/image.jpg"
	assert.False(t, services.IsValidFirebaseURL(url))
}

func TestIsValidFirebaseURL_InvalidHTTPURL(t *testing.T) {
	// Must be HTTPS
	url := "http://firebasestorage.googleapis.com/v0/b/project/o/image.jpg"
	assert.False(t, services.IsValidFirebaseURL(url))
}

func TestIsValidFirebaseURL_EmptyString(t *testing.T) {
	assert.False(t, services.IsValidFirebaseURL(""))
}

func TestIsValidFirebaseURL_PartialMatch(t *testing.T) {
	// URL containing Firebase host but not starting with it
	url := "https://malicious.com/https://firebasestorage.googleapis.com/image.jpg"
	assert.False(t, services.IsValidFirebaseURL(url))
}

// ============================================
// IsEvidenceLimitReached Tests
// ============================================

func TestIsEvidenceLimitReached_AtLimit(t *testing.T) {
	assert.True(t, services.IsEvidenceLimitReached(5))
}

func TestIsEvidenceLimitReached_OverLimit(t *testing.T) {
	assert.True(t, services.IsEvidenceLimitReached(10))
}

func TestIsEvidenceLimitReached_UnderLimit(t *testing.T) {
	assert.False(t, services.IsEvidenceLimitReached(3))
}

func TestIsEvidenceLimitReached_ZeroCount(t *testing.T) {
	assert.False(t, services.IsEvidenceLimitReached(0))
}

func TestIsEvidenceLimitReached_JustUnderLimit(t *testing.T) {
	assert.False(t, services.IsEvidenceLimitReached(services.MaxEvidencePerMotorcycle-1))
}

// ============================================
// Constants Tests
// ============================================

func TestConstants_FirebaseStorageHost(t *testing.T) {
	assert.Equal(t, "firebasestorage.googleapis.com", services.FirebaseStorageHost)
}

func TestConstants_MaxEvidencePerMotorcycle(t *testing.T) {
	assert.Equal(t, 4, services.MaxEvidencePerMotorcycle)
}
