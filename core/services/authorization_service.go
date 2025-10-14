package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/Nerzal/gocloak/v13"
)

type authorizationService struct {
	keycloakClient ports.KeycloakClient
	repository     ports.Repository
	
	// Cache en memoria para mapear PersonID -> KeycloakUserID
	// En producción, considera usar Redis o base de datos
	userMapping map[string]string
	mu          sync.RWMutex
}

// NewAuthorizationService crea un nuevo servicio de autorización
func NewAuthorizationService(keycloakClient ports.KeycloakClient, repository ports.Repository) ports.AuthorizationService {
	return &authorizationService{
		keycloakClient: keycloakClient,
		repository:     repository,
		userMapping:    make(map[string]string),
	
	}
}

// SyncUserToKeycloak sincroniza un usuario de tu aplicación con Keycloak
func (a *authorizationService) SyncUserToKeycloak(ctx context.Context, person *domain.Person) (string, error) {
	// Verificar si ya está sincronizado
	a.mu.RLock()
	if keycloakUserID, exists := a.userMapping[person.ID]; exists {
		a.mu.RUnlock()
		return keycloakUserID, nil
	}
	a.mu.RUnlock()

	// Crear usuario en Keycloak
	keycloakUser := &gocloak.User{
		Email:         &person.Email,
		FirstName:     &person.FirstName,
		LastName:      &person.LastName,
		EmailVerified: &person.EmailVerified,
		Enabled:       gocloak.BoolP(true),
		Username:      &person.Email, // Usar email como username
	}

	keycloakUserID, err := a.keycloakClient.CreateUser(ctx, keycloakUser)
	if err != nil {
		// Si el usuario ya existe, intentar obtenerlo
		existingUser, getErr := a.keycloakClient.GetUserByEmail(ctx, person.Email)
		if getErr != nil {
			return "", fmt.Errorf("failed to create or get user in keycloak: %w", err)
		}
		keycloakUserID = *existingUser.ID
	}

	// Guardar mapeo en cache
	a.mu.Lock()
	a.userMapping[person.ID] = keycloakUserID
	a.mu.Unlock()

	return keycloakUserID, nil
}

// AssignRole asigna un rol a un usuario
func (a *authorizationService) AssignRole(ctx context.Context, personID string, roleName string) error {
	// Obtener el usuario de la base de datos
	person, err := a.repository.GetPersonByID(personID)
	if err != nil {
		return fmt.Errorf("person not found: %w", err)
	}

	// Sincronizar con Keycloak si es necesario
	keycloakUserID, err := a.SyncUserToKeycloak(ctx, person)
	if err != nil {
		return fmt.Errorf("failed to sync user to keycloak: %w", err)
	}

	// Asignar rol en Keycloak
	err = a.keycloakClient.AssignRole(ctx, keycloakUserID, roleName)
	if err != nil {
		return fmt.Errorf("failed to assign role in keycloak: %w", err)
	}

	// Actualizar rol en la base de datos local también
	person.Role = roleName
	err = a.repository.Update(*person)
	if err != nil {
		// Log warning pero no fallar, el rol ya está en Keycloak
		fmt.Printf("Warning: failed to update role in local database: %v\n", err)
	}

	return nil
}

// RemoveRole remueve un rol de un usuario
func (a *authorizationService) RemoveRole(ctx context.Context, personID string, roleName string) error {
	keycloakUserID, err := a.getKeycloakUserID(ctx, personID)
	if err != nil {
		return err
	}

	return a.keycloakClient.RemoveRole(ctx, keycloakUserID, roleName)
}

// GetUserRoles obtiene todos los roles de un usuario
func (a *authorizationService) GetUserRoles(ctx context.Context, personID string) ([]string, error) {
	keycloakUserID, err := a.getKeycloakUserID(ctx, personID)
	if err != nil {
		return nil, err
	}

	roles, err := a.keycloakClient.GetUserRoles(ctx, keycloakUserID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}

	return roleNames, nil
}

// HasRole verifica si un usuario tiene un rol específico
func (a *authorizationService) HasRole(ctx context.Context, personID string, roleName string) (bool, error) {
	roles, err := a.GetUserRoles(ctx, personID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role == roleName {
			return true, nil
		}
	}

	return false, nil
}

// HasPermission verifica si un usuario tiene permisos para una acción específica
func (a *authorizationService) HasPermission(ctx context.Context, personID string, resource, action string) (bool, error) {
	// Implementación básica basada en roles
	// Puedes expandir esto según tus necesidades de negocio
	
	roles, err := a.GetUserRoles(ctx, personID)
	if err != nil {
		return false, err
	}

	// Ejemplo de lógica de permisos básica
	for _, role := range roles {
		switch role {
		case "admin":
			return true, nil // Admin tiene todos los permisos
		case "moderator":
			if resource == "users" && (action == "read" || action == "update") {
				return true, nil
			}
		case "user":
			if resource == "profile" && action == "read" {
				return true, nil
			}
		}
	}

	return false, nil
}

// CreateRole crea un nuevo rol en Keycloak
func (a *authorizationService) CreateRole(ctx context.Context, roleName, description string) error {
	// Primero verificar si el rol ya existe
	roles, err := a.keycloakClient.GetUserRoles(ctx, "dummy") // Esto fallará pero nos permite usar la interfaz existente
	if err == nil {
		for _, role := range roles {
			if role.Name != nil && *role.Name == roleName {
				return fmt.Errorf("role %s already exists", roleName)
			}
		}
	}

	// Crear el rol usando la funcionalidad de Keycloak
	// Nota: Necesitarías agregar este método a la interfaz KeycloakClient si no existe
	return fmt.Errorf("create role functionality needs to be implemented in KeycloakClient")
}

// GetAllRoles obtiene todos los roles disponibles
func (a *authorizationService) GetAllRoles(ctx context.Context) ([]string, error) {
	// Similar al método anterior, necesitarías implementar esto en KeycloakClient
	return nil, fmt.Errorf("get all roles functionality needs to be implemented in KeycloakClient")
}

// getKeycloakUserID obtiene el ID de Keycloak para un PersonID
func (a *authorizationService) getKeycloakUserID(ctx context.Context, personID string) (string, error) {
	// Verificar cache primero
	a.mu.RLock()
	if keycloakUserID, exists := a.userMapping[personID]; exists {
		a.mu.RUnlock()
		return keycloakUserID, nil
	}
	a.mu.RUnlock()

	// Obtener persona de la base de datos y sincronizar
	person, err := a.repository.GetPersonByID(personID)
	if err != nil {
		return "", fmt.Errorf("person not found: %w", err)
	}

	return a.SyncUserToKeycloak(ctx, person)
}
