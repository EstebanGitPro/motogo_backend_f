package output

import "context"

// CustomTokenProvider defines the contract for creating custom authentication tokens
// Implementations live in platform/firebase (e.g., Client)
type CustomTokenProvider interface {
	// CreateCustomToken generates a custom token for the given UID
	CreateCustomToken(ctx context.Context, uid string) (string, error)

	// CreateCustomTokenWithClaims generates a token with additional claims
	CreateCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error)
}
