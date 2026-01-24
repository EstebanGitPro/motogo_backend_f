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

	// DeleteBranch deletes a branch by ID (HU61)
	DeleteBranch(ctx context.Context, tx output.Tx, branchID string) error
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

// FranchiseService - Use Cases for Franchise operations (HU26-29)
type FranchiseService interface {
	// Transactions
	BeginTx(ctx context.Context) (output.Tx, error)

	// Franchise CRUD
	CreateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) (*domain.Franchise, error)
	GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error)
	GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error)
	UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error
	DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error

	// Branch association
	AssociateBranches(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error
	DissociateBranches(ctx context.Context, tx output.Tx, franchiseID string) error
	DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error
	CountBranches(ctx context.Context, franchiseID string) (int, error)
}

// ServiceCatalogService - Use Cases for Service catalog operations (HU63, HU68, HU75)
type ServiceCatalogService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// GetAllServiceTypes returns all available service types (HU75)
	GetAllServiceTypes() []domain.ServiceType

	// GetAllServices retrieves all services from catalog (HU63)
	GetAllServices(ctx context.Context) ([]domain.Service, error)

	// GetServicesByType retrieves services filtered by type (HU63)
	GetServicesByType(ctx context.Context, serviceType domain.ServiceType) ([]domain.Service, error)

	// GetServiceByID retrieves a service by its UUID (HU68)
	GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error)

	// GetServicesByBranch retrieves services associated with a specific branch
	GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error)

	// AssociateBranchServices associates services to a branch
	AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error

	// DissociateBranchService removes a service from a branch
	DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error

	// ValidateServiceIDs checks if all provided service IDs exist
	ValidateServiceIDs(ctx context.Context, serviceIDs []string) error

	// CheckServiceAssociation checks if a service is already associated with a branch
	CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error)

	// UpdateService updates an existing service in the catalog (HU68 - Admin only)
	UpdateService(ctx context.Context, tx output.Tx, service domain.Service) error
}

// ScheduleService - Use Cases for Schedule operations (HU30-35)
type ScheduleService interface {
	// Transactions
	BeginTx(ctx context.Context) (output.Tx, error)

	// Schedule CRUD (HU30, HU31, HU32, HU33)
	CreateSchedule(ctx context.Context, tx output.Tx, branchID string) (*domain.BranchSchedule, error)
	GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error)
	GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error)
	UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error
	DeleteSchedule(ctx context.Context, tx output.Tx, scheduleID string) error

	// Activation (HU34, HU35)
	ActivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error
	DeactivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error
}

// ScheduleDetailService - Use Cases for Schedule Detail operations (HU6-9, HU20-25)
type ScheduleDetailService interface {
	// Transactions
	BeginTx(ctx context.Context) (output.Tx, error)

	// Schedule Detail CRUD (HU6, HU7, HU8, HU9)
	CreateDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) (*domain.ScheduleDetail, error)
	GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error)
	GetDetailsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error)
	UpdateDetail(ctx context.Context, tx output.Tx, detail domain.ScheduleDetail) error
	DeleteDetail(ctx context.Context, tx output.Tx, detailID string) error

	// Validation
	ValidateTimeRange(openingTime, closingTime string) error
	CheckTimeConflict(ctx context.Context, scheduleID string, dayOfWeek int, openingTime, closingTime string, excludeDetailID string) (bool, error)

	// Schedule Exception operations (HU20-25)
	CreateException(ctx context.Context, tx output.Tx, exception domain.ScheduleDetail) (*domain.ScheduleDetail, error)
	GetExceptionsByScheduleID(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error)
	GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error)
	UpdateException(ctx context.Context, tx output.Tx, exception domain.ScheduleDetail) error
	DeleteException(ctx context.Context, tx output.Tx, exceptionID string) error
	CheckExceptionDateConflict(ctx context.Context, scheduleID, excludeExceptionID, startDate, endDate string) (bool, error)
}
