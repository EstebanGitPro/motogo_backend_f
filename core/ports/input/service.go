package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/Nerzal/gocloak/v13"
)

// Service - Use Cases atómicos que el Interactor orquesta (Person)
type Service interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Person - Validaciones y consultas
	RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error)
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
	GetPersonByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Person, error)
	CheckAndCleanInconsistentState(ctx context.Context, email string) error

	// Person - Operaciones transaccionales de BD
	SavePersonToDB(ctx context.Context, tx output.Tx, person domain.Person) error
	UpdatePersonKeycloakID(ctx context.Context, tx output.Tx, personID string, keycloakUserID string) error

	// Person - Operaciones de Keycloak
	CreateUserInKeycloak(ctx context.Context, person *domain.Person) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	SendVerificationEmail(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	Login(ctx context.Context, email, password string) (*gocloak.JWT, error)
	// RefreshToken obtains a new access token using the refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
	// VerifyEmailByToken receives a JWT token, extracts the email, and verifies it in Keycloak
	// Returns the extracted email on success
	VerifyEmailByToken(ctx context.Context, token string) (string, error)
	// ResetPasswordWithToken receives a JWT token, extracts the email, and updates the password in Keycloak
	// Returns nil on success
	ResetPasswordWithToken(ctx context.Context, token string, newPassword string) error
	// ChangePassword verifies the current password and updates to the new password (HU57)
	// Returns nil on success, ErrInvalidCredentials if current password is wrong
	ChangePassword(ctx context.Context, keycloakUserID, currentPassword, newPassword string) error
	// UpdatePersonProfile updates person data in DB and optionally syncs to Keycloak (HU52)
	// Returns updated person on success
	UpdatePersonProfile(ctx context.Context, tx output.Tx, person domain.Person) error

	// Person - Compensaciones (rollback)
	RollbackPerson(ctx context.Context, personID string) error
	RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error
}

// MessageService - Use Cases para Messages (separado de Person)
type MessageService interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Messages - Validaciones y operaciones
	ValidateMessage(ctx context.Context, message domain.Message) error
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessageByCode(ctx context.Context, code string) (*domain.Message, error)
	ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error)
	ListActiveMessages(ctx context.Context) ([]domain.Message, error)

	// Messages - Operaciones transaccionales de BD
	SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error
	UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error
	DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error
}

// BranchService - Use Cases for Branch operations (HU59)
type BranchService interface {
	// Transactions
	BeginTx(ctx context.Context) (output.Tx, error)

	// Branch - Validations and queries
	RegisterBranch(ctx context.Context, tx output.Tx, branch domain.Branch) (*domain.Branch, error)
	GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error)
	ValidateBrands(ctx context.Context, brands []string) error

	// Geocoding - Attempts to geocode location if coordinates not provided
	// Returns (coordsGenerated, error) - true if coordinates were generated
	GeocodeLocation(ctx context.Context, location *domain.Location) (bool, error)

	// Branch - Location operations
	SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error

	// Branch - Brands operations
	SaveBranchBrands(ctx context.Context, tx output.Tx, branchID string, brands []string) error

	// GetBranchesByRepresentative retrieves all branches for a representative (HU62)
	GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error)

	// UpdateBranch updates an existing branch (HU60)
	UpdateBranch(ctx context.Context, tx output.Tx, branch domain.Branch) error
}

// BrandService - Use Cases for Brand catalog operations
type BrandService interface {
	// GetAllBrands retrieves all brands from the catalog
	GetAllBrands(ctx context.Context) ([]domain.Brand, error)

	// ValidateBrandIDs checks if all provided brand IDs exist
	ValidateBrandIDs(ctx context.Context, brandIDs []string) error
}

// LocationService - Use Cases for geographic catalog operations
type LocationService interface {
	// GetAllDepartments retrieves all departments
	GetAllDepartments(ctx context.Context) ([]domain.Department, error)

	// GetCitiesByDepartment retrieves all cities for a specific department
	GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error)

	// ValidateCityInDepartment checks if the city belongs to the specified department
	ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error
}
