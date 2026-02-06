package domain

import "time"

// MotorcycleEvidence represents a photographic evidence of a motorcycle (HU16-19)
type MotorcycleEvidence struct {
	ID           string    // UUID primary key
	MotorcycleID string    // FK to motorcycle
	Angle        *string   // Optional: FRONT, SIDE, BACK
	ImageURL     string    // Firebase Storage URL
	UploadDate   time.Time // Upload timestamp
}

// Evidence angle constants
const (
	EvidenceAngleFront = "FRONT"
	EvidenceAngleSide  = "SIDE"
	EvidenceAngleBack  = "BACK"
)

// ValidEvidenceAngles contains all valid angle values
var ValidEvidenceAngles = []string{EvidenceAngleFront, EvidenceAngleSide, EvidenceAngleBack}
