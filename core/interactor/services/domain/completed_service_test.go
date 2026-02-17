package domain_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// IsValidServiceStatus
// ============================================

func TestIsValidServiceStatus_ValidStatuses(t *testing.T) {
	validStatuses := []string{"PENDIENTE", "EN_PROCESO", "FINALIZADO", "CANCELADO"}
	for _, s := range validStatuses {
		assert.True(t, domain.IsValidServiceStatus(s), "expected %q to be valid", s)
	}
}

func TestIsValidServiceStatus_InvalidStatus(t *testing.T) {
	assert.False(t, domain.IsValidServiceStatus("INVALIDO"))
	assert.False(t, domain.IsValidServiceStatus(""))
	assert.False(t, domain.IsValidServiceStatus("pendiente")) // case-sensitive
}

// ============================================
// AllServiceStatuses
// ============================================

func TestAllServiceStatuses_ReturnsFour(t *testing.T) {
	statuses := domain.AllServiceStatuses()
	assert.Len(t, statuses, 4)
	assert.Contains(t, statuses, domain.ServiceStatusPending)
	assert.Contains(t, statuses, domain.ServiceStatusInProgress)
	assert.Contains(t, statuses, domain.ServiceStatusCompleted)
	assert.Contains(t, statuses, domain.ServiceStatusCancelled)
}

// ============================================
// SetID — CompletedService
// ============================================

func TestCompletedService_SetID(t *testing.T) {
	cs := &domain.CompletedService{}
	cs.SetID()
	assert.NotEmpty(t, cs.ID)
	assert.Len(t, cs.ID, 36) // UUID format
}

func TestCompletedService_SetID_Unique(t *testing.T) {
	cs1 := &domain.CompletedService{}
	cs2 := &domain.CompletedService{}
	cs1.SetID()
	cs2.SetID()
	assert.NotEqual(t, cs1.ID, cs2.ID)
}

// ============================================
// SetID — CompletedServiceItem
// ============================================

func TestCompletedServiceItem_SetID(t *testing.T) {
	item := &domain.CompletedServiceItem{}
	item.SetID()
	assert.NotEmpty(t, item.ID)
	assert.Len(t, item.ID, 36)
}

// ============================================
// SetID — ServiceStatusHistory
// ============================================

func TestServiceStatusHistory_SetID(t *testing.T) {
	h := &domain.ServiceStatusHistory{}
	h.SetID()
	assert.NotEmpty(t, h.ID)
	assert.Len(t, h.ID, 36)
}

// ============================================
// IsValidTransition — State Machine Matrix
// ============================================

func TestIsValidTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from domain.ServiceStatus
		to   domain.ServiceStatus
	}{
		{"PENDIENTE -> EN_PROCESO", domain.ServiceStatusPending, domain.ServiceStatusInProgress},
		{"PENDIENTE -> CANCELADO", domain.ServiceStatusPending, domain.ServiceStatusCancelled},
		{"EN_PROCESO -> FINALIZADO", domain.ServiceStatusInProgress, domain.ServiceStatusCompleted},
		{"EN_PROCESO -> CANCELADO", domain.ServiceStatusInProgress, domain.ServiceStatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, domain.IsValidTransition(tt.from, tt.to))
		})
	}
}

func TestIsValidTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from domain.ServiceStatus
		to   domain.ServiceStatus
	}{
		{"PENDIENTE -> FINALIZADO", domain.ServiceStatusPending, domain.ServiceStatusCompleted},
		{"PENDIENTE -> PENDIENTE", domain.ServiceStatusPending, domain.ServiceStatusPending},
		{"EN_PROCESO -> PENDIENTE", domain.ServiceStatusInProgress, domain.ServiceStatusPending},
		{"EN_PROCESO -> EN_PROCESO", domain.ServiceStatusInProgress, domain.ServiceStatusInProgress},
		{"FINALIZADO -> PENDIENTE", domain.ServiceStatusCompleted, domain.ServiceStatusPending},
		{"FINALIZADO -> EN_PROCESO", domain.ServiceStatusCompleted, domain.ServiceStatusInProgress},
		{"FINALIZADO -> CANCELADO", domain.ServiceStatusCompleted, domain.ServiceStatusCancelled},
		{"FINALIZADO -> FINALIZADO", domain.ServiceStatusCompleted, domain.ServiceStatusCompleted},
		{"CANCELADO -> PENDIENTE", domain.ServiceStatusCancelled, domain.ServiceStatusPending},
		{"CANCELADO -> EN_PROCESO", domain.ServiceStatusCancelled, domain.ServiceStatusInProgress},
		{"CANCELADO -> FINALIZADO", domain.ServiceStatusCancelled, domain.ServiceStatusCompleted},
		{"CANCELADO -> CANCELADO", domain.ServiceStatusCancelled, domain.ServiceStatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.False(t, domain.IsValidTransition(tt.from, tt.to))
		})
	}
}

func TestIsValidTransition_UnknownStatus(t *testing.T) {
	assert.False(t, domain.IsValidTransition("UNKNOWN", domain.ServiceStatusPending))
	assert.False(t, domain.IsValidTransition(domain.ServiceStatusPending, "UNKNOWN"))
}

// ============================================
// ServiceStatusLabels
// ============================================

func TestServiceStatusLabels_AllStatusesHaveLabels(t *testing.T) {
	statuses := domain.AllServiceStatuses()
	for _, s := range statuses {
		label, exists := domain.ServiceStatusLabels[s]
		assert.True(t, exists, "missing label for status %q", s)
		assert.NotEmpty(t, label, "label for status %q should not be empty", s)
	}
}

func TestServiceStatusLabels_CorrectValues(t *testing.T) {
	assert.Equal(t, "Pendiente", domain.ServiceStatusLabels[domain.ServiceStatusPending])
	assert.Equal(t, "En Proceso", domain.ServiceStatusLabels[domain.ServiceStatusInProgress])
	assert.Equal(t, "Finalizado", domain.ServiceStatusLabels[domain.ServiceStatusCompleted])
	assert.Equal(t, "Cancelado", domain.ServiceStatusLabels[domain.ServiceStatusCancelled])
}
