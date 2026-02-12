package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
)

// ============================================
// JWKSConfig Tests
// ============================================

func TestJWKSConfig_Defaults(t *testing.T) {
	config := JWKSConfig{
		JWKSURL: "https://example.com/jwks",
		Issuer:  "https://example.com/realms/test",
	}

	assert.Equal(t, "https://example.com/jwks", config.JWKSURL)
	assert.Equal(t, "https://example.com/realms/test", config.Issuer)
	assert.Equal(t, time.Duration(0), config.RefreshInterval) // Zero before initialization
}

// ============================================
// NewJWKSValidator Tests
// ============================================

func TestNewJWKSValidator_EmptyURL_ReturnsError(t *testing.T) {
	config := JWKSConfig{
		JWKSURL: "",
		Issuer:  "https://example.com/realms/test",
	}

	validator, err := NewJWKSValidator(context.Background(), config)

	assert.Error(t, err)
	assert.Nil(t, validator)
	assert.Contains(t, err.Error(), "JWKS URL cannot be empty")
}

func TestNewJWKSValidator_InvalidURL_ReturnsError(t *testing.T) {
	config := JWKSConfig{
		JWKSURL:         "http://localhost:99999/invalid-jwks", // Invalid port/unreachable
		Issuer:          "https://example.com/realms/test",
		RefreshInterval: time.Millisecond * 100,
	}

	// This should fail because it tries to fetch from an unreachable URL
	validator, err := NewJWKSValidator(context.Background(), config)

	assert.Error(t, err)
	assert.Nil(t, validator)
	assert.ErrorIs(t, err, ErrJWKSUnavailable)
}

// ============================================
// Custom Error Tests
// ============================================

func TestCustomErrors_AreDistinct(t *testing.T) {
	errors := []error{
		ErrTokenExpired,
		ErrTokenNotValidYet,
		ErrInvalidSignature,
		ErrInvalidIssuer,
		ErrInvalidClaims,
		ErrJWKSUnavailable,
		ErrTokenMalformed,
	}

	// Verify all errors are distinct
	seen := make(map[string]bool)
	for _, err := range errors {
		msg := err.Error()
		assert.False(t, seen[msg], "Duplicate error message: %s", msg)
		seen[msg] = true
	}
}

func TestErrTokenExpired_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token has expired", ErrTokenExpired.Error())
}

func TestErrTokenNotValidYet_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token is not valid yet", ErrTokenNotValidYet.Error())
}

func TestErrInvalidSignature_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token signature is invalid", ErrInvalidSignature.Error())
}

func TestErrInvalidIssuer_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token issuer is invalid", ErrInvalidIssuer.Error())
}

func TestErrInvalidClaims_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token claims are invalid", ErrInvalidClaims.Error())
}

func TestErrJWKSUnavailable_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "JWKS endpoint is unavailable", ErrJWKSUnavailable.Error())
}

func TestErrTokenMalformed_HasCorrectMessage(t *testing.T) {
	assert.Equal(t, "token is malformed", ErrTokenMalformed.Error())
}

// ============================================
// Interface Implementation Tests
// ============================================

func TestValidator_InterfaceCompliance(t *testing.T) {
	// Verify JWKSValidator implements output.JWTValidator interface
	var _ output.JWTValidator = (*JWKSValidator)(nil)
}
