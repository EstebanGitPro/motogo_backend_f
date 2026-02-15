package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================
// Diagnostic Tests
// ============================================

func TestDiagnostic_SetID_GeneratesUUID(t *testing.T) {
	d := &Diagnostic{}
	d.SetID()

	assert.NotEmpty(t, d.ID)
	assert.Len(t, d.ID, 36)
}

func TestDiagnostic_SetID_GeneratesUniqueIDs(t *testing.T) {
	d1 := &Diagnostic{}
	d2 := &Diagnostic{}
	d1.SetID()
	d2.SetID()

	assert.NotEqual(t, d1.ID, d2.ID)
}

func TestDiagnostic_StructAllFields(t *testing.T) {
	desc := "Engine noise"
	solution := "Replace bearing"

	d := Diagnostic{
		ID:                 "diag-001",
		MotorcycleID:       "moto-001",
		BranchID:           "branch-001",
		Date:               time.Now(),
		ProblemDescription: &desc,
		PossibleSolution:   &solution,
	}

	assert.Equal(t, "diag-001", d.ID)
	assert.Equal(t, "moto-001", d.MotorcycleID)
	assert.Equal(t, "branch-001", d.BranchID)
	assert.NotNil(t, d.ProblemDescription)
	assert.Equal(t, "Engine noise", *d.ProblemDescription)
	assert.NotNil(t, d.PossibleSolution)
	assert.Equal(t, "Replace bearing", *d.PossibleSolution)
	assert.Empty(t, d.Evidence)
}

func TestDiagnostic_NullableFieldsDefaults(t *testing.T) {
	d := Diagnostic{
		ID:           "diag-002",
		MotorcycleID: "moto-001",
		BranchID:     "branch-001",
	}

	assert.Nil(t, d.ProblemDescription)
	assert.Nil(t, d.PossibleSolution)
	assert.Nil(t, d.Evidence)
}

func TestDiagnostic_WithEvidence(t *testing.T) {
	evidence := []DiagnosticEvidence{
		{ID: "ev-001", DiagnosticID: "diag-001", ImageURL: "https://firebase.com/img1.jpg"},
		{ID: "ev-002", DiagnosticID: "diag-001", ImageURL: "https://firebase.com/img2.jpg"},
	}

	d := Diagnostic{
		ID:           "diag-001",
		MotorcycleID: "moto-001",
		BranchID:     "branch-001",
		Evidence:     evidence,
	}

	assert.Len(t, d.Evidence, 2)
	assert.Equal(t, "ev-001", d.Evidence[0].ID)
	assert.Equal(t, "ev-002", d.Evidence[1].ID)
}

// ============================================
// DiagnosticEvidence Tests
// ============================================

func TestDiagnosticEvidence_SetID_GeneratesUUID(t *testing.T) {
	e := &DiagnosticEvidence{}
	e.SetID()

	assert.NotEmpty(t, e.ID)
	assert.Len(t, e.ID, 36)
}

func TestDiagnosticEvidence_SetID_GeneratesUniqueIDs(t *testing.T) {
	e1 := &DiagnosticEvidence{}
	e2 := &DiagnosticEvidence{}
	e1.SetID()
	e2.SetID()

	assert.NotEqual(t, e1.ID, e2.ID)
}

func TestDiagnosticEvidence_StructAllFields(t *testing.T) {
	desc := "Foto lateral del daño"

	e := DiagnosticEvidence{
		ID:           "ev-001",
		DiagnosticID: "diag-001",
		ImageURL:     "https://firebasestorage.googleapis.com/v0/b/test/img.jpg",
		Description:  &desc,
		CreatedAt:    time.Now(),
	}

	assert.Equal(t, "ev-001", e.ID)
	assert.Equal(t, "diag-001", e.DiagnosticID)
	assert.Contains(t, e.ImageURL, "firebasestorage")
	assert.NotNil(t, e.Description)
	assert.Equal(t, "Foto lateral del daño", *e.Description)
}

func TestDiagnosticEvidence_NullDescription(t *testing.T) {
	e := DiagnosticEvidence{
		ID:           "ev-002",
		DiagnosticID: "diag-001",
		ImageURL:     "https://firebase.com/img.jpg",
	}

	assert.Nil(t, e.Description)
}
