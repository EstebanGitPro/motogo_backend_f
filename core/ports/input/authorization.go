package input

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/Nerzal/gocloak/v13"
)

// AuthorizationService maneja la autorización y roles usando Keycloak
// Se integra con el flujo existente de registro de personas
type AuthorizationService interface {
	// Autenticación
	LoginUser(ctx context.Context, email, password string) (*gocloak.JWT, error)

	// Sincronización con Keycloak (todos los usuarios para autenticación)
	SyncUserToKeycloak(ctx context.Context, person *domain.Person) (string, error) // Retorna Keycloak UserID
	DeleteUserFromKeycloak(ctx context.Context, keycloakUserID string) error       // Para rollback

	// Configuración de contraseña en Keycloak para autenticación
	SetUserPassword(ctx context.Context, keycloakUserID string, password string) error

	// Gestión de Roles
	AssignRole(ctx context.Context, personID string, roleName string) error
	RemoveRole(ctx context.Context, personID string, roleName string) error
	GetUserRoles(ctx context.Context, personID string) ([]string, error)

	// Validación de permisos
	HasRole(ctx context.Context, personID string, roleName string) (bool, error)
	HasPermission(ctx context.Context, personID string, resource, action string) (bool, error)

	// Gestión de roles en Keycloak
	CreateRole(ctx context.Context, roleName, description string) error
	GetAllRoles(ctx context.Context) ([]string, error)
}
