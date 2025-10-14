package ports

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

// KeycloakClient define las operaciones con Keycloak como servicio externo.
// Esta interfaz abstrae la autenticación y gestión de identidades.
type KeycloakClient interface {

	// Autenticación
	LoginAdmin(ctx context.Context) (*gocloak.JWT, error)

	// Gestión de usuarios
	CreateUser(ctx context.Context, user *gocloak.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	GetUserByID(ctx context.Context, userID string) (*gocloak.User, error)
	UpdateUser(ctx context.Context, user *gocloak.User) error
	DeleteUser(ctx context.Context, userID string) error
	SetPassword(ctx context.Context, userID string, password string, temporary bool) error

	// Roles
	AssignRole(ctx context.Context, userID string, roleName string) error
	RemoveRole(ctx context.Context, userID string, roleName string) error
	GetUserRoles(ctx context.Context, userID string) ([]*gocloak.Role, error)

	// Verificación
	SendVerificationEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, userID string) error

	// Sesiones
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
}
