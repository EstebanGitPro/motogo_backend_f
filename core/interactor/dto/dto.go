package dto

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type RegistrationResult struct {
	Person  domain.Person `json:"person"`
	Message string        `json:"message"`
}

type UserSyncStatus struct {
	PersonID       string `json:"person_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
}


type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
