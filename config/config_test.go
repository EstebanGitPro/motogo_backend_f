package config_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/stretchr/testify/assert"
)

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
