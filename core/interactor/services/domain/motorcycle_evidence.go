package domain

import (
	"time"

	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// MotorcycleEvidence represents a photographic evidence of a motorcycle (HU16-19)
type MotorcycleEvidence struct {
	ID           string    // UUID primary key
	MotorcycleID string    // FK to motorcycle
	Angle        *string   // Optional: FRONTAL, LATERAL, REAR
	ImageURL     string    // Firebase Storage URL
	Description  *string   // Optional description
	CreatedAt    time.Time // Upload timestamp
}

func (e *MotorcycleEvidence) SetID() {
	e.ID = uuid.Generate()
}

// Evidence angle constants
const (
	EvidenceAngleFrontal = "FRONTAL"
	EvidenceAngleLateral = "LATERAL"
	EvidenceAngleRear    = "REAR"
)

// ValidEvidenceAngles contains all valid angle values
var ValidEvidenceAngles = []string{EvidenceAngleFrontal, EvidenceAngleLateral, EvidenceAngleRear}
