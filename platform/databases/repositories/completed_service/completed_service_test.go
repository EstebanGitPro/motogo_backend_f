package completed_service_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	cs "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/completed_service"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ToDomain Tests
// ============================================

func TestToDomain_AllFieldsPopulated(t *testing.T) {
	now := time.Now()
	diagID := "diag-1"
	notes := "test notes"
	branchName := "Sucursal Centro"
	quotedPrice := 150000.0
	finalPrice := 145000.0

	dbModel := &cs.CompletedService{
		ID:                  "cs-1",
		BranchID:            "branch-1",
		BranchName:          sql.NullString{String: branchName, Valid: true},
		MotorcycleID:        "moto-1",
		DiagnosticID:        sql.NullString{String: diagID, Valid: true},
		RequestDate:         now,
		CompletionDate:      sql.NullTime{Time: now, Valid: true},
		Status:              "EN_PROCESO",
		QuotedPrice:         sql.NullFloat64{Float64: quotedPrice, Valid: true},
		FinalPrice:          sql.NullFloat64{Float64: finalPrice, Valid: true},
		RepresentativeNotes: sql.NullString{String: notes, Valid: true},
		DeletedAt:           sql.NullTime{Time: now, Valid: true},
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result := dbModel.ToDomain()

	assert.Equal(t, "cs-1", result.ID)
	assert.Equal(t, "branch-1", result.BranchID)
	assert.NotNil(t, result.BranchName)
	assert.Equal(t, branchName, *result.BranchName)
	assert.Equal(t, "moto-1", result.MotorcycleID)
	assert.NotNil(t, result.DiagnosticID)
	assert.Equal(t, diagID, *result.DiagnosticID)
	assert.NotNil(t, result.CompletionDate)
	assert.Equal(t, domain.ServiceStatusInProgress, result.Status)
	assert.NotNil(t, result.QuotedPrice)
	assert.Equal(t, quotedPrice, *result.QuotedPrice)
	assert.NotNil(t, result.FinalPrice)
	assert.Equal(t, finalPrice, *result.FinalPrice)
	assert.NotNil(t, result.RepresentativeNotes)
	assert.Equal(t, notes, *result.RepresentativeNotes)
	assert.NotNil(t, result.DeletedAt)
}

func TestToDomain_NullableFieldsEmpty(t *testing.T) {
	now := time.Now()

	dbModel := &cs.CompletedService{
		ID:           "cs-1",
		BranchID:     "branch-1",
		MotorcycleID: "moto-1",
		RequestDate:  now,
		Status:       "PENDIENTE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := dbModel.ToDomain()

	assert.Equal(t, "cs-1", result.ID)
	assert.Nil(t, result.BranchName)
	assert.Nil(t, result.DiagnosticID)
	assert.Nil(t, result.CompletionDate)
	assert.Nil(t, result.QuotedPrice)
	assert.Nil(t, result.FinalPrice)
	assert.Nil(t, result.RepresentativeNotes)
	assert.Nil(t, result.DeletedAt)
	assert.Equal(t, domain.ServiceStatusPending, result.Status)
}

// ============================================
// FromDomain Tests
// ============================================

func TestFromDomain_AllFieldsPopulated(t *testing.T) {
	now := time.Now()
	diagID := "diag-1"
	notes := "test notes"
	quotedPrice := 150000.0
	finalPrice := 145000.0

	domainCS := &domain.CompletedService{
		ID:                  "cs-1",
		BranchID:            "branch-1",
		MotorcycleID:        "moto-1",
		DiagnosticID:        &diagID,
		RequestDate:         now,
		CompletionDate:      &now,
		Status:              domain.ServiceStatusCompleted,
		QuotedPrice:         &quotedPrice,
		FinalPrice:          &finalPrice,
		RepresentativeNotes: &notes,
	}

	dbModel := cs.FromDomain(domainCS)

	assert.Equal(t, "cs-1", dbModel.ID)
	assert.Equal(t, "branch-1", dbModel.BranchID)
	assert.Equal(t, "moto-1", dbModel.MotorcycleID)
	assert.True(t, dbModel.DiagnosticID.Valid)
	assert.Equal(t, diagID, dbModel.DiagnosticID.String)
	assert.True(t, dbModel.CompletionDate.Valid)
	assert.Equal(t, "FINALIZADO", dbModel.Status)
	assert.True(t, dbModel.QuotedPrice.Valid)
	assert.Equal(t, quotedPrice, dbModel.QuotedPrice.Float64)
	assert.True(t, dbModel.FinalPrice.Valid)
	assert.Equal(t, finalPrice, dbModel.FinalPrice.Float64)
	assert.True(t, dbModel.RepresentativeNotes.Valid)
	assert.Equal(t, notes, dbModel.RepresentativeNotes.String)
}

func TestFromDomain_NullableFieldsNil(t *testing.T) {
	now := time.Now()

	domainCS := &domain.CompletedService{
		ID:           "cs-1",
		BranchID:     "branch-1",
		MotorcycleID: "moto-1",
		RequestDate:  now,
		Status:       domain.ServiceStatusPending,
	}

	dbModel := cs.FromDomain(domainCS)

	assert.False(t, dbModel.DiagnosticID.Valid)
	assert.False(t, dbModel.CompletionDate.Valid)
	assert.False(t, dbModel.QuotedPrice.Valid)
	assert.False(t, dbModel.FinalPrice.Valid)
	assert.False(t, dbModel.RepresentativeNotes.Valid)
}

// ============================================
// ItemFromDomain Tests
// ============================================

func TestItemFromDomain(t *testing.T) {
	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		ServiceID:          "svc-1",
	}

	dbItem := cs.ItemFromDomain(item)

	assert.Equal(t, "item-1", dbItem.ID)
	assert.Equal(t, "cs-1", dbItem.CompletedServiceID)
	assert.Equal(t, "svc-1", dbItem.ServiceID)
}

// ============================================
// ItemToDomain Tests
// ============================================

func TestItemToDomain_AllFieldsPopulated(t *testing.T) {
	now := time.Now()

	dbItem := &cs.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		ServiceID:          "svc-1",
		Rating:             sql.NullInt32{Int32: 5, Valid: true},
		Comment:            sql.NullString{String: "excelente", Valid: true},
		RatedAt:            sql.NullTime{Time: now, Valid: true},
		IsOffensiveComment: false,
	}

	result := dbItem.ItemToDomain()

	assert.Equal(t, "item-1", result.ID)
	assert.Equal(t, "cs-1", result.CompletedServiceID)
	assert.Equal(t, "svc-1", result.ServiceID)
	assert.NotNil(t, result.Rating)
	assert.Equal(t, 5, *result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, "excelente", *result.Comment)
	assert.NotNil(t, result.RatedAt)
	assert.False(t, result.IsOffensiveComment)
}

func TestItemToDomain_NullableFieldsEmpty(t *testing.T) {
	dbItem := &cs.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		ServiceID:          "svc-1",
	}

	result := dbItem.ItemToDomain()

	assert.Nil(t, result.Rating)
	assert.Nil(t, result.Comment)
	assert.Nil(t, result.RatedAt)
}

// ============================================
// HistoryFromDomain Tests
// ============================================

func TestHistoryFromDomain_WithPreviousStatus(t *testing.T) {
	prev := domain.ServiceStatusPending
	h := &domain.ServiceStatusHistory{
		ID:                 "h-1",
		CompletedServiceID: "cs-1",
		PreviousStatus:     &prev,
		NewStatus:          domain.ServiceStatusInProgress,
		CreatedBy:          "person-1",
	}

	dbH := cs.HistoryFromDomain(h)

	assert.Equal(t, "h-1", dbH.ID)
	assert.Equal(t, "cs-1", dbH.CompletedServiceID)
	assert.True(t, dbH.PreviousStatus.Valid)
	assert.Equal(t, "PENDIENTE", dbH.PreviousStatus.String)
	assert.Equal(t, "EN_PROCESO", dbH.NewStatus)
	assert.Equal(t, "person-1", dbH.CreatedBy)
}

func TestHistoryFromDomain_NoPreviousStatus(t *testing.T) {
	h := &domain.ServiceStatusHistory{
		ID:                 "h-1",
		CompletedServiceID: "cs-1",
		NewStatus:          domain.ServiceStatusPending,
		CreatedBy:          "person-1",
	}

	dbH := cs.HistoryFromDomain(h)

	assert.False(t, dbH.PreviousStatus.Valid)
	assert.Equal(t, "PENDIENTE", dbH.NewStatus)
}

// ============================================
// HistoryToDomain Tests
// ============================================

func TestHistoryToDomain_WithPreviousStatus(t *testing.T) {
	now := time.Now()

	dbH := &cs.ServiceStatusHistory{
		ID:                 "h-1",
		CompletedServiceID: "cs-1",
		PreviousStatus:     sql.NullString{String: "PENDIENTE", Valid: true},
		NewStatus:          "EN_PROCESO",
		CreatedBy:          "person-1",
		CreatedAt:          now,
	}

	result := dbH.HistoryToDomain()

	assert.Equal(t, "h-1", result.ID)
	assert.Equal(t, "cs-1", result.CompletedServiceID)
	assert.NotNil(t, result.PreviousStatus)
	assert.Equal(t, domain.ServiceStatusPending, *result.PreviousStatus)
	assert.Equal(t, domain.ServiceStatusInProgress, result.NewStatus)
	assert.Equal(t, "person-1", result.CreatedBy)
	assert.Equal(t, now, result.CreatedAt)
}

func TestHistoryToDomain_NoPreviousStatus(t *testing.T) {
	now := time.Now()

	dbH := &cs.ServiceStatusHistory{
		ID:                 "h-1",
		CompletedServiceID: "cs-1",
		NewStatus:          "PENDIENTE",
		CreatedBy:          "person-1",
		CreatedAt:          now,
	}

	result := dbH.HistoryToDomain()

	assert.Nil(t, result.PreviousStatus)
	assert.Equal(t, domain.ServiceStatusPending, result.NewStatus)
}
