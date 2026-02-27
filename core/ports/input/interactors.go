package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// ============================================
// Interactor Interfaces for Handler Layer
// These interfaces allow dependency injection and testing
// ============================================

// BrandLister defines the contract for brand catalog operations
type BrandLister interface {
	GetAllBrands(ctx context.Context) ([]domain.Brand, error)
}

// LocationInteractorInterface defines the contract for geographic catalog operations
type LocationInteractorInterface interface {
	GetAllDepartments(ctx context.Context) ([]domain.Department, error)
	GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error)
}

// MotorcycleInteractorInterface defines the contract for motorcycle operations
// Note: Signatures match the actual MotorcycleInteractor implementation
type MotorcycleInteractorInterface interface {
	RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error)
	GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error)
	GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error)
	GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Motorcycle, error)
	UpdateMotorcycle(ctx context.Context, motorcycleID string, ownerID string, updates *domain.Motorcycle) (*domain.Motorcycle, error)
	DeleteMotorcycle(ctx context.Context, motorcycleID string, ownerID string) error
	DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) error // HU39
	GetMotorcycleReferences(ctx context.Context) ([]domain.MotorcycleReference, error)
	GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error)
	GetDistinctCategories(ctx context.Context) ([]domain.MotorcycleCategory, error)
	GetLinesByCategory(ctx context.Context, categoryName string) ([]domain.CategoryLine, error)
	GetDistinctDisplacements(ctx context.Context) ([]domain.EngineDisplacementRange, error)
	GetRatingRanges(ctx context.Context) ([]domain.RatingRange, error)
	GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string, active bool) (*domain.DiagnosticPermission, error)
	RevokeDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) error
	ListDiagnosticPermissions(ctx context.Context, motorcycleID, ownerID string) ([]domain.DiagnosticPermission, error)
	LookupPermissions(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error)
}

// EvidenceInteractorInterface defines the contract for motorcycle evidence operations (HU16-19)
type EvidenceInteractorInterface interface {
	CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error)
	GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error)
	ListEvidenceByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.MotorcycleEvidence, error)
	UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) // HU17
	DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error
	LookupEvidence(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) // Workshop lookup (no ownership check)
}

// CompletedServiceInteractorInterface defines the contract for completed service operations (HU64)
type CompletedServiceInteractorInterface interface {
	RegisterCompletedService(ctx context.Context, service *domain.CompletedService, serviceIDs []string, personID string) (*domain.CompletedService, error)
	GetCompletedServiceByID(ctx context.Context, serviceID string) (*domain.CompletedService, error)
	GetCompletedServicesByBranch(ctx context.Context, branchID string) ([]domain.CompletedService, error)
	GetCompletedServicesByMotorcycle(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error)
	TransitionStatus(ctx context.Context, serviceID string, newStatus string, personID string, finalPrice *float64) error
	UpdateCompletedServiceDetails(ctx context.Context, serviceID string, quotedPrice, finalPrice *float64, notes *string) error
}
