package domain_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// DiagnosticPermission.SetID Tests
// ============================================

func TestDiagnosticPermission_SetID(t *testing.T) {
	p := &domain.DiagnosticPermission{}
	p.SetID()
	assert.NotEmpty(t, p.ID)
	assert.Len(t, p.ID, 36)
}

func TestDiagnosticPermission_SetID_Unique(t *testing.T) {
	p1 := &domain.DiagnosticPermission{}
	p2 := &domain.DiagnosticPermission{}
	p1.SetID()
	p2.SetID()
	assert.NotEqual(t, p1.ID, p2.ID)
}

// ============================================
// Franchise.SetID Tests
// ============================================

func TestFranchise_SetID(t *testing.T) {
	f := &domain.Franchise{}
	f.SetID()
	assert.NotEmpty(t, f.ID)
	assert.Len(t, f.ID, 36)
}

func TestFranchise_SetID_Unique(t *testing.T) {
	f1 := &domain.Franchise{}
	f2 := &domain.Franchise{}
	f1.SetID()
	f2.SetID()
	assert.NotEqual(t, f1.ID, f2.ID)
}

// ============================================
// Franchise.ToLogger Tests
// ============================================

func TestFranchise_ToLogger(t *testing.T) {
	f := &domain.Franchise{ID: "franchise-1", Name: "Mi Taller"}
	result := f.ToLogger()

	assert.Len(t, result, 2)
	assert.Contains(t, result[0], "id:franchise-1")
	assert.Contains(t, result[1], "name:Mi Taller")
}

// ============================================
// Motorcycle.SetID Tests
// ============================================

func TestMotorcycle_SetID(t *testing.T) {
	m := &domain.Motorcycle{}
	m.SetID()
	assert.NotEmpty(t, m.ID)
	assert.Len(t, m.ID, 36)
}

func TestMotorcycle_SetID_Unique(t *testing.T) {
	m1 := &domain.Motorcycle{}
	m2 := &domain.Motorcycle{}
	m1.SetID()
	m2.SetID()
	assert.NotEqual(t, m1.ID, m2.ID)
}

// ============================================
// MotorcycleEvidence.SetID Tests
// ============================================

func TestMotorcycleEvidence_SetID(t *testing.T) {
	e := &domain.MotorcycleEvidence{}
	e.SetID()
	assert.NotEmpty(t, e.ID)
	assert.Len(t, e.ID, 36)
}

func TestMotorcycleEvidence_SetID_Unique(t *testing.T) {
	e1 := &domain.MotorcycleEvidence{}
	e2 := &domain.MotorcycleEvidence{}
	e1.SetID()
	e2.SetID()
	assert.NotEqual(t, e1.ID, e2.ID)
}

// ============================================
// BranchSchedule.SetID Tests
// ============================================

func TestBranchSchedule_SetID(t *testing.T) {
	s := &domain.BranchSchedule{}
	s.SetID()
	assert.NotEmpty(t, s.ID)
	assert.Len(t, s.ID, 36)
}

func TestBranchSchedule_SetID_Unique(t *testing.T) {
	s1 := &domain.BranchSchedule{}
	s2 := &domain.BranchSchedule{}
	s1.SetID()
	s2.SetID()
	assert.NotEqual(t, s1.ID, s2.ID)
}

// ============================================
// ScheduleDetail.SetID Tests
// ============================================

func TestScheduleDetail_SetID(t *testing.T) {
	d := &domain.ScheduleDetail{}
	d.SetID()
	assert.NotEmpty(t, d.ID)
	assert.Len(t, d.ID, 36)
}

func TestScheduleDetail_SetID_Unique(t *testing.T) {
	d1 := &domain.ScheduleDetail{}
	d2 := &domain.ScheduleDetail{}
	d1.SetID()
	d2.SetID()
	assert.NotEqual(t, d1.ID, d2.ID)
}
