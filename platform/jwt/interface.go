package jwt

// Validator defines the interface for JWT validation
// This allows for easy mocking in tests
type Validator interface {
	// ValidateToken validates a JWT token and returns the claims if valid
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

// Ensure JWKSValidator implements Validator
var _ Validator = (*JWKSValidator)(nil)
