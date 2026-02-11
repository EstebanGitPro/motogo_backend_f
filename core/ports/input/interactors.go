package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// ============================================
// Interactor Interfaces for Handler Layer
// These interfaces allow dependency injection and testing
// ============================================

// BrandInteractorInterface defines the contract for brand catalog operations
type BrandInteractorInterface interface {
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
	GetMotorcycleByLicensePlate(ctx context.Context, licensePlate string, branchIDs []string) (*domain.Motorcycle, error)
	UpdateMotorcycle(ctx context.Context, motorcycleID string, ownerID string, updates *domain.Motorcycle) (*domain.Motorcycle, error)
	DeleteMotorcycle(ctx context.Context, motorcycleID string, ownerID string) error
	DeleteProfileImage(ctx context.Context, motorcycleID string, ownerID string) error // HU39
	GetMotorcycleReferences(ctx context.Context) ([]domain.MotorcycleReference, error)
	GetReferencesByBrandID(ctx context.Context, brandID string) ([]domain.MotorcycleReference, error)
	// Diagnostic Permissions
	GrantDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string, active bool) (*domain.DiagnosticPermission, error)
	RevokeDiagnosticPermission(ctx context.Context, motorcycleID, branchID, ownerID string) error
	ListDiagnosticPermissions(ctx context.Context, motorcycleID, ownerID string) ([]domain.DiagnosticPermission, error)
}

// EvidenceInteractorInterface defines the contract for motorcycle evidence operations (HU16-19)
type EvidenceInteractorInterface interface {
	CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error)
	GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error)
	ListEvidenceByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.MotorcycleEvidence, error)
	UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) // HU17
	DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error
	ListEvidenceByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error)
}
