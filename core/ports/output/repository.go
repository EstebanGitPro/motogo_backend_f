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

	// Location operations - transactional
	SaveLocation(ctx context.Context, tx Tx, location domain.Location) error
	UpdateLocation(ctx context.Context, tx Tx, location domain.Location) error

	// Branch brands operations - transactional
	SaveBranchBrands(ctx context.Context, tx Tx, branchID string, brands []string) error
	DeleteBranchBrands(ctx context.Context, tx Tx, branchID string) error

	// Branch operations - read
	GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error)
	GetBranchByFranchiseAndName(ctx context.Context, franchiseID, name string) (*domain.Branch, error)
	GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error)

	// Brand validation - read
	ValidateBrands(ctx context.Context, brands []string) error
}

// BrandRepository interface for Brand catalog operations
type BrandRepository interface {
	// GetAllBrands retrieves all brands from the catalog ordered by name
	GetAllBrands(ctx context.Context) ([]domain.Brand, error)

	// ValidateBrandIDs checks if all provided brand IDs exist in the brands table
	ValidateBrandIDs(ctx context.Context, brandIDs []string) error
}
