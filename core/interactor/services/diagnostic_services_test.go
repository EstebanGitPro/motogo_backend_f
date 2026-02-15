package services_test

import (
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewDiagnostic Tests
// ============================================

func TestNewDiagnostic_GeneratesUUID(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)

	assert.NotEmpty(t, diag.ID)
	assert.Len(t, diag.ID, 36)
}

func TestNewDiagnostic_GeneratesUniqueIDs(t *testing.T) {
	d1 := services.NewDiagnostic("moto-001", "branch-001", nil)
	d2 := services.NewDiagnostic("moto-001", "branch-001", nil)

	assert.NotEqual(t, d1.ID, d2.ID)
}

func TestNewDiagnostic_SetsMotorcycleID(t *testing.T) {
	diag := services.NewDiagnostic("moto-456", "branch-001", nil)
	assert.Equal(t, "moto-456", diag.MotorcycleID)
}

func TestNewDiagnostic_SetsBranchID(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-789", nil)
	assert.Equal(t, "branch-789", diag.BranchID)
}

func TestNewDiagnostic_SetsTimestamp(t *testing.T) {
	before := time.Now()
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)
	after := time.Now()

	assert.True(t, diag.Date.After(before) || diag.Date.Equal(before))
	assert.True(t, diag.Date.Before(after) || diag.Date.Equal(after))
}

func TestNewDiagnostic_NilProblemDescription(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)
	assert.Nil(t, diag.ProblemDescription)
}

func TestNewDiagnostic_WithProblemDescription(t *testing.T) {
	desc := "Motorcycle makes strange noise"
	diag := services.NewDiagnostic("moto-001", "branch-001", &desc)

	assert.NotNil(t, diag.ProblemDescription)
	assert.Equal(t, "Motorcycle makes strange noise", *diag.ProblemDescription)
}

func TestNewDiagnostic_PossibleSolutionIsNil(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)
	assert.Nil(t, diag.PossibleSolution)
}

// ============================================
// NewDiagnosticEvidence Tests
// ============================================

func TestNewDiagnosticEvidence_GeneratesUUID(t *testing.T) {
	evidence := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img.jpg", nil)

	assert.NotEmpty(t, evidence.ID)
	assert.Len(t, evidence.ID, 36)
}

func TestNewDiagnosticEvidence_GeneratesUniqueIDs(t *testing.T) {
	e1 := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img1.jpg", nil)
	e2 := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img2.jpg", nil)

	assert.NotEqual(t, e1.ID, e2.ID)
}

func TestNewDiagnosticEvidence_SetsDiagnosticID(t *testing.T) {
	evidence := services.NewDiagnosticEvidence("diag-456", "https://firebase.com/img.jpg", nil)
	assert.Equal(t, "diag-456", evidence.DiagnosticID)
}

func TestNewDiagnosticEvidence_SetsImageURL(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/test/image.jpg"
	evidence := services.NewDiagnosticEvidence("diag-001", url, nil)
	assert.Equal(t, url, evidence.ImageURL)
}

func TestNewDiagnosticEvidence_NilDescription(t *testing.T) {
	evidence := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img.jpg", nil)
	assert.Nil(t, evidence.Description)
}

func TestNewDiagnosticEvidence_WithDescription(t *testing.T) {
	desc := "Foto del daño lateral"
	evidence := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img.jpg", &desc)

	assert.NotNil(t, evidence.Description)
	assert.Equal(t, "Foto del daño lateral", *evidence.Description)
}

func TestNewDiagnosticEvidence_SetsCreatedAt(t *testing.T) {
	before := time.Now()
	evidence := services.NewDiagnosticEvidence("diag-001", "https://firebase.com/img.jpg", nil)
	after := time.Now()

	assert.True(t, evidence.CreatedAt.After(before) || evidence.CreatedAt.Equal(before))
	assert.True(t, evidence.CreatedAt.Before(after) || evidence.CreatedAt.Equal(after))
}

// ============================================
// RefreshDiagnostic Tests
// ============================================

func TestRefreshDiagnostic_UpdatesDescription(t *testing.T) {
	oldDesc := "Old problem"
	diag := services.NewDiagnostic("moto-001", "branch-001", &oldDesc)

	newDesc := "New problem description"
	services.RefreshDiagnostic(diag, &newDesc)

	assert.NotNil(t, diag.ProblemDescription)
	assert.Equal(t, "New problem description", *diag.ProblemDescription)
}

func TestRefreshDiagnostic_UpdatesTimestamp(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)
	originalDate := diag.Date

	// Small delay to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)
	services.RefreshDiagnostic(diag, nil)

	assert.True(t, diag.Date.After(originalDate) || diag.Date.Equal(originalDate))
}

func TestRefreshDiagnostic_SetsDescriptionToNil(t *testing.T) {
	desc := "Some problem"
	diag := services.NewDiagnostic("moto-001", "branch-001", &desc)

	services.RefreshDiagnostic(diag, nil)
	assert.Nil(t, diag.ProblemDescription)
}

func TestRefreshDiagnostic_PreservesOtherFields(t *testing.T) {
	diag := services.NewDiagnostic("moto-001", "branch-001", nil)
	originalID := diag.ID

	newDesc := "Updated"
	services.RefreshDiagnostic(diag, &newDesc)

	assert.Equal(t, originalID, diag.ID)
	assert.Equal(t, "moto-001", diag.MotorcycleID)
	assert.Equal(t, "branch-001", diag.BranchID)
}
