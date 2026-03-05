package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

var log = logger.SlogLogger{}

const localConfigFile = "local-config.json"

type Config struct {
	Environment  string          `json:"environment"`
	Database     Database        `json:"database"`
	Server       Server          `json:"server"`
	Resend       Resend          `json:"resend"`
	Verification Verification    `json:"verification"`
	Keycloak     KeycloakConfig  `json:"keycloak"`
	IDEncoder    IDEncoder       `json:"id_encoder"`
	Firebase     FirebaseConfig  `json:"firebase"`
	Geocoding    GeocodingConfig `json:"geocoding"`
	Cookie       CookieConfig    `json:"cookie"`
	CORS         CORSConfig      `json:"cors"`
}

// CORSConfig contiene los orígenes permitidos para Cross-Origin Resource Sharing.
// En local se permiten todos los puertos de desarrollo; en producción solo los dominios reales.
type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

// CookieConfig contiene los parámetros de configuración para cookies HttpOnly.
// Los valores varían entre entornos (local vs producción).
type CookieConfig struct {
	Domain   string `json:"domain"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"same_site"`
}

// ParseSameSite convierte el string de configuración a http.SameSite.
func (cc CookieConfig) ParseSameSite() http.SameSite {
	switch strings.ToLower(cc.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

type FirebaseConfig struct {
	CredentialsPath string `json:"credentials_path"`
}

type GeocodingConfig struct {
	Provider       string `json:"provider"`        // Primary: "google", "mapbox", "opencage"
	APIKey         string `json:"api_key"`         // Primary provider API key
	BaseURL        string `json:"base_url"`        // Primary provider base URL
	TimeoutSeconds int    `json:"timeout_seconds"` // Request timeout
	CountryCode    string `json:"country_code"`    // e.g., "co" for Colombia

	// Fallback provider (optional) - used when primary hits quota limits
	FallbackProvider string `json:"fallback_provider,omitempty"` // "mapbox", "opencage"
	FallbackAPIKey   string `json:"fallback_api_key,omitempty"`
	FallbackBaseURL  string `json:"fallback_base_url,omitempty"`
}

type IDEncoder struct {
	Secret    string `json:"secret"`
	MinLength int    `json:"min_length"`
}

type Verification struct {
	BaseURL string `json:"base_url"`
}

type Database struct {
	Driver          string `json:"driver"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	SSL             string `json:"ssl,omitempty"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time"`
}

type Server struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type Resend struct {
	APIKey    string `json:"api_key"`
	FromEmail string `json:"from_email"`
}

type KeycloakConfig struct {
	ServerURL    string `json:"server_url"`
	Realm        string `json:"realm"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AdminUser    string `json:"admin_user"`
	AdminPass    string `json:"admin_pass"`
}

func LoadConfig() (*Config, error) {
	root, err := utils.FindModuleRoot()
	if err != nil {
		log.Error("error finding module root", err)
		return nil, err
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	var configFile string
	switch env {
	case "railway":
		configFile = "railway-config.json"
	case "production":
		configFile = "prod-config.json"
	default:
		configFile = localConfigFile
	}

	configPath := filepath.Join(root, "config", configFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Warn("Config file not found, falling back to default",
			slog.String("requested_file", configFile),
			slog.String("fallback_file", localConfigFile))
		configPath = filepath.Join(root, "config", localConfigFile)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing JSON configuration: %w", err)
	}

	// ── Sobrescribir con variables de entorno si existen ─────────
	// Patrón: JSON = defaults locales, ENV = override para producción

	// Environment
	config.Environment = envOrStr("APP_ENV", config.Environment)

	// Database
	config.Database.Driver = envOrStr("DB_DRIVER", config.Database.Driver)
	config.Database.Host = envOrStr("DB_HOST", config.Database.Host)
	config.Database.Port = envOrStr("DB_PORT", config.Database.Port)
	config.Database.Username = envOrStr("DB_USERNAME", config.Database.Username)
	config.Database.Password = envOrStr("DB_PASSWORD", config.Database.Password)
	config.Database.Name = envOrStr("DB_NAME", config.Database.Name)
	config.Database.SSL = envOrStr("DB_SSL", config.Database.SSL)
	config.Database.MaxOpenConns = envOrInt("DB_MAX_OPEN_CONNS", config.Database.MaxOpenConns)
	config.Database.MaxIdleConns = envOrInt("DB_MAX_IDLE_CONNS", config.Database.MaxIdleConns)
	config.Database.ConnMaxLifetime = envOrInt("DB_CONN_MAX_LIFETIME", config.Database.ConnMaxLifetime)
	config.Database.ConnMaxIdleTime = envOrInt("DB_CONN_MAX_IDLE_TIME", config.Database.ConnMaxIdleTime)

	// Server
	config.Server.Port = envOrStr("SERVER_PORT", config.Server.Port)
	config.Server.Host = envOrStr("SERVER_HOST", config.Server.Host)

	// Resend
	config.Resend.APIKey = envOrStr("RESEND_API_KEY", config.Resend.APIKey)
	config.Resend.FromEmail = envOrStr("RESEND_FROM_EMAIL", config.Resend.FromEmail)

	// Keycloak
	config.Keycloak.ServerURL = envOrStr("KEYCLOAK_SERVER_URL", config.Keycloak.ServerURL)
	config.Keycloak.Realm = envOrStr("KEYCLOAK_REALM", config.Keycloak.Realm)
	config.Keycloak.ClientID = envOrStr("KEYCLOAK_CLIENT_ID", config.Keycloak.ClientID)
	config.Keycloak.ClientSecret = envOrStr("KEYCLOAK_CLIENT_SECRET", config.Keycloak.ClientSecret)
	config.Keycloak.AdminUser = envOrStr("KEYCLOAK_ADMIN", config.Keycloak.AdminUser)
	config.Keycloak.AdminPass = envOrStr("KEYCLOAK_ADMIN_PASSWORD", config.Keycloak.AdminPass)

	// Verification
	config.Verification.BaseURL = envOrStr("VERIFICATION_BASE_URL", config.Verification.BaseURL)

	// ID Encoder
	config.IDEncoder.Secret = envOrStr("ID_ENCODER_SECRET", config.IDEncoder.Secret)
	config.IDEncoder.MinLength = envOrInt("ID_ENCODER_MIN_LENGTH", config.IDEncoder.MinLength)

	// Firebase
	config.Firebase.CredentialsPath = envOrStr("FIREBASE_CREDENTIALS_PATH", config.Firebase.CredentialsPath)

	// Geocoding
	config.Geocoding.Provider = envOrStr("GEOCODING_PROVIDER", config.Geocoding.Provider)
	config.Geocoding.APIKey = envOrStr("GEOCODING_API_KEY", config.Geocoding.APIKey)
	config.Geocoding.BaseURL = envOrStr("GEOCODING_BASE_URL", config.Geocoding.BaseURL)
	config.Geocoding.TimeoutSeconds = envOrInt("GEOCODING_TIMEOUT", config.Geocoding.TimeoutSeconds)
	config.Geocoding.CountryCode = envOrStr("GEOCODING_COUNTRY_CODE", config.Geocoding.CountryCode)
	config.Geocoding.FallbackProvider = envOrStr("GEOCODING_FALLBACK_PROVIDER", config.Geocoding.FallbackProvider)
	config.Geocoding.FallbackAPIKey = envOrStr("GEOCODING_FALLBACK_API_KEY", config.Geocoding.FallbackAPIKey)
	config.Geocoding.FallbackBaseURL = envOrStr("GEOCODING_FALLBACK_BASE_URL", config.Geocoding.FallbackBaseURL)

	// Cookie
	config.Cookie.Domain = envOrStr("COOKIE_DOMAIN", config.Cookie.Domain)
	config.Cookie.Secure = envOrBool("COOKIE_SECURE", config.Cookie.Secure)
	config.Cookie.SameSite = envOrStr("COOKIE_SAME_SITE", config.Cookie.SameSite)

	// CORS — comma-separated origins override (e.g. "https://app.example.com,https://admin.example.com")
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		parsed := strings.Split(origins, ",")
		trimmed := make([]string, 0, len(parsed))
		for _, o := range parsed {
			if s := strings.TrimSpace(o); s != "" {
				trimmed = append(trimmed, s)
			}
		}
		config.CORS.AllowedOrigins = trimmed
	}

	slog.Info("Configuration loaded successfully",
		slog.String("config_file", configFile),
		slog.String("environment", config.Environment),
		slog.String("config_path", configPath),
		slog.String("keycloak_server", config.Keycloak.ServerURL),
		slog.String("keycloak_realm", config.Keycloak.Realm))

	return &config, nil
}

// envOrStr retorna el valor de la variable de entorno si existe, o el fallback.
func envOrStr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// envOrInt retorna el valor (int) de la variable de entorno si existe, o el fallback.
func envOrInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.Atoi(val)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

// envOrBool retorna el valor (bool) de la variable de entorno si existe, o el fallback.
func envOrBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		return strings.EqualFold(val, "true") || val == "1"
	}
	return fallback
}

func (c *Config) GetMySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	if c.Database.SSL != "" {
		dsn += "&tls=" + c.Database.SSL
	}

	return dsn
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "railway"
}

// Helper para obtener la URL completa del auth endpoint de Keycloak
func (c *Config) GetKeycloakAuthURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

// Helper para obtener la URL del admin API
func (c *Config) GetKeycloakAdminURL() string {
	return fmt.Sprintf("%s/admin/realms/%s",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

// Helper para obtener la URL del endpoint JWKS de Keycloak (para validación JWT)
func (c *Config) GetKeycloakJWKSURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

// Helper para obtener la URL del issuer de Keycloak (para validación JWT)
func (c *Config) GetKeycloakIssuerURL() string {
	return fmt.Sprintf("%s/realms/%s",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}
