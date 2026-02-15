package domain

import (
	"time"

	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// MotorcycleEvidence represents a photographic evidence of a motorcycle (HU16-19)
type MotorcycleEvidence struct {
	ID           string         // UUID primary key
	MotorcycleID string         // FK to motorcycle
	Angle        *EvidenceAngle // Optional: FRONTAL, LATERAL, REAR
	ImageURL     string         // Firebase Storage URL
	Description  *string        // Optional description
	CreatedAt    time.Time      // Upload timestamp
}

func (e *MotorcycleEvidence) SetID() {
	e.ID = uuid.Generate()
}

// EvidenceAngle represents the photographic angle of motorcycle evidence
type EvidenceAngle string

const (
	EvidenceAngleFrontal EvidenceAngle = "FRONTAL"
	EvidenceAngleLateral EvidenceAngle = "LATERAL"
	EvidenceAngleRear    EvidenceAngle = "REAR"
)

// ValidEvidenceAngles contains all valid angle values
var ValidEvidenceAngles = []EvidenceAngle{EvidenceAngleFrontal, EvidenceAngleLateral, EvidenceAngleRear}
