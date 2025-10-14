package ports

import (
	"context"
	"github.com/EstebanGitPro/motogo-backend/core/domain"
)

// AuthorizationService maneja la autorización y roles usando Keycloak
// Se integra con el flujo existente de registro de personas
type AuthorizationService interface {
	// Sincronización con Keycloak (solo cuando sea necesario para roles)
	SyncUserToKeycloak(ctx context.Context, person *domain.Person) (string, error) // Retorna Keycloak UserID
	
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

// UserSyncStatus representa el estado de sincronización con Keycloak
type UserSyncStatus struct {
	PersonID       string `json:"person_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
}
