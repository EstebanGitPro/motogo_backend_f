package keycloak

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/Nerzal/gocloak/v13"
)

type client struct {
	gocloak *gocloak.GoCloak
	config  *config.KeycloakConfig
	token   *gocloak.JWT
	mu      sync.RWMutex // Protección para acceso concurrente al token
}

// NewClient crea una nueva instancia del cliente de Keycloak
func NewClient(cfg *config.KeycloakConfig) (ports.KeycloakClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("keycloak config cannot be nil")
	}

	gc := gocloak.NewClient(cfg.ServerURL)

	c := &client{
		gocloak: gc,
		config:  cfg,
	}

	// Inicializar token de admin
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := c.loginAdminInternal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin token: %w", err)
	}
	c.token = token

	return c, nil
}

// loginAdminInternal es el método interno que realmente hace login
func (c *client) loginAdminInternal(ctx context.Context) (*gocloak.JWT, error) {
	token, err := c.gocloak.LoginAdmin(
		ctx,
		c.config.AdminUser,
		c.config.AdminPass,
		"master", // Admin users exist in master realm
	)
	if err != nil {
		return nil, fmt.Errorf("keycloak admin login failed: %w", err)
	}
	return token, nil
}

// LoginAdmin expone el método públicamente
func (c *client) LoginAdmin(ctx context.Context) (*gocloak.JWT, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	token, err := c.loginAdminInternal(ctx)
	if err != nil {
		return nil, err
	}

	c.token = token
	return token, nil
}

// getToken obtiene el token actual de forma thread-safe
// Si el token está por expirar, lo refresca automáticamente
func (c *client) getToken(ctx context.Context) (string, error) {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	if token == nil {
		return "", fmt.Errorf("no admin token available")
	}

	// Si el token expira en menos de 30 segundos, refrescarlo
	if token.ExpiresIn < 30 {
		c.mu.Lock()
		newToken, err := c.loginAdminInternal(ctx)
		if err != nil {
			c.mu.Unlock()
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}
		c.token = newToken
		token = newToken
		c.mu.Unlock()
	}

	return token.AccessToken, nil
}

func (c *client) CreateUser(ctx context.Context, user *gocloak.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return "", err
	}

	userID, err := c.gocloak.CreateUser(
		ctx,
		accessToken,
		c.config.Realm,
		*user,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	return userID, nil
}

func (c *client) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	users, err := c.gocloak.GetUsers(
		ctx,
		accessToken,
		c.config.Realm,
		gocloak.GetUsersParams{
			Email: &email,
			Exact: gocloak.BoolP(true), // Búsqueda exacta
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return users[0], nil
}

func (c *client) GetUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := c.gocloak.GetUserByID(
		ctx,
		accessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (c *client) UpdateUser(ctx context.Context, user *gocloak.User) error {
	if user == nil || user.ID == nil {
		return fmt.Errorf("user or user ID cannot be nil")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	err = c.gocloak.UpdateUser(
		ctx,
		accessToken,
		c.config.Realm,
		*user,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *client) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	err = c.gocloak.DeleteUser(
		ctx,
		accessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (c *client) SetPassword(ctx context.Context, userID string, password string, temporary bool) error {
	if userID == "" || password == "" {
		return fmt.Errorf("userID and password cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	err = c.gocloak.SetPassword(
		ctx,
		accessToken,
		userID,
		c.config.Realm,
		password,
		temporary,
	)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

func (c *client) AssignRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	// Obtener el role por nombre
	role, err := c.gocloak.GetRealmRole(
		ctx,
		accessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	// Asignar el role al usuario
	err = c.gocloak.AddRealmRoleToUser(
		ctx,
		accessToken,
		c.config.Realm,
		userID,
		[]gocloak.Role{*role},
	)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (c *client) RemoveRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	role, err := c.gocloak.GetRealmRole(
		ctx,
		accessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	err = c.gocloak.DeleteRealmRoleFromUser(
		ctx,
		accessToken,
		c.config.Realm,
		userID,
		[]gocloak.Role{*role},
	)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

func (c *client) GetUserRoles(ctx context.Context, userID string) ([]*gocloak.Role, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := c.gocloak.GetRealmRolesByUserID(
		ctx,
		accessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

func (c *client) SendVerificationEmail(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	params := gocloak.ExecuteActionsEmail{
		UserID:   &userID,
		Actions:  &[]string{"VERIFY_EMAIL"},
		Lifespan: gocloak.IntP(86400), // 24 horas
	}

	err = c.gocloak.ExecuteActionsEmail(
		ctx,
		accessToken,
		c.config.Realm,
		params,
	)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (c *client) VerifyEmail(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	accessToken, err := c.getToken(ctx)
	if err != nil {
		return err
	}

	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	emailVerified := true
	user.EmailVerified = &emailVerified

	err = c.gocloak.UpdateUser(
		ctx,
		accessToken,
		c.config.Realm,
		*user,
	)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (c *client) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("refreshToken cannot be empty")
	}

	err := c.gocloak.Logout(
		ctx,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (c *client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refreshToken cannot be empty")
	}

	token, err := c.gocloak.RefreshToken(
		ctx,
		refreshToken,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return token, nil
}
