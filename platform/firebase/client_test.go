package firebase

import (
	"context"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Interface Compliance Tests
// ============================================

func TestFirebaseClient_InterfaceCompliance(t *testing.T) {
	// Verify Client implements output.CustomTokenProvider interface
	var _ output.CustomTokenProvider = (*Client)(nil)
}

// ============================================
// NewClient Tests
// ============================================

func TestNewClient_InvalidCredentialsPath_ReturnsError(t *testing.T) {
	client, err := NewClient("/nonexistent/path/credentials.json")

	assert.Error(t, err)
	assert.Nil(t, client)
	// Firebase returns "cannot read credentials file" for invalid path
	assert.Contains(t, err.Error(), "credentials")
}

func TestNewClient_EmptyCredentialsPath_ReturnsError(t *testing.T) {
	client, err := NewClient("")

	assert.Error(t, err)
	assert.Nil(t, client)
}

// ============================================
// parseFirebaseStorageURL Tests
// ============================================

func TestParseFirebaseStorageURL_ValidURL_Success(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/motorcycles%2Fprofile.jpg?alt=media"

	bucketName, objectPath, err := parseFirebaseStorageURL(url)

	assert.NoError(t, err)
	assert.Equal(t, "motogo-xxx.appspot.com", bucketName)
	assert.Equal(t, "motorcycles/profile.jpg", objectPath)
}

func TestParseFirebaseStorageURL_NestedPath_Success(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/project.appspot.com/o/users%2Fuser-123%2Fmotorcycles%2Fmoto-456%2Fprofile.png?alt=media&token=abc"

	bucketName, objectPath, err := parseFirebaseStorageURL(url)

	assert.NoError(t, err)
	assert.Equal(t, "project.appspot.com", bucketName)
	assert.Equal(t, "users/user-123/motorcycles/moto-456/profile.png", objectPath)
}

func TestParseFirebaseStorageURL_InvalidURL_Error(t *testing.T) {
	_, _, err := parseFirebaseStorageURL("not-a-valid-url")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected Firebase Storage URL format")
}

func TestParseFirebaseStorageURL_WrongFormat_Error(t *testing.T) {
	url := "https://storage.googleapis.com/bucket/object.jpg"

	_, _, err := parseFirebaseStorageURL(url)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected Firebase Storage URL format")
}

func TestParseFirebaseStorageURL_EmptyURL_Error(t *testing.T) {
	_, _, err := parseFirebaseStorageURL("")

	assert.Error(t, err)
}

// ============================================
// extractBucketName Tests
// ============================================

func TestExtractBucketName_ReturnsEmptyString(t *testing.T) {
	// Current implementation returns empty string
	result := extractBucketName("/path/to/credentials.json")
	assert.Equal(t, "", result)
}

// ============================================
// DeleteStorageFile Edge Case Tests
// ============================================

func TestDeleteStorageFile_NilStorageClient_ReturnsNil(t *testing.T) {
	c := &Client{
		storageClient: nil, // Not configured
	}

	err := c.DeleteStorageFile(context.Background(), "https://firebasestorage.googleapis.com/v0/b/bucket/o/file.jpg")

	assert.NoError(t, err)
}

func TestDeleteStorageFile_EmptyURL_ReturnsNil(t *testing.T) {
	c := &Client{
		storageClient: nil,
	}

	err := c.DeleteStorageFile(context.Background(), "")

	assert.NoError(t, err)
}
