package keycloak

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/Nerzal/gocloak/v13"
)

var log logger.Logger = logger.NewSlogLogger()

type client struct {
	gocloak        *gocloak.GoCloak
	config         *config.KeycloakConfig
	token          *gocloak.JWT
	tokenExpiresAt time.Time
	tokenMutex     sync.RWMutex
}

func NewClient(cfg *config.KeycloakConfig) (output.AuthClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("keycloak config cannot be nil")
	}

	log.Info(logger.LogKeycloakClientInit, "server_url", cfg.ServerURL, "realm", cfg.Realm)

	gc := gocloak.NewClient(cfg.ServerURL)

	authClient := &client{
		gocloak: gc,
		config:  cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Debug(logger.LogKeycloakAdminAuth, "admin_user", cfg.AdminUser, "realm", cfg.Realm)
	token, err := authClient.gocloak.LoginAdmin(ctx, authClient.config.AdminUser, authClient.config.AdminPass, authClient.config.Realm)
	if err != nil {
		log.Error(logger.LogKeycloakAdminAuthError, "error", err, "realm", cfg.Realm)
		return nil, fmt.Errorf("failed to initialize admin token: %w", err)
	}
	authClient.token = token
	authClient.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	log.Success(logger.LogKeycloakClientOK, "realm", cfg.Realm, "expires_in", token.ExpiresIn)

	return authClient, nil
}

func (c *client) ensureValidToken(ctx context.Context) (string, error) {
	c.tokenMutex.RLock()

	needsRefresh := time.Now().Add(30 * time.Second).After(c.tokenExpiresAt)
	currentToken := c.token.AccessToken
	c.tokenMutex.RUnlock()

	if !needsRefresh {
		return currentToken, nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if time.Now().Add(30 * time.Second).Before(c.tokenExpiresAt) {
		return c.token.AccessToken, nil
	}

	log.Info(logger.LogKeycloakTokenRefresh,
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"token_expires_at", c.tokenExpiresAt.Format(time.RFC3339))

	token, err := c.gocloak.LoginAdmin(ctx, c.config.AdminUser, c.config.AdminPass, c.config.Realm)
	if err != nil {
		log.Error(logger.LogKeycloakTokenRefreshErr,
			"realm", c.config.Realm,
			"admin_user", c.config.AdminUser,
			"error", err)
		return "", fmt.Errorf("failed to refresh admin token: %w", err)
	}

	c.token = token
	c.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	log.Success(logger.LogKeycloakTokenRefreshOK,
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"new_expires_at", c.tokenExpiresAt.Format(time.RFC3339),
		"expires_in_seconds", token.ExpiresIn)

	return c.token.AccessToken, nil
}

func (c *client) LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	log.Info(logger.LogKeycloakUserLogin, "username", username, "realm", c.config.Realm)

	//TODO: Me dice que no he compeltado los campos y yo le he mandado todo lo necesario.
	token, err := c.gocloak.Login(
		ctx,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
		username,
		password,
	)
	if err != nil {
		log.Error(logger.LogKeycloakUserLoginError, "username", username, "error", err)
		return nil, fmt.Errorf("user login failed: %w", err)
	}

	log.Success(logger.LogKeycloakUserLoginOK, "username", username)
	return token, nil
}

func (c *client) CreateUser(ctx context.Context, person *domain.Person) (string, error) {
	if person == nil {
		return "", fmt.Errorf("person cannot be nil")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return "", err
	}

	log.Info(logger.LogKeycloakUserCreate, "email", person.Email, "realm", c.config.Realm)

	keycloakUser := gocloak.User{
		Email:     &person.Email,
		FirstName: &person.FirstName,
		LastName:  &person.LastName,
		Enabled:   gocloak.BoolP(true),
		Username:  &person.Email,
	}

	userID, err := c.gocloak.CreateUser(
		ctx,
		token,
		c.config.Realm,
		keycloakUser,
	)
	if err != nil {
		log.Error(logger.LogKeycloakUserCreateError,
			"email", person.Email,
			"error", err,
			"error_details", err.Error(),
			"realm", c.config.Realm)
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	log.Success(logger.LogKeycloakUserCreateOK, "email", person.Email, "user_id", userID)
	return userID, nil
}

func (c *client) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return nil, err
	}

	users, err := c.gocloak.GetUsers(
		ctx,
		token,
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

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := c.gocloak.GetUserByID(
		ctx,
		token,
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

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	err = c.gocloak.UpdateUser(
		ctx,
		token,
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

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	log.Warn(logger.LogKeycloakUserDelete, "user_id", userID)

	err = c.gocloak.DeleteUser(
		ctx,
		token,
		c.config.Realm,
		userID,
	)
	if err != nil {
		log.Error(logger.LogKeycloakUserDeleteError, "user_id", userID, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	log.Info(logger.LogKeycloakUserDeleteOK, "user_id", userID)
	return nil
}

func (c *client) SetPassword(ctx context.Context, userID string, password string, temporary bool) error {
	if userID == "" || password == "" {
		return fmt.Errorf("userID and password cannot be empty")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	log.Debug(logger.LogKeycloakPasswordSet, "user_id", userID, "temporary", temporary)

	err = c.gocloak.SetPassword(
		ctx,
		token,
		userID,
		c.config.Realm,
		password,
		temporary,
	)
	if err != nil {
		log.Error(logger.LogKeycloakPasswordSetError, "user_id", userID, "error", err)
		return fmt.Errorf("failed to set password: %w", err)
	}

	log.Success(logger.LogKeycloakPasswordSetOK, "user_id", userID)
	return nil
}

func (c *client) AssignRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	log.Info(logger.LogKeycloakRoleAssign, "user_id", userID, "role", roleName)

	// Obtener el role por nombre
	role, err := c.gocloak.GetRealmRole(
		ctx,
		token,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		log.Error(logger.LogKeycloakRoleGetError, "role", roleName, "error", err)
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	// Asignar el role al usuario
	err = c.gocloak.AddRealmRoleToUser(
		ctx,
		token,
		c.config.Realm,
		userID,
		[]gocloak.Role{*role},
	)
	if err != nil {
		log.Error(logger.LogKeycloakRoleAssignError, "user_id", userID, "role", roleName, "error", err)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	log.Success(logger.LogKeycloakRoleAssignOK, "user_id", userID, "role", roleName)
	return nil
}

func (c *client) RemoveRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	role, err := c.gocloak.GetRealmRole(
		ctx,
		token,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	err = c.gocloak.DeleteRealmRoleFromUser(
		ctx,
		token,
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

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := c.gocloak.GetRealmRolesByUserID(
		ctx,
		token,
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

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	log.Info(logger.LogKeycloakSendVerificationEmail, "user_id", userID, "realm", c.config.Realm)

	params := gocloak.ExecuteActionsEmail{
		UserID:   &userID,
		Actions:  &[]string{"VERIFY_EMAIL"},
		Lifespan: gocloak.IntP(86400), // 24 horas
	}

	err = c.gocloak.ExecuteActionsEmail(
		ctx,
		token,
		c.config.Realm,
		params,
	)
	if err != nil {
		log.Error(logger.LogKeycloakSendVerificationEmailError, "user_id", userID, "error", err)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Success(logger.LogKeycloakSendVerificationEmailOK, "user_id", userID)
	return nil
}

// SendPasswordResetEmail sends a password reset email to the user
// It searches for the user by email first, then sends the reset email
func (c *client) SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	token, err := c.ensureValidToken(ctx)
	if err != nil {
		return err
	}

	log.Info(logger.LogKeycloakSendPasswordReset, "email", email, "realm", c.config.Realm)

	// Buscar usuario por email
	log.Debug(logger.LogKeycloakSearchUserByEmail, "email", email)

	users, err := c.gocloak.GetUsers(
		ctx,
		token,
		c.config.Realm,
		gocloak.GetUsersParams{
			Email: &email,
			Exact: gocloak.BoolP(true),
		},
	)
	if err != nil {
		log.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return fmt.Errorf("failed to search user: %w", err)
	}

	if len(users) == 0 {
		log.Warn(logger.LogKeycloakUserNotFound, "email", email)
		return fmt.Errorf("user with email %s not found", email)
	}

	log.Debug(logger.LogKeycloakSearchUserByEmailOK, "email", email, "user_id", *users[0].ID)

	// Enviar email de reset de contraseña
	params := gocloak.ExecuteActionsEmail{
		UserID:   users[0].ID,
		Actions:  &[]string{"UPDATE_PASSWORD"},
		Lifespan: gocloak.IntP(43200), // 12 horas
	}

	err = c.gocloak.ExecuteActionsEmail(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		params,
	)
	if err != nil {
		log.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Success(logger.LogKeycloakSendPasswordResetOK, "email", email, "user_id", *users[0].ID)
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

	log.Info(logger.LogKeycloakUserTokenRefresh,
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
		log.Error(logger.LogKeycloakUserTokenRefreshErr,
			"realm", c.config.Realm,
			"client_id", c.config.ClientID,
			"error", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	log.Success(logger.LogKeycloakUserTokenRefreshOK,
		"realm", c.config.Realm,
		"client_id", c.config.ClientID,
		"expires_in_seconds", token.ExpiresIn,
		"refresh_expires_in_seconds", token.RefreshExpiresIn)

	return token, nil
}
