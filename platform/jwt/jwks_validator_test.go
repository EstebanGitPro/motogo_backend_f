package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/MicahParks/keyfunc"
	jwtlib "github.com/golang-jwt/jwt/v4"
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

// ============================================
// ValidateToken Tests
// ============================================

func createTestValidator(t *testing.T, issuer string) (*JWKSValidator, *rsa.PrivateKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	givenKey := keyfunc.NewGivenRSACustomWithOptions(&privateKey.PublicKey, keyfunc.GivenKeyOptions{Algorithm: "RS256"})
	jwks := keyfunc.NewGiven(map[string]keyfunc.GivenKey{"test-kid": givenKey})

	return &JWKSValidator{
		jwks:           jwks,
		expectedIssuer: issuer,
	}, privateKey
}

func signToken(t *testing.T, claims jwtlib.MapClaims, key *rsa.PrivateKey) string {
	t.Helper()
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	token.Header["kid"] = "test-kid"
	signed, err := token.SignedString(key)
	assert.NoError(t, err)
	return signed
}

func TestValidateToken_MalformedToken(t *testing.T) {
	v, _ := createTestValidator(t, "")
	_, err := v.ValidateToken("not-a-valid-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTokenMalformed)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	v, key := createTestValidator(t, "")
	tokenStr := signToken(t, jwtlib.MapClaims{
		"exp": time.Now().Add(-time.Hour).Unix(),
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
		"sub": "user-1",
	}, key)

	_, err := v.ValidateToken(tokenStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestValidateToken_InvalidIssuer(t *testing.T) {
	v, key := createTestValidator(t, "https://expected.com/realms/test")
	tokenStr := signToken(t, jwtlib.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"iss": "https://wrong.com/realms/other",
		"sub": "user-1",
	}, key)

	_, err := v.ValidateToken(tokenStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidIssuer)
}

func TestValidateToken_ValidToken_WithIssuer(t *testing.T) {
	v, key := createTestValidator(t, "https://example.com/realms/test")
	tokenStr := signToken(t, jwtlib.MapClaims{
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "https://example.com/realms/test",
		"sub":   "user-1",
		"email": "test@example.com",
	}, key)

	claims, err := v.ValidateToken(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-1", claims["sub"])
	assert.Equal(t, "test@example.com", claims["email"])
	assert.Equal(t, "https://example.com/realms/test", claims["iss"])
}

func TestValidateToken_ValidToken_NoIssuerCheck(t *testing.T) {
	v, key := createTestValidator(t, "") // No issuer validation
	tokenStr := signToken(t, jwtlib.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"sub": "user-1",
	}, key)

	claims, err := v.ValidateToken(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-1", claims["sub"])
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	v, _ := createTestValidator(t, "")
	// Sign with a different key
	otherKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenStr := signToken(t, jwtlib.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"sub": "user-1",
	}, otherKey)

	_, err := v.ValidateToken(tokenStr)
	assert.Error(t, err)
}

// ============================================
// Close Tests
// ============================================

func TestClose_NilJWKS(t *testing.T) {
	v := &JWKSValidator{jwks: nil}
	// Should not panic
	v.Close()
}

func TestClose_WithJWKS(t *testing.T) {
	givenKey := keyfunc.NewGivenRSACustomWithOptions(&rsa.PublicKey{}, keyfunc.GivenKeyOptions{Algorithm: "RS256"})
	jwks := keyfunc.NewGiven(map[string]keyfunc.GivenKey{"test-kid": givenKey})
	v := &JWKSValidator{jwks: jwks}
	// Should not panic
	v.Close()
}
