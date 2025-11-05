package dto

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
)

// RegistrationResult contiene el resultado del registro de usuario
type RegistrationResult struct {
	Person  domain.Person `json:"person"`
	Message string        `json:"message"`
}

// UserSyncStatus representa el estado de sincronización con Keycloak
type UserSyncStatus struct {
	PersonID       string `json:"person_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
}