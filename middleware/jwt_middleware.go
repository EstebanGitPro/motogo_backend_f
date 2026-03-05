package middleware

import (
	"errors"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	cookiepkg "github.com/EstebanGitPro/motogo-backend/platform/cookie"
	"github.com/EstebanGitPro/motogo-backend/platform/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuth(personService input.Service, msgCache *messaging.MessageCache, jwtValidator output.JWTValidator) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()

	return func(c *gin.Context) {
		token, ok := extractBearerToken(c)
		if !ok {
			return
		}

		claims, err := resolveTokenClaims(c, token, jwtValidator, tokenParser)
		if err != nil {
			return
		}

		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			_ = c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		person, err := personService.GetPersonByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			_ = c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		c.Set("authenticated_user", person)

		c.Next()
	}
}

// extractBearerToken parses the Authorization header and returns the raw token.
// Falls back to the mg_access_token cookie if no header is present (web clients).
// Returns false and aborts if neither source provides a token.
func extractBearerToken(c *gin.Context) (string, bool) {
	// 1. Try Authorization header first (mobile clients)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			_ = c.Error(domain.ErrInvalidToken)
			c.Abort()
			return "", false
		}
		return parts[1], true
	}

	// 2. Fallback: try mg_access_token cookie (web clients)
	token, err := cookiepkg.GetAccessToken(c)
	if err == nil && token != "" {
		return token, true
	}

	// 3. Neither source provided a token
	_ = c.Error(domain.ErrInvalidToken)
	c.Abort()
	return "", false
}

// resolveTokenClaims validates the token and extracts claims using either the JWKS validator
// or the fallback token parser. Aborts and returns an error if validation fails.
func resolveTokenClaims(c *gin.Context, token string, jwtValidator output.JWTValidator, tokenParser *jwt.TokenParser) (map[string]interface{}, error) {
	if jwtValidator == nil {
		claims, err := tokenParser.ExtractClaimsFromToken(token)
		if err != nil {
			_ = c.Error(domain.ErrInvalidToken)
			c.Abort()
			return nil, err
		}
		return claims, nil
	}

	claims, err := jwtValidator.ValidateToken(token)
	if err != nil {
		_ = c.Error(mapJWTError(err))
		c.Abort()
		return nil, err
	}
	return claims, nil
}

// mapJWTError converts JWT validation errors to the appropriate domain error.
func mapJWTError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return domain.ErrTokenExpired
	default:
		return domain.ErrInvalidToken
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

func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		person, exists := GetAuthenticatedUser(c)
		if !exists {
			_ = c.Error(domain.ErrUserNotFound)
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
		_ = c.Error(domain.ErrRoleRequired)
		c.Abort()
	}
}
