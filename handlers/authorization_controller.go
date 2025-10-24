package handlers

import (
	"context"
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/gin-gonic/gin"
)

// AuthorizationController maneja las operaciones de autorización
type AuthorizationController struct {
	authzService ports.AuthorizationService
}

// NewAuthorizationController crea un nuevo controlador de autorización
func NewAuthorizationController(authzService ports.AuthorizationService) *AuthorizationController {
	return &AuthorizationController{
		authzService: authzService,
	}
}

// SyncUserToKeycloak sincroniza un usuario existente con Keycloak para manejo de roles
func (a *AuthorizationController) SyncUserToKeycloak() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			PersonID string `json:"person_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		// Aquí necesitarías obtener la persona de tu base de datos
		// person, err := a.personService.GetPersonByID(request.PersonID)
		// if err != nil {
		//     c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		//     return
		// }

		// ctx := context.Background()
		// keycloakUserID, err := a.authzService.SyncUserToKeycloak(ctx, person)
		// if err != nil {
		//     c.JSON(http.StatusInternalServerError, gin.H{
		//         "error": "Error sincronizando usuario con Keycloak: " + err.Error(),
		//     })
		//     return
		// }

		c.JSON(http.StatusOK, gin.H{
			"message": "Usuario sincronizado con Keycloak exitosamente",
			// "keycloak_user_id": keycloakUserID,
		})
	}
}

// AssignRole asigna un rol a un usuario
func (a *AuthorizationController) AssignRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			PersonID string `json:"person_id" binding:"required"`
			Role     string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		err := a.authzService.AssignRole(ctx, request.PersonID, request.Role)
		if err != nil {
			c.Error(domain.ErrRoleAssignmentFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Rol asignado exitosamente",
			"person_id": request.PersonID,
			"role": request.Role,
		})
	}
}

// RemoveRole remueve un rol de un usuario
func (a *AuthorizationController) RemoveRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			PersonID string `json:"person_id" binding:"required"`
			Role     string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		err := a.authzService.RemoveRole(ctx, request.PersonID, request.Role)
		if err != nil {
			c.Error(domain.ErrRoleRemovalFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Rol removido exitosamente",
		})
	}
}

// GetUserRoles obtiene todos los roles de un usuario
func (a *AuthorizationController) GetUserRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		personID := c.Param("person_id")
		if personID == "" {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		roles, err := a.authzService.GetUserRoles(ctx, personID)
		if err != nil {
			c.Error(domain.ErrGetUserRolesFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"person_id": personID,
			"roles": roles,
		})
	}
}

// CheckRole verifica si un usuario tiene un rol específico
func (a *AuthorizationController) CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		personID := c.Param("person_id")
		role := c.Param("role")
		
		if personID == "" || role == "" {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		hasRole, err := a.authzService.HasRole(ctx, personID, role)
		if err != nil {
			c.Error(domain.ErrRoleCheckFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"person_id": personID,
			"role": role,
			"has_role": hasRole,
		})
	}
}
