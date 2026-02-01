package firebase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// Interface Compliance Tests
// ============================================

func TestFirebaseClient_InterfaceCompliance(t *testing.T) {
	// Verify Client implements FirebaseClient interface
	var _ FirebaseClient = (*Client)(nil)
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
