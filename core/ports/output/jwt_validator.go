package output

// JWTValidator defines the contract for JWT token validation
// Implementations live in platform/jwt (e.g., JWKSValidator)
type JWTValidator interface {
	// ValidateToken validates a JWT token and returns the claims if valid
	ValidateToken(tokenString string) (map[string]interface{}, error)
}
