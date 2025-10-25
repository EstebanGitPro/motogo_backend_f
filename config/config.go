package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

type Config struct {
	Environment  string         `json:"environment"`
	Database     Database       `json:"database"`
	Server       Server         `json:"server"`
	Resend       Resend         `json:"resend"`
	Verification Verification   `json:"verification"`
	Keycloak     KeycloakConfig `json:"keycloak"`
}

type Verification struct {
	BaseURL string `json:"base_url"`
}

type Database struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSL      string `json:"ssl,omitempty"`
	MaxOpenConns int `json:"max_open_conns"`
	MaxIdleConns int `json:"max_idle_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	ConnMaxIdleTime int `json:"conn_max_idle_time"`
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
		return nil, fmt.Errorf("error finding module root: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	var configFile string
	switch env {
	case "railway":
		configFile = "railway-config.json"
	default:
		configFile = "local-config.json"
	}

	configPath := filepath.Join(root, "config", configFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Warn("Config file not found, falling back to default",
			slog.String("requested_file", configFile),
			slog.String("fallback_file", "local-config.json"))
		configPath = filepath.Join(root, "config", "local-config.json")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing JSON configuration: %w", err)
	}

	// Sobrescribir configuración de Keycloak con variables de entorno si existen
	// Esto permite: JSON para defaults locales, ENV para producción/secrets
	if serverURL := os.Getenv("KEYCLOAK_SERVER_URL"); serverURL != "" {
		config.Keycloak.ServerURL = serverURL
	}
	if realm := os.Getenv("KEYCLOAK_REALM"); realm != "" {
		config.Keycloak.Realm = realm
	}
	if clientID := os.Getenv("KEYCLOAK_CLIENT_ID"); clientID != "" {
		config.Keycloak.ClientID = clientID
	}
	if clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET"); clientSecret != "" {
		config.Keycloak.ClientSecret = clientSecret
	}
	if adminUser := os.Getenv("KEYCLOAK_ADMIN"); adminUser != "" {
		config.Keycloak.AdminUser = adminUser
	}
	if adminPass := os.Getenv("KEYCLOAK_ADMIN_PASSWORD"); adminPass != "" {
		config.Keycloak.AdminPass = adminPass
	}

	slog.Info("Configuration loaded successfully",
		slog.String("config_file", configFile),
		slog.String("environment", config.Environment),
		slog.String("config_path", configPath),
		slog.String("keycloak_server", config.Keycloak.ServerURL),
		slog.String("keycloak_realm", config.Keycloak.Realm))

	return &config, nil
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