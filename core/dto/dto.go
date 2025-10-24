package dto

import (
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/Nerzal/gocloak/v13"
)

// RegistrationResult contiene el resultado del registro con token de Keycloak
type RegistrationResult struct {
	Person domain.Person
	Token  *gocloak.JWT
}

// UserSyncStatus representa el estado de sincronización con Keycloak
type UserSyncStatus struct {
	PersonID       string `json:"person_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
}
