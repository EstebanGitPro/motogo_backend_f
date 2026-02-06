package domain

import "time"

// MotorcycleEvidence represents a photographic evidence of a motorcycle (HU16-19)
type MotorcycleEvidence struct {
	ID           string    // UUID primary key
	MotorcycleID string    // FK to motorcycle
	Angle        *string   // Optional: FRONTAL, LATERAL, REAR
	ImageURL     string    // Firebase Storage URL
	Description  *string   // Optional description
	CreatedAt    time.Time // Upload timestamp
}

// Evidence angle constants
const (
	EvidenceAngleFrontal = "FRONTAL"
	EvidenceAngleLateral = "LATERAL"
	EvidenceAngleRear    = "REAR"
)

// ValidEvidenceAngles contains all valid angle values
var ValidEvidenceAngles = []string{EvidenceAngleFrontal, EvidenceAngleLateral, EvidenceAngleRear}
