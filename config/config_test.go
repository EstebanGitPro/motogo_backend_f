package config_test

import (
	"net/http"
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
	// so it should fall back to local-config.json for defaults,
	// but APP_ENV override sets environment to "railway"
	t.Setenv("APP_ENV", "railway")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// APP_ENV env var overrides the JSON environment field
	assert.Equal(t, "railway", cfg.Environment)
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

func TestLoadConfig_DatabaseEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("DB_HOST", "db.production.com")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USERNAME", "prod_user")
	t.Setenv("DB_PASSWORD", "prod_secret")
	t.Setenv("DB_NAME", "motogo_prod")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "db.production.com", cfg.Database.Host)
	assert.Equal(t, "3306", cfg.Database.Port)
	assert.Equal(t, "prod_user", cfg.Database.Username)
	assert.Equal(t, "prod_secret", cfg.Database.Password)
	assert.Equal(t, "motogo_prod", cfg.Database.Name)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
}

func TestLoadConfig_ServerEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_HOST", "127.0.0.1")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
}

func TestLoadConfig_ResendEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("RESEND_API_KEY", "re_prod_key")
	t.Setenv("RESEND_FROM_EMAIL", "no-reply@motogo.com")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "re_prod_key", cfg.Resend.APIKey)
	assert.Equal(t, "no-reply@motogo.com", cfg.Resend.FromEmail)
}

func TestLoadConfig_GeocodingEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("GEOCODING_API_KEY", "prod-google-key")
	t.Setenv("GEOCODING_TIMEOUT", "10")
	t.Setenv("GEOCODING_FALLBACK_API_KEY", "prod-mapbox-key")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "prod-google-key", cfg.Geocoding.APIKey)
	assert.Equal(t, 10, cfg.Geocoding.TimeoutSeconds)
	assert.Equal(t, "prod-mapbox-key", cfg.Geocoding.FallbackAPIKey)
}

func TestLoadConfig_InvalidIntEnvFallsBackToJSON(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "not-a-number")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	// Should keep the JSON default (25) since the env var is invalid
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
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

// ============================================
// ParseSameSite Tests
// ============================================

func TestParseSameSite_Strict(t *testing.T) {
	cc := config.CookieConfig{SameSite: "strict"}
	assert.Equal(t, http.SameSiteStrictMode, cc.ParseSameSite())
}

func TestParseSameSite_StrictUpperCase(t *testing.T) {
	cc := config.CookieConfig{SameSite: "STRICT"}
	assert.Equal(t, http.SameSiteStrictMode, cc.ParseSameSite())
}

func TestParseSameSite_None(t *testing.T) {
	cc := config.CookieConfig{SameSite: "none"}
	assert.Equal(t, http.SameSiteNoneMode, cc.ParseSameSite())
}

func TestParseSameSite_Lax(t *testing.T) {
	cc := config.CookieConfig{SameSite: "lax"}
	assert.Equal(t, http.SameSiteLaxMode, cc.ParseSameSite())
}

func TestParseSameSite_DefaultFallback(t *testing.T) {
	cc := config.CookieConfig{SameSite: "unknown"}
	assert.Equal(t, http.SameSiteLaxMode, cc.ParseSameSite())
}

func TestParseSameSite_Empty(t *testing.T) {
	cc := config.CookieConfig{SameSite: ""}
	assert.Equal(t, http.SameSiteLaxMode, cc.ParseSameSite())
}

// ============================================
// envOrBool Tests (via LoadConfig Cookie.Secure override)
// ============================================

func TestLoadConfig_CookieEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("COOKIE_DOMAIN", ".motogo.com")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAME_SITE", "none")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, ".motogo.com", cfg.Cookie.Domain)
	assert.True(t, cfg.Cookie.Secure)
	assert.Equal(t, "none", cfg.Cookie.SameSite)
}

func TestLoadConfig_CookieSecureFalse(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("COOKIE_SECURE", "false")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.False(t, cfg.Cookie.Secure)
}

func TestLoadConfig_CookieSecureNumeric(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("COOKIE_SECURE", "1")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.True(t, cfg.Cookie.Secure)
}

// ============================================
// CORS env override Tests
// ============================================

func TestLoadConfig_CORSEnvOverride(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.motogo.com, https://admin.motogo.com")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Len(t, cfg.CORS.AllowedOrigins, 2)
	assert.Equal(t, "https://app.motogo.com", cfg.CORS.AllowedOrigins[0])
	assert.Equal(t, "https://admin.motogo.com", cfg.CORS.AllowedOrigins[1])
}

func TestLoadConfig_CORSEnvOverrideSingleOrigin(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.motogo.com")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Len(t, cfg.CORS.AllowedOrigins, 1)
	assert.Equal(t, "https://app.motogo.com", cfg.CORS.AllowedOrigins[0])
}

func TestLoadConfig_CORSEnvOverrideWithEmptyEntries(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.motogo.com,,, https://admin.motogo.com, ")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.Len(t, cfg.CORS.AllowedOrigins, 2)
}

// ============================================
// Production config fallback Tests
// ============================================

func TestLoadConfig_ProductionEnv_FallsBackToLocal(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "production", cfg.Environment)
}
