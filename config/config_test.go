package config_test

import (
	"os"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/stretchr/testify/assert"
)

// ============================================
// LoadConfig Tests
// ============================================

func TestLoadConfig_DefaultLocal_Success(t *testing.T) {
	// Ensure APP_ENV is empty so it defaults to "local"
	t.Setenv("APP_ENV", "")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "local", cfg.Environment)
	assert.Equal(t, "mysql", cfg.Database.Driver)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.NotEmpty(t, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Keycloak.ServerURL)
	assert.NotEmpty(t, cfg.Keycloak.Realm)
}

func TestLoadConfig_RailwayEnv_FallsBackToLocal(t *testing.T) {
	// Set APP_ENV to "railway" — railway-config.json doesn't exist,
	// so it should fall back to local-config.json
	t.Setenv("APP_ENV", "railway")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// Even though env was "railway", the file loaded is local-config.json
	// which has environment: "local"
	assert.Equal(t, "local", cfg.Environment)
}

func TestLoadConfig_KeycloakEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")

	// Set Keycloak env overrides
	t.Setenv("KEYCLOAK_SERVER_URL", "https://kc.override.com")
	t.Setenv("KEYCLOAK_REALM", "override-realm")
	t.Setenv("KEYCLOAK_CLIENT_ID", "override-client")
	t.Setenv("KEYCLOAK_CLIENT_SECRET", "override-secret")
	t.Setenv("KEYCLOAK_ADMIN", "override-admin")
	t.Setenv("KEYCLOAK_ADMIN_PASSWORD", "override-pass")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "https://kc.override.com", cfg.Keycloak.ServerURL)
	assert.Equal(t, "override-realm", cfg.Keycloak.Realm)
	assert.Equal(t, "override-client", cfg.Keycloak.ClientID)
	assert.Equal(t, "override-secret", cfg.Keycloak.ClientSecret)
	assert.Equal(t, "override-admin", cfg.Keycloak.AdminUser)
	assert.Equal(t, "override-pass", cfg.Keycloak.AdminPass)
}

func TestLoadConfig_PartialKeycloakOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")

	// Only override ServerURL — others should keep JSON values
	t.Setenv("KEYCLOAK_SERVER_URL", "https://partial.override.com")
	// Ensure other overrides are cleared
	os.Unsetenv("KEYCLOAK_REALM")
	os.Unsetenv("KEYCLOAK_CLIENT_ID")
	os.Unsetenv("KEYCLOAK_CLIENT_SECRET")
	os.Unsetenv("KEYCLOAK_ADMIN")
	os.Unsetenv("KEYCLOAK_ADMIN_PASSWORD")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "https://partial.override.com", cfg.Keycloak.ServerURL)
	// These should retain the values from local-config.json
	assert.Equal(t, "motogo", cfg.Keycloak.Realm)
	assert.Equal(t, "stifler", cfg.Keycloak.ClientID)
}

func TestGetMySQLDSN_WithSSL(t *testing.T) {
	c := &config.Config{
		Database: config.Database{
			Username: "user",
			Password: "pass",
			Host:     "localhost",
			Port:     "3306",
			Name:     "testdb",
			SSL:      "true",
		},
	}

	dsn := c.GetMySQLDSN()

	assert.Contains(t, dsn, "user:pass@tcp(localhost:3306)/testdb")
	assert.Contains(t, dsn, "&tls=true")
}

func TestGetMySQLDSN_WithoutSSL(t *testing.T) {
	c := &config.Config{
		Database: config.Database{
			Username: "user",
			Password: "pass",
			Host:     "localhost",
			Port:     "3306",
			Name:     "testdb",
		},
	}

	dsn := c.GetMySQLDSN()

	assert.Contains(t, dsn, "user:pass@tcp(localhost:3306)/testdb")
	assert.NotContains(t, dsn, "tls=")
}

func TestGetServerAddress(t *testing.T) {
	c := &config.Config{
		Server: config.Server{Host: "0.0.0.0", Port: "8080"},
	}

	assert.Equal(t, "0.0.0.0:8080", c.GetServerAddress())
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"production", true},
		{"railway", true},
		{"local", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			c := &config.Config{Environment: tt.env}
			assert.Equal(t, tt.expected, c.IsProduction())
		})
	}
}

func TestGetKeycloakAuthURL(t *testing.T) {
	c := &config.Config{
		Keycloak: config.KeycloakConfig{
			ServerURL: "https://keycloak.example.com",
			Realm:     "motogo",
		},
	}

	expected := "https://keycloak.example.com/realms/motogo/protocol/openid-connect/token"
	assert.Equal(t, expected, c.GetKeycloakAuthURL())
}

func TestGetKeycloakAdminURL(t *testing.T) {
	c := &config.Config{
		Keycloak: config.KeycloakConfig{
			ServerURL: "https://keycloak.example.com",
			Realm:     "motogo",
		},
	}

	expected := "https://keycloak.example.com/admin/realms/motogo"
	assert.Equal(t, expected, c.GetKeycloakAdminURL())
}

func TestGetKeycloakJWKSURL(t *testing.T) {
	c := &config.Config{
		Keycloak: config.KeycloakConfig{
			ServerURL: "https://keycloak.example.com",
			Realm:     "motogo",
		},
	}

	expected := "https://keycloak.example.com/realms/motogo/protocol/openid-connect/certs"
	assert.Equal(t, expected, c.GetKeycloakJWKSURL())
}

func TestGetKeycloakIssuerURL(t *testing.T) {
	c := &config.Config{
		Keycloak: config.KeycloakConfig{
			ServerURL: "https://keycloak.example.com",
			Realm:     "motogo",
		},
	}

	expected := "https://keycloak.example.com/realms/motogo"
	assert.Equal(t, expected, c.GetKeycloakIssuerURL())
}
