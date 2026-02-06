package firebase

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"google.golang.org/api/option"
)

var log logger.Logger = logger.NewSlogLogger()

// Client wraps Firebase Auth and Storage for token generation and file management
type Client struct {
	authClient    *auth.Client
	storageClient *storage.Client
	bucketName    string
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

	// Initialize Storage client (optional - may fail if not configured)
	storageClient, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Warn(logger.LogFirebaseStorageClientWarn, "error", err)
		// Continue without storage - some deployments may not need it
		storageClient = nil
	}

	// Extract bucket name from credentials or use default pattern
	bucketName := extractBucketName(credentialsPath)

	log.Success(logger.LogFirebaseInitOK)
	return &Client{
		authClient:    authClient,
		storageClient: storageClient,
		bucketName:    bucketName,
	}, nil
}

// extractBucketName attempts to derive the bucket name from project
func extractBucketName(credentialsPath string) string {
	// Default bucket follows pattern: {project-id}.appspot.com
	// This will be overridden if we can parse it from the URL
	return ""
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

// DeleteStorageFile deletes a file from Firebase Storage given its full URL
// Example URL: https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/motorcycles%2Fuser%2Fprofile.jpg?alt=media&token=xxx
func (c *Client) DeleteStorageFile(ctx context.Context, fileURL string) error {
	if c.storageClient == nil {
		log.Warn(logger.LogFirebaseStorageNotConfigured, "url", fileURL)
		return nil // Silently ignore if storage not configured
	}

	if fileURL == "" {
		return nil // Nothing to delete
	}

	// Parse the Firebase Storage URL to extract bucket and object path
	bucketName, objectPath, err := parseFirebaseStorageURL(fileURL)
	if err != nil {
		log.Error(logger.LogFirebaseStorageParseError, "url", fileURL, "error", err)
		return fmt.Errorf("error parsing storage URL: %w", err)
	}

	// Get bucket handle and delete the object
	bucket := c.storageClient.Bucket(bucketName)
	obj := bucket.Object(objectPath)

	if err := obj.Delete(ctx); err != nil {
		// Check if object doesn't exist (already deleted)
		if err == storage.ErrObjectNotExist {
			log.Info(logger.LogFirebaseStorageAlreadyDeleted, "path", objectPath)
			return nil
		}
		log.Error(logger.LogFirebaseStorageDeleteError, "path", objectPath, "error", err)
		return fmt.Errorf("error deleting storage file: %w", err)
	}

	log.Info(logger.LogFirebaseStorageDeleteOK, "path", objectPath)
	return nil
}

// parseFirebaseStorageURL extracts bucket name and object path from Firebase Storage URL
// Input: https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/motorcycles%2Fprofile.jpg?alt=media
// Output: bucketName="motogo-xxx.appspot.com", objectPath="motorcycles/profile.jpg"
func parseFirebaseStorageURL(fileURL string) (bucketName, objectPath string, err error) {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	// Expected path: /v0/b/{bucket}/o/{encoded-path}
	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 5 || parts[1] != "v0" || parts[2] != "b" || parts[4] != "o" {
		return "", "", fmt.Errorf("unexpected Firebase Storage URL format")
	}

	bucketName = parts[3]

	// The object path is everything after /o/, URL-encoded
	objectPathEncoded := strings.Join(parts[5:], "/")
	objectPath, err = url.PathUnescape(objectPathEncoded)
	if err != nil {
		return "", "", fmt.Errorf("error decoding object path: %w", err)
	}

	return bucketName, objectPath, nil
}
