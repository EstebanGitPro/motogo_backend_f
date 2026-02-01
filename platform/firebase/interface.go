package firebase

import "context"

// FirebaseClient defines the interface for Firebase custom token operations
// This allows for easy mocking in tests
type FirebaseClient interface {
	// CreateCustomToken generates a Firebase custom token for the given UID
	CreateCustomToken(ctx context.Context, uid string) (string, error)

	// CreateCustomTokenWithClaims generates a token with additional claims
	CreateCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error)
}

// Ensure Client implements FirebaseClient
var _ FirebaseClient = (*Client)(nil)
