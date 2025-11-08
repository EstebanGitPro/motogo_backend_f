package output

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/Nerzal/gocloak/v13"
)

type AuthClient interface {

	// Autenticación
	LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) // Login de usuario normal

	// Gestión de usuarios
	CreateUser(ctx context.Context,person *domain.Person) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	GetUserByID(ctx context.Context, userID string) (*gocloak.User, error)
	UpdateUser(ctx context.Context, user *gocloak.User) error
	DeleteUser(ctx context.Context, userID string) error
	SetPassword(ctx context.Context,userID string, password string, temporary bool) error

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

type Tx interface {
	Commit() error
	Rollback() error
}
