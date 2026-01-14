package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"google.golang.org/api/option"
)

var log logger.Logger = logger.NewSlogLogger()

// Client wraps Firebase Auth for custom token generation
type Client struct {
	authClient *auth.Client
}

// NewClient initializes Firebase Admin SDK with service account credentials
func NewClient(credentialsPath string) (*Client, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Error(logger.LogFirebaseInitAppError, "error", err)
		return nil, fmt.Errorf("error initializing Firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Error(logger.LogFirebaseAuthClientError, "error", err)
		return nil, fmt.Errorf("error getting Firebase Auth client: %w", err)
	}

	log.Success(logger.LogFirebaseInitOK)
	return &Client{authClient: authClient}, nil
}

// CreateCustomToken generates a Firebase custom token for the given UID
// The UID should be the Keycloak user ID to maintain consistency
func (c *Client) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	token, err := c.authClient.CustomToken(ctx, uid)
	if err != nil {
		log.Error(logger.LogFirebaseTokenCreateError, "uid", uid, "error", err)
		return "", fmt.Errorf("error creating custom token: %w", err)
	}

	log.Info(logger.LogFirebaseTokenCreateOK, "uid", uid)
	return token, nil
}

// CreateCustomTokenWithClaims generates a token with additional claims
// Useful for setting roles or other custom data
func (c *Client) CreateCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	token, err := c.authClient.CustomTokenWithClaims(ctx, uid, claims)
	if err != nil {
		log.Error(logger.LogFirebaseTokenClaimsError, "uid", uid, "error", err)
		return "", fmt.Errorf("error creating custom token with claims: %w", err)
	}

	log.Info(logger.LogFirebaseTokenClaimsOK, "uid", uid)
	return token, nil
}
