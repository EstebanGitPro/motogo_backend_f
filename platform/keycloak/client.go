package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/Nerzal/gocloak/v13"
)

type client struct {
	gocloak        *gocloak.GoCloak
	config         *config.KeycloakConfig
	token          *gocloak.JWT
	tokenExpiresAt time.Time
	tokenMutex     sync.RWMutex
}

func NewClient(cfg *config.KeycloakConfig) (ports.AuthClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("keycloak config cannot be nil")
	}

	gc := gocloak.NewClient(cfg.ServerURL)

	authClient := &client{
		gocloak: gc,
		config:  cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	token, err := authClient.gocloak.LoginAdmin(ctx, authClient.config.AdminUser, authClient.config.AdminPass, authClient.config.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin token: %w", err)
	}
	authClient.token = token
	authClient.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return authClient, nil
}


func (c *client) ensureValidToken(ctx context.Context) error {
	c.tokenMutex.RLock()
	
	needsRefresh := time.Now().Add(30 * time.Second).After(c.tokenExpiresAt)
	c.tokenMutex.RUnlock()

	if !needsRefresh {
		return nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	
	if time.Now().Add(30 * time.Second).Before(c.tokenExpiresAt) {
		return nil
	}

	slog.Info("Refreshing Keycloak admin token",
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"token_expires_at", c.tokenExpiresAt.Format(time.RFC3339))

	
	token, err := c.gocloak.LoginAdmin(ctx, c.config.AdminUser, c.config.AdminPass, c.config.Realm)
	if err != nil {
		slog.Error("Failed to refresh Keycloak admin token",
			"realm", c.config.Realm,
			"admin_user", c.config.AdminUser,
			"error", err)
		return fmt.Errorf("failed to refresh admin token: %w", err)
	}

	c.token = token
	c.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	slog.Info("Keycloak admin token refreshed successfully",
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"new_expires_at", c.tokenExpiresAt.Format(time.RFC3339),
		"expires_in_seconds", token.ExpiresIn)

	return nil
}

func (c *client) LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	token, err := c.gocloak.Login(
		ctx,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
		username, 
		password,
	)
	if err != nil {
		return nil, fmt.Errorf("user login failed: %w", err)
	}

	return token, nil
}

func (c *client) CreateUser(ctx context.Context, person *domain.Person) (string, error) {
	if person == nil {
		return "", fmt.Errorf("person cannot be nil")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return "", err
	}

	keycloakUser := gocloak.User{
		Email:         &person.Email,
		FirstName:     &person.FirstName,
		LastName:      &person.LastName,
		Enabled:       gocloak.BoolP(true),
		Username:      &person.Email,
	}

	userID, err := c.gocloak.CreateUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		keycloakUser,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	
	users, err := c.gocloak.GetUsers(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		gocloak.GetUsersParams{
			Email: &email,
			Exact: gocloak.BoolP(true),
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

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	user, err := c.gocloak.GetUserByID(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	err := c.gocloak.UpdateUser(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	err := c.gocloak.DeleteUser(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	err := c.gocloak.SetPassword(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	// Obtener el role por nombre
	role, err := c.gocloak.GetRealmRole(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	// Asignar el role al usuario
	err = c.gocloak.AddRealmRoleToUser(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	role, err := c.gocloak.GetRealmRole(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	err = c.gocloak.DeleteRealmRoleFromUser(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	roles, err := c.gocloak.GetRealmRolesByUserID(
		ctx,
		c.token.AccessToken,
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

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	params := gocloak.ExecuteActionsEmail{
		UserID:   &userID,
		Actions:  &[]string{"VERIFY_EMAIL"},
		Lifespan: gocloak.IntP(86400), // 24 horas
	}

	err := c.gocloak.ExecuteActionsEmail(
		ctx,
		c.token.AccessToken,
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


	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	emailVerified := true
	user.EmailVerified = &emailVerified

	err = c.gocloak.UpdateUser(
		ctx,
		c.token.AccessToken,
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

	slog.Info("Refreshing user token",
		"realm", c.config.Realm,
		"client_id", c.config.ClientID)

	token, err := c.gocloak.RefreshToken(
		ctx,
		refreshToken,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
	)
	if err != nil {
		slog.Error("Failed to refresh user token",
			"realm", c.config.Realm,
			"client_id", c.config.ClientID,
			"error", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	slog.Info("User token refreshed successfully",
		"realm", c.config.Realm,
		"client_id", c.config.ClientID,
		"expires_in_seconds", token.ExpiresIn,
		"refresh_expires_in_seconds", token.RefreshExpiresIn)

	return token, nil
}
