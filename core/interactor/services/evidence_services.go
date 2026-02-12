package services

import (
	"strings"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

const (
	// FirebaseStorageHost is the required host for valid Firebase Storage URLs
	FirebaseStorageHost = "firebasestorage.googleapis.com"
	// MaxEvidencePerMotorcycle is the maximum number of evidence allowed per motorcycle
	MaxEvidencePerMotorcycle = 4
)

// NewEvidence creates a new MotorcycleEvidence with generated ID and timestamp
// This is a factory function that encapsulates ID generation and timestamp assignment
func NewEvidence(motorcycleID, imageURL string, angle, description *string) *domain.MotorcycleEvidence {
	e := &domain.MotorcycleEvidence{
		MotorcycleID: motorcycleID,
		Angle:        angle,
		ImageURL:     imageURL,
		Description:  description,
		CreatedAt:    time.Now(),
	}
	e.SetID()
	return e
}

// IsValidAngle checks if the angle is valid (FRONTAL, LATERAL, REAR)
func IsValidAngle(angle string) bool {
	for _, valid := range domain.ValidEvidenceAngles {
		if angle == valid {
			return true
		}
	}
	return false
}

// IsValidFirebaseURL validates that the URL is from Firebase Storage
func IsValidFirebaseURL(url string) bool {
	return strings.HasPrefix(url, "https://"+FirebaseStorageHost)
}

// IsEvidenceLimitReached checks if the evidence limit per motorcycle is reached
func IsEvidenceLimitReached(count int) bool {
	return count >= MaxEvidencePerMotorcycle
}
