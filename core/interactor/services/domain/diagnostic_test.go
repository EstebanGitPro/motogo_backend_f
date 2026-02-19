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
}

func TestDiagnostic_NullableFieldsDefaults(t *testing.T) {
	d := Diagnostic{
		ID:           "diag-002",
		MotorcycleID: "moto-001",
		BranchID:     "branch-001",
	}

	assert.Nil(t, d.ProblemDescription)
	assert.Nil(t, d.PossibleSolution)
}
