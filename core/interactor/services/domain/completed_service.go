package domain

import (
	"time"

	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// ServiceStatus represents the status of a completed service (HU64)
type ServiceStatus string

const (
	ServiceStatusPending    ServiceStatus = "PENDIENTE"
	ServiceStatusInProgress ServiceStatus = "EN_PROCESO"
	ServiceStatusCompleted  ServiceStatus = "FINALIZADO"
	ServiceStatusCancelled  ServiceStatus = "CANCELADO"
)

// ServiceStatusLabels maps status values to Spanish display labels
var ServiceStatusLabels = map[ServiceStatus]string{
	ServiceStatusPending:    "Pendiente",
	ServiceStatusInProgress: "En Proceso",
	ServiceStatusCompleted:  "Finalizado",
	ServiceStatusCancelled:  "Cancelado",
}

// AllServiceStatuses returns all valid service statuses
func AllServiceStatuses() []ServiceStatus {
	return []ServiceStatus{
		ServiceStatusPending,
		ServiceStatusInProgress,
		ServiceStatusCompleted,
		ServiceStatusCancelled,
	}
}

// IsValidServiceStatus checks if a status string is valid
func IsValidServiceStatus(s string) bool {
	switch ServiceStatus(s) {
	case ServiceStatusPending, ServiceStatusInProgress,
		ServiceStatusCompleted, ServiceStatusCancelled:
		return true
	}
	return false
}

// IsValidTransition checks if a status transition is allowed (HU74)
func IsValidTransition(from, to ServiceStatus) bool {
	switch from {
	case ServiceStatusPending:
		return to == ServiceStatusInProgress || to == ServiceStatusCancelled
	case ServiceStatusInProgress:
		return to == ServiceStatusCompleted || to == ServiceStatusCancelled
	default:
		// FINALIZADO and CANCELADO are terminal states
		return false
	}
}

// CompletedService represents a service visit at a branch (HU64)
type CompletedService struct {
	ID                  string                 // UUID primary key
	BranchID            string                 // FK to branch (from representative's token)
	BranchName          *string                // Branch name (populated via JOIN for listing endpoints)
	MotorcycleID        string                 // FK to motorcycle
	DiagnosticID        *string                // Optional FK to diagnostic
	RequestDate         time.Time              // Service request date
	CompletionDate      *time.Time             // Completion date (set when status=FINALIZADO)
	Status              ServiceStatus          // PENDIENTE, EN_PROCESO, FINALIZADO, CANCELADO
	QuotedPrice         *float64               // Quoted price before service
	FinalPrice          *float64               // Final price charged
	RepresentativeNotes *string                // Notes from representative
	Services            []CompletedServiceItem // Associated services (from pivot table)
	DeletedAt           *time.Time             // Soft delete timestamp (HU65 hybrid)
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// SetID generates a UUID for the completed service
func (cs *CompletedService) SetID() {
	cs.ID = uuid.Generate()
}

// CompletedServiceItem represents a service associated with a completed service (pivot)
type CompletedServiceItem struct {
	ID                 string     // UUID primary key
	CompletedServiceID string     // FK to completed_services
	ServiceID          string     // FK to services
	Rating             *int       // Rating 1-5 (NULL until rated by motorcyclist)
	Comment            *string    // Optional customer review text
	RatedAt            *time.Time // When rating was submitted
	IsOffensiveComment bool       // Flag for offensive comments
}

// SetID generates a UUID for the item
func (csi *CompletedServiceItem) SetID() {
	csi.ID = uuid.Generate()
}

// ServiceStatusHistory tracks status changes for audit trail (HU73)
type ServiceStatusHistory struct {
	ID                 string         // UUID primary key
	CompletedServiceID string         // FK to completed_services
	PreviousStatus     *ServiceStatus // Previous status (nil on creation)
	NewStatus          ServiceStatus  // New status
	CreatedBy          string         // FK to person who made the change
	CreatedAt          time.Time
}

// SetID generates a UUID for the history entry
func (h *ServiceStatusHistory) SetID() {
	h.ID = uuid.Generate()
}
