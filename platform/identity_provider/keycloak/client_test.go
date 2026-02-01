package keycloak

import (
	"context"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewClient Tests
// ============================================

func TestNewClient_NilConfig(t *testing.T) {
	client, err := NewClient(nil, nil)

	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

func TestNewClient_InvalidServerURL_FailsAdminAuth(t *testing.T) {
	cfg := &config.KeycloakConfig{
		ServerURL:    "http://localhost:9999", // Non-existent server
		Realm:        "test-realm",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		AdminUser:    "admin",
		AdminPass:    "admin",
	}
	log := logger.NewSlogLogger()

	// Should fail to connect to admin API for token
	client, err := NewClient(cfg, log)

	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize admin token")
}

// ============================================
// Client Struct Tests
// ============================================

func TestClient_TokenExpired(t *testing.T) {
	c := &client{
		tokenExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}

	// Token should be considered expired
	assert.True(t, time.Now().After(c.tokenExpiresAt))
}

func TestClient_TokenValid(t *testing.T) {
	c := &client{
		tokenExpiresAt: time.Now().Add(1 * time.Hour), // Valid for 1 more hour
	}

	// Token should be considered valid
	assert.True(t, time.Now().Before(c.tokenExpiresAt))
}

// ============================================
// Implementation-specific validation Tests
// Note: These test input validation without requiring a real Keycloak server
// ============================================

func TestClient_LoginUser_EmptyCredentials(t *testing.T) {
	// Create a minimal client struct for testing input validation
	cfg := &config.KeycloakConfig{
		ServerURL:    "http://localhost:9999",
		Realm:        "test-realm",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	log := logger.NewSlogLogger()

	c := &client{
		config: cfg,
		logger: log,
	}

	// Test empty username
	_, err := c.LoginUser(context.Background(), "", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Test empty password
	_, err = c.LoginUser(context.Background(), "user", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_CreateUser_NilPerson(t *testing.T) {
	cfg := &config.KeycloakConfig{
		ServerURL:    "http://localhost:9999",
		Realm:        "test-realm",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	log := logger.NewSlogLogger()

	c := &client{
		config: cfg,
		logger: log,
	}

	_, err := c.CreateUser(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestClient_GetUserByEmail_EmptyEmail(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	_, err := c.GetUserByEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_GetUserByID_EmptyID(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	_, err := c.GetUserByID(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_UpdateUser_NilUser(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.UpdateUser(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestClient_DeleteUser_EmptyID(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.DeleteUser(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_SetPassword_EmptyCredentials(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	// Empty userID
	err := c.SetPassword(context.Background(), "", "password", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Empty password
	err = c.SetPassword(context.Background(), "user-id", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_AssignRole_EmptyParams(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	// Empty userID
	err := c.AssignRole(context.Background(), "", "role")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Empty roleName
	err = c.AssignRole(context.Background(), "user-id", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_RemoveRole_EmptyParams(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	// Empty userID
	err := c.RemoveRole(context.Background(), "", "role")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Empty roleName
	err = c.RemoveRole(context.Background(), "user-id", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_GetUserRoles_EmptyID(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	_, err := c.GetUserRoles(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_SendVerificationEmail_EmptyID(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.SendVerificationEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_SendPasswordResetEmail_EmptyEmail(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.SendPasswordResetEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_VerifyEmail_EmptyID(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.VerifyEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_Logout_EmptyRefreshToken(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	err := c.Logout(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestClient_RefreshToken_EmptyRefreshToken(t *testing.T) {
	c := &client{
		config: &config.KeycloakConfig{},
		logger: logger.NewSlogLogger(),
	}

	_, err := c.RefreshToken(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}
