package middleware

import (
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/jwt"
	"github.com/gin-gonic/gin"
)

// RequireAuth creates a middleware that validates JWT tokens from Keycloak
// and injects the authenticated user into the Gin context
func RequireAuth(personService input.Service, msgCache *messaging.MessageCache) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()

	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]

		// Decode JWT and extract claims
		// Note: We extract the "sub" claim which contains the Keycloak User ID
		claims, err := tokenParser.ExtractClaimsFromToken(token)
		if err != nil {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Extract Keycloak User ID from "sub" claim
		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Find user in database by Keycloak ID
		person, err := personService.GetPersonByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			// User not found in our database
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		// Inject authenticated user into context
		c.Set("authenticated_user", person)

		c.Next()
	}
}

// GetAuthenticatedUser extracts the authenticated user from the Gin context
func GetAuthenticatedUser(c *gin.Context) (*domain.Person, bool) {
	user, exists := c.Get("authenticated_user")
	if !exists {
		return nil, false
	}

	person, ok := user.(*domain.Person)
	return person, ok
}

// RequireRole creates a middleware that validates the user has the required role
// Must be used AFTER RequireAuth middleware
// Example usage: router.POST("/branches", RequireRole(domain.RoleRepresentative), handler.RegisterBranch())
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		person, exists := GetAuthenticatedUser(c)
		if !exists {
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		for _, role := range allowedRoles {
			if person.Role == role {
				c.Next()
				return
			}
		}

		// Role not allowed
		c.Error(domain.ErrRoleRequired)
		c.Abort()
	}
}
