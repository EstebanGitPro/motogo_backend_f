package output

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Tx interface for database transactions
type Tx interface {
	Commit() error
	Rollback() error
}

// Repository interface for Person operations
type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Person operations - transactional
	SavePerson(ctx context.Context, tx Tx, person domain.Person) error
	UpdatePerson(ctx context.Context, tx Tx, person domain.Person) error
	PatchPerson(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeletePerson(ctx context.Context, tx Tx, id string) error

	// Person operations - read
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
	GetPersonByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Person, error)
}

// MessageRepository interface for Message operations
type MessageRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Message operations - transactional
	SaveMessage(ctx context.Context, tx Tx, message domain.Message) error
	UpdateMessage(ctx context.Context, tx Tx, message domain.Message) error
	DeleteMessage(ctx context.Context, tx Tx, id string) error

	// Message operations - read
	GetAllActive(ctx context.Context) ([]domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	GetByCode(ctx context.Context, code string) (*domain.Message, error)
	GetByType(ctx context.Context, msgType string) ([]domain.Message, error)
	GetByModule(ctx context.Context, module string) ([]domain.Message, error)
}

// BranchRepository interface for Branch operations (HU59)
type BranchRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Branch operations - transactional
	SaveBranch(ctx context.Context, tx Tx, branch domain.Branch) error
	UpdateBranch(ctx context.Context, tx Tx, branch domain.Branch) error
	DeleteBranch(ctx context.Context, tx Tx, branchID string) error

	// Branch brands operations - transactional
	SaveBranchBrands(ctx context.Context, tx Tx, branchID string, brands []string) error
	DeleteBranchBrands(ctx context.Context, tx Tx, branchID string) error

	// Branch operations - read
	GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error)
	GetBranchByFranchiseAndName(ctx context.Context, franchiseID, name string) (*domain.Branch, error)
	GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error)
	HasBranchesByRepresentative(ctx context.Context, representativeID string) (bool, error) // HU53

	// Brand validation - read
	ValidateBrands(ctx context.Context, brands []string) error

	// GetBranchesNearby retrieves branches within radius of given coordinates (HU89)
	GetBranchesNearby(ctx context.Context, lat, lng, radiusKm float64, establishmentType string, latMin, latMax, lngMin, lngMax float64) ([]domain.NearbyBranch, error)
}

// BrandRepository interface for Brand catalog operations
type BrandRepository interface {
	// GetAllBrands retrieves all brands from the catalog ordered by name
	GetAllBrands(ctx context.Context) ([]domain.Brand, error)

	// ValidateBrandIDs checks if all provided brand IDs exist in the brands table
	ValidateBrandIDs(ctx context.Context, brandIDs []string) error
}

// LocationRepository interface for geographic catalog and branch location operations
type LocationRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Location operations - transactional
	SaveLocation(ctx context.Context, tx Tx, location domain.Location) error
	UpdateLocation(ctx context.Context, tx Tx, location domain.Location) error

	// Location validation - read
	CheckAddressExists(ctx context.Context, address string) (bool, error)

	// Geographic catalog - read
	GetAllDepartments(ctx context.Context) ([]domain.Department, error)
	GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error)
	ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error
	GetDepartmentByID(ctx context.Context, departmentID string) (*domain.Department, error)
}

// FranchiseRepository interface for Franchise operations (HU26-29)
type FranchiseRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Franchise operations - transactional
	SaveFranchise(ctx context.Context, tx Tx, franchise domain.Franchise) error
	UpdateFranchise(ctx context.Context, tx Tx, franchise domain.Franchise) error
	DeleteFranchise(ctx context.Context, tx Tx, franchiseID string) error

	// Franchise operations - read
	GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error)
	GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error)
	GetFranchiseByName(ctx context.Context, name string) (*domain.Franchise, error)

	// Branch association
	AssociateBranchesToFranchise(ctx context.Context, tx Tx, franchiseID string, branchIDs []string) error
	DissociateBranchesFromFranchise(ctx context.Context, tx Tx, franchiseID string) error
	DissociateSingleBranch(ctx context.Context, tx Tx, branchID string) error
	CountBranchesByFranchise(ctx context.Context, franchiseID string) (int, error)
}

// ServiceRepository interface for Service catalog operations (HU63, HU75, HU68)
type ServiceRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// GetAllServices retrieves all services from the catalog ordered by name
	GetAllServices(ctx context.Context) ([]domain.Service, error)

	// GetServicesByType retrieves services filtered by type
	GetServicesByType(ctx context.Context, serviceType string) ([]domain.Service, error)

	// GetServiceByID retrieves a service by its UUID (HU68)
	GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error)

	// GetServicesByBranch retrieves services associated with a specific branch
	// Returns the services with their added_at timestamp from branch_services table
	GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error)

	// AssociateBranchServices associates services to a branch
	AssociateBranchServices(ctx context.Context, tx Tx, branchID string, serviceIDs []string) error

	// DissociateBranchService removes a service from a branch
	DissociateBranchService(ctx context.Context, tx Tx, branchID, serviceID string) error

	// ValidateServiceIDs checks if all provided service IDs exist in the services table
	ValidateServiceIDs(ctx context.Context, serviceIDs []string) error

	// CheckServiceAssociation checks if a service is already associated with a branch
	CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error)

	// UpdateService updates an existing service in the catalog (HU68 - Admin only)
	UpdateService(ctx context.Context, tx Tx, service domain.Service) error
}

// ScheduleRepository interface for Schedule operations (HU30-35)
type ScheduleRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Schedule operations - transactional
	SaveSchedule(ctx context.Context, tx Tx, schedule domain.BranchSchedule) error
	UpdateSchedule(ctx context.Context, tx Tx, schedule domain.BranchSchedule) error
	DeleteSchedule(ctx context.Context, tx Tx, scheduleID string) error
	SetActive(ctx context.Context, tx Tx, scheduleID string, active bool) error

	// Schedule operations - read
	GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error)
	GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error)
}

// ScheduleDetailRepository interface for Schedule Detail operations (HU6-9, HU20-25)
type ScheduleDetailRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Schedule Detail operations - transactional
	SaveScheduleDetail(ctx context.Context, tx Tx, detail domain.ScheduleDetail) error
	UpdateScheduleDetail(ctx context.Context, tx Tx, detail domain.ScheduleDetail) error
	DeleteScheduleDetail(ctx context.Context, tx Tx, detailID string) error

	// Schedule Detail operations - read
	GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error)
	GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error)
	GetDetailsByScheduleAndDay(ctx context.Context, scheduleID string, dayOfWeek int) ([]domain.ScheduleDetail, error)

	// Conflict detection
	CheckTimeConflict(ctx context.Context, scheduleID string, dayOfWeek int, openingTime, closingTime string, excludeDetailID string) (bool, error)

	// Schedule Exception operations - read (HU20-25)
	GetExceptionsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error)
	GetExceptionsByScheduleIDForUpdate(ctx context.Context, tx Tx, scheduleID string) ([]domain.ScheduleDetail, error)
	GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error)

	// Exception conflict detection (HU20)
	CheckExceptionDateConflict(ctx context.Context, scheduleID string, excludeExceptionID string, startDate string, endDate string) (bool, error)

	// Duplicate detection for REGULAR entries (Validation R1, R2, R3)
	CheckDayIsClosed(ctx context.Context, scheduleID string, dayOfWeek int, excludeDetailID string) (bool, error)
	CheckDayHasTimeSlots(ctx context.Context, scheduleID string, dayOfWeek int, excludeDetailID string) (bool, error)

	// Redundancy detection for EXCEPTION entries (Validation E1)
	CheckExceptionIsRedundant(ctx context.Context, scheduleID string, dayOfWeek int) (bool, error)
}

// MotorcycleRepository interface for Motorcycle operations (HU43-47, HU50, HU40)
type MotorcycleRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Motorcycle operations - write (HU43, HU44, HU45)
	Save(ctx context.Context, tx Tx, motorcycle *domain.Motorcycle) error
	Update(ctx context.Context, tx Tx, motorcycle *domain.Motorcycle) error
	Delete(ctx context.Context, tx Tx, motorcycleID string) error               // Soft delete
	HardDelete(ctx context.Context, tx Tx, motorcycleID string) error           // Hard delete (HU45 hybrid)
	ClearProfileImageURL(ctx context.Context, tx Tx, motorcycleID string) error // HU39

	// Motorcycle operations - read (HU46, HU47)
	GetByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error)
	GetByOwnerID(ctx context.Context, ownerID string) ([]domain.Motorcycle, error)
	GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error)

	// Reference catalog (HU50, HU40)
	GetAllReferences(ctx context.Context) ([]domain.MotorcycleReference, error)
	GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error)

	// Validation methods (HU43, HU44)
	ValidateReferenceExists(ctx context.Context, referenceID string) (bool, error)
	CheckLicensePlateExists(ctx context.Context, licensePlate string) (bool, error)

	// History check (HU45 hybrid delete)
	HasServiceHistory(ctx context.Context, motorcycleID string) (bool, error)
}

// EvidenceRepository interface for Motorcycle Evidence operations (HU16-19)
type EvidenceRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Evidence operations - write (HU16, HU17, HU19)
	Save(ctx context.Context, tx Tx, evidence *domain.MotorcycleEvidence) error
	Update(ctx context.Context, tx Tx, evidence *domain.MotorcycleEvidence) error
	Delete(ctx context.Context, tx Tx, evidenceID string) error

	// Evidence operations - read (HU18)
	GetByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error)
	GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error)

	// Validation methods (HU16)
	CountByMotorcycleID(ctx context.Context, motorcycleID string) (int, error)
}

// DiagnosticRepository interface for Diagnostic operations (HU11-14)
type DiagnosticRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Diagnostic operations - write (HU11, HU12, HU13)
	Save(ctx context.Context, tx Tx, diagnostic *domain.Diagnostic) error
	Update(ctx context.Context, tx Tx, diagnostic *domain.Diagnostic) error
	Delete(ctx context.Context, tx Tx, diagnosticID string) error

	// Diagnostic Evidence operations - write (HU11)
	SaveEvidence(ctx context.Context, tx Tx, evidence *domain.DiagnosticEvidence) error

	// Diagnostic operations - read (HU14)
	GetByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error)
	GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error)
	GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.Diagnostic, error)

	// Diagnostic Evidence operations - read
	GetEvidenceByDiagnosticID(ctx context.Context, diagnosticID string) ([]domain.DiagnosticEvidence, error)

	// Diagnostic Evidence operations - write (cleanup for UPSERT)
	DeleteEvidenceByDiagnosticID(ctx context.Context, tx Tx, diagnosticID string) error
}

// DiagnosticPermissionRepository interface for per-branch diagnostic permission operations
type DiagnosticPermissionRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Permission operations - write
	Save(ctx context.Context, tx Tx, permission *domain.DiagnosticPermission) error
	Deactivate(ctx context.Context, tx Tx, motorcycleID, branchID string) error

	// Permission operations - read
	GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.DiagnosticPermission, error)
	GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error)
}
