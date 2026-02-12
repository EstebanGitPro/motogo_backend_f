package middleware

import (
	"errors"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuth(personService input.Service, msgCache *messaging.MessageCache, jwtValidator output.JWTValidator) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()
	_ = tokenParser

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]
		var claims map[string]interface{}
		var err error

		if jwtValidator != nil {
			claims, err = jwtValidator.ValidateToken(token)
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					c.Error(domain.ErrTokenExpired)
				case errors.Is(err, jwt.ErrInvalidSignature):
					c.Error(domain.ErrInvalidToken)
				case errors.Is(err, jwt.ErrInvalidIssuer):
					c.Error(domain.ErrInvalidToken)
				default:
					c.Error(domain.ErrInvalidToken)
				}
				c.Abort()
				return
			}
		} else {
			claims, err = tokenParser.ExtractClaimsFromToken(token)
			if err != nil {
				c.Error(domain.ErrInvalidToken)
				c.Abort()
				return
			}
		}

		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		person, err := personService.GetPersonByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		c.Set("authenticated_user", person)

		c.Next()
	}
}

func GetAuthenticatedUser(c *gin.Context) (*domain.Person, bool) {
	user, exists := c.Get("authenticated_user")
	if !exists {
		return nil, false
	}

	person, ok := user.(*domain.Person)
	return person, ok
}

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
